package controllers

import (
	"errors"
	"net/http"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabelBlogController struct {
	service        services.LabelBlogService
	reorderService *services.ReorderService
	activityLog    services.ActivityLogService
}

func NewLabelBlogController(service services.LabelBlogService, reorderService *services.ReorderService, activityLog services.ActivityLogService) *LabelBlogController {
	return &LabelBlogController{service: service, reorderService: reorderService, activityLog: activityLog}
}

func (c *LabelBlogController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateLabelBlogRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}

	label, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat label", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionCreate, "label_blog", "Label blog berhasil dibuat")
	return utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Label berhasil dibuat", label)
}

func (c *LabelBlogController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req dto.UpdateLabelBlogRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}

	label, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Label tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal memperbarui label", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "label_blog", "Label blog berhasil diupdate")
	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Label berhasil diperbarui", label)
}

func (c *LabelBlogController) Delete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Label tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus label", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionDelete, "label_blog", "Label blog berhasil dihapus")
	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Label berhasil dihapus", nil)
}

func (c *LabelBlogController) GetByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	label, err := c.service.GetByID(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Label tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil label", err.Error())
	}

	return utils.SuccessResponse(ctx, "Label berhasil diambil", label)
}

func (c *LabelBlogController) GetAll(ctx *fiber.Ctx) error {
	var params dto.LabelBlogFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}

	labels, meta, err := c.service.GetAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil label", err.Error())
	}

	return utils.PaginatedSuccessResponse(ctx, "Label berhasil diambil", labels, meta)
}

func (c *LabelBlogController) Reorder(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req dto.ReorderDirectionRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}

	result, err := c.reorderService.Reorder(
		ctx.UserContext(),
		"label_blog",
		id,
		req.Direction,
		"",
		nil,
	)
	if err != nil {
		if err.Error() == "item sudah berada di urutan paling atas" ||
			err.Error() == "item sudah berada di urutan paling bawah" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, err.Error(), "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah urutan", err.Error())
	}

	return 	utils.SuccessResponse(ctx, "Urutan label blog berhasil diubah", fiber.Map{
		"item": fiber.Map{
			"id":     result.ItemID.String(),
			"urutan": result.ItemUrutan,
		},
		"swapped": fiber.Map{
			"id":     result.SwappedID.String(),
			"urutan": result.SwappedUrutan,
		},
	})
}

func (c *LabelBlogController) GetDropdownOptions(ctx *fiber.Ctx) error {
	labels, err := c.service.GetAllActive(ctx.UserContext())
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil label", err.Error())
	}

	return utils.SuccessResponse(ctx, "Label berhasil diambil", labels)
}

func (c *LabelBlogController) GetAllPublic(ctx *fiber.Ctx) error {
	labels, err := c.service.GetAllPublicWithCount(ctx.UserContext())
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil label", err.Error())
	}

	return utils.SuccessResponse(ctx, "Label berhasil diambil", labels)
}
