package main

import (
	"fmt"
	"log"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/pkg/database"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
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

	fmt.Println("=== Seeding Buyer Accounts ===")
	fmt.Println()

	// Data buyer untuk seed
	buyers := []struct {
		Nama     string
		Username string
		Email    string
		Password string
		Telepon  string
	}{
		{
			Nama:     "John Doe",
			Username: "johndoe",
			Email:    "john.doe@example.com",
			Password: "buyer123",
			Telepon:  "081234567890",
		},
		{
			Nama:     "Jane Smith",
			Username: "janesmith",
			Email:    "jane.smith@example.com",
			Password: "buyer123",
			Telepon:  "081234567891",
		},
		{
			Nama:     "Bob Wilson",
			Username: "bobwilson",
			Email:    "bob.wilson@example.com",
			Password: "buyer123",
			Telepon:  "081234567892",
		},
	}

	successCount := 0
	skipCount := 0

	for _, buyerData := range buyers {
		// Check if email or username already exists
		var count int64
		db.Model(&models.Buyer{}).
			Where("email = ? OR username = ?", buyerData.Email, buyerData.Username).
			Count(&count)

		if count > 0 {
			fmt.Printf("⊘ Skip: %s (email/username sudah terdaftar)\n", buyerData.Email)
			skipCount++
			continue
		}

		// Hash password using configured bcrypt cost
		hashedPassword, err := utils.HashPasswordWithCost(buyerData.Password, cfg.BcryptCost)
		if err != nil {
			log.Printf("✗ Gagal hash password untuk %s: %v\n", buyerData.Email, err)
			continue
		}

		// Create buyer
		buyer := &models.Buyer{
			ID:         uuid.New(),
			Nama:       buyerData.Nama,
			Username:   buyerData.Username,
			Email:      buyerData.Email,
			Password:   hashedPassword,
			Telepon:    &buyerData.Telepon,
			IsActive:   true,
			IsVerified: false, // Set false by default, bisa diverifikasi manual
		}

		if err := db.Create(buyer).Error; err != nil {
			log.Printf("✗ Gagal membuat buyer %s: %v\n", buyerData.Email, err)
			continue
		}

		fmt.Printf("✓ Buyer berhasil dibuat: %s (%s)\n", buyer.Nama, buyer.Email)
		successCount++
	}

	fmt.Println()
	fmt.Println("=== Seeding Summary ===")
	fmt.Printf("✓ Berhasil: %d buyer\n", successCount)
	fmt.Printf("⊘ Dilewati: %d buyer (sudah ada)\n", skipCount)
	fmt.Printf("✗ Gagal: %d buyer\n", len(buyers)-successCount-skipCount)

	if successCount > 0 {
		fmt.Println()
		fmt.Println("📝 Info Login:")
		fmt.Println("   Username/Email: johndoe / john.doe@example.com")
		fmt.Println("   Password: buyer123")
		fmt.Println()
		fmt.Println("   Username/Email: janesmith / jane.smith@example.com")
		fmt.Println("   Password: buyer123")
		fmt.Println()
		fmt.Println("   Username/Email: bobwilson / bob.wilson@example.com")
		fmt.Println("   Password: buyer123")
	}
}
