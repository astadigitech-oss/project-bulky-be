package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DokumenKebijakan struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Judul     string         `gorm:"type:varchar(100);not null" json:"judul"`
	Slug      string         `gorm:"type:varchar(120);not null;unique" json:"slug"`
	Konten    string         `gorm:"type:text;not null" json:"konten"` // HTML content from rich editor
	Urutan    int            `gorm:"default:0" json:"urutan"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DokumenKebijakan) TableName() string {
	return "dokumen_kebijakan"
}

// Request DTOs
type UpdateDokumenKebijakanRequest struct {
	Judul    *string `json:"judul" binding:"omitempty,min=2,max=100"`
	Konten   *string `json:"konten"` // HTML content
	IsActive *bool   `json:"is_active"`
}

// Response DTOs
type DokumenKebijakanListResponse struct {
	ID        string    `json:"id"`
	Judul     string    `json:"judul"`
	Slug      string    `json:"slug"`
	Urutan    int       `json:"urutan"`
	IsActive  bool      `json:"is_active"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DokumenKebijakanDetailResponse struct {
	ID        string    `json:"id"`
	Judul     string    `json:"judul"`
	Slug      string    `json:"slug"`
	Konten    string    `json:"konten"` // HTML content
	Urutan    int       `json:"urutan"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DokumenKebijakanPublicResponse struct {
	Judul  string `json:"judul"`
	Slug   string `json:"slug"`
	Konten string `json:"konten"` // HTML content
}
