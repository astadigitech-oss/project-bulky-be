package controllers

import (
	"net/http"
	"strconv"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type FAQController struct {
	service services.FAQService
}

func NewFAQController(service services.FAQService) *FAQController {
	return &FAQController{service: service}
}

// Get - Get FAQ with all items for admin panel
// GET /api/panel/faq
func (c *FAQController) Get(ctx *gin.Context) {
	result, err := c.service.Get(ctx.Request.Context())
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data FAQ berhasil diambil", result)
}

// Update - Update FAQ (bulk: title + all items)
// PUT /api/panel/faq
func (c *FAQController) Update(ctx *gin.Context) {
	var req models.FAQUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "FAQ berhasil diupdate", result)
}

// AddItem - Add new FAQ item (append to end)
// POST /api/panel/faq/items
func (c *FAQController) AddItem(ctx *gin.Context) {
	var req models.FAQItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.AddItem(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "FAQ item berhasil ditambahkan",
		"data":    result,
	})
}

// UpdateItem - Update FAQ item by index
// PUT /api/panel/faq/items/:index
func (c *FAQController) UpdateItem(ctx *gin.Context) {
	indexStr := ctx.Param("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Index harus berupa angka", nil)
		return
	}

	var req models.FAQItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.UpdateItem(ctx.Request.Context(), index, &req)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		// Check if error message contains "tidak ditemukan"
		if len(err.Error()) > 0 && err.Error()[:3] == "FAQ" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "FAQ item berhasil diupdate", result)
}

// DeleteItem - Delete FAQ item by index
// DELETE /api/panel/faq/items/:index
func (c *FAQController) DeleteItem(ctx *gin.Context) {
	indexStr := ctx.Param("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Index harus berupa angka", nil)
		return
	}

	err = c.service.DeleteItem(ctx.Request.Context(), index)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		// Check if error message contains "tidak ditemukan"
		if len(err.Error()) > 0 && err.Error()[:3] == "FAQ" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "FAQ item berhasil dihapus", nil)
}

// ReorderItem - Reorder FAQ items (up/down)
// PATCH /api/panel/faq/items/reorder
func (c *FAQController) ReorderItem(ctx *gin.Context) {
	var req models.FAQReorderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.ReorderItem(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		// Check for edge case errors
		if err.Error() == "item sudah berada di posisi paling atas" ||
			err.Error() == "item sudah berada di posisi paling bawah" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		// Check if error message contains "tidak ditemukan"
		if len(err.Error()) > 0 && err.Error()[:3] == "FAQ" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "FAQ item berhasil direorder", result)
}
