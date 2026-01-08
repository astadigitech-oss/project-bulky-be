package main

import (
	"fmt"
	"log"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDB(cfg)
	db := database.GetDB()

	fmt.Println("=== Seeding Kategori Produk ===")
	fmt.Println()

	// Data kategori produk untuk seed
	kategoris := []struct {
		Nama    string
		Slug    string
		IconURL string
	}{
		{Nama: "Elektronik", Slug: "elektronik", IconURL: "product-categories/elektronik.png"},
		{Nama: "Ibu & Anak", Slug: "ibu-anak", IconURL: "product-categories/ibu-anak.png"},
		{Nama: "Kosmetik", Slug: "kosmetik", IconURL: "product-categories/kosmetik.png"},
		{Nama: "Otomotif", Slug: "otomotif", IconURL: "product-categories/otomotif.png"},
		{Nama: "Alat Rumah Tangga", Slug: "alat-rumah-tangga", IconURL: "product-categories/alat-rumah-tangga.png"},
		{Nama: "FMCG", Slug: "fmcg", IconURL: "product-categories/fmcg.png"},
		{Nama: "Tools", Slug: "tools", IconURL: "product-categories/tools.png"},
		{Nama: "Redknot", Slug: "redknot", IconURL: "product-categories/redknot.png"},
		{Nama: "Sepatu", Slug: "sepatu", IconURL: "product-categories/sepatu.png"},
		{Nama: "Aksesoris", Slug: "aksesoris", IconURL: "product-categories/aksesoris.png"},
		{Nama: "Buku", Slug: "buku", IconURL: "product-categories/buku.png"},
		{Nama: "Tas", Slug: "tas", IconURL: "product-categories/tas.png"},
		{Nama: "Fashion", Slug: "fashion", IconURL: "product-categories/fashion.png"},
		{Nama: "Fashion & Tas", Slug: "fashion-tas", IconURL: "product-categories/fashion-tas.png"},
		{Nama: "Fashion & Aksesoris", Slug: "fashion-aksesoris", IconURL: "product-categories/fashion-aksesoris.png"},
		{Nama: "Kulkas", Slug: "kulkas", IconURL: "product-categories/kulkas.png"},
		{Nama: "Mesin Cuci", Slug: "mesin-cuci", IconURL: "product-categories/mesin-cuci.png"},
		{Nama: "TV", Slug: "tv", IconURL: "product-categories/tv.png"},
		{Nama: "Lainnya", Slug: "lainnya", IconURL: "product-categories/lainnya.png"},
		{Nama: "Unggulan", Slug: "unggulan", IconURL: "product-categories/unggulan.png"},
		{Nama: "Toys", Slug: "toys", IconURL: "product-categories/toys.png"},
	}

	successCount := 0
	skipCount := 0

	for i, kategoriData := range kategoris {
		// Check if slug already exists
		var count int64
		db.Model(&models.KategoriProduk{}).
			Where("slug = ?", kategoriData.Slug).
			Count(&count)

		if count > 0 {
			fmt.Printf("⊘ Skip: %s (slug sudah terdaftar)\n", kategoriData.Nama)
			skipCount++
			continue
		}

		// Create kategori produk
		kategori := &models.KategoriProduk{
			Nama:                    kategoriData.Nama,
			Slug:                    kategoriData.Slug,
			IconURL:                 &kategoriData.IconURL,
			MemilikiKondisiTambahan: false,
			IsActive:                true,
			CreatedAt:               time.Now().Add(time.Duration(i) * time.Second), // Ordering by created_at
		}

		if err := db.Create(kategori).Error; err != nil {
			log.Printf("✗ Gagal membuat kategori %s: %v\n", kategoriData.Nama, err)
			continue
		}

		fmt.Printf("✓ Kategori berhasil dibuat: %s (%s)\n", kategori.Nama, kategori.Slug)
		successCount++
	}

	fmt.Println()
	fmt.Println("=== Seeding Summary ===")
	fmt.Printf("✓ Berhasil: %d kategori\n", successCount)
	fmt.Printf("⊘ Dilewati: %d kategori (sudah ada)\n", skipCount)
	fmt.Printf("✗ Gagal: %d kategori\n", len(kategoris)-successCount-skipCount)
	fmt.Printf("📊 Total: %d kategori produk\n", len(kategoris))
}
