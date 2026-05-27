package controllers

import (
	"net/http"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FAQController struct {
	service services.FAQService
}

func NewFAQController(service services.FAQService) *FAQController {
	return &FAQController{service: service}
}

// GetAll - Get all FAQ items for admin panel with search and pagination
// GET /api/panel/faq
func (c *FAQController) GetAll(ctx *fiber.Ctx) error {
	var params models.FAQFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	result, meta, err := c.service.GetAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data FAQ berhasil diambil", result, *meta)
}

// GetByID - Get FAQ by ID
// GET /api/panel/faq/:id
func (c *FAQController) GetByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
	}

	result, err := c.service.GetByID(ctx.UserContext(), id)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data FAQ berhasil diambil", result)
}

// GetPublic - Get active FAQ for public
// GET /api/public/faq
func (c *FAQController) GetPublic(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "id")
	if lang != "id" && lang != "en" {
		lang = "id"
	}

	result, err := c.service.GetPublic(ctx.UserContext(), lang)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	message := "Data FAQ berhasil diambil"
	if lang == "en" {
		message = "FAQ data retrieved successfully"
	}

	return utils.SuccessResponse(ctx, message, result)
}

// Create - Create new FAQ item
// POST /api/panel/faq
func (c *FAQController) Create(ctx *fiber.Ctx) error {
	var req models.FAQCreateRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "FAQ berhasil dibuat",
		"data":    result,
	})
}

// Update - Update FAQ item
// PUT /api/panel/faq/:id
func (c *FAQController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
	}

	var req models.FAQUpdateRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "FAQ berhasil diupdate", result)
}

// Delete - Delete FAQ item
// DELETE /api/panel/faq/:id
func (c *FAQController) Delete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
	}

	err = c.service.Delete(ctx.UserContext(), id)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "FAQ berhasil dihapus", nil)
}

// ToggleStatus - Toggle FAQ is_active status
// PATCH /api/panel/faq/:id/toggle-status
func (c *FAQController) ToggleStatus(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
	}

	result, err := c.service.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		if err.Error() == "FAQ tidak ditemukan" {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Status FAQ berhasil diubah", result)
}

// Reorder - Reorder FAQ item (up/down)
// PATCH /api/panel/faq/:id/reorder
func (c *FAQController) Reorder(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
	}

	var req dto.ReorderDirectionRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Reorder(ctx.UserContext(), id, req.Direction)
	if err != nil {
		if err.Error() == "item sudah berada di urutan paling atas" ||
			err.Error() == "item sudah berada di urutan paling bawah" {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	// Standardized reorder response
	return 	utils.SuccessResponse(ctx, "Urutan berhasil diubah", fiber.Map{
		"item": fiber.Map{
			"id":     result.ItemID.String(),
			"urutan": result.ItemUrutan,
		},
		"swapped_with": fiber.Map{
			"id":     result.SwappedID.String(),
			"urutan": result.SwappedUrutan,
		},
	})
}
