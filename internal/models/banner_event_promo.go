package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BannerEventPromo struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama           string         `gorm:"type:varchar(100);not null" json:"nama"`
	Gambar         string         `gorm:"type:varchar(255);not null" json:"gambar"`
	UrlTujuan      *string        `gorm:"type:varchar(255)" json:"url_tujuan"`
	Urutan         int            `gorm:"default:0" json:"urutan"`
	IsActive       bool           `gorm:"default:false" json:"is_active"`
	TanggalMulai   *time.Time     `json:"tanggal_mulai"`
	TanggalSelesai *time.Time     `json:"tanggal_selesai"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BannerEventPromo) TableName() string {
	return "banner_event_promo"
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

	if b.TanggalSelesai != nil && now.After(*b.TanggalSelesai) {
		return false
	}

	return true
}
