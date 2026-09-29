package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type DokumenKebijakanController struct {
	service     services.DokumenKebijakanService
	activityLog services.ActivityLogService
}

func NewDokumenKebijakanController(service services.DokumenKebijakanService, activityLog services.ActivityLogService) *DokumenKebijakanController {
	return &DokumenKebijakanController{service: service, activityLog: activityLog}
}

// ========================================
// PANEL ENDPOINTS (Admin)
// ========================================

// GetAll - Get all dokumen kebijakan (7 fixed pages)
func (c *DokumenKebijakanController) GetAll(ctx *fiber.Ctx) error {
	items, err := c.service.FindAll(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data dokumen kebijakan berhasil diambil", items)
}

// GetByID - Get single dokumen kebijakan by ID or slug (for edit form)
func (c *DokumenKebijakanController) GetByID(ctx *fiber.Ctx) error {
	idOrSlug := ctx.Params("id")

	// Try to determine if it's a UUID or slug
	var result *models.DokumenKebijakanDetailResponse
	var err error

	// Check if it looks like a UUID (contains hyphens and is 36 chars)
	if len(idOrSlug) == 36 && utils.IsValidUUID(idOrSlug) {
		result, err = c.service.FindByID(ctx.UserContext(), idOrSlug)
	} else {
		// Treat as slug
		result, err = c.service.FindBySlug(ctx.UserContext(), idOrSlug)
	}

	if err != nil {
		if err.Error() == "dokumen kebijakan tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}
	if isSyaratKetentuanLelang(result) && !memilikiIzin(ctx, "syarat_ketentuan_lelang:read") {
		return utils.ErrorResponse(ctx, http.StatusForbidden, "Akses syarat dan ketentuan lelang ditolak", nil)
	}

	return utils.SuccessResponse(ctx, "Detail dokumen kebijakan berhasil diambil", result)
}

// Update - Update dokumen kebijakan by ID or slug
func (c *DokumenKebijakanController) Update(ctx *fiber.Ctx) error {
	idOrSlug := ctx.Params("id")
	if utils.IsValidUUID(idOrSlug) && len(idOrSlug) == 36 {
		existing, err := c.service.FindByID(ctx.UserContext(), idOrSlug)
		if err != nil {
			if err.Error() == "dokumen kebijakan tidak ditemukan" {
				return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			}
			return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		}
		if isSyaratKetentuanLelang(existing) && !memilikiIzin(ctx, "syarat_ketentuan_lelang:manage") {
			return utils.ErrorResponse(ctx, http.StatusForbidden, "Akses syarat dan ketentuan lelang ditolak", nil)
		}
	} else if isSlugSyaratKetentuanLelang(idOrSlug) && !memilikiIzin(ctx, "syarat_ketentuan_lelang:manage") {
		return utils.ErrorResponse(ctx, http.StatusForbidden, "Akses syarat dan ketentuan lelang ditolak", nil)
	}

	var req models.UpdateDokumenKebijakanRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	// Try to determine if it's a UUID or slug
	var result *models.DokumenKebijakanDetailResponse
	var err error

	// Check if it looks like a UUID (contains hyphens and is 36 chars)
	if len(idOrSlug) == 36 && utils.IsValidUUID(idOrSlug) {
		result, err = c.service.Update(ctx.UserContext(), idOrSlug, &req)
	} else {
		// Treat as slug
		result, err = c.service.UpdateBySlug(ctx.UserContext(), idOrSlug, &req)
	}

	if err != nil {
		if err.Error() == "dokumen kebijakan tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "dokumen_kebijakan", "Dokumen kebijakan berhasil diupdate")
	return utils.SuccessResponse(ctx, "Dokumen kebijakan berhasil diupdate", result)
}

// AmbilSyaratKetentuanLelang hanya mengambil dokumen S&K lelang untuk panel.
func (c *DokumenKebijakanController) AmbilSyaratKetentuanLelang(ctx *fiber.Ctx) error {
	result, err := c.service.FindBySlug(ctx.UserContext(), "syarat-ketentuan-lelang")
	if err != nil {
		if err.Error() == "dokumen kebijakan tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Syarat dan ketentuan lelang berhasil diambil", result)
}

// PerbaruiSyaratKetentuanLelang hanya memperbarui dokumen S&K lelang.
func (c *DokumenKebijakanController) PerbaruiSyaratKetentuanLelang(ctx *fiber.Ctx) error {
	var req models.UpdateDokumenKebijakanRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.UpdateBySlug(ctx.UserContext(), "syarat-ketentuan-lelang", &req)
	if err != nil {
		if err.Error() == "dokumen kebijakan tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "dokumen_kebijakan", "Syarat dan ketentuan lelang berhasil diperbarui")
	return utils.SuccessResponse(ctx, "Syarat dan ketentuan lelang berhasil diperbarui", result)
}

func isSyaratKetentuanLelang(dokumen *models.DokumenKebijakanDetailResponse) bool {
	return dokumen != nil &&
		((dokumen.SlugID != nil && *dokumen.SlugID == "syarat-ketentuan-lelang") ||
			(dokumen.SlugEN != nil && *dokumen.SlugEN == "auction-terms-and-conditions"))
}

func isSlugSyaratKetentuanLelang(slug string) bool {
	return slug == "syarat-ketentuan-lelang" || slug == "auction-terms-and-conditions"
}

func memilikiIzin(ctx *fiber.Ctx, izin string) bool {
	permissions, ok := ctx.Locals("user_permissions").([]string)
	if !ok {
		return false
	}
	for _, permission := range permissions {
		if permission == izin {
			return true
		}
	}
	return false
}

// ========================================
// PUBLIC ENDPOINTS (Buyer/Guest)
// ========================================

// GetAllPublic - Get list of active dokumen kebijakan for public
func (c *DokumenKebijakanController) GetAllPublic(ctx *fiber.Ctx) error {
	items, err := c.service.GetActiveListPublic(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data dokumen kebijakan berhasil diambil", items)
}

// GetByIDPublic - Get single dokumen kebijakan by ID or slug for public
func (c *DokumenKebijakanController) GetByIDPublic(ctx *fiber.Ctx) error {
	idOrSlug := ctx.Params("id")
	lang := ctx.Query("lang", "id") // default to Indonesian

	// Validate lang parameter
	if lang != "id" && lang != "en" {
		lang = "id"
	}

	// Try to determine if it's a UUID or slug
	var result *models.DokumenKebijakanPublicResponse
	var err error

	// Check if it looks like a UUID (contains hyphens and is 36 chars)
	if len(idOrSlug) == 36 && utils.IsValidUUID(idOrSlug) {
		result, err = c.service.GetByIDPublic(ctx.UserContext(), idOrSlug, lang)
	} else {
		// Treat as slug
		result, err = c.service.GetBySlugPublic(ctx.UserContext(), idOrSlug, lang)
	}

	if err != nil {
		if err.Error() == "dokumen kebijakan tidak ditemukan" || err.Error() == "dokumen kebijakan tidak aktif" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail dokumen kebijakan berhasil diambil", result)
}

// GetFAQ - Get FAQ in accordion format for public
func (c *DokumenKebijakanController) GetFAQ(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "id") // default to Indonesian

	// Validate lang parameter
	if lang != "id" && lang != "en" {
		lang = "id"
	}

	result, err := c.service.GetFAQ(ctx.UserContext(), lang)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" || err.Error() == "FAQ tidak aktif" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data FAQ berhasil diambil", result)
}
