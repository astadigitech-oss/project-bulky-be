package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TipeProdukController struct {
	service services.TipeProdukService
}

func NewTipeProdukController(service services.TipeProdukService) *TipeProdukController {
	return &TipeProdukController{service: service}
}

// FindAll retrieves all tipe produk with pagination
func (c *TipeProdukController) FindAll(ctx *gin.Context) {
	var params models.PaginationRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data tipe produk berhasil diambil", items, *meta)
}

// FindByID retrieves a single tipe produk by ID
func (c *TipeProdukController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail tipe produk berhasil diambil", result)
}

// FindBySlug retrieves a single tipe produk by slug
// func (c *TipeProdukController) FindBySlug(ctx *gin.Context) {
// 	slug := ctx.Param("slug")

// 	result, err := c.service.FindBySlug(ctx.Request.Context(), slug)
// 	if err != nil {
// 		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
// 		return
// 	}

// 	utils.SuccessResponse(ctx, "Detail tipe produk berhasil diambil", result)
// }
