package controllers

import (
	"fmt"
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// WMSController menyediakan endpoint admin panel untuk memverifikasi integrasi
// OAuth ke WMS (Warehouse Management System) — fondasi untuk fitur sync produk
// palet dari inventory WMS jadi cargo online.
type WMSController struct {
	service     services.WMSService
	produkRepo  repositories.ProdukRepository
	activityLog services.ActivityLogService
}

func NewWMSController(service services.WMSService, produkRepo repositories.ProdukRepository, activityLog services.ActivityLogService) *WMSController {
	return &WMSController{
		service:     service,
		produkRepo:  produkRepo,
		activityLog: activityLog,
	}
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

// CountReadyToPriceCargos memanggil GET /api/integration/cargos/ready-to-price/count
// di WMS untuk mendapatkan jumlah cargo yang siap ditetapkan harga jualnya,
// tanpa menarik seluruh isi daftar — dipakai untuk badge notifikasi.
// Hanya bisa diakses role dengan permission wms_integration:manage.
func (c *WMSController) CountReadyToPriceCargos(ctx *fiber.Ctx) error {
	ready, err := c.service.CountReadyToPriceCargos(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Jumlah cargo siap harga WMS berhasil diambil", fiber.Map{"ready": ready})
}

// SetCargoPrice memanggil POST /api/integration/cargos/{id}/price di WMS
// untuk menetapkan harga jual cargo — tipe "discount" (persentase dari
// total_price) atau "fix" (harga jual/sale_price akhir secara langsung).
// Hanya bisa diakses role dengan permission wms_integration:manage.
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

// ListAlreadyPricedCargos memanggil GET /api/integration/cargos/already-priced
// di WMS untuk mendapatkan daftar cargo yang sudah diberi harga tapi belum
// dikonfirmasi sinkron — sumber dropdown "ID Cargo" saat create produk.
func (c *WMSController) ListAlreadyPricedCargos(ctx *fiber.Ctx) error {
	search := ctx.Query("search")

	items, meta, err := c.service.ListAlreadyPricedCargos(ctx.UserContext(), search)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	paginationMeta := models.NewPaginationMeta(meta.Page, meta.Limit, meta.TotalItems)
	return utils.PaginatedSuccessResponse(ctx, "Daftar cargo sudah diberi harga berhasil diambil", items, paginationMeta)
}

// DownloadCargoPricingPDF meneruskan (proxy) PDF harga cargo dari WMS ke
// client sebagai response application/pdf mentah — FE mengunduhnya lalu
// mengunggahnya kembali sebagai dokumen produk seolah file diupload manual.
func (c *WMSController) DownloadCargoPricingPDF(ctx *fiber.Ctx) error {
	cargoID := ctx.Params("id")
	if cargoID == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID cargo tidak boleh kosong", nil)
	}

	data, err := c.service.DownloadCargoPricingPDF(ctx.UserContext(), cargoID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	ctx.Set("Content-Type", "application/pdf")
	return ctx.Send(data)
}

// MarkCargoSynced menandai cargo sudah dikonfirmasi sinkron (is_sync = true)
// di WMS setelah produk lokal berhasil dibuat dari cargo terkait. Idempotent
// — aman dipanggil berkali-kali.
func (c *WMSController) MarkCargoSynced(ctx *fiber.Ctx) error {
	cargoID := ctx.Params("id")
	if cargoID == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID cargo tidak boleh kosong", nil)
	}

	result, err := c.service.MarkCargoSynced(ctx.UserContext(), cargoID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Cargo berhasil ditandai sinkron", result)
}

// UpdateCargoActualPrice memanggil POST /api/integration/cargos/{id}/actual-price di WMS
// untuk mencatat harga jual final (actual price) cargo berdasarkan cargo ID (WMS UUID).
func (c *WMSController) UpdateCargoActualPrice(ctx *fiber.Ctx) error {
	cargoID := ctx.Params("id")
	if cargoID == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID cargo tidak boleh kosong", nil)
	}

	var req models.SetWMSCargoActualPriceRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.UpdateCargoActualPrice(ctx.UserContext(), cargoID, req.Value)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	if c.activityLog != nil {
		c.activityLog.Log(ctx, models.ActionUpdate, "wms_penjualan", fmt.Sprintf("Update penjualan WMS untuk cargo '%s' sebesar Rp%.0f", cargoID, req.Value))
	}

	return utils.SuccessResponse(ctx, "Penjualan cargo WMS berhasil diperbarui", result)
}

// UpdateProdukPenjualan mencatat harga jual final ke WMS berdasarkan produk Bulky (produk_id).
// Dipanggil oleh Storefront BE via internal API saat produk dibeli dan dibayar.
// Jika produk bukan berasal dari WMS (id_cargo dan reference_code kosong/nil), endpoint ini
// mengembalikan 200 OK dengan status dilewati (skip) agar tidak mengganggu alur checkout.
func (c *WMSController) UpdateProdukPenjualan(ctx *fiber.Ctx) error {
	var req models.InternalUpdatePenjualanProdukRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	produkID := ctx.Params("id")
	if produkID == "" {
		produkID = req.ProdukID
	}
	if produkID == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID produk tidak boleh kosong", nil)
	}

	val := req.GetValue()
	if val <= 0 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Harga penjualan (actual price / value) harus lebih besar dari 0", nil)
	}

	if c.produkRepo == nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Repository produk tidak tersedia", nil)
	}

	produk, err := c.produkRepo.FindByID(ctx.UserContext(), produkID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Produk tidak ditemukan", nil)
	}

	// Cek apakah produk bersumber dari WMS (dicek dari id_cargo atau reference_code)
	targetCargoID := ""
	if produk.IDCargo != nil && *produk.IDCargo != "" {
		targetCargoID = *produk.IDCargo
	} else if produk.ReferenceCode != nil && *produk.ReferenceCode != "" {
		targetCargoID = *produk.ReferenceCode
	}

	if targetCargoID == "" {
		return utils.SuccessResponse(ctx, "Produk bukan berasal dari WMS, pencatatan penjualan WMS dilewati", fiber.Map{
			"is_wms":       false,
			"produk_id":    produk.ID.String(),
			"actual_price": val,
		})
	}

	// Update ke WMS menggunakan targetCargoID
	result, err := c.service.UpdateCargoActualPrice(ctx.UserContext(), targetCargoID, val)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	if c.activityLog != nil {
		refCode := "-"
		if produk.ReferenceCode != nil && *produk.ReferenceCode != "" {
			refCode = *produk.ReferenceCode
		}
		c.activityLog.Log(ctx, models.ActionUpdate, "wms_penjualan", fmt.Sprintf("Update penjualan WMS untuk produk '%s' (Ref: %s / Cargo: %s) sebesar Rp%.0f", produk.NamaID, refCode, targetCargoID, val))
	}

	return utils.SuccessResponse(ctx, "Penjualan cargo WMS berhasil diperbarui", fiber.Map{
		"is_wms":                  true,
		"produk_id":               produk.ID.String(),
		"cargo_id":                result.ID,
		"code":                    result.Code,
		"sale_price":              result.SalePrice,
		"actual_price":            result.ActualPrice,
		"actual_price_updated_at": result.ActualPriceUpdatedAt,
	})
}

