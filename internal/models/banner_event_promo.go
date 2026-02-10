package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BannerEventPromo struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama           string         `gorm:"type:varchar(100);not null" json:"nama"`
	GambarURLID    string         `gorm:"column:gambar_url_id;type:varchar(500);not null" json:"-"`
	GambarURLEN    string         `gorm:"column:gambar_url_en;type:varchar(500);not null" json:"-"`
	Urutan         int            `gorm:"default:0" json:"urutan"`
	TanggalMulai   *time.Time     `json:"tanggal_mulai"`
	TanggalSelesai *time.Time     `gorm:"column:tanggal_selesai" json:"tanggal_selesai"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	KategoriList []BannerEPKategoriProduk `gorm:"foreignKey:BannerID" json:"kategori_list,omitempty"`
}

func (BannerEventPromo) TableName() string {
	return "banner_event_promo"
}

func (b *BannerEventPromo) GetGambarURL() TranslatableImage {
	return TranslatableImage{
		ID: b.GambarURLID,
		EN: &b.GambarURLEN,
	}
}

// GetKategoriIDs returns slice of kategori IDs
func (b *BannerEventPromo) GetKategoriIDs() []uuid.UUID {
	ids := make([]uuid.UUID, len(b.KategoriList))
	for i, k := range b.KategoriList {
		ids[i] = k.KategoriProdukID
	}
	return ids
}

// GetKategoriIDStrings returns slice of kategori ID strings
func (b *BannerEventPromo) GetKategoriIDStrings() []string {
	ids := make([]string, len(b.KategoriList))
	for i, k := range b.KategoriList {
		ids[i] = k.KategoriProdukID.String()
	}
	return ids
}

// IsCurrentlyVisible checks if banner should be displayed based on schedule only
func (b *BannerEventPromo) IsCurrentlyVisible() bool {
	// Jika tidak ada tanggal sama sekali, tidak tampil
	if b.TanggalMulai == nil && b.TanggalSelesai == nil {
		return false
	}

	now := time.Now()

	// Not visible if NOW is before start date
	if b.TanggalMulai != nil && now.Before(*b.TanggalMulai) {
		return false
	}

	// Not visible if NOW is >= end date (includes exact match)
	if b.TanggalSelesai != nil && !now.Before(*b.TanggalSelesai) {
		return false
	}

	return true
}
