package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/pkg/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// App menampung koneksi, mode eksekusi, report, dan state mapping antar-fase.
type App struct {
	my      *sql.DB  // sumber: MySQL v1 (hasil restore dump)
	pg      *gorm.DB // target: Postgres v2
	tx      *gorm.DB // transaksi aktif per fase (nil saat dry-run)
	execute bool
	rep     *Report

	// state target (di-preload sebelum fase berjalan)
	tgt *TargetState

	// mapping v1 -> v2 hasil Fase 1
	kategoriKnown map[string]bool   // id kategori yang valid di v2 (existing + inserted)
	merekKnown    map[string]bool
	condMap       map[string]string // product_conditions.id -> kondisi_produk.id
	pakMap        map[string]string // status_packages.id -> kondisi_paket.id
	sumMap        map[string]string // product_statuses.id -> sumber_produk.id
	whMap         map[string]string // warehouses.id -> warehouse.id
	tipeMap       map[string]string // packaging_type enum -> tipe_produk.id
	fallbackKategoriID string
	fallbackWarehouseID string

	// state hasil Fase 2/3 (untuk fase berikutnya)
	produkKnown map[string]bool
	buyerKnown  map[string]bool
}

const (
	fallbackKondisiID      = "9fffffff-0000-4000-a000-000000000001"
	fallbackKondisiPaketID = "9fffffff-0000-4000-a000-000000000002"
	// ID v1 kategori "Lainnya" (slug v1: uncategorized) — dipakai sebagai fallback kategori
	v1KategoriLainnyaID = "9cc2e714-28ca-4e45-a09a-17aedea3c95e"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	execute := flag.Bool("execute", false, "tulis ke Postgres (tanpa flag ini = dry-run)")
	reportPath := flag.String("report", "migrate-v1-report.json", "path file report JSON")
	flag.Parse()

	app := &App{
		execute:     *execute,
		rep:         NewReport(!*execute),
		produkKnown: map[string]bool{},
		buyerKnown:  map[string]bool{},
	}

	mode := "DRY-RUN (tidak menulis apa pun; jalankan dengan -execute untuk eksekusi)"
	if app.execute {
		mode = "EXECUTE"
	}
	log.Printf("=== Migrasi Data Bulky v1 -> v2 (produk & buyer) | mode: %s ===", mode)

	// Koneksi target Postgres (pakai config/env yang sama dengan aplikasi)
	cfg := config.LoadConfig()
	database.InitDB(cfg)
	app.pg = database.GetDB()

	// Koneksi sumber MySQL v1
	myDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC&charset=utf8mb4",
		envOr("V1_DB_USER", "root"),
		os.Getenv("V1_DB_PASSWORD"),
		envOr("V1_DB_HOST", "127.0.0.1"),
		envOr("V1_DB_PORT", "3306"),
		envOr("V1_DB_NAME", "bulky_v1"),
	)
	var err error
	app.my, err = sql.Open("mysql", myDSN)
	if err != nil {
		log.Fatalf("gagal membuka koneksi MySQL v1: %v", err)
	}
	if err := app.my.Ping(); err != nil {
		log.Fatalf("gagal konek MySQL v1 (%s@%s/%s): %v — pastikan dump sudah di-restore",
			envOr("V1_DB_USER", "root"), envOr("V1_DB_HOST", "127.0.0.1"), envOr("V1_DB_NAME", "bulky_v1"), err)
	}
	defer app.my.Close()

	app.tgt = LoadTargetState(app.pg)

	phases := []struct {
		name string
		fn   func() error
	}{
		{"Fase 1: master referensi", app.phaseMaster},
		{"Fase 2: produk + gambar + dokumen + merek pivot", app.phaseProduk},
		{"Fase 3: buyer + admin", app.phaseBuyer},
		{"Fase 4: alamat buyer", app.phaseAlamat},
	}
	for _, p := range phases {
		log.Printf("--- %s ---", p.name)
		app.beginPhase()
		if err := p.fn(); err != nil {
			app.rollbackPhase()
			log.Fatalf("%s GAGAL: %v (transaksi fase ini di-rollback)", p.name, err)
		}
		app.commitPhase()
	}

	log.Print("--- Fase 6: validasi ---")
	app.validate()

	app.rep.FinishedAt = time.Now().UTC()
	if err := app.rep.WriteFile(*reportPath); err != nil {
		log.Fatalf("gagal menulis report: %v", err)
	}
	app.rep.PrintSummary()
	log.Printf("Report lengkap: %s", *reportPath)
	if !app.execute {
		log.Print("Selesai DRY-RUN. Tinjau report, lalu jalankan ulang dengan -execute.")
	} else {
		log.Print("Selesai. Jangan lupa Fase 5 (file fisik) setelah file export prod tersedia.")
	}
}

// beginPhase membuka transaksi Postgres per fase (hanya saat execute).
func (a *App) beginPhase() {
	if a.execute {
		a.tx = a.pg.Begin()
		if a.tx.Error != nil {
			log.Fatalf("gagal membuka transaksi: %v", a.tx.Error)
		}
	}
}

func (a *App) commitPhase() {
	if a.execute && a.tx != nil {
		if err := a.tx.Commit().Error; err != nil {
			log.Fatalf("gagal commit transaksi: %v", err)
		}
		a.tx = nil
	}
}

func (a *App) rollbackPhase() {
	if a.execute && a.tx != nil {
		a.tx.Rollback()
		a.tx = nil
	}
}

// exec menjalankan statement tulis; di mode dry-run hanya menghitung.
func (a *App) exec(query string, args ...interface{}) error {
	if !a.execute {
		return nil
	}
	return a.tx.Exec(query, args...).Error
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
