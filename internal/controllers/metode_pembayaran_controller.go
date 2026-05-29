package controllers

import (
	"net/http"
	"strconv"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type MetodePembayaranController struct {
	service        services.MetodePembayaranService
	reorderService *services.ReorderService
	activityLog    services.ActivityLogService
}

func NewMetodePembayaranController(service services.MetodePembayaranService, reorderService *services.ReorderService, activityLog services.ActivityLogService) *MetodePembayaranController {
	return &MetodePembayaranController{
		service:        service,
		reorderService: reorderService,
		activityLog:    activityLog,
	}
}

// GetAllGrouped - Admin endpoint with grouped response
func (c *MetodePembayaranController) GetAllGrouped(ctx *fiber.Ctx) error {
	result, err := c.service.GetAllGrouped(ctx.UserContext(), true) // isAdmin = true
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data", nil)
	}

	return utils.SuccessResponse(ctx, "Data metode pembayaran berhasil diambil", result)
}

// GetAllGroupedPublic - Public endpoint with grouped response (active only)
func (c *MetodePembayaranController) GetAllGroupedPublic(ctx *fiber.Ctx) error {
	result, err := c.service.GetAllGrouped(ctx.UserContext(), false) // isAdmin = false
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data", nil)
	}

	return utils.SuccessResponse(ctx, "Data metode pembayaran berhasil diambil", result)
}

// ToggleMethodStatus - Toggle status of individual payment method
func (c *MetodePembayaranController) ToggleMethodStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.ToggleMethodStatus(ctx.UserContext(), id)
	if err != nil {
		if err.Error() == "Metode pembayaran tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "metode_pembayaran", "Status metode pembayaran berhasil diubah")
	return utils.SuccessResponse(ctx, "Status metode pembayaran berhasil diubah", result)
}

// ToggleGroupStatus - Toggle status of payment group by urutan
func (c *MetodePembayaranController) ToggleGroupStatus(ctx *fiber.Ctx) error {
	urutanStr := ctx.Params("urutan")
	urutan, err := strconv.Atoi(urutanStr)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Urutan tidak valid", nil)
	}

	result, err := c.service.ToggleGroupStatus(ctx.UserContext(), urutan)
	if err != nil {
		if err.Error() == "Group tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "metode_pembayaran", "Status group metode pembayaran berhasil diubah")
	return utils.SuccessResponse(ctx, "Status group berhasil diubah", result)
}

// Legacy methods - kept for backward compatibility if needed
func (c *MetodePembayaranController) GetAll(ctx *fiber.Ctx) error {
	// Parse optional filters
	var groupID *string
	if groupIDStr := ctx.Query("group_id"); groupIDStr != "" {
		groupID = &groupIDStr
	}

	var isActive *bool
	if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
		isActiveBool := isActiveStr == "true"
		isActive = &isActiveBool
	}

	items, err := c.service.GetAll(ctx.UserContext(), groupID, isActive)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data metode pembayaran berhasil diambil", items)
}

func (c *MetodePembayaranController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateMetodePembayaranRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		if err.Error() == "Metode pembayaran tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "Kode sudah digunakan" || err.Error() == "Group tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "metode_pembayaran", "Metode pembayaran berhasil diupdate")
	return utils.SuccessResponse(ctx, "Metode pembayaran berhasil diupdate", result)
}

func (c *MetodePembayaranController) ReorderByDirection(ctx *fiber.Ctx) error {
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
		"metode_pembayaran",
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
