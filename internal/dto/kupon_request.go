package dto

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// CreateKuponRequest DTO for creating kupon
type CreateKuponRequest struct {
	Kode              string      `json:"kode" validate:"required,min=3,max=50"`
	Nama              *string     `json:"nama" validate:"omitempty,max=255"`
	Deskripsi         *string     `json:"deskripsi"`
	JenisDiskon       string      `json:"jenis_diskon" validate:"required,oneof=persentase jumlah_tetap"`
	NilaiDiskon       float64     `json:"nilai_diskon" validate:"required,gt=0"`
	MinimalPembelian  float64     `json:"minimal_pembelian" validate:"required,gte=0"`
	LimitPemakaian    *int        `json:"limit_pemakaian" validate:"omitempty,gt=0"`
	TanggalKedaluarsa string      `json:"tanggal_kedaluarsa" validate:"required"`
	IsAllKategori     *bool       `json:"is_all_kategori" validate:"required"`
	KategoriIDs       []uuid.UUID `json:"kategori" validate:"omitempty,dive,uuid"`
}

// Validate performs custom validation for CreateKuponRequest
func (r *CreateKuponRequest) Validate() error {
	// Validate persentase max 100
	if r.JenisDiskon == "persentase" && r.NilaiDiskon > 100 {
		return errors.New("nilai persentase maksimal 100")
	}

	// Validate date format and not in the past
	tanggal, err := time.Parse("2006-01-02", r.TanggalKedaluarsa)
	if err != nil {
		return errors.New("format tanggal kedaluarsa tidak valid (gunakan YYYY-MM-DD)")
	}

	today := time.Now().Truncate(24 * time.Hour)
	if tanggal.Before(today) {
		return errors.New("tanggal kedaluarsa harus hari ini atau setelahnya")
	}

	// Validate is_all_kategori required
	if r.IsAllKategori == nil {
		return errors.New("is_all_kategori wajib diisi")
	}

	// Validate is_all_kategori and kategori consistency
	if *r.IsAllKategori && len(r.KategoriIDs) > 0 {
		return errors.New("jika kupon berlaku untuk semua kategori, kategori harus kosong")
	}

	// Validate kategori required if not all kategori
	if !*r.IsAllKategori && len(r.KategoriIDs) == 0 {
		return errors.New("kategori wajib dipilih jika tidak berlaku semua kategori")
	}

	return nil
}

// UpdateKuponRequest DTO for updating kupon
type UpdateKuponRequest struct {
	Kode              string      `json:"kode" validate:"required,min=3,max=50"`
	Nama              *string     `json:"nama" validate:"omitempty,max=255"`
	Deskripsi         *string     `json:"deskripsi"`
	JenisDiskon       string      `json:"jenis_diskon" validate:"required,oneof=persentase jumlah_tetap"`
	NilaiDiskon       float64     `json:"nilai_diskon" validate:"required,gt=0"`
	MinimalPembelian  float64     `json:"minimal_pembelian" validate:"required,gte=0"`
	LimitPemakaian    *int        `json:"limit_pemakaian" validate:"omitempty,gt=0"`
	TanggalKedaluarsa string      `json:"tanggal_kedaluarsa" validate:"required"`
	IsAllKategori     *bool       `json:"is_all_kategori" validate:"required"`
	KategoriIDs       []uuid.UUID `json:"kategori" validate:"omitempty,dive,uuid"`
}

// Validate performs custom validation for UpdateKuponRequest
func (r *UpdateKuponRequest) Validate() error {
	// Validate persentase max 100
	if r.JenisDiskon == "persentase" && r.NilaiDiskon > 100 {
		return errors.New("nilai persentase maksimal 100")
	}

	// Validate date format and not in the past
	tanggal, err := time.Parse("2006-01-02", r.TanggalKedaluarsa)
	if err != nil {
		return errors.New("format tanggal kedaluarsa tidak valid (gunakan YYYY-MM-DD)")
	}

	today := time.Now().Truncate(24 * time.Hour)
	if tanggal.Before(today) {
		return errors.New("tanggal kedaluarsa harus hari ini atau setelahnya")
	}

	// Validate is_all_kategori required
	if r.IsAllKategori == nil {
		return errors.New("is_all_kategori wajib diisi")
	}

	// Validate is_all_kategori and kategori consistency
	if *r.IsAllKategori && len(r.KategoriIDs) > 0 {
		return errors.New("jika kupon berlaku untuk semua kategori, kategori harus kosong")
	}

	// Validate kategori required if not all kategori
	if !*r.IsAllKategori && len(r.KategoriIDs) == 0 {
		return errors.New("kategori wajib dipilih jika tidak berlaku semua kategori")
	}

	return nil
}

// GenerateKodeRequest DTO for generating random kode kupon
type GenerateKodeRequest struct {
	Prefix string `json:"prefix" validate:"omitempty,max=20"`
	Length int    `json:"length" validate:"omitempty,min=4,max=20"`
}

// SetDefaults sets default values for GenerateKodeRequest
func (r *GenerateKodeRequest) SetDefaults() {
	if r.Length == 0 {
		r.Length = 8
	}
}

// KuponQueryParams DTO for query parameters in list kupon
type KuponQueryParams struct {
	Page        int     `form:"page" validate:"required,min=1"`
	PerPage     int     `form:"per_page" validate:"required,min=1,max=100"`
	Search      string  `form:"search"`
	JenisDiskon *string `form:"jenis_diskon" validate:"omitempty,oneof=persentase jumlah_tetap"`
	IsActive    *bool   `form:"is_active"`
	IsExpired   *bool   `form:"is_expired"`
	SortBy      string  `form:"sort_by" validate:"omitempty,oneof=tanggal_kedaluarsa updated_at"`
	Order       string  `form:"order" validate:"omitempty,oneof=asc desc"`
}

// SetDefaults sets default values for optional query params
func (q *KuponQueryParams) SetDefaults() {
	// Page and PerPage are required, no defaults
	// Set defaults for optional params
	if q.SortBy == "" {
		q.SortBy = "updated_at"
	}
	if q.Order == "" {
		q.Order = "desc"
	}
}

// KuponUsagesQueryParams DTO for query parameters in get kupon usages
type KuponUsagesQueryParams struct {
	Page    int `form:"page" validate:"required,min=1"`
	PerPage int `form:"per_page" validate:"required,min=1,max=100"`
}
