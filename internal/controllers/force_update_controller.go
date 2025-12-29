package controllers

import (
	"net/http"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ForceUpdateController struct {
	service services.ForceUpdateService
}

func NewForceUpdateController(service services.ForceUpdateService) *ForceUpdateController {
	return &ForceUpdateController{service: service}
}

// CreateForceUpdate creates force update configuration
func (c *ForceUpdateController) CreateForceUpdate(ctx *gin.Context) {
	var req models.CreateForceUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	forceUpdate, err := c.service.CreateForceUpdate(&req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to create force update", err.Error())
		return
	}

	response := models.ForceUpdateDetailResponse{
		ID:              uint(forceUpdate.ID.ID()),
		KodeVersi:       forceUpdate.KodeVersi,
		UpdateType:      string(forceUpdate.UpdateType),
		InformasiUpdate: forceUpdate.InformasiUpdate,
		IsActive:        forceUpdate.IsActive,
		CreatedAt:       forceUpdate.CreatedAt,
		UpdatedAt:       forceUpdate.UpdatedAt,
	}

	utils.CreatedResponse(ctx, "Force update created successfully", response)
}

// UpdateForceUpdate updates force update configuration
func (c *ForceUpdateController) UpdateForceUpdate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid force update ID", err.Error())
		return
	}

	var req models.UpdateForceUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	forceUpdate, err := c.service.UpdateForceUpdate(uint(id), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "force update not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Failed to update force update", err.Error())
		return
	}

	response := models.ForceUpdateDetailResponse{
		ID:              uint(forceUpdate.ID.ID()),
		KodeVersi:       forceUpdate.KodeVersi,
		UpdateType:      string(forceUpdate.UpdateType),
		InformasiUpdate: forceUpdate.InformasiUpdate,
		IsActive:        forceUpdate.IsActive,
		CreatedAt:       forceUpdate.CreatedAt,
		UpdatedAt:       forceUpdate.UpdatedAt,
	}

	utils.SuccessResponse(ctx, "Force update updated successfully", response)
}

// DeleteForceUpdate deletes force update configuration
func (c *ForceUpdateController) DeleteForceUpdate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid force update ID", err.Error())
		return
	}

	err = c.service.DeleteForceUpdate(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "force update not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Failed to delete force update", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Force update deleted successfully", nil)
}

// GetForceUpdateByID gets force update by ID
func (c *ForceUpdateController) GetForceUpdateByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid force update ID", err.Error())
		return
	}

	forceUpdate, err := c.service.GetForceUpdateByID(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "force update not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Failed to get force update", err.Error())
		return
	}

	response := models.ForceUpdateDetailResponse{
		ID:              uint(forceUpdate.ID.ID()),
		KodeVersi:       forceUpdate.KodeVersi,
		UpdateType:      string(forceUpdate.UpdateType),
		InformasiUpdate: forceUpdate.InformasiUpdate,
		IsActive:        forceUpdate.IsActive,
		CreatedAt:       forceUpdate.CreatedAt,
		UpdatedAt:       forceUpdate.UpdatedAt,
	}

	utils.SuccessResponse(ctx, "Force update retrieved successfully", response)
}

// GetAllForceUpdates gets all force update configurations
func (c *ForceUpdateController) GetAllForceUpdates(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	forceUpdates, total, err := c.service.GetAllForceUpdates(page, limit)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get force updates", err.Error())
		return
	}

	var response []models.ForceUpdateListResponse
	for _, fu := range forceUpdates {
		response = append(response, models.ForceUpdateListResponse{
			ID:         uint(fu.ID.ID()),
			KodeVersi:  fu.KodeVersi,
			UpdateType: string(fu.UpdateType),
			IsActive:   fu.IsActive,
			CreatedAt:  fu.CreatedAt,
		})
	}

	meta := models.PaginationMeta{
		Halaman:      page,
		PerHalaman:   limit,
		TotalData:    total,
		TotalHalaman: int64((int(total) + limit - 1) / limit),
	}

	utils.PaginatedSuccessResponse(ctx, "Force updates retrieved successfully", response, meta)
}

// SetActiveForceUpdate sets active force update configuration
func (c *ForceUpdateController) SetActiveForceUpdate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid force update ID", err.Error())
		return
	}

	err = c.service.SetActiveForceUpdate(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "force update not found" {
			status = http.StatusNotFound
		}
		utils.SimpleErrorResponse(ctx, status, "Failed to set active force update", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Force update activated successfully", nil)
}
