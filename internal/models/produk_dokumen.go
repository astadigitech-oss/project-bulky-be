package models

import (
	"time"

	"github.com/google/uuid"
)

type ProdukDokumen struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProdukID    uuid.UUID `gorm:"type:uuid;not null" json:"produk_id"`
	NamaDokumen string    `gorm:"type:varchar(255);not null" json:"nama_dokumen"`
	FileURL     string    `gorm:"type:varchar(500);not null;column:file_url" json:"file_url"`
	TipeFile    string    `gorm:"type:varchar(50);not null" json:"tipe_file"`
	UkuranFile  *int      `gorm:"type:integer" json:"ukuran_file"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (ProdukDokumen) TableName() string {
	return "produk_dokumen"
}
