package main

import (
	"log"
	"os"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/controllers"
	"project-bulky-be/internal/middleware"
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

	// Initialize JWT config with 24 hour access token (single token)
	utils.SetJWTConfig(cfg.JWTSecret, cfg.JWTAccessDuration)

	// Initialize custom validators
	utils.InitCustomValidators()

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
	pesananRepo := repositories.NewPesananRepository(db)
	pesananItemRepo := repositories.NewPesananItemRepository(db)
	ulasanRepo := repositories.NewUlasanRepository(db)
	forceUpdateRepo := repositories.NewForceUpdateRepository(db)
	modeMaintenanceRepo := repositories.NewModeMaintenanceRepository(db)
	ppnRepo := repositories.NewPPNRepository(db)

	// Auth V2 repositories
	authRepo := repositories.NewAuthRepository(db)
	activityLogRepo := repositories.NewActivityLogRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)

	// Initialize services
	kategoriService := services.NewKategoriProdukService(kategoriRepo, cfg)
	merekService := services.NewMerekProdukService(merekRepo, cfg)
	kondisiService := services.NewKondisiProdukService(kondisiRepo)
	kondisiPaketService := services.NewKondisiPaketService(kondisiPaketRepo)
	sumberService := services.NewSumberProdukService(sumberRepo)
	warehouseService := services.NewWarehouseService(warehouseRepo)
	tipeProdukService := services.NewTipeProdukService(tipeProdukRepo)
	diskonKategoriService := services.NewDiskonKategoriService(diskonKategoriRepo)
	bannerTipeProdukService := services.NewBannerTipeProdukService(bannerTipeProdukRepo, cfg)
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
	ulasanService := services.NewUlasanService(ulasanRepo, pesananItemRepo, pesananRepo, cfg.UploadPath, cfg.BaseURL)
	forceUpdateService := services.NewForceUpdateService(forceUpdateRepo, cfg.PlayStoreURL, cfg.AppStoreURL)
	modeMaintenanceService := services.NewModeMaintenanceService(modeMaintenanceRepo)
	ppnService := services.NewPPNService(ppnRepo)

	// Auth V2 services
	authV2Service := services.NewAuthV2Service(authRepo, activityLogRepo)
	activityLogService := services.NewActivityLogService(activityLogRepo)
	roleService := services.NewRoleService(roleRepo)
	permissionService := services.NewPermissionService(permissionRepo)

	// Initialize controllers
	kategoriController := controllers.NewKategoriProdukController(kategoriService, cfg)
	merekController := controllers.NewMerekProdukController(merekService, cfg)
	kondisiController := controllers.NewKondisiProdukController(kondisiService)
	kondisiPaketController := controllers.NewKondisiPaketController(kondisiPaketService)
	sumberController := controllers.NewSumberProdukController(sumberService)
	warehouseController := controllers.NewWarehouseController(warehouseService)
	tipeProdukController := controllers.NewTipeProdukController(tipeProdukService)
	diskonKategoriController := controllers.NewDiskonKategoriController(diskonKategoriService)
	bannerTipeProdukController := controllers.NewBannerTipeProdukController(bannerTipeProdukService, cfg)
	produkController := controllers.NewProdukController(produkService, produkGambarService, produkDokumenService)
	authController := controllers.NewAuthController(authService)
	adminController := controllers.NewAdminController(adminService)
	masterController := controllers.NewMasterController(masterService)
	buyerController := controllers.NewBuyerController(buyerService)
	alamatBuyerController := controllers.NewAlamatBuyerController(alamatBuyerService)
	heroSectionController := controllers.NewHeroSectionController(heroSectionService, cfg)
	bannerEventPromoController := controllers.NewBannerEventPromoController(bannerEventPromoService, cfg)
	ulasanController := controllers.NewUlasanController(ulasanService)
	forceUpdateController := controllers.NewForceUpdateController(forceUpdateService)
	modeMaintenanceController := controllers.NewModeMaintenanceController(modeMaintenanceService)
	appStatusController := controllers.NewAppStatusController(forceUpdateService, modeMaintenanceService)
	ppnController := controllers.NewPPNController(ppnService)

	// Auth V2 controllers
	authV2Controller := controllers.NewAuthV2Controller(authV2Service, adminService, buyerService)
	activityLogController := controllers.NewActivityLogController(activityLogService)
	roleController := controllers.NewRoleController(roleService)
	permissionController := controllers.NewPermissionController(permissionService)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	routes.SetupRoutes(
		router,
		kategoriController, merekController, kondisiController, kondisiPaketController, sumberController,
		warehouseController, tipeProdukController, diskonKategoriController, bannerTipeProdukController,
		produkController, authController, adminController, masterController,
		buyerController, alamatBuyerController,
		heroSectionController, bannerEventPromoController,
		ulasanController,
		forceUpdateController, modeMaintenanceController, appStatusController,
		ppnController,
	)

	// Setup Auth V2 routes (new authentication system with roles & permissions)
	routes.SetupAuthV2Routes(
		router,
		authV2Controller,
		roleController,
		permissionController,
		activityLogController,
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
