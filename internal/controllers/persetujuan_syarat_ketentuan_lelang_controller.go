package controllers

import (
	"errors"
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PersetujuanSyaratKetentuanLelangController struct {
	service services.PersetujuanSyaratKetentuanLelangService
}

func NewPersetujuanSyaratKetentuanLelangController(
	service services.PersetujuanSyaratKetentuanLelangService,
) *PersetujuanSyaratKetentuanLelangController {
	return &PersetujuanSyaratKetentuanLelangController{service: service}
}

func (c *PersetujuanSyaratKetentuanLelangController) AmbilSemua(ctx *fiber.Ctx) error {
	var params models.PaginationRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}
	items, meta, err := c.service.AmbilSemua(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data persetujuan syarat ketentuan lelang", nil)
	}
	return utils.PaginatedSuccessResponse(ctx, "Data persetujuan syarat ketentuan lelang berhasil diambil", items, *meta)
}

func (c *PersetujuanSyaratKetentuanLelangController) AmbilBerdasarkanID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID persetujuan diperlukan", nil)
	}
	item, err := c.service.AmbilBerdasarkanID(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(ctx, http.StatusNotFound, "Data persetujuan syarat ketentuan lelang tidak ditemukan", nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail persetujuan syarat ketentuan lelang", nil)
	}
	return utils.SuccessResponse(ctx, "Detail persetujuan syarat ketentuan lelang berhasil diambil", item)
}
