package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type AuctionController struct {
	service services.AuctionService
}

func NewAuctionController(service services.AuctionService) *AuctionController {
	return &AuctionController{service: service}
}

// helper adminID dari context.
func auctionAdminID(c *fiber.Ctx) (uuid.UUID, bool) {
	adminIDStr := localsString(c, "admin_id")
	if adminIDStr == "" {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(adminIDStr)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

// helper idempotency key dari header.
func idempotencyKey(c *fiber.Ctx) string {
	key := c.Get("Idempotency-Key")
	return strings.TrimSpace(key)
}

// handleAuctionError memetakan AuctionError ke response HTTP sesuai status.
func handleAuctionError(c *fiber.Ctx, err error) error {
	var ae *services.AuctionError
	if errors.As(err, &ae) {
		return utils.ErrorResponse(c, ae.Status, ae.Message, ae.FieldErrors)
	}
	return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
}

// ============================================================
// Batch
// ============================================================

func (ctrl *AuctionController) List(c *fiber.Ctx) error {
	var params dto.AuctionListQueryParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
	}
	params.SetDefaults()

	items, meta, err := ctrl.service.ListBatches(c.UserContext(), &params)
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.PaginatedSuccessResponse(c, "Data batch berhasil diambil", items, *meta)
}

func (ctrl *AuctionController) Create(c *fiber.Ctx) error {
	adminID, ok := auctionAdminID(c)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Admin tidak valid", nil)
	}

	var req dto.AuctionDraftInput
	if err := BindJSON(c, &req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := ctrl.service.CreateDraft(c.UserContext(), &req, adminID, idempotencyKey(c))
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.CreatedResponse(c, "Draft batch berhasil dibuat", result)
}

func (ctrl *AuctionController) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid", nil)
	}

	result, err := ctrl.service.GetBatchDetail(c.UserContext(), id)
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.SuccessResponse(c, "Detail batch berhasil diambil", result)
}

func (ctrl *AuctionController) Update(c *fiber.Ctx) error {
	adminID, ok := auctionAdminID(c)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Admin tidak valid", nil)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid", nil)
	}

	var req dto.AuctionDraftInput
	if err := BindJSON(c, &req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := ctrl.service.UpdateDraft(c.UserContext(), id, &req, adminID, req.Version)
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.SuccessResponse(c, "Draft batch berhasil diperbarui", result)
}

func (ctrl *AuctionController) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid", nil)
	}

	if err := ctrl.service.DeleteDraft(c.UserContext(), id); err != nil {
		return handleAuctionError(c, err)
	}
	return utils.SuccessResponse(c, "Draft batch berhasil dihapus", nil)
}

func (ctrl *AuctionController) Publish(c *fiber.Ctx) error {
	adminID, ok := auctionAdminID(c)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Admin tidak valid", nil)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid", nil)
	}

	var req dto.AuctionPublishRequest
	if err := BindJSON(c, &req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := ctrl.service.Publish(c.UserContext(), id, req.Version, adminID, idempotencyKey(c))
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.SuccessResponse(c, "Batch berhasil dibuka untuk lelang", result)
}

func (ctrl *AuctionController) ListBids(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid", nil)
	}

	var params dto.AuctionBidsQueryParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
	}
	params.SetDefaults()

	items, meta, err := ctrl.service.ListBids(c.UserContext(), id, &params)
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.PaginatedSuccessResponse(c, "Data bid berhasil diambil", items, *meta)
}

func (ctrl *AuctionController) ExportBids(c *fiber.Ctx) error {
	var params dto.AuctionBidsExportQueryParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
	}
	params.SetDefaults()

	items, err := ctrl.service.ExportBids(c.UserContext(), &params)
	if err != nil {
		return handleAuctionError(c, err)
	}

	f := excelize.NewFile()
	defer f.Close()
	sheet := "Bid"
	f.SetSheetName("Sheet1", sheet)
	headers := []string{"No", "Batch ID", "Buyer", "Telepon", "Bid ke-", "Mode", "Persentase Input", "Nominal Bid", "Ongkir", "PPN", "Total Estimasi", "Catatan", "Waktu", "Pemenang"}
	for column, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(column+1, 1)
		_ = f.SetCellValue(sheet, cell, header)
	}
	for index, item := range items {
		row := index + 2
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), index+1)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.BatchID)
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.Buyer.Nama)
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Buyer.Telepon)
		_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.Sequence)
		_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.InputMode)
		_ = f.SetCellValue(sheet, fmt.Sprintf("G%d", row), nullableAuctionPercent(item.InputPercent))
		_ = f.SetCellValue(sheet, fmt.Sprintf("H%d", row), item.Amount)
		_ = f.SetCellValue(sheet, fmt.Sprintf("I%d", row), item.ShippingAmountSnapshot)
		_ = f.SetCellValue(sheet, fmt.Sprintf("J%d", row), item.PPNAmountSnapshot)
		_ = f.SetCellValue(sheet, fmt.Sprintf("K%d", row), item.EstimatedTotalSnapshot)
		_ = f.SetCellValue(sheet, fmt.Sprintf("L%d", row), item.Note)
		_ = f.SetCellValue(sheet, fmt.Sprintf("M%d", row), formatAuctionExportTime(item.CreatedAt))
		winner := "Tidak"
		if item.IsSelected {
			winner = "Ya"
		}
		_ = f.SetCellValue(sheet, fmt.Sprintf("N%d", row), winner)
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"bid-lelang-bulky-%s.xlsx\"", time.Now().Format("20060102")))
	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat file Excel", nil)
	}
	return nil
}

func formatAuctionExportTime(value time.Time) string {
	wib := time.FixedZone("WIB", 7*60*60)
	return value.In(wib).Format("02 Jan 2006 15:04 WIB")
}

func nullableAuctionPercent(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (ctrl *AuctionController) SelectWinner(c *fiber.Ctx) error {
	adminID, ok := auctionAdminID(c)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Admin tidak valid", nil)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid", nil)
	}

	var req dto.AuctionWinnerRequest
	if err := BindJSON(c, &req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := ctrl.service.SelectWinner(c.UserContext(), id, &req, adminID, idempotencyKey(c))
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.SuccessResponse(c, "Pemenang berhasil dipilih", result)
}

func (ctrl *AuctionController) UpdateOperations(c *fiber.Ctx) error {
	adminID, ok := auctionAdminID(c)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Admin tidak valid", nil)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid", nil)
	}

	var req dto.AuctionOperationRequest
	if err := BindJSON(c, &req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := ctrl.service.UpdateOperations(c.UserContext(), id, &req, adminID, idempotencyKey(c))
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.SuccessResponse(c, "Status operasional berhasil diperbarui", result)
}

// ============================================================
// Assets
// ============================================================

func (ctrl *AuctionController) UploadAsset(c *fiber.Ctx) error {
	adminID, ok := auctionAdminID(c)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Admin tidak valid", nil)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "File tidak ditemukan", nil)
	}

	kind := c.FormValue("kind")
	if kind == "" {
		kind = "IMAGE"
	}

	result, err := ctrl.service.UploadAsset(c.UserContext(), file, kind, adminID)
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.CreatedResponse(c, "Aset berhasil diupload", result)
}

func (ctrl *AuctionController) PreviewSupplierExcel(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "File Excel tidak ditemukan", nil)
	}
	result, err := ctrl.service.PreviewSupplierExcel(c.UserContext(), file)
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.SuccessResponse(c, "Kolom Excel berhasil dibaca", result)
}

func (ctrl *AuctionController) ImportSupplierExcel(c *fiber.Ctx) error {
	adminID, ok := auctionAdminID(c)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Admin tidak valid", nil)
	}
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "File Excel tidak ditemukan", nil)
	}
	parseColumn := func(name string) (int, error) {
		value, parseErr := strconv.Atoi(c.FormValue(name))
		if parseErr != nil || value < 0 {
			return 0, fmt.Errorf("Kolom %s tidak valid", name)
		}
		return value, nil
	}
	nameColumn, err := parseColumn("name_column")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	priceColumn, err := parseColumn("price_column")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	quantityColumn, err := parseColumn("quantity_column")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	headerRow, err := parseColumn("header_row")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	result, err := ctrl.service.ImportSupplierExcel(c.UserContext(), file, dto.AuctionSupplierExcelMapping{
		NameColumn: nameColumn, PriceColumn: priceColumn, QuantityColumn: quantityColumn, HeaderRow: headerRow,
	}, c.FormValue("title"), adminID)
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.CreatedResponse(c, "Item supplier berhasil diimpor dan PDF disimpan", result)
}

// ============================================================
// Product Options
// ============================================================

func (ctrl *AuctionController) ListProductOptions(c *fiber.Ctx) error {
	var params dto.AuctionProductOptionsQueryParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
	}
	params.SetDefaults()

	items, meta, err := ctrl.service.ListProductOptions(c.UserContext(), &params)
	if err != nil {
		return handleAuctionError(c, err)
	}
	return utils.PaginatedSuccessResponse(c, "Kandidat produk berhasil diambil", items, *meta)
}
