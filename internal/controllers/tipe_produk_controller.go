package controllers

import (
	"net/http"

	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type TipeProdukController struct {
	service services.TipeProdukService
}

func NewTipeProdukController(service services.TipeProdukService) *TipeProdukController {
	return &TipeProdukController{service: service}
}

// FindAll retrieves all tipe produk without pagination
func (c *TipeProdukController) FindAll(ctx *fiber.Ctx) error {
	items, err := c.service.FindAll(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data tipe produk berhasil diambil", items)
}

// FindAllWithProduk retrieves all tipe produk with their products
func (c *TipeProdukController) FindAllWithProduk(ctx *fiber.Ctx) error {
	items, err := c.service.FindAllWithProduk(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data tipe produk dengan produk berhasil diambil", items)
}

func (c *TipeProdukController) Dropdown(ctx *fiber.Ctx) error {
	// Get all tipe produk for dropdown
	items, err := c.service.FindAll(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data tipe produk", nil)
	}

	// Convert to simple dropdown response
	response := make([]map[string]interface{}, len(items))
	for i, t := range items {
		response[i] = map[string]interface{}{
			"id":   t.ID,
			"nama": t.Nama,
		}
	}

	return utils.SuccessResponse(ctx, "Data dropdown tipe produk berhasil diambil", response)
}
