package controllers

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"
	"strconv"
	"strings"
	"time"
)

type AuctionEducationBannerController struct {
	service  services.AuctionEducationBannerService
	cfg      *config.Config
	activity services.ActivityLogService
}

func NewAuctionEducationBannerController(s services.AuctionEducationBannerService, c *config.Config, a services.ActivityLogService) *AuctionEducationBannerController {
	return &AuctionEducationBannerController{s, c, a}
}
func auctionBannerDates(c *fiber.Ctx) (*time.Time, *time.Time, error) {
	parse := func(n string) (*time.Time, error) {
		v := strings.TrimSpace(c.FormValue(n))
		if v == "" {
			return nil, nil
		}
		t, e := time.Parse(time.RFC3339, v)
		if e != nil {
			return nil, e
		}
		return &t, nil
	}
	a, e := parse("tanggal_mulai")
	if e != nil {
		return nil, nil, e
	}
	b, e := parse("tanggal_selesai")
	return a, b, e
}
func (c *AuctionEducationBannerController) upload(ctx *fiber.Ctx, name string, required bool) (*string, error) {
	f, e := ctx.FormFile(name)
	if e != nil {
		if required {
			return nil, auctionBannerError("file " + name + " wajib diupload")
		}
		return nil, nil
	}
	if !utils.IsValidBannerImageType(f) {
		return nil, auctionBannerError("format gambar harus JPG, PNG, atau WebP")
	}
	p, e := utils.CompressAndSaveImageWebP(f, "auction-education-banners", c.cfg)
	if e != nil {
		return nil, e
	}
	return &p, nil
}

func auctionBannerOrder(ctx *fiber.Ctx) (*int, error) {
	v := strings.TrimSpace(ctx.FormValue("urutan"))
	if v == "" {
		return nil, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return nil, auctionBannerError("urutan harus berupa bilangan bulat tidak negatif")
	}
	return &n, nil
}

type auctionBannerError string

func (e auctionBannerError) Error() string { return string(e) }
func (c *AuctionEducationBannerController) Create(ctx *fiber.Ctx) error {
	start, end, e := auctionBannerDates(ctx)
	if e != nil {
		return utils.ErrorResponse(ctx, 400, "format tanggal tidak valid", nil)
	}
	id, e := c.upload(ctx, "gambar_id", true)
	if e != nil {
		return utils.ErrorResponse(ctx, 400, e.Error(), nil)
	}
	en, e := c.upload(ctx, "gambar_en", true)
	if e != nil {
		utils.DeleteFile(*id, c.cfg)
		return utils.ErrorResponse(ctx, 400, e.Error(), nil)
	}
	order, e := auctionBannerOrder(ctx)
	if e != nil {
		utils.DeleteFile(*id, c.cfg)
		utils.DeleteFile(*en, c.cfg)
		return utils.ErrorResponse(ctx, 400, e.Error(), nil)
	}
	urutan := 0
	if order != nil {
		urutan = *order
	}
	r, e := c.service.Create(ctx.UserContext(), ctx.FormValue("nama"), *id, *en, urutan, start, end)
	if e != nil {
		utils.DeleteFile(*id, c.cfg)
		utils.DeleteFile(*en, c.cfg)
		return utils.ErrorResponse(ctx, 400, e.Error(), nil)
	}
	c.activity.Log(ctx, models.ActionCreate, "auction_education_banner", "Banner edukasi lelang dibuat")
	return utils.CreatedResponse(ctx, "Banner edukasi lelang berhasil dibuat", r)
}
func (c *AuctionEducationBannerController) List(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	per := ctx.QueryInt("per_page", 20)
	r, m, e := c.service.List(ctx.UserContext(), ctx.Query("search"), page, per)
	if e != nil {
		return utils.ErrorResponse(ctx, 500, "Gagal mengambil banner", nil)
	}
	return utils.PaginatedSuccessResponse(ctx, "Data banner edukasi lelang berhasil diambil", r, *m)
}
func (c *AuctionEducationBannerController) Get(ctx *fiber.Ctx) error {
	r, e := c.service.Get(ctx.UserContext(), ctx.Params("id"))
	if e != nil {
		return utils.ErrorResponse(ctx, 404, e.Error(), nil)
	}
	return utils.SuccessResponse(ctx, "Detail banner edukasi lelang berhasil diambil", r)
}
func (c *AuctionEducationBannerController) Update(ctx *fiber.Ctx) error {
	start, end, e := auctionBannerDates(ctx)
	if e != nil {
		return utils.ErrorResponse(ctx, 400, "format tanggal tidak valid", nil)
	}
	var nama *string
	if v := ctx.FormValue("nama"); v != "" {
		nama = &v
	}
	id, e := c.upload(ctx, "gambar_id", false)
	if e != nil {
		return utils.ErrorResponse(ctx, 400, e.Error(), nil)
	}
	en, e := c.upload(ctx, "gambar_en", false)
	if e != nil {
		if id != nil {
			utils.DeleteFile(*id, c.cfg)
		}
		return utils.ErrorResponse(ctx, 400, e.Error(), nil)
	}
	order, e := auctionBannerOrder(ctx)
	if e != nil {
		if id != nil {
			utils.DeleteFile(*id, c.cfg)
		}
		if en != nil {
			utils.DeleteFile(*en, c.cfg)
		}
		return utils.ErrorResponse(ctx, 400, e.Error(), nil)
	}
	r, e := c.service.Update(ctx.UserContext(), ctx.Params("id"), nama, id, en, order, start, end)
	if e != nil {
		if id != nil {
			utils.DeleteFile(*id, c.cfg)
		}
		if en != nil {
			utils.DeleteFile(*en, c.cfg)
		}
		return utils.ErrorResponse(ctx, 400, e.Error(), nil)
	}
	c.activity.Log(ctx, models.ActionUpdate, "auction_education_banner", "Banner edukasi lelang diperbarui")
	return utils.SuccessResponse(ctx, "Banner edukasi lelang berhasil diperbarui", r)
}
func (c *AuctionEducationBannerController) Reorder(ctx *fiber.Ctx) error {
	var request struct {
		IDs []string `json:"ids"`
	}
	if err := ctx.BodyParser(&request); err != nil {
		return utils.ErrorResponse(ctx, 400, "payload reorder tidak valid", nil)
	}
	result, err := c.service.Reorder(ctx.UserContext(), request.IDs)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, err.Error(), nil)
	}
	c.activity.Log(ctx, models.ActionUpdate, "auction_education_banner", "Urutan banner edukasi lelang diperbarui")
	return utils.SuccessResponse(ctx, "Urutan banner edukasi lelang berhasil diperbarui", result)
}
func (c *AuctionEducationBannerController) Delete(ctx *fiber.Ctx) error {
	if e := c.service.Delete(ctx.UserContext(), ctx.Params("id")); e != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, e.Error(), nil)
	}
	c.activity.Log(ctx, models.ActionDelete, "auction_education_banner", "Banner edukasi lelang dihapus")
	return utils.SuccessResponse(ctx, "Banner edukasi lelang berhasil dihapus", nil)
}
