package dto

import (
	"time"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
)

// Banner Event Promo Request
type BannerEventPromoRequest struct {
	Nama         string     `json:"nama" binding:"required,min=2,max=100"`
	LinkURL      *string    `json:"link_url,omitempty"`
	Urutan       int        `json:"urutan,omitempty"`
	IsActive     *bool      `json:"is_active,omitempty"`
	TanggalMulai *time.Time `json:"tanggal_mulai,omitempty"`
	TanggalAkhir *time.Time `json:"tanggal_akhir,omitempty"`
}

// Banner Event Promo Response
type BannerEventPromoListResponse struct {
	ID        uuid.UUID                `json:"id"`
	Nama      string                   `json:"nama"`
	GambarURL models.TranslatableImage `json:"gambar_url"`
	LinkURL   *string                  `json:"link_url,omitempty"`
	Urutan    int                      `json:"urutan"`
	IsActive  bool                     `json:"is_active"`
	UpdatedAt time.Time                `json:"updated_at"`
}

type BannerEventPromoDetailResponse struct {
	ID           uuid.UUID                `json:"id"`
	Nama         string                   `json:"nama"`
	GambarURL    models.TranslatableImage `json:"gambar_url"`
	LinkURL      *string                  `json:"link_url,omitempty"`
	Urutan       int                      `json:"urutan"`
	IsActive     bool                     `json:"is_active"`
	TanggalMulai *time.Time               `json:"tanggal_mulai,omitempty"`
	TanggalAkhir *time.Time               `json:"tanggal_akhir,omitempty"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

// Converters
func ToBannerEventPromoListResponse(b *models.BannerEventPromo) BannerEventPromoListResponse {
	return BannerEventPromoListResponse{
		ID:        b.ID,
		Nama:      b.Nama,
		GambarURL: b.GetGambarURL(),
		LinkURL:   b.LinkURL,
		Urutan:    b.Urutan,
		IsActive:  b.IsActive,
		UpdatedAt: b.UpdatedAt,
	}
}

func ToBannerEventPromoDetailResponse(b *models.BannerEventPromo) BannerEventPromoDetailResponse {
	return BannerEventPromoDetailResponse{
		ID:           b.ID,
		Nama:         b.Nama,
		GambarURL:    b.GetGambarURL(),
		LinkURL:      b.LinkURL,
		Urutan:       b.Urutan,
		IsActive:     b.IsActive,
		TanggalMulai: b.TanggalMulai,
		TanggalAkhir: b.TanggalAkhir,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
}

// Hero Section Request
type HeroSectionRequest struct {
	Nama     string  `json:"nama" binding:"required,min=2,max=100"`
	LinkURL  *string `json:"link_url,omitempty"`
	Urutan   int     `json:"urutan,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// Hero Section Response
type HeroSectionListResponse struct {
	ID        uuid.UUID                `json:"id"`
	Nama      string                   `json:"nama"`
	GambarURL models.TranslatableImage `json:"gambar_url"`
	LinkURL   *string                  `json:"link_url,omitempty"`
	Urutan    int                      `json:"urutan"`
	IsActive  bool                     `json:"is_active"`
	UpdatedAt time.Time                `json:"updated_at"`
}

type HeroSectionDetailResponse struct {
	ID        uuid.UUID                `json:"id"`
	Nama      string                   `json:"nama"`
	GambarURL models.TranslatableImage `json:"gambar_url"`
	LinkURL   *string                  `json:"link_url,omitempty"`
	Urutan    int                      `json:"urutan"`
	IsActive  bool                     `json:"is_active"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

// Converters
func ToHeroSectionListResponse(h *models.HeroSection) HeroSectionListResponse {
	return HeroSectionListResponse{
		ID:        h.ID,
		Nama:      h.Nama,
		GambarURL: h.GetGambarURL(),
		// LinkURL:   h.LinkURL,
		Urutan:    h.Urutan,
		IsActive:  h.IsActive,
		UpdatedAt: h.UpdatedAt,
	}
}

func ToHeroSectionDetailResponse(h *models.HeroSection) HeroSectionDetailResponse {
	return HeroSectionDetailResponse{
		ID:        h.ID,
		Nama:      h.Nama,
		GambarURL: h.GetGambarURL(),
		// LinkURL:   h.LinkURL,
		Urutan:    h.Urutan,
		IsActive:  h.IsActive,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}
