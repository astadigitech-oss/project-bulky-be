package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TujuanKategori represents single kategori target
// TypeScript equivalent:
//
//	interface TujuanKategori {
//	    id: string;
//	    slug: string;
//	}
type TujuanKategori struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

// TujuanList is array of TujuanKategori with JSONB support
// TypeScript equivalent: TujuanKategori[] | null
type TujuanList []TujuanKategori

// Scan implements sql.Scanner for GORM
func (t *TujuanList) Scan(value interface{}) error {
	if value == nil {
		*t = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan TujuanList: not []byte")
	}

	if len(bytes) == 0 {
		*t = nil
		return nil
	}

	return json.Unmarshal(bytes, t)
}

// Value implements driver.Valuer for GORM
func (t TujuanList) Value() (driver.Value, error) {
	if len(t) == 0 {
		return nil, nil
	}
	return json.Marshal(t)
}

type BannerEventPromo struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama           string         `gorm:"type:varchar(100);not null" json:"nama"`
	GambarURLID    string         `gorm:"column:gambar_url_id;type:varchar(500);not null" json:"-"`
	GambarURLEN    *string        `gorm:"column:gambar_url_en;type:varchar(500);not null" json:"-"`
	Tujuan         TujuanList     `gorm:"type:jsonb" json:"tujuan"`
	Urutan         int            `gorm:"default:0" json:"urutan"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	TanggalMulai   *time.Time     `json:"tanggal_mulai,omitempty"`
	TanggalSelesai *time.Time     `gorm:"column:tanggal_selesai" json:"tanggal_selesai,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (BannerEventPromo) TableName() string {
	return "banner_event_promo"
}

func (b *BannerEventPromo) GetGambarURL() TranslatableImage {
	return TranslatableImage{
		ID: b.GambarURLID,
		EN: b.GambarURLEN,
	}
}

// IsCurrentlyVisible checks if banner should be displayed based on schedule
func (b *BannerEventPromo) IsCurrentlyVisible() bool {
	if !b.IsActive {
		return false
	}

	now := time.Now()

	if b.TanggalMulai != nil && now.Before(*b.TanggalMulai) {
		return false
	}

	if b.TanggalSelesai != nil && now.After(*b.TanggalSelesai) {
		return false
	}

	return true
}
