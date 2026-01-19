package main

import (
	"fmt"
	"log"
	"os"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/pkg/database"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type MetodePembayaranGroupSeed struct {
	Nama   string
	Urutan int
}

var groupSeeds = []MetodePembayaranGroupSeed{
	{
		Nama:   "Bank Transfer / VA",
		Urutan: 1,
	},
	{
		Nama:   "E-Wallet",
		Urutan: 2,
	},
	{
		Nama:   "Kartu Kredit",
		Urutan: 3,
	},
	{
		Nama:   "QRIS",
		Urutan: 4,
	},
}

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

	log.Println("🌱 Seeding Metode Pembayaran Group...")

	successCount := 0
	skipCount := 0

	for _, seed := range groupSeeds {
		// Check if exists
		var existing models.MetodePembayaranGroup
		err := db.Where("nama = ?", seed.Nama).First(&existing).Error

		if err == nil {
			// Already exists, update urutan if needed
			if existing.Urutan != seed.Urutan {
				existing.Urutan = seed.Urutan
				if err := db.Save(&existing).Error; err != nil {
					log.Printf("❌ Failed to update %s: %v", seed.Nama, err)
					continue
				}
				log.Printf("✓ Updated: %s", seed.Nama)
				successCount++
			} else {
				log.Printf("⊘ Skip: %s (already exists)", seed.Nama)
				skipCount++
			}
			continue
		}

		// Create new
		group := models.MetodePembayaranGroup{
			ID:       uuid.New(),
			Nama:     seed.Nama,
			Urutan:   seed.Urutan,
			IsActive: true,
		}

		if err := db.Create(&group).Error; err != nil {
			log.Printf("❌ Failed to create %s: %v", seed.Nama, err)
			os.Exit(1)
		}

		log.Printf("✓ Created: %s", seed.Nama)
		successCount++
	}

	fmt.Println()
	log.Println("=== Seeding Summary ===")
	log.Printf("✓ Created/Updated: %d groups", successCount)
	log.Printf("⊘ Skipped: %d groups (already exists)", skipCount)
	fmt.Println()
	log.Println("✅ Metode Pembayaran Group seeding completed!")
}
