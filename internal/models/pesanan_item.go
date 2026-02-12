package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PesananItem struct {
	ID           uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PesananID    uuid.UUID       `gorm:"type:uuid;not null" json:"pesanan_id"`
	ProdukID     uuid.UUID       `gorm:"type:uuid;not null" json:"produk_id"`
	NamaProduk   string          `gorm:"type:varchar(200);not null" json:"nama_produk"`
	SKU          *string         `gorm:"type:varchar(50)" json:"sku"`
	Qty          int             `gorm:"not null" json:"qty"`
	HargaSatuan  decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"harga_satuan"`
	DiskonSatuan decimal.Decimal `gorm:"type:decimal(15,2);default:0" json:"diskon_satuan"`
	Subtotal     decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	CreatedAt    time.Time       `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`

	// Relations
	Pesanan Pesanan `gorm:"foreignKey:PesananID" json:"pesanan,omitempty"`
	Produk  Produk  `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
}

func (PesananItem) TableName() string {
	return "pesanan_item"
}

// Request DTOs
type CreatePesananItemRequest struct {
	ProdukID string `json:"produk_id" binding:"required,uuid"`
	Qty      int    `json:"qty" binding:"required,min=1"`
}

// Response DTOs
type PesananItemResponse struct {
	ID           string          `json:"id"`
	PesananID    string          `json:"pesanan_id"`
	ProdukID     string          `json:"produk_id"`
	NamaProduk   string          `json:"nama_produk"`
	SKU          *string         `json:"sku"`
	Qty          int             `json:"qty"`
	HargaSatuan  float64         `json:"harga_satuan"`
	DiskonSatuan float64         `json:"diskon_satuan"`
	Subtotal     float64         `json:"subtotal"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Produk       *ProdukResponse `json:"produk,omitempty"`
}
