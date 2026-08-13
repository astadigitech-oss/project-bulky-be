package controllers

import (
	"fmt"
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// ForwarderMappingController mengelola master data mapping kota & kecamatan Forwarder
// (forwarder_city_mapping & forwarder_subdistrict_mapping) di admin panel:
// lihat data hasil sync + tombol Sync dari API Forwarder.
type ForwarderMappingController struct {
	service     services.ForwarderMappingService
	activityLog services.ActivityLogService
}

func NewForwarderMappingController(service services.ForwarderMappingService, activityLog services.ActivityLogService) *ForwarderMappingController {
	return &ForwarderMappingController{
		service:     service,
		activityLog: activityLog,
	}
}

func (c *ForwarderMappingController) FindCities(ctx *fiber.Ctx) error {
	var params models.ForwarderCityFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}
	params.SetDefaults()

	items, meta, err := c.service.FindCities(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data mapping kota Forwarder berhasil diambil", items, *meta)
}

func (c *ForwarderMappingController) FindSubdistricts(ctx *fiber.Ctx) error {
	var params models.ForwarderSubdistrictFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}
	params.SetDefaults()

	items, meta, err := c.service.FindSubdistricts(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data mapping kecamatan Forwarder berhasil diambil", items, *meta)
}

// Sync menarik data terbaru dari API Forwarder (citylist &/atau subdistrictlist)
// dan upsert ke tabel mapping. Body opsional; default sync keduanya.
// Hanya bisa diakses role dengan permission forwarder_mapping:manage (Super Admin).
func (c *ForwarderMappingController) Sync(ctx *fiber.Ctx) error {
	var req models.SyncForwarderMappingRequest
	// Body opsional — kalau kosong, kedua bagian di-sync.
	if len(ctx.Body()) > 0 {
		if err := BindJSON(ctx, &req); err != nil {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		}
	}

	result, err := c.service.Sync(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadGateway, err.Error(), nil)
	}

	detail := fmt.Sprintf(
		"Sync mapping Forwarder: %d kota dibuat, %d diupdate (%d dari API); %d kecamatan dibuat, %d diupdate (%d dari API)",
		result.CityCreated, result.CityUpdated, result.CityTotalFromAPI,
		result.SubdistrictCreated, result.SubdistrictUpdated, result.SubdistrictTotalFromAPI,
	)
	c.activityLog.Log(ctx, models.ActionImport, "forwarder_mapping", detail)
	return utils.SuccessResponse(ctx, "Sync mapping Forwarder berhasil", result)
}
