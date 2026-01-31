package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
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
func (c *DokumenKebijakanController) GetAll(ctx *gin.Context) {
	items, err := c.service.FindAll(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data dokumen kebijakan berhasil diambil", items)
}

// GetByID - Get single dokumen kebijakan by ID or slug (for edit form)
func (c *DokumenKebijakanController) GetByID(ctx *gin.Context) {
	idOrSlug := ctx.Param("id")

	// Try to determine if it's a UUID or slug
	var result *models.DokumenKebijakanDetailResponse
	var err error

	// Check if it looks like a UUID (contains hyphens and is 36 chars)
	if len(idOrSlug) == 36 && utils.IsValidUUID(idOrSlug) {
		result, err = c.service.FindByID(ctx.Request.Context(), idOrSlug)
	} else {
		// Treat as slug
		result, err = c.service.FindBySlug(ctx.Request.Context(), idOrSlug)
	}

	if err != nil {
		if err.Error() == "dokumen kebijakan tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail dokumen kebijakan berhasil diambil", result)
}

// Update - Update dokumen kebijakan by ID or slug
func (c *DokumenKebijakanController) Update(ctx *gin.Context) {
	idOrSlug := ctx.Param("id")

	var req models.UpdateDokumenKebijakanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	// Try to determine if it's a UUID or slug
	var result *models.DokumenKebijakanDetailResponse
	var err error

	// Check if it looks like a UUID (contains hyphens and is 36 chars)
	if len(idOrSlug) == 36 && utils.IsValidUUID(idOrSlug) {
		result, err = c.service.Update(ctx.Request.Context(), idOrSlug, &req)
	} else {
		// Treat as slug
		result, err = c.service.UpdateBySlug(ctx.Request.Context(), idOrSlug, &req)
	}

	if err != nil {
		if err.Error() == "dokumen kebijakan tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Dokumen kebijakan berhasil diupdate", result)
}

// ========================================
// PUBLIC ENDPOINTS (Buyer/Guest)
// ========================================

// GetAllPublic - Get list of active dokumen kebijakan for public
func (c *DokumenKebijakanController) GetAllPublic(ctx *gin.Context) {
	items, err := c.service.GetActiveListPublic(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data dokumen kebijakan berhasil diambil", items)
}

// GetByIDPublic - Get single dokumen kebijakan by ID or slug for public
func (c *DokumenKebijakanController) GetByIDPublic(ctx *gin.Context) {
	idOrSlug := ctx.Param("id")
	lang := ctx.DefaultQuery("lang", "id") // default to Indonesian

	// Validate lang parameter
	if lang != "id" && lang != "en" {
		lang = "id"
	}

	// Try to determine if it's a UUID or slug
	var result *models.DokumenKebijakanPublicResponse
	var err error

	// Check if it looks like a UUID (contains hyphens and is 36 chars)
	if len(idOrSlug) == 36 && utils.IsValidUUID(idOrSlug) {
		result, err = c.service.GetByIDPublic(ctx.Request.Context(), idOrSlug, lang)
	} else {
		// Treat as slug
		result, err = c.service.GetBySlugPublic(ctx.Request.Context(), idOrSlug, lang)
	}

	if err != nil {
		if err.Error() == "dokumen kebijakan tidak ditemukan" || err.Error() == "dokumen kebijakan tidak aktif" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail dokumen kebijakan berhasil diambil", result)
}

// GetFAQ - Get FAQ in accordion format for public
func (c *DokumenKebijakanController) GetFAQ(ctx *gin.Context) {
	lang := ctx.DefaultQuery("lang", "id") // default to Indonesian

	// Validate lang parameter
	if lang != "id" && lang != "en" {
		lang = "id"
	}

	result, err := c.service.GetFAQ(ctx.Request.Context(), lang)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" || err.Error() == "FAQ tidak aktif" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data FAQ berhasil diambil", result)
}
