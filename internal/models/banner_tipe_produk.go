package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BannerTipeProduk struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TipeProdukID uuid.UUID      `gorm:"type:uuid;not null" json:"tipe_produk_id"`
	Nama         string         `gorm:"type:varchar(100);not null" json:"nama"`
	GambarURL    string         `gorm:"type:varchar(500);not null;column:gambar_url" json:"gambar_url"`
	Urutan       int            `gorm:"default:0" json:"urutan"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`

	// Relations
	TipeProduk TipeProduk `gorm:"foreignKey:TipeProdukID" json:"tipe_produk,omitempty"`
}

func (BannerTipeProduk) TableName() string {
	return "banner_tipe_produk"
}
