package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type KotaController struct {
	service services.KotaService
}

func NewKotaController(service services.KotaService) *KotaController {
	return &KotaController{service: service}
}

func (c *KotaController) Create(ctx *gin.Context) {
	var req models.CreateKotaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Kota berhasil dibuat", result)
}

func (c *KotaController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Detail kota berhasil diambil", result)
}

func (c *KotaController) FindAll(ctx *gin.Context) {
	var params models.KotaFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedResponse(ctx, http.StatusOK, "Data kota berhasil diambil", items, meta)
}

func (c *KotaController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.UpdateKotaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Kota berhasil diupdate", result)
}

func (c *KotaController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Kota berhasil dihapus", nil)
}
