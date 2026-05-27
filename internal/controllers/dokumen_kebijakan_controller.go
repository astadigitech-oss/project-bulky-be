package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type DokumenKebijakanController struct {
	service services.DokumenKebijakanService
}

func NewDokumenKebijakanController(service services.DokumenKebijakanService) *DokumenKebijakanController {
	return &DokumenKebijakanController{service: service}
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

	return utils.SuccessResponse(ctx, "Detail dokumen kebijakan berhasil diambil", result)
}

// Update - Update dokumen kebijakan by ID or slug
func (c *DokumenKebijakanController) Update(ctx *fiber.Ctx) error {
	idOrSlug := ctx.Params("id")

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

	return utils.SuccessResponse(ctx, "Dokumen kebijakan berhasil diupdate", result)
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
