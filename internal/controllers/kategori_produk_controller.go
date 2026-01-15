package controllers

import (
	"net/http"
	"strings"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type KategoriProdukController struct {
	service services.KategoriProdukService
	cfg     *config.Config
}

func NewKategoriProdukController(service services.KategoriProdukService, cfg *config.Config) *KategoriProdukController {
	return &KategoriProdukController{
		service: service,
		cfg:     cfg,
	}
}

func (c *KategoriProdukController) Create(ctx *gin.Context) {
	var req models.CreateKategoriProdukRequest
	var iconURL *string
	var gambarKondisiURL *string

	// Check content type
	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.NamaID = ctx.PostForm("nama_id")
		if namaEN := ctx.PostForm("nama_en"); namaEN != "" {
			req.NamaEN = &namaEN
		}
		if deskripsi := ctx.PostForm("deskripsi"); deskripsi != "" {
			req.Deskripsi = &deskripsi
		}
		if memilikiKondisi := ctx.PostForm("memiliki_kondisi_tambahan"); memilikiKondisi != "" {
			hasKondisi := memilikiKondisi == "true"
			req.MemilikiKondisiTambahan = hasKondisi
		}
		if tipeKondisi := ctx.PostForm("tipe_kondisi_tambahan"); tipeKondisi != "" {
			req.TipeKondisiTambahan = &tipeKondisi
		}
		if teksKondisi := ctx.PostForm("teks_kondisi"); teksKondisi != "" {
			req.TeksKondisi = &teksKondisi
		}

		// Validate required field
		if req.NamaID == "" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Nama kategori wajib diisi", nil)
			return
		}

		// Handle icon upload or URL
		if file, err := ctx.FormFile("icon"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file icon tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-categories", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan icon: "+err.Error(), nil)
				return
			}
			iconURL = &savedPath
		} else if iconStr := ctx.PostForm("icon"); iconStr != "" {
			iconURL = &iconStr
		}

		// Handle gambar kondisi upload or URL
		if file, err := ctx.FormFile("gambar_kondisi"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar kondisi tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-categories/kondisi", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan gambar kondisi: "+err.Error(), nil)
				return
			}
			gambarKondisiURL = &savedPath
		} else if gambarKondisiStr := ctx.PostForm("gambar_kondisi"); gambarKondisiStr != "" {
			gambarKondisiURL = &gambarKondisiStr
		}

		// Create with icon
		result, err := c.service.CreateWithIcon(ctx.Request.Context(), &req, iconURL, gambarKondisiURL)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.CreatedResponse(ctx, "Kategori produk berhasil dibuat", result)
		return
	}

	// Handle application/json (no file upload)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Kategori produk berhasil dibuat", result)
}

func (c *KategoriProdukController) FindAll(ctx *gin.Context) {
	var params models.KategoriProdukFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data kategori produk berhasil diambil", items, *meta)
}

func (c *KategoriProdukController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail kategori produk berhasil diambil", result)
}

func (c *KategoriProdukController) FindBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	result, err := c.service.FindBySlug(ctx.Request.Context(), slug)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail kategori produk berhasil diambil", result)
}

func (c *KategoriProdukController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateKategoriProdukRequest
	var iconURL *string
	var gambarKondisiURL *string

	// Check content type
	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		if namaID := ctx.PostForm("nama_id"); namaID != "" {
			req.NamaID = &namaID
		}
		if namaEN := ctx.PostForm("nama_en"); namaEN != "" {
			req.NamaEN = &namaEN
		}
		if deskripsi := ctx.PostForm("deskripsi"); deskripsi != "" {
			req.Deskripsi = &deskripsi
		}
		if isActiveStr := ctx.PostForm("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true"
			req.IsActive = &isActive
		}
		if memilikiKondisi := ctx.PostForm("memiliki_kondisi_tambahan"); memilikiKondisi != "" {
			hasKondisi := memilikiKondisi == "true"
			req.MemilikiKondisiTambahan = &hasKondisi
		}
		if tipeKondisi := ctx.PostForm("tipe_kondisi_tambahan"); tipeKondisi != "" {
			req.TipeKondisiTambahan = &tipeKondisi
		}
		if teksKondisi := ctx.PostForm("teks_kondisi"); teksKondisi != "" {
			req.TeksKondisi = &teksKondisi
		}

		// Handle icon upload or URL
		if file, err := ctx.FormFile("icon"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file icon tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-categories", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan icon: "+err.Error(), nil)
				return
			}
			iconURL = &savedPath
		} else if iconStr := ctx.PostForm("icon"); iconStr != "" {
			iconURL = &iconStr
		}

		// Handle gambar kondisi upload or URL
		if file, err := ctx.FormFile("gambar_kondisi"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar kondisi tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-categories/kondisi", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan gambar kondisi: "+err.Error(), nil)
				return
			}
			gambarKondisiURL = &savedPath
		} else if gambarKondisiStr := ctx.PostForm("gambar_kondisi"); gambarKondisiStr != "" {
			gambarKondisiURL = &gambarKondisiStr
		}

		// Use UpdateWithIcon for multipart
		result, err := c.service.UpdateWithIcon(ctx.Request.Context(), id, &req, iconURL, gambarKondisiURL)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.SuccessResponse(ctx, "Kategori produk berhasil diupdate", result)
		return
	}

	// Handle application/json (no file upload)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Kategori produk berhasil diupdate", result)
}

func (c *KategoriProdukController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "kategori produk tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Kategori produk berhasil dihapus", nil)
}

func (c *KategoriProdukController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status kategori berhasil diubah", result)
}
