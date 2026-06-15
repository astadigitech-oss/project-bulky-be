package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/pkg/database"
	"project-bulky-be/pkg/utils"

	"github.com/joho/godotenv"
	"gorm.io/gorm/clause"
)

type cityItem struct {
	ItemID   int    `json:"item_id"`
	ItemName string `json:"item_name"`
}

type cityListResponse struct {
	Data []cityItem `json:"data"`
}

type subdistrictItem struct {
	ItemID   int    `json:"item_id"`
	ItemName string `json:"item_name"`
}

type subdistrictListResponse struct {
	Data []subdistrictItem `json:"data"`
}

// dirOfThisFile mengembalikan path direktori file Go ini saat compile time.
func dirOfThisFile() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

func readCityListFromFile(path string) ([]cityItem, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baca file city list gagal: %w", err)
	}
	var clr cityListResponse
	if err := json.Unmarshal(data, &clr); err != nil {
		return nil, fmt.Errorf("parse city list gagal: %w", err)
	}
	return clr.Data, nil
}

func readSubdistrictListFromFile(path string) ([]subdistrictItem, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baca file subdistrict list gagal: %w", err)
	}
	var slr subdistrictListResponse
	if err := json.Unmarshal(data, &slr); err != nil {
		return nil, fmt.Errorf("parse subdistrict list gagal: %w", err)
	}
	return slr.Data, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()
	database.InitDB(cfg)
	db := database.GetDB()

	baseDir := dirOfThisFile()
	cityFile := filepath.Join(baseDir, "forwarder_city_list.json")
	subdistrictFile := filepath.Join(baseDir, "forwarder_subdistrict_list.json")

	// --- Seed forwarder_city_mapping ---
	fmt.Println("=== Seed Forwarder City Mapping ===")

	cities, err := readCityListFromFile(cityFile)
	if err != nil {
		log.Fatalf("Gagal membaca city list: %v", err)
	}

	seeded := 0
	for _, city := range cities {
		pattern := utils.NormalizeKota(city.ItemName)
		record := models.ForwarderCityMapping{
			KotaPattern:       pattern,
			ForwarderCityID:   city.ItemID,
			ForwarderCityName: city.ItemName,
		}
		result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "kota_pattern"}},
			DoUpdates: clause.AssignmentColumns([]string{"forwarder_city_id", "forwarder_city_name", "updated_at"}),
		}).Create(&record)
		if result.Error != nil {
			log.Printf("Gagal upsert kota '%s': %v", city.ItemName, result.Error)
			continue
		}
		seeded++
	}

	fmt.Printf("Selesai: %d/%d kota berhasil di-seed.\n\n", seeded, len(cities))

	// --- Seed forwarder_subdistrict_mapping ---
	withSubdistricts := false
	for _, arg := range os.Args[1:] {
		if arg == "--with-subdistricts" {
			withSubdistricts = true
		}
	}

	if !withSubdistricts {
		fmt.Println("Tip: jalankan dengan flag --with-subdistricts untuk seed tabel forwarder_subdistrict_mapping juga.")
		return
	}

	fmt.Println("=== Seed Forwarder Subdistrict Mapping ===")

	subdistricts, err := readSubdistrictListFromFile(subdistrictFile)
	if err != nil {
		log.Fatalf("Gagal membaca subdistrict list: %v", err)
	}

	// File JSON global tidak menyertakan city_id per kecamatan.
	// Gunakan forwarder_city_id = 0 sebagai penanda "global/semua kota".
	totalSubSeeded := 0
	for _, sub := range subdistricts {
		pattern := utils.NormalizeKecamatan(sub.ItemName)
		record := models.ForwarderSubdistrictMapping{
			KecamatanPattern:         pattern,
			ForwarderCityID:          0,
			ForwarderSubdistrictID:   sub.ItemID,
			ForwarderSubdistrictName: sub.ItemName,
		}
		result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "kecamatan_pattern"}, {Name: "forwarder_city_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"forwarder_subdistrict_id", "forwarder_subdistrict_name", "updated_at"}),
		}).Create(&record)
		if result.Error != nil {
			log.Printf("Gagal upsert kecamatan '%s': %v", sub.ItemName, result.Error)
			continue
		}
		totalSubSeeded++
	}

	fmt.Printf("Selesai: %d/%d kecamatan berhasil di-seed.\n", totalSubSeeded, len(subdistricts))
}
