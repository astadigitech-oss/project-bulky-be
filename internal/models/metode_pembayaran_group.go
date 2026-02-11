package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetodePembayaranGroup struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama      string         `gorm:"type:varchar(50);not null;unique" json:"nama"`
	Urutan    int            `gorm:"default:0" json:"urutan"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`

	// Relations
	MetodePembayaran []MetodePembayaran `gorm:"foreignKey:GroupID" json:"metode_pembayaran,omitempty"`
}

func (MetodePembayaranGroup) TableName() string {
	return "metode_pembayaran_group"
}

// Request DTOs
type CreateMetodePembayaranGroupRequest struct {
	Nama string `json:"nama" binding:"required,min=1,max=50"`
}

type UpdateMetodePembayaranGroupRequest struct {
	Nama     *string `json:"nama" binding:"omitempty,min=1,max=50"`
	IsActive *bool   `json:"is_active"`
}

// Response DTOs
type MetodePembayaranGroupListResponse struct {
	ID           string    `json:"id"`
	Nama         string    `json:"nama"`
	Urutan       int       `json:"urutan"`
	IsActive     bool      `json:"is_active"`
	JumlahMetode int       `json:"jumlah_metode"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MetodePembayaranGroupDetailResponse struct {
	ID               string                     `json:"id"`
	Nama             string                     `json:"nama"`
	Urutan           int                        `json:"urutan"`
	IsActive         bool                       `json:"is_active"`
	JumlahMetode     int                        `json:"jumlah_metode"`
	MetodePembayaran []MetodePembayaranResponse `json:"metode_pembayaran,omitempty"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
}

type MetodePembayaranGroupResponse struct {
	ID        string    `json:"id"`
	Nama      string    `json:"nama"`
	Urutan    int       `json:"urutan"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
