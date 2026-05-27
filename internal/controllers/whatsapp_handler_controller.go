package controllers

import (
	"net/http"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type WhatsAppHandlerController struct {
	service services.WhatsAppHandlerService
}

func NewWhatsAppHandlerController(service services.WhatsAppHandlerService) *WhatsAppHandlerController {
	return &WhatsAppHandlerController{service: service}
}

// Get - Get single WhatsApp handler (no list)
func (c *WhatsAppHandlerController) Get(ctx *fiber.Ctx) error {
	result, err := c.service.Get(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data WhatsApp handler berhasil diambil", result)
}

// Update - Update WhatsApp handler (no ID in path)
func (c *WhatsAppHandlerController) Update(ctx *fiber.Ctx) error {
	var req models.UpdateWhatsAppHandlerRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), &req)
	if err != nil {
		if err.Error() == "WhatsApp handler tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "WhatsApp handler berhasil diupdate", result)
}

// GetActive - Get active WhatsApp handler for public (floating button)
func (c *WhatsAppHandlerController) GetActive(ctx *fiber.Ctx) error {
	result, err := c.service.GetActive(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	if result == nil {
		return utils.SuccessResponse(ctx, "Tidak ada WhatsApp handler aktif", nil)
	}

	return utils.SuccessResponse(ctx, "WhatsApp handler aktif berhasil diambil", result)
}
