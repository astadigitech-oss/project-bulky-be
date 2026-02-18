package dto

import (
	"project-bulky-be/internal/models"
	"time"

	"github.com/google/uuid"
)

// KuponListResponse DTO for kupon list item
type KuponListResponse struct {
	ID          uuid.UUID `json:"id"`
	Status      bool      `json:"status"`
	JenisDiskon string    `json:"jenis_diskon"`
	Kode        string    `json:"kode"`
	Nama        *string   `json:"nama"`
	NilaiDiskon float64   `json:"nilai_diskon"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// KuponDetailResponse DTO for kupon detail
type KuponDetailResponse struct {
	ID                uuid.UUID               `json:"id"`
	Kode              string                  `json:"kode"`
	Nama              *string                 `json:"nama"`
	Deskripsi         *string                 `json:"deskripsi"`
	JenisDiskon       string                  `json:"jenis_diskon"`
	NilaiDiskon       float64                 `json:"nilai_diskon"`
	MinimalPembelian  float64                 `json:"minimal_pembelian"`
	LimitPemakaian    *int                    `json:"limit_pemakaian"`
	TotalUsage        int                     `json:"total_usage"`
	RemainingUsage    *int                    `json:"remaining_usage"`
	TanggalKedaluarsa time.Time               `json:"tanggal_kedaluarsa"`
	IsAllKategori     bool                    `json:"is_all_kategori"`
	Kategori          []KuponKategoriResponse `json:"kategori"`
	IsActive          bool                    `json:"is_active"`
	IsExpired         bool                    `json:"is_expired"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
}

// KuponKategoriResponse DTO for kategori in kupon
type KuponKategoriResponse struct {
	ID   uuid.UUID                 `json:"id"`
	Nama models.TranslatableString `json:"nama"`
	Slug string                    `json:"slug"`
}

// KuponUsageItemResponse DTO for kupon usage item
type KuponUsageItemResponse struct {
	ID            uuid.UUID             `json:"id"`
	Buyer         KuponUsageBuyerInfo   `json:"buyer"`
	Pesanan       KuponUsagePesananInfo `json:"pesanan"`
	NilaiPotongan float64               `json:"nilai_potongan"`
	CreatedAt     time.Time             `json:"created_at"`
}

// KuponUsageBuyerInfo DTO for buyer info in usage
type KuponUsageBuyerInfo struct {
	ID    uuid.UUID `json:"id"`
	Nama  string    `json:"nama"`
	Email string    `json:"email"`
}

// KuponUsagePesananInfo DTO for pesanan info in usage
type KuponUsagePesananInfo struct {
	ID   uuid.UUID `json:"id"`
	Kode string    `json:"kode"`
}

// KuponUsageListResponse DTO for list of kupon usages
type KuponUsageListResponse struct {
	Kupon  KuponUsageSummary        `json:"kupon"`
	Usages []KuponUsageItemResponse `json:"usages"`
}

// KuponUsageSummary DTO for kupon summary in usage list
type KuponUsageSummary struct {
	ID         uuid.UUID `json:"id"`
	Kode       string    `json:"kode"`
	TotalUsage int       `json:"total_usage"`
}

// GeneratedKodeResponse DTO for generated kode response
type GeneratedKodeResponse struct {
	Kode string `json:"kode"`
}

// KategoriDropdownResponse DTO for kategori dropdown
type KategoriDropdownResponse struct {
	ID   uuid.UUID                 `json:"id"`
	Nama models.TranslatableString `json:"nama"`
}

// ToggleStatusResponse DTO for toggle status response
type ToggleStatusResponse struct {
	ID        uuid.UUID `json:"id"`
	Kode      string    `json:"kode"`
	IsActive  bool      `json:"is_active"`
	UpdatedAt time.Time `json:"updated_at"`
}
