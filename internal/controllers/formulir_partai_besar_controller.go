package controllers

import (
	"net/http"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FormulirPartaiBesarController struct {
	service        services.FormulirPartaiBesarService
	reorderService *services.ReorderService
}

func NewFormulirPartaiBesarController(service services.FormulirPartaiBesarService, reorderService *services.ReorderService) *FormulirPartaiBesarController {
	return &FormulirPartaiBesarController{
		service:        service,
		reorderService: reorderService,
	}
}

// ========================================
// Config (Admin)
// ========================================

func (c *FormulirPartaiBesarController) GetConfig(ctx *fiber.Ctx) error {
	result, err := c.service.GetConfig(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Konfigurasi formulir berhasil diambil", result)
}

func (c *FormulirPartaiBesarController) UpdateConfig(ctx *fiber.Ctx) error {
	var req models.UpdateFormulirConfigRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.UpdateConfig(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Konfigurasi formulir berhasil diupdate", result)
}

// ========================================
// Anggaran (Admin)
// ========================================

func (c *FormulirPartaiBesarController) GetAnggaranList(ctx *fiber.Ctx) error {
	items, err := c.service.FindAllAnggaran(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data anggaran berhasil diambil", items)
}

func (c *FormulirPartaiBesarController) GetAnggaranByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindAnggaranByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail anggaran berhasil diambil", result)
}

func (c *FormulirPartaiBesarController) CreateAnggaran(ctx *fiber.Ctx) error {
	var req models.CreateAnggaranRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.CreateAnggaran(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.CreatedResponse(ctx, "Anggaran berhasil ditambahkan", result)
}

func (c *FormulirPartaiBesarController) UpdateAnggaran(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateAnggaranRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.UpdateAnggaran(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Anggaran berhasil diupdate", result)
}

func (c *FormulirPartaiBesarController) DeleteAnggaran(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.DeleteAnggaran(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "anggaran tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Anggaran berhasil dihapus", nil)
}

func (c *FormulirPartaiBesarController) ReorderAnggaran(ctx *fiber.Ctx) error {
	var req models.ReorderRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	if err := c.service.ReorderAnggaran(ctx.UserContext(), &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Urutan anggaran berhasil diubah", nil)
}

// ========================================
// Submission (Admin)
// ========================================

func (c *FormulirPartaiBesarController) GetSubmissionList(ctx *fiber.Ctx) error {
	var params models.FormulirSubmissionFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAllSubmission(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data submission berhasil diambil", items, *meta)
}

func (c *FormulirPartaiBesarController) GetSubmissionDetail(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindSubmissionByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail submission berhasil diambil", result)
}

func (c *FormulirPartaiBesarController) ResendEmail(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.ResendEmail(ctx.UserContext(), id); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Email berhasil dikirim ulang", nil)
}

// ========================================
// Buyer Endpoints
// ========================================

func (c *FormulirPartaiBesarController) GetOptions(ctx *fiber.Ctx) error {
	result, err := c.service.GetOptions(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Opsi formulir berhasil diambil", result)
}

func (c *FormulirPartaiBesarController) Submit(ctx *fiber.Ctx) error {
	var req models.CreateFormulirSubmissionRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	// Get buyer ID from context (set by auth middleware)
	var buyerID *uuid.UUID
	if buyerIDValue := ctx.Locals("buyer_id"); buyerIDValue != nil {
		if id, ok := buyerIDValue.(uuid.UUID); ok {
			buyerID = &id
		}
	}

	submissionID, err := c.service.SubmitFormulir(ctx.UserContext(), buyerID, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.CreatedResponse(ctx, "Formulir berhasil dikirim. Tim kami akan segera menghubungi Anda.", fiber.Map{
		"id": submissionID,
	})
}

func (c *FormulirPartaiBesarController) ReorderAnggaranByDirection(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.ReorderByDirectionRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
	}

	result, err := c.reorderService.Reorder(
		ctx.UserContext(),
		"formulir_partai_besar_anggaran",
		idUUID,
		req.Direction,
		"",
		nil,
	)

	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Urutan berhasil diubah", fiber.Map{
		"item": fiber.Map{
			"id":     result.ItemID,
			"urutan": result.ItemUrutan,
		},
		"swapped_with": fiber.Map{
			"id":     result.SwappedID,
			"urutan": result.SwappedUrutan,
		},
	})
}
