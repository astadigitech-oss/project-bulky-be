package controllers

import (
	"net/http"
	"strconv"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type ForceUpdateController struct {
	service services.ForceUpdateService
}

func NewForceUpdateController(service services.ForceUpdateService) *ForceUpdateController {
	return &ForceUpdateController{service: service}
}

// CreateForceUpdate creates force update configuration
func (c *ForceUpdateController) CreateForceUpdate(ctx *fiber.Ctx) error {
	var req models.CreateForceUpdateRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Data request tidak valid", utils.GetValidationErrorMessage(err))
	}

	forceUpdate, err := c.service.CreateForceUpdate(&req)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat konfigurasi force update", err.Error())
	}

	response := models.ForceUpdateDetailResponse{
		ID:              forceUpdate.ID.String(),
		KodeVersi:       forceUpdate.KodeVersi,
		UpdateType:      string(forceUpdate.UpdateType),
		InformasiUpdate: forceUpdate.InformasiUpdate,
		IsActive:        forceUpdate.IsActive,
		CreatedAt:       forceUpdate.CreatedAt,
		UpdatedAt:       forceUpdate.UpdatedAt,
	}

	return utils.CreatedResponse(ctx, "Konfigurasi force update berhasil dibuat", response)
}

// UpdateForceUpdate updates force update configuration
func (c *ForceUpdateController) UpdateForceUpdate(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateForceUpdateRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Data request tidak valid", utils.GetValidationErrorMessage(err))
	}

	forceUpdate, err := c.service.UpdateForceUpdate(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "force update not found" {
			status = http.StatusNotFound
		}
		return utils.SimpleErrorResponse(ctx, status, "Gagal memperbarui konfigurasi force update", err.Error())
	}

	response := models.ForceUpdateDetailResponse{
		ID:              forceUpdate.ID.String(),
		KodeVersi:       forceUpdate.KodeVersi,
		UpdateType:      string(forceUpdate.UpdateType),
		InformasiUpdate: forceUpdate.InformasiUpdate,
		IsActive:        forceUpdate.IsActive,
		CreatedAt:       forceUpdate.CreatedAt,
		UpdatedAt:       forceUpdate.UpdatedAt,
	}

	return utils.SuccessResponse(ctx, "Konfigurasi force update berhasil diperbarui", response)
}

// DeleteForceUpdate deletes force update configuration
func (c *ForceUpdateController) DeleteForceUpdate(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteForceUpdate(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "force update not found" {
			status = http.StatusNotFound
		}
		return utils.SimpleErrorResponse(ctx, status, "Gagal menghapus konfigurasi force update", err.Error())
	}

	return utils.SuccessResponse(ctx, "Konfigurasi force update berhasil dihapus", nil)
}

// GetForceUpdateByID gets force update by ID
func (c *ForceUpdateController) GetForceUpdateByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	forceUpdate, err := c.service.GetForceUpdateByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "force update not found" {
			status = http.StatusNotFound
		}
		return utils.SimpleErrorResponse(ctx, status, "Gagal mengambil data force update", err.Error())
	}

	response := models.ForceUpdateDetailResponse{
		ID:              forceUpdate.ID.String(),
		KodeVersi:       forceUpdate.KodeVersi,
		UpdateType:      string(forceUpdate.UpdateType),
		InformasiUpdate: forceUpdate.InformasiUpdate,
		IsActive:        forceUpdate.IsActive,
		CreatedAt:       forceUpdate.CreatedAt,
		UpdatedAt:       forceUpdate.UpdatedAt,
	}

	return utils.SuccessResponse(ctx, "Data force update berhasil diambil", response)
}

// GetAllForceUpdates gets all force update configurations
func (c *ForceUpdateController) GetAllForceUpdates(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	forceUpdates, total, err := c.service.GetAllForceUpdates(page, limit)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data force update", err.Error())
	}

	var response []models.ForceUpdateListResponse
	for _, fu := range forceUpdates {
		response = append(response, models.ForceUpdateListResponse{
			ID:         fu.ID.String(),
			KodeVersi:  fu.KodeVersi,
			UpdateType: string(fu.UpdateType),
			IsActive:   fu.IsActive,
			CreatedAt:  fu.CreatedAt,
		})
	}

	meta := models.NewPaginationMeta(page, limit, total)

	return utils.PaginatedSuccessResponse(ctx, "Data force update berhasil diambil", response, meta)
}

// SetActiveForceUpdate sets active force update configuration
func (c *ForceUpdateController) SetActiveForceUpdate(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.SetActiveForceUpdate(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "force update not found" {
			status = http.StatusNotFound
		}
		return utils.SimpleErrorResponse(ctx, status, "Gagal mengaktifkan force update", err.Error())
	}

	return utils.SuccessResponse(ctx, "Force update berhasil diaktifkan", nil)
}
