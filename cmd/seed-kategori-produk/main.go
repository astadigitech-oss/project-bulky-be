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
		NamaID  string
		NamaEN  string
		Slug    string
		IconURL string
	}{
		{NamaID: "Elektronik", NamaEN: "Electronics", Slug: "elektronik", IconURL: "product-categories/elektronik.png"},
		{NamaID: "Ibu & Anak", NamaEN: "Mother & Baby", Slug: "ibu-anak", IconURL: "product-categories/ibu-anak.png"},
		{NamaID: "Kosmetik", NamaEN: "Cosmetics", Slug: "kosmetik", IconURL: "product-categories/kosmetik.png"},
		{NamaID: "Otomotif", NamaEN: "Automotive", Slug: "otomotif", IconURL: "product-categories/otomotif.png"},
		{NamaID: "Alat Rumah Tangga", NamaEN: "Household Appliances", Slug: "alat-rumah-tangga", IconURL: "product-categories/alat-rumah-tangga.png"},
		{NamaID: "FMCG", NamaEN: "FMCG", Slug: "fmcg", IconURL: "product-categories/fmcg.png"},
		{NamaID: "Tools", NamaEN: "Tools", Slug: "tools", IconURL: "product-categories/tools.png"},
		{NamaID: "Redknot", NamaEN: "Redknot", Slug: "redknot", IconURL: "product-categories/redknot.png"},
		{NamaID: "Sepatu", NamaEN: "Shoes", Slug: "sepatu", IconURL: "product-categories/sepatu.png"},
		{NamaID: "Aksesoris", NamaEN: "Accessories", Slug: "aksesoris", IconURL: "product-categories/aksesoris.png"},
		{NamaID: "Buku", NamaEN: "Books", Slug: "buku", IconURL: "product-categories/buku.png"},
		{NamaID: "Tas", NamaEN: "Bags", Slug: "tas", IconURL: "product-categories/tas.png"},
		{NamaID: "Fashion", NamaEN: "Fashion", Slug: "fashion", IconURL: "product-categories/fashion.png"},
		{NamaID: "Fashion & Tas", NamaEN: "Fashion & Bags", Slug: "fashion-tas", IconURL: "product-categories/fashion-tas.png"},
		{NamaID: "Fashion & Aksesoris", NamaEN: "Fashion & Accessories", Slug: "fashion-aksesoris", IconURL: "product-categories/fashion-aksesoris.png"},
		{NamaID: "Kulkas", NamaEN: "Refrigerator", Slug: "kulkas", IconURL: "product-categories/kulkas.png"},
		{NamaID: "Mesin Cuci", NamaEN: "Washing Machine", Slug: "mesin-cuci", IconURL: "product-categories/mesin-cuci.png"},
		{NamaID: "TV", NamaEN: "TV", Slug: "tv", IconURL: "product-categories/tv.png"},
		{NamaID: "Lainnya", NamaEN: "Others", Slug: "lainnya", IconURL: "product-categories/lainnya.png"},
		{NamaID: "Unggulan", NamaEN: "Featured", Slug: "unggulan", IconURL: "product-categories/unggulan.png"},
		{NamaID: "Toys", NamaEN: "Toys", Slug: "toys", IconURL: "product-categories/toys.png"},
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
			fmt.Printf("⊘ Skip: %s (slug sudah terdaftar)\n", kategoriData.NamaID)
			skipCount++
			continue
		}

		// Create kategori produk
		namaEN := kategoriData.NamaEN
		kategori := &models.KategoriProduk{
			NamaID:    kategoriData.NamaID,
			NamaEN:    &namaEN,
			Slug:      kategoriData.Slug,
			IconURL:   &kategoriData.IconURL,
			IsActive:  true,
			CreatedAt: time.Now().Add(time.Duration(i) * time.Second), // Ordering by created_at
		}

		if err := db.Create(kategori).Error; err != nil {
			log.Printf("✗ Gagal membuat kategori %s: %v\n", kategoriData.NamaID, err)
			continue
		}

		fmt.Printf("✓ Kategori berhasil dibuat: %s (%s)\n", kategori.NamaID, kategori.Slug)
		successCount++
	}

	fmt.Println()
	fmt.Println("=== Seeding Summary ===")
	fmt.Printf("✓ Berhasil: %d kategori\n", successCount)
	fmt.Printf("⊘ Dilewati: %d kategori (sudah ada)\n", skipCount)
	fmt.Printf("✗ Gagal: %d kategori\n", len(kategoris)-successCount-skipCount)
	fmt.Printf("📊 Total: %d kategori produk\n", len(kategoris))
}
