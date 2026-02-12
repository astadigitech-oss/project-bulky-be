package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiskonKategori struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	KategoriID       uuid.UUID      `gorm:"type:uuid;not null" json:"kategori_id"`
	PersentaseDiskon float64        `gorm:"type:decimal(5,2);not null" json:"persentase_diskon"`
	NominalDiskon    float64        `gorm:"type:decimal(15,2);default:0" json:"nominal_diskon"`
	TanggalMulai     *time.Time     `gorm:"type:date" json:"tanggal_mulai"`
	TanggalSelesai   *time.Time     `gorm:"type:date" json:"tanggal_selesai"`
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`

	// Relations
	Kategori KategoriProduk `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
}

func (DiskonKategori) TableName() string {
	return "diskon_kategori"
}
