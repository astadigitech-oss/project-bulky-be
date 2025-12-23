package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SumberProduk struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama      string         `gorm:"type:varchar(100);not null" json:"nama"`
	Slug      string         `gorm:"type:varchar(120);uniqueIndex;not null" json:"slug"`
	Deskripsi *string        `gorm:"type:text" json:"deskripsi"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SumberProduk) TableName() string {
	return "sumber_produk"
}
