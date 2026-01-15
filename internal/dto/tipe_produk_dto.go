package dto

import (
	"github.com/google/uuid"
)

// TipeProdukListDTO represents simplified response for GET All (List)
// Used for table display without pagination
type TipeProdukListDTO struct {
	ID   uuid.UUID `json:"id"`
	Nama string    `json:"nama"`
	// Slug      string    `json:"slug"`
	Urutan int `json:"urutan"`
	// IsActive  bool      `json:"is_active"`
	// UpdatedAt time.Time `json:"updated_at"`
}

// TipeProdukWithProdukDTO represents tipe produk with array of products
type TipeProdukWithProdukDTO struct {
	ID   uuid.UUID `json:"id"`
	Nama string    `json:"nama"`
	// Slug      string           `json:"slug"`
	// Deskripsi *string          `json:"deskripsi"`
	Urutan int `json:"urutan"`
	// IsActive bool             `json:"is_active"`
	Produk []ProdukBasicDTO `json:"produk"`
}

// ProdukBasicDTO represents basic product information nested in tipe produk
type ProdukBasicDTO struct {
	ID                 uuid.UUID `json:"id"`
	Nama               string    `json:"nama"`
	Slug               string    `json:"slug"`
	HargaSebelumDiskon float64   `json:"harga_sebelum_diskon"`
	PersentaseDiskon   float64   `json:"persentase_diskon"`
	HargaSesudahDiskon float64   `json:"harga_sesudah_diskon"`
	Quantity           int       `json:"quantity"`
	IsActive           bool      `json:"is_active"`
}
