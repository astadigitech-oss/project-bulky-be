package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateType string

const (
	UpdateTypeOptional  UpdateType = "OPTIONAL"
	UpdateTypeMandatory UpdateType = "MANDATORY"
)

type ForceUpdateApp struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	KodeVersi       string         `gorm:"type:varchar(20);not null" json:"kode_versi"`
	UpdateType      UpdateType     `gorm:"type:update_type;not null" json:"update_type"`
	InformasiUpdate string         `gorm:"type:text;not null" json:"informasi_update"`
	IsActive        bool           `gorm:"default:false" json:"is_active"`
	CreatedAt       time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`
}

func (ForceUpdateApp) TableName() string {
	return "force_update_app"
}
