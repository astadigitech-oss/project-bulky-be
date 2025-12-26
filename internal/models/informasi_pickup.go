package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InformasiPickup struct {
	ID             uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Alamat         string           `gorm:"type:text;not null" json:"alamat"`
	JamOperasional string           `gorm:"type:varchar(100);not null" json:"jam_operasional"`
	NomorWhatsapp  string           `gorm:"type:varchar(20);not null" json:"nomor_whatsapp"`
	Latitude       *decimal.Decimal `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude      *decimal.Decimal `gorm:"type:decimal(11,8)" json:"longitude"`
	GoogleMapsURL  *string          `gorm:"type:text" json:"google_maps_url"`
	CreatedAt      time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	JadwalGudang []JadwalGudang `gorm:"foreignKey:InformasiPickupID" json:"jadwal_gudang,omitempty"`
}

func (InformasiPickup) TableName() string {
	return "informasi_pickup"
}

// Request DTOs
type UpdateInformasiPickupRequest struct {
	Alamat         string   `json:"alamat" binding:"required"`
	JamOperasional string   `json:"jam_operasional" binding:"required,max=100"`
	NomorWhatsapp  string   `json:"nomor_whatsapp" binding:"required,max=20"`
	Latitude       *float64 `json:"latitude"`
	Longitude      *float64 `json:"longitude"`
	GoogleMapsURL  *string  `json:"google_maps_url"`
}

// Response DTOs
type InformasiPickupResponse struct {
	ID             string                 `json:"id"`
	Alamat         string                 `json:"alamat"`
	JamOperasional string                 `json:"jam_operasional"`
	NomorWhatsapp  string                 `json:"nomor_whatsapp"`
	Latitude       *float64               `json:"latitude"`
	Longitude      *float64               `json:"longitude"`
	GoogleMapsURL  *string                `json:"google_maps_url"`
	JadwalGudang   []JadwalGudangResponse `json:"jadwal_gudang"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}
