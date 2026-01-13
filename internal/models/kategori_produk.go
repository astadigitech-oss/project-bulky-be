package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KategoriProduk struct {
	ID                      uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama                    string         `gorm:"type:varchar(100);not null" json:"nama"`
	Slug                    string         `gorm:"type:varchar(120);uniqueIndex;not null" json:"slug"`
	Deskripsi               *string        `gorm:"type:text" json:"deskripsi"`
	IconURL                 *string        `gorm:"type:varchar(500);column:icon_url" json:"icon_url"`
	MemilikiKondisiTambahan bool           `gorm:"default:false" json:"memiliki_kondisi_tambahan"`
	TipeKondisiTambahan     *string        `gorm:"type:varchar(20)" json:"tipe_kondisi_tambahan"`
	GambarKondisiURL        *string        `gorm:"type:varchar(500);column:gambar_kondisi_url" json:"gambar_kondisi_url"`
	TeksKondisi             *string        `gorm:"type:text" json:"teks_kondisi"`
	IsActive                bool           `gorm:"default:true" json:"is_active"`
	CreatedAt               time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"-"`
}

type KategoriProdukSimpleResponse struct {
	ID                      string    `json:"id"`
	Nama                    string    `json:"nama"`
	IconURL                 *string   `json:"icon_url"`
	IsActive                bool      `json:"is_active"`
	MemilikiKondisiTambahan bool      `json:"memiliki_kondisi_tambahan"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (KategoriProduk) TableName() string {
	return "kategori_produk"
}
