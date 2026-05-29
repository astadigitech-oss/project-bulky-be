package controllers

import (
	"net/http"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ModeMaintenanceController struct {
	service     services.ModeMaintenanceService
	activityLog services.ActivityLogService
}

func NewModeMaintenanceController(service services.ModeMaintenanceService, activityLog services.ActivityLogService) *ModeMaintenanceController {
	return &ModeMaintenanceController{service: service, activityLog: activityLog}
}

// CreateMaintenance creates maintenance mode configuration
func (c *ModeMaintenanceController) CreateMaintenance(ctx *fiber.Ctx) error {
	var req models.CreateMaintenanceRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Data request tidak valid", utils.GetValidationErrorMessage(err))
	}

	maintenance, err := c.service.CreateMaintenance(&req)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat konfigurasi maintenance", err.Error())
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

	c.activityLog.Log(ctx, models.ActionCreate, "mode_maintenance", "Konfigurasi maintenance berhasil dibuat")
	return utils.CreatedResponse(ctx, "Konfigurasi maintenance berhasil dibuat", response)
}

// UpdateMaintenance updates maintenance mode configuration
func (c *ModeMaintenanceController) UpdateMaintenance(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateMaintenanceRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Data request tidak valid", utils.GetValidationErrorMessage(err))
	}

	maintenance, err := c.service.UpdateMaintenance(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		return utils.SimpleErrorResponse(ctx, status, "Gagal memperbarui konfigurasi maintenance", err.Error())
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

	c.activityLog.Log(ctx, models.ActionUpdate, "mode_maintenance", "Konfigurasi maintenance berhasil diperbarui")
	return utils.SuccessResponse(ctx, "Konfigurasi maintenance berhasil diperbarui", response)
}

// DeleteMaintenance deletes maintenance mode configuration
func (c *ModeMaintenanceController) DeleteMaintenance(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteMaintenance(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		return utils.SimpleErrorResponse(ctx, status, "Gagal menghapus konfigurasi maintenance", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionDelete, "mode_maintenance", "Konfigurasi maintenance berhasil dihapus")
	return utils.SuccessResponse(ctx, "Konfigurasi maintenance berhasil dihapus", nil)
}

// GetMaintenanceByID gets maintenance by ID
func (c *ModeMaintenanceController) GetMaintenanceByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	maintenance, err := c.service.GetMaintenanceByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		return utils.SimpleErrorResponse(ctx, status, "Gagal mengambil data maintenance", err.Error())
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

	return utils.SuccessResponse(ctx, "Data maintenance berhasil diambil", response)
}

// GetAllMaintenances gets all maintenance mode configurations
func (c *ModeMaintenanceController) GetAllMaintenances(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	maintenances, total, err := c.service.GetAllMaintenances(page, limit)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data maintenance", err.Error())
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

	return utils.PaginatedSuccessResponse(ctx, "Data maintenance berhasil diambil", response, meta)
}

// ActivateMaintenance activates maintenance mode
func (c *ModeMaintenanceController) ActivateMaintenance(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.ActivateMaintenance(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		return utils.SimpleErrorResponse(ctx, status, "Gagal mengaktifkan maintenance", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "mode_maintenance", "Maintenance berhasil diaktifkan")
	return utils.SuccessResponse(ctx, "Maintenance berhasil diaktifkan", nil)
}

// DeactivateMaintenance deactivates maintenance mode
func (c *ModeMaintenanceController) DeactivateMaintenance(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeactivateMaintenance(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "maintenance not found" {
			status = http.StatusNotFound
		}
		return utils.SimpleErrorResponse(ctx, status, "Gagal menonaktifkan maintenance", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "mode_maintenance", "Maintenance berhasil dinonaktifkan")
	return utils.SuccessResponse(ctx, "Maintenance berhasil dinonaktifkan", nil)
}
