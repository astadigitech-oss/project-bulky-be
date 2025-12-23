package routes

import (
	"net/http"

	"project-bulky-be/internal/controllers"

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
	masterController *controllers.MasterController,
) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Server is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
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

			// Produk Gambar
			produk.POST("/:id/gambar", produkController.AddGambar)
			produk.PUT("/:id/gambar/:id", produkController.UpdateGambar)
			produk.DELETE("/:id/gambar/:id", produkController.DeleteGambar)
			produk.PUT("/:id/gambar/reorder", produkController.ReorderGambar)

			// Produk Dokumen
			produk.POST("/:id/dokumen", produkController.AddDokumen)
			produk.DELETE("/:id/dokumen/:id", produkController.DeleteDokumen)
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
			endpointList = append(endpointList, gin.H{
				"method": route.Method,
				"path":   route.Path,
			})
		}
		c.JSON(http.StatusOK, endpointList)
	})
}
