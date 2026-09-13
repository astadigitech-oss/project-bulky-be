package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var errUnsupportedSeasonalAsset = errors.New("unsupported seasonal campaign asset type")

type SeasonalCampaignController struct {
	service     services.SeasonalCampaignService
	cfg         *config.Config
	activityLog services.ActivityLogService
}

func NewSeasonalCampaignController(service services.SeasonalCampaignService, cfg *config.Config, activityLog services.ActivityLogService) *SeasonalCampaignController {
	return &SeasonalCampaignController{service: service, cfg: cfg, activityLog: activityLog}
}

func (c *SeasonalCampaignController) Create(ctx *fiber.Ctx) error {
	var req models.CreateSeasonalCampaignRequest
	var uploaded []string
	if strings.Contains(ctx.Get("Content-Type"), "multipart/form-data") {
		req.Nama = ctx.FormValue("nama")
		req.TanggalMulai = optionalFormValue(ctx, "tanggal_mulai")
		req.TanggalSelesai = optionalFormValue(ctx, "tanggal_selesai")
		var err error
		if req.WebLogoURL, req.MobileLoadingLogoURL, err = c.saveLogoVariants(ctx); err != nil {
			return assetError(ctx, "web_logo", err)
		}
		if req.WebLogoURL != nil {
			uploaded = append(uploaded, *req.WebLogoURL)
		}
		if req.MobileLoadingLogoURL != nil {
			uploaded = append(uploaded, *req.MobileLoadingLogoURL)
		}
		if req.WebNavbarDecorationURL, err = c.saveOptionalAsset(ctx, "web_navbar_decoration"); err != nil {
			c.deleteUploaded(uploaded)
			return assetError(ctx, "web_navbar_decoration", err)
		}
		if req.WebNavbarDecorationURL != nil {
			uploaded = append(uploaded, *req.WebNavbarDecorationURL)
		}
		if req.MobileTopAppBarOrnamentURL, err = c.saveOptionalAsset(ctx, "mobile_top_app_bar_ornament"); err != nil {
			c.deleteUploaded(uploaded)
			return assetError(ctx, "mobile_top_app_bar_ornament", err)
		}
		if req.MobileTopAppBarOrnamentURL != nil {
			uploaded = append(uploaded, *req.MobileTopAppBarOrnamentURL)
		}
	} else if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}
	result, err := c.service.Create(ctx.UserContext(), &req, adminID(ctx))
	if err != nil {
		c.deleteUploaded(uploaded)
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}
	c.logCreate(ctx, result, "Draft campaign seasonal berhasil dibuat")
	return utils.CreatedResponse(ctx, "Draft campaign seasonal berhasil dibuat", result)
}

func (c *SeasonalCampaignController) FindAll(ctx *fiber.Ctx) error {
	var params models.SeasonalCampaignFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}
	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}
	return utils.PaginatedSuccessResponse(ctx, "Data campaign seasonal berhasil diambil", items, *meta)
}

func (c *SeasonalCampaignController) FindByID(ctx *fiber.Ctx) error {
	result, err := c.service.FindByID(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return seasonalError(ctx, err)
	}
	return utils.SuccessResponse(ctx, "Detail campaign seasonal berhasil diambil", result)
}

func (c *SeasonalCampaignController) Update(ctx *fiber.Ctx) error {
	var req models.UpdateSeasonalCampaignRequest
	var uploaded []string
	if strings.Contains(ctx.Get("Content-Type"), "multipart/form-data") {
		if value := optionalFormValue(ctx, "nama"); value != nil {
			req.Nama = value
		}
		if value := optionalFormValue(ctx, "tanggal_mulai"); value != nil {
			req.TanggalMulai = value
		}
		if value := optionalFormValue(ctx, "tanggal_selesai"); value != nil {
			req.TanggalSelesai = value
		}
		if value, present, err := optionalBoolFormValue(ctx, "remove_web_logo"); err != nil {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		} else if present {
			req.RemoveWebLogo = value
		}
		if value, present, err := optionalBoolFormValue(ctx, "remove_web_navbar_decoration"); err != nil {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		} else if present {
			req.RemoveWebNavbarDecoration = value
		}
		if value, present, err := optionalBoolFormValue(ctx, "remove_mobile_top_app_bar_ornament"); err != nil {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		} else if present {
			req.RemoveMobileTopAppBarOrnament = value
		}
		var err error
		if req.WebLogoURL, req.MobileLoadingLogoURL, err = c.saveLogoVariants(ctx); err != nil {
			return assetError(ctx, "web_logo", err)
		}
		if req.WebLogoURL != nil {
			uploaded = append(uploaded, *req.WebLogoURL)
		}
		if req.MobileLoadingLogoURL != nil {
			uploaded = append(uploaded, *req.MobileLoadingLogoURL)
		}
		if req.WebNavbarDecorationURL, err = c.saveOptionalAsset(ctx, "web_navbar_decoration"); err != nil {
			c.deleteUploaded(uploaded)
			return assetError(ctx, "web_navbar_decoration", err)
		}
		if req.WebNavbarDecorationURL != nil {
			uploaded = append(uploaded, *req.WebNavbarDecorationURL)
		}
		if req.MobileTopAppBarOrnamentURL, err = c.saveOptionalAsset(ctx, "mobile_top_app_bar_ornament"); err != nil {
			c.deleteUploaded(uploaded)
			return assetError(ctx, "mobile_top_app_bar_ornament", err)
		}
		if req.MobileTopAppBarOrnamentURL != nil {
			uploaded = append(uploaded, *req.MobileTopAppBarOrnamentURL)
		}
	} else if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}
	result, err := c.service.Update(ctx.UserContext(), ctx.Params("id"), &req, adminID(ctx))
	if err != nil {
		c.deleteUploaded(uploaded)
		return seasonalError(ctx, err)
	}
	c.logUpdate(ctx, result, "Campaign seasonal berhasil diupdate")
	return utils.SuccessResponse(ctx, "Campaign seasonal berhasil diupdate", result)
}

func (c *SeasonalCampaignController) Publish(ctx *fiber.Ctx) error {
	result, err := c.service.Publish(ctx.UserContext(), ctx.Params("id"), adminID(ctx))
	if err != nil {
		return seasonalError(ctx, err)
	}
	c.logAction(ctx, models.ActionPublish, result, "Campaign seasonal berhasil dipublish")
	return utils.SuccessResponse(ctx, "Campaign seasonal berhasil dipublish", result)
}

func (c *SeasonalCampaignController) Cancel(ctx *fiber.Ctx) error {
	result, err := c.service.Cancel(ctx.UserContext(), ctx.Params("id"), adminID(ctx))
	if err != nil {
		return seasonalError(ctx, err)
	}
	c.logAction(ctx, models.ActionCancel, result, "Campaign seasonal berhasil dibatalkan")
	return utils.SuccessResponse(ctx, "Campaign seasonal berhasil dibatalkan", result)
}

func (c *SeasonalCampaignController) Delete(ctx *fiber.Ctx) error {
	result, err := c.service.Delete(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return seasonalError(ctx, err)
	}
	if id, err := uuid.Parse(result.ID); err == nil {
		c.activityLog.LogDelete(ctx, "seasonal_campaign", "seasonal_campaign", id, "Campaign seasonal berhasil dihapus", result)
	}
	return utils.SuccessResponse(ctx, "Campaign seasonal berhasil dihapus", nil)
}

func (c *SeasonalCampaignController) Preview(ctx *fiber.Ctx) error {
	result, err := c.service.Preview(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return seasonalError(ctx, err)
	}
	return utils.SuccessResponse(ctx, "Preview campaign seasonal berhasil diambil", result)
}

func (c *SeasonalCampaignController) saveOptionalAsset(ctx *fiber.Ctx, field string) (*string, error) {
	file, err := ctx.FormFile(field)
	if err != nil {
		return nil, nil
	}
	target, ok := seasonalAssetTarget(field)
	if !ok {
		return nil, errUnsupportedSeasonalAsset
	}
	path, err := utils.SaveSeasonalAsset(file, c.cfg, target)
	if err != nil {
		return nil, err
	}
	return &path, nil
}

func seasonalAssetTarget(field string) (utils.SeasonalAssetTarget, bool) {
	switch field {
	case "web_navbar_decoration":
		return utils.SeasonalWebNavbarDecorationTarget, true
	case "mobile_top_app_bar_ornament":
		return utils.SeasonalMobileTopAppBarTarget, true
	default:
		return utils.SeasonalAssetTarget{}, false
	}
}

func (c *SeasonalCampaignController) saveLogoVariants(ctx *fiber.Ctx) (*string, *string, error) {
	file, err := ctx.FormFile("web_logo")
	if err != nil {
		return nil, nil, nil
	}
	variants, err := utils.SaveSeasonalLogoVariants(file, c.cfg)
	if err != nil {
		return nil, nil, err
	}
	return &variants.WebLogoPath, &variants.MobileLoadingLogoPath, nil
}

func (c *SeasonalCampaignController) deleteUploaded(paths []string) {
	for _, path := range paths {
		_ = utils.DeleteFile(path, c.cfg)
	}
}

func (c *SeasonalCampaignController) logCreate(ctx *fiber.Ctx, result *models.SeasonalCampaignResponse, description string) {
	if id, err := uuid.Parse(result.ID); err == nil {
		c.activityLog.LogCreate(ctx, "seasonal_campaign", "seasonal_campaign", id, description, result)
	}
}
func (c *SeasonalCampaignController) logUpdate(ctx *fiber.Ctx, result *models.SeasonalCampaignResponse, description string) {
	if id, err := uuid.Parse(result.ID); err == nil {
		c.activityLog.LogUpdate(ctx, "seasonal_campaign", "seasonal_campaign", id, description, nil, result)
	}
}
func (c *SeasonalCampaignController) logAction(ctx *fiber.Ctx, action models.ActivityAction, result *models.SeasonalCampaignResponse, description string) {
	if id, err := uuid.Parse(result.ID); err == nil {
		c.activityLog.Log(ctx, action, "seasonal_campaign", description, services.WithEntity("seasonal_campaign", id))
	}
}

func optionalFormValue(ctx *fiber.Ctx, key string) *string {
	form, err := ctx.MultipartForm()
	if err != nil || form == nil {
		return nil
	}
	values, ok := form.Value[key]
	if !ok || len(values) == 0 {
		return nil
	}
	return &values[0]
}
func optionalBoolFormValue(ctx *fiber.Ctx, key string) (*bool, bool, error) {
	value := optionalFormValue(ctx, key)
	if value == nil {
		return nil, false, nil
	}
	parsed, err := strconv.ParseBool(*value)
	if err != nil {
		return nil, true, fmtError("%s harus bernilai true atau false", key)
	}
	return &parsed, true, nil
}
func adminID(ctx *fiber.Ctx) *uuid.UUID {
	id, err := uuid.Parse(localsString(ctx, "user_id"))
	if err != nil {
		return nil
	}
	return &id
}
func assetError(ctx *fiber.Ctx, field string, err error) error {
	if errors.Is(err, errUnsupportedSeasonalAsset) {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file "+field+" tidak didukung. Gunakan jpg, png, atau webp (SVG tidak diizinkan)", nil)
	}
	return utils.ErrorResponse(ctx, http.StatusBadRequest, "Gagal menyimpan file "+field+": "+err.Error(), nil)
}
func seasonalError(ctx *fiber.Ctx, err error) error {
	if err.Error() == "campaign seasonal tidak ditemukan" {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}
	return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
}
func fmtError(format string, args ...interface{}) error { return fmt.Errorf(format, args...) }
