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

type WhatsAppHandlerSeed struct {
	NomorWA   string
	PesanAwal string
	IsActive  bool
}

var whatsappSeeds = []WhatsAppHandlerSeed{
	{
		NomorWA:   "62811833164", // Actual Bulky number from Doc 19
		PesanAwal: "Halo, ada yang bisa kami bantu?",
		IsActive:  true,
	},
	// Backup handler (inactive by default)
	{
		NomorWA:   "6289876543210", // Backup number
		PesanAwal: "Halo, tim Bulky siap membantu Anda!",
		IsActive:  false,
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

	log.Println("🌱 Seeding WhatsApp Handler...")

	successCount := 0
	skipCount := 0

	for _, seed := range whatsappSeeds {
		// Check if exists
		var existing models.WhatsAppHandler
		err := db.Where("nomor_wa = ?", seed.NomorWA).First(&existing).Error

		if err == nil {
			// Already exists, update if needed
			if existing.PesanAwal != seed.PesanAwal || existing.IsActive != seed.IsActive {
				existing.PesanAwal = seed.PesanAwal
				existing.IsActive = seed.IsActive
				if err := db.Save(&existing).Error; err != nil {
					log.Printf("❌ Failed to update %s: %v", seed.NomorWA, err)
					continue
				}
				log.Printf("✓ Updated: %s", seed.NomorWA)
			} else {
				log.Printf("⊘ Skip: %s (already exists)", seed.NomorWA)
				skipCount++
			}
			continue
		}

		// Create new
		handler := models.WhatsAppHandler{
			ID:        uuid.New(),
			NomorWA:   seed.NomorWA,
			PesanAwal: seed.PesanAwal,
			IsActive:  seed.IsActive,
		}

		if err := db.Create(&handler).Error; err != nil {
			log.Printf("❌ Failed to create %s: %v", seed.NomorWA, err)
			os.Exit(1)
		}

		log.Printf("✓ Created: %s", seed.NomorWA)
		successCount++
	}

	fmt.Println()
	log.Println("=== Seeding Summary ===")
	log.Printf("✓ Created/Updated: %d handlers", successCount)
	log.Printf("⊘ Skipped: %d handlers (already exists)", skipCount)
	fmt.Println()
	log.Println("✅ WhatsApp Handler seeding completed!")
}
