package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ulasan struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PesananID     uuid.UUID `gorm:"type:uuid;not null" json:"pesanan_id"`
	PesananItemID uuid.UUID `gorm:"type:uuid;not null;unique" json:"pesanan_item_id"`
	BuyerID       uuid.UUID `gorm:"type:uuid;not null" json:"buyer_id"`
	ProdukID      uuid.UUID `gorm:"type:uuid;not null" json:"produk_id"`

	Rating   int     `gorm:"not null" json:"rating"`
	Komentar *string `gorm:"type:text" json:"komentar"`
	Gambar   *string `gorm:"type:varchar(255)" json:"gambar"`

	IsApproved bool       `gorm:"default:false" json:"is_approved"`
	ApprovedAt *time.Time `gorm:"type:timestamptz" json:"approved_at"`
	ApprovedBy *uuid.UUID `gorm:"type:uuid" json:"approved_by"`

	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`

	// Relations
	Pesanan     Pesanan     `gorm:"foreignKey:PesananID" json:"pesanan,omitempty"`
	PesananItem PesananItem `gorm:"foreignKey:PesananItemID" json:"pesanan_item,omitempty"`
	Buyer       Buyer       `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Produk      Produk      `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
	Approver    *Admin      `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

func (Ulasan) TableName() string {
	return "ulasan"
}
