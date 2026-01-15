package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HeroSection struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama        string    `gorm:"type:varchar(100);not null" json:"nama"`
	GambarURLID string    `gorm:"column:gambar_url_id;type:varchar(500);not null" json:"-"`
	GambarURLEN *string   `gorm:"column:gambar_url_en;type:varchar(500)" json:"-"`
	// LinkURL     *string        `gorm:"type:varchar(500)" json:"link_url,omitempty"`
	Urutan    int            `gorm:"default:0" json:"urutan"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (HeroSection) TableName() string {
	return "hero_section"
}

func (h *HeroSection) GetGambarURL() TranslatableImage {
	return TranslatableImage{
		ID: h.GambarURLID,
		EN: h.GambarURLEN,
	}
}
