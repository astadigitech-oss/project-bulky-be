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
	provinsiController *controllers.ProvinsiController,
	kotaController *controllers.KotaController,
	kecamatanController *controllers.KecamatanController,
	kelurahanController *controllers.KelurahanController,
	wilayahController *controllers.WilayahController,
) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "Server is running"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public Auth Routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
		}

		// Protected Auth Routes
		authProtected := v1.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware())
		{
			authProtected.POST("/logout", authController.Logout)
			authProtected.GET("/me", authController.Me)
			authProtected.PUT("/profile", authController.UpdateProfile)
			authProtected.PUT("/change-password", authController.ChangePassword)
		}

		// Admin Management Routes (Protected)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		{
			admin.GET("", adminController.FindAll)
			admin.GET("/:id", adminController.FindByID)
			admin.POST("", adminController.Create)
			admin.PUT("/:id", adminController.Update)
			admin.DELETE("/:id", adminController.Delete)
			admin.PATCH("/:id/toggle-status", adminController.ToggleStatus)
			admin.PUT("/:id/reset-password", adminController.ResetPassword)
		}

		// Buyer Management Routes (Protected - RUD only)
		buyer := v1.Group("/buyer")
		buyer.Use(middleware.AuthMiddleware())
		{
			buyer.GET("", buyerController.FindAll)
			buyer.GET("/statistik", buyerController.GetStatistik)
			buyer.GET("/:id", buyerController.FindByID)
			buyer.PUT("/:id", buyerController.Update)
			buyer.DELETE("/:id", buyerController.Delete)
			buyer.PATCH("/:id/toggle-status", buyerController.ToggleStatus)
			buyer.PUT("/:id/reset-password", buyerController.ResetPassword)
		}

		// Alamat Buyer Routes (Protected)
		alamatBuyer := v1.Group("/alamat-buyer")
		alamatBuyer.Use(middleware.AuthMiddleware())
		{
			alamatBuyer.GET("", alamatBuyerController.FindAll)
			alamatBuyer.GET("/:id", alamatBuyerController.FindByID)
			alamatBuyer.POST("", alamatBuyerController.Create)
			alamatBuyer.PUT("/:id", alamatBuyerController.Update)
			alamatBuyer.DELETE("/:id", alamatBuyerController.Delete)
			alamatBuyer.PATCH("/:id/set-default", alamatBuyerController.SetDefault)
		}

		// Wilayah Dropdown Routes (Public - no auth)
		wilayah := v1.Group("/wilayah")
		{
			wilayah.GET("/provinsi", wilayahController.DropdownProvinsi)
			wilayah.GET("/kota", wilayahController.DropdownKota)
			wilayah.GET("/kecamatan", wilayahController.DropdownKecamatan)
			wilayah.GET("/kelurahan", wilayahController.DropdownKelurahan)
		}

		// Provinsi CRUD Routes (Protected)
		provinsi := v1.Group("/provinsi")
		provinsi.Use(middleware.AuthMiddleware())
		{
			provinsi.GET("", provinsiController.FindAll)
			provinsi.GET("/:id", provinsiController.FindByID)
			provinsi.POST("", provinsiController.Create)
			provinsi.PUT("/:id", provinsiController.Update)
			provinsi.DELETE("/:id", provinsiController.Delete)
		}

		// Kota CRUD Routes (Protected)
		kota := v1.Group("/kota")
		kota.Use(middleware.AuthMiddleware())
		{
			kota.GET("", kotaController.FindAll)
			kota.GET("/:id", kotaController.FindByID)
			kota.POST("", kotaController.Create)
			kota.PUT("/:id", kotaController.Update)
			kota.DELETE("/:id", kotaController.Delete)
		}

		// Kecamatan CRUD Routes (Protected)
		kecamatan := v1.Group("/kecamatan")
		kecamatan.Use(middleware.AuthMiddleware())
		{
			kecamatan.GET("", kecamatanController.FindAll)
			kecamatan.GET("/:id", kecamatanController.FindByID)
			kecamatan.POST("", kecamatanController.Create)
			kecamatan.PUT("/:id", kecamatanController.Update)
			kecamatan.DELETE("/:id", kecamatanController.Delete)
		}

		// Kelurahan CRUD Routes (Protected)
		kelurahan := v1.Group("/kelurahan")
		kelurahan.Use(middleware.AuthMiddleware())
		{
			kelurahan.GET("", kelurahanController.FindAll)
			kelurahan.GET("/:id", kelurahanController.FindByID)
			kelurahan.POST("", kelurahanController.Create)
			kelurahan.PUT("/:id", kelurahanController.Update)
			kelurahan.DELETE("/:id", kelurahanController.Delete)
		}

		// Kategori Produk
		kategori := v1.Group("/kategori-produk")
		{
			kategori.GET("", kategoriController.FindAll)
			kategori.GET("/:id", kategoriController.FindByID)
			kategori.GET("/slug/:slug", kategoriController.FindBySlug)
			kategori.POST("", kategoriController.Create)
			kategori.PUT("/:id", kategoriController.Update)
			kategori.DELETE("/:id", kategoriController.Delete)
			kategori.PATCH("/:id/toggle-status", kategoriController.ToggleStatus)
		}

		// Merek Produk
		merek := v1.Group("/merek-produk")
		{
			merek.GET("", merekController.FindAll)
			merek.GET("/:id", merekController.FindByID)
			merek.GET("/slug/:slug", merekController.FindBySlug)
			merek.POST("", merekController.Create)
			merek.PUT("/:id", merekController.Update)
			merek.DELETE("/:id", merekController.Delete)
			merek.PATCH("/:id/toggle-status", merekController.ToggleStatus)
		}

		// Kondisi Produk
		kondisi := v1.Group("/kondisi-produk")
		{
			kondisi.GET("", kondisiController.FindAll)
			kondisi.GET("/:id", kondisiController.FindByID)
			kondisi.GET("/slug/:slug", kondisiController.FindBySlug)
			kondisi.POST("", kondisiController.Create)
			kondisi.PUT("/:id", kondisiController.Update)
			kondisi.DELETE("/:id", kondisiController.Delete)
			kondisi.PATCH("/:id/toggle-status", kondisiController.ToggleStatus)
			kondisi.PUT("/reorder", kondisiController.Reorder)
		}

		// Kondisi Paket
		paket := v1.Group("/kondisi-paket")
		{
			paket.GET("", kondisiPaketController.FindAll)
			paket.GET("/:id", kondisiPaketController.FindByID)
			paket.GET("/slug/:slug", kondisiPaketController.FindBySlug)
			paket.POST("", kondisiPaketController.Create)
			paket.PUT("/:id", kondisiPaketController.Update)
			paket.DELETE("/:id", kondisiPaketController.Delete)
			paket.PATCH("/:id/toggle-status", kondisiPaketController.ToggleStatus)
			paket.PUT("/reorder", kondisiPaketController.Reorder)
		}

		// Sumber Produk
		sumber := v1.Group("/sumber-produk")
		{
			sumber.GET("", sumberController.FindAll)
			sumber.GET("/:id", sumberController.FindByID)
			sumber.GET("/slug/:slug", sumberController.FindBySlug)
			sumber.POST("", sumberController.Create)
			sumber.PUT("/:id", sumberController.Update)
			sumber.DELETE("/:id", sumberController.Delete)
			sumber.PATCH("/:id/toggle-status", sumberController.ToggleStatus)
		}

		// Warehouse
		warehouse := v1.Group("/warehouse")
		{
			warehouse.GET("", warehouseController.FindAll)
			warehouse.GET("/:id", warehouseController.FindByID)
			warehouse.POST("", warehouseController.Create)
			warehouse.PUT("/:id", warehouseController.Update)
			warehouse.DELETE("/:id", warehouseController.Delete)
			warehouse.PATCH("/:id/toggle-status", warehouseController.ToggleStatus)
		}

		// Tipe Produk
		tipeProduk := v1.Group("/tipe-produk")
		{
			tipeProduk.GET("", tipeProdukController.FindAll)
			tipeProduk.GET("/:id", tipeProdukController.FindByID)
			tipeProduk.GET("/slug/:slug", tipeProdukController.FindBySlug)
			tipeProduk.POST("", tipeProdukController.Create)
			tipeProduk.PUT("/:id", tipeProdukController.Update)
			tipeProduk.DELETE("/:id", tipeProdukController.Delete)
			tipeProduk.PATCH("/:id/toggle-status", tipeProdukController.ToggleStatus)
			tipeProduk.PUT("/reorder", tipeProdukController.Reorder)
		}

		// Diskon Kategori
		diskonKategori := v1.Group("/diskon-kategori")
		{
			diskonKategori.GET("", diskonKategoriController.FindAll)
			diskonKategori.GET("/:id", diskonKategoriController.FindByID)
			diskonKategori.GET("/by-kategori/:kategori_id", diskonKategoriController.FindActiveByKategoriID)
			diskonKategori.POST("", diskonKategoriController.Create)
			diskonKategori.PUT("/:id", diskonKategoriController.Update)
			diskonKategori.DELETE("/:id", diskonKategoriController.Delete)
			diskonKategori.PATCH("/:id/toggle-status", diskonKategoriController.ToggleStatus)
		}

		// Banner Tipe Produk
		bannerTipeProduk := v1.Group("/banner-tipe-produk")
		{
			bannerTipeProduk.GET("", bannerTipeProdukController.FindAll)
			bannerTipeProduk.GET("/:id", bannerTipeProdukController.FindByID)
			bannerTipeProduk.GET("/by-tipe/:tipe_produk_id", bannerTipeProdukController.FindByTipeProdukID)
			bannerTipeProduk.POST("", bannerTipeProdukController.Create)
			bannerTipeProduk.PUT("/:id", bannerTipeProdukController.Update)
			bannerTipeProduk.DELETE("/:id", bannerTipeProdukController.Delete)
			bannerTipeProduk.PATCH("/:id/toggle-status", bannerTipeProdukController.ToggleStatus)
			bannerTipeProduk.PUT("/reorder", bannerTipeProdukController.Reorder)
		}

		// Produk
		produk := v1.Group("/produk")
		{
			produk.GET("", produkController.FindAll)
			produk.GET("/:id", produkController.FindByID)
			produk.GET("/slug/:slug", produkController.FindBySlug)
			produk.POST("", produkController.Create)
			produk.PUT("/:id", produkController.Update)
			produk.DELETE("/:id", produkController.Delete)
			produk.PATCH("/:id/toggle-status", produkController.ToggleStatus)
			produk.PATCH("/:id/update-stock", produkController.UpdateStock)
			produk.POST("/:id/gambar", produkController.AddGambar)
			produk.PUT("/:id/gambar/:gambar_id", produkController.UpdateGambar)
			produk.DELETE("/:id/gambar/:gambar_id", produkController.DeleteGambar)
			produk.PUT("/:id/gambar/reorder", produkController.ReorderGambar)
			produk.POST("/:id/dokumen", produkController.AddDokumen)
			produk.DELETE("/:id/dokumen/:dokumen_id", produkController.DeleteDokumen)
		}

		// Master (Dropdown)
		master := v1.Group("/master")
		{
			master.GET("/dropdown", masterController.GetDropdown)
		}
	}

	// Endpoint to list all registered routes
	router.GET("/routes", func(c *gin.Context) {
		var endpointList []gin.H
		for _, route := range router.Routes() {
			endpointList = append(endpointList, gin.H{"method": route.Method, "path": route.Path})
		}
		c.JSON(http.StatusOK, endpointList)
	})
}
