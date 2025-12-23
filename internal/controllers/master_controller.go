package controllers

import (
	"net/http"

	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MasterController struct {
	service services.MasterService
}

func NewMasterController(service services.MasterService) *MasterController {
	return &MasterController{service: service}
}

func (c *MasterController) GetDropdown(ctx *gin.Context) {
	result, err := c.service.GetDropdown(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data dropdown berhasil diambil", result)
}
