package controllers

import (
	"net/http"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
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
func (c *AppStatusController) CheckVersion(ctx *fiber.Ctx) error {
	currentVersion := ctx.Query("version")
	if currentVersion == "" {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Version parameter is required", "version query parameter is missing")
	}

	response, err := c.forceUpdateService.CheckVersion(currentVersion)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to check version", err.Error())
	}

	return utils.SuccessResponse(ctx, "Version check completed", response)
}

// CheckMaintenance checks maintenance mode status
func (c *AppStatusController) CheckMaintenance(ctx *fiber.Ctx) error {
	response, err := c.maintenanceService.CheckMaintenance()
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to check maintenance", err.Error())
	}

	return utils.SuccessResponse(ctx, "Maintenance check completed", response)
}

// AppStatus checks app version and maintenance status
func (c *AppStatusController) AppStatus(ctx *fiber.Ctx) error {
	currentVersion := ctx.Query("version")
	if currentVersion == "" {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Version parameter is required", "version query parameter is missing")
	}

	// Check version
	versionResponse, err := c.forceUpdateService.CheckVersion(currentVersion)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to check version", err.Error())
	}

	// Check maintenance
	maintenanceResponse, err := c.maintenanceService.CheckMaintenance()
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to check maintenance", err.Error())
	}

	// Combine responses
	response := map[string]interface{}{
		"version":     versionResponse,
		"maintenance": maintenanceResponse,
	}

	return utils.SuccessResponse(ctx, "App status check completed", response)
}
