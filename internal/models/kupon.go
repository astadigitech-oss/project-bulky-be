package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JenisDiskon string

const (
	JenisDiskonPersentase  JenisDiskon = "persentase"
	JenisDiskonJumlahTetap JenisDiskon = "jumlah_tetap"
)

type Kupon struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Kode              string         `gorm:"type:varchar(50);not null" json:"kode"`
	Nama              *string        `gorm:"type:varchar(255)" json:"nama"`
	Deskripsi         *string        `gorm:"type:text" json:"deskripsi"`
	JenisDiskon       JenisDiskon    `gorm:"type:varchar(20);not null" json:"jenis_diskon"`
	NilaiDiskon       float64        `gorm:"type:decimal(15,2);not null" json:"nilai_diskon"`
	MinimalPembelian  float64        `gorm:"type:decimal(15,2);not null;default:0" json:"minimal_pembelian"`
	LimitPemakaian    *int           `gorm:"type:integer" json:"limit_pemakaian"`
	TanggalKedaluarsa time.Time      `gorm:"type:date;not null" json:"tanggal_kedaluarsa"`
	IsAllKategori     bool           `gorm:"not null;default:true" json:"is_all_kategori"`
	IsActive          bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt         time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"type:timestamptz;index" json:"deleted_at,omitempty"`

	// Relations
	Kategori []KuponKategori `gorm:"foreignKey:KuponID" json:"kategori,omitempty"`
	Usages   []KuponUsage    `gorm:"foreignKey:KuponID" json:"usages,omitempty"`
}

func (Kupon) TableName() string {
	return "kupon"
}

// GetUsageCount returns current usage count
func (k *Kupon) GetUsageCount(db *gorm.DB) int64 {
	var count int64
	db.Model(&KuponUsage{}).Where("kupon_id = ?", k.ID).Count(&count)
	return count
}

// IsExpired checks if kupon is expired
func (k *Kupon) IsExpired() bool {
	now := time.Now()
	// Set to end of day for comparison
	endOfDay := time.Date(
		k.TanggalKedaluarsa.Year(),
		k.TanggalKedaluarsa.Month(),
		k.TanggalKedaluarsa.Day(),
		23, 59, 59, 999999999,
		time.Local,
	)
	return now.After(endOfDay)
}

// IsLimitReached checks if usage limit is reached
func (k *Kupon) IsLimitReached(db *gorm.DB) bool {
	if k.LimitPemakaian == nil {
		return false
	}
	return k.GetUsageCount(db) >= int64(*k.LimitPemakaian)
}

// GetRemainingUsage returns remaining usage count
func (k *Kupon) GetRemainingUsage(db *gorm.DB) *int {
	if k.LimitPemakaian == nil {
		return nil
	}
	remaining := *k.LimitPemakaian - int(k.GetUsageCount(db))
	return &remaining
}
