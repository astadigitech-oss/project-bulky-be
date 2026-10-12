package controllers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

type wmsCargoIDSyncCandidate struct {
	productID uuid.UUID
	legacyID  int64
	idCargo   string
	code      string
}

type wmsCargoIDSyncFingerprintRow struct {
	LegacyID             int64   `json:"legacy_id"`
	ProductID            string  `json:"product_id"`
	ProductName          string  `json:"product_name"`
	CurrentIDCargo       *string `json:"current_id_cargo"`
	CurrentLegacyIDCargo *int64  `json:"current_legacy_id_cargo"`
	CurrentReferenceCode *string `json:"current_reference_code"`
	IncomingID           string  `json:"incoming_id"`
	IncomingCode         string  `json:"incoming_code"`
	IDCargoConflict      bool    `json:"id_cargo_conflict"`
	LookupFailed         bool    `json:"lookup_failed"`
}

type wmsCargoIDSyncPlan struct {
	preview         models.WMSCargoIDSyncPreview
	candidates      []wmsCargoIDSyncCandidate
	fingerprintRows []wmsCargoIDSyncFingerprintRow
}

func (c *WMSController) buildCargoIDSyncPlan(ctx context.Context) (*wmsCargoIDSyncPlan, error) {
	items, err := c.service.ListCargoIDsForSync(ctx)
	if err != nil {
		return nil, err
	}

	plan := &wmsCargoIDSyncPlan{
		preview: models.WMSCargoIDSyncPreview{
			TotalFromWMS: len(items),
			Failures:     make([]models.WMSCargoIDSyncFailure, 0),
			Candidates:   make([]models.WMSCargoIDSyncCandidate, 0),
		},
		candidates:      make([]wmsCargoIDSyncCandidate, 0),
		fingerprintRows: make([]wmsCargoIDSyncFingerprintRow, 0),
	}
	sort.Slice(items, func(i, j int) bool {
		left, right := items[i], items[j]
		if left.LegacyID == nil && right.LegacyID != nil {
			return false
		}
		if left.LegacyID != nil && right.LegacyID == nil {
			return true
		}
		if left.LegacyID != nil && right.LegacyID != nil && *left.LegacyID != *right.LegacyID {
			return *left.LegacyID < *right.LegacyID
		}
		if left.ID != right.ID {
			return left.ID < right.ID
		}
		return left.Code < right.Code
	})
	legacyCounts := make(map[int64]int)
	for _, item := range items {
		if item.LegacyID != nil {
			legacyCounts[*item.LegacyID]++
		}
	}

	for _, item := range items {
		if item.LegacyID == nil {
			plan.preview.SkippedNoLegacyID++
			continue
		}

		legacyID := *item.LegacyID
		fail := func(reason string) {
			plan.preview.Failed++
			plan.preview.Failures = append(plan.preview.Failures, models.WMSCargoIDSyncFailure{
				LegacyID: legacyID,
				Code:     item.Code,
				Reason:   reason,
			})
		}

		if legacyID <= 0 {
			fail("legacy_id harus berupa angka positif")
			continue
		}
		if legacyCounts[legacyID] > 1 {
			fail("legacy_id duplikat pada respons WMS; pemetaan dilewati")
			continue
		}

		wmsID, err := uuid.Parse(strings.TrimSpace(item.ID))
		if err != nil {
			fail("id WMS bukan UUID yang valid")
			continue
		}
		code := strings.TrimSpace(item.Code)
		if code == "" || len(code) > 100 {
			fail("code WMS kosong atau melebihi 100 karakter")
			continue
		}

		products, err := c.produkRepo.FindByLegacyIDCargo(ctx, legacyID)
		if err != nil {
			plan.fingerprintRows = append(plan.fingerprintRows, wmsCargoIDSyncFingerprintRow{
				LegacyID:     legacyID,
				IncomingID:   wmsID.String(),
				IncomingCode: code,
				LookupFailed: true,
			})
			fail("gagal mencari produk berdasarkan legacy_id")
			continue
		}
		sort.Slice(products, func(i, j int) bool {
			return products[i].ID.String() < products[j].ID.String()
		})
		for _, product := range products {
			plan.fingerprintRows = append(plan.fingerprintRows, wmsCargoIDSyncFingerprintRow{
				LegacyID:             legacyID,
				ProductID:            product.ID.String(),
				ProductName:          product.NamaID,
				CurrentIDCargo:       product.IDCargo,
				CurrentLegacyIDCargo: product.LegacyIDCargo,
				CurrentReferenceCode: product.ReferenceCode,
				IncomingID:           wmsID.String(),
				IncomingCode:         code,
			})
		}

		if len(products) == 0 {
			plan.preview.NotFound++
			plan.preview.UnmatchedLegacyIDs = append(plan.preview.UnmatchedLegacyIDs, legacyID)
			continue
		}
		if len(products) > 1 {
			fail("legacy_id cocok ke lebih dari satu produk Bulky")
			continue
		}

		product := products[0]
		plan.preview.Matched++
		if product.IDCargo != nil && *product.IDCargo == strconv.FormatInt(legacyID, 10) {
			plan.preview.LegacyIDCargoMatches++
		}
		if product.LegacyIDCargo != nil && *product.LegacyIDCargo != legacyID {
			fail("produk sudah terhubung ke legacy_id_cargo yang berbeda")
			continue
		}

		productID := product.ID.String()
		idCargoConflict, err := c.produkRepo.ExistsByIDCargo(ctx, wmsID.String(), &productID)
		if err != nil {
			plan.fingerprintRows[len(plan.fingerprintRows)-1].LookupFailed = true
			fail("gagal memeriksa benturan id_cargo WMS")
			continue
		}
		plan.fingerprintRows[len(plan.fingerprintRows)-1].IDCargoConflict = idCargoConflict
		if idCargoConflict {
			fail("id_cargo WMS sudah digunakan oleh produk lain")
			continue
		}

		currentIDCargo := ""
		if product.IDCargo != nil {
			currentIDCargo = *product.IDCargo
		}
		currentReferenceCode := ""
		if product.ReferenceCode != nil {
			currentReferenceCode = *product.ReferenceCode
		}
		isAlreadyCurrent := currentIDCargo == wmsID.String() &&
			currentReferenceCode == code &&
			product.LegacyIDCargo != nil && *product.LegacyIDCargo == legacyID
		if isAlreadyCurrent {
			plan.preview.AlreadyCurrent++
			continue
		}

		plan.preview.WillUpdate++
		plan.candidates = append(plan.candidates, wmsCargoIDSyncCandidate{
			productID: product.ID,
			legacyID:  legacyID,
			idCargo:   wmsID.String(),
			code:      code,
		})
		plan.preview.Candidates = append(plan.preview.Candidates, models.WMSCargoIDSyncCandidate{
			ProductID:            product.ID.String(),
			ProductName:          product.NamaID,
			LegacyID:             legacyID,
			CurrentIDCargo:       product.IDCargo,
			CurrentReferenceCode: product.ReferenceCode,
			WMSID:                wmsID.String(),
			WMSCode:              code,
		})
	}

	fingerprint, err := json.Marshal(struct {
		Items []models.WMSCargoSyncID        `json:"items"`
		Rows  []wmsCargoIDSyncFingerprintRow `json:"rows"`
	}{Items: items, Rows: plan.fingerprintRows})
	if err != nil {
		return nil, err
	}
	token := sha256.Sum256(fingerprint)
	plan.preview.PreviewToken = hex.EncodeToString(token[:])
	return plan, nil
}

// PreviewSyncCargoIDs membaca pemetaan WMS dan menghitung hasil pencocokan
// tanpa mengubah data produk Bulky.
func (c *WMSController) PreviewSyncCargoIDs(ctx *fiber.Ctx) error {
	if c.produkRepo == nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Repository produk tidak tersedia", nil)
	}
	plan, err := c.buildCargoIDSyncPlan(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}
	return utils.SuccessResponse(ctx, "Preview re-sync ID cargo WMS berhasil dibuat", plan.preview)
}

// SyncCargoIDs mengambil ulang pemetaan cargo dari WMS untuk memvalidasi
// preview_token, lalu hanya menerapkan kandidat produk terpilih. Produk tanpa
// pasangan lokal tidak dicocokkan berdasarkan nama produk.
func (c *WMSController) SyncCargoIDs(ctx *fiber.Ctx) error {
	if c.produkRepo == nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Repository produk tidak tersedia", nil)
	}
	var request models.SyncWMSCargoIDsRequest
	if err := BindJSON(ctx, &request); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Body preview_token wajib diisi", nil)
	}
	if strings.TrimSpace(request.PreviewToken) == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "preview_token wajib diisi", nil)
	}
	if len(request.SelectedProductIDs) == 0 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Pilih minimal satu produk untuk di-re-sync", nil)
	}

	plan, err := c.buildCargoIDSyncPlan(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}
	if request.PreviewToken != plan.preview.PreviewToken {
		return utils.ErrorResponse(ctx, http.StatusConflict, "Data WMS atau produk Bulky berubah sejak preview. Muat ulang preview sebelum melanjutkan.", nil)
	}

	candidatesByProductID := make(map[uuid.UUID]wmsCargoIDSyncCandidate, len(plan.candidates))
	for _, candidate := range plan.candidates {
		candidatesByProductID[candidate.productID] = candidate
	}
	selectedCandidates := make([]wmsCargoIDSyncCandidate, 0, len(request.SelectedProductIDs))
	selectedProductIDs := make(map[uuid.UUID]struct{}, len(request.SelectedProductIDs))
	for _, productIDString := range request.SelectedProductIDs {
		productID, err := uuid.Parse(strings.TrimSpace(productIDString))
		if err != nil {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "selected_product_ids berisi ID produk yang tidak valid", nil)
		}
		if _, duplicate := selectedProductIDs[productID]; duplicate {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "selected_product_ids tidak boleh berisi ID duplikat", nil)
		}
		candidate, exists := candidatesByProductID[productID]
		if !exists {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "Produk yang dipilih bukan kandidat update pada preview ini", nil)
		}
		selectedProductIDs[productID] = struct{}{}
		selectedCandidates = append(selectedCandidates, candidate)
	}

	result := models.WMSCargoIDSyncResult{
		TotalFromWMS:         plan.preview.TotalFromWMS,
		Matched:              plan.preview.Matched,
		LegacyIDCargoMatches: plan.preview.LegacyIDCargoMatches,
		AlreadyCurrent:       plan.preview.AlreadyCurrent,
		Selected:             len(selectedCandidates),
		NotSelected:          len(plan.candidates) - len(selectedCandidates),
		SkippedNoLegacyID:    plan.preview.SkippedNoLegacyID,
		NotFound:             plan.preview.NotFound,
		UnmatchedLegacyIDs:   plan.preview.UnmatchedLegacyIDs,
		Failed:               plan.preview.Failed,
		Failures:             append([]models.WMSCargoIDSyncFailure(nil), plan.preview.Failures...),
	}
	for _, candidate := range selectedCandidates {
		if err := c.produkRepo.SyncCargoID(ctx.UserContext(), candidate.productID, candidate.legacyID, candidate.idCargo, candidate.code); err != nil {
			result.Failed++
			result.Failures = append(result.Failures, models.WMSCargoIDSyncFailure{
				LegacyID: candidate.legacyID,
				Code:     candidate.code,
				Reason:   "gagal menyimpan id_cargo dan reference_code: " + err.Error(),
			})
			continue
		}
		result.Updated++
	}

	if c.activityLog != nil {
		c.activityLog.Log(ctx, models.ActionUpdate, "wms_cargo_id_sync", fmt.Sprintf(
			"Re-sync ID cargo WMS selesai: selected=%d, updated=%d, not_selected=%d, no_legacy_id=%d, not_found=%d, failed=%d",
			result.Selected, result.Updated, result.NotSelected, result.SkippedNoLegacyID, result.NotFound, result.Failed,
		))
	}

	message := "Re-sync ID cargo WMS selesai"
	if result.NotFound > 0 || result.Failed > 0 {
		message = "Re-sync ID cargo WMS selesai dengan sebagian cargo tidak diperbarui"
	}
	return utils.SuccessResponse(ctx, message, result)
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

// DownloadProdukPricingPDF mengambil PDF harga (pricing PDF) terbaru dari WMS
// untuk sebuah produk pada panel admin — dipakai di halaman edit produk.
//
// ID yang dikirim ke API WMS adalah nilai kolom `id_cargo` (UUID inventory WMS)
// pada produk tersebut, sesuai komentar di form create/edit produk.
//
// Perilaku status:
//   - Produk tidak ditemukan di DB Bulky              → 404 "Produk tidak ditemukan"
//   - Produk diinsert manual (id_cargo kosong/nil)    → 404 (bukan dari sync WMS)
//   - Cargo tidak ditemukan di WMS (status 404)       → 404 "Cargo tidak ditemukan di WMS"
//   - Error lain dari WMS                             → 502 Bad Gateway
func (c *WMSController) DownloadProdukPricingPDF(ctx *fiber.Ctx) error {
	identifier := ctx.Params("id")
	if identifier == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID produk tidak boleh kosong", nil)
	}

	if c.produkRepo == nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Repository produk tidak tersedia", nil)
	}

	// Query cepat ke DB Bulky mencakup id, id_cargo, dan reference_code.
	produk, err := c.produkRepo.FindByIdentifier(ctx.UserContext(), identifier)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Produk tidak ditemukan di database Bulky", nil)
	}

	// API WMS hanya menerima id_cargo (UUID inventory WMS). Jika id_cargo
	// kosong/nil, produk ini bukan produk inventory WMS (diinsert manual).
	if produk.IDCargo == nil || *produk.IDCargo == "" {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Produk ini diinsert manual, tidak memiliki id_cargo WMS (bukan berasal dari sinkronisasi WMS)", nil)
	}

	data, err := c.service.DownloadCargoPricingPDF(ctx.UserContext(), *produk.IDCargo)
	if err != nil {
		if errors.Is(err, services.ErrWMSCargoNotFound) {
			return utils.ErrorResponse(ctx, http.StatusNotFound, "Cargo tidak ditemukan di WMS", nil)
		}
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	ctx.Set("Content-Type", "application/pdf")
	return ctx.Send(data)
}

// RefreshProdukPricingPDF menetapkan ulang harga jual cargo WMS (type "fix",
// value = harga_sesudah_diskon produk) lalu mengembalikan PDF harga terbaru.
//
// Dipakai tombol "Ambil PDF terbaru" di halaman edit produk: karena WMS hanya
// me-render PDF saat harga di-set (termasuk diset ulang), GET pricing-pdf biasa
// hanya mengembalikan PDF terakhir yang tersimpan. Endpoint ini memaksa WMS
// membuat PDF baru yang mencerminkan harga saat ini, lalu FE memperlakukannya
// sebagai dokumen produk (diunggah ulang lewat dokumen[] saat save).
//
// Nilai harga diambil dari body (live form). Kalau body kosong / value <= 0,
// fallback ke harga produk di DB (produk.HargaSesudahDiskon).
//
// Perilaku status:
//   - Produk tidak ditemukan                              → 404
//   - Produk manual (id_cargo kosong)                     → 404
//   - value tidak valid (<= 0)                            → 400
//   - SetCargoPrice gagal (mis. cargo terjual/terkunci)  → 502 (pesan WMS)
//   - Download PDF gagal (cargo 404 di WMS)               → 404
func (c *WMSController) RefreshProdukPricingPDF(ctx *fiber.Ctx) error {
	identifier := ctx.Params("id")
	if identifier == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID produk tidak boleh kosong", nil)
	}

	var req models.RefreshProdukPricingPDFRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	if c.produkRepo == nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Repository produk tidak tersedia", nil)
	}

	produk, err := c.produkRepo.FindByIdentifier(ctx.UserContext(), identifier)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Produk tidak ditemukan di database Bulky", nil)
	}

	// API WMS hanya menerima id_cargo (UUID inventory WMS). Jika id_cargo
	// kosong/nil, produk ini bukan produk inventory WMS (diinsert manual).
	if produk.IDCargo == nil || *produk.IDCargo == "" {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Produk ini diinsert manual, tidak memiliki id_cargo WMS (bukan berasal dari sinkronisasi WMS)", nil)
	}

	// Nilai harga: dari live form (req.Value) atau fallback ke harga produk di DB.
	value := req.Value
	if value <= 0 {
		value = produk.HargaSesudahDiskon
	}
	if value <= 0 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Harga jual (harga sesudah diskon) harus lebih besar dari 0", nil)
	}

	// Set ulang harga cargo → WMS me-render PDF baru yang mencerminkan harga ini.
	// type "fix" berarti sale_price = value (bukan potongan dari total_price).
	priceResult, err := c.service.SetCargoPrice(ctx.UserContext(), *produk.IDCargo, &models.SetWMSCargoPriceRequest{
		Type:  "fix",
		Value: value,
	})
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	// Ambil PDF terbaru hasil render ulang. Ikuti pricing_pdf_url yang
	// dikembalikan WMS dari respons POST price; fallback ke path hasil
	// konstruksi dari id_cargo kalau url kosong.
	var data []byte
	if priceResult.PricingPDFURL != "" {
		data, err = c.service.DownloadCargoPricingPDFByURL(ctx.UserContext(), priceResult.PricingPDFURL)
	} else {
		data, err = c.service.DownloadCargoPricingPDF(ctx.UserContext(), *produk.IDCargo)
	}
	if err != nil {
		if errors.Is(err, services.ErrWMSCargoNotFound) {
			return utils.ErrorResponse(ctx, http.StatusNotFound, "Cargo tidak ditemukan di WMS", nil)
		}
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	if c.activityLog != nil {
		c.activityLog.Log(ctx, models.ActionUpdate, "wms_pricing_pdf", fmt.Sprintf("Generate ulang PDF harga WMS untuk produk '%s' sebesar Rp%.0f", produk.NamaID, value))
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
