package controllers

import (
	"net/http"
	"strings"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
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

func (c *HeroSectionController) Create(ctx *fiber.Ctx) error {
	var req models.CreateHeroSectionRequest
	var gambarIDURL *string
	var gambarENURL *string

	contentType := ctx.Get("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.Nama = ctx.FormValue("nama")
		req.TanggalMulai = nil
		if tm := ctx.FormValue("tanggal_mulai"); tm != "" {
			req.TanggalMulai = &tm
		}
		req.TanggalSelesai = nil
		if ts := ctx.FormValue("tanggal_selesai"); ts != "" {
			req.TanggalSelesai = &ts
		}

		// Note: is_default cannot be set on create, use toggle endpoint after creation

		// Validate required fields
		if req.Nama == "" {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "nama wajib diisi", nil)
		}

		// Handle gambar_id upload (required)
		if file, err := ctx.FormFile("gambar_id"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar_id tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "hero-section", c.cfg)
			if err != nil {
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file gambar_id: "+err.Error(), nil)
			}
			gambarIDURL = &savedPath
			req.GambarID = savedPath
		} else {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "File gambar_id wajib diupload", nil)
		}

		// Handle gambar_en upload (optional)
		if file, err := ctx.FormFile("gambar_en"); err == nil {
			if !utils.IsValidImageType(file) {
				// Rollback gambar_id if gambar_en is invalid
				if gambarIDURL != nil {
					utils.DeleteFile(*gambarIDURL, c.cfg)
				}
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar_en tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "hero-section", c.cfg)
			if err != nil {
				// Rollback gambar_id if gambar_en upload fails
				if gambarIDURL != nil {
					utils.DeleteFile(*gambarIDURL, c.cfg)
				}
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file gambar_en: "+err.Error(), nil)
			}
			gambarENURL = &savedPath
			req.GambarEN = &savedPath
		}

		result, err := c.service.Create(ctx.UserContext(), &req)
		if err != nil {
			// Rollback: delete uploaded files if creation fails
			if gambarIDURL != nil {
				utils.DeleteFile(*gambarIDURL, c.cfg)
			}
			if gambarENURL != nil {
				utils.DeleteFile(*gambarENURL, c.cfg)
			}
			return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		}

		return utils.CreatedResponse(ctx, "Hero section berhasil dibuat", result)
	}

	// Handle JSON request (for backward compatibility)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.CreatedResponse(ctx, "Hero section berhasil dibuat", result)
}

func (c *HeroSectionController) FindAll(ctx *fiber.Ctx) error {
	var params models.HeroSectionFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data hero section berhasil diambil", items, *meta)
}

func (c *HeroSectionController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail hero section berhasil diambil", result)
}

func (c *HeroSectionController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateHeroSectionRequest
	var gambarIDURL *string
	var gambarENURL *string
	var oldGambarID *string
	var oldGambarEN *string

	contentType := ctx.Get("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Get old data for rollback
		oldData, err := c.service.FindByID(ctx.UserContext(), id)
		if err != nil {
			return utils.ErrorResponse(ctx, http.StatusNotFound, "Hero section tidak ditemukan", nil)
		}
		oldGambarID = &oldData.GambarURL.ID
		if oldData.GambarURL.EN != nil {
			oldGambarEN = oldData.GambarURL.EN
		}

		// Parse form data
		if nama := ctx.FormValue("nama"); nama != "" {
			req.Nama = &nama
		}
		// Note: is_default cannot be updated here, use toggle endpoint instead
		if tm := ctx.FormValue("tanggal_mulai"); tm != "" {
			req.TanggalMulai = &tm
		}
		if ts := ctx.FormValue("tanggal_selesai"); ts != "" {
			req.TanggalSelesai = &ts
		}

		// Handle gambar_id upload (optional for update)
		if file, err := ctx.FormFile("gambar_id"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar_id tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "hero-section", c.cfg)
			if err != nil {
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file gambar_id: "+err.Error(), nil)
			}
			gambarIDURL = &savedPath
			req.GambarID = &savedPath
		}

		// Handle gambar_en upload (optional for update)
		if file, err := ctx.FormFile("gambar_en"); err == nil {
			if !utils.IsValidImageType(file) {
				// Rollback gambar_id if gambar_en is invalid
				if gambarIDURL != nil {
					utils.DeleteFile(*gambarIDURL, c.cfg)
				}
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar_en tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "hero-section", c.cfg)
			if err != nil {
				// Rollback gambar_id if gambar_en upload fails
				if gambarIDURL != nil {
					utils.DeleteFile(*gambarIDURL, c.cfg)
				}
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file gambar_en: "+err.Error(), nil)
			}
			gambarENURL = &savedPath
			req.GambarEN = &savedPath
		}

		result, err := c.service.Update(ctx.UserContext(), id, &req)
		if err != nil {
			// Rollback: delete newly uploaded files if update fails
			if gambarIDURL != nil {
				utils.DeleteFile(*gambarIDURL, c.cfg)
			}
			if gambarENURL != nil {
				utils.DeleteFile(*gambarENURL, c.cfg)
			}
			return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		}

		// Delete old files after successful update (only if new files were uploaded)
		if gambarIDURL != nil && oldGambarID != nil {
			utils.DeleteFile(*oldGambarID, c.cfg)
		}
		if gambarENURL != nil && oldGambarEN != nil {
			utils.DeleteFile(*oldGambarEN, c.cfg)
		}

		return utils.SuccessResponse(ctx, "Hero section berhasil diupdate", result)
	}

	// Handle JSON request (for backward compatibility)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Hero section berhasil diupdate", result)
}

func (c *HeroSectionController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "hero section tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Hero section berhasil dihapus", nil)
}

func (c *HeroSectionController) ToggleStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Status default hero section berhasil diubah", result)
}

func (c *HeroSectionController) GetActive(ctx *fiber.Ctx) error {
	result, err := c.service.GetVisibleHero(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	if result == nil {
		return utils.SuccessResponse(ctx, "Tidak ada hero section aktif", nil)
	}

	return utils.SuccessResponse(ctx, "Hero section aktif berhasil diambil", result)
}

func (c *HeroSectionController) GetSchedules(ctx *fiber.Ctx) error {
	result, err := c.service.GetSchedules(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil jadwal hero section: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Jadwal hero section berhasil diambil", result)
}
