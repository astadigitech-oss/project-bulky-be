package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Kelurahan struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	KecamatanID uuid.UUID      `gorm:"type:uuid;not null" json:"kecamatan_id"`
	Nama        string         `gorm:"type:varchar(100);not null" json:"nama"`
	Kode        *string        `gorm:"type:varchar(15);uniqueIndex" json:"kode"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Kecamatan Kecamatan `gorm:"foreignKey:KecamatanID" json:"kecamatan,omitempty"`
}

func (Kelurahan) TableName() string {
	return "kelurahan"
}
