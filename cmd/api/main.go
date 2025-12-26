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
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()
	utils.SetJWTSecret(cfg.JWTSecret)

	database.InitDB(cfg)
	db := database.GetDB()

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
	adminRepo := repositories.NewAdminRepository(db)
	adminSessionRepo := repositories.NewAdminSessionRepository(db)
	buyerRepo := repositories.NewBuyerRepository(db)
	alamatBuyerRepo := repositories.NewAlamatBuyerRepository(db)
	heroSectionRepo := repositories.NewHeroSectionRepository(db)
	bannerEventPromoRepo := repositories.NewBannerEventPromoRepository(db)

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
	authService := services.NewAuthService(adminRepo, adminSessionRepo)
	adminService := services.NewAdminService(adminRepo, adminSessionRepo)
	masterService := services.NewMasterService(kategoriRepo, merekRepo, kondisiRepo, kondisiPaketRepo, sumberRepo)
	buyerService := services.NewBuyerService(buyerRepo, alamatBuyerRepo)
	alamatBuyerService := services.NewAlamatBuyerService(alamatBuyerRepo, buyerRepo)
	heroSectionService := services.NewHeroSectionService(heroSectionRepo)
	bannerEventPromoService := services.NewBannerEventPromoService(bannerEventPromoRepo)

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
	authController := controllers.NewAuthController(authService)
	adminController := controllers.NewAdminController(adminService)
	masterController := controllers.NewMasterController(masterService)
	buyerController := controllers.NewBuyerController(buyerService)
	alamatBuyerController := controllers.NewAlamatBuyerController(alamatBuyerService)
	heroSectionController := controllers.NewHeroSectionController(heroSectionService)
	bannerEventPromoController := controllers.NewBannerEventPromoController(bannerEventPromoService)

	router := gin.Default()

	routes.SetupRoutes(
		router,
		kategoriController, merekController, kondisiController, kondisiPaketController, sumberController,
		warehouseController, tipeProdukController, diskonKategoriController, bannerTipeProdukController,
		produkController, authController, adminController, masterController,
		buyerController, alamatBuyerController,
		heroSectionController, bannerEventPromoController,
	)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
