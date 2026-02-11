package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TipeProduk struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama      string         `gorm:"type:varchar(100);not null" json:"nama"`
	Slug      string         `gorm:"type:varchar(120);uniqueIndex;not null" json:"slug"`
	Deskripsi *string        `gorm:"type:text" json:"deskripsi"`
	Urutan    int            `gorm:"default:0" json:"urutan"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`
}

func (TipeProduk) TableName() string {
	return "tipe_produk"
}
