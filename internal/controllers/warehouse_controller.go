package controllers

import (
	"net/http"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type WarehouseController struct {
	service services.WarehouseService
}

func NewWarehouseController(service services.WarehouseService) *WarehouseController {
	return &WarehouseController{service: service}
}

func (c *WarehouseController) Create(ctx *fiber.Ctx) error {
	var req models.CreateWarehouseRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
	}

	return utils.CreatedResponse(ctx, "Warehouse berhasil dibuat", result)
}

func (c *WarehouseController) FindAll(ctx *fiber.Ctx) error {
	var params models.PaginationRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	kota := ctx.Query("kota")

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params, kota)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data warehouse berhasil diambil", items, *meta)
}

func (c *WarehouseController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail warehouse berhasil diambil", result)
}

func (c *WarehouseController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateWarehouseRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Warehouse berhasil diupdate", result)
}

func (c *WarehouseController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "warehouse tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Warehouse berhasil dihapus", nil)
}

func (c *WarehouseController) ToggleStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Status warehouse berhasil diubah", result)
}

// Get returns the first active warehouse (singleton pattern)
func (c *WarehouseController) Get(ctx *fiber.Ctx) error {
	result, err := c.service.Get(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data warehouse berhasil diambil", result)
}

// UpdateSingleton updates the first active warehouse (singleton pattern)
func (c *WarehouseController) UpdateSingleton(ctx *fiber.Ctx) error {
	var req dto.WarehouseUpdateRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.UpdateSingleton(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data warehouse berhasil diupdate", result)
}

// GetPublic returns simplified warehouse data for public
func (c *WarehouseController) GetPublic(ctx *fiber.Ctx) error {
	result, err := c.service.GetPublic(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data warehouse berhasil diambil", result)
}

// GetInformasiPickup returns warehouse + jadwal for public informasi pickup endpoint
func (c *WarehouseController) GetInformasiPickup(ctx *fiber.Ctx) error {
	result, err := c.service.GetInformasiPickup(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data informasi pickup berhasil diambil", result)
}

// GetJadwal returns jadwal gudang as array
func (c *WarehouseController) GetJadwal(ctx *fiber.Ctx) error {
	result, err := c.service.GetJadwal(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Data jadwal gudang berhasil diambil", result)
}

// UpdateJadwal updates jadwal gudang and returns updated array
func (c *WarehouseController) UpdateJadwal(ctx *fiber.Ctx) error {
	var req dto.UpdateJadwalRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.UpdateJadwal(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Jadwal gudang berhasil diupdate", result)
}

// Dropdown returns all active warehouses for dropdown selection
// This endpoint is prepared for future scalability when multiple warehouses exist
func (c *WarehouseController) Dropdown(ctx *fiber.Ctx) error {
	// Get all active warehouses for dropdown
	var params models.PaginationRequest
	params.Page = 1
	params.PerPage = 1000 // Get all
	isActive := true
	params.IsActive = &isActive

	warehouseList, _, err := c.service.FindAll(ctx.UserContext(), &params, "")
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data warehouse", nil)
	}

	// Convert to simple dropdown response
	response := make([]map[string]interface{}, len(warehouseList))
	for i, w := range warehouseList {
		response[i] = map[string]interface{}{
			"id":   w.ID,
			"nama": w.Nama,
		}
	}

	return utils.SuccessResponse(ctx, "Data dropdown warehouse berhasil diambil", response)
}
