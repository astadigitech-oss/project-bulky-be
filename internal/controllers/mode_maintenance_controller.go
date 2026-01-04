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
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	maintenance, err := c.service.CreateMaintenance(&req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to create maintenance", err.Error())
		return
	}

	response := models.MaintenanceDetailResponse{
		ID:              uint(maintenance.ID.ID()),
		Judul:           maintenance.Judul,
		TipeMaintenance: string(maintenance.TipeMaintenance),
		Deskripsi:       maintenance.Deskripsi,
		IsActive:        maintenance.IsActive,
		CreatedAt:       maintenance.CreatedAt,
		UpdatedAt:       maintenance.UpdatedAt,
	}

	utils.CreatedResponse(ctx, "Maintenance created successfully", response)
}

// UpdateMaintenance updates maintenance mode configuration
func (c *ModeMaintenanceController) UpdateMaintenance(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid maintenance ID", err.Error())
		return
	}

	var req models.UpdateMaintenanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	maintenance, err := c.service.UpdateMaintenance(uint(id), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Failed to update maintenance", err.Error())
		return
	}

	response := models.MaintenanceDetailResponse{
		ID:              uint(maintenance.ID.ID()),
		Judul:           maintenance.Judul,
		TipeMaintenance: string(maintenance.TipeMaintenance),
		Deskripsi:       maintenance.Deskripsi,
		IsActive:        maintenance.IsActive,
		CreatedAt:       maintenance.CreatedAt,
		UpdatedAt:       maintenance.UpdatedAt,
	}

	utils.SuccessResponse(ctx, "Maintenance updated successfully", response)
}

// DeleteMaintenance deletes maintenance mode configuration
func (c *ModeMaintenanceController) DeleteMaintenance(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid maintenance ID", err.Error())
		return
	}

	err = c.service.DeleteMaintenance(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Failed to delete maintenance", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Maintenance deleted successfully", nil)
}

// GetMaintenanceByID gets maintenance by ID
func (c *ModeMaintenanceController) GetMaintenanceByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid maintenance ID", err.Error())
		return
	}

	maintenance, err := c.service.GetMaintenanceByID(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Failed to get maintenance", err.Error())
		return
	}

	response := models.MaintenanceDetailResponse{
		ID:              uint(maintenance.ID.ID()),
		Judul:           maintenance.Judul,
		TipeMaintenance: string(maintenance.TipeMaintenance),
		Deskripsi:       maintenance.Deskripsi,
		IsActive:        maintenance.IsActive,
		CreatedAt:       maintenance.CreatedAt,
		UpdatedAt:       maintenance.UpdatedAt,
	}

	utils.SuccessResponse(ctx, "Maintenance retrieved successfully", response)
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
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get maintenances", err.Error())
		return
	}

	var response []models.MaintenanceListResponse
	for _, m := range maintenances {
		response = append(response, models.MaintenanceListResponse{
			ID:              uint(m.ID.ID()),
			Judul:           m.Judul,
			TipeMaintenance: string(m.TipeMaintenance),
			IsActive:        m.IsActive,
			CreatedAt:       m.CreatedAt,
		})
	}

	meta := models.NewPaginationMeta(page, limit, total)

	utils.PaginatedSuccessResponse(ctx, "Maintenances retrieved successfully", response, meta)
}

// ActivateMaintenance activates maintenance mode
func (c *ModeMaintenanceController) ActivateMaintenance(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid maintenance ID", err.Error())
		return
	}

	err = c.service.ActivateMaintenance(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Failed to activate maintenance", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Maintenance activated successfully", nil)
}

// DeactivateMaintenance deactivates maintenance mode
func (c *ModeMaintenanceController) DeactivateMaintenance(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid maintenance ID", err.Error())
		return
	}

	err = c.service.DeactivateMaintenance(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Failed to deactivate maintenance", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Maintenance deactivated successfully", nil)
}
