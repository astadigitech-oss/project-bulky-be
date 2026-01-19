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

type MetodePembayaranSeed struct {
	GroupName string // Reference to group
	Nama      string
	Kode      string // Xendit channel code
	LogoValue string
	Urutan    int
}

var metodePembayaranSeeds = []MetodePembayaranSeed{
	// ============================================
	// BANK TRANSFER / VA
	// ============================================
	{
		GroupName: "Bank Transfer / VA",
		Nama:      "BCA",
		Kode:      "BCA",
		LogoValue: "bca",
		Urutan:    1,
	},
	{
		GroupName: "Bank Transfer / VA",
		Nama:      "MANDIRI",
		Kode:      "MANDIRI",
		LogoValue: "mandiri",
		Urutan:    2,
	},
	{
		GroupName: "Bank Transfer / VA",
		Nama:      "BNI",
		Kode:      "BNI",
		LogoValue: "bni",
		Urutan:    3,
	},
	{
		GroupName: "Bank Transfer / VA",
		Nama:      "BRI",
		Kode:      "BRI",
		LogoValue: "bri",
		Urutan:    4,
	},
	{
		GroupName: "Bank Transfer / VA",
		Nama:      "PERMATA",
		Kode:      "PERMATA",
		LogoValue: "permata",
		Urutan:    5,
	},
	{
		GroupName: "Bank Transfer / VA",
		Nama:      "CIMB NIAGA",
		Kode:      "CIMB",
		LogoValue: "cimb",
		Urutan:    6,
	},
	{
		GroupName: "Bank Transfer / VA",
		Nama:      "BSI",
		Kode:      "BSI",
		LogoValue: "bsi",
		Urutan:    7,
	},
	{
		GroupName: "Bank Transfer / VA",
		Nama:      "BJB",
		Kode:      "BJB",
		LogoValue: "bjb",
		Urutan:    8,
	},

	// ============================================
	// E-WALLET
	// ============================================
	{
		GroupName: "E-Wallet",
		Nama:      "LINK AJA",
		Kode:      "LINKAJA",
		LogoValue: "linkaja",
		Urutan:    1,
	},
	{
		GroupName: "E-Wallet",
		Nama:      "DANA",
		Kode:      "DANA",
		LogoValue: "dana",
		Urutan:    2,
	},
	{
		GroupName: "E-Wallet",
		Nama:      "OVO",
		Kode:      "ID_OVO",
		LogoValue: "ovo",
		Urutan:    3,
	},
	{
		GroupName: "E-Wallet",
		Nama:      "ShopeePay",
		Kode:      "ID_SHOPEEPAY",
		LogoValue: "shopeepay",
		Urutan:    4,
	},

	// ============================================
	// KARTU KREDIT
	// ============================================
	{
		GroupName: "Kartu Kredit",
		Nama:      "KARTU KREDIT",
		Kode:      "CREDIT_CARD",
		LogoValue: "credit-card",
		Urutan:    1,
	},

	// ============================================
	// QRIS
	// ============================================
	{
		GroupName: "QRIS",
		Nama:      "QRIS",
		Kode:      "ID_QRIS",
		LogoValue: "qris",
		Urutan:    1,
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

	log.Println("🌱 Seeding Metode Pembayaran...")

	successCount := 0
	skipCount := 0
	errorCount := 0

	for _, seed := range metodePembayaranSeeds {
		// Get group ID
		var group models.MetodePembayaranGroup
		if err := db.Where("nama = ?", seed.GroupName).First(&group).Error; err != nil {
			log.Printf("❌ Error: Group '%s' not found for %s", seed.GroupName, seed.Nama)
			errorCount++
			continue
		}

		// Check if exists
		var existing models.MetodePembayaran
		err := db.Where("kode = ?", seed.Kode).First(&existing).Error

		if err == nil {
			// Already exists, update if needed
			needUpdate := false
			if existing.Nama != seed.Nama {
				existing.Nama = seed.Nama
				needUpdate = true
			}
			if existing.GroupID != group.ID {
				existing.GroupID = group.ID
				needUpdate = true
			}
			if existing.Urutan != seed.Urutan {
				existing.Urutan = seed.Urutan
				needUpdate = true
			}
			if existing.LogoValue == nil || *existing.LogoValue != seed.LogoValue {
				existing.LogoValue = &seed.LogoValue
				needUpdate = true
			}

			if needUpdate {
				if err := db.Save(&existing).Error; err != nil {
					log.Printf("❌ Failed to update %s: %v", seed.Nama, err)
					errorCount++
					continue
				}
				log.Printf("✓ Updated: %s (Group: %s)", seed.Nama, seed.GroupName)
				successCount++
			} else {
				log.Printf("⊘ Skip: %s (already exists)", seed.Nama)
				skipCount++
			}
			continue
		}

		// Create new
		metode := models.MetodePembayaran{
			ID:        uuid.New(),
			GroupID:   group.ID,
			Nama:      seed.Nama,
			Kode:      seed.Kode,
			LogoValue: &seed.LogoValue,
			Urutan:    seed.Urutan,
			IsActive:  true,
		}

		if err := db.Create(&metode).Error; err != nil {
			log.Printf("❌ Failed to create %s: %v", seed.Nama, err)
			errorCount++
			continue
		}

		log.Printf("✓ Created: %s (Group: %s)", seed.Nama, seed.GroupName)
		successCount++
	}

	fmt.Println()
	log.Println("=== Seeding Summary ===")
	log.Printf("✓ Created/Updated: %d payment methods", successCount)
	log.Printf("⊘ Skipped: %d payment methods (already exists)", skipCount)
	log.Printf("❌ Errors: %d payment methods", errorCount)
	log.Printf("📊 Total: %d payment methods", len(metodePembayaranSeeds))

	if errorCount > 0 {
		fmt.Println()
		log.Println("⚠️  Some payment methods failed to seed. Please check the errors above.")
		log.Println("💡 Make sure to run 'seed-metode-pembayaran-group' first!")
		os.Exit(1)
	}

	fmt.Println()
	log.Println("✅ Metode Pembayaran seeding completed!")
}
