package routes

import (
	"os"

	"project-bulky-be/internal/controllers"
	"project-bulky-be/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	router *fiber.App,
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
	ulasanAdminController *controllers.UlasanAdminController,
	pesananAdminController *controllers.PesananAdminController,
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
	kuponController *controllers.KuponController,
) {
	// Health check
	router.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK", "message": "Server is running"})
	})

	// Serve static files from uploads folder
	uploadPath := os.Getenv("UPLOAD_PATH")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}
	router.Static("/uploads", uploadPath)

	// API routes
	v1 := router.Group("/api")

	// Admin Management Routes (Protected)
	admin := v1.Group("/panel/admin",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		middleware.RequirePermission("admin:manage"),
	)
	admin.Get("", adminController.FindAll)
	admin.Get("/:id", adminController.FindByID)
	admin.Post("", adminController.Create)
	admin.Put("/:id", adminController.Update)
	admin.Delete("/:id", adminController.Delete)
	admin.Patch("/:id/toggle-status", adminController.ToggleStatus)
	admin.Put("/:id/reset-password", adminController.ResetPassword)

	// Buyer Management Routes (Admin Side)
	buyerManagement := v1.Group("/panel/buyer",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	buyerManagement.Get("", middleware.RequirePermission("buyer:read"), buyerController.FindAll)
	buyerManagement.Get("/statistik", middleware.RequirePermission("buyer:read"), buyerController.GetStatistik)
	buyerManagement.Get("/chart", middleware.RequirePermission("buyer:read"), buyerController.GetChart)
	buyerManagement.Get("/:id", middleware.RequirePermission("buyer:read"), buyerController.FindByID)
	buyerManagement.Delete("/:id", middleware.RequirePermission("buyer:manage"), buyerController.Delete)
	buyerManagement.Put("/:id/reset-password", middleware.RequirePermission("buyer:manage"), buyerController.ResetPassword)

	// Alamat Buyer Routes (Buyer Only)
	alamatBuyer := v1.Group("/buyer/alamat",
		middleware.AuthMiddleware(),
		middleware.BuyerOnly(),
	)
	alamatBuyer.Get("", alamatBuyerController.FindAll)
	alamatBuyer.Get("/:id", alamatBuyerController.FindByID)
	alamatBuyer.Post("", alamatBuyerController.Create)
	alamatBuyer.Put("/:id", alamatBuyerController.Update)
	alamatBuyer.Delete("/:id", alamatBuyerController.Delete)
	alamatBuyer.Patch("/:id/set-default", alamatBuyerController.SetDefault)

	// Kategori Produk - Public
	kategoriPublic := v1.Group("/kategori-produk")
	kategoriPublic.Get("", kategoriController.FindAll)
	kategoriPublic.Get("/:id", kategoriController.FindByID)
	kategoriPublic.Get("/slug/:slug", kategoriController.FindBySlug)

	// Kategori Produk - Admin
	kategoriAdmin := v1.Group("/panel/kategori-produk",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	kategoriAdmin.Get("", middleware.RequirePermission("kategori:read"), kategoriController.FindAll)
	kategoriAdmin.Get("/dropdown", middleware.RequirePermission("kategori:read"), kategoriController.Dropdown)
	kategoriAdmin.Get("/:id", middleware.RequirePermission("kategori:read"), kategoriController.FindByID)
	kategoriAdmin.Post("", middleware.RequirePermission("kategori:manage"), kategoriController.Create)
	kategoriAdmin.Put("/:id", middleware.RequirePermission("kategori:manage"), kategoriController.Update)
	kategoriAdmin.Delete("/:id", middleware.RequirePermission("kategori:manage"), kategoriController.Delete)
	kategoriAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("kategori:manage"), kategoriController.ToggleStatus)

	// Merek Produk - Public
	merekPublic := v1.Group("/merek-produk")
	merekPublic.Get("", merekController.FindAll)
	merekPublic.Get("/:id", merekController.FindByID)
	merekPublic.Get("/slug/:slug", merekController.FindBySlug)

	// Merek Produk - Admin
	merekAdmin := v1.Group("/panel/merek-produk",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	merekAdmin.Get("", middleware.RequirePermission("brand:read"), merekController.FindAll)
	merekAdmin.Get("/dropdown", middleware.RequirePermission("brand:read"), merekController.Dropdown)
	merekAdmin.Get("/:id", middleware.RequirePermission("brand:read"), merekController.FindByID)
	merekAdmin.Post("", middleware.RequirePermission("brand:manage"), merekController.Create)
	merekAdmin.Put("/:id", middleware.RequirePermission("brand:manage"), merekController.Update)
	merekAdmin.Delete("/:id", middleware.RequirePermission("brand:manage"), merekController.Delete)
	merekAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("brand:manage"), merekController.ToggleStatus)

	// Kondisi Produk - Public
	kondisiPublic := v1.Group("/kondisi-produk")
	kondisiPublic.Get("", kondisiController.FindAll)
	kondisiPublic.Get("/:id", kondisiController.FindByID)
	kondisiPublic.Get("/slug/:slug", kondisiController.FindBySlug)

	// Kondisi Produk - Admin
	kondisiAdmin := v1.Group("/panel/kondisi-produk",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	kondisiAdmin.Get("", middleware.RequirePermission("kondisi:read"), kondisiController.FindAll)
	kondisiAdmin.Get("/dropdown", middleware.RequirePermission("kondisi:read"), kondisiController.Dropdown)
	kondisiAdmin.Get("/:id", middleware.RequirePermission("kondisi:read"), kondisiController.FindByID)
	kondisiAdmin.Post("", middleware.RequirePermission("kondisi:manage"), kondisiController.Create)
	kondisiAdmin.Put("/:id", middleware.RequirePermission("kondisi:manage"), kondisiController.Update)
	kondisiAdmin.Delete("/:id", middleware.RequirePermission("kondisi:manage"), kondisiController.Delete)
	kondisiAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("kondisi:manage"), kondisiController.ToggleStatus)
	kondisiAdmin.Put("/reorder", middleware.RequirePermission("kondisi:manage"), kondisiController.Reorder)
	kondisiAdmin.Patch("/:id/reorder", middleware.RequirePermission("kondisi:manage"), kondisiController.ReorderByDirection)

	// Kondisi Paket - Public
	paketPublic := v1.Group("/kondisi-paket")
	paketPublic.Get("", kondisiPaketController.FindAll)
	paketPublic.Get("/:id", kondisiPaketController.FindByID)
	paketPublic.Get("/slug/:slug", kondisiPaketController.FindBySlug)

	// Kondisi Paket - Admin
	paketAdmin := v1.Group("/panel/kondisi-paket",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	paketAdmin.Get("", middleware.RequirePermission("kondisi:read"), kondisiPaketController.FindAll)
	paketAdmin.Get("/dropdown", middleware.RequirePermission("kondisi:read"), kondisiPaketController.Dropdown)
	paketAdmin.Get("/:id", middleware.RequirePermission("kondisi:read"), kondisiPaketController.FindByID)
	paketAdmin.Post("", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.Create)
	paketAdmin.Put("/:id", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.Update)
	paketAdmin.Delete("/:id", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.Delete)
	paketAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.ToggleStatus)
	paketAdmin.Put("/reorder", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.Reorder)
	paketAdmin.Patch("/:id/reorder", middleware.RequirePermission("kondisi:manage"), kondisiPaketController.ReorderByDirection)

	// Sumber Produk - Public
	sumberPublic := v1.Group("/sumber-produk")
	sumberPublic.Get("", sumberController.FindAll)
	sumberPublic.Get("/:id", sumberController.FindByID)
	sumberPublic.Get("/slug/:slug", sumberController.FindBySlug)

	// Sumber Produk - Admin
	sumberAdmin := v1.Group("/panel/sumber-produk",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	sumberAdmin.Get("", middleware.RequirePermission("kondisi:read"), sumberController.FindAll)
	sumberAdmin.Get("/dropdown", middleware.RequirePermission("kondisi:read"), sumberController.Dropdown)
	sumberAdmin.Get("/:id", middleware.RequirePermission("kondisi:read"), sumberController.FindByID)
	sumberAdmin.Post("", middleware.RequirePermission("kondisi:manage"), sumberController.Create)
	sumberAdmin.Put("/:id", middleware.RequirePermission("kondisi:manage"), sumberController.Update)
	sumberAdmin.Delete("/:id", middleware.RequirePermission("kondisi:manage"), sumberController.Delete)
	sumberAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("kondisi:manage"), sumberController.ToggleStatus)

	// Warehouse - Public
	v1.Get("/public/warehouse", warehouseController.GetPublic)

	// Warehouse - Admin
	warehouseAdmin := v1.Group("/panel/warehouse",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	warehouseAdmin.Get("", middleware.RequirePermission("operasional:read"), warehouseController.Get)
	warehouseAdmin.Get("/dropdown", middleware.RequirePermission("operasional:read"), warehouseController.Dropdown)
	warehouseAdmin.Put("", middleware.RequirePermission("operasional:manage"), warehouseController.UpdateSingleton)

	// Tipe Produk - Public
	tipeProdukPublic := v1.Group("/tipe-produk")
	tipeProdukPublic.Get("", tipeProdukController.FindAll)
	tipeProdukPublic.Get("/with-produk", tipeProdukController.FindAllWithProduk)

	// Tipe Produk - Admin
	tipeProdukAdmin := v1.Group("/panel/tipe-produk",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	tipeProdukAdmin.Get("", middleware.RequirePermission("system:read"), tipeProdukController.FindAll)
	tipeProdukAdmin.Get("/dropdown", middleware.RequirePermission("system:read"), tipeProdukController.Dropdown)
	tipeProdukAdmin.Get("/with-produk", middleware.RequirePermission("system:read"), tipeProdukController.FindAllWithProduk)

	// Diskon Kategori - Public
	diskonKategoriPublic := v1.Group("/diskon-kategori")
	diskonKategoriPublic.Get("", diskonKategoriController.FindAll)
	diskonKategoriPublic.Get("/:id", diskonKategoriController.FindByID)
	diskonKategoriPublic.Get("/by-kategori/:kategori_id", diskonKategoriController.FindActiveByKategoriID)

	// Diskon Kategori - Admin
	diskonKategoriAdmin := v1.Group("/panel/diskon-kategori",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	diskonKategoriAdmin.Get("", middleware.RequirePermission("diskon:read"), diskonKategoriController.FindAll)
	diskonKategoriAdmin.Get("/:id", middleware.RequirePermission("diskon:read"), diskonKategoriController.FindByID)
	diskonKategoriAdmin.Post("", middleware.RequirePermission("diskon:manage"), diskonKategoriController.Create)
	diskonKategoriAdmin.Put("/:id", middleware.RequirePermission("diskon:manage"), diskonKategoriController.Update)
	diskonKategoriAdmin.Delete("/:id", middleware.RequirePermission("diskon:manage"), diskonKategoriController.Delete)
	diskonKategoriAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("diskon:manage"), diskonKategoriController.ToggleStatus)

	// Kupon - Admin
	kuponAdmin := v1.Group("/panel/kupon",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	kuponAdmin.Get("", middleware.RequirePermission("kupon:read"), kuponController.GetAll)
	kuponAdmin.Get("/dropdown/kategori", middleware.RequirePermission("kupon:read"), kuponController.GetKategoriDropdown)
	kuponAdmin.Get("/:id", middleware.RequirePermission("kupon:read"), kuponController.GetByID)
	kuponAdmin.Get("/:id/usages", middleware.RequirePermission("kupon:read"), kuponController.GetUsages)
	kuponAdmin.Post("", middleware.RequirePermission("kupon:manage"), kuponController.Create)
	kuponAdmin.Post("/generate-kode", middleware.RequirePermission("kupon:manage"), kuponController.GenerateKode)
	kuponAdmin.Put("/:id", middleware.RequirePermission("kupon:manage"), kuponController.Update)
	kuponAdmin.Delete("/:id", middleware.RequirePermission("kupon:manage"), kuponController.Delete)
	kuponAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("kupon:manage"), kuponController.ToggleStatus)

	// Banner Tipe Produk - Public
	bannerTipeProdukPublic := v1.Group("/banner-tipe-produk")
	bannerTipeProdukPublic.Get("", bannerTipeProdukController.FindAll)
	bannerTipeProdukPublic.Get("/:id", bannerTipeProdukController.FindByID)
	bannerTipeProdukPublic.Get("/by-tipe/:tipe_produk_id", bannerTipeProdukController.FindByTipeProdukID)

	// Banner Tipe Produk - Admin
	bannerTipeProdukAdmin := v1.Group("/panel/banner-tipe-produk",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	bannerTipeProdukAdmin.Get("", middleware.RequirePermission("marketing:read"), bannerTipeProdukController.FindAll)
	bannerTipeProdukAdmin.Get("/:id", middleware.RequirePermission("marketing:read"), bannerTipeProdukController.FindByID)
	bannerTipeProdukAdmin.Post("", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.Create)
	bannerTipeProdukAdmin.Put("/:id", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.Update)
	bannerTipeProdukAdmin.Delete("/:id", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.Delete)
	bannerTipeProdukAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.ToggleStatus)
	bannerTipeProdukAdmin.Put("/reorder", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.Reorder)
	bannerTipeProdukAdmin.Patch("/:id/reorder", middleware.RequirePermission("marketing:manage"), bannerTipeProdukController.ReorderByDirection)

	// Produk - Public
	produkPublic := v1.Group("/produk")
	produkPublic.Get("", produkController.FindAll)
	produkPublic.Get("/:id", produkController.FindByID)
	produkPublic.Get("/slug/:slug", produkController.FindBySlug)

	// Produk - Admin
	produkAdmin := v1.Group("/panel/produk",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	produkAdmin.Get("", middleware.RequirePermission("produk:read"), produkController.FindAll)
	produkAdmin.Get("/:id", middleware.RequirePermission("produk:read"), produkController.FindByID)
	produkAdmin.Post("", middleware.RequirePermission("produk:create"), produkController.Create)
	produkAdmin.Put("/:id", middleware.RequirePermission("produk:update"), produkController.Update)
	produkAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("produk:update"), produkController.ToggleStatus)
	produkAdmin.Patch("/:id/update-stock", middleware.RequirePermission("produk:update"), produkController.UpdateStock)
	produkAdmin.Delete("/:id", middleware.RequirePermission("produk:delete"), produkController.Delete)
	produkAdmin.Post("/:id/gambar", middleware.RequirePermission("produk:update"), produkController.AddGambar)
	produkAdmin.Delete("/:id/gambar/:gambar_id", middleware.RequirePermission("produk:update"), produkController.DeleteGambar)
	produkAdmin.Patch("/:id/gambar/:gambar_id/reorder", middleware.RequirePermission("produk:update"), produkController.ReorderGambar)
	produkAdmin.Post("/:id/dokumen", middleware.RequirePermission("produk:update"), produkController.AddDokumen)
	produkAdmin.Delete("/:id/dokumen/:dokumen_id", middleware.RequirePermission("produk:update"), produkController.DeleteDokumen)

	// Master (Dropdown)
	v1.Get("/master/dropdown", masterController.GetDropdown)

	// Hero Section - Admin
	heroSectionAdmin := v1.Group("/panel/hero-section",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	heroSectionAdmin.Get("", middleware.RequirePermission("marketing:read"), heroSectionController.FindAll)
	heroSectionAdmin.Get("/schedule", middleware.RequirePermission("marketing:read"), heroSectionController.GetSchedules)
	heroSectionAdmin.Get("/:id", middleware.RequirePermission("marketing:read"), heroSectionController.FindByID)
	heroSectionAdmin.Post("", middleware.RequirePermission("marketing:manage"), heroSectionController.Create)
	heroSectionAdmin.Put("/:id", middleware.RequirePermission("marketing:manage"), heroSectionController.Update)
	heroSectionAdmin.Delete("/:id", middleware.RequirePermission("marketing:manage"), heroSectionController.Delete)
	heroSectionAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), heroSectionController.ToggleStatus)

	// Hero Section - Public
	v1.Get("/hero-section/active", heroSectionController.GetActive)

	// Banner Event Promo - Admin
	bannerEventPromoAdmin := v1.Group("/panel/banner-event-promo",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	bannerEventPromoAdmin.Get("", middleware.RequirePermission("marketing:read"), bannerEventPromoController.FindAll)
	bannerEventPromoAdmin.Get("/:id", middleware.RequirePermission("marketing:read"), bannerEventPromoController.FindByID)
	bannerEventPromoAdmin.Post("", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.Create)
	bannerEventPromoAdmin.Put("/:id", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.Update)
	bannerEventPromoAdmin.Delete("/:id", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.Delete)
	bannerEventPromoAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.ToggleStatus)
	bannerEventPromoAdmin.Put("/reorder", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.Reorder)
	bannerEventPromoAdmin.Patch("/:id/reorder", middleware.RequirePermission("marketing:manage"), bannerEventPromoController.ReorderByDirection)

	// Banner Event Promo - Public
	v1.Get("/banner-event-promo/active", bannerEventPromoController.GetActive)

	// Ulasan - Admin
	ulasanAdmin := v1.Group("/panel/ulasan",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	ulasanAdmin.Get("", middleware.RequirePermission("ulasan:read"), ulasanAdminController.GetAll)
	ulasanAdmin.Get("/:id", middleware.RequirePermission("ulasan:read"), ulasanAdminController.GetByID)
	ulasanAdmin.Patch("/:id/approve", middleware.RequirePermission("ulasan:manage"), ulasanAdminController.Approve)
	ulasanAdmin.Patch("/:id/reject", middleware.RequirePermission("ulasan:manage"), ulasanAdminController.Reject)
	ulasanAdmin.Patch("/bulk-approve", middleware.RequirePermission("ulasan:manage"), ulasanAdminController.BulkApprove)
	ulasanAdmin.Delete("/:id", middleware.RequirePermission("ulasan:manage"), ulasanAdminController.Delete)

	// Pesanan - Admin
	pesananAdmin := v1.Group("/panel/pesanan",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	pesananAdmin.Get("", middleware.RequirePermission("pesanan:read"), pesananAdminController.GetAll)
	pesananAdmin.Get("/statistics", middleware.RequirePermission("pesanan:read"), pesananAdminController.GetStatistics)
	pesananAdmin.Get("/:id", middleware.RequirePermission("pesanan:read"), pesananAdminController.GetByID)
	pesananAdmin.Patch("/:id/update-status", middleware.RequirePermission("pesanan:update_status"), pesananAdminController.UpdateStatus)
	pesananAdmin.Delete("/:id", middleware.RequirePermission("pesanan:delete"), pesananAdminController.Delete)

	// Ulasan - Buyer
	ulasanBuyer := v1.Group("/buyer/ulasan",
		middleware.AuthMiddleware(),
		middleware.BuyerOnly(),
	)
	ulasanBuyer.Get("/pending", ulasanController.GetPendingReviews)
	ulasanBuyer.Get("", ulasanController.BuyerFindAll)
	ulasanBuyer.Post("", ulasanController.Create)

	// Ulasan - Public
	ulasanPublic := v1.Group("/public/produk")
	ulasanPublic.Get("/:produk_id/ulasan", ulasanController.GetProdukUlasan)
	ulasanPublic.Get("/:produk_id/rating", ulasanController.GetProdukRating)

	// Force Update - Super Admin Only
	forceUpdateAdmin := v1.Group("/panel/force-update",
		middleware.AuthMiddleware(),
		middleware.SuperAdminOnly(),
	)
	forceUpdateAdmin.Get("", forceUpdateController.GetAllForceUpdates)
	forceUpdateAdmin.Get("/:id", forceUpdateController.GetForceUpdateByID)
	forceUpdateAdmin.Post("", forceUpdateController.CreateForceUpdate)
	forceUpdateAdmin.Put("/:id", forceUpdateController.UpdateForceUpdate)
	forceUpdateAdmin.Delete("/:id", forceUpdateController.DeleteForceUpdate)
	forceUpdateAdmin.Post("/:id/set-active", forceUpdateController.SetActiveForceUpdate)

	// Mode Maintenance - Super Admin Only
	modeMaintenanceAdmin := v1.Group("/panel/mode-maintenance",
		middleware.AuthMiddleware(),
		middleware.SuperAdminOnly(),
	)
	modeMaintenanceAdmin.Get("", modeMaintenanceController.GetAllMaintenances)
	modeMaintenanceAdmin.Get("/:id", modeMaintenanceController.GetMaintenanceByID)
	modeMaintenanceAdmin.Post("", modeMaintenanceController.CreateMaintenance)
	modeMaintenanceAdmin.Put("/:id", modeMaintenanceController.UpdateMaintenance)
	modeMaintenanceAdmin.Delete("/:id", modeMaintenanceController.DeleteMaintenance)
	modeMaintenanceAdmin.Post("/:id/activate", modeMaintenanceController.ActivateMaintenance)
	modeMaintenanceAdmin.Post("/:id/deactivate", modeMaintenanceController.DeactivateMaintenance)

	// App Status - Public
	appStatus := v1.Group("/public")
	appStatus.Get("/check-version", appStatusController.CheckVersion)
	appStatus.Get("/check-maintenance", appStatusController.CheckMaintenance)
	appStatus.Get("/app-status", appStatusController.AppStatus)

	// PPN - Admin
	ppnAdmin := v1.Group("/panel/ppn",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	ppnAdmin.Get("", middleware.RequirePermission("system:read"), ppnController.GetAll)
	ppnAdmin.Get("/:id", middleware.RequirePermission("system:read"), ppnController.GetByID)
	ppnAdmin.Post("", middleware.RequirePermission("system:manage"), ppnController.Create)
	ppnAdmin.Put("/:id", middleware.RequirePermission("system:manage"), ppnController.Update)
	ppnAdmin.Delete("/:id", middleware.RequirePermission("system:manage"), ppnController.Delete)
	ppnAdmin.Patch("/:id/set-active", middleware.RequirePermission("system:manage"), ppnController.SetActive)

	// Metode Pembayaran - Admin
	metodePembayaranAdmin := v1.Group("/panel/metode-pembayaran",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	metodePembayaranAdmin.Get("", middleware.RequirePermission("pembayaran:read"), metodePembayaranController.GetAllGrouped)
	metodePembayaranAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("pembayaran:manage"), metodePembayaranController.ToggleMethodStatus)
	metodePembayaranAdmin.Patch("/group/:urutan/toggle-status", middleware.RequirePermission("pembayaran:manage"), metodePembayaranController.ToggleGroupStatus)

	// Metode Pembayaran - Public
	v1.Get("/public/metode-pembayaran", metodePembayaranController.GetAllGroupedPublic)

	// Dokumen Kebijakan - Admin
	dokumenKebijakanAdmin := v1.Group("/panel/dokumen-kebijakan",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	dokumenKebijakanAdmin.Get("", middleware.RequirePermission("system:read"), dokumenKebijakanController.GetAll)
	dokumenKebijakanAdmin.Get("/:id", middleware.RequirePermission("system:read"), dokumenKebijakanController.GetByID)
	dokumenKebijakanAdmin.Put("/:id", middleware.RequirePermission("system:manage"), dokumenKebijakanController.Update)

	// Dokumen Kebijakan - Public
	dokumenKebijakanPublic := v1.Group("/public/dokumen-kebijakan")
	dokumenKebijakanPublic.Get("", dokumenKebijakanController.GetAllPublic)
	dokumenKebijakanPublic.Get("/:id", dokumenKebijakanController.GetByIDPublic)

	// FAQ - Public
	v1.Get("/public/faq", faqController.GetPublic)

	// FAQ - Admin
	faqAdmin := v1.Group("/panel/faq",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	faqAdmin.Get("", middleware.RequirePermission("operasional:read"), faqController.GetAll)
	faqAdmin.Get("/:id", middleware.RequirePermission("operasional:read"), faqController.GetByID)
	faqAdmin.Post("", middleware.RequirePermission("operasional:manage"), faqController.Create)
	faqAdmin.Put("/:id", middleware.RequirePermission("operasional:manage"), faqController.Update)
	faqAdmin.Delete("/:id", middleware.RequirePermission("operasional:manage"), faqController.Delete)
	faqAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("operasional:manage"), faqController.ToggleStatus)
	faqAdmin.Patch("/:id/reorder", middleware.RequirePermission("operasional:manage"), faqController.Reorder)

	// Disclaimer - Admin
	disclaimerAdmin := v1.Group("/panel/disclaimer",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	disclaimerAdmin.Get("", middleware.RequirePermission("system:read"), disclaimerController.FindAll)
	disclaimerAdmin.Get("/:id", middleware.RequirePermission("system:read"), disclaimerController.FindByID)
	disclaimerAdmin.Post("", middleware.RequirePermission("system:manage"), disclaimerController.Create)
	disclaimerAdmin.Put("/:id", middleware.RequirePermission("system:manage"), disclaimerController.Update)
	disclaimerAdmin.Delete("/:id", middleware.RequirePermission("system:manage"), disclaimerController.Delete)
	disclaimerAdmin.Patch("/:id/set-active", middleware.RequirePermission("system:manage"), disclaimerController.SetActive)

	// Disclaimer - Public
	v1.Get("/public/disclaimer", disclaimerController.GetActive)

	// Formulir Partai Besar - Config (Admin)
	formulirConfigAdmin := v1.Group("/panel/formulir-partai-besar/config",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	formulirConfigAdmin.Get("", middleware.RequirePermission("system:read"), formulirPartaiBesarController.GetConfig)
	formulirConfigAdmin.Put("", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.UpdateConfig)

	// Formulir Partai Besar - Anggaran (Admin)
	formulirAnggaranAdmin := v1.Group("/panel/formulir-partai-besar/anggaran",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	formulirAnggaranAdmin.Get("", middleware.RequirePermission("system:read"), formulirPartaiBesarController.GetAnggaranList)
	formulirAnggaranAdmin.Get("/:id", middleware.RequirePermission("system:read"), formulirPartaiBesarController.GetAnggaranByID)
	formulirAnggaranAdmin.Post("", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.CreateAnggaran)
	formulirAnggaranAdmin.Put("/:id", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.UpdateAnggaran)
	formulirAnggaranAdmin.Delete("/:id", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.DeleteAnggaran)
	formulirAnggaranAdmin.Put("/reorder", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.ReorderAnggaran)
	formulirAnggaranAdmin.Patch("/:id/reorder", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.ReorderAnggaranByDirection)

	// Formulir Partai Besar - Submission (Admin)
	formulirSubmissionAdmin := v1.Group("/panel/formulir-partai-besar/submission",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	formulirSubmissionAdmin.Get("", middleware.RequirePermission("system:read"), formulirPartaiBesarController.GetSubmissionList)
	formulirSubmissionAdmin.Get("/:id", middleware.RequirePermission("system:read"), formulirPartaiBesarController.GetSubmissionDetail)
	formulirSubmissionAdmin.Post("/:id/resend-email", middleware.RequirePermission("system:manage"), formulirPartaiBesarController.ResendEmail)

	// Formulir Partai Besar - Buyer
	formulirBuyer := v1.Group("/buyer/formulir-partai-besar",
		middleware.AuthMiddleware(),
		middleware.BuyerOnly(),
	)
	formulirBuyer.Get("/options", formulirPartaiBesarController.GetOptions)
	formulirBuyer.Post("/submit", formulirPartaiBesarController.Submit)

	// WhatsApp Handler - Admin
	whatsappAdmin := v1.Group("/panel/whatsapp-handler",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	whatsappAdmin.Get("", middleware.RequirePermission("system:read"), whatsappHandlerController.Get)
	whatsappAdmin.Put("", middleware.RequirePermission("system:manage"), whatsappHandlerController.Update)

	// WhatsApp Handler - Public
	v1.Get("/public/whatsapp-handler", whatsappHandlerController.GetActive)

	// Informasi Pickup - Public
	v1.Get("/public/informasi-pickup", warehouseController.GetInformasiPickup)

	// Informasi Pickup - Admin
	informasiPickupAdmin := v1.Group("/panel/informasi-pickup",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	informasiPickupAdmin.Get("", middleware.RequirePermission("operasional:read"), warehouseController.Get)
	informasiPickupAdmin.Get("/jadwal", middleware.RequirePermission("operasional:read"), warehouseController.GetJadwal)
	informasiPickupAdmin.Put("/jadwal", middleware.RequirePermission("operasional:manage"), warehouseController.UpdateJadwal)

	// Blog - Admin
	blogAdmin := v1.Group("/panel/blog",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	blogAdmin.Get("", middleware.RequirePermission("marketing:read"), blogController.GetAll)
	blogAdmin.Get("/statistik", middleware.RequirePermission("marketing:read"), blogController.GetStatistics)
	blogAdmin.Get("/search", middleware.RequirePermission("marketing:read"), blogController.Search)
	blogAdmin.Get("/:id", middleware.RequirePermission("marketing:read"), blogController.GetByID)
	blogAdmin.Post("", middleware.RequirePermission("marketing:manage"), blogController.Create)
	blogAdmin.Put("/:id", middleware.RequirePermission("marketing:manage"), blogController.Update)
	blogAdmin.Delete("/:id", middleware.RequirePermission("marketing:manage"), blogController.Delete)
	blogAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), blogController.ToggleStatus)

	// Blog - Public
	blogPublic := v1.Group("/public/blog")
	blogPublic.Get("", blogController.GetAll)
	blogPublic.Get("/populer", blogController.GetPopular)
	blogPublic.Get("/slug/:slug", blogController.GetBySlug)
	blogPublic.Get("/search", blogController.Search)

	// Kategori Blog - Admin
	kategoriBlogAdmin := v1.Group("/panel/kategori-blog",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	kategoriBlogAdmin.Get("", middleware.RequirePermission("marketing:read"), kategoriBlogController.GetAll)
	kategoriBlogAdmin.Get("/dropdown", middleware.RequirePermission("marketing:read"), kategoriBlogController.GetDropdownOptions)
	kategoriBlogAdmin.Get("/:id", middleware.RequirePermission("marketing:read"), kategoriBlogController.GetByID)
	kategoriBlogAdmin.Post("", middleware.RequirePermission("marketing:manage"), kategoriBlogController.Create)
	kategoriBlogAdmin.Put("/:id", middleware.RequirePermission("marketing:manage"), kategoriBlogController.Update)
	kategoriBlogAdmin.Delete("/:id", middleware.RequirePermission("marketing:manage"), kategoriBlogController.Delete)
	kategoriBlogAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), kategoriBlogController.ToggleStatus)
	kategoriBlogAdmin.Patch("/:id/reorder", middleware.RequirePermission("marketing:manage"), kategoriBlogController.Reorder)

	// Kategori Blog - Public
	v1.Get("/public/kategori-blog", kategoriBlogController.GetAllPublic)

	// Label Blog - Admin
	labelBlogAdmin := v1.Group("/panel/label-blog",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	labelBlogAdmin.Get("", middleware.RequirePermission("marketing:read"), labelBlogController.GetAll)
	labelBlogAdmin.Get("/dropdown", middleware.RequirePermission("marketing:read"), labelBlogController.GetDropdownOptions)
	labelBlogAdmin.Get("/:id", middleware.RequirePermission("marketing:read"), labelBlogController.GetByID)
	labelBlogAdmin.Post("", middleware.RequirePermission("marketing:manage"), labelBlogController.Create)
	labelBlogAdmin.Put("/:id", middleware.RequirePermission("marketing:manage"), labelBlogController.Update)
	labelBlogAdmin.Delete("/:id", middleware.RequirePermission("marketing:manage"), labelBlogController.Delete)
	labelBlogAdmin.Patch("/:id/reorder", middleware.RequirePermission("marketing:manage"), labelBlogController.Reorder)

	// Label Blog - Public
	v1.Get("/public/label-blog", labelBlogController.GetAllPublic)

	// Video - Admin
	videoAdmin := v1.Group("/panel/video",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	videoAdmin.Get("", middleware.RequirePermission("marketing:read"), videoController.GetAll)
	videoAdmin.Get("/statistik", middleware.RequirePermission("marketing:read"), videoController.GetStatistics)
	videoAdmin.Get("/dropdown", middleware.RequirePermission("marketing:read"), videoController.GetDropdownOptions)
	videoAdmin.Get("/search", middleware.RequirePermission("marketing:read"), videoController.Search)
	videoAdmin.Get("/:id", middleware.RequirePermission("marketing:read"), videoController.GetByID)
	videoAdmin.Post("", middleware.RequirePermission("marketing:manage"), videoController.Create)
	videoAdmin.Put("/:id", middleware.RequirePermission("marketing:manage"), videoController.Update)
	videoAdmin.Delete("/:id", middleware.RequirePermission("marketing:manage"), videoController.Delete)
	videoAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), videoController.ToggleStatus)

	// Video - Public
	videoPublic := v1.Group("/public/video")
	videoPublic.Get("", videoController.GetAll)
	videoPublic.Get("/slug/:slug", videoController.GetBySlug)
	videoPublic.Get("/popular", videoController.GetPopular)
	videoPublic.Get("/search", videoController.Search)

	// Kategori Video - Admin
	kategoriVideoAdmin := v1.Group("/panel/kategori-video",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	kategoriVideoAdmin.Get("", middleware.RequirePermission("marketing:read"), kategoriVideoController.GetAll)
	kategoriVideoAdmin.Get("/dropdown", middleware.RequirePermission("marketing:read"), kategoriVideoController.GetDropdownOptions)
	kategoriVideoAdmin.Get("/:id", middleware.RequirePermission("marketing:read"), kategoriVideoController.GetByID)
	kategoriVideoAdmin.Post("", middleware.RequirePermission("marketing:manage"), kategoriVideoController.Create)
	kategoriVideoAdmin.Put("/:id", middleware.RequirePermission("marketing:manage"), kategoriVideoController.Update)
	kategoriVideoAdmin.Delete("/:id", middleware.RequirePermission("marketing:manage"), kategoriVideoController.Delete)
	kategoriVideoAdmin.Patch("/:id/toggle-status", middleware.RequirePermission("marketing:manage"), kategoriVideoController.ToggleStatus)
	kategoriVideoAdmin.Patch("/:id/reorder", middleware.RequirePermission("marketing:manage"), kategoriVideoController.Reorder)

	// Kategori Video - Public
	v1.Get("/public/kategori-video", kategoriVideoController.GetAllPublic)
	v1.Get("/public/kategori-video/dropdown", kategoriVideoController.GetDropdownOptions)

	// Routes list endpoint
	router.Get("/api/routes", func(c *fiber.Ctx) error {
		var endpointList []fiber.Map
		for _, route := range router.GetRoutes() {
			endpointList = append(endpointList, fiber.Map{"method": route.Method, "path": route.Path})
		}
		return c.JSON(endpointList)
	})
}
