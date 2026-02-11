package models

import (
	"time"

	"github.com/google/uuid"
)

type ProdukGambar struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProdukID  uuid.UUID `gorm:"type:uuid;not null" json:"produk_id"`
	GambarURL string    `gorm:"type:varchar(500);not null;column:gambar_url" json:"gambar_url"`
	Urutan    int       `gorm:"default:0" json:"urutan"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
}

func (ProdukGambar) TableName() string {
	return "produk_gambar"
}
