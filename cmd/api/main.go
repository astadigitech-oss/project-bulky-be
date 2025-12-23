package main

import (
	"log"
	"os"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/controllers"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/internal/routes"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/database"

	"github.com/gin-gonic/gin"
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

	// Auto migrate models (only in development)
	database.AutoMigrate(cfg.AppEnv)

	// Set Gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize repositories
	kategoriRepo := repositories.NewKategoriProdukRepository(db)
	merekRepo := repositories.NewMerekProdukRepository(db)
	kondisiRepo := repositories.NewKondisiProdukRepository(db)
	kondisiPaketRepo := repositories.NewKondisiPaketRepository(db)
	sumberRepo := repositories.NewSumberProdukRepository(db)
	warehouseRepo := repositories.NewWarehouseRepository(db)
	tipeProdukRepo := repositories.NewTipeProdukRepository(db)
	diskonKategoriRepo := repositories.NewDiskonKategoriRepository(db)
	bannerTipeProdukRepo := repositories.NewBannerTipeProdukRepository(db)
	produkRepo := repositories.NewProdukRepository(db)
	produkGambarRepo := repositories.NewProdukGambarRepository(db)
	produkDokumenRepo := repositories.NewProdukDokumenRepository(db)

	// Initialize services
	kategoriService := services.NewKategoriProdukService(kategoriRepo)
	merekService := services.NewMerekProdukService(merekRepo)
	kondisiService := services.NewKondisiProdukService(kondisiRepo)
	kondisiPaketService := services.NewKondisiPaketService(kondisiPaketRepo)
	sumberService := services.NewSumberProdukService(sumberRepo)
	warehouseService := services.NewWarehouseService(warehouseRepo)
	tipeProdukService := services.NewTipeProdukService(tipeProdukRepo)
	diskonKategoriService := services.NewDiskonKategoriService(diskonKategoriRepo)
	bannerTipeProdukService := services.NewBannerTipeProdukService(bannerTipeProdukRepo)
	produkGambarService := services.NewProdukGambarService(produkGambarRepo)
	produkDokumenService := services.NewProdukDokumenService(produkDokumenRepo)
	produkService := services.NewProdukService(produkRepo, produkGambarRepo)
	masterService := services.NewMasterService(
		kategoriRepo, merekRepo, kondisiRepo, kondisiPaketRepo, sumberRepo,
	)

	// Initialize controllers
	kategoriController := controllers.NewKategoriProdukController(kategoriService)
	merekController := controllers.NewMerekProdukController(merekService)
	kondisiController := controllers.NewKondisiProdukController(kondisiService)
	kondisiPaketController := controllers.NewKondisiPaketController(kondisiPaketService)
	sumberController := controllers.NewSumberProdukController(sumberService)
	warehouseController := controllers.NewWarehouseController(warehouseService)
	tipeProdukController := controllers.NewTipeProdukController(tipeProdukService)
	diskonKategoriController := controllers.NewDiskonKategoriController(diskonKategoriService)
	bannerTipeProdukController := controllers.NewBannerTipeProdukController(bannerTipeProdukService)
	produkController := controllers.NewProdukController(produkService, produkGambarService, produkDokumenService)
	masterController := controllers.NewMasterController(masterService)

	// Initialize router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(
		router,
		kategoriController,
		merekController,
		kondisiController,
		kondisiPaketController,
		sumberController,
		warehouseController,
		tipeProdukController,
		diskonKategoriController,
		bannerTipeProdukController,
		produkController,
		masterController,
	)

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
