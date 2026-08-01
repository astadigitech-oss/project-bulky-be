package main

import (
	"log"
	"strings"

	"gorm.io/gorm"
)

// TargetState menampung snapshot isi Postgres v2 yang relevan,
// dipakai untuk mapping, dedup, dan idempotensi (skip baris yang sudah ada).
type TargetState struct {
	ProdukIDs    map[string]bool
	ProdukSlugs  map[string]bool
	ProdukSlugEN map[string]bool
	IDCargo      map[string]bool

	KategoriIDs     map[string]bool
	KategoriSlugs   map[string]bool
	KategoriSlugsEN map[string]bool
	KategoriBySlug  map[string]string

	MerekIDs     map[string]bool
	MerekSlugs   map[string]bool
	MerekSlugsEN map[string]bool

	KondisiIDs     map[string]bool
	KondisiBySlug  map[string]string
	KondisiSlugs   map[string]bool
	KondisiSlugsEN map[string]bool

	PaketIDs     map[string]bool
	PaketBySlug  map[string]string
	PaketSlugs   map[string]bool
	PaketSlugsEN map[string]bool

	SumberBySlug map[string]string
	TipeBySlug   map[string]string

	WarehouseBySlug map[string]string

	BuyerIDs       map[string]bool
	BuyerEmails    map[string]bool // lowercase
	BuyerTelepon   map[string]bool
	BuyerUsernames map[string]bool

	AdminIDs    map[string]bool
	AdminEmails map[string]bool // lowercase
	RoleAdminID string

	AlamatIDs map[string]bool

	DisclaimerIDs   map[string]bool
	DisclaimerSlugs map[string]bool

	PesananIDs     map[string]bool
	PesananKode    map[string]bool
	PesananItemIDs map[string]bool
	PembayaranIDs  map[string]bool

	KuponIDs      map[string]bool
	KuponUsageIDs map[string]bool
	UlasanIDs     map[string]bool
	ConsentIDs    map[string]bool

	KeranjangByBuyer map[string]string
}

func LoadTargetState(db *gorm.DB) *TargetState {
	t := &TargetState{}

	t.ProdukIDs = loadSet(db, `SELECT id::text FROM produk`)
	t.ProdukSlugs = loadSet(db, `SELECT slug FROM produk`)
	t.ProdukSlugEN = loadSet(db, `SELECT slug_en FROM produk WHERE slug_en IS NOT NULL`)
	t.IDCargo = loadSet(db, `SELECT id_cargo FROM produk WHERE id_cargo IS NOT NULL`)

	t.KategoriIDs = loadSet(db, `SELECT id::text FROM kategori_produk`)
	t.KategoriSlugs = loadSet(db, `SELECT slug FROM kategori_produk`)
	t.KategoriSlugsEN = loadSet(db, `SELECT slug_en FROM kategori_produk WHERE slug_en IS NOT NULL`)
	t.KategoriBySlug = loadMap(db, `SELECT slug, id::text FROM kategori_produk WHERE deleted_at IS NULL`)

	t.MerekIDs = loadSet(db, `SELECT id::text FROM merek_produk`)
	t.MerekSlugs = loadSet(db, `SELECT slug FROM merek_produk`)
	t.MerekSlugsEN = loadSet(db, `SELECT slug_en FROM merek_produk WHERE slug_en IS NOT NULL`)

	t.KondisiIDs = loadSet(db, `SELECT id::text FROM kondisi_produk`)
	t.KondisiBySlug = loadMap(db, `SELECT slug, id::text FROM kondisi_produk WHERE deleted_at IS NULL`)
	t.KondisiSlugs = loadSet(db, `SELECT slug FROM kondisi_produk`)
	t.KondisiSlugsEN = loadSet(db, `SELECT slug_en FROM kondisi_produk WHERE slug_en IS NOT NULL`)

	t.PaketIDs = loadSet(db, `SELECT id::text FROM kondisi_paket`)
	t.PaketBySlug = loadMap(db, `SELECT slug, id::text FROM kondisi_paket WHERE deleted_at IS NULL`)
	t.PaketSlugs = loadSet(db, `SELECT slug FROM kondisi_paket`)
	t.PaketSlugsEN = loadSet(db, `SELECT slug_en FROM kondisi_paket WHERE slug_en IS NOT NULL`)

	t.SumberBySlug = loadMap(db, `SELECT slug, id::text FROM sumber_produk WHERE deleted_at IS NULL`)
	t.TipeBySlug = loadMap(db, `SELECT slug, id::text FROM tipe_produk WHERE deleted_at IS NULL`)
	t.WarehouseBySlug = loadMap(db, `SELECT slug, id::text FROM warehouse WHERE deleted_at IS NULL`)

	t.BuyerIDs = loadSet(db, `SELECT id::text FROM buyer`)
	t.BuyerEmails = loadSetLower(db, `SELECT email FROM buyer WHERE email IS NOT NULL`)
	t.BuyerTelepon = loadSet(db, `SELECT telepon FROM buyer`)
	t.BuyerUsernames = loadSet(db, `SELECT username FROM buyer WHERE username IS NOT NULL`)

	t.AdminIDs = loadSet(db, `SELECT id::text FROM admin`)
	t.AdminEmails = loadSetLower(db, `SELECT email FROM admin`)

	var roleID []string
	if err := db.Raw(`SELECT id::text FROM role WHERE kode = 'ADMIN'`).Scan(&roleID).Error; err != nil || len(roleID) == 0 {
		log.Fatalf("role ADMIN tidak ditemukan di target (err=%v) — pastikan migrasi v2 sudah dijalankan penuh", err)
	}
	t.RoleAdminID = roleID[0]

	t.AlamatIDs = loadSet(db, `SELECT id::text FROM alamat_buyer`)

	t.DisclaimerIDs = loadSet(db, `SELECT id::text FROM disclaimer`)
	// slug, slug_id, slug_en berbagi satu ruang nama agar alokasi slug baru
	// tidak bentrok dengan kolom mana pun yang sudah terisi di target.
	t.DisclaimerSlugs = loadSet(db, `SELECT slug FROM disclaimer WHERE slug IS NOT NULL
		UNION SELECT slug_id FROM disclaimer WHERE slug_id IS NOT NULL
		UNION SELECT slug_en FROM disclaimer WHERE slug_en IS NOT NULL`)

	t.PesananIDs = loadSet(db, `SELECT id::text FROM pesanan`)
	t.PesananKode = loadSet(db, `SELECT kode FROM pesanan`)
	t.PesananItemIDs = loadSet(db, `SELECT id::text FROM pesanan_item`)
	t.PembayaranIDs = loadSet(db, `SELECT id::text FROM pesanan_pembayaran`)

	t.KuponIDs = loadSet(db, `SELECT id::text FROM kupon`)
	t.KuponUsageIDs = loadSet(db, `SELECT id::text FROM kupon_usage`)
	t.UlasanIDs = loadSet(db, `SELECT id::text FROM ulasan`)
	t.ConsentIDs = loadSet(db, `SELECT id::text FROM buyer_disclaimer_consent`)

	t.KeranjangByBuyer = loadMap(db, `SELECT buyer_id::text, id::text FROM keranjang`)

	log.Printf("target: %d produk, %d kategori, %d merek, %d kondisi, %d kondisi_paket, %d buyer, %d admin, %d alamat",
		len(t.ProdukIDs), len(t.KategoriIDs), len(t.MerekIDs), len(t.KondisiIDs), len(t.PaketIDs),
		len(t.BuyerIDs), len(t.AdminIDs), len(t.AlamatIDs))
	return t
}

func loadSet(db *gorm.DB, query string) map[string]bool {
	var vals []string
	if err := db.Raw(query).Scan(&vals).Error; err != nil {
		log.Fatalf("gagal preload target (%s): %v", query, err)
	}
	out := make(map[string]bool, len(vals))
	for _, v := range vals {
		out[v] = true
	}
	return out
}

func loadSetLower(db *gorm.DB, query string) map[string]bool {
	var vals []string
	if err := db.Raw(query).Scan(&vals).Error; err != nil {
		log.Fatalf("gagal preload target (%s): %v", query, err)
	}
	out := make(map[string]bool, len(vals))
	for _, v := range vals {
		out[strings.ToLower(v)] = true
	}
	return out
}

func loadMap(db *gorm.DB, query string) map[string]string {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		log.Fatalf("gagal preload target (%s): %v", query, err)
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			log.Fatalf("gagal scan preload (%s): %v", query, err)
		}
		out[k] = v
	}
	return out
}
