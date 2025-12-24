package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AlamatBuyerController struct {
	service services.AlamatBuyerService
}

func NewAlamatBuyerController(service services.AlamatBuyerService) *AlamatBuyerController {
	return &AlamatBuyerController{service: service}
}

func (c *AlamatBuyerController) Create(ctx *gin.Context) {
	var req models.CreateAlamatBuyerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Alamat berhasil ditambahkan", result)
}

func (c *AlamatBuyerController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Detail alamat berhasil diambil", result)
}

func (c *AlamatBuyerController) FindAll(ctx *gin.Context) {
	var params models.AlamatBuyerFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedResponse(ctx, http.StatusOK, "Data alamat berhasil diambil", items, meta)
}

func (c *AlamatBuyerController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.UpdateAlamatBuyerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Alamat berhasil diupdate", result)
}

func (c *AlamatBuyerController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Alamat berhasil dihapus", nil)
}

func (c *AlamatBuyerController) SetDefault(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.SetDefault(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Alamat berhasil diset sebagai default", result)
}
