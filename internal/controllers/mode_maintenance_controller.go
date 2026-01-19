package controllers

import (
	"net/http"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ModeMaintenanceController struct {
	service services.ModeMaintenanceService
}

func NewModeMaintenanceController(service services.ModeMaintenanceService) *ModeMaintenanceController {
	return &ModeMaintenanceController{service: service}
}

// CreateMaintenance creates maintenance mode configuration
func (c *ModeMaintenanceController) CreateMaintenance(ctx *gin.Context) {
	var req models.CreateMaintenanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Data request tidak valid", utils.GetValidationErrorMessage(err))
		return
	}

	maintenance, err := c.service.CreateMaintenance(&req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat konfigurasi maintenance", err.Error())
		return
	}

	response := models.MaintenanceDetailResponse{
		ID:              maintenance.ID.String(),
		Judul:           maintenance.Judul,
		TipeMaintenance: string(maintenance.TipeMaintenance),
		Deskripsi:       maintenance.Deskripsi,
		IsActive:        maintenance.IsActive,
		CreatedAt:       maintenance.CreatedAt,
		UpdatedAt:       maintenance.UpdatedAt,
	}

	utils.CreatedResponse(ctx, "Konfigurasi maintenance berhasil dibuat", response)
}

// UpdateMaintenance updates maintenance mode configuration
func (c *ModeMaintenanceController) UpdateMaintenance(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateMaintenanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Data request tidak valid", utils.GetValidationErrorMessage(err))
		return
	}

	maintenance, err := c.service.UpdateMaintenance(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Gagal memperbarui konfigurasi maintenance", err.Error())
		return
	}

	response := models.MaintenanceDetailResponse{
		ID:              maintenance.ID.String(),
		Judul:           maintenance.Judul,
		TipeMaintenance: string(maintenance.TipeMaintenance),
		Deskripsi:       maintenance.Deskripsi,
		IsActive:        maintenance.IsActive,
		CreatedAt:       maintenance.CreatedAt,
		UpdatedAt:       maintenance.UpdatedAt,
	}

	utils.SuccessResponse(ctx, "Konfigurasi maintenance berhasil diperbarui", response)
}

// DeleteMaintenance deletes maintenance mode configuration
func (c *ModeMaintenanceController) DeleteMaintenance(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeleteMaintenance(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Gagal menghapus konfigurasi maintenance", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Konfigurasi maintenance berhasil dihapus", nil)
}

// GetMaintenanceByID gets maintenance by ID
func (c *ModeMaintenanceController) GetMaintenanceByID(ctx *gin.Context) {
	id := ctx.Param("id")

	maintenance, err := c.service.GetMaintenanceByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Gagal mengambil data maintenance", err.Error())
		return
	}

	response := models.MaintenanceDetailResponse{
		ID:              maintenance.ID.String(),
		Judul:           maintenance.Judul,
		TipeMaintenance: string(maintenance.TipeMaintenance),
		Deskripsi:       maintenance.Deskripsi,
		IsActive:        maintenance.IsActive,
		CreatedAt:       maintenance.CreatedAt,
		UpdatedAt:       maintenance.UpdatedAt,
	}

	utils.SuccessResponse(ctx, "Data maintenance berhasil diambil", response)
}

// GetAllMaintenances gets all maintenance mode configurations
func (c *ModeMaintenanceController) GetAllMaintenances(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	maintenances, total, err := c.service.GetAllMaintenances(page, limit)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data maintenance", err.Error())
		return
	}

	var response []models.MaintenanceListResponse
	for _, m := range maintenances {
		response = append(response, models.MaintenanceListResponse{
			ID:              m.ID.String(),
			Judul:           m.Judul,
			TipeMaintenance: string(m.TipeMaintenance),
			IsActive:        m.IsActive,
			CreatedAt:       m.CreatedAt,
		})
	}

	meta := models.NewPaginationMeta(page, limit, total)

	utils.PaginatedSuccessResponse(ctx, "Data maintenance berhasil diambil", response, meta)
}

// ActivateMaintenance activates maintenance mode
func (c *ModeMaintenanceController) ActivateMaintenance(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.ActivateMaintenance(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Gagal mengaktifkan maintenance", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Maintenance berhasil diaktifkan", nil)
}

// DeactivateMaintenance deactivates maintenance mode
func (c *ModeMaintenanceController) DeactivateMaintenance(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeactivateMaintenance(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Gagal menonaktifkan maintenance", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Maintenance berhasil dinonaktifkan", nil)
}
