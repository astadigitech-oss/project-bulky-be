package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type KondisiProdukController struct {
	service        services.KondisiProdukService
	reorderService *services.ReorderService
}

func NewKondisiProdukController(service services.KondisiProdukService, reorderService *services.ReorderService) *KondisiProdukController {
	return &KondisiProdukController{
		service:        service,
		reorderService: reorderService,
	}
}

func (c *KondisiProdukController) Create(ctx *fiber.Ctx) error {
	var req models.CreateKondisiProdukRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
	}

	return utils.CreatedResponse(ctx, "Kondisi produk berhasil dibuat", result)
}

func (c *KondisiProdukController) FindAll(ctx *fiber.Ctx) error {
	var params models.KondisiProdukFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data kondisi produk berhasil diambil", items, *meta)
}

func (c *KondisiProdukController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail kondisi produk berhasil diambil", result)
}

func (c *KondisiProdukController) FindBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	result, err := c.service.FindBySlug(ctx.UserContext(), slug)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail kondisi produk berhasil diambil", result)
}

func (c *KondisiProdukController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateKondisiProdukRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Kondisi produk berhasil diupdate", result)
}

func (c *KondisiProdukController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "kondisi produk tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Kondisi produk berhasil dihapus", nil)
}

func (c *KondisiProdukController) ToggleStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Status kondisi produk berhasil diubah", result)
}

func (c *KondisiProdukController) Reorder(ctx *fiber.Ctx) error {
	var req models.ReorderRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	if err := c.service.Reorder(ctx.UserContext(), &req); err != nil {
		if err.Error() == "Data kondisi produk tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Urutan kondisi produk berhasil diubah", nil)
}

func (c *KondisiProdukController) ReorderByDirection(ctx *fiber.Ctx) error {
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
		"kondisi_produk",
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

func (c *KondisiProdukController) Dropdown(ctx *fiber.Ctx) error {
	response, err := c.service.GetAllForDropdown(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data kondisi", nil)
	}

	return utils.SuccessResponse(ctx, "Data dropdown kondisi produk berhasil diambil", response)
}
