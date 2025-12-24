package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type KecamatanController struct {
	service services.KecamatanService
}

func NewKecamatanController(service services.KecamatanService) *KecamatanController {
	return &KecamatanController{service: service}
}

func (c *KecamatanController) Create(ctx *gin.Context) {
	var req models.CreateKecamatanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Kecamatan berhasil dibuat", result)
}

func (c *KecamatanController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Detail kecamatan berhasil diambil", result)
}

func (c *KecamatanController) FindAll(ctx *gin.Context) {
	var params models.KecamatanFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedResponse(ctx, http.StatusOK, "Data kecamatan berhasil diambil", items, meta)
}

func (c *KecamatanController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.UpdateKecamatanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Kecamatan berhasil diupdate", result)
}

func (c *KecamatanController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Kecamatan berhasil dihapus", nil)
}
