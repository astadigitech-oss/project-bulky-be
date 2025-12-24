package controllers

import (
	"net/http"

	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type WilayahController struct {
	provinsiService  services.ProvinsiService
	kotaService      services.KotaService
	kecamatanService services.KecamatanService
	kelurahanService services.KelurahanService
}

func NewWilayahController(
	provinsiService services.ProvinsiService,
	kotaService services.KotaService,
	kecamatanService services.KecamatanService,
	kelurahanService services.KelurahanService,
) *WilayahController {
	return &WilayahController{
		provinsiService:  provinsiService,
		kotaService:      kotaService,
		kecamatanService: kecamatanService,
		kelurahanService: kelurahanService,
	}
}

func (c *WilayahController) DropdownProvinsi(ctx *gin.Context) {
	items, err := c.provinsiService.FindAllDropdown(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Data provinsi berhasil diambil", items)
}

func (c *WilayahController) DropdownKota(ctx *gin.Context) {
	provinsiID := ctx.Query("provinsi_id")
	if provinsiID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "provinsi_id wajib diisi")
		return
	}

	items, err := c.kotaService.FindByProvinsiID(ctx.Request.Context(), provinsiID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Data kota berhasil diambil", items)
}

func (c *WilayahController) DropdownKecamatan(ctx *gin.Context) {
	kotaID := ctx.Query("kota_id")
	if kotaID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "kota_id wajib diisi")
		return
	}

	items, err := c.kecamatanService.FindByKotaID(ctx.Request.Context(), kotaID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Data kecamatan berhasil diambil", items)
}

func (c *WilayahController) DropdownKelurahan(ctx *gin.Context) {
	kecamatanID := ctx.Query("kecamatan_id")
	if kecamatanID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "kecamatan_id wajib diisi")
		return
	}

	items, err := c.kelurahanService.FindByKecamatanID(ctx.Request.Context(), kecamatanID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Data kelurahan berhasil diambil", items)
}
