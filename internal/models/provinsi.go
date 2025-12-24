package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Provinsi struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama      string         `gorm:"type:varchar(100);not null" json:"nama"`
	Kode      *string        `gorm:"type:varchar(10);uniqueIndex" json:"kode"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Kota []Kota `gorm:"foreignKey:ProvinsiID" json:"kota,omitempty"`
}

func (Provinsi) TableName() string {
	return "provinsi"
}
