package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/controllers"
	"project-bulky-be/internal/middleware"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/internal/routes"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/database"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
		log.Println("Running in production mode")
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
	metodePembayaranGroupRepo := repositories.NewMetodePembayaranGroupRepository(db)
	metodePembayaranRepo := repositories.NewMetodePembayaranRepository(db)
	dokumenKebijakanRepo := repositories.NewDokumenKebijakanRepository(db)
	disclaimerRepo := repositories.NewDisclaimerRepository(db)
	formulirPartaiBesarRepo := repositories.NewFormulirPartaiBesarRepository(db)
	whatsappHandlerRepo := repositories.NewWhatsAppHandlerRepository(db)
	faqRepo := repositories.NewFAQRepository(db)
	jadwalGudangRepo := repositories.NewJadwalGudangRepository(db)
	blogRepo := repositories.NewBlogRepository(db)
	kategoriBlogRepo := repositories.NewKategoriBlogRepository(db)
	labelBlogRepo := repositories.NewLabelBlogRepository(db)
	videoRepo := repositories.NewVideoRepository(db)
	kategoriVideoRepo := repositories.NewKategoriVideoRepository(db)
	kuponRepo := repositories.NewKuponRepository(db)
	dasborRepo := repositories.NewDasborRepository(db)
	disclaimerConsentRepo := repositories.NewBuyerDisclaimerConsentRepository(db)
	delivereeVehicleTypeRepo := repositories.NewDelivereeVehicleTypeRepository(db)
	forwarderMappingRepo := repositories.NewForwarderMappingRepository(db)

	// Auth V2 repositories
	authRepo := repositories.NewAuthRepository(db)
	activityLogRepo := repositories.NewActivityLogRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)

	// Initialize services
	reorderService := services.NewReorderService(db)
	kategoriService := services.NewKategoriProdukService(kategoriRepo, cfg)
	merekService := services.NewMerekProdukService(merekRepo, cfg)
	kondisiService := services.NewKondisiProdukService(kondisiRepo, reorderService)
	kondisiPaketService := services.NewKondisiPaketService(kondisiPaketRepo, reorderService)
	sumberService := services.NewSumberProdukService(sumberRepo)
	warehouseService := services.NewWarehouseService(warehouseRepo, jadwalGudangRepo)
	tipeProdukService := services.NewTipeProdukService(tipeProdukRepo)
	diskonKategoriService := services.NewDiskonKategoriService(diskonKategoriRepo)
	bannerTipeProdukService := services.NewBannerTipeProdukService(bannerTipeProdukRepo, tipeProdukRepo, reorderService, cfg)
	produkGambarService := services.NewProdukGambarService(produkGambarRepo, cfg)
	produkDokumenService := services.NewProdukDokumenService(produkDokumenRepo, cfg)
	produkService := services.NewProdukService(produkRepo, produkGambarRepo, produkDokumenRepo, warehouseRepo, tipeProdukRepo, cfg, db)
	authService := services.NewAuthService(adminRepo, adminSessionRepo)
	adminService := services.NewAdminService(adminRepo, adminSessionRepo, roleRepo)
	masterService := services.NewMasterService(kategoriRepo, merekRepo, kondisiRepo, kondisiPaketRepo, sumberRepo)
	buyerService := services.NewBuyerService(buyerRepo, alamatBuyerRepo)
	alamatBuyerService := services.NewAlamatBuyerService(alamatBuyerRepo, buyerRepo)
	heroSectionService := services.NewHeroSectionService(heroSectionRepo, cfg)
	bannerEventPromoService := services.NewBannerEventPromoService(bannerEventPromoRepo, reorderService, kategoriService, cfg)
	ulasanService := services.NewUlasanService(ulasanRepo, pesananItemRepo, pesananRepo, cfg.UploadPath, cfg.BaseURL)
	ulasanAdminService := services.NewUlasanAdminService(ulasanRepo)
	activityLogService := services.NewActivityLogService(activityLogRepo)
	delivereeVehicleTypeService := services.NewDelivereeVehicleTypeService(delivereeVehicleTypeRepo, warehouseRepo, activityLogService)
	forwarderMappingService := services.NewForwarderMappingService(forwarderMappingRepo)
	shippingService := services.NewShippingService(db, delivereeVehicleTypeService)
	pesananAdminService := services.NewPesananAdminService(pesananRepo, shippingService, db)
	forceUpdateService := services.NewForceUpdateService(forceUpdateRepo)
	modeMaintenanceService := services.NewModeMaintenanceService(modeMaintenanceRepo)
	ppnService := services.NewPPNService(ppnRepo)
	metodePembayaranService := services.NewMetodePembayaranService(metodePembayaranRepo, metodePembayaranGroupRepo)
	dokumenKebijakanService := services.NewDokumenKebijakanService(dokumenKebijakanRepo)
	disclaimerService := services.NewDisclaimerService(disclaimerRepo)
	disclaimerConsentService := services.NewBuyerDisclaimerConsentService(disclaimerConsentRepo)
	emailService := services.NewEmailService()
	formulirPartaiBesarService := services.NewFormulirPartaiBesarService(formulirPartaiBesarRepo, kategoriRepo, reorderService, emailService)
	whatsappHandlerService := services.NewWhatsAppHandlerService(whatsappHandlerRepo)
	faqService := services.NewFAQService(faqRepo, reorderService)
	blogService := services.NewBlogService(blogRepo, kategoriBlogRepo, labelBlogRepo, cfg)
	kategoriBlogService := services.NewKategoriBlogService(kategoriBlogRepo)
	labelBlogService := services.NewLabelBlogService(labelBlogRepo)
	videoService := services.NewVideoService(videoRepo, kategoriVideoRepo, cfg)
	kategoriVideoService := services.NewKategoriVideoService(kategoriVideoRepo)

	// Recovery: kembalikan video yang stuck di status 'processing' ke 'failed' saat startup
	videoService.RecoverStuckJobs(context.Background())
	kuponService := services.NewKuponService(kuponRepo, kategoriRepo, db)
	dasborService := services.NewDasborService(dasborRepo)

	// Auth V2 services
	authV2Service := services.NewAuthV2Service(authRepo, activityLogRepo)
	roleService := services.NewRoleService(roleRepo)
	permissionService := services.NewPermissionService(permissionRepo)

	// Initialize controllers
	kategoriController := controllers.NewKategoriProdukController(kategoriService, cfg, activityLogService)
	merekController := controllers.NewMerekProdukController(merekService, cfg, activityLogService)
	kondisiController := controllers.NewKondisiProdukController(kondisiService, reorderService, activityLogService)
	kondisiPaketController := controllers.NewKondisiPaketController(kondisiPaketService, reorderService, activityLogService)
	sumberController := controllers.NewSumberProdukController(sumberService, activityLogService)
	warehouseController := controllers.NewWarehouseController(warehouseService, activityLogService)
	tipeProdukController := controllers.NewTipeProdukController(tipeProdukService)
	diskonKategoriController := controllers.NewDiskonKategoriController(diskonKategoriService, activityLogService)
	bannerTipeProdukController := controllers.NewBannerTipeProdukController(bannerTipeProdukService, reorderService, cfg, activityLogService)
	produkController := controllers.NewProdukController(produkService, produkGambarService, produkDokumenService, activityLogService)
	authController := controllers.NewAuthController(authService)
	adminController := controllers.NewAdminController(adminService, activityLogService)
	masterController := controllers.NewMasterController(masterService)
	buyerController := controllers.NewBuyerController(buyerService, activityLogService)
	alamatBuyerController := controllers.NewAlamatBuyerController(alamatBuyerService, activityLogService)
	heroSectionController := controllers.NewHeroSectionController(heroSectionService, cfg, activityLogService)
	bannerEventPromoController := controllers.NewBannerEventPromoController(bannerEventPromoService, reorderService, cfg, activityLogService)
	ulasanController := controllers.NewUlasanController(ulasanService)
	ulasanAdminController := controllers.NewUlasanAdminController(ulasanAdminService, activityLogService)
	pesananAdminController := controllers.NewPesananAdminController(pesananAdminService, activityLogService)
	forceUpdateController := controllers.NewForceUpdateController(forceUpdateService, activityLogService)
	modeMaintenanceController := controllers.NewModeMaintenanceController(modeMaintenanceService, activityLogService)
	ppnController := controllers.NewPPNController(ppnService, activityLogService)
	metodePembayaranController := controllers.NewMetodePembayaranController(metodePembayaranService, reorderService, activityLogService)
	dokumenKebijakanController := controllers.NewDokumenKebijakanController(dokumenKebijakanService, activityLogService)
	disclaimerController := controllers.NewDisclaimerController(disclaimerService, activityLogService)
	disclaimerConsentController := controllers.NewBuyerDisclaimerConsentController(disclaimerConsentService)
	formulirPartaiBesarController := controllers.NewFormulirPartaiBesarController(formulirPartaiBesarService, reorderService, activityLogService)
	whatsappHandlerController := controllers.NewWhatsAppHandlerController(whatsappHandlerService)
	faqController := controllers.NewFAQController(faqService, activityLogService)
	blogController := controllers.NewBlogController(blogService, cfg, activityLogService)
	kategoriBlogController := controllers.NewKategoriBlogController(kategoriBlogService, reorderService, activityLogService)
	labelBlogController := controllers.NewLabelBlogController(labelBlogService, reorderService, activityLogService)
	videoController := controllers.NewVideoController(videoService, kategoriVideoService, cfg, activityLogService)
	kategoriVideoController := controllers.NewKategoriVideoController(kategoriVideoService, reorderService, activityLogService)
	kuponController := controllers.NewKuponController(kuponService, activityLogService)
	dasborController := controllers.NewDasborController(dasborService)
	internalUploadController := controllers.NewInternalUploadController(cfg)
	assetMigrationController := controllers.NewAssetMigrationController(db, cfg)
	delivereeWebhookService := services.NewDelivereeWebhookService(pesananRepo, db)
	delivereeWebhookController := controllers.NewDelivereeWebhookController(delivereeWebhookService, cfg.DelivereeWebhookAuthorization)
	forwarderWebhookService := services.NewForwarderWebhookService(pesananRepo, db)
	forwarderWebhookController := controllers.NewForwarderWebhookController(forwarderWebhookService, cfg.ForwarderWebhookAuthorization)
	delivereeVehicleTypeController := controllers.NewDelivereeVehicleTypeController(delivereeVehicleTypeService, activityLogService)
	forwarderMappingController := controllers.NewForwarderMappingController(forwarderMappingService, activityLogService)

	// Auth V2 controllers
	authV2Controller := controllers.NewAuthV2Controller(authV2Service, adminService, buyerService)
	activityLogController := controllers.NewActivityLogController(activityLogService)
	roleController := controllers.NewRoleController(roleService)
	permissionController := controllers.NewPermissionController(permissionService)

	router := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024, // 500MB untuk upload file video
	})
	router.Use(logger.New(logger.Config{
		// Jangan log request health check agar log sistem tidak tenggelam
		// oleh ping berkala dari orchestrator (Docker HEALTHCHECK / Dokploy).
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api/health")
		},
	}))
	router.Use(middleware.CORSMiddleware())

	routes.SetupRoutes(
		router,
		cfg,
		dasborController,
		kategoriController, merekController, kondisiController, kondisiPaketController, sumberController,
		warehouseController, tipeProdukController, diskonKategoriController, bannerTipeProdukController,
		produkController, authController, adminController, masterController,
		buyerController, alamatBuyerController,
		heroSectionController, bannerEventPromoController,
		ulasanController,
		ulasanAdminController, pesananAdminController,
		forceUpdateController, modeMaintenanceController,
		ppnController,
		metodePembayaranController,
		dokumenKebijakanController, disclaimerController, disclaimerConsentController,
		formulirPartaiBesarController, whatsappHandlerController,
		faqController,
		blogController, kategoriBlogController, labelBlogController,
		videoController, kategoriVideoController,
		kuponController,
		internalUploadController,
		assetMigrationController,
		delivereeWebhookController,
		forwarderWebhookController,
		delivereeVehicleTypeController,
		forwarderMappingController,
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

	// Jalankan server di goroutine agar sinyal shutdown bisa ditangkap.
	// Catatan: saat Shutdown dipanggil, fasthttp mengembalikan nil dari Serve,
	// sehingga hanya error selain shutdown yang dianggap fatal.
	go func() {
		log.Printf("Server is running on port %s", port)
		if err := router.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown: tunggu sinyal SIGTERM/SIGINT, lalu stop server
	// setelah request yang sedang berjalan (termasuk upload video) selesai
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")
	if err := router.ShutdownWithContext(context.Background()); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}
