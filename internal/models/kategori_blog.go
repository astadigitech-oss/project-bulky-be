package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KategoriBlog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	NamaID    string         `gorm:"type:varchar(100);not null" json:"nama_id"`
	NamaEN    *string        `gorm:"type:varchar(100)" json:"nama_en"`
	Slug      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	Urutan    int            `gorm:"default:0" json:"urutan"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (KategoriBlog) TableName() string {
	return "kategori_blog"
}
