package controllers

import (
	"errors"
	"net/http"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
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
func (c *PesananAdminController) GetAll(ctx *gin.Context) {
	var params dto.PesananAdminQueryParams

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

	pesanan, meta, err := c.pesananService.GetAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data pesanan", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data pesanan berhasil diambil", pesanan, *meta)
}

// GetByID retrieves a pesanan by ID (admin)
func (c *PesananAdminController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	pesanan, err := c.pesananService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "pesanan tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Pesanan tidak ditemukan", "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail pesanan", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Detail pesanan berhasil diambil", pesanan)
}

// UpdateStatus updates pesanan status
func (c *PesananAdminController) UpdateStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	var req dto.UpdatePesananStatusRequest
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

	result, err := c.pesananService.UpdateStatus(ctx.Request.Context(), id, &req, adminID)
	if err != nil {
		if err.Error() == "pesanan tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Pesanan tidak ditemukan", "")
			return
		}
		// Check if it's a validation error for status transition
		if err.Error() == "tidak dapat mengubah status dari COMPLETED" ||
			err.Error() == "tidak dapat mengubah status dari CANCELLED" ||
			len(err.Error()) > 25 && err.Error()[:25] == "tidak dapat mengubah status" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Transisi status tidak valid", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal update status pesanan", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Status pesanan berhasil diupdate", result)
}

// Delete deletes a pesanan (soft delete)
func (c *PesananAdminController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	if err := c.pesananService.Delete(ctx.Request.Context(), id); err != nil {
		if err.Error() == "pesanan tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Pesanan tidak ditemukan", "")
			return
		}
		if err.Error() == "hanya pesanan dengan status CANCELLED yang dapat dihapus" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Hanya pesanan dengan status CANCELLED yang dapat dihapus", "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus pesanan", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Pesanan berhasil dihapus", nil)
}

// GetStatistics retrieves pesanan statistics
func (c *PesananAdminController) GetStatistics(ctx *gin.Context) {
	tanggalDari := ctx.Query("tanggal_dari")
	tanggalSampai := ctx.Query("tanggal_sampai")

	var dari, sampai *string
	if tanggalDari != "" {
		dari = &tanggalDari
	}
	if tanggalSampai != "" {
		sampai = &tanggalSampai
	}

	stats, err := c.pesananService.GetStatistics(ctx.Request.Context(), dari, sampai)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil statistik pesanan", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Statistik pesanan berhasil diambil", stats)
}
