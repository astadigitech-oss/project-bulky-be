package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BannerEventPromoController struct {
	service        services.BannerEventPromoService
	reorderService *services.ReorderService
	cfg            *config.Config
}

func NewBannerEventPromoController(service services.BannerEventPromoService, reorderService *services.ReorderService, cfg *config.Config) *BannerEventPromoController {
	return &BannerEventPromoController{
		service:        service,
		reorderService: reorderService,
		cfg:            cfg,
	}
}

func (c *BannerEventPromoController) Create(ctx *gin.Context) {
	var req models.CreateBannerEventPromoRequest
	var gambarIDURL *string
	var gambarENURL *string

	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.Nama = ctx.PostForm("nama")

		// Parse tujuan JSON (optional)
		if tujuanStr := ctx.PostForm("tujuan"); tujuanStr != "" {
			var tujuan []models.TujuanKategoriInput
			if err := json.Unmarshal([]byte(tujuanStr), &tujuan); err == nil {
				req.Tujuan = tujuan
			}
		}

		req.TanggalMulai = nil
		if tm := ctx.PostForm("tanggal_mulai"); tm != "" {
			req.TanggalMulai = &tm
		}
		req.TanggalSelesai = nil
		if ts := ctx.PostForm("tanggal_selesai"); ts != "" {
			req.TanggalSelesai = &ts
		}

		// Validate required fields
		if req.Nama == "" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "nama wajib diisi", nil)
			return
		}

		// Handle gambar_id upload (required)
		if file, err := ctx.FormFile("gambar_id"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar_id tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "banner-event-promo", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file gambar_id: "+err.Error(), nil)
				return
			}
			gambarIDURL = &savedPath
		} else {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "File gambar_id wajib diupload", nil)
			return
		}

		// Handle gambar_en upload (required)
		if file, err := ctx.FormFile("gambar_en"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar_en tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "banner-event-promo", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file gambar_en: "+err.Error(), nil)
				return
			}
			gambarENURL = &savedPath
		} else {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "File gambar_en wajib diupload", nil)
			return
		}

		req.GambarID = *gambarIDURL
		req.GambarEN = *gambarENURL

		result, err := c.service.Create(ctx.Request.Context(), &req)
		if err != nil {
			// Rollback: delete uploaded files if creation fails
			if gambarIDURL != nil {
				utils.DeleteFile(*gambarIDURL, c.cfg)
			}
			if gambarENURL != nil {
				utils.DeleteFile(*gambarENURL, c.cfg)
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

func (c *BannerEventPromoController) FindAll(ctx *gin.Context) {
	var params models.BannerEventPromoFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data banner berhasil diambil", items, *meta)
}

func (c *BannerEventPromoController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail banner berhasil diambil", result)
}

func (c *BannerEventPromoController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateBannerEventPromoRequest
	var gambarURL *string
	var oldGambar *string

	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Get old data for rollback
		oldData, err := c.service.FindByID(ctx.Request.Context(), id)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Banner tidak ditemukan", nil)
			return
		}
		oldGambar = &oldData.GambarURL.ID

		// Parse form data
		if nama := ctx.PostForm("nama"); nama != "" {
			req.Nama = &nama
		}

		// Parse tujuan JSON (optional)
		if tujuanStr := ctx.PostForm("tujuan"); tujuanStr != "" {
			var tujuan []models.TujuanKategoriInput
			if err := json.Unmarshal([]byte(tujuanStr), &tujuan); err == nil {
				req.Tujuan = tujuan
			}
		}

		if tm := ctx.PostForm("tanggal_mulai"); tm != "" {
			req.TanggalMulai = &tm
		}
		if ts := ctx.PostForm("tanggal_selesai"); ts != "" {
			req.TanggalSelesai = &ts
		}

		// Handle gambar_id upload (optional for update)
		if file, err := ctx.FormFile("gambar_id"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar_id tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "banner-event-promo", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file gambar_id: "+err.Error(), nil)
				return
			}
			gambarURL = &savedPath
			req.GambarID = gambarURL
		}

		// Handle gambar_en upload (optional for update)
		var gambarENURL *string
		if file, err := ctx.FormFile("gambar_en"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar_en tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "banner-event-promo", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file gambar_en: "+err.Error(), nil)
				return
			}
			gambarENURL = &savedPath
			req.GambarEN = gambarENURL
		}

		result, err := c.service.Update(ctx.Request.Context(), id, &req)
		if err != nil {
			// Rollback: delete newly uploaded file if update fails
			if gambarURL != nil {
				utils.DeleteFile(*gambarURL, c.cfg)
			}
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		// Delete old file after successful update (only if new file was uploaded)
		if gambarURL != nil && oldGambar != nil {
			utils.DeleteFile(*oldGambar, c.cfg)
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
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Banner berhasil diupdate", result)
}

func (c *BannerEventPromoController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "banner tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Banner berhasil dihapus", nil)
}

func (c *BannerEventPromoController) Reorder(ctx *gin.Context) {
	var req models.ReorderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	if err := c.service.Reorder(ctx.Request.Context(), &req); err != nil {
		// Check if error is "not found"
		if strings.Contains(err.Error(), "tidak ditemukan") {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Urutan banner berhasil diupdate", nil)
}

func (c *BannerEventPromoController) GetActive(ctx *gin.Context) {
	result, err := c.service.GetVisibleBanners(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data banner aktif berhasil diambil", result)
}

func (c *BannerEventPromoController) ReorderByDirection(ctx *gin.Context) {
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
		"banner_event_promo",
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
