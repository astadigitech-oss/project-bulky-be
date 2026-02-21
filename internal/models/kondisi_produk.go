package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KondisiProduk struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	NamaID    string         `gorm:"column:nama_id;type:varchar(100);not null" json:"-"`
	NamaEN    *string        `gorm:"column:nama_en;type:varchar(100)" json:"-"`
	Slug      string         `gorm:"type:varchar(120);uniqueIndex;not null" json:"slug"`
	SlugID    *string        `gorm:"type:varchar(120);uniqueIndex" json:"slug_id"`
	SlugEN    *string        `gorm:"type:varchar(120);uniqueIndex" json:"slug_en"`
	Deskripsi *string        `gorm:"type:text" json:"deskripsi,omitempty"`
	Urutan    int            `gorm:"default:0" json:"urutan"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"deleted_at,omitempty"`
}

func (KondisiProduk) TableName() string {
	return "kondisi_produk"
}

func (k *KondisiProduk) GetNama() TranslatableString {
	return TranslatableString{
		ID: k.NamaID,
		EN: k.NamaEN,
	}
}
