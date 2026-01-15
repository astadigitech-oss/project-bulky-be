package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BannerEventPromo struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama         string         `gorm:"type:varchar(100);not null" json:"nama"`
	GambarURLID  string         `gorm:"column:gambar_url_id;type:varchar(500);not null" json:"-"`
	GambarURLEN  *string        `gorm:"column:gambar_url_en;type:varchar(500)" json:"-"`
	LinkURL      *string        `gorm:"column:url_tujuan;type:varchar(500)" json:"url_tujuan,omitempty"`
	Urutan       int            `gorm:"default:0" json:"urutan"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	TanggalMulai *time.Time     `json:"tanggal_mulai,omitempty"`
	TanggalAkhir *time.Time     `gorm:"column:tanggal_selesai" json:"tanggal_selesai,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (BannerEventPromo) TableName() string {
	return "banner_event_promo"
}

func (b *BannerEventPromo) GetGambarURL() TranslatableImage {
	return TranslatableImage{
		ID: b.GambarURLID,
		EN: b.GambarURLEN,
	}
}

// IsCurrentlyVisible checks if banner should be displayed based on schedule
func (b *BannerEventPromo) IsCurrentlyVisible() bool {
	if !b.IsActive {
		return false
	}

	now := time.Now()

	if b.TanggalMulai != nil && now.Before(*b.TanggalMulai) {
		return false
	}

	if b.TanggalAkhir != nil && now.After(*b.TanggalAkhir) {
		return false
	}

	return true
}
