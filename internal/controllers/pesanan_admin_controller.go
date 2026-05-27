package controllers

import (
	"errors"
	"net/http"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PesananAdminController struct {
	pesananService services.PesananAdminService
}

func NewPesananAdminController(pesananService services.PesananAdminService) *PesananAdminController {
	return &PesananAdminController{
		pesananService: pesananService,
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

	return utils.SuccessResponse(ctx, "Pesanan berhasil dihapus", nil)
}

// GetStatistics retrieves pesanan statistics
func (c *PesananAdminController) GetStatistics(ctx *fiber.Ctx) error {
	tanggalDari := ctx.Query("tanggal_dari")
	tanggalSampai := ctx.Query("tanggal_sampai")

	var dari, sampai *string
	if tanggalDari != "" {
		dari = &tanggalDari
	}
	if tanggalSampai != "" {
		sampai = &tanggalSampai
	}

	stats, err := c.pesananService.GetStatistics(ctx.UserContext(), dari, sampai)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil statistik pesanan", err.Error())
	}

	return utils.SuccessResponse(ctx, "Statistik pesanan berhasil diambil", stats)
}
