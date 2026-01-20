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

// GetByID - Get single dokumen kebijakan by ID (for edit form)
func (c *DokumenKebijakanController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
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

// Update - Update dokumen kebijakan by ID
func (c *DokumenKebijakanController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateDokumenKebijakanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
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

// GetByIDPublic - Get single dokumen kebijakan by ID for public
func (c *DokumenKebijakanController) GetByIDPublic(ctx *gin.Context) {
	id := ctx.Param("id")
	lang := ctx.DefaultQuery("lang", "id") // default to Indonesian

	// Validate lang parameter
	if lang != "id" && lang != "en" {
		lang = "id"
	}

	result, err := c.service.GetByIDPublic(ctx.Request.Context(), id, lang)
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
