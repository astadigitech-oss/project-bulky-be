package controllers

import (
	"net/http"
	"strings"

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
	search := strings.TrimSpace(ctx.Query("search"))
	if len(search) > 100 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Pencarian maksimal 100 karakter", nil)
	}
	recipients, err := c.service.FindRecipients(ctx.UserContext(), search)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadGateway, "Gagal mengambil daftar penerima tes", err.Error())
	}
	return utils.SuccessResponse(ctx, "Daftar penerima tes berhasil diambil", recipients)
}

func (c *PushTestController) Send(ctx *fiber.Ctx) error {
	var req struct {
		BuyerID string `json:"buyer_id"`
	}
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Payload tes push tidak valid", parseValidationErrors(err))
	}
	if _, err := uuid.Parse(req.BuyerID); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Buyer tidak valid", "buyer_id harus UUID")
	}
	result, err := c.service.Send(ctx.UserContext(), req.BuyerID)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadGateway, "Gagal mengirim tes push notification", err.Error())
	}
	return utils.SuccessResponse(ctx, "Tes push notification diproses", result)
}
