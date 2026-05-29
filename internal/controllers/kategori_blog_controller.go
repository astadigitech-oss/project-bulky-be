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

type KategoriBlogController struct {
	service        services.KategoriBlogService
	reorderService *services.ReorderService
	activityLog    services.ActivityLogService
}

func NewKategoriBlogController(service services.KategoriBlogService, reorderService *services.ReorderService, activityLog services.ActivityLogService) *KategoriBlogController {
	return &KategoriBlogController{
		service:        service,
		reorderService: reorderService,
		activityLog:    activityLog,
	}
}

func (c *KategoriBlogController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateKategoriBlogRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Request tidak valid", err.Error())
	}

	kategori, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat kategori", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionCreate, "kategori_blog", "Kategori blog berhasil dibuat")
	return utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Kategori berhasil dibuat", kategori)
}

func (c *KategoriBlogController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req dto.UpdateKategoriBlogRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Request tidak valid", err.Error())
	}

	kategori, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kategori tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate kategori", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "kategori_blog", "Kategori blog berhasil diupdate")
	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Kategori berhasil diupdate", kategori)
}

func (c *KategoriBlogController) Delete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kategori tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus kategori", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionDelete, "kategori_blog", "Kategori blog berhasil dihapus")
	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Kategori berhasil dihapus", nil)
}

func (c *KategoriBlogController) GetByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	kategori, err := c.service.GetByID(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kategori tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
	}

	return utils.SuccessResponse(ctx, "Kategori berhasil diambil", kategori)
}

func (c *KategoriBlogController) GetAll(ctx *fiber.Ctx) error {
	var params dto.KategoriBlogFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}

	kategoris, meta, err := c.service.GetAllPaginated(ctx.UserContext(), &params)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
	}

	return utils.PaginatedSuccessResponse(ctx, "Kategori berhasil diambil", kategoris, meta)
}

func (c *KategoriBlogController) ToggleStatus(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	if err := c.service.ToggleStatus(ctx.UserContext(), id); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah status", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "kategori_blog", "Status kategori blog berhasil diubah")
	return utils.SuccessResponse(ctx, "Status kategori berhasil diubah", nil)
}

func (c *KategoriBlogController) Reorder(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req dto.ReorderDirectionRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Request tidak valid", err.Error())
	}

	result, err := c.reorderService.Reorder(
		ctx.UserContext(),
		"kategori_blog",
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

	return 	utils.SuccessResponse(ctx, "Urutan kategori berhasil diubah", fiber.Map{
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

func (c *KategoriBlogController) GetDropdownOptions(ctx *fiber.Ctx) error {
	kategoris, err := c.service.GetAllActive(ctx.UserContext())
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
	}

	return utils.SuccessResponse(ctx, "Kategori berhasil diambil", kategoris)
}

func (c *KategoriBlogController) GetAllPublic(ctx *fiber.Ctx) error {
	kategoris, err := c.service.GetAllPublicWithCount(ctx.UserContext())
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
	}

	return utils.SuccessResponse(ctx, "Kategori berhasil diambil", kategoris)
}
