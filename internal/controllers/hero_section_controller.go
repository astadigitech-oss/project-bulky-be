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

type HeroSectionController struct {
	service services.HeroSectionService
	cfg     *config.Config
}

func NewHeroSectionController(service services.HeroSectionService, cfg *config.Config) *HeroSectionController {
	return &HeroSectionController{
		service: service,
		cfg:     cfg,
	}
}

func (c *HeroSectionController) Create(ctx *gin.Context) {
	var req models.CreateHeroSectionRequest
	var gambarURL *string

	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.Nama = ctx.PostForm("nama")
		req.TanggalMulai = nil
		if tm := ctx.PostForm("tanggal_mulai"); tm != "" {
			req.TanggalMulai = &tm
		}
		req.TanggalSelesai = nil
		if ts := ctx.PostForm("tanggal_selesai"); ts != "" {
			req.TanggalSelesai = &ts
		}

		// Parse urutan (optional, default 0)
		// if urutanStr := ctx.PostForm("urutan"); urutanStr != "" {
		// 	if urutan, err := strconv.Atoi(urutanStr); err == nil {
		// 		req.Urutan = urutan
		// 	}
		// }

		// Parse is_active (optional, default false)
		if isActiveStr := ctx.PostForm("is_active"); isActiveStr != "" {
			req.IsActive = isActiveStr == "true" || isActiveStr == "1"
		}

		// Validate required fields
		if req.Nama == "" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "nama wajib diisi", nil)
			return
		}

		// Handle file upload
		if file, err := ctx.FormFile("gambar"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "hero-section", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error(), nil)
				return
			}
			gambarURL = &savedPath
		} else {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "File gambar wajib diupload", nil)
			return
		}

		req.Gambar = *gambarURL

		result, err := c.service.Create(ctx.Request.Context(), &req)
		if err != nil {
			// Rollback: delete uploaded file if creation fails
			if gambarURL != nil {
				utils.DeleteFile(*gambarURL, c.cfg)
			}
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		utils.CreatedResponse(ctx, "Hero section berhasil dibuat", result)
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

	utils.CreatedResponse(ctx, "Hero section berhasil dibuat", result)
}

func (c *HeroSectionController) FindAll(ctx *gin.Context) {
	var params models.HeroSectionFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data hero section berhasil diambil", items, *meta)
}

func (c *HeroSectionController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail hero section berhasil diambil", result)
}

func (c *HeroSectionController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateHeroSectionRequest
	var gambarURL *string
	var oldGambar *string

	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Get old data for rollback
		oldData, err := c.service.FindByID(ctx.Request.Context(), id)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Hero section tidak ditemukan", nil)
			return
		}
		oldGambar = &oldData.Gambar

		// Parse form data
		if nama := ctx.PostForm("nama"); nama != "" {
			req.Nama = &nama
		}
		if tm := ctx.PostForm("tanggal_mulai"); tm != "" {
			req.TanggalMulai = &tm
		}
		if ts := ctx.PostForm("tanggal_selesai"); ts != "" {
			req.TanggalSelesai = &ts
		}
		// if urutanStr := ctx.PostForm("urutan"); urutanStr != "" {
		// 	if urutan, err := strconv.Atoi(urutanStr); err == nil {
		// 		req.Urutan = &urutan
		// 	}
		// }
		if isActiveStr := ctx.PostForm("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true" || isActiveStr == "1"
			req.IsActive = &isActive
		}

		// Handle file upload (optional for update)
		if file, err := ctx.FormFile("gambar"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file tidak didukung", nil)
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "hero-section", c.cfg)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error(), nil)
				return
			}
			gambarURL = &savedPath
			req.Gambar = gambarURL
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

		utils.SuccessResponse(ctx, "Hero section berhasil diupdate", result)
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

	utils.SuccessResponse(ctx, "Hero section berhasil diupdate", result)
}

func (c *HeroSectionController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "hero section tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Hero section berhasil dihapus", nil)
}

func (c *HeroSectionController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status hero section berhasil diubah", result)
}

func (c *HeroSectionController) Reorder(ctx *gin.Context) {
	var req models.ReorderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	if err := c.service.Reorder(ctx.Request.Context(), &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Urutan hero section berhasil diupdate", nil)
}

func (c *HeroSectionController) GetActive(ctx *gin.Context) {
	result, err := c.service.GetVisibleHero(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if result == nil {
		utils.SuccessResponse(ctx, "Tidak ada hero section aktif", nil)
		return
	}

	utils.SuccessResponse(ctx, "Hero section aktif berhasil diambil", result)
}
