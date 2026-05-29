package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type PPNController struct {
	service     services.PPNService
	activityLog services.ActivityLogService
}

func NewPPNController(service services.PPNService, activityLog services.ActivityLogService) *PPNController {
	return &PPNController{service: service, activityLog: activityLog}
}

func (c *PPNController) GetAll(ctx *fiber.Ctx) error {
	var params models.PaginationRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	// Set default values
	params.SetDefaults()

	items, meta, err := c.service.GetAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data PPN berhasil diambil", items, *meta)
}

func (c *PPNController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.GetByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail PPN berhasil diambil", result)
}

func (c *PPNController) Create(ctx *fiber.Ctx) error {
	var req models.CreatePPNRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionCreate, "ppn", "PPN berhasil dibuat")
	return utils.CreatedResponse(ctx, "PPN berhasil dibuat", result)
}

func (c *PPNController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdatePPNRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "ppn", "PPN berhasil diupdate")
	return utils.SuccessResponse(ctx, "PPN berhasil diupdate", result)
}

func (c *PPNController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		if err.Error() == "Tidak dapat menghapus PPN yang sedang aktif" {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		}
		if err.Error() == "PPN tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionDelete, "ppn", "PPN berhasil dihapus")
	return utils.SuccessResponse(ctx, "PPN berhasil dihapus", nil)
}

func (c *PPNController) SetActive(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.SetActive(ctx.UserContext(), id)
	if err != nil {
		if err.Error() == "PPN tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "ppn", "PPN berhasil diaktifkan")
	return utils.SuccessResponse(ctx, "PPN berhasil diaktifkan", result)
}
