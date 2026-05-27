package controllers

import (
	"net/http"

	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type MasterController struct {
	service services.MasterService
}

func NewMasterController(service services.MasterService) *MasterController {
	return &MasterController{service: service}
}

func (c *MasterController) GetDropdown(ctx *fiber.Ctx) error {
	result, err := c.service.GetDropdown(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data dropdown berhasil diambil", result)
}
