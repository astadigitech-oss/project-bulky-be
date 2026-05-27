package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type DisclaimerController struct {
	service services.DisclaimerService
}

func NewDisclaimerController(service services.DisclaimerService) *DisclaimerController {
	return &DisclaimerController{service: service}
}

func (c *DisclaimerController) Create(ctx *fiber.Ctx) error {
	var req models.CreateDisclaimerRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.CreatedResponse(ctx, "Disclaimer berhasil dibuat", result)
}

func (c *DisclaimerController) FindAll(ctx *fiber.Ctx) error {
	var params models.PaginationRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data disclaimer berhasil diambil", items, *meta)
}

func (c *DisclaimerController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail disclaimer berhasil diambil", result)
}

func (c *DisclaimerController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateDisclaimerRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Disclaimer berhasil diupdate", result)
}

func (c *DisclaimerController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "disclaimer tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Disclaimer berhasil dihapus", nil)
}

func (c *DisclaimerController) SetActive(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.SetActive(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Disclaimer berhasil diaktifkan", result)
}

// Public endpoint
func (c *DisclaimerController) GetActive(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "id")

	result, err := c.service.GetActive(ctx.UserContext(), lang)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	if result == nil {
		msg := "Tidak ada disclaimer aktif"
		if lang == "en" {
			msg = "No active disclaimer"
		}
		return utils.SuccessResponse(ctx, msg, nil)
	}

	msg := "Disclaimer aktif berhasil diambil"
	if lang == "en" {
		msg = "Active disclaimer retrieved successfully"
	}
	return utils.SuccessResponse(ctx, msg, result)
}
