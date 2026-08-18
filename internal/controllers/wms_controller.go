package controllers

import (
	"net/http"

	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// WMSController menyediakan endpoint admin panel untuk memverifikasi integrasi
// OAuth ke WMS (Warehouse Management System) — fondasi untuk fitur sync produk
// palet dari inventory WMS jadi cargo online.
type WMSController struct {
	service services.WMSService
}

func NewWMSController(service services.WMSService) *WMSController {
	return &WMSController{service: service}
}

// TestConnection menukar client_id/client_secret jadi access token lalu
// memanggil GET /api/integration/me untuk memverifikasi kredensial WMS aktif
// & terkoneksi. Hanya bisa diakses role dengan permission wms_integration:manage.
func (c *WMSController) TestConnection(ctx *fiber.Ctx) error {
	result, err := c.service.TestConnection(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Koneksi WMS berhasil diverifikasi", result)
}
