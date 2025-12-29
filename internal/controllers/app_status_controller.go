package controllers

import (
	"net/http"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AppStatusController struct {
	forceUpdateService services.ForceUpdateService
	maintenanceService services.ModeMaintenanceService
}

func NewAppStatusController(
	forceUpdateService services.ForceUpdateService,
	maintenanceService services.ModeMaintenanceService,
) *AppStatusController {
	return &AppStatusController{
		forceUpdateService: forceUpdateService,
		maintenanceService: maintenanceService,
	}
}

// CheckVersion checks app version for updates
func (c *AppStatusController) CheckVersion(ctx *gin.Context) {
	currentVersion := ctx.Query("version")
	if currentVersion == "" {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Version parameter is required", "version query parameter is missing")
		return
	}

	response, err := c.forceUpdateService.CheckVersion(currentVersion)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to check version", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Version check completed", response)
}

// CheckMaintenance checks maintenance mode status
func (c *AppStatusController) CheckMaintenance(ctx *gin.Context) {
	response, err := c.maintenanceService.CheckMaintenance()
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to check maintenance", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Maintenance check completed", response)
}

// AppStatus checks app version and maintenance status
func (c *AppStatusController) AppStatus(ctx *gin.Context) {
	currentVersion := ctx.Query("version")
	if currentVersion == "" {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Version parameter is required", "version query parameter is missing")
		return
	}

	// Check version
	versionResponse, err := c.forceUpdateService.CheckVersion(currentVersion)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to check version", err.Error())
		return
	}

	// Check maintenance
	maintenanceResponse, err := c.maintenanceService.CheckMaintenance()
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to check maintenance", err.Error())
		return
	}

	// Combine responses
	response := map[string]interface{}{
		"version":     versionResponse,
		"maintenance": maintenanceResponse,
	}

	utils.SuccessResponse(ctx, "App status check completed", response)
}
