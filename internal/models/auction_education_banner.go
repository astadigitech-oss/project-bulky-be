package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuctionEducationBanner struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama        string         `gorm:"type:varchar(100);not null" json:"nama"`
	GambarURLID string         `gorm:"type:varchar(500);not null" json:"-"`
	GambarURLEN string         `gorm:"type:varchar(500);not null" json:"-"`
	Urutan      int            `gorm:"not null;default:0" json:"urutan"`
	IsPublished bool           `gorm:"not null;default:false" json:"is_published"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AuctionEducationBanner) TableName() string { return "auction_education_banners" }

type AuctionEducationBannerResponse struct {
	ID        string            `json:"id"`
	Nama      string            `json:"nama"`
	GambarURL TranslatableImage `json:"gambar_url"`
	Urutan    int               `json:"urutan"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
