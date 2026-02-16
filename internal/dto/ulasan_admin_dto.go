package dto

import (
	"time"

	"github.com/google/uuid"
)

// UlasanAdminQueryParams query parameters for admin ulasan list
type UlasanAdminQueryParams struct {
	Page       int    `form:"page" binding:"required,min=1"`
	PerPage    int    `form:"per_page" binding:"required,min=1,max=100"`
	Cari       string `form:"cari"`
	Rating     *int   `form:"rating" binding:"omitempty,min=1,max=5"`
	IsApproved *bool  `form:"is_approved"`
	SortBy     string `form:"sort_by"`
	SortOrder  string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// SetDefaults sets default values for query params
func (p *UlasanAdminQueryParams) SetDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = 10
	}
	if p.PerPage > 100 {
		p.PerPage = 100
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.SortOrder == "" {
		p.SortOrder = "desc"
	}
}

// UlasanAdminListResponse simple response for ulasan list (admin)
type UlasanAdminListResponse struct {
	ID         uuid.UUID                    `json:"id"`
	Buyer      UlasanAdminBuyerResponse     `json:"buyer"`
	Pesanan    UlasanAdminPesananResponse   `json:"pesanan"`
	Produk     UlasanAdminProdukResponse    `json:"produk"`
	Rating     int                          `json:"rating"`
	Komentar   *string                      `json:"komentar"`
	GambarURL  *string                      `json:"gambar_url"`
	IsApproved bool                         `json:"is_approved"`
	ApprovedAt *time.Time                   `json:"approved_at"`
	ApprovedBy *UlasanAdminApproverResponse `json:"approved_by"`
	CreatedAt  time.Time                    `json:"created_at"`
}

// UlasanAdminDetailResponse detailed response for ulasan (admin)
type UlasanAdminDetailResponse struct {
	ID         uuid.UUID                        `json:"id"`
	Buyer      UlasanAdminBuyerDetailResponse   `json:"buyer"`
	Pesanan    UlasanAdminPesananDetailResponse `json:"pesanan"`
	Produk     UlasanAdminProdukDetailResponse  `json:"produk"`
	Rating     int                              `json:"rating"`
	Komentar   *string                          `json:"komentar"`
	GambarURL  *string                          `json:"gambar_url"`
	IsApproved bool                             `json:"is_approved"`
	ApprovedAt *time.Time                       `json:"approved_at"`
	ApprovedBy *UlasanAdminApproverResponse     `json:"approved_by"`
	CreatedAt  time.Time                        `json:"created_at"`
	UpdatedAt  time.Time                        `json:"updated_at"`
}

// UlasanAdminBuyerResponse buyer info for list
type UlasanAdminBuyerResponse struct {
	ID    uuid.UUID `json:"id"`
	Nama  string    `json:"nama"`
	Email string    `json:"email"`
}

// UlasanAdminBuyerDetailResponse buyer info with phone for detail
type UlasanAdminBuyerDetailResponse struct {
	ID      uuid.UUID `json:"id"`
	Nama    string    `json:"nama"`
	Email   string    `json:"email"`
	Telepon string    `json:"telepon"`
}

// UlasanAdminPesananResponse pesanan info for list
type UlasanAdminPesananResponse struct {
	ID   uuid.UUID `json:"id"`
	Kode string    `json:"kode"`
}

// UlasanAdminPesananDetailResponse pesanan info for detail
type UlasanAdminPesananDetailResponse struct {
	ID          uuid.UUID `json:"id"`
	Kode        string    `json:"kode"`
	OrderStatus string    `json:"order_status"`
	CreatedAt   time.Time `json:"created_at"`
}

// UlasanAdminProdukResponse produk info for list
type UlasanAdminProdukResponse struct {
	ID   uuid.UUID `json:"id"`
	Nama string    `json:"nama"`
	Slug string    `json:"slug"`
}

// UlasanAdminProdukDetailResponse produk info with image for detail
type UlasanAdminProdukDetailResponse struct {
	ID        uuid.UUID `json:"id"`
	Nama      string    `json:"nama"`
	Slug      string    `json:"slug"`
	GambarURL *string   `json:"gambar_url"`
}

// UlasanAdminApproverResponse approver info
type UlasanAdminApproverResponse struct {
	ID   uuid.UUID `json:"id"`
	Nama string    `json:"nama"`
}

// UlasanApproveResponse response after approve/reject action
type UlasanApproveResponse struct {
	ID         uuid.UUID  `json:"id"`
	IsApproved bool       `json:"is_approved"`
	ApprovedAt *time.Time `json:"approved_at"`
	ApprovedBy *uuid.UUID `json:"approved_by"`
}

// BulkApproveUlasanRequest request for bulk approve
type BulkApproveUlasanRequest struct {
	IDs []uuid.UUID `json:"ids" binding:"required,min=1,dive,required"`
}

// BulkApproveUlasanResponse response for bulk approve
type BulkApproveUlasanResponse struct {
	ApprovedCount int         `json:"approved_count"`
	ApprovedIDs   []uuid.UUID `json:"approved_ids"`
}
