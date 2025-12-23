package database

import (
	"fmt"
	"log"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

func GetDB() *gorm.DB {
	return DB
}

// AutoMigrate - Only use in development!
// For production, use manual migration via golang-migrate
func AutoMigrate(env string) {
	if env == "production" {
		log.Println("Skipping auto migrate in production. Use manual migration.")
		return
	}

	// Enable UUID extension
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	err := DB.AutoMigrate(
		&models.KategoriProduk{},
		&models.MerekProduk{},
		&models.KondisiProduk{},
		&models.KondisiPaket{},
		&models.SumberProduk{},
		&models.Warehouse{},
		&models.TipeProduk{},
		&models.DiskonKategori{},
		&models.BannerTipeProduk{},
		&models.Produk{},
		&models.ProdukGambar{},
		&models.ProdukDokumen{},
	)
	if err != nil {
		log.Fatal("Failed to auto migrate:", err)
	}
	log.Println("Database migrated successfully")
}
