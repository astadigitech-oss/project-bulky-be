package controllers

import (
	"errors"
	"net/http"
	"strings"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
