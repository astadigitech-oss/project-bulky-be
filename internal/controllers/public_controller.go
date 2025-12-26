package controllers

import (
	"net/http"

	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PublicController struct {
	heroService   services.HeroSectionService
	bannerService services.BannerEventPromoService
}

func NewPublicController(heroService services.HeroSectionService, bannerService services.BannerEventPromoService) *PublicController {
	return &PublicController{
		heroService:   heroService,
		bannerService: bannerService,
	}
}

func (c *PublicController) GetActiveHero(ctx *gin.Context) {
	result, err := c.heroService.GetVisibleHero(ctx.Request.Context())
	if err != nil || result == nil {
		utils.SuccessResponse(ctx, "Tidak ada hero section aktif", nil)
		return
	}

	utils.SuccessResponse(ctx, "Hero section berhasil diambil", result)
}

func (c *PublicController) GetActiveBanners(ctx *gin.Context) {
	result, err := c.bannerService.GetVisibleBanners(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Banner berhasil diambil", result)
}
