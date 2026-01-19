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

// GetBySlug - Get single dokumen kebijakan by slug (for edit form)
func (c *DokumenKebijakanController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	result, err := c.service.FindBySlug(ctx.Request.Context(), slug)
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

// Update - Update dokumen kebijakan by slug
func (c *DokumenKebijakanController) Update(ctx *gin.Context) {
	slug := ctx.Param("slug")

	var req models.UpdateDokumenKebijakanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), slug, &req)
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

// GetBySlugPublic - Get single dokumen kebijakan by slug for public
func (c *DokumenKebijakanController) GetBySlugPublic(ctx *gin.Context) {
	slug := ctx.Param("slug")

	result, err := c.service.GetBySlugPublic(ctx.Request.Context(), slug)
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
