package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaintenanceType string

const (
	MaintenanceTypeBug       MaintenanceType = "BUG"
	MaintenanceTypeError     MaintenanceType = "ERROR"
	MaintenanceTypeBigUpdate MaintenanceType = "BIG_UPDATE"
	MaintenanceTypeOther     MaintenanceType = "OTHER"
)

type ModeMaintenance struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Judul           string          `gorm:"type:varchar(100);not null" json:"judul"`
	TipeMaintenance MaintenanceType `gorm:"type:maintenance_type;not null" json:"tipe_maintenance"`
	Deskripsi       string          `gorm:"type:text;not null" json:"deskripsi"`
	IsActive        bool            `gorm:"default:false" json:"is_active"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (ModeMaintenance) TableName() string {
	return "mode_maintenance"
}
