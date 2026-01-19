package controllers

import (
	"net/http"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type WhatsAppHandlerController struct {
	service services.WhatsAppHandlerService
}

func NewWhatsAppHandlerController(service services.WhatsAppHandlerService) *WhatsAppHandlerController {
	return &WhatsAppHandlerController{service: service}
}

// GetAll - Get all WhatsApp handlers (no pagination)
func (c *WhatsAppHandlerController) FindAll(ctx *gin.Context) {
	items, err := c.service.FindAll(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data WhatsApp handler berhasil diambil", items)
}

// Update - Update WhatsApp handler by ID
func (c *WhatsAppHandlerController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateWhatsAppHandlerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "WhatsApp handler tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "WhatsApp handler berhasil diupdate", result)
}

// GetActive - Get active WhatsApp handler for public (floating button)
func (c *WhatsAppHandlerController) GetActive(ctx *gin.Context) {
	result, err := c.service.GetActive(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if result == nil {
		utils.SuccessResponse(ctx, "Tidak ada WhatsApp handler aktif", nil)
		return
	}

	utils.SuccessResponse(ctx, "WhatsApp handler aktif berhasil diambil", result)
}
