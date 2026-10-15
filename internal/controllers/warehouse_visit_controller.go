package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WarehouseVisitController struct {
	service     services.WarehouseVisitService
	activityLog services.ActivityLogService
}

func NewWarehouseVisitController(service services.WarehouseVisitService, activityLog services.ActivityLogService) *WarehouseVisitController {
	return &WarehouseVisitController{service: service, activityLog: activityLog}
}

func (c *WarehouseVisitController) GetPublicAvailability(ctx *fiber.Ctx) error {
	result, err := c.service.ListPublicAvailability(ctx.UserContext(), ctx.Query("from"), ctx.Query("to"))
	if err != nil {
		return warehouseVisitErrorResponse(ctx, err)
	}
	return utils.SuccessResponse(ctx, "Ketersediaan kunjungan gudang berhasil diambil", result)
}

func (c *WarehouseVisitController) GetDates(ctx *fiber.Ctx) error {
	result, err := c.service.ListAdminDates(ctx.UserContext(), ctx.Query("from"), ctx.Query("to"))
	if err != nil {
		return warehouseVisitErrorResponse(ctx, err)
	}
	return utils.SuccessResponse(ctx, "Tanggal kunjungan gudang berhasil diambil", result)
}

func (c *WarehouseVisitController) CreateDate(ctx *fiber.Ctx) error {
	var req models.WarehouseVisitDateInput
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}
	adminID, err := warehouseVisitAdminID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusUnauthorized, "Identitas Admin tidak tersedia", nil)
	}
	result, err := c.service.CreateDate(ctx.UserContext(), &req, adminID)
	if err != nil {
		return warehouseVisitErrorResponse(ctx, err)
	}
	c.activityLog.LogCreate(ctx, "warehouse_visit", "warehouse_visit_date", result.ID, "Tanggal kunjungan gudang dibuka", result)
	return utils.CreatedResponse(ctx, "Tanggal kunjungan gudang berhasil dibuat", result)
}

func (c *WarehouseVisitController) UpdateDate(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tanggal tidak valid", nil)
	}
	var req models.WarehouseVisitDateUpdateInput
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}
	adminID, err := warehouseVisitAdminID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusUnauthorized, "Identitas Admin tidak tersedia", nil)
	}
	previous, result, err := c.service.UpdateDate(ctx.UserContext(), id, &req, adminID)
	if err != nil {
		return warehouseVisitErrorResponse(ctx, err)
	}
	c.activityLog.LogUpdate(ctx, "warehouse_visit", "warehouse_visit_date", result.ID, "Pengaturan tanggal kunjungan gudang diperbarui", previous, result)
	return utils.SuccessResponse(ctx, "Tanggal kunjungan gudang berhasil diperbarui", result)
}

func (c *WarehouseVisitController) DeleteDate(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tanggal tidak valid", nil)
	}
	result, err := c.service.DeleteDate(ctx.UserContext(), id)
	if err != nil {
		return warehouseVisitErrorResponse(ctx, err)
	}
	c.activityLog.LogDelete(ctx, "warehouse_visit", "warehouse_visit_date", result.ID, "Jadwal kunjungan gudang dihapus", result)
	return utils.SuccessResponse(ctx, "Jadwal kunjungan gudang berhasil dihapus", result)
}

func (c *WarehouseVisitController) GetBookings(ctx *fiber.Ctx) error {
	var req models.WarehouseVisitBookingListRequest
	if err := ctx.QueryParser(&req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}
	result, meta, err := c.service.ListBookings(ctx.UserContext(), &req)
	if err != nil {
		return warehouseVisitErrorResponse(ctx, err)
	}
	return utils.PaginatedSuccessResponse(ctx, "Booking kunjungan gudang berhasil diambil", result, *meta)
}

func (c *WarehouseVisitController) RescheduleBooking(ctx *fiber.Ctx) error {
	bookingID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID booking tidak valid", nil)
	}
	var req models.WarehouseVisitRescheduleRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}
	targetDateID, err := uuid.Parse(req.VisitDateID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tanggal tujuan tidak valid", nil)
	}
	previous, result, err := c.service.RescheduleBooking(ctx.UserContext(), bookingID, targetDateID)
	if err != nil {
		return warehouseVisitErrorResponse(ctx, err)
	}
	c.activityLog.LogUpdate(ctx, "warehouse_visit", "warehouse_visit_booking", result.ID, "Booking kunjungan gudang dipindahkan", previous, result)
	return utils.SuccessResponse(ctx, "Booking berhasil dipindahkan", result)
}

func (c *WarehouseVisitController) CancelBooking(ctx *fiber.Ctx) error {
	bookingID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID booking tidak valid", nil)
	}
	previous, result, err := c.service.CancelBooking(ctx.UserContext(), bookingID)
	if err != nil {
		return warehouseVisitErrorResponse(ctx, err)
	}
	oldData, _ := json.Marshal(previous)
	newData, _ := json.Marshal(result)
	c.activityLog.Log(ctx, models.ActionCancel, "warehouse_visit", "Booking kunjungan gudang dibatalkan",
		services.WithEntity("warehouse_visit_booking", result.ID),
		services.WithOldData(oldData),
		services.WithNewData(newData),
	)
	return utils.SuccessResponse(ctx, "Booking berhasil dibatalkan", result)
}

func (c *WarehouseVisitController) SetAttendance(ctx *fiber.Ctx) error {
	bookingID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID booking tidak valid", nil)
	}
	var req models.WarehouseVisitAttendanceRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}
	previous, result, err := c.service.SetBookingAttendance(ctx.UserContext(), bookingID, req.Status)
	if err != nil {
		return warehouseVisitErrorResponse(ctx, err)
	}
	c.activityLog.LogUpdate(ctx, "warehouse_visit", "warehouse_visit_booking", result.ID, "Kehadiran booking kunjungan gudang dicatat", previous, result)
	return utils.SuccessResponse(ctx, "Kehadiran booking berhasil dicatat", result)
}

func warehouseVisitAdminID(ctx *fiber.Ctx) (uuid.UUID, error) {
	value, ok := ctx.Locals("admin_id").(string)
	if !ok {
		return uuid.Nil, errors.New("admin_id tidak tersedia")
	}
	return uuid.Parse(value)
}

func warehouseVisitErrorResponse(ctx *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, services.ErrWarehouseVisitValidation):
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, repositories.ErrWarehouseVisitNotFound):
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	case errors.Is(err, repositories.ErrWarehouseVisitHasBookings):
		return utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
	case errors.Is(err, repositories.ErrWarehouseVisitConflict):
		return utils.ErrorResponse(ctx, http.StatusConflict, "Perubahan ditolak karena kapasitas tidak cukup atau status booking sudah berubah", nil)
	default:
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal memproses booking kunjungan gudang", nil)
	}
}
