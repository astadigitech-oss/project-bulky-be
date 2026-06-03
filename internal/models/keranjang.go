package models

import (
	"time"

	"github.com/google/uuid"
)

// ========================================
// Keranjang (Cart Header)
// ========================================

type Keranjang struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	BuyerID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:keranjang_buyer_id_unique" json:"buyer_id"`
	CreatedAt time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`

	// Relations
	Buyer *Buyer          `gorm:"foreignKey:BuyerID" json:"-"`
	Items []KeranjangItem `gorm:"foreignKey:KeranjangID" json:"items,omitempty"`
}

func (Keranjang) TableName() string {
	return "keranjang"
}

// ========================================
// KeranjangItem (Cart Line Item)
// ========================================

type KeranjangItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	KeranjangID uuid.UUID `gorm:"type:uuid;not null;index" json:"keranjang_id"`
	ProdukID    uuid.UUID `gorm:"type:uuid;not null;index" json:"produk_id"`
	Quantity    int       `gorm:"not null;default:1" json:"quantity"`
	IsSelected  bool      `gorm:"not null;default:true" json:"is_selected"`
	CreatedAt   time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`

	// Relations
	Keranjang *Keranjang `gorm:"foreignKey:KeranjangID" json:"-"`
	Produk    *Produk    `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
}

func (KeranjangItem) TableName() string {
	return "keranjang_item"
}

// ========================================
// Request Types
// ========================================

type AddKeranjangItemRequest struct {
	ProdukID string `json:"produk_id" validate:"required,uuid"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

type UpdateKeranjangItemRequest struct {
	Quantity   *int  `json:"quantity" validate:"omitempty,min=1"`
	IsSelected *bool `json:"is_selected"`
}

type SelectAllKeranjangRequest struct {
	IsSelected bool `json:"is_selected"`
}

// ========================================
// Response Types
// ========================================

type ProdukKeranjangResponse struct {
	ID                 string  `json:"id"`
	NamaID             string  `json:"nama_id"`
	NamaEN             string  `json:"nama_en"`
	HargaSebelumDiskon float64 `json:"harga_sebelum_diskon"`
	HargaSesudahDiskon float64 `json:"harga_sesudah_diskon"`
	Stok               int     `json:"stok"`
	IsActive           bool    `json:"is_active"`
	GambarURL          *string `json:"gambar_url"`
}

type KeranjangItemResponse struct {
	ID         string                   `json:"id"`
	ProdukID   string                   `json:"produk_id"`
	Quantity   int                      `json:"quantity"`
	IsSelected bool                     `json:"is_selected"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
	Produk     *ProdukKeranjangResponse `json:"produk,omitempty"`
}

type KeranjangResponse struct {
	ID        string                  `json:"id"`
	BuyerID   string                  `json:"buyer_id"`
	Items     []KeranjangItemResponse `json:"items"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
}
