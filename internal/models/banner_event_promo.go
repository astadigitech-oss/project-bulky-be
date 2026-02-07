package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BannerEventPromo struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama           string         `gorm:"type:varchar(100);not null" json:"nama"`
	GambarURLID    string         `gorm:"column:gambar_url_id;type:varchar(500);not null" json:"-"`
	GambarURLEN    string         `gorm:"column:gambar_url_en;type:varchar(500);not null" json:"-"`
	Tujuan         *string        `gorm:"type:varchar(1000)" json:"tujuan"` // Comma-separated kategori IDs
	Urutan         int            `gorm:"default:0" json:"urutan"`
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
		EN: &b.GambarURLEN,
	}
}

// GetTujuanIDs parses comma-separated string to slice of UUIDs
func (b *BannerEventPromo) GetTujuanIDs() []uuid.UUID {
	if b.Tujuan == nil || *b.Tujuan == "" {
		return nil
	}

	ids := strings.Split(*b.Tujuan, ",")
	result := make([]uuid.UUID, 0, len(ids))

	for _, idStr := range ids {
		idStr = strings.TrimSpace(idStr)
		if id, err := uuid.Parse(idStr); err == nil {
			result = append(result, id)
		}
	}

	return result
}

// SetTujuanFromIDs sets tujuan from slice of UUIDs
func (b *BannerEventPromo) SetTujuanFromIDs(ids []uuid.UUID) {
	if len(ids) == 0 {
		b.Tujuan = nil
		return
	}

	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = id.String()
	}

	tujuan := strings.Join(strIDs, ",")
	b.Tujuan = &tujuan
}

// IsCurrentlyVisible checks if banner should be displayed based on schedule only
func (b *BannerEventPromo) IsCurrentlyVisible() bool {
	// Jika tidak ada tanggal sama sekali, tidak tampil
	if b.TanggalMulai == nil && b.TanggalSelesai == nil {
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
