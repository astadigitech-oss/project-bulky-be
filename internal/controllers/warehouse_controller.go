package controllers

import (
	"net/http"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type WarehouseController struct {
	service services.WarehouseService
}

func NewWarehouseController(service services.WarehouseService) *WarehouseController {
	return &WarehouseController{service: service}
}

func (c *WarehouseController) Create(ctx *gin.Context) {
	var req models.CreateWarehouseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Warehouse berhasil dibuat", result)
}

func (c *WarehouseController) FindAll(ctx *gin.Context) {
	var params models.PaginationRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	kota := ctx.Query("kota")

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params, kota)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data warehouse berhasil diambil", items, *meta)
}


func (c *WarehouseController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail warehouse berhasil diambil", result)
}

func (c *WarehouseController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateWarehouseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Warehouse berhasil diupdate", result)
}

func (c *WarehouseController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "warehouse tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Warehouse berhasil dihapus", nil)
}

func (c *WarehouseController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status warehouse berhasil diubah", result)
}

// Get returns the first active warehouse (singleton pattern)
func (c *WarehouseController) Get(ctx *gin.Context) {
	result, err := c.service.Get(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data warehouse berhasil diambil", result)
}

// UpdateSingleton updates the first active warehouse (singleton pattern)
func (c *WarehouseController) UpdateSingleton(ctx *gin.Context) {
	var req dto.WarehouseUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.UpdateSingleton(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data warehouse berhasil diupdate", result)
}

// GetPublic returns simplified warehouse data for public
func (c *WarehouseController) GetPublic(ctx *gin.Context) {
	result, err := c.service.GetPublic(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data warehouse berhasil diambil", result)
}

// GetInformasiPickup returns warehouse + jadwal for public informasi pickup endpoint
func (c *WarehouseController) GetInformasiPickup(ctx *gin.Context) {
	result, err := c.service.GetInformasiPickup(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data informasi pickup berhasil diambil", result)
}

// UpdateJadwal updates jadwal gudang
func (c *WarehouseController) UpdateJadwal(ctx *gin.Context) {
	var req dto.UpdateJadwalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	if err := c.service.UpdateJadwal(ctx.Request.Context(), &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Jadwal berhasil diupdate", nil)
}
