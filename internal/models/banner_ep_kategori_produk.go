package models

import (
	"time"

	"github.com/google/uuid"
)

type BannerEPKategoriProduk struct {
	BannerID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"banner_id"`
	KategoriProdukID uuid.UUID `gorm:"type:uuid;primaryKey" json:"kategori_produk_id"`
	CreatedAt        time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`

	// Relations (for preload)
	KategoriProduk *KategoriProduk `gorm:"foreignKey:KategoriProdukID" json:"kategori_produk,omitempty"`
}

func (BannerEPKategoriProduk) TableName() string {
	return "banner_ep_kategori_produk"
}
