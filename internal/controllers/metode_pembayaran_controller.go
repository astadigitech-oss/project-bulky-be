package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MetodePembayaranController struct {
	service services.MetodePembayaranService
	cfg     *config.Config
}

func NewMetodePembayaranController(service services.MetodePembayaranService, cfg *config.Config) *MetodePembayaranController {
	return &MetodePembayaranController{
		service: service,
		cfg:     cfg,
	}
}

func (c *MetodePembayaranController) GetAll(ctx *gin.Context) {
	var params models.PaginationRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	// Set default values
	params.SetDefaults()

	// Parse optional filters
	var groupID *string
	if groupIDStr := ctx.Query("group_id"); groupIDStr != "" {
		groupID = &groupIDStr
	}

	var isActive *bool
	if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
		isActiveBool := isActiveStr == "true"
		isActive = &isActiveBool
	}

	items, meta, err := c.service.GetAll(ctx.Request.Context(), &params, groupID, isActive)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data metode pembayaran berhasil diambil", items, *meta)
}

func (c *MetodePembayaranController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail metode pembayaran berhasil diambil", result)
}

func (c *MetodePembayaranController) Create(ctx *gin.Context) {
	var req models.CreateMetodePembayaranRequest
	var logoURL *string

	// Check content type
	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.GroupID = ctx.PostForm("group_id")
		req.Nama = ctx.PostForm("nama")
		req.Kode = ctx.PostForm("kode")

		// Parse urutan (optional)
		if urutanStr := ctx.PostForm("urutan"); urutanStr != "" {
			if urutan, err := strconv.Atoi(urutanStr); err == nil {
				req.Urutan = urutan
			}
		}

		// Parse is_active (optional, default true)
		req.IsActive = true
		if isActiveStr := ctx.PostForm("is_active"); isActiveStr != "" {
			req.IsActive = isActiveStr == "true"
		}

		// Validate required fields
		if req.GroupID == "" || req.Nama == "" || req.Kode == "" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Group ID, nama, dan kode wajib diisi", nil)
			return
		}

		// Handle logo upload
		if file, err := ctx.FormFile("logo"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "File harus berupa gambar (jpg, png, webp)", nil)
				return
			}

			// Check file size (max 2MB)
			if file.Size > 2*1024*1024 {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Ukuran file maksimal 2MB", nil)
				return
			}

			savedPath, err := utils.SaveUploadedFile(file, "payment-methods", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan logo: "+err.Error(), nil)
				return
			}
			logoURL = &savedPath
		}

		// Create metode pembayaran
		result, err := c.service.Create(ctx.Request.Context(), &req, logoURL)
		if err != nil {
			if err.Error() == "Kode sudah digunakan" || err.Error() == "Group tidak ditemukan" {
				utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
				return
			}
			utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		utils.CreatedResponse(ctx, "Metode pembayaran berhasil dibuat", result)
		return
	}

	// Handle application/json (no file upload)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req, nil)
	if err != nil {
		if err.Error() == "Kode sudah digunakan" || err.Error() == "Group tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Metode pembayaran berhasil dibuat", result)
}

func (c *MetodePembayaranController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateMetodePembayaranRequest
	var logoURL *string

	// Check content type
	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		if groupID := ctx.PostForm("group_id"); groupID != "" {
			req.GroupID = &groupID
		}
		if nama := ctx.PostForm("nama"); nama != "" {
			req.Nama = &nama
		}
		if kode := ctx.PostForm("kode"); kode != "" {
			req.Kode = &kode
		}
		if urutanStr := ctx.PostForm("urutan"); urutanStr != "" {
			if urutan, err := strconv.Atoi(urutanStr); err == nil {
				req.Urutan = &urutan
			}
		}
		if isActiveStr := ctx.PostForm("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true"
			req.IsActive = &isActive
		}

		// Handle logo upload
		if file, err := ctx.FormFile("logo"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "File harus berupa gambar (jpg, png, webp)", nil)
				return
			}

			// Check file size (max 2MB)
			if file.Size > 2*1024*1024 {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Ukuran file maksimal 2MB", nil)
				return
			}

			savedPath, err := utils.SaveUploadedFile(file, "payment-methods", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan logo: "+err.Error(), nil)
				return
			}
			logoURL = &savedPath
		}

		// Update metode pembayaran
		result, err := c.service.Update(ctx.Request.Context(), id, &req, logoURL)
		if err != nil {
			if err.Error() == "Metode pembayaran tidak ditemukan" {
				utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
				return
			}
			if err.Error() == "Kode sudah digunakan" || err.Error() == "Group tidak ditemukan" {
				utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
				return
			}
			utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		utils.SuccessResponse(ctx, "Metode pembayaran berhasil diupdate", result)
		return
	}

	// Handle application/json (no file upload)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req, nil)
	if err != nil {
		if err.Error() == "Metode pembayaran tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err.Error() == "Kode sudah digunakan" || err.Error() == "Group tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Metode pembayaran berhasil diupdate", result)
}

func (c *MetodePembayaranController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		if err.Error() == "Metode pembayaran tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err.Error() == "Tidak dapat menghapus metode pembayaran yang sudah digunakan dalam transaksi" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Metode pembayaran berhasil dihapus", nil)
}

func (c *MetodePembayaranController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "Metode pembayaran tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status metode pembayaran berhasil diubah", result)
}
