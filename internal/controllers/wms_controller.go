package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
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

// ListReadyToPriceCargos memanggil GET /api/integration/cargos/ready-to-price
// di WMS untuk mendapatkan daftar cargo (ukuran fisik lengkap, belum pernah
// dihargai, belum disinkronkan) yang siap ditetapkan harga jualnya.
// Hanya bisa diakses role dengan permission wms_integration:manage.
func (c *WMSController) ListReadyToPriceCargos(ctx *fiber.Ctx) error {
	var params models.WMSCargoListFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}
	params.SetDefaults()

	items, meta, err := c.service.ListReadyToPriceCargos(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	paginationMeta := models.NewPaginationMeta(meta.Page, meta.Limit, meta.TotalItems)
	return utils.PaginatedSuccessResponse(ctx, "Daftar cargo siap harga WMS berhasil diambil", items, paginationMeta)
}

// SetCargoPrice memanggil POST /api/integration/cargos/{id}/price di WMS
// untuk menetapkan harga jual cargo (diskon persentase atau potongan rupiah
// fix dari total_price). Hanya bisa diakses role dengan permission
// wms_integration:manage.
func (c *WMSController) SetCargoPrice(ctx *fiber.Ctx) error {
	cargoID := ctx.Params("id")
	if cargoID == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID cargo tidak boleh kosong", nil)
	}

	var req models.SetWMSCargoPriceRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.SetCargoPrice(ctx.UserContext(), cargoID, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Harga cargo berhasil ditetapkan", result)
}
