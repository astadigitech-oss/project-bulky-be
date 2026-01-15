package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KategoriProduk struct {
	ID                      uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	NamaID                  string         `gorm:"column:nama_id;type:varchar(100);not null" json:"-"`
	NamaEN                  *string        `gorm:"column:nama_en;type:varchar(100)" json:"-"`
	Slug                    string         `gorm:"type:varchar(120);uniqueIndex;not null" json:"slug"`
	Deskripsi               *string        `gorm:"type:text" json:"deskripsi,omitempty"`
	IconURL                 *string        `gorm:"type:varchar(500);column:icon_url" json:"icon_url,omitempty"`
	MemilikiKondisiTambahan bool           `gorm:"default:false" json:"memiliki_kondisi_tambahan"`
	TipeKondisiTambahan     *string        `gorm:"type:varchar(20)" json:"tipe_kondisi_tambahan,omitempty"`
	GambarKondisiURL        *string        `gorm:"type:varchar(500);column:gambar_kondisi_url" json:"gambar_kondisi_url,omitempty"`
	TeksKondisi             *string        `gorm:"type:text" json:"teks_kondisi,omitempty"`
	IsActive                bool           `gorm:"default:true" json:"is_active"`
	CreatedAt               time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (KategoriProduk) TableName() string {
	return "kategori_produk"
}

func (k *KategoriProduk) GetNama() TranslatableString {
	return TranslatableString{
		ID: k.NamaID,
		EN: k.NamaEN,
	}
}
