package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Produk struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama               string         `gorm:"type:varchar(255);not null" json:"nama"`
	Slug               string         `gorm:"type:varchar(280);unique;not null" json:"slug"`
	IDCargo            *string        `gorm:"type:varchar(50);unique;column:id_cargo" json:"id_cargo"`
	KategoriID         uuid.UUID      `gorm:"type:uuid;not null" json:"kategori_id"`
	MerekID            *uuid.UUID     `gorm:"type:uuid" json:"merek_id"`
	KondisiID          uuid.UUID      `gorm:"type:uuid;not null" json:"kondisi_id"`
	KondisiPaketID     uuid.UUID      `gorm:"type:uuid;not null" json:"kondisi_paket_id"`
	SumberID           *uuid.UUID     `gorm:"type:uuid" json:"sumber_id"`
	WarehouseID        uuid.UUID      `gorm:"type:uuid;not null" json:"warehouse_id"`
	TipeProdukID       uuid.UUID      `gorm:"type:uuid;not null" json:"tipe_produk_id"`
	HargaSebelumDiskon float64        `gorm:"type:decimal(15,2);not null" json:"harga_sebelum_diskon"`
	PersentaseDiskon   float64        `gorm:"type:decimal(5,2);default:0" json:"persentase_diskon"`
	HargaSesudahDiskon float64        `gorm:"type:decimal(15,2);not null" json:"harga_sesudah_diskon"`
	Quantity           int            `gorm:"not null;default:0" json:"quantity"`
	QuantityTerjual    int            `gorm:"default:0" json:"quantity_terjual"`
	Discrepancy        *string        `gorm:"type:text" json:"discrepancy"`
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Kategori     KategoriProduk  `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	Merek        *MerekProduk    `gorm:"foreignKey:MerekID" json:"merek,omitempty"`
	Kondisi      KondisiProduk   `gorm:"foreignKey:KondisiID" json:"kondisi,omitempty"`
	KondisiPaket KondisiPaket    `gorm:"foreignKey:KondisiPaketID" json:"kondisi_paket,omitempty"`
	Sumber       *SumberProduk   `gorm:"foreignKey:SumberID" json:"sumber,omitempty"`
	Warehouse    Warehouse       `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	TipeProduk   TipeProduk      `gorm:"foreignKey:TipeProdukID" json:"tipe_produk,omitempty"`
	Gambar       []ProdukGambar  `gorm:"foreignKey:ProdukID" json:"gambar,omitempty"`
	Dokumen      []ProdukDokumen `gorm:"foreignKey:ProdukID" json:"dokumen,omitempty"`
}

func (Produk) TableName() string {
	return "produk"
}

// Response DTO
type ProdukResponse struct {
	ID                 string    `json:"id"`
	Nama               string    `json:"nama"`
	Slug               string    `json:"slug"`
	IDCargo            *string   `json:"id_cargo"`
	KategoriID         string    `json:"kategori_id"`
	MerekID            *string   `json:"merek_id"`
	KondisiID          string    `json:"kondisi_id"`
	KondisiPaketID     string    `json:"kondisi_paket_id"`
	SumberID           *string   `json:"sumber_id"`
	WarehouseID        string    `json:"warehouse_id"`
	TipeProdukID       string    `json:"tipe_produk_id"`
	HargaSebelumDiskon float64   `json:"harga_sebelum_diskon"`
	PersentaseDiskon   float64   `json:"persentase_diskon"`
	HargaSesudahDiskon float64   `json:"harga_sesudah_diskon"`
	Quantity           int       `json:"quantity"`
	QuantityTerjual    int       `json:"quantity_terjual"`
	Discrepancy        *string   `json:"discrepancy"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
