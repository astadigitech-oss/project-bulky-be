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

	// Cari di DB Bulky untuk memastikan kita mendapatkan id_cargo (UUID inventory WMS).
	// Jika input berupa reference_code / kode kargo / ID produk, kita resolve ke id_cargo aslinya.
	targetWMSID := cargoID
	if c.produkRepo != nil {
		if produk, err := c.produkRepo.FindByIdentifier(ctx.UserContext(), cargoID); err == nil && produk != nil {
			if produk.IDCargo != nil && *produk.IDCargo != "" {
				targetWMSID = *produk.IDCargo
			}
		}
	}

	result, err := c.service.UpdateCargoActualPrice(ctx.UserContext(), targetWMSID, req.Value)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	if c.activityLog != nil {
		c.activityLog.Log(ctx, models.ActionUpdate, "wms_penjualan", fmt.Sprintf("Update penjualan WMS untuk cargo '%s' sebesar Rp%.0f", targetWMSID, req.Value))
	}

	return utils.SuccessResponse(ctx, "Penjualan cargo WMS berhasil diperbarui", result)
}

// UpdateProdukPenjualan mencatat harga jual final ke WMS berdasarkan produk Bulky.
// Identifier produk dapat berupa UUID produk, id_cargo (UUID WMS), atau reference_code (kode kargo WMS).
// Database mencari kecocokan di ketiga field tersebut (id / id_cargo / reference_code).
// Untuk request ke API WMS, sistem SELALU menggunakan id_cargo (UUID inventory WMS).
// Jika produk bukan dari WMS (id_cargo kosong/nil), endpoint mengembalikan 200 OK
// dengan status dilewati (skip) agar tidak mengganggu alur checkout.
func (c *WMSController) UpdateProdukPenjualan(ctx *fiber.Ctx) error {
	var req models.InternalUpdatePenjualanProdukRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	identifier := ctx.Params("id")
	if identifier == "" {
		identifier = req.ProdukID
	}
	if identifier == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID / identifier produk tidak boleh kosong", nil)
	}

	val := req.GetValue()
	if val <= 0 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Harga penjualan (actual price / value) harus lebih besar dari 0", nil)
	}

	if c.produkRepo == nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Repository produk tidak tersedia", nil)
	}

	// Query cepat ke DB Bulky mencakup id, id_cargo, dan reference_code
	produk, err := c.produkRepo.FindByIdentifier(ctx.UserContext(), identifier)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Produk tidak ditemukan di database Bulky", nil)
	}

	// API WMS hanya menerima id_cargo (UUID inventory WMS).
	// Jika id_cargo kosong/nil, produk ini bukan produk inventory WMS.
	if produk.IDCargo == nil || *produk.IDCargo == "" {
		return utils.SuccessResponse(ctx, "Produk bukan berasal dari WMS atau belum memiliki id_cargo WMS, pencatatan penjualan WMS dilewati", fiber.Map{
			"is_wms":       false,
			"produk_id":    produk.ID.String(),
			"actual_price": val,
		})
	}

	// Kirim id_cargo (UUID inventory WMS) ke API WMS
	wmsInventoryID := *produk.IDCargo
	result, err := c.service.UpdateCargoActualPrice(ctx.UserContext(), wmsInventoryID, val)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	if c.activityLog != nil {
		refCode := "-"
		if produk.ReferenceCode != nil && *produk.ReferenceCode != "" {
			refCode = *produk.ReferenceCode
		}
		c.activityLog.Log(ctx, models.ActionUpdate, "wms_penjualan", fmt.Sprintf("Update penjualan WMS untuk produk '%s' (Ref: %s / Cargo UUID: %s) sebesar Rp%.0f", produk.NamaID, refCode, wmsInventoryID, val))
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

