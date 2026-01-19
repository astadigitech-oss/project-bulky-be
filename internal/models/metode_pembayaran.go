package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetodePembayaran struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	GroupID   uuid.UUID      `gorm:"type:uuid;not null" json:"group_id"`
	Nama      string         `gorm:"type:varchar(50);not null" json:"nama"`
	Kode      string         `gorm:"type:varchar(30);not null;unique" json:"kode"`
	LogoValue *string        `gorm:"column:logo_value;type:varchar(50)" json:"logo_value"`
	Urutan    int            `gorm:"default:0" json:"urutan"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Group MetodePembayaranGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}

func (MetodePembayaran) TableName() string {
	return "metode_pembayaran"
}

// Request DTOs
type UpdateMetodePembayaranRequest struct {
	GroupID   *string `json:"group_id" binding:"omitempty,uuid"`
	Nama      *string `json:"nama" binding:"omitempty,min=1,max=50"`
	Kode      *string `json:"kode" binding:"omitempty,min=1,max=30"`
	LogoValue *string `json:"logo_value" binding:"omitempty,max=50"`
	Urutan    *int    `json:"urutan"`
	IsActive  *bool   `json:"is_active"`
}

// Response DTOs
type MetodePembayaranGroupSimple struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

type MetodePembayaranListResponse struct {
	ID        string                      `json:"id"`
	Nama      string                      `json:"nama"`
	Kode      string                      `json:"kode"`
	LogoValue *string                     `json:"logo_value"`
	Urutan    int                         `json:"urutan"`
	IsActive  bool                        `json:"is_active"`
	Group     MetodePembayaranGroupSimple `json:"group"`
	UpdatedAt time.Time                   `json:"updated_at"`
}

type MetodePembayaranDetailResponse struct {
	ID        string                      `json:"id"`
	Nama      string                      `json:"nama"`
	Kode      string                      `json:"kode"`
	LogoValue *string                     `json:"logo_value"`
	Urutan    int                         `json:"urutan"`
	IsActive  bool                        `json:"is_active"`
	Group     MetodePembayaranGroupSimple `json:"group"`
	CreatedAt time.Time                   `json:"created_at"`
	UpdatedAt time.Time                   `json:"updated_at"`
}

type MetodePembayaranResponse struct {
	ID        string  `json:"id"`
	Nama      string  `json:"nama"`
	Kode      string  `json:"kode"`
	LogoValue *string `json:"logo_value"`
	Urutan    int     `json:"urutan"`
	IsActive  bool    `json:"is_active"`
}
