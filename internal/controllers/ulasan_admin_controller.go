package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UlasanAdminController struct {
	ulasanService services.UlasanAdminService
}

func NewUlasanAdminController(ulasanService services.UlasanAdminService) *UlasanAdminController {
	return &UlasanAdminController{
		ulasanService: ulasanService,
	}
}

// GetAll retrieves all ulasan with pagination and filters (admin)
func (c *UlasanAdminController) GetAll(ctx *gin.Context) {
	var params dto.UlasanAdminQueryParams

	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
		return
	}

	// Manual validation for required fields
	if params.Page < 1 {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'page' wajib diisi dan minimal 1")
		return
	}
	if params.PerPage < 1 {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'per_page' wajib diisi dan minimal 1")
		return
	}
	if params.PerPage > 100 {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'per_page' maksimal 100")
		return
	}

	params.SetDefaults()

	ulasan, meta, err := c.ulasanService.GetAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data ulasan", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data ulasan berhasil diambil", ulasan, *meta)
}

// GetByID retrieves an ulasan by ID (admin)
func (c *UlasanAdminController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	ulasan, err := c.ulasanService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "ulasan tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Ulasan tidak ditemukan", "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail ulasan", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Detail ulasan berhasil diambil", ulasan)
}

// Approve approves an ulasan
func (c *UlasanAdminController) Approve(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	// Get admin ID from context
	adminIDStr := ctx.GetString("admin_id")
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak valid", err.Error())
		return
	}

	if err := c.ulasanService.Approve(ctx.Request.Context(), id, adminID); err != nil {
		if err.Error() == "ulasan tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Ulasan tidak ditemukan", "")
			return
		}
		if err.Error() == "ulasan sudah di-approve" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Ulasan sudah di-approve", "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal approve ulasan", err.Error())
		return
	}

	// Return simple response with updated data
	ulasan, err := c.ulasanService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data ulasan", err.Error())
		return
	}

	response := dto.UlasanApproveResponse{
		ID:         ulasan.ID,
		IsApproved: ulasan.IsApproved,
		ApprovedAt: ulasan.ApprovedAt,
	}
	if ulasan.ApprovedBy != nil {
		response.ApprovedBy = &ulasan.ApprovedBy.ID
	}

	utils.SuccessResponse(ctx, "Ulasan berhasil di-approve", response)
}

// Reject rejects an ulasan
func (c *UlasanAdminController) Reject(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	if err := c.ulasanService.Reject(ctx.Request.Context(), id); err != nil {
		if err.Error() == "ulasan tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Ulasan tidak ditemukan", "")
			return
		}
		if err.Error() == "ulasan sudah di-reject" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Ulasan sudah di-reject", "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal reject ulasan", err.Error())
		return
	}

	// Return simple response with updated data
	ulasan, err := c.ulasanService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data ulasan", err.Error())
		return
	}

	response := dto.UlasanApproveResponse{
		ID:         ulasan.ID,
		IsApproved: ulasan.IsApproved,
		ApprovedAt: ulasan.ApprovedAt,
		ApprovedBy: nil,
	}

	utils.SuccessResponse(ctx, "Ulasan berhasil di-reject", response)
}

// BulkApprove approves multiple ulasan
func (c *UlasanAdminController) BulkApprove(ctx *gin.Context) {
	var req dto.BulkApproveUlasanRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	// Get admin ID from context
	adminIDStr := ctx.GetString("admin_id")
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak valid", err.Error())
		return
	}

	result, err := c.ulasanService.BulkApprove(ctx.Request.Context(), req.IDs, adminID)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal bulk approve ulasan", err.Error())
		return
	}

	message := fmt.Sprintf("%d ulasan berhasil di-approve", result.ApprovedCount)
	utils.SuccessResponse(ctx, message, result)
}

// Delete deletes an ulasan (soft delete)
func (c *UlasanAdminController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	if err := c.ulasanService.Delete(ctx.Request.Context(), id); err != nil {
		if err.Error() == "ulasan tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Ulasan tidak ditemukan", "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus ulasan", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Ulasan berhasil dihapus", nil)
}
