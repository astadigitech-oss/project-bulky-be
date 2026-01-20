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

type BannerTipeProdukController struct {
	service        services.BannerTipeProdukService
	reorderService *services.ReorderService
	cfg            *config.Config
}

func NewBannerTipeProdukController(service services.BannerTipeProdukService, reorderService *services.ReorderService, cfg *config.Config) *BannerTipeProdukController {
	return &BannerTipeProdukController{
		service:        service,
		reorderService: reorderService,
		cfg:            cfg,
	}
}

func (c *BannerTipeProdukController) Create(ctx *gin.Context) {
	var req models.CreateBannerTipeProdukRequest
	var gambarURL *string

	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.TipeProdukID = ctx.PostForm("tipe_produk_id")
		req.Nama = ctx.PostForm("nama")

		// Validate required fields
		if req.TipeProdukID == "" || req.Nama == "" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "tipe_produk_id dan nama wajib diisi", nil)
			return
		}

		// Handle file upload
		if file, err := ctx.FormFile("file"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "banners/tipe-produk", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error(), nil)
				return
			}
			gambarURL = &savedPath
		} else {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "File banner wajib diupload", nil)
			return
		}

		req.GambarURL = *gambarURL

		result, err := c.service.Create(ctx.Request.Context(), &req)
		if err != nil {
			// Rollback: delete uploaded file if creation fails
			if gambarURL != nil {
				utils.DeleteFile(*gambarURL, c.cfg)
			}
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		utils.CreatedResponse(ctx, "Banner berhasil dibuat", result)
		return
	}

	// Handle JSON request (for backward compatibility)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Banner berhasil dibuat", result)
}

func (c *BannerTipeProdukController) FindAll(ctx *gin.Context) {
	// Check if grouped response is requested (no pagination params)
	_, hasPage := ctx.GetQuery("page")
	_, hasPerPage := ctx.GetQuery("per_page")
	search := ctx.Query("search")

	// If no pagination params, return grouped response
	if !hasPage && !hasPerPage {
		grouped, meta, err := c.service.FindAllGrouped(ctx.Request.Context(), search)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Data banner tipe produk berhasil diambil",
			"data":    grouped,
			"meta":    meta,
		})
		return
	}

	// Otherwise, return paginated response (old behavior)
	var params models.BannerTipeProdukFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	tipeProdukID := ctx.Query("tipe_produk_id")

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params, tipeProdukID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data banner berhasil diambil", items, *meta)
}

func (c *BannerTipeProdukController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail banner berhasil diambil", result)
}

func (c *BannerTipeProdukController) FindByTipeProdukID(ctx *gin.Context) {
	tipeProdukID := ctx.Param("tipe_produk_id")

	items, err := c.service.FindByTipeProdukID(ctx.Request.Context(), tipeProdukID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Banner berhasil diambil", items)
}

func (c *BannerTipeProdukController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateBannerTipeProdukRequest
	var newGambarURL *string

	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with optional file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		if tipeProdukID := ctx.PostForm("tipe_produk_id"); tipeProdukID != "" {
			req.TipeProdukID = &tipeProdukID
		}
		if nama := ctx.PostForm("nama"); nama != "" {
			req.Nama = &nama
		}
		if isActiveStr := ctx.PostForm("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true"
			req.IsActive = &isActive
		}

		// Handle file upload (optional)
		if file, err := ctx.FormFile("file"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "banners/tipe-produk", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error(), nil)
				return
			}
			newGambarURL = &savedPath
			req.GambarURL = newGambarURL
		}

		result, err := c.service.UpdateWithFile(ctx.Request.Context(), id, &req, newGambarURL)
		if err != nil {
			// Rollback: delete new uploaded file if update fails
			if newGambarURL != nil {
				utils.DeleteFile(*newGambarURL, c.cfg)
			}
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}

		utils.SuccessResponse(ctx, "Banner berhasil diupdate", result)
		return
	}

	// Handle JSON request (for backward compatibility)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Banner berhasil diupdate", result)
}

func (c *BannerTipeProdukController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteWithFile(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "banner tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Banner berhasil dihapus", nil)
}

func (c *BannerTipeProdukController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status banner berhasil diubah", result)
}

func (c *BannerTipeProdukController) Reorder(ctx *gin.Context) {
	var req models.ReorderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	if err := c.service.Reorder(ctx.Request.Context(), &req); err != nil {
		status := http.StatusInternalServerError
		// Return 404 if banner not found
		if err.Error() == "salah satu atau lebih banner tidak ditemukan" {
			status = http.StatusNotFound
		} else if err.Error() == "items tidak boleh kosong" {
			status = http.StatusBadRequest
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Urutan banner berhasil diubah", nil)
}

func (c *BannerTipeProdukController) ReorderByDirection(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.ReorderByDirectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	idUUID, err := utils.ParseUUID(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	result, err := c.reorderService.Reorder(
		ctx.Request.Context(),
		"banner_tipe_produk",
		idUUID,
		req.Direction,
		"",
		nil,
	)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Urutan berhasil diubah", gin.H{
		"item": gin.H{
			"id":     result.ItemID,
			"urutan": result.ItemUrutan,
		},
		"swapped_with": gin.H{
			"id":     result.SwappedID,
			"urutan": result.SwappedUrutan,
		},
	})
}
