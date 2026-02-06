package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Blog struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	JudulID            string         `gorm:"type:varchar(200);not null" json:"judul_id"`
	JudulEN            *string        `gorm:"type:varchar(200)" json:"judul_en"`
	Slug               string         `gorm:"type:varchar(250);uniqueIndex;not null" json:"slug"`
	KontenID           string         `gorm:"type:text;not null" json:"konten_id"`
	KontenEN           *string        `gorm:"type:text" json:"konten_en"`
	DeskripsiSingkatID string         `gorm:"type:text;not null" json:"deskripsi_singkat_id"`
	DeskripsiSingkatEN *string        `gorm:"type:text" json:"deskripsi_singkat_en"`
	FeaturedImageURL   *string        `gorm:"type:varchar(500)" json:"featured_image_url"`
	KategoriID         uuid.UUID      `gorm:"type:uuid;not null" json:"kategori_id"`
	MetaTitleID        *string        `gorm:"type:varchar(200)" json:"meta_title_id"`
	MetaTitleEN        *string        `gorm:"type:varchar(200)" json:"meta_title_en"`
	MetaDescriptionID  *string        `gorm:"type:text" json:"meta_description_id"`
	MetaDescriptionEN  *string        `gorm:"type:text" json:"meta_description_en"`
	MetaKeywords       *string        `gorm:"type:text" json:"meta_keywords"`
	IsActive           bool           `gorm:"default:false" json:"is_active"`
	ViewCount          int            `gorm:"default:0" json:"view_count"`
	PublishedAt        *time.Time     `json:"published_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Kategori *KategoriBlog `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	Labels   []LabelBlog   `gorm:"many2many:blog_label;" json:"labels,omitempty"`
}

func (Blog) TableName() string {
	return "blog"
}

type BlogLabel struct {
	BlogID  uuid.UUID `gorm:"type:uuid;primaryKey" json:"blog_id"`
	LabelID uuid.UUID `gorm:"type:uuid;primaryKey" json:"label_id"`
}

func (BlogLabel) TableName() string {
	return "blog_label"
}
