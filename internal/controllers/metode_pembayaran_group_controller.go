package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type MetodePembayaranGroupController struct {
	service        services.MetodePembayaranGroupService
	reorderService *services.ReorderService
}

func NewMetodePembayaranGroupController(service services.MetodePembayaranGroupService, reorderService *services.ReorderService) *MetodePembayaranGroupController {
	return &MetodePembayaranGroupController{
		service:        service,
		reorderService: reorderService,
	}
}

func (c *MetodePembayaranGroupController) GetAll(ctx *fiber.Ctx) error {
	items, err := c.service.GetAll(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data group metode pembayaran berhasil diambil", items)
}

func (c *MetodePembayaranGroupController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateMetodePembayaranGroupRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		if err.Error() == "Group metode pembayaran tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "Nama group sudah digunakan" {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Group metode pembayaran berhasil diupdate", result)
}

func (c *MetodePembayaranGroupController) ReorderByDirection(ctx *fiber.Ctx) error {
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
		"metode_pembayaran_group",
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
