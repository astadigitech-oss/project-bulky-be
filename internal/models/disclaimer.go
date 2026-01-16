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
	Konten    string         `gorm:"type:text;not null" json:"konten"`
	KontenEn  string         `gorm:"type:text;not null" json:"konten_en"`
	IsActive  bool           `gorm:"default:false" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Disclaimer) TableName() string {
	return "disclaimer"
}

// Request DTOs
type CreateDisclaimerRequest struct {
	Judul    string  `json:"judul" binding:"required,min=1,max=200"`
	JudulEn  string  `json:"judul_en" binding:"required,min=1,max=200"`
	Slug     *string `json:"slug" binding:"omitempty,max=200"`
	Konten   string  `json:"konten" binding:"required"`
	KontenEn string  `json:"konten_en" binding:"required"`
	IsActive bool    `json:"is_active"`
}

type UpdateDisclaimerRequest struct {
	Judul    *string `json:"judul" binding:"omitempty,min=1,max=200"`
	JudulEn  *string `json:"judul_en" binding:"omitempty,min=1,max=200"`
	Slug     *string `json:"slug" binding:"omitempty,max=200"`
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
	Slug      string             `json:"slug"`
	Konten    TranslatableString `json:"konten"`
	IsActive  bool               `json:"is_active"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type DisclaimerPublicResponse struct {
	Judul  string `json:"judul"`
	Slug   string `json:"slug"`
	Konten string `json:"konten"`
}
