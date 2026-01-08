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

type MerekProdukController struct {
	service services.MerekProdukService
	cfg     *config.Config
}

func NewMerekProdukController(service services.MerekProdukService, cfg *config.Config) *MerekProdukController {
	return &MerekProdukController{
		service: service,
		cfg:     cfg,
	}
}

func (c *MerekProdukController) Create(ctx *gin.Context) {
	var req models.CreateMerekProdukRequest
	var logoURL *string

	// Check content type
	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.Nama = ctx.PostForm("nama")

		// Validate required field
		if req.Nama == "" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Nama merek wajib diisi", nil)
			return
		}

		// Handle logo upload or URL
		if file, err := ctx.FormFile("logo"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file logo tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-brands", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan logo: "+err.Error(), nil)
				return
			}
			logoURL = &savedPath
		} else if logoStr := ctx.PostForm("logo"); logoStr != "" {
			logoURL = &logoStr
		}

		// Create with logo
		result, err := c.service.CreateWithLogo(ctx.Request.Context(), &req, logoURL)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.CreatedResponse(ctx, "Merek produk berhasil dibuat", result)
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

	utils.CreatedResponse(ctx, "Merek produk berhasil dibuat", result)
}

func (c *MerekProdukController) FindAll(ctx *gin.Context) {
	var params models.PaginationRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data merek produk berhasil diambil", items, *meta)
}

func (c *MerekProdukController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail merek produk berhasil diambil", result)
}

func (c *MerekProdukController) FindBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	result, err := c.service.FindBySlug(ctx.Request.Context(), slug)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail merek produk berhasil diambil", result)
}

func (c *MerekProdukController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateMerekProdukRequest
	var logoURL *string

	// Check content type
	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		if nama := ctx.PostForm("nama"); nama != "" {
			req.Nama = &nama
		}
		if isActiveStr := ctx.PostForm("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true"
			req.IsActive = &isActive
		}

		// Handle logo upload or URL
		if file, err := ctx.FormFile("logo"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file logo tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-brands", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan logo: "+err.Error(), nil)
				return
			}
			logoURL = &savedPath
		} else if logoStr := ctx.PostForm("logo"); logoStr != "" {
			logoURL = &logoStr
		}

		// Use UpdateWithLogo for multipart
		result, err := c.service.UpdateWithLogo(ctx.Request.Context(), id, &req, logoURL)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.SuccessResponse(ctx, "Merek produk berhasil diupdate", result)
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

	utils.SuccessResponse(ctx, "Merek produk berhasil diupdate", result)
}

func (c *MerekProdukController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "merek produk tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Merek produk berhasil dihapus", nil)
}

func (c *MerekProdukController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status merek berhasil diubah", result)
}
