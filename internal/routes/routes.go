package routes

import (
	"net/http"

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
) {
	// Health check
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "Server is running"})
	})

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

		// Kategori Produk - Admin (Write)
		kategoriAdmin := v1.Group("/panel/kategori-produk")
		kategoriAdmin.Use(middleware.AuthMiddleware())
		kategoriAdmin.Use(middleware.AdminOnly())
		kategoriAdmin.Use(middleware.RequirePermission("kategori:manage"))
		{
			kategoriAdmin.POST("", kategoriController.Create)
			kategoriAdmin.PUT("/:id", kategoriController.Update) // Support JSON & multipart/form-data
			kategoriAdmin.DELETE("/:id", kategoriController.Delete)
			kategoriAdmin.PATCH("/:id/toggle-status", kategoriController.ToggleStatus)
		}

		// Merek Produk - Public (Read Only)
		merekPublic := v1.Group("/merek-produk")
		{
			merekPublic.GET("", merekController.FindAll)
			merekPublic.GET("/:id", merekController.FindByID)
			merekPublic.GET("/slug/:slug", merekController.FindBySlug)
		}

		// Merek Produk - Admin (Write)
		merekAdmin := v1.Group("/panel/merek-produk")
		merekAdmin.Use(middleware.AuthMiddleware())
		merekAdmin.Use(middleware.AdminOnly())
		merekAdmin.Use(middleware.RequirePermission("brand:manage"))
		{
			merekAdmin.POST("", merekController.Create)
			merekAdmin.PUT("/:id", merekController.Update)
			merekAdmin.DELETE("/:id", merekController.Delete)
			merekAdmin.PATCH("/:id/toggle-status", merekController.ToggleStatus)
		}

		// Kondisi Produk - Public (Read Only)
		kondisiPublic := v1.Group("/kondisi-produk")
		{
			kondisiPublic.GET("", kondisiController.FindAll)
			kondisiPublic.GET("/:id", kondisiController.FindByID)
			kondisiPublic.GET("/slug/:slug", kondisiController.FindBySlug)
		}

		// Kondisi Produk - Admin (Write)
		kondisiAdmin := v1.Group("/panel/kondisi-produk")
		kondisiAdmin.Use(middleware.AuthMiddleware())
		kondisiAdmin.Use(middleware.AdminOnly())
		kondisiAdmin.Use(middleware.RequirePermission("kondisi:manage"))
		{
			kondisiAdmin.POST("", kondisiController.Create)
			kondisiAdmin.PUT("/:id", kondisiController.Update)
			kondisiAdmin.DELETE("/:id", kondisiController.Delete)
			kondisiAdmin.PATCH("/:id/toggle-status", kondisiController.ToggleStatus)
			kondisiAdmin.PUT("/reorder", kondisiController.Reorder)
		}

		// Kondisi Paket - Public (Read Only)
		paketPublic := v1.Group("/kondisi-paket")
		{
			paketPublic.GET("", kondisiPaketController.FindAll)
			paketPublic.GET("/:id", kondisiPaketController.FindByID)
			paketPublic.GET("/slug/:slug", kondisiPaketController.FindBySlug)
		}

		// Kondisi Paket - Admin (Write)
		paketAdmin := v1.Group("/panel/kondisi-paket")
		paketAdmin.Use(middleware.AuthMiddleware())
		paketAdmin.Use(middleware.AdminOnly())
		paketAdmin.Use(middleware.RequirePermission("kondisi:manage"))
		{
			paketAdmin.POST("", kondisiPaketController.Create)
			paketAdmin.PUT("/:id", kondisiPaketController.Update)
			paketAdmin.DELETE("/:id", kondisiPaketController.Delete)
			paketAdmin.PATCH("/:id/toggle-status", kondisiPaketController.ToggleStatus)
			paketAdmin.PUT("/reorder", kondisiPaketController.Reorder)
		}

		// Sumber Produk - Public (Read Only)
		sumberPublic := v1.Group("/sumber-produk")
		{
			sumberPublic.GET("", sumberController.FindAll)
			sumberPublic.GET("/:id", sumberController.FindByID)
			sumberPublic.GET("/slug/:slug", sumberController.FindBySlug)
		}

		// Sumber Produk - Admin (Write)
		sumberAdmin := v1.Group("/panel/sumber-produk")
		sumberAdmin.Use(middleware.AuthMiddleware())
		sumberAdmin.Use(middleware.AdminOnly())
		sumberAdmin.Use(middleware.RequirePermission("kondisi:manage"))
		{
			sumberAdmin.POST("", sumberController.Create)
			sumberAdmin.PUT("/:id", sumberController.Update)
			sumberAdmin.DELETE("/:id", sumberController.Delete)
			sumberAdmin.PATCH("/:id/toggle-status", sumberController.ToggleStatus)
		}

		// Warehouse - Public (Read Only)
		warehousePublic := v1.Group("/warehouse")
		{
			warehousePublic.GET("", warehouseController.FindAll)
			warehousePublic.GET("/:id", warehouseController.FindByID)
		}

		// Warehouse - Admin (Write)
		warehouseAdmin := v1.Group("/panel/warehouse")
		warehouseAdmin.Use(middleware.AuthMiddleware())
		warehouseAdmin.Use(middleware.AdminOnly())
		warehouseAdmin.Use(middleware.RequirePermission("kondisi:manage"))
		{
			warehouseAdmin.POST("", warehouseController.Create)
			warehouseAdmin.PUT("/:id", warehouseController.Update)
			warehouseAdmin.DELETE("/:id", warehouseController.Delete)
			warehouseAdmin.PATCH("/:id/toggle-status", warehouseController.ToggleStatus)
		}

		// Tipe Produk - Public (Read Only)
		// Note: Tipe produk is read-only (Paletbox, Container, Truckload) - managed via migration only
		tipeProdukPublic := v1.Group("/tipe-produk")
		{
			tipeProdukPublic.GET("", tipeProdukController.FindAll)
			tipeProdukPublic.GET("/:id", tipeProdukController.FindByID)
			tipeProdukPublic.GET("/slug/:slug", tipeProdukController.FindBySlug)
		}

		// Diskon Kategori - Public (Read Only)
		diskonKategoriPublic := v1.Group("/diskon-kategori")
		{
			diskonKategoriPublic.GET("", diskonKategoriController.FindAll)
			diskonKategoriPublic.GET("/:id", diskonKategoriController.FindByID)
			diskonKategoriPublic.GET("/by-kategori/:kategori_id", diskonKategoriController.FindActiveByKategoriID)
		}

		// Diskon Kategori - Admin (Write)
		diskonKategoriAdmin := v1.Group("/panel/diskon-kategori")
		diskonKategoriAdmin.Use(middleware.AuthMiddleware())
		diskonKategoriAdmin.Use(middleware.AdminOnly())
		diskonKategoriAdmin.Use(middleware.RequirePermission("diskon:manage"))
		{
			diskonKategoriAdmin.POST("", diskonKategoriController.Create)
			diskonKategoriAdmin.PUT("/:id", diskonKategoriController.Update)
			diskonKategoriAdmin.DELETE("/:id", diskonKategoriController.Delete)
			diskonKategoriAdmin.PATCH("/:id/toggle-status", diskonKategoriController.ToggleStatus)
		}

		// Banner Tipe Produk - Public (Read Only)
		bannerTipeProdukPublic := v1.Group("/banner-tipe-produk")
		{
			bannerTipeProdukPublic.GET("", bannerTipeProdukController.FindAll)
			bannerTipeProdukPublic.GET("/:id", bannerTipeProdukController.FindByID)
			bannerTipeProdukPublic.GET("/by-tipe/:tipe_produk_id", bannerTipeProdukController.FindByTipeProdukID)
		}

		// Banner Tipe Produk - Admin (Write)
		bannerTipeProdukAdmin := v1.Group("/panel/banner-tipe-produk")
		bannerTipeProdukAdmin.Use(middleware.AuthMiddleware())
		bannerTipeProdukAdmin.Use(middleware.AdminOnly())
		bannerTipeProdukAdmin.Use(middleware.RequirePermission("banner:manage"))
		{
			bannerTipeProdukAdmin.POST("", bannerTipeProdukController.Create)
			bannerTipeProdukAdmin.PUT("/:id", bannerTipeProdukController.Update)
			bannerTipeProdukAdmin.DELETE("/:id", bannerTipeProdukController.Delete)
			bannerTipeProdukAdmin.PATCH("/:id/toggle-status", bannerTipeProdukController.ToggleStatus)
			bannerTipeProdukAdmin.PUT("/reorder", bannerTipeProdukController.Reorder)
		}

		// Produk - Public (Read Only)
		produkPublic := v1.Group("/produk")
		{
			produkPublic.GET("", produkController.FindAll)
			produkPublic.GET("/:id", produkController.FindByID)
			produkPublic.GET("/slug/:slug", produkController.FindBySlug)
		}

		// Produk - Admin (Write)
		produkAdmin := v1.Group("/panel/produk")
		produkAdmin.Use(middleware.AuthMiddleware())
		produkAdmin.Use(middleware.AdminOnly())
		{
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
			produkAdmin.PUT("/:id/gambar/:gambar_id", middleware.RequirePermission("produk:update"), produkController.UpdateGambar)
			produkAdmin.DELETE("/:id/gambar/:gambar_id", middleware.RequirePermission("produk:update"), produkController.DeleteGambar)
			produkAdmin.PUT("/:id/gambar/reorder", middleware.RequirePermission("produk:update"), produkController.ReorderGambar)

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
			heroSectionAdmin.GET("", middleware.RequirePermission("hero_section:read"), heroSectionController.FindAll)
			heroSectionAdmin.GET("/:id", middleware.RequirePermission("hero_section:read"), heroSectionController.FindByID)
			heroSectionAdmin.POST("", middleware.RequirePermission("hero_section:manage"), heroSectionController.Create)
			heroSectionAdmin.PUT("/:id", middleware.RequirePermission("hero_section:manage"), heroSectionController.Update)
			heroSectionAdmin.DELETE("/:id", middleware.RequirePermission("hero_section:manage"), heroSectionController.Delete)
			heroSectionAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("hero_section:manage"), heroSectionController.ToggleStatus)
			heroSectionAdmin.PUT("/reorder", middleware.RequirePermission("hero_section:manage"), heroSectionController.Reorder)
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
			bannerEventPromoAdmin.GET("", middleware.RequirePermission("banner:read"), bannerEventPromoController.FindAll)
			bannerEventPromoAdmin.GET("/:id", middleware.RequirePermission("banner:read"), bannerEventPromoController.FindByID)
			bannerEventPromoAdmin.POST("", middleware.RequirePermission("banner:manage"), bannerEventPromoController.Create)
			bannerEventPromoAdmin.PUT("/:id", middleware.RequirePermission("banner:manage"), bannerEventPromoController.Update)
			bannerEventPromoAdmin.DELETE("/:id", middleware.RequirePermission("banner:manage"), bannerEventPromoController.Delete)
			bannerEventPromoAdmin.PATCH("/:id/toggle-status", middleware.RequirePermission("banner:manage"), bannerEventPromoController.ToggleStatus)
			bannerEventPromoAdmin.PUT("/reorder", middleware.RequirePermission("banner:manage"), bannerEventPromoController.Reorder)
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
