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

type KategoriVideoController struct {
	service        services.KategoriVideoService
	reorderService *services.ReorderService
	activityLog    services.ActivityLogService
}

func NewKategoriVideoController(service services.KategoriVideoService, reorderService *services.ReorderService, activityLog services.ActivityLogService) *KategoriVideoController {
	return &KategoriVideoController{
		service:        service,
		reorderService: reorderService,
		activityLog:    activityLog,
	}
}

func (c *KategoriVideoController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateKategoriVideoRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Request tidak valid", utils.GetValidationErrorMessage(err))
	}

	kategori, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat kategori", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionCreate, "kategori_video", "Kategori video berhasil dibuat")
	return utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Kategori berhasil dibuat", kategori)
}

func (c *KategoriVideoController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req dto.UpdateKategoriVideoRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Request tidak valid", utils.GetValidationErrorMessage(err))
	}

	kategori, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kategori tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate kategori", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "kategori_video", "Kategori video berhasil diupdate")
	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Kategori berhasil diupdate", kategori)
}

func (c *KategoriVideoController) Delete(ctx *fiber.Ctx) error {
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

	c.activityLog.Log(ctx, models.ActionDelete, "kategori_video", "Kategori video berhasil dihapus")
	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Kategori berhasil dihapus", nil)
}

func (c *KategoriVideoController) GetByID(ctx *fiber.Ctx) error {
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

func (c *KategoriVideoController) GetAll(ctx *fiber.Ctx) error {
	var params dto.KategoriVideoFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}

	kategoris, meta, err := c.service.GetAllPaginated(ctx.UserContext(), &params)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
	}

	return utils.PaginatedSuccessResponse(ctx, "Kategori berhasil diambil", kategoris, meta)
}

func (c *KategoriVideoController) ToggleStatus(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	if err := c.service.ToggleStatus(ctx.UserContext(), id); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah status", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "kategori_video", "Status kategori video berhasil diubah")
	return utils.SuccessResponse(ctx, "Status kategori berhasil diubah", nil)
}

func (c *KategoriVideoController) GetDropdownOptions(ctx *fiber.Ctx) error {
	kategoris, err := c.service.GetAllActive(ctx.UserContext())
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data kategori", err.Error())
	}

	return utils.SuccessResponse(ctx, "Data dropdown kategori video berhasil diambil", kategoris)
}

func (c *KategoriVideoController) Reorder(ctx *fiber.Ctx) error {
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
		"kategori_video",
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

func (c *KategoriVideoController) GetAllPublic(ctx *fiber.Ctx) error {
	kategoris, err := c.service.GetAllPublicWithCount(ctx.UserContext())
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
	}

	return utils.SuccessResponse(ctx, "Kategori berhasil diambil", kategoris)
}
