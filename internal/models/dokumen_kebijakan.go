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
	Slug      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	SlugID    *string        `gorm:"type:varchar(100);uniqueIndex" json:"slug_id"`
	SlugEN    *string        `gorm:"type:varchar(100);uniqueIndex" json:"slug_en"`
	Konten    string         `gorm:"type:text;not null" json:"konten"`                     // HTML content (Indonesian) or JSON for FAQ
	KontenEN  string         `gorm:"column:konten_en;type:text;not null" json:"konten_en"` // HTML content (English)
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`
}

func (DokumenKebijakan) TableName() string {
	return "dokumen_kebijakan"
}

// Request DTOs
type UpdateDokumenKebijakanRequest struct {
	Judul    *string `json:"judul" binding:"omitempty,min=2,max=100"`
	JudulEN  *string `json:"judul_en" binding:"omitempty,min=2,max=100"`
	SlugID   *string `json:"slug_id" binding:"omitempty,max=100"`
	SlugEN   *string `json:"slug_en" binding:"omitempty,max=100"`
	Konten   *string `json:"konten"`    // HTML content (Indonesian)
	KontenEN *string `json:"konten_en"` // HTML content (English)
	IsActive *bool   `json:"is_active"`
}

// Response DTOs
type DokumenKebijakanListResponse struct {
	ID        string    `json:"id"`
	Judul     string    `json:"judul"`
	JudulEN   string    `json:"judul_en"`
	SlugID    *string   `json:"slug_id"`
	SlugEN    *string   `json:"slug_en"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DokumenKebijakanDetailResponse struct {
	ID        string    `json:"id"`
	Judul     string    `json:"judul"`
	JudulEN   string    `json:"judul_en"`
	SlugID    *string   `json:"slug_id"`
	SlugEN    *string   `json:"slug_en"`
	Konten    string    `json:"konten"`    // HTML content (Indonesian) or JSON for FAQ
	KontenEN  string    `json:"konten_en"` // HTML content (English)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DokumenKebijakanPublicResponse struct {
	ID     string `json:"id"`
	Judul  string `json:"judul"` // Based on lang param
	SlugID *string `json:"slug_id"`
	SlugEN *string `json:"slug_en"`
	Konten string `json:"konten"` // Based on lang param
}

// FAQ DTOs (Legacy - untuk backward compatibility dengan dokumen_kebijakan)
type FAQContentItem struct {
	Question   string `json:"question"`
	QuestionEN string `json:"question_en"`
	Answer     string `json:"answer"`
	AnswerEN   string `json:"answer_en"`
}

type FAQItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type FAQLegacyResponse struct {
	Judul string    `json:"judul"`
	Items []FAQItem `json:"items"`
}
