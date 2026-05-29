package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type DiskonKategoriController struct {
	service     services.DiskonKategoriService
	activityLog services.ActivityLogService
}

func NewDiskonKategoriController(service services.DiskonKategoriService, activityLog services.ActivityLogService) *DiskonKategoriController {
	return &DiskonKategoriController{service: service, activityLog: activityLog}
}

func (c *DiskonKategoriController) Create(ctx *fiber.Ctx) error {
	var req models.CreateDiskonKategoriRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionCreate, "diskon_kategori", "Diskon kategori berhasil dibuat")
	return utils.CreatedResponse(ctx, "Diskon kategori berhasil dibuat", result)
}

func (c *DiskonKategoriController) FindAll(ctx *fiber.Ctx) error {
	var params models.PaginationRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	kategoriID := ctx.Query("kategori_id")
	berlakuHariIni := ctx.Query("berlaku_hari_ini") == "true"

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params, kategoriID, berlakuHariIni)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data diskon kategori berhasil diambil", items, *meta)
}


func (c *DiskonKategoriController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail diskon kategori berhasil diambil", result)
}

func (c *DiskonKategoriController) FindActiveByKategoriID(ctx *fiber.Ctx) error {
	kategoriID := ctx.Params("kategori_id")

	result, err := c.service.FindActiveByKategoriID(ctx.UserContext(), kategoriID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	if result == nil {
		return utils.SuccessResponse(ctx, "Tidak ada diskon aktif untuk kategori ini", nil)
	}

	return utils.SuccessResponse(ctx, "Diskon kategori ditemukan", result)
}

func (c *DiskonKategoriController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateDiskonKategoriRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "diskon_kategori", "Diskon kategori berhasil diupdate")
	return utils.SuccessResponse(ctx, "Diskon kategori berhasil diupdate", result)
}

func (c *DiskonKategoriController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "diskon kategori tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionDelete, "diskon_kategori", "Diskon kategori berhasil dihapus")
	return utils.SuccessResponse(ctx, "Diskon kategori berhasil dihapus", nil)
}

func (c *DiskonKategoriController) ToggleStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "diskon_kategori", "Status diskon kategori berhasil diubah")
	return utils.SuccessResponse(ctx, "Status diskon kategori berhasil diubah", result)
}
