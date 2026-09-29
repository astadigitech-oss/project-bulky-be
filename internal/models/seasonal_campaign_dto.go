package models

import "time"

type CreateSeasonalCampaignRequest struct {
	Nama                       string  `json:"nama" binding:"required,min=1,max=100"`
	TanggalMulai               *string `json:"tanggal_mulai"`
	TanggalSelesai             *string `json:"tanggal_selesai"`
	WebLogoURL                 *string `json:"-"`
	MobileLoadingLogoURL       *string `json:"-"`
	WebNavbarDecorationURL     *string `json:"-"`
	MobileTopAppBarOrnamentURL *string `json:"-"`
}

type UpdateSeasonalCampaignRequest struct {
	Nama                          *string `json:"nama" binding:"omitempty,min=1,max=100"`
	TanggalMulai                  *string `json:"tanggal_mulai"`
	TanggalSelesai                *string `json:"tanggal_selesai"`
	WebLogoURL                    *string `json:"-"`
	MobileLoadingLogoURL          *string `json:"-"`
	WebNavbarDecorationURL        *string `json:"-"`
	MobileTopAppBarOrnamentURL    *string `json:"-"`
	RemoveWebLogo                 *bool   `json:"remove_web_logo"`
	RemoveWebNavbarDecoration     *bool   `json:"remove_web_navbar_decoration"`
	RemoveMobileTopAppBarOrnament *bool   `json:"remove_mobile_top_app_bar_ornament"`
}

type SeasonalCampaignFilterRequest struct {
	PaginationRequest
	Status         string     `query:"status"`
	TanggalMulai   *time.Time `query:"tanggal_mulai"`
	TanggalSelesai *time.Time `query:"tanggal_selesai"`
}

func (p *SeasonalCampaignFilterRequest) SetDefaults() {
	p.PaginationRequest.SetDefaults()
	if p.SortBy != "created_at" && p.SortBy != "updated_at" && p.SortBy != "nama" && p.SortBy != "tanggal_mulai" && p.SortBy != "tanggal_selesai" {
		p.SortBy = "created_at"
	}
	if p.Order != "asc" && p.Order != "desc" {
		p.Order = "desc"
	}
	switch p.Status {
	case "", "draft", "cancelled", "scheduled", "active", "ended":
	default:
		p.Status = ""
	}
}

type SeasonalCampaignAssetsResponse struct {
	WebLogoURL                 *string `json:"web_logo_url"`
	MobileLoadingLogoURL       *string `json:"mobile_loading_logo_url"`
	WebNavbarDecorationURL     *string `json:"web_navbar_decoration_url"`
	MobileTopAppBarOrnamentURL *string `json:"mobile_top_app_bar_ornament_url"`
}

type SeasonalCampaignResponse struct {
	ID             string                         `json:"id"`
	Nama           string                         `json:"nama"`
	Status         string                         `json:"status"`
	TanggalMulai   *time.Time                     `json:"tanggal_mulai"`
	TanggalSelesai *time.Time                     `json:"tanggal_selesai"`
	Assets         SeasonalCampaignAssetsResponse `json:"assets"`
	CreatedAt      time.Time                      `json:"created_at"`
	UpdatedAt      time.Time                      `json:"updated_at"`
}
