package controllers

import (
	"log"
	"net/http"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// ForwarderWebhookController menerima event webhook dari Forwarder.
// Endpoint ini dipanggil oleh provider (inbound), bukan oleh admin,
// sehingga tidak memakai middleware auth admin — verifikasi dilakukan
// via header Authorization yang dibandingkan dengan
// FORWARDER_WEBHOOK_AUTHORIZATION.
type ForwarderWebhookController struct {
	webhookService services.ForwarderWebhookService
	authKey        string
}

func NewForwarderWebhookController(webhookService services.ForwarderWebhookService, authKey string) *ForwarderWebhookController {
	return &ForwarderWebhookController{
		webhookService: webhookService,
		authKey:        authKey,
	}
}

// Handle memproses webhook dari Forwarder.
//
// Payload yang diterima (JSON) sama dengan webhook Deliveree karena kedua
// provider berbagi platform on-demand yang sama:
//   - status: delivery_completed | locating_driver | driver_accept_booking |
//     delivery_in_progress | locating_driver_timeout | canceled
//   - id / no_booking: identifier booking/tracking Forwarder (minimal salah satu
//     wajib), dicocokkan dengan kolom forwarder_tracking_no pada pesanan
//   - tracking_url: opsional, link tracking dari provider
//
// Selalu merespons 200 selama payload valid, bahkan jika identifier tidak
// ditemukan (idempotent) — provider tidak perlu mengirim ulang event lama.
func (c *ForwarderWebhookController) Handle(ctx *fiber.Ctx) error {
	// Verifikasi Authorization header
	key := ctx.Get("Authorization")
	if c.authKey == "" {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false,
			"message": "Webhook Forwarder belum dikonfigurasi (FORWARDER_WEBHOOK_AUTHORIZATION kosong)",
		})
	}
	if key != c.authKey {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	var req dto.ForwarderWebhookRequest
	if err := BindJSON(ctx, &req); err != nil {
		log.Printf("[forwarder-webhook] parse body gagal: err=%v raw=%s", err, string(ctx.Body()))
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

	return ctx.Status(fiber.StatusOK).JSON(dto.ForwarderWebhookResponse{
		Success: true,
		Message: "Webhook diproses",
	})
}
