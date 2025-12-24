package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ProvinsiController struct {
	service services.ProvinsiService
}

func NewProvinsiController(service services.ProvinsiService) *ProvinsiController {
	return &ProvinsiController{service: service}
}

func (c *ProvinsiController) Create(ctx *gin.Context) {
	var req models.CreateProvinsiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Provinsi berhasil dibuat", result)
}

func (c *ProvinsiController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Detail provinsi berhasil diambil", result)
}

func (c *ProvinsiController) FindAll(ctx *gin.Context) {
	var params models.PaginationRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedResponse(ctx, http.StatusOK, "Data provinsi berhasil diambil", items, meta)
}

func (c *ProvinsiController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.UpdateProvinsiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Provinsi berhasil diupdate", result)
}

func (c *ProvinsiController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Provinsi berhasil dihapus", nil)
}
