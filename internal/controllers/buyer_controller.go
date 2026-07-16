package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type BuyerController struct {
	service     services.BuyerService
	activityLog services.ActivityLogService
}

func NewBuyerController(service services.BuyerService, activityLog services.ActivityLogService) *BuyerController {
	return &BuyerController{service: service, activityLog: activityLog}
}

func (c *BuyerController) FindAll(ctx *fiber.Ctx) error {
	var params models.BuyerFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data buyer berhasil diambil", items, *meta)
}

func (c *BuyerController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}
	return utils.SuccessResponse(ctx, "Detail buyer berhasil diambil", result)
}

func (c *BuyerController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "buyer tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}
	c.activityLog.Log(ctx, models.ActionDelete, "buyer", "Buyer berhasil dihapus")
	return utils.SuccessResponse(ctx, "Buyer berhasil dihapus", nil)
}

func (c *BuyerController) ResetPassword(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req models.ResetBuyerPasswordRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	if err := c.service.ResetPassword(ctx.UserContext(), id, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "buyer", "Password buyer berhasil direset")
	return utils.SuccessResponse(ctx, "Password buyer berhasil direset", nil)
}

func (c *BuyerController) GetStatistik(ctx *fiber.Ctx) error {
	result, err := c.service.GetStatistik(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}
	return utils.SuccessResponse(ctx, "Statistik buyer berhasil diambil", result)
}

func (c *BuyerController) GetChart(ctx *fiber.Ctx) error {
	var params models.ChartParams
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	result, err := c.service.GetChart(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}
	return utils.SuccessResponse(ctx, "Data chart berhasil diambil", result)
}
