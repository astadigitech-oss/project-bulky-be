package controllers

import (
	"net/http"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// DelivereeWebhookController menerima event webhook dari Deliveree.
// Endpoint ini dipanggil oleh provider (inbound), bukan oleh admin,
// sehingga tidak memakai middleware auth admin — verifikasi dilakukan
// via header Authorization yang dibandingkan dengan
// DELIVEREE_WEBHOOK_AUTHORIZATION.
type DelivereeWebhookController struct {
	webhookService services.DelivereeWebhookService
	authKey        string
}

func NewDelivereeWebhookController(webhookService services.DelivereeWebhookService, authKey string) *DelivereeWebhookController {
	return &DelivereeWebhookController{
		webhookService: webhookService,
		authKey:        authKey,
	}
}

// Handle memproses webhook dari Deliveree.
//
// Payload yang diterima (JSON):
//   - status: delivery_completed | locating_driver | driver_accept_booking |
//     delivery_in_progress | locating_driver_timeout | canceled
//   - id / no_booking: booking ID Deliveree (minimal salah satu wajib)
//   - tracking_url: opsional, link tracking dari provider
//
// Selalu merespons 200 selama payload valid, bahkan jika booking ID tidak
// ditemukan (idempotent) — provider tidak perlu mengirim ulang event lama.
func (c *DelivereeWebhookController) Handle(ctx *fiber.Ctx) error {
	// Verifikasi Authorization header
	key := ctx.Get("Authorization")
	if c.authKey == "" {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false,
			"message": "Webhook Deliveree belum dikonfigurasi (DELIVEREE_WEBHOOK_AUTHORIZATION kosong)",
		})
	}
	if key != c.authKey {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	var req dto.DelivereeWebhookRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Validasi gagal", parseValidationErrors(err))
	}

	// Validasi manual: status wajib, dan id/no_booking minimal salah satu
	if req.Status == "" {
		return utils.SimpleErrorResponse(ctx, http.StatusUnprocessableEntity, "Validasi gagal", "field status wajib diisi")
	}
	if req.ID == "" && req.NoBooking == "" {
		return utils.SimpleErrorResponse(ctx, http.StatusUnprocessableEntity, "Validasi gagal", "field id atau no_booking wajib diisi")
	}

	_, err := c.webhookService.Handle(ctx.UserContext(), &req)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal memproses webhook", err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.DelivereeWebhookResponse{
		Success: true,
		Message: "Webhook diproses",
	})
}
