package controllers

import (
	"net/http"
	"strings"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PushTestController struct {
	service *services.PushTestService
}

func NewPushTestController(service *services.PushTestService) *PushTestController {
	return &PushTestController{service: service}
}

func (c *PushTestController) FindRecipients(ctx *fiber.Ctx) error {
	var params models.PaginationRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter penerima tidak valid", err.Error())
	}
	params.SetDefaults()
	params.Search = strings.TrimSpace(params.Search)
	if len(params.Search) > 100 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Pencarian maksimal 100 karakter", nil)
	}
	recipients, err := c.service.FindRecipients(ctx.UserContext(), params.Search, params.Page, params.PerPage)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadGateway, "Gagal mengambil daftar penerima tes", err.Error())
	}
	return utils.SuccessResponse(ctx, "Daftar penerima tes berhasil diambil", recipients)
}

func (c *PushTestController) Send(ctx *fiber.Ctx) error {
	var req struct {
		BuyerID  string   `json:"buyer_id,omitempty"`
		BuyerIDs []string `json:"buyer_ids,omitempty"`
	}
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Payload tes push tidak valid", parseValidationErrors(err))
	}
	if req.BuyerID != "" && len(req.BuyerIDs) > 0 {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Penerima tidak valid", "kirim buyer_id atau buyer_ids")
	}
	buyerIDs := req.BuyerIDs
	if len(buyerIDs) == 0 && req.BuyerID != "" {
		buyerIDs = []string{req.BuyerID}
	}
	if len(buyerIDs) == 0 || len(buyerIDs) > 100 {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Penerima tidak valid", "pilih antara 1 sampai 100 buyer")
	}
	validatedIDs := make([]string, 0, len(buyerIDs))
	seen := make(map[string]struct{}, len(buyerIDs))
	for _, buyerID := range buyerIDs {
		buyerID = strings.TrimSpace(buyerID)
		if _, err := uuid.Parse(buyerID); err != nil {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Buyer tidak valid", "setiap buyer_id harus berupa UUID")
		}
		if _, exists := seen[buyerID]; exists {
			continue
		}
		seen[buyerID] = struct{}{}
		validatedIDs = append(validatedIDs, buyerID)
	}
	result, err := c.service.Send(ctx.UserContext(), validatedIDs)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadGateway, "Gagal mengirim tes push notification", err.Error())
	}
	return utils.SuccessResponse(ctx, "Tes push notification diproses", result)
}
