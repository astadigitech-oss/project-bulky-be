package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	JudulID           string         `gorm:"type:varchar(200);not null" json:"judul_id"`
	JudulEN           *string        `gorm:"type:varchar(200)" json:"judul_en"`
	Slug              string         `gorm:"type:varchar(250);uniqueIndex;not null" json:"slug"`
	DeskripsiID       string         `gorm:"type:text;not null" json:"deskripsi_id"`
	DeskripsiEN       *string        `gorm:"type:text" json:"deskripsi_en"`
	VideoURL          string         `gorm:"type:varchar(500);not null" json:"video_url"`
	ThumbnailURL      *string        `gorm:"type:varchar(500)" json:"thumbnail_url"`
	KategoriID        uuid.UUID      `gorm:"type:uuid;not null" json:"kategori_id"`
	DurasiDetik       int            `gorm:"not null" json:"durasi_detik"`
	MetaTitleID       *string        `gorm:"type:varchar(200)" json:"meta_title_id"`
	MetaTitleEN       *string        `gorm:"type:varchar(200)" json:"meta_title_en"`
	MetaDescriptionID *string        `gorm:"type:text" json:"meta_description_id"`
	MetaDescriptionEN *string        `gorm:"type:text" json:"meta_description_en"`
	MetaKeywords      *string        `gorm:"type:text" json:"meta_keywords"`
	IsActive          bool           `gorm:"default:false" json:"is_active"`
	ViewCount         int            `gorm:"default:0" json:"view_count"`
	PublishedAt       *time.Time     `json:"published_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Kategori *KategoriVideo `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
}

func (Video) TableName() string {
	return "video"
}
