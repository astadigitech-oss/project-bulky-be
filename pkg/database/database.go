package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	// Configure logger based on environment
	var logLevel logger.LogLevel
	if cfg.AppEnv == "development" {
		logLevel = logger.Info // Log all queries in development
	} else {
		logLevel = logger.Warn // Log only slow queries and errors in production
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Register callback for slow query logging
	DB.Callback().Query().After("gorm:query").Register("slow_query_log", func(db *gorm.DB) {
		if db.Error != nil {
			return
		}

		elapsed := time.Since(db.Statement.Context.Value("query_start_time").(time.Time))
		if elapsed > 200*time.Millisecond {
			log.Printf("[SLOW QUERY] Duration: %v | SQL: %s", elapsed, db.Statement.SQL.String())
		}
	})

	// Set query start time before query execution
	DB.Callback().Query().Before("gorm:query").Register("query_start_time", func(db *gorm.DB) {
		db.Statement.Context = context.WithValue(db.Statement.Context, "query_start_time", time.Now())
	})

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
		&models.Admin{},
		&models.AdminSession{},
	)
	if err != nil {
		log.Fatal("Failed to auto migrate:", err)
	}
	log.Println("Database migrated successfully")
}
