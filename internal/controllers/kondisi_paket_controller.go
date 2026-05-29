package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type KondisiPaketController struct {
	service        services.KondisiPaketService
	reorderService *services.ReorderService
	activityLog    services.ActivityLogService
}

func NewKondisiPaketController(service services.KondisiPaketService, reorderService *services.ReorderService, activityLog services.ActivityLogService) *KondisiPaketController {
	return &KondisiPaketController{
		service:        service,
		reorderService: reorderService,
		activityLog:    activityLog,
	}
}

func (c *KondisiPaketController) Create(ctx *fiber.Ctx) error {
	var req models.CreateKondisiPaketRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionCreate, "kondisi_paket", "Kondisi paket berhasil dibuat")
	return utils.CreatedResponse(ctx, "Kondisi paket berhasil dibuat", result)
}

func (c *KondisiPaketController) FindAll(ctx *fiber.Ctx) error {
	var params models.PaginationRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data kondisi paket berhasil diambil", items, *meta)
}

func (c *KondisiPaketController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail kondisi paket berhasil diambil", result)
}

func (c *KondisiPaketController) FindBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	result, err := c.service.FindBySlug(ctx.UserContext(), slug)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail kondisi paket berhasil diambil", result)
}

func (c *KondisiPaketController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateKondisiPaketRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "kondisi_paket", "Kondisi paket berhasil diupdate")
	return utils.SuccessResponse(ctx, "Kondisi paket berhasil diupdate", result)
}

func (c *KondisiPaketController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "kondisi paket tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionDelete, "kondisi_paket", "Kondisi paket berhasil dihapus")
	return utils.SuccessResponse(ctx, "Kondisi paket berhasil dihapus", nil)
}

func (c *KondisiPaketController) ToggleStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "kondisi_paket", "Status kondisi paket berhasil diubah")
	return utils.SuccessResponse(ctx, "Status kondisi paket berhasil diubah", result)
}

func (c *KondisiPaketController) Reorder(ctx *fiber.Ctx) error {
	var req models.ReorderRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	if err := c.service.Reorder(ctx.UserContext(), &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Urutan kondisi paket berhasil diubah", nil)
}

func (c *KondisiPaketController) ReorderByDirection(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.ReorderByDirectionRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	idUUID, err := utils.ParseUUID(id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
	}

	result, err := c.reorderService.Reorder(
		ctx.UserContext(),
		"kondisi_paket",
		idUUID,
		req.Direction,
		"",
		nil,
	)

	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return 	utils.SuccessResponse(ctx, "Urutan berhasil diubah", fiber.Map{
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

func (c *KondisiPaketController) Dropdown(ctx *fiber.Ctx) error {
	response, err := c.service.GetAllForDropdown(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data kondisi paket", nil)
	}

	return utils.SuccessResponse(ctx, "Data dropdown kondisi paket berhasil diambil", response)
}
