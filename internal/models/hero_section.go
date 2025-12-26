package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HeroSection struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama           string         `gorm:"type:varchar(100);not null" json:"nama"`
	Gambar         string         `gorm:"type:varchar(255);not null" json:"gambar"`
	Urutan         int            `gorm:"default:0" json:"urutan"`
	IsActive       bool           `gorm:"default:false" json:"is_active"`
	TanggalMulai   *time.Time     `json:"tanggal_mulai"`
	TanggalSelesai *time.Time     `json:"tanggal_selesai"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (HeroSection) TableName() string {
	return "hero_section"
}

// IsCurrentlyVisible checks if hero should be displayed based on schedule
func (h *HeroSection) IsCurrentlyVisible() bool {
	if !h.IsActive {
		return false
	}

	now := time.Now()

	if h.TanggalMulai != nil && now.Before(*h.TanggalMulai) {
		return false
	}

	if h.TanggalSelesai != nil && now.After(*h.TanggalSelesai) {
		return false
	}

	return true
}
