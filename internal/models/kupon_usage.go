package models

import (
	"time"

	"github.com/google/uuid"
)

type KuponUsage struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	KuponID       uuid.UUID `gorm:"type:uuid;not null" json:"kupon_id"`
	BuyerID       uuid.UUID `gorm:"type:uuid;not null" json:"buyer_id"`
	PesananID     uuid.UUID `gorm:"type:uuid;not null" json:"pesanan_id"`
	KodeKupon     string    `gorm:"type:varchar(50);not null" json:"kode_kupon"`
	NilaiPotongan float64   `gorm:"type:decimal(15,2);not null" json:"nilai_potongan"`
	CreatedAt     time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`

	// Relations
	Kupon   *Kupon   `gorm:"foreignKey:KuponID" json:"kupon,omitempty"`
	Buyer   *Buyer   `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Pesanan *Pesanan `gorm:"foreignKey:PesananID" json:"pesanan,omitempty"`
}

func (KuponUsage) TableName() string {
	return "kupon_usage"
}
