package controllers

import (
	"fmt"
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type DelivereeVehicleTypeController struct {
	service     services.DelivereeVehicleTypeService
	activityLog services.ActivityLogService
}

func NewDelivereeVehicleTypeController(service services.DelivereeVehicleTypeService, activityLog services.ActivityLogService) *DelivereeVehicleTypeController {
	return &DelivereeVehicleTypeController{
		service:     service,
		activityLog: activityLog,
	}
}

func (c *DelivereeVehicleTypeController) FindAll(ctx *fiber.Ctx) error {
	var params models.DelivereeVehicleTypeFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}
	params.SetDefaults()

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data kendaraan Deliveree berhasil diambil", items, *meta)
}

func (c *DelivereeVehicleTypeController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail kendaraan Deliveree berhasil diambil", result)
}

func (c *DelivereeVehicleTypeController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateDelivereeVehicleTypeRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "deliveree_vehicle_type", "Master data kendaraan Deliveree berhasil diupdate")
	return utils.SuccessResponse(ctx, "Kendaraan Deliveree berhasil diupdate", result)
}

// Sync menarik data terbaru dari API Deliveree (mengikuti DELIVEREE_BASE_URL/
// DELIVEREE_API_KEY yang aktif di environment deployment saat ini) dan
// menyimpannya sebagai master data. Hanya bisa diakses Super Admin.
func (c *DelivereeVehicleTypeController) Sync(ctx *fiber.Ctx) error {
	result, err := c.service.Sync(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionImport, "deliveree_vehicle_type",
		fmt.Sprintf("Sync master data kendaraan Deliveree (%s): %d dibuat, %d diupdate, %d dinonaktifkan", result.Environment, result.Created, result.Updated, result.Deactivated))
	return utils.SuccessResponse(ctx, "Sync kendaraan Deliveree berhasil", result)
}
