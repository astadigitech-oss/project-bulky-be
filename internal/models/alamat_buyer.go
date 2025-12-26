package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlamatBuyer struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BuyerID         uuid.UUID      `gorm:"type:uuid;not null" json:"buyer_id"`

	// Identitas Penerima
	Label           string `gorm:"type:varchar(50);not null" json:"label"`
	NamaPenerima    string `gorm:"type:varchar(100);not null" json:"nama_penerima"`
	TeleponPenerima string `gorm:"type:varchar(20);not null" json:"telepon_penerima"`

	// Wilayah (dari Google Maps)
	Provinsi  string  `gorm:"type:varchar(100);not null" json:"provinsi"`
	Kota      string  `gorm:"type:varchar(100);not null" json:"kota"`
	Kecamatan *string `gorm:"type:varchar(100)" json:"kecamatan"`
	Kelurahan *string `gorm:"type:varchar(100)" json:"kelurahan"`
	KodePos   *string `gorm:"type:varchar(10)" json:"kode_pos"`

	// Alamat Detail
	AlamatLengkap string  `gorm:"type:text;not null" json:"alamat_lengkap"`
	Catatan       *string `gorm:"type:text" json:"catatan"`

	// Koordinat
	Latitude  *float64 `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude *float64 `gorm:"type:decimal(11,8)" json:"longitude"`

	// Google Reference
	GooglePlaceID *string `gorm:"type:varchar(255)" json:"google_place_id"`

	// Status
	IsDefault bool `gorm:"default:false" json:"is_default"`

	// Timestamps
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Buyer Buyer `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
}

func (AlamatBuyer) TableName() string {
	return "alamat_buyer"
}
