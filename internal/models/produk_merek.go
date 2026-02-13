package models

import (
	"time"

	"github.com/google/uuid"
)

// ProdukMerek is the pivot table for many-to-many relationship between Produk and MerekProduk
type ProdukMerek struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProdukID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_produk_merek;column:produk_id" json:"produk_id"`
	MerekID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_produk_merek;column:merek_id" json:"merek_id"`
	CreatedAt time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`

	// Relations
	Produk *Produk      `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
	Merek  *MerekProduk `gorm:"foreignKey:MerekID" json:"merek,omitempty"`
}

func (ProdukMerek) TableName() string {
	return "produk_merek"
}
