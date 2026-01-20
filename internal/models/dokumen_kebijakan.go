package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DokumenKebijakan struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Judul     string         `gorm:"type:varchar(100);not null" json:"judul"`
	JudulEN   string         `gorm:"column:judul_en;type:varchar(100);not null" json:"judul_en"`
	Konten    string         `gorm:"type:text;not null" json:"konten"`                     // HTML content (Indonesian)
	KontenEN  string         `gorm:"column:konten_en;type:text;not null" json:"konten_en"` // HTML content (English)
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
	JudulEN  *string `json:"judul_en" binding:"omitempty,min=2,max=100"`
	Konten   *string `json:"konten"`    // HTML content (Indonesian)
	KontenEN *string `json:"konten_en"` // HTML content (English)
	IsActive *bool   `json:"is_active"`
}

// Response DTOs
type DokumenKebijakanListResponse struct {
	ID      string `json:"id"`
	Judul   string `json:"judul"`
	JudulEN string `json:"judul_en"`
	Urutan  int    `json:"urutan"`
	// IsActive  bool      `json:"is_active"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DokumenKebijakanDetailResponse struct {
	ID       string `json:"id"`
	Judul    string `json:"judul"`
	JudulEN  string `json:"judul_en"`
	Konten   string `json:"konten"`    // HTML content (Indonesian)
	KontenEN string `json:"konten_en"` // HTML content (English)
	Urutan   int    `json:"urutan"`
	// IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DokumenKebijakanPublicResponse struct {
	ID     string `json:"id"`
	Judul  string `json:"judul"`  // Based on lang param
	Konten string `json:"konten"` // Based on lang param
	Urutan int    `json:"urutan"`
}
