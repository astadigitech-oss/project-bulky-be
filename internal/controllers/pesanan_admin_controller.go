package controllers

import (
	"errors"
	"net/http"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PesananAdminController struct {
	pesananService services.PesananAdminService
	activityLog    services.ActivityLogService
}

func NewPesananAdminController(pesananService services.PesananAdminService, activityLog services.ActivityLogService) *PesananAdminController {
	return &PesananAdminController{
		pesananService: pesananService,
		activityLog:    activityLog,
	}
}

// GetAll retrieves all pesanan with pagination and filters (admin)
func (c *PesananAdminController) GetAll(ctx *fiber.Ctx) error {
	var params dto.PesananAdminQueryParams

	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
	}

	params.SetDefaults()

	pesanan, meta, err := c.pesananService.GetAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data pesanan", err.Error())
	}

	return utils.PaginatedSuccessResponse(ctx, "Data pesanan berhasil diambil", pesanan, *meta)
}

// GetByID retrieves a pesanan by ID (admin)
func (c *PesananAdminController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	pesanan, err := c.pesananService.GetByID(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "pesanan tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Pesanan tidak ditemukan", "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail pesanan", err.Error())
	}

	return utils.SuccessResponse(ctx, "Detail pesanan berhasil diambil", pesanan)
}

// UpdateStatus updates pesanan status
func (c *PesananAdminController) UpdateStatus(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req dto.UpdatePesananStatusRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	// Get admin ID from context
	adminIDStr := localsString(ctx, "admin_id")
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak valid", err.Error())
	}

	result, err := c.pesananService.UpdateStatus(ctx.UserContext(), id, &req, adminID)
	if err != nil {
		if err.Error() == "pesanan tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Pesanan tidak ditemukan", "")
		}
		// Check if it's a validation error for status transition
		if err.Error() == "tidak dapat mengubah status dari COMPLETED" ||
			err.Error() == "tidak dapat mengubah status dari CANCELLED" ||
			len(err.Error()) > 25 && err.Error()[:25] == "tidak dapat mengubah status" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Transisi status tidak valid", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal update status pesanan", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "pesanan", "Status pesanan berhasil diupdate")
	return utils.SuccessResponse(ctx, "Status pesanan berhasil diupdate", result)
}

// Delete deletes a pesanan (soft delete)
func (c *PesananAdminController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	if err := c.pesananService.Delete(ctx.UserContext(), id); err != nil {
		if err.Error() == "pesanan tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Pesanan tidak ditemukan", "")
		}
		if err.Error() == "hanya pesanan dengan status CANCELLED yang dapat dihapus" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Hanya pesanan dengan status CANCELLED yang dapat dihapus", "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus pesanan", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionDelete, "pesanan", "Pesanan berhasil dihapus")
	return utils.SuccessResponse(ctx, "Pesanan berhasil dihapus", nil)
}

// RetryBooking retries failed shipping booking for a pesanan
func (c *PesananAdminController) RetryBooking(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	result, err := c.pesananService.RetryBooking(ctx.UserContext(), id)
	if err != nil {
		msg := err.Error()
		switch {
		case len(msg) > 17 && msg[:17] == "retry:bad_request":
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest,
				"Retry tidak diperlukan. Booking sudah berhasil atau pesanan bukan tipe delivery.", "")
		case len(msg) > 19 && msg[:19] == "retry:already_booked":
			bookingRef := msg[20:]
			return ctx.Status(http.StatusOK).JSON(fiber.Map{
				"success": false,
				"message": "Booking sudah ada, tidak perlu retry",
				"data":    fiber.Map{"booking_id": bookingRef},
			})
		case len(msg) > 20 && msg[:20] == "retry:city_not_mapped":
			// extract kota name from error message
			return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": msg[21:],
				"data":    fiber.Map{},
			})
		case len(msg) > 20 && msg[:20] == "retry:provider_error":
			parts := strings.SplitN(msg[21:], ":", 2)
			provider := ""
			errDetail := msg[21:]
			if len(parts) == 2 {
				provider = parts[0]
				errDetail = parts[1]
			}
			return ctx.Status(http.StatusBadGateway).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menghubungi layanan pengiriman. Silakan coba beberapa saat lagi.",
				"data":    fiber.Map{"provider": provider, "error": errDetail},
			})
		case msg == "pesanan tidak ditemukan":
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Pesanan tidak ditemukan", "")
		default:
			return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal melakukan retry booking", msg)
		}
	}

	return utils.SuccessResponse(ctx, "Booking pengiriman berhasil dibuat", result)
}

// TrackDelivery retrieves live tracking info from shipping provider
func (c *PesananAdminController) TrackDelivery(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	result, err := c.pesananService.TrackDelivery(ctx.UserContext(), id)
	if err != nil {
		msg := err.Error()
		switch {
		case msg == "pesanan tidak ditemukan":
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Pesanan tidak ditemukan", "")
		case len(msg) > 20 && msg[:20] == "tracking:not_applicable":
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, msg[23:], "")
		case strings.Contains(msg, "belum memiliki"):
			return utils.SimpleErrorResponse(ctx, http.StatusConflict, msg, "")
		default:
			return utils.SimpleErrorResponse(ctx, http.StatusBadGateway, "Gagal mengambil data tracking", msg)
		}
	}

	return utils.SuccessResponse(ctx, "Data tracking berhasil diambil", result)
}

// GetStatistics retrieves pesanan statistics
func (c *PesananAdminController) GetStatistics(ctx *fiber.Ctx) error {
	var params dto.StatisticsQueryParams
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
	}

	stats, err := c.pesananService.GetStatistics(ctx.UserContext(), &params)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil statistik pesanan", err.Error())
	}

	return utils.SuccessResponse(ctx, "Statistik pesanan berhasil diambil", stats)
}
