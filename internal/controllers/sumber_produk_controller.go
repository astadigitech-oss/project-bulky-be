package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type SumberProdukController struct {
	service     services.SumberProdukService
	activityLog services.ActivityLogService
}

func NewSumberProdukController(service services.SumberProdukService, activityLog services.ActivityLogService) *SumberProdukController {
	return &SumberProdukController{service: service, activityLog: activityLog}
}

func (c *SumberProdukController) Create(ctx *fiber.Ctx) error {
	var req models.CreateSumberProdukRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionCreate, "sumber_produk", "Sumber produk berhasil dibuat")
	return utils.CreatedResponse(ctx, "Sumber produk berhasil dibuat", result)
}

func (c *SumberProdukController) FindAll(ctx *fiber.Ctx) error {
	var params models.SumberProdukFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data sumber produk berhasil diambil", items, *meta)
}

func (c *SumberProdukController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail sumber produk berhasil diambil", result)
}

func (c *SumberProdukController) FindBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	result, err := c.service.FindBySlug(ctx.UserContext(), slug)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail sumber produk berhasil diambil", result)
}

func (c *SumberProdukController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateSumberProdukRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "sumber_produk", "Sumber produk berhasil diupdate")
	return utils.SuccessResponse(ctx, "Sumber produk berhasil diupdate", result)
}

func (c *SumberProdukController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "sumber produk tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionDelete, "sumber_produk", "Sumber produk berhasil dihapus")
	return utils.SuccessResponse(ctx, "Sumber produk berhasil dihapus", nil)
}

func (c *SumberProdukController) ToggleStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "sumber_produk", "Status sumber produk berhasil diubah")
	return utils.SuccessResponse(ctx, "Status sumber produk berhasil diubah", result)
}

func (c *SumberProdukController) Dropdown(ctx *fiber.Ctx) error {
	var params models.SumberProdukFilterRequest
	params.Page = 1
	params.PerPage = 1000
	isActive := true
	params.IsActive = &isActive
	params.SortBy = "updated_at"
	params.Order = "asc"

	sumberList, _, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data sumber", nil)
	}

	response := make([]map[string]interface{}, len(sumberList))
	for i, s := range sumberList {
		response[i] = map[string]interface{}{
			"id":   s.ID,
			"nama": s.Nama.ID,
		}
	}

	return utils.SuccessResponse(ctx, "Data dropdown sumber produk berhasil diambil", response)
}
