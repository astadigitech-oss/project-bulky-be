package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlamatBuyer struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BuyerID         uuid.UUID      `gorm:"type:uuid;not null" json:"buyer_id"`
	KelurahanID     uuid.UUID      `gorm:"type:uuid;not null" json:"kelurahan_id"`
	Label           string         `gorm:"type:varchar(50);not null" json:"label"`
	NamaPenerima    string         `gorm:"type:varchar(100);not null" json:"nama_penerima"`
	TeleponPenerima string         `gorm:"type:varchar(20);not null" json:"telepon_penerima"`
	KodePos         string         `gorm:"type:varchar(10);not null" json:"kode_pos"`
	AlamatLengkap   string         `gorm:"type:text;not null" json:"alamat_lengkap"`
	Catatan         *string        `gorm:"type:text" json:"catatan"`
	IsDefault       bool           `gorm:"default:false" json:"is_default"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Buyer     Buyer     `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Kelurahan Kelurahan `gorm:"foreignKey:KelurahanID" json:"kelurahan,omitempty"`
}

func (AlamatBuyer) TableName() string {
	return "alamat_buyer"
}
