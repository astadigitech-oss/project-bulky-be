package models

import (
	"time"

	"github.com/google/uuid"
)

type KuponKategori struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	KuponID    uuid.UUID `gorm:"type:uuid;not null" json:"kupon_id"`
	KategoriID uuid.UUID `gorm:"type:uuid;not null" json:"kategori_id"`
	CreatedAt  time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`

	// Relations
	Kupon    *Kupon          `gorm:"foreignKey:KuponID" json:"kupon,omitempty"`
	Kategori *KategoriProduk `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
}

func (KuponKategori) TableName() string {
	return "kupon_kategori"
}
