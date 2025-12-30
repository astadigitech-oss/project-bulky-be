package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Create Admin Account ===")
	fmt.Println()

	// Get nama
	fmt.Print("Nama: ")
	nama, _ := reader.ReadString('\n')
	nama = strings.TrimSpace(nama)
	if nama == "" {
		log.Fatal("Nama tidak boleh kosong")
	}

	// Get email
	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)
	if email == "" {
		log.Fatal("Email tidak boleh kosong")
	}

	// Check if email exists
	var count int64
	db.Model(&models.Admin{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		log.Fatal("Email sudah terdaftar")
	}

	// Get password
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)
	if len(password) < 6 {
		log.Fatal("Password minimal 6 karakter")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal("Gagal hash password:", err)
	}

	// Set role_id to super admin role
	var adminRole models.Role
	if err := db.Where("kode = ?", "SUPER_ADMIN").First(&adminRole).Error; err != nil {
		log.Fatal("Gagal mendapatkan role SUPER_ADMIN:", err)
	}

	// Create admin
	admin := &models.Admin{
		ID:       uuid.New(),
		Nama:     nama,
		Email:    email,
		Password: hashedPassword,
		RoleID:   adminRole.ID,
		IsActive: true,
	}

	if err := db.Create(admin).Error; err != nil {
		log.Fatal("Gagal membuat admin:", err)
	}

	fmt.Println()
	fmt.Println("✓ Admin berhasil dibuat!")
	fmt.Printf("  ID: %s\n", admin.ID)
	fmt.Printf("  Nama: %s\n", admin.Nama)
	fmt.Printf("  Email: %s\n", admin.Email)
}
