package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AlamatBuyerController struct {
	service     services.AlamatBuyerService
	activityLog services.ActivityLogService
}

func NewAlamatBuyerController(service services.AlamatBuyerService, activityLog services.ActivityLogService) *AlamatBuyerController {
	return &AlamatBuyerController{service: service, activityLog: activityLog}
}

func (c *AlamatBuyerController) Create(ctx *fiber.Ctx) error {
	var req models.CreateAlamatBuyerRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionCreate, "alamat_buyer", "Alamat berhasil ditambahkan")
	return utils.CreatedResponse(ctx, "Alamat berhasil ditambahkan", result)
}

func (c *AlamatBuyerController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}
	return utils.SuccessResponse(ctx, "Detail alamat berhasil diambil", result)
}

func (c *AlamatBuyerController) FindAll(ctx *fiber.Ctx) error {
	var params models.AlamatBuyerFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data alamat berhasil diambil", items, *meta)
}

func (c *AlamatBuyerController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req models.UpdateAlamatBuyerRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "alamat_buyer", "Alamat berhasil diupdate")
	return utils.SuccessResponse(ctx, "Alamat berhasil diupdate", result)
}

func (c *AlamatBuyerController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "alamat tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}
	c.activityLog.Log(ctx, models.ActionDelete, "alamat_buyer", "Alamat berhasil dihapus")
	return utils.SuccessResponse(ctx, "Alamat berhasil dihapus", nil)
}

func (c *AlamatBuyerController) SetDefault(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	result, err := c.service.SetDefault(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}
	return utils.SuccessResponse(ctx, "Alamat berhasil dijadikan default", result)
}
