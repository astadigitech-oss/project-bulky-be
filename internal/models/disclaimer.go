package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Disclaimer struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Judul     string         `gorm:"type:varchar(200);not null" json:"judul"`
	JudulEn   string         `gorm:"type:varchar(200);not null" json:"judul_en"`
	Slug      *string        `gorm:"type:varchar(200);unique" json:"slug"`
	SlugID    *string        `gorm:"type:varchar(200);uniqueIndex" json:"slug_id"`
	SlugEN    *string        `gorm:"type:varchar(200);uniqueIndex" json:"slug_en"`
	Konten    string         `gorm:"type:text;not null" json:"konten"`
	KontenEn  string         `gorm:"type:text;not null" json:"konten_en"`
	IsActive  bool           `gorm:"default:false" json:"is_active"`
	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`
}

func (Disclaimer) TableName() string {
	return "disclaimer"
}

// Request DTOs
type CreateDisclaimerRequest struct {
	Judul    string  `json:"judul" binding:"required,min=1,max=200"`
	JudulEn  string  `json:"judul_en" binding:"required,min=1,max=200"`
	SlugID   *string `json:"slug_id" binding:"omitempty,max=200"`
	SlugEN   *string `json:"slug_en" binding:"omitempty,max=200"`
	Konten   string  `json:"konten" binding:"required"`
	KontenEn string  `json:"konten_en" binding:"required"`
	IsActive bool    `json:"is_active"`
}

type UpdateDisclaimerRequest struct {
	Judul    *string `json:"judul" binding:"omitempty,min=1,max=200"`
	JudulEn  *string `json:"judul_en" binding:"omitempty,min=1,max=200"`
	SlugID   *string `json:"slug_id" binding:"omitempty,max=200"`
	SlugEN   *string `json:"slug_en" binding:"omitempty,max=200"`
	Konten   *string `json:"konten"`
	KontenEn *string `json:"konten_en"`
	IsActive *bool   `json:"is_active"`
}

// Response DTOs
type DisclaimerListResponse struct {
	ID    string             `json:"id"`
	Judul TranslatableString `json:"judul"`
	// Slug      string             `json:"slug"`
	IsActive bool `json:"is_active"`
	// CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DisclaimerDetailResponse struct {
	ID        string             `json:"id"`
	Judul     TranslatableString `json:"judul"`
	SlugID    *string            `json:"slug_id"`
	SlugEN    *string            `json:"slug_en"`
	Konten    TranslatableString `json:"konten"`
	IsActive  bool               `json:"is_active"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type DisclaimerPublicResponse struct {
	Judul  string `json:"judul"`
	SlugID *string `json:"slug_id"`
	SlugEN *string `json:"slug_en"`
	Konten string `json:"konten"`
}
