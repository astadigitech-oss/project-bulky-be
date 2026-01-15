package dto

import (
	"time"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
)

// Response untuk list (simplified)
type MerekProdukListResponse struct {
	ID           uuid.UUID                 `json:"id"`
	Nama         models.TranslatableString `json:"nama"`
	Slug         string                    `json:"slug"`
	LogoURL      *string                   `json:"logo_url,omitempty"`
	JumlahProduk int                       `json:"jumlah_produk"`
	IsActive     bool                      `json:"is_active"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

// Response untuk detail
type MerekProdukDetailResponse struct {
	ID           uuid.UUID                 `json:"id"`
	Nama         models.TranslatableString `json:"nama"`
	Slug         string                    `json:"slug"`
	LogoURL      *string                   `json:"logo_url,omitempty"`
	JumlahProduk int                       `json:"jumlah_produk"`
	IsActive     bool                      `json:"is_active"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

// Converter dari model ke response
func ToMerekProdukListResponse(m *models.MerekProduk, jumlahProduk int) MerekProdukListResponse {
	return MerekProdukListResponse{
		ID:           m.ID,
		Nama:         m.GetNama(),
		Slug:         m.Slug,
		LogoURL:      m.LogoURL,
		JumlahProduk: jumlahProduk,
		IsActive:     m.IsActive,
		UpdatedAt:    m.UpdatedAt,
	}
}

func ToMerekProdukDetailResponse(m *models.MerekProduk, jumlahProduk int) MerekProdukDetailResponse {
	return MerekProdukDetailResponse{
		ID:           m.ID,
		Nama:         m.GetNama(),
		Slug:         m.Slug,
		LogoURL:      m.LogoURL,
		JumlahProduk: jumlahProduk,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// Kategori Produk Response
type KategoriProdukListResponse struct {
	ID                      uuid.UUID                 `json:"id"`
	Nama                    models.TranslatableString `json:"nama"`
	Slug                    string                    `json:"slug"`
	IconURL                 *string                   `json:"icon_url,omitempty"`
	MemilikiKondisiTambahan bool                      `json:"memiliki_kondisi_tambahan"`
	JumlahProduk            int                       `json:"jumlah_produk"`
	IsActive                bool                      `json:"is_active"`
	UpdatedAt               time.Time                 `json:"updated_at"`
}

type KategoriProdukDetailResponse struct {
	ID                      uuid.UUID                 `json:"id"`
	Nama                    models.TranslatableString `json:"nama"`
	Slug                    string                    `json:"slug"`
	Deskripsi               *string                   `json:"deskripsi,omitempty"`
	IconURL                 *string                   `json:"icon_url,omitempty"`
	MemilikiKondisiTambahan bool                      `json:"memiliki_kondisi_tambahan"`
	TipeKondisiTambahan     *string                   `json:"tipe_kondisi_tambahan,omitempty"`
	GambarKondisiURL        *string                   `json:"gambar_kondisi_url,omitempty"`
	TeksKondisi             *string                   `json:"teks_kondisi,omitempty"`
	JumlahProduk            int                       `json:"jumlah_produk"`
	IsActive                bool                      `json:"is_active"`
	CreatedAt               time.Time                 `json:"created_at"`
	UpdatedAt               time.Time                 `json:"updated_at"`
}

func ToKategoriProdukListResponse(k *models.KategoriProduk, jumlahProduk int) KategoriProdukListResponse {
	return KategoriProdukListResponse{
		ID:                      k.ID,
		Nama:                    k.GetNama(),
		Slug:                    k.Slug,
		IconURL:                 k.IconURL,
		MemilikiKondisiTambahan: k.MemilikiKondisiTambahan,
		JumlahProduk:            jumlahProduk,
		IsActive:                k.IsActive,
		UpdatedAt:               k.UpdatedAt,
	}
}

func ToKategoriProdukDetailResponse(k *models.KategoriProduk, jumlahProduk int) KategoriProdukDetailResponse {
	return KategoriProdukDetailResponse{
		ID:                      k.ID,
		Nama:                    k.GetNama(),
		Slug:                    k.Slug,
		Deskripsi:               k.Deskripsi,
		IconURL:                 k.IconURL,
		MemilikiKondisiTambahan: k.MemilikiKondisiTambahan,
		TipeKondisiTambahan:     k.TipeKondisiTambahan,
		GambarKondisiURL:        k.GambarKondisiURL,
		TeksKondisi:             k.TeksKondisi,
		JumlahProduk:            jumlahProduk,
		IsActive:                k.IsActive,
		CreatedAt:               k.CreatedAt,
		UpdatedAt:               k.UpdatedAt,
	}
}

// Kondisi Produk Response
type KondisiProdukListResponse struct {
	ID        uuid.UUID                 `json:"id"`
	Nama      models.TranslatableString `json:"nama"`
	Slug      string                    `json:"slug"`
	Urutan    int                       `json:"urutan"`
	IsActive  bool                      `json:"is_active"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

type KondisiProdukDetailResponse struct {
	ID        uuid.UUID                 `json:"id"`
	Nama      models.TranslatableString `json:"nama"`
	Slug      string                    `json:"slug"`
	Deskripsi *string                   `json:"deskripsi,omitempty"`
	Urutan    int                       `json:"urutan"`
	IsActive  bool                      `json:"is_active"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

func ToKondisiProdukListResponse(k *models.KondisiProduk) KondisiProdukListResponse {
	return KondisiProdukListResponse{
		ID:        k.ID,
		Nama:      k.GetNama(),
		Slug:      k.Slug,
		Urutan:    k.Urutan,
		IsActive:  k.IsActive,
		UpdatedAt: k.UpdatedAt,
	}
}

func ToKondisiProdukDetailResponse(k *models.KondisiProduk) KondisiProdukDetailResponse {
	return KondisiProdukDetailResponse{
		ID:        k.ID,
		Nama:      k.GetNama(),
		Slug:      k.Slug,
		Deskripsi: k.Deskripsi,
		Urutan:    k.Urutan,
		IsActive:  k.IsActive,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

// Sumber Produk Response
type SumberProdukListResponse struct {
	ID        uuid.UUID                 `json:"id"`
	Nama      models.TranslatableString `json:"nama"`
	Slug      string                    `json:"slug"`
	IsActive  bool                      `json:"is_active"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

type SumberProdukDetailResponse struct {
	ID        uuid.UUID                 `json:"id"`
	Nama      models.TranslatableString `json:"nama"`
	Slug      string                    `json:"slug"`
	Deskripsi *string                   `json:"deskripsi,omitempty"`
	IsActive  bool                      `json:"is_active"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

func ToSumberProdukListResponse(s *models.SumberProduk) SumberProdukListResponse {
	return SumberProdukListResponse{
		ID:        s.ID,
		Nama:      s.GetNama(),
		Slug:      s.Slug,
		IsActive:  s.IsActive,
		UpdatedAt: s.UpdatedAt,
	}
}

func ToSumberProdukDetailResponse(s *models.SumberProduk) SumberProdukDetailResponse {
	return SumberProdukDetailResponse{
		ID:        s.ID,
		Nama:      s.GetNama(),
		Slug:      s.Slug,
		Deskripsi: s.Deskripsi,
		IsActive:  s.IsActive,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// Kondisi Paket Response
type KondisiPaketListResponse struct {
	ID        uuid.UUID                 `json:"id"`
	Nama      models.TranslatableString `json:"nama"`
	Slug      string                    `json:"slug"`
	Urutan    int                       `json:"urutan"`
	IsActive  bool                      `json:"is_active"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

type KondisiPaketDetailResponse struct {
	ID        uuid.UUID                 `json:"id"`
	Nama      models.TranslatableString `json:"nama"`
	Slug      string                    `json:"slug"`
	Deskripsi *string                   `json:"deskripsi,omitempty"`
	Urutan    int                       `json:"urutan"`
	IsActive  bool                      `json:"is_active"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

func ToKondisiPaketListResponse(k *models.KondisiPaket) KondisiPaketListResponse {
	return KondisiPaketListResponse{
		ID:        k.ID,
		Nama:      k.GetNama(),
		Slug:      k.Slug,
		Urutan:    k.Urutan,
		IsActive:  k.IsActive,
		UpdatedAt: k.UpdatedAt,
	}
}

func ToKondisiPaketDetailResponse(k *models.KondisiPaket) KondisiPaketDetailResponse {
	return KondisiPaketDetailResponse{
		ID:        k.ID,
		Nama:      k.GetNama(),
		Slug:      k.Slug,
		Deskripsi: k.Deskripsi,
		Urutan:    k.Urutan,
		IsActive:  k.IsActive,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}
