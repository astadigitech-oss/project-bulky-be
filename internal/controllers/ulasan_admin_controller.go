package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UlasanAdminController struct {
	ulasanService services.UlasanAdminService
	activityLog   services.ActivityLogService
}

func NewUlasanAdminController(ulasanService services.UlasanAdminService, activityLog services.ActivityLogService) *UlasanAdminController {
	return &UlasanAdminController{
		ulasanService: ulasanService,
		activityLog:   activityLog,
	}
}

// GetAll retrieves all ulasan with pagination and filters (admin)
func (c *UlasanAdminController) GetAll(ctx *fiber.Ctx) error {
	var params dto.UlasanAdminQueryParams

	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
	}

	params.SetDefaults()

	ulasan, meta, err := c.ulasanService.GetAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data ulasan", err.Error())
	}

	return utils.PaginatedSuccessResponse(ctx, "Data ulasan berhasil diambil", ulasan, *meta)
}

// GetByID retrieves an ulasan by ID (admin)
func (c *UlasanAdminController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	ulasan, err := c.ulasanService.GetByID(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "ulasan tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Ulasan tidak ditemukan", "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail ulasan", err.Error())
	}

	return utils.SuccessResponse(ctx, "Detail ulasan berhasil diambil", ulasan)
}

// Approve approves an ulasan
func (c *UlasanAdminController) Approve(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	// Get admin ID from context
	adminIDStr := localsString(ctx, "admin_id")
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak valid", err.Error())
	}

	if err := c.ulasanService.Approve(ctx.UserContext(), id, adminID); err != nil {
		if err.Error() == "ulasan tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Ulasan tidak ditemukan", "")
		}
		if err.Error() == "ulasan sudah di-approve" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Ulasan sudah di-approve", "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal approve ulasan", err.Error())
	}

	// Return simple response with updated data
	ulasan, err := c.ulasanService.GetByID(ctx.UserContext(), id)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data ulasan", err.Error())
	}

	response := dto.UlasanApproveResponse{
		ID:         ulasan.ID,
		IsApproved: ulasan.IsApproved,
		ApprovedAt: ulasan.ApprovedAt,
	}
	if ulasan.ApprovedBy != nil {
		response.ApprovedBy = &ulasan.ApprovedBy.ID
	}

	c.activityLog.Log(ctx, models.ActionApprove, "ulasan", "Ulasan berhasil di-approve")
	return utils.SuccessResponse(ctx, "Ulasan berhasil di-approve", response)
}

// Reject rejects an ulasan
func (c *UlasanAdminController) Reject(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	if err := c.ulasanService.Reject(ctx.UserContext(), id); err != nil {
		if err.Error() == "ulasan tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Ulasan tidak ditemukan", "")
		}
		if err.Error() == "ulasan sudah di-reject" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Ulasan sudah di-reject", "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal reject ulasan", err.Error())
	}

	// Return simple response with updated data
	ulasan, err := c.ulasanService.GetByID(ctx.UserContext(), id)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data ulasan", err.Error())
	}

	response := dto.UlasanApproveResponse{
		ID:         ulasan.ID,
		IsApproved: ulasan.IsApproved,
		ApprovedAt: ulasan.ApprovedAt,
		ApprovedBy: nil,
	}

	c.activityLog.Log(ctx, models.ActionReject, "ulasan", "Ulasan berhasil di-reject")
	return utils.SuccessResponse(ctx, "Ulasan berhasil di-reject", response)
}

// BulkApprove approves multiple ulasan
func (c *UlasanAdminController) BulkApprove(ctx *fiber.Ctx) error {
	var req dto.BulkApproveUlasanRequest

	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	// Get admin ID from context
	adminIDStr := localsString(ctx, "admin_id")
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak valid", err.Error())
	}

	result, err := c.ulasanService.BulkApprove(ctx.UserContext(), req.IDs, adminID)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal bulk approve ulasan", err.Error())
	}

	message := fmt.Sprintf("%d ulasan berhasil di-approve", result.ApprovedCount)
	c.activityLog.Log(ctx, models.ActionApprove, "ulasan", message)
	return utils.SuccessResponse(ctx, message, result)
}

// Delete deletes an ulasan (soft delete)
func (c *UlasanAdminController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	if err := c.ulasanService.Delete(ctx.UserContext(), id); err != nil {
		if err.Error() == "ulasan tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Ulasan tidak ditemukan", "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus ulasan", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionDelete, "ulasan", "Ulasan berhasil dihapus")
	return utils.SuccessResponse(ctx, "Ulasan berhasil dihapus", nil)
}
