package controllers

import (
	"net/http"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FAQController struct {
	service services.FAQService
}

func NewFAQController(service services.FAQService) *FAQController {
	return &FAQController{service: service}
}

// GetAll - Get all FAQ items for admin panel
// GET /api/panel/faq
func (c *FAQController) GetAll(ctx *gin.Context) {
	result, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data FAQ berhasil diambil", result)
}

// GetByID - Get FAQ by ID
// GET /api/panel/faq/:id
func (c *FAQController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	result, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data FAQ berhasil diambil", result)
}

// GetPublic - Get active FAQ for public
// GET /api/public/faq
func (c *FAQController) GetPublic(ctx *gin.Context) {
	lang := ctx.DefaultQuery("lang", "id")
	if lang != "id" && lang != "en" {
		lang = "id"
	}

	result, err := c.service.GetPublic(ctx.Request.Context(), lang)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	message := "Data FAQ berhasil diambil"
	if lang == "en" {
		message = "FAQ data retrieved successfully"
	}

	utils.SuccessResponse(ctx, message, result)
}

// Create - Create new FAQ item
// POST /api/panel/faq
func (c *FAQController) Create(ctx *gin.Context) {
	var req models.FAQCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "FAQ berhasil dibuat",
		"data":    result,
	})
}

// Update - Update FAQ item
// PUT /api/panel/faq/:id
func (c *FAQController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	var req models.FAQUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "FAQ berhasil diupdate", result)
}

// Delete - Delete FAQ item
// DELETE /api/panel/faq/:id
func (c *FAQController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	err = c.service.Delete(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "FAQ berhasil dihapus", nil)
}

// Reorder - Reorder FAQ item (up/down)
// PATCH /api/panel/faq/:id/reorder
func (c *FAQController) Reorder(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	var req dto.ReorderDirectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Reorder(ctx.Request.Context(), id, req.Direction)
	if err != nil {
		if err.Error() == "item sudah berada di urutan paling atas" ||
			err.Error() == "item sudah berada di urutan paling bawah" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Standardized reorder response
	utils.SuccessResponse(ctx, "Urutan berhasil diubah", gin.H{
		"item": gin.H{
			"id":     result.ItemID.String(),
			"urutan": result.ItemUrutan,
		},
		"swapped_with": gin.H{
			"id":     result.SwappedID.String(),
			"urutan": result.SwappedUrutan,
		},
	})
}
