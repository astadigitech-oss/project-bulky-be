package routes

import (
	"net/http"
	"os"

	"project-bulky-be/internal/controllers"
	"project-bulky-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	kategoriController *controllers.KategoriProdukController,
	merekController *controllers.MerekProdukController,
	kondisiController *controllers.KondisiProdukController,
	kondisiPaketController *controllers.KondisiPaketController,
	sumberController *controllers.SumberProdukController,
	warehouseController *controllers.WarehouseController,
	tipeProdukController *controllers.TipeProdukController,
	diskonKategoriController *controllers.DiskonKategoriController,
	bannerTipeProdukController *controllers.BannerTipeProdukController,
	produkController *controllers.ProdukController,
	authController *controllers.AuthController,
	adminController *controllers.AdminController,
	masterController *controllers.MasterController,
	buyerController *controllers.BuyerController,
	alamatBuyerController *controllers.AlamatBuyerController,
	heroSectionController *controllers.HeroSectionController,
	bannerEventPromoController *controllers.BannerEventPromoController,
	ulasanController *controllers.UlasanController,
	forceUpdateController *controllers.ForceUpdateController,
	modeMaintenanceController *controllers.ModeMaintenanceController,
	appStatusController *controllers.AppStatusController,
	ppnController *controllers.PPNController,
	metodePembayaranController *controllers.MetodePembayaranController,
	dokumenKebijakanController *controllers.DokumenKebijakanController,
	disclaimerController *controllers.DisclaimerController,
	formulirPartaiBesarController *controllers.FormulirPartaiBesarController,
	whatsappHandlerController *controllers.WhatsAppHandlerController,
	faqController *controllers.FAQController,
	blogController *controllers.BlogController,
	kategoriBlogController *controllers.KategoriBlogController,
	labelBlogController *controllers.LabelBlogController,
	videoController *controllers.VideoController,
	kategoriVideoController *controllers.KategoriVideoController,
) {
	// Health check
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "Server is running"})
	})

	// Serve static files from uploads folder
	uploadPath := os.Getenv("UPLOAD_PATH")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}
	router.Static("/uploads", uploadPath)

	// API routes (tanpa versioning)
	v1 := router.Group("/api")
	{
		// Public Auth Routes (Legacy - Commented out, using AuthV2 instead)
		// Uncomment if you want to use old auth system
		/*
			auth := v1.Group("/auth")
			{
				auth.POST("/login", authController.Login)
				auth.POST("/refresh", authController.RefreshToken)
			}
		*/

		// Protected Auth Routes (Legacy - Commented out, using AuthV2 instead)
		/*
			authProtected := v1.Group("/auth")
			authProtected.Use(middleware.AuthMiddleware())
			{
				authProtected.POST("/logout", authController.Logout)
				authProtected.GET("/me", authController.Me)
				authProtected.PUT("/profile", authController.UpdateProfile)
				authProtected.PUT("/change-password", authController.ChangePassword)
			}
		*/

		// Admin Management Routes (Protected)
		admin := v1.Group("/panel/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminOnly())
		admin.Use(middleware.RequirePermission("admin:manage"))
		{
			admin.GET("", adminController.FindAll)
			admin.GET("/:id", adminController.FindByID)
			admin.POST("", adminController.Create)
			admin.PUT("/:id", adminController.Update)
			admin.DELETE("/:id", adminController.Delete)
			admin.PATCH("/:id/toggle-status", adminController.ToggleStatus)
			admin.PUT("/:id/reset-password", adminController.ResetPassword)
		}

		// Buyer Management Routes (Admin Side)
		buyerManagement := v1.Group("/panel/buyer")
		buyerManagement.Use(middleware.AuthMiddleware())
		buyerManagement.Use(middleware.AdminOnly())
		{
			// Read permission untuk list, detail, statistik, chart
			buyerManagement.GET("", middleware.RequirePermission("buyer:read"), buyerController.FindAll)
			buyerManagement.GET("/statistik", middleware.RequirePermission("buyer:read"), buyerController.GetStatistik)
			buyerManagement.GET("/chart", middleware.RequirePermission("buyer:read"), buyerController.GetChart)
			buyerManagement.GET("/:id", middleware.RequirePermission("buyer:read"), buyerController.FindByID)

			// Manage permission untuk delete dan reset password
			buyerManagement.DELETE("/:id", middleware.RequirePermission("buyer:manage"), buyerController.Delete)
			buyerManagement.PUT("/:id/reset-password", middleware.RequirePermission("buyer:manage"), buyerController.ResetPassword)
		}

		// Alamat Buyer Routes (Buyer Only)
		// Data wilayah dari Google Maps API (frontend)
		alamatBuyer := v1.Group("/buyer/alamat")
		alamatBuyer.Use(middleware.AuthMiddleware())
		alamatBuyer.Use(middleware.BuyerOnly())
		{
			alamatBuyer.GET("", alamatBuyerController.FindAll)
			alamatBuyer.GET("/:id", alamatBuyerController.FindByID)
			alamatBuyer.POST("", alamatBuyerController.Create)
			alamatBuyer.PUT("/:id", alamatBuyerController.Update)
			alamatBuyer.DELETE("/:id", alamatBuyerController.Delete)
			alamatBuyer.PATCH("/:id/set-default", alamatBuyerController.SetDefault)
		}

		// Kategori Produk - Public (Read Only)
		kategoriPublic := v1.Group("/kategori-produk")
		{
			kategoriPublic.GET("", kategoriController.FindAll)
			kategoriPublic.GET("/:id", kategoriController.FindByID)
			kategoriPublic.GET("/slug/:slug", kategoriController.FindBySlug)
		}

		// Kategori Produk - Admin
		kategoriAdmin := v1.Group("/panel/kategori-produk")
		kategoriAdmin.Use(middleware.AuthMiddleware())
		kategoriAdmin.Use(middleware.AdminOnly())
		{
			kategoriAdmin.GET("", middleware.RequirePermission("kategori:read"), kategoriController.FindAll)
			kategoriAdmin.GET("/dropdown", middleware.RequirePermission("kategori:read"), kategoriController.Dropdown)
			kategoriAdmin.GET("/:id", middleware.RequirePermission("kategori:read"), kategoriController.FindByID)
			kategoriAdmin.POST("", middleware.RequirePermission("kategori:manage"), kategoriController.Create)
			kategoriAdmin.PUT("/:id", middleware.RequirePermission("kategori:manage"), kategoriController.Update)
			kategoriAdmin.DELETE("/:id", middleware.RequirePermission("kategori:manage"), kategoriController.Delete)
			kategoriAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("kategori:manage"), kategoriController.ToggleStatus)
		}

		// Merek Produk - Public (Read Only)
		merekPublic := v1.Group("/merek-produk")
		{
			merekPublic.GET("", merekController.FindAll)
			merekPublic.GET("/:id", merekController.FindByID)
			merekPublic.GET("/slug/:slug", merekController.FindBySlug)
		}

		// Merek Produk - Admin
		merekAdmin := v1.Group("/panel/merek-produk")
		merekAdmin.Use(middleware.AuthMiddleware())
		merekAdmin.Use(middleware.AdminOnly())
		{
			merekAdmin.GET("", middleware.RequirePermission("brand:read"), merekController.FindAll)
			merekAdmin.GET("/dropdown", middleware.RequirePermission("brand:read"), merekController.Dropdown)
			merekAdmin.GET("/:id", middleware.RequirePermission("brand:read"), merekController.FindByID)
			merekAdmin.POST("", middleware.RequirePermission("brand:manage"), merekController.Create)
			merekAdmin.PUT("/:id", middleware.RequirePermission("brand:manage"), merekController.Update)
			merekAdmin.DELETE("/:id", middleware.RequirePermission("brand:manage"), merekController.Delete)
			merekAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("brand:manage"), merekController.ToggleStatus)
		}

		// Kondisi Produk - Public (Read Only)
		kondisiPublic := v1.Group("/kondisi-produk")
		{
			kondisiPublic.GET("", kondisiController.FindAll)
			kondisiPublic.GET("/:id", kondisiController.FindByID)
			kondisiPublic.GET("/slug/:slug", kondisiController.FindBySlug)
		}

		// Kondisi Produk - Admin
		kondisiAdmin := v1.Group("/panel/kondisi-produk")
		kondisiAdmin.Use(middleware.AuthMiddleware())
		kondisiAdmin.Use(middleware.AdminOnly())
		{
			kondisiAdmin.GET("", middleware.RequirePermission("kondisi:read"), kondisiController.FindAll)
			kondisiAdmin.GET("/dropdown", middleware.RequirePermission("kondisi:read"), kondisiController.Dropdown)
			kondisiAdmin.GET("/:id", middleware.RequirePermission("kondisi:read"), kondisiController.FindByID)
			kondisiAdmin.POST("", middleware.RequirePermission("kondisi:manage"), kondisiController.Create)
			kondisiAdmin.PUT("/:id", middleware.RequirePermission("kondisi:manage"), kondisiController.Update)
			kondisiAdmin.DELETE("/:id", middleware.RequirePermission("kondisi:manage"), kondisiController.Delete)
			kondisiAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("kondisi:manage"), kondisiController.ToggleStatus)
			kondisiAdmin.PUT("/reorder", middleware.RequirePermission("kondisi:manage"), kondisiController.Reorder)
			kondisiAdmin.PATCH("/:id/reorder", middleware.RequirePermission("kondisi:manage"), kondisiController.ReorderByDirection)
		}

		// Kondisi Paket - Public (Read Only)
		paketPublic := v1.Group("/kondisi-paket")
		{
			paketPublic.GET("", kondisiPaketController.FindAll)
			paketPublic.GET("/:id", kondisiPaketController.FindByID)
			paketPublic.GET("/slug/:slug", kondisiPaketController.FindBySlug)
		}

		// Kondisi Paket - Admin
		paketAdmin := v1.Group("/panel/kondisi-paket")
		paketAdmin.Use(middleware.AuthMiddleware())
		paketAdmin.Use(middleware.AdminOnly())
		{
			paketAdmin.GET("", middleware.RequirePermission("kondisi:read"), kondisiPaketController.FindAll)
			paketAdmin.GET("/dropdown", middleware.RequirePermission("kondisi:read"), kondisiPaketController.Dropdown)
			paketAdmin.GET("/:id", middleware.RequirePermission("kondisi:read"), kondisiPaketController.FindByID)
			paketAdmin.POST("", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.Create)
			paketAdmin.PUT("/:id", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.Update)
			paketAdmin.DELETE("/:id", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.Delete)
			paketAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.ToggleStatus)
			paketAdmin.PUT("/reorder", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.Reorder)
			paketAdmin.PATCH("/:id/reorder", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.ReorderByDirection)
		}

		// Sumber Produk - Public (Read Only)
		sumberPublic := v1.Group("/sumber-produk")
		{
			sumberPublic.GET("", sumberController.FindAll)
			sumberPublic.GET("/:id", sumberController.FindByID)
			sumberPublic.GET("/slug/:slug", sumberController.FindBySlug)
		}

		// Sumber Produk - Admin
		sumberAdmin := v1.Group("/panel/sumber-produk")
		sumberAdmin.Use(middleware.AuthMiddleware())
		sumberAdmin.Use(middleware.AdminOnly())
		{
			sumberAdmin.GET("", middleware.RequirePermission("kondisi:read"), sumberController.FindAll)
			sumberAdmin.GET("/dropdown", middleware.RequirePermission("kondisi:read"), sumberController.Dropdown)
			sumberAdmin.GET("/:id", middleware.RequirePermission("kondisi:read"), sumberController.FindByID)
			sumberAdmin.POST("", middleware.RequirePermission("kondisi:manage"), sumberController.Create)
			sumberAdmin.PUT("/:id", middleware.RequirePermission("kondisi:manage"), sumberController.Update)
			sumberAdmin.DELETE("/:id", middleware.RequirePermission("kondisi:manage"), sumberController.Delete)
			sumberAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("kondisi:manage"), sumberController.ToggleStatus)
		}

		// Warehouse - Public (Singleton - Simplified)
		v1.GET("/public/warehouse", warehouseController.GetPublic)

		// Warehouse - Admin (Singleton Pattern)
		warehouseAdmin := v1.Group("/panel/warehouse")
		warehouseAdmin.Use(middleware.AuthMiddleware())
		warehouseAdmin.Use(middleware.AdminOnly())
		{
			warehouseAdmin.GET("", middleware.RequirePermission("operasional:read"), warehouseController.Get)
			warehouseAdmin.GET("/dropdown", middleware.RequirePermission("operasional:read"), warehouseController.Dropdown)
			warehouseAdmin.PUT("", middleware.RequirePermission("operasional:manage"), warehouseController.UpdateSingleton)
		}

		// Tipe Produk - Public (Read Only)
		// Note: Tipe produk is read-only (Paletbox, Container, Truckload) - managed via migration only
		tipeProdukPublic := v1.Group("/tipe-produk")
		{
			tipeProdukPublic.GET("", tipeProdukController.FindAll)
			tipeProdukPublic.GET("/with-produk", tipeProdukController.FindAllWithProduk)
		}

		// Tipe Produk - Admin (Read Only)
		tipeProdukAdmin := v1.Group("/panel/tipe-produk")
		tipeProdukAdmin.Use(middleware.AuthMiddleware())
		tipeProdukAdmin.Use(middleware.AdminOnly())
		{
			tipeProdukAdmin.GET("", middleware.RequirePermission("system:read"), tipeProdukController.FindAll)
			tipeProdukAdmin.GET("/dropdown", middleware.RequirePermission("system:read"), tipeProdukController.Dropdown)
			tipeProdukAdmin.GET("/with-produk", middleware.RequirePermission("system:read"), tipeProdukController.FindAllWithProduk)
		}

		// Diskon Kategori - Public (Read Only)
		diskonKategoriPublic := v1.Group("/diskon-kategori")
		{
			diskonKategoriPublic.GET("", diskonKategoriController.FindAll)
			diskonKategoriPublic.GET("/:id", diskonKategoriController.FindByID)
			diskonKategoriPublic.GET("/by-kategori/:kategori_id", diskonKategoriController.FindActiveByKategoriID)
		}

		// Diskon Kategori - Admin
		diskonKategoriAdmin := v1.Group("/panel/diskon-kategori")
		diskonKategoriAdmin.Use(middleware.AuthMiddleware())
		diskonKategoriAdmin.Use(middleware.AdminOnly())
		{
			diskonKategoriAdmin.GET("", middleware.RequirePermission("diskon:read"), diskonKategoriController.FindAll)
			diskonKategoriAdmin.GET("/:id", middleware.RequirePermission("diskon:read"), diskonKategoriController.FindByID)
			diskonKategoriAdmin.POST("", middleware.RequirePermission("diskon:manage"), diskonKategoriController.Create)
			diskonKategoriAdmin.PUT("/:id", middleware.RequirePermission("diskon:manage"), diskonKategoriController.Update)
			diskonKategoriAdmin.DELETE("/:id", middleware.RequirePermission("diskon:manage"), diskonKategoriController.Delete)
			diskonKategoriAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("diskon:manage"), diskonKategoriController.ToggleStatus)
		}

		// Banner Tipe Produk - Public (Read Only)
		bannerTipeProdukPublic := v1.Group("/banner-tipe-produk")
		{
			bannerTipeProdukPublic.GET("", bannerTipeProdukController.FindAll)
			bannerTipeProdukPublic.GET("/:id", bannerTipeProdukController.FindByID)
			bannerTipeProdukPublic.GET("/by-tipe/:tipe_produk_id", bannerTipeProdukController.FindByTipeProdukID)
		}

		// Banner Tipe Produk - Admin
		bannerTipeProdukAdmin := v1.Group("/panel/banner-tipe-produk")
		bannerTipeProdukAdmin.Use(middleware.AuthMiddleware())
		bannerTipeProdukAdmin.Use(middleware.AdminOnly())
		{
			bannerTipeProdukAdmin.GET("", middleware.RequirePermission("marketing:read"), bannerTipeProdukController.FindAll)
			bannerTipeProdukAdmin.GET("/:id", middleware.RequirePermission("marketing:read"), bannerTipeProdukController.FindByID)
			bannerTipeProdukAdmin.POST("", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.Create)
			bannerTipeProdukAdmin.PUT("/:id", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.Update)
			bannerTipeProdukAdmin.DELETE("/:id", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.Delete)
			bannerTipeProdukAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.ToggleStatus)
			bannerTipeProdukAdmin.PUT("/reorder", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.Reorder)
			bannerTipeProdukAdmin.PATCH("/:id/reorder", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.ReorderByDirection)
		}

		// Produk - Public (Read Only)
		produkPublic := v1.Group("/produk")
		{
			produkPublic.GET("", produkController.FindAll)
			produkPublic.GET("/:id", produkController.FindByID)
			produkPublic.GET("/slug/:slug", produkController.FindBySlug)
		}

		// Produk - Admin
		produkAdmin := v1.Group("/panel/produk")
		produkAdmin.Use(middleware.AuthMiddleware())
		produkAdmin.Use(middleware.AdminOnly())
		{
			// Read - produk:read
			produkAdmin.GET("", middleware.RequirePermission("produk:read"), produkController.FindAll)
			produkAdmin.GET("/:id", middleware.RequirePermission("produk:read"), produkController.FindByID)

			// Create - produk:create
			produkAdmin.POST("", middleware.RequirePermission("produk:create"), produkController.Create)

			// Update - produk:update
			produkAdmin.PUT("/:id", middleware.RequirePermission("produk:update"), produkController.Update)
			produkAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("produk:update"), produkController.ToggleStatus)
			produkAdmin.PATCH("/:id/update-stock", middleware.RequirePermission("produk:update"), produkController.UpdateStock)

			// Delete - produk:delete
			produkAdmin.DELETE("/:id", middleware.RequirePermission("produk:delete"), produkController.Delete)

			// Gambar - produk:update
			produkAdmin.POST("/:id/gambar", middleware.RequirePermission("produk:update"), produkController.AddGambar)
			produkAdmin.DELETE("/:id/gambar/:gambar_id", middleware.RequirePermission("produk:update"), produkController.DeleteGambar)
			produkAdmin.PATCH("/:id/gambar/:gambar_id/reorder", middleware.RequirePermission("produk:update"), produkController.ReorderGambar)
			produkAdmin.PATCH("/:id/gambar/:gambar_id/set-primary", middleware.RequirePermission("produk:update"), produkController.SetPrimaryGambar)

			// Dokumen - produk:update
			produkAdmin.POST("/:id/dokumen", middleware.RequirePermission("produk:update"), produkController.AddDokumen)
			produkAdmin.DELETE("/:id/dokumen/:dokumen_id", middleware.RequirePermission("produk:update"), produkController.DeleteDokumen)
		}

		// Master (Dropdown)
		master := v1.Group("/master")
		{
			master.GET("/dropdown", masterController.GetDropdown)
		}

		// Hero Section (Admin)
		heroSectionAdmin := v1.Group("/panel/hero-section")
		heroSectionAdmin.Use(middleware.AuthMiddleware())
		heroSectionAdmin.Use(middleware.AdminOnly())
		{
			heroSectionAdmin.GET("", middleware.RequirePermission("marketing:read"), heroSectionController.FindAll)
			heroSectionAdmin.GET("/:id", middleware.RequirePermission("marketing:read"), heroSectionController.FindByID)
			heroSectionAdmin.POST("", middleware.RequirePermission("marketing:manage"), heroSectionController.Create)
			heroSectionAdmin.PUT("/:id", middleware.RequirePermission("marketing:manage"), heroSectionController.Update)
			heroSectionAdmin.DELETE("/:id", middleware.RequirePermission("marketing:manage"), heroSectionController.Delete)
			heroSectionAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), heroSectionController.ToggleStatus)
			heroSectionAdmin.GET("/schedule", middleware.RequirePermission("marketing:read"), heroSectionController.GetSchedules)
		}

		// Hero Section (Public)
		heroSectionPublic := v1.Group("/hero-section")
		{
			heroSectionPublic.GET("/active", heroSectionController.GetActive)
		}

		// Banner Event Promo (Admin)
		bannerEventPromoAdmin := v1.Group("/panel/banner-event-promo")
		bannerEventPromoAdmin.Use(middleware.AuthMiddleware())
		bannerEventPromoAdmin.Use(middleware.AdminOnly())
		{
			bannerEventPromoAdmin.GET("", middleware.RequirePermission("marketing:read"), bannerEventPromoController.FindAll)
			// bannerEventPromoAdmin.GET("/schedule", middleware.RequirePermission("marketing:read"), bannerEventPromoController.GetSchedules)
			bannerEventPromoAdmin.GET("/:id", middleware.RequirePermission("marketing:read"), bannerEventPromoController.FindByID)
			bannerEventPromoAdmin.POST("", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.Create)
			bannerEventPromoAdmin.PUT("/:id", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.Update)
			bannerEventPromoAdmin.DELETE("/:id", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.Delete)
			bannerEventPromoAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.ToggleStatus)
			bannerEventPromoAdmin.PUT("/reorder", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.Reorder)
			bannerEventPromoAdmin.PATCH("/:id/reorder", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.ReorderByDirection)
		}

		// Banner Event Promo (Public)
		bannerEventPromoPublic := v1.Group("/banner-event-promo")
		{
			bannerEventPromoPublic.GET("/active", bannerEventPromoController.GetActive)
		}

		// Ulasan - Admin
		ulasanAdmin := v1.Group("/panel/ulasan")
		ulasanAdmin.Use(middleware.AuthMiddleware())
		ulasanAdmin.Use(middleware.AdminOnly())
		{
			ulasanAdmin.GET("", middleware.RequirePermission("ulasan:read"), ulasanController.AdminFindAll)
			ulasanAdmin.GET("/:id", middleware.RequirePermission("ulasan:read"), ulasanController.AdminFindByID)
			ulasanAdmin.PATCH("/:id/approve", middleware.RequirePermission("ulasan:manage"), ulasanController.Approve)
			ulasanAdmin.PATCH("/bulk-approve", middleware.RequirePermission("ulasan:manage"), ulasanController.BulkApprove)
			ulasanAdmin.DELETE("/:id", middleware.RequirePermission("ulasan:manage"), ulasanController.Delete)
		}

		// Ulasan - Buyer
		ulasanBuyer := v1.Group("/buyer/ulasan")
		ulasanBuyer.Use(middleware.AuthMiddleware())
		ulasanBuyer.Use(middleware.BuyerOnly())
		{
			ulasanBuyer.GET("/pending", ulasanController.GetPendingReviews)
			ulasanBuyer.GET("", ulasanController.BuyerFindAll)
			ulasanBuyer.POST("", ulasanController.Create)
		}

		// Ulasan - Public
		ulasanPublic := v1.Group("/public/produk")
		{
			ulasanPublic.GET("/:produk_id/ulasan", ulasanController.GetProdukUlasan)
			ulasanPublic.GET("/:produk_id/rating", ulasanController.GetProdukRating)
		}

		// Force Update - Super Admin Only
		forceUpdateAdmin := v1.Group("/panel/force-update")
		forceUpdateAdmin.Use(middleware.AuthMiddleware())
		forceUpdateAdmin.Use(middleware.SuperAdminOnly())
		{
			forceUpdateAdmin.GET("", forceUpdateController.GetAllForceUpdates)
			forceUpdateAdmin.GET("/:id", forceUpdateController.GetForceUpdateByID)
			forceUpdateAdmin.POST("", forceUpdateController.CreateForceUpdate)
			forceUpdateAdmin.PUT("/:id", forceUpdateController.UpdateForceUpdate)
			forceUpdateAdmin.DELETE("/:id", forceUpdateController.DeleteForceUpdate)
			forceUpdateAdmin.POST("/:id/set-active", forceUpdateController.SetActiveForceUpdate)
		}

		// Mode Maintenance - Super Admin Only
		modeMaintenanceAdmin := v1.Group("/panel/mode-maintenance")
		modeMaintenanceAdmin.Use(middleware.AuthMiddleware())
		modeMaintenanceAdmin.Use(middleware.SuperAdminOnly())
		{
			modeMaintenanceAdmin.GET("", modeMaintenanceController.GetAllMaintenances)
			modeMaintenanceAdmin.GET("/:id", modeMaintenanceController.GetMaintenanceByID)
			modeMaintenanceAdmin.POST("", modeMaintenanceController.CreateMaintenance)
			modeMaintenanceAdmin.PUT("/:id", modeMaintenanceController.UpdateMaintenance)
			modeMaintenanceAdmin.DELETE("/:id", modeMaintenanceController.DeleteMaintenance)
			modeMaintenanceAdmin.POST("/:id/activate", modeMaintenanceController.ActivateMaintenance)
			modeMaintenanceAdmin.POST("/:id/deactivate", modeMaintenanceController.DeactivateMaintenance)
		}

		// App Status - Public
		appStatus := v1.Group("/public")
		{
			appStatus.GET("/check-version", appStatusController.CheckVersion)
			appStatus.GET("/check-maintenance", appStatusController.CheckMaintenance)
			appStatus.GET("/app-status", appStatusController.AppStatus)
		}

		// PPN - Admin Only
		ppnAdmin := v1.Group("/panel/ppn")
		ppnAdmin.Use(middleware.AuthMiddleware())
		ppnAdmin.Use(middleware.AdminOnly())
		{
			ppnAdmin.GET("", middleware.RequirePermission("system:read"), ppnController.GetAll)
			ppnAdmin.GET("/:id", middleware.RequirePermission("system:read"), ppnController.GetByID)
			ppnAdmin.POST("", middleware.RequirePermission("system:manage"), ppnController.Create)
			ppnAdmin.PUT("/:id", middleware.RequirePermission("system:manage"), ppnController.Update)
			ppnAdmin.DELETE("/:id", middleware.RequirePermission("system:manage"), ppnController.Delete)
			ppnAdmin.PATCH("/:id/set-active", middleware.RequirePermission("system:manage"), ppnController.SetActive)
		}

		// Metode Pembayaran - Simplified Grouped API
		metodePembayaranAdmin := v1.Group("/panel/metode-pembayaran")
		metodePembayaranAdmin.Use(middleware.AuthMiddleware())
		metodePembayaranAdmin.Use(middleware.AdminOnly())
		{
			metodePembayaranAdmin.GET("", middleware.RequirePermission("pembayaran:read"), metodePembayaranController.GetAllGrouped)
			metodePembayaranAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("pembayaran:manage"), metodePembayaranController.ToggleMethodStatus)
			metodePembayaranAdmin.PATCH("/group/:urutan/toggle-status", middleware.RequirePermission("pembayaran:manage"), metodePembayaranController.ToggleGroupStatus)
		}

		// Metode Pembayaran - Public (Grouped, Active Only)
		v1.GET("/public/metode-pembayaran", metodePembayaranController.GetAllGroupedPublic)

		// Dokumen Kebijakan - Admin (Simplified - Fixed Pages)
		dokumenKebijakanAdmin := v1.Group("/panel/dokumen-kebijakan")
		dokumenKebijakanAdmin.Use(middleware.AuthMiddleware())
		dokumenKebijakanAdmin.Use(middleware.AdminOnly())
		{
			dokumenKebijakanAdmin.GET("", middleware.RequirePermission("system:read"), dokumenKebijakanController.GetAll)
			dokumenKebijakanAdmin.GET("/:id", middleware.RequirePermission("system:read"), dokumenKebijakanController.GetByID)
			dokumenKebijakanAdmin.PUT("/:id", middleware.RequirePermission("system:manage"), dokumenKebijakanController.Update)
		}

		// Dokumen Kebijakan - Public
		dokumenKebijakanPublic := v1.Group("/public/dokumen-kebijakan")
		{
			dokumenKebijakanPublic.GET("", dokumenKebijakanController.GetAllPublic)
			dokumenKebijakanPublic.GET("/:id", dokumenKebijakanController.GetByIDPublic)
		}

		// FAQ - Public
		v1.GET("/public/faq", faqController.GetPublic)

		// FAQ - Admin (Tabel terpisah dengan UUID)
		faqAdmin := v1.Group("/panel/faq")
		faqAdmin.Use(middleware.AuthMiddleware())
		faqAdmin.Use(middleware.AdminOnly())
		{
			faqAdmin.GET("", middleware.RequirePermission("operasional:read"), faqController.GetAll)
			faqAdmin.GET("/:id", middleware.RequirePermission("operasional:read"), faqController.GetByID)
			faqAdmin.POST("", middleware.RequirePermission("operasional:manage"), faqController.Create)
			faqAdmin.PUT("/:id", middleware.RequirePermission("operasional:manage"), faqController.Update)
			faqAdmin.DELETE("/:id", middleware.RequirePermission("operasional:manage"), faqController.Delete)
			faqAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("operasional:manage"), faqController.ToggleStatus)
			faqAdmin.PATCH("/:id/reorder", middleware.RequirePermission("operasional:manage"), faqController.Reorder)
		}

		// Disclaimer - Admin
		disclaimerAdmin := v1.Group("/panel/disclaimer")
		disclaimerAdmin.Use(middleware.AuthMiddleware())
		disclaimerAdmin.Use(middleware.AdminOnly())
		{
			disclaimerAdmin.GET("", middleware.RequirePermission("system:read"), disclaimerController.FindAll)
			disclaimerAdmin.GET("/:id", middleware.RequirePermission("system:read"), disclaimerController.FindByID)
			disclaimerAdmin.POST("", middleware.RequirePermission("system:manage"), disclaimerController.Create)
			disclaimerAdmin.PUT("/:id", middleware.RequirePermission("system:manage"), disclaimerController.Update)
			disclaimerAdmin.DELETE("/:id", middleware.RequirePermission("system:manage"), disclaimerController.Delete)
			disclaimerAdmin.PATCH("/:id/set-active", middleware.RequirePermission("system:manage"), disclaimerController.SetActive)
		}

		// Disclaimer - Public
		disclaimerPublic := v1.Group("/public/disclaimer")
		{
			disclaimerPublic.GET("", disclaimerController.GetActive)
		}

		// Formulir Partai Besar - Config (Admin)
		formulirConfigAdmin := v1.Group("/panel/formulir-partai-besar/config")
		formulirConfigAdmin.Use(middleware.AuthMiddleware())
		formulirConfigAdmin.Use(middleware.AdminOnly())
		{
			formulirConfigAdmin.GET("", middleware.RequirePermission("system:read"), formulirPartaiBesarController.GetConfig)
			formulirConfigAdmin.PUT("", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.UpdateConfig)
		}

		// Formulir Partai Besar - Anggaran (Admin)
		formulirAnggaranAdmin := v1.Group("/panel/formulir-partai-besar/anggaran")
		formulirAnggaranAdmin.Use(middleware.AuthMiddleware())
		formulirAnggaranAdmin.Use(middleware.AdminOnly())
		{
			formulirAnggaranAdmin.GET("", middleware.RequirePermission("system:read"), formulirPartaiBesarController.GetAnggaranList)
			formulirAnggaranAdmin.GET("/:id", middleware.RequirePermission("system:read"), formulirPartaiBesarController.GetAnggaranByID)
			formulirAnggaranAdmin.POST("", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.CreateAnggaran)
			formulirAnggaranAdmin.PUT("/:id", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.UpdateAnggaran)
			formulirAnggaranAdmin.DELETE("/:id", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.DeleteAnggaran)
			formulirAnggaranAdmin.PUT("/reorder", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.ReorderAnggaran)
			formulirAnggaranAdmin.PATCH("/:id/reorder", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.ReorderAnggaranByDirection)
		}

		// Formulir Partai Besar - Submission (Admin)
		formulirSubmissionAdmin := v1.Group("/panel/formulir-partai-besar/submission")
		formulirSubmissionAdmin.Use(middleware.AuthMiddleware())
		formulirSubmissionAdmin.Use(middleware.AdminOnly())
		{
			formulirSubmissionAdmin.GET("", middleware.RequirePermission("system:read"), formulirPartaiBesarController.GetSubmissionList)
			formulirSubmissionAdmin.GET("/:id", middleware.RequirePermission("system:read"), formulirPartaiBesarController.GetSubmissionDetail)
			formulirSubmissionAdmin.POST("/:id/resend-email", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.ResendEmail)
		}

		// Formulir Partai Besar - Buyer
		formulirBuyer := v1.Group("/buyer/formulir-partai-besar")
		formulirBuyer.Use(middleware.AuthMiddleware())
		formulirBuyer.Use(middleware.BuyerOnly())
		{
			formulirBuyer.GET("/options", formulirPartaiBesarController.GetOptions)
			formulirBuyer.POST("/submit", formulirPartaiBesarController.Submit)
		}

		// WhatsApp Handler - Admin (Simplified - Single Record)
		whatsappAdmin := v1.Group("/panel/whatsapp-handler")
		whatsappAdmin.Use(middleware.AuthMiddleware())
		whatsappAdmin.Use(middleware.AdminOnly())
		{
			whatsappAdmin.GET("", middleware.RequirePermission("system:read"), whatsappHandlerController.Get)
			whatsappAdmin.PUT("", middleware.RequirePermission("system:manage"), whatsappHandlerController.Update)
		}

		// WhatsApp Handler - Public
		whatsappPublic := v1.Group("/public/whatsapp-handler")
		{
			whatsappPublic.GET("", whatsappHandlerController.GetActive)
		}

		// Informasi Pickup - Public (Warehouse + Jadwal)
		v1.GET("/public/informasi-pickup", warehouseController.GetInformasiPickup)

		// Informasi Pickup - Admin (Jadwal via Warehouse)
		informasiPickupAdmin := v1.Group("/panel/informasi-pickup")
		informasiPickupAdmin.Use(middleware.AuthMiddleware())
		informasiPickupAdmin.Use(middleware.AdminOnly())
		{
			informasiPickupAdmin.GET("", middleware.RequirePermission("operasional:read"), warehouseController.Get)
			informasiPickupAdmin.GET("/jadwal", middleware.RequirePermission("operasional:read"), warehouseController.GetJadwal)
			informasiPickupAdmin.PUT("/jadwal", middleware.RequirePermission("operasional:manage"), warehouseController.UpdateJadwal)
		}

		// Blog - Admin
		blogAdmin := v1.Group("/panel/blog")
		blogAdmin.Use(middleware.AuthMiddleware())
		blogAdmin.Use(middleware.AdminOnly())
		{
			blogAdmin.GET("", middleware.RequirePermission("marketing:read"), blogController.GetAll)
			blogAdmin.GET("/statistik", middleware.RequirePermission("marketing:read"), blogController.GetStatistics)
			blogAdmin.GET("/dropdown", middleware.RequirePermission("marketing:read"), blogController.GetDropdownOptions)
			blogAdmin.GET("/:id", middleware.RequirePermission("marketing:read"), blogController.GetByID)
			blogAdmin.POST("", middleware.RequirePermission("marketing:manage"), blogController.Create)
			blogAdmin.PUT("/:id", middleware.RequirePermission("marketing:manage"), blogController.Update)
			blogAdmin.DELETE("/:id", middleware.RequirePermission("marketing:manage"), blogController.Delete)
			blogAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), blogController.ToggleStatus)
			blogAdmin.GET("/search", middleware.RequirePermission("marketing:read"), blogController.Search)
		}

		// Blog - Public
		blogPublic := v1.Group("/public/blog")
		{
			blogPublic.GET("", blogController.GetAll)
			blogPublic.GET("/populer", blogController.GetPopular)
			blogPublic.GET("/slug/:slug", blogController.GetBySlug)
			blogPublic.GET("/search", blogController.Search)
		}

		// Kategori Blog - Admin
		kategoriBlogAdmin := v1.Group("/panel/kategori-blog")
		kategoriBlogAdmin.Use(middleware.AuthMiddleware())
		kategoriBlogAdmin.Use(middleware.AdminOnly())
		{
			kategoriBlogAdmin.GET("", middleware.RequirePermission("marketing:read"), kategoriBlogController.GetAll)
			kategoriBlogAdmin.GET("/:id", middleware.RequirePermission("marketing:read"), kategoriBlogController.GetByID)
			kategoriBlogAdmin.POST("", middleware.RequirePermission("marketing:manage"), kategoriBlogController.Create)
			kategoriBlogAdmin.PUT("/:id", middleware.RequirePermission("marketing:manage"), kategoriBlogController.Update)
			kategoriBlogAdmin.DELETE("/:id", middleware.RequirePermission("marketing:manage"), kategoriBlogController.Delete)
			kategoriBlogAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), kategoriBlogController.ToggleStatus)
			kategoriBlogAdmin.PATCH("/:id/reorder", middleware.RequirePermission("marketing:manage"), kategoriBlogController.Reorder)
		}

		// Kategori Blog - Public
		v1.GET("/public/kategori-blog", kategoriBlogController.GetAllPublic)

		// Label Blog - Admin
		labelBlogAdmin := v1.Group("/panel/label-blog")
		labelBlogAdmin.Use(middleware.AuthMiddleware())
		labelBlogAdmin.Use(middleware.AdminOnly())
		{
			labelBlogAdmin.GET("", middleware.RequirePermission("marketing:read"), labelBlogController.GetAll)
			labelBlogAdmin.GET("/:id", middleware.RequirePermission("marketing:read"), labelBlogController.GetByID)
			labelBlogAdmin.POST("", middleware.RequirePermission("marketing:manage"), labelBlogController.Create)
			labelBlogAdmin.PUT("/:id", middleware.RequirePermission("marketing:manage"), labelBlogController.Update)
			labelBlogAdmin.DELETE("/:id", middleware.RequirePermission("marketing:manage"), labelBlogController.Delete)
			labelBlogAdmin.GET("/dropdown", middleware.RequirePermission("marketing:read"), labelBlogController.GetDropdownOptions)
			labelBlogAdmin.PATCH("/:id/reorder", middleware.RequirePermission("marketing:manage"), labelBlogController.Reorder)
		}

		// Label Blog - Public
		v1.GET("/public/label-blog", labelBlogController.GetAllPublic)

		// Video - Admin
		videoAdmin := v1.Group("/panel/video")
		videoAdmin.Use(middleware.AuthMiddleware())
		videoAdmin.Use(middleware.AdminOnly())
		{
			videoAdmin.GET("", middleware.RequirePermission("marketing:read"), videoController.GetAll)
			videoAdmin.GET("/statistik", middleware.RequirePermission("marketing:read"), videoController.GetStatistics)
			videoAdmin.GET("/dropdown", middleware.RequirePermission("marketing:read"), videoController.GetDropdownOptions)
			videoAdmin.GET("/:id", middleware.RequirePermission("marketing:read"), videoController.GetByID)
			videoAdmin.POST("", middleware.RequirePermission("marketing:manage"), videoController.Create)
			videoAdmin.PUT("/:id", middleware.RequirePermission("marketing:manage"), videoController.Update)
			videoAdmin.DELETE("/:id", middleware.RequirePermission("marketing:manage"), videoController.Delete)
			videoAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), videoController.ToggleStatus)
			videoAdmin.GET("/search", middleware.RequirePermission("marketing:read"), videoController.Search)
		}

		// Video - Public
		videoPublic := v1.Group("/public/video")
		{
			videoPublic.GET("", videoController.GetAll)
			videoPublic.GET("/slug/:slug", videoController.GetBySlug)
			videoPublic.GET("/popular", videoController.GetPopular)
			videoPublic.GET("/search", videoController.Search)
		}

		// Kategori Video - Admin
		kategoriVideoAdmin := v1.Group("/panel/kategori-video")
		kategoriVideoAdmin.Use(middleware.AuthMiddleware())
		kategoriVideoAdmin.Use(middleware.AdminOnly())
		{
			kategoriVideoAdmin.GET("", middleware.RequirePermission("marketing:read"), kategoriVideoController.GetAll)
			kategoriVideoAdmin.GET("/:id", middleware.RequirePermission("marketing:read"), kategoriVideoController.GetByID)
			kategoriVideoAdmin.POST("", middleware.RequirePermission("marketing:manage"), kategoriVideoController.Create)
			kategoriVideoAdmin.PUT("/:id", middleware.RequirePermission("marketing:manage"), kategoriVideoController.Update)
			kategoriVideoAdmin.DELETE("/:id", middleware.RequirePermission("marketing:manage"), kategoriVideoController.Delete)
			kategoriVideoAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), kategoriVideoController.ToggleStatus)
			kategoriVideoAdmin.PATCH("/:id/reorder", middleware.RequirePermission("marketing:manage"), kategoriVideoController.Reorder)
		}

		// Kategori Video - Public
		v1.GET("/public/kategori-video", kategoriVideoController.GetAllPublic)
	}

	// Endpoint to list all registered routes
	router.GET("/api/routes", func(c *gin.Context) {
		var endpointList []gin.H
		for _, route := range router.Routes() {
			endpointList = append(endpointList, gin.H{"method": route.Method, "path": route.Path})
		}
		c.JSON(http.StatusOK, endpointList)
	})
}
