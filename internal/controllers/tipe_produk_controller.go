package controllers

import (
	"net/http"

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

// FindAll retrieves all tipe produk without pagination
func (c *TipeProdukController) FindAll(ctx *gin.Context) {
	items, err := c.service.FindAll(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data tipe produk berhasil diambil", items)
}

// FindAllWithProduk retrieves all tipe produk with their products
func (c *TipeProdukController) FindAllWithProduk(ctx *gin.Context) {
	items, err := c.service.FindAllWithProduk(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data tipe produk dengan produk berhasil diambil", items)
}
