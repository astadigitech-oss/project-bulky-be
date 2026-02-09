package dto

import (
	"time"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
)

// TujuanKategoriRequest for input validation
type TujuanKategoriRequest struct {
	ID   uuid.UUID `json:"id" binding:"required"`
	Slug string    `json:"slug" binding:"required,max=100"`
}

// TujuanKategoriResponse for API response
// TypeScript: interface TujuanKategori { id: string; slug: string; }
type TujuanKategoriResponse struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

// Banner Event Promo Request
// Note: gambar_id and gambar_en are handled via multipart form files
type BannerEventPromoRequest struct {
	Nama           string     `json:"nama" binding:"required,min=2,max=100"`
	Tujuan         []string   `json:"tujuan"` // Array of kategori IDs
	Urutan         int        `json:"urutan,omitempty"`
	TanggalMulai   *time.Time `json:"tanggal_mulai,omitempty"`
	TanggalSelesai *time.Time `json:"tanggal_selesai,omitempty"`
}

// Banner Event Promo Response
type BannerEventPromoListResponse struct {
	ID        uuid.UUID                `json:"id"`
	Nama      string                   `json:"nama"`
	GambarURL models.TranslatableImage `json:"gambar_url"`
	Tujuan    []string                 `json:"tujuan"` // Array of kategori ID strings
	Urutan    int                      `json:"urutan"`
	IsVisible bool                     `json:"is_visible"` // computed dari schedule
	UpdatedAt time.Time                `json:"updated_at"`
}

type BannerEventPromoDetailResponse struct {
	ID             uuid.UUID                `json:"id"`
	Nama           string                   `json:"nama"`
	GambarURL      models.TranslatableImage `json:"gambar_url"`
	Tujuan         []string                 `json:"tujuan"` // Array of kategori ID strings
	Urutan         int                      `json:"urutan"`
	IsVisible      bool                     `json:"is_visible"`
	TanggalMulai   *time.Time               `json:"tanggal_mulai,omitempty"`
	TanggalSelesai *time.Time               `json:"tanggal_selesai,omitempty"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

// Converters
func ToBannerEventPromoListResponse(b *models.BannerEventPromo) BannerEventPromoListResponse {
	return BannerEventPromoListResponse{
		ID:        b.ID,
		Nama:      b.Nama,
		GambarURL: b.GetGambarURL(),
		Tujuan:    b.GetKategoriIDStrings(),
		Urutan:    b.Urutan,
		IsVisible: b.IsCurrentlyVisible(),
		UpdatedAt: b.UpdatedAt,
	}
}

func ToBannerEventPromoDetailResponse(b *models.BannerEventPromo) BannerEventPromoDetailResponse {
	return BannerEventPromoDetailResponse{
		ID:             b.ID,
		Nama:           b.Nama,
		GambarURL:      b.GetGambarURL(),
		Tujuan:         b.GetKategoriIDStrings(),
		Urutan:         b.Urutan,
		IsVisible:      b.IsCurrentlyVisible(),
		TanggalMulai:   b.TanggalMulai,
		TanggalSelesai: b.TanggalSelesai,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
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
		Urutan:    0, // Removed field
		IsActive:  h.IsDefault,
		UpdatedAt: h.UpdatedAt,
	}
}

func ToHeroSectionDetailResponse(h *models.HeroSection) HeroSectionDetailResponse {
	return HeroSectionDetailResponse{
		ID:        h.ID,
		Nama:      h.Nama,
		GambarURL: h.GetGambarURL(),
		// LinkURL:   h.LinkURL,
		Urutan:    0, // Removed field
		IsActive:  h.IsDefault,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}
