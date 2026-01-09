package dto

import (
	"time"

	"github.com/google/uuid"
)

// TipeProdukListDTO represents simplified response for GET All (List)
// Used for table display with reduced payload (8 fields)
type TipeProdukListDTO struct {
	ID   uuid.UUID `json:"id"`
	Nama string    `json:"nama"`
	// Slug         string    `json:"slug"`
	// IconURL *string `json:"icon_url"`
	Urutan int `json:"urutan"`
	// JumlahProduk int       `json:"jumlah_produk"`
	// IsActive     bool      `json:"is_active"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TipeProdukDetailDTO represents complete response for GET by ID (Detail)
// Includes all fields for detail view (10 fields)
type TipeProdukDetailDTO struct {
	ID        uuid.UUID `json:"id"`
	Nama      string    `json:"nama"`
	Slug      string    `json:"slug"`
	Deskripsi *string   `json:"deskripsi"`
	// IconURL   *string   `json:"icon_url"`
	Urutan int `json:"urutan"`
	// IsActive     bool      `json:"is_active"`
	// JumlahProduk int       `json:"jumlah_produk"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
