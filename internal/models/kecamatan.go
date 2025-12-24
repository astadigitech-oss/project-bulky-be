package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Kecamatan struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	KotaID    uuid.UUID      `gorm:"type:uuid;not null" json:"kota_id"`
	Nama      string         `gorm:"type:varchar(100);not null" json:"nama"`
	Kode      *string        `gorm:"type:varchar(10);uniqueIndex" json:"kode"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Kota      Kota        `gorm:"foreignKey:KotaID" json:"kota,omitempty"`
	Kelurahan []Kelurahan `gorm:"foreignKey:KecamatanID" json:"kelurahan,omitempty"`
}

func (Kecamatan) TableName() string {
	return "kecamatan"
}
