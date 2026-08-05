package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type BuyerDisclaimerConsentController struct {
	service services.BuyerDisclaimerConsentService
}

func NewBuyerDisclaimerConsentController(service services.BuyerDisclaimerConsentService) *BuyerDisclaimerConsentController {
	return &BuyerDisclaimerConsentController{service: service}
}

// GetAllConsents — GET /api/v1/panel/disclaimer-consent
// Admin: seluruh audit log persetujuan disclaimer (paginated).
func (c *BuyerDisclaimerConsentController) GetAllConsents(ctx *fiber.Ctx) error {
	var params models.PaginationRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.GetAllConsents(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data audit disclaimer consent berhasil diambil", items, *meta)
}

// GetConsentByPesanan — GET /api/v1/panel/disclaimer-consent/:id
// Admin: detail consent berdasarkan pesanan_id.
func (c *BuyerDisclaimerConsentController) GetConsentByPesanan(ctx *fiber.Ctx) error {
	pesananID := ctx.Params("id")
	if pesananID == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "id diperlukan", nil)
	}

	result, err := c.service.GetConsentByPesanan(ctx.UserContext(), pesananID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail persetujuan disclaimer berhasil diambil", result)
}
