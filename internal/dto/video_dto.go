package dto

import (
	"project-bulky-be/internal/models"
	"time"

	"github.com/google/uuid"
)

// Video DTOs
type CreateVideoRequest struct {
	JudulID           string    `json:"judul_id" validate:"required,max=200"`
	JudulEN           string    `json:"judul_en" validate:"required,max=200"`
	Slug              string    `json:"slug" validate:"required,max=250"`
	DeskripsiID       string    `json:"deskripsi_id" validate:"required"`
	DeskripsiEN       string    `json:"deskripsi_en" validate:"required"`
	VideoURL          string    `json:"video_url" validate:"required,max=500"`
	ThumbnailURL      *string   `json:"thumbnail_url" validate:"omitempty,max=500"`
	KategoriID        uuid.UUID `json:"kategori_id" validate:"required"`
	DurasiDetik       int       `json:"durasi_detik" validate:"omitempty,gt=0"` // Optional: auto-detected from uploaded video file
	MetaTitleID       *string   `json:"meta_title_id" validate:"omitempty,max=200"`
	MetaTitleEN       *string   `json:"meta_title_en" validate:"omitempty,max=200"`
	MetaDescriptionID *string   `json:"meta_description_id"`
	MetaDescriptionEN *string   `json:"meta_description_en"`
	MetaKeywords      *string   `json:"meta_keywords"`
	IsActive          string    `json:"is_active" validate:"required,oneof=true false"`
}

type UpdateVideoRequest struct {
	JudulID           *string    `json:"judul_id" validate:"omitempty,max=200"`
	JudulEN           *string    `json:"judul_en" validate:"omitempty,max=200"`
	Slug              *string    `json:"slug" validate:"omitempty,max=250"`
	DeskripsiID       *string    `json:"deskripsi_id"`
	DeskripsiEN       *string    `json:"deskripsi_en"`
	VideoURL          *string    `json:"video_url" validate:"omitempty,max=500"`
	ThumbnailURL      *string    `json:"thumbnail_url" validate:"omitempty,max=500"`
	KategoriID        *uuid.UUID `json:"kategori_id"`
	DurasiDetik       *int       `json:"durasi_detik" validate:"omitempty,gt=0"`
	MetaTitleID       *string    `json:"meta_title_id" validate:"omitempty,max=200"`
	MetaTitleEN       *string    `json:"meta_title_en" validate:"omitempty,max=200"`
	MetaDescriptionID *string    `json:"meta_description_id"`
	MetaDescriptionEN *string    `json:"meta_description_en"`
	MetaKeywords      *string    `json:"meta_keywords"`
	IsActive          *string    `json:"is_active" validate:"omitempty,oneof=true false"`
}

type VideoResponse struct {
	ID                uuid.UUID           `json:"id"`
	JudulID           string              `json:"judul_id"`
	JudulEN           *string             `json:"judul_en"`
	Slug              string              `json:"slug"`
	DeskripsiID       string              `json:"deskripsi_id"`
	DeskripsiEN       *string             `json:"deskripsi_en"`
	VideoURL          string              `json:"video_url"`
	ThumbnailURL      *string             `json:"thumbnail_url"`
	KategoriID        uuid.UUID           `json:"kategori_id"`
	Kategori          *KategoriVideoBrief `json:"kategori,omitempty"`
	DurasiDetik       int                 `json:"durasi_detik"`
	MetaTitleID       *string             `json:"meta_title_id"`
	MetaTitleEN       *string             `json:"meta_title_en"`
	MetaDescriptionID *string             `json:"meta_description_id"`
	MetaDescriptionEN *string             `json:"meta_description_en"`
	MetaKeywords      *string             `json:"meta_keywords"`
	IsActive          bool                `json:"is_active"`
	ViewCount         int                 `json:"view_count"`
	PublishedAt       *time.Time          `json:"published_at"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

type VideoListResponse struct {
	ID           uuid.UUID           `json:"id"`
	JudulID      string              `json:"judul_id"`
	JudulEN      *string             `json:"judul_en"`
	Slug         string              `json:"slug"`
	DeskripsiID  string              `json:"deskripsi_id"`
	DeskripsiEN  *string             `json:"deskripsi_en"`
	ThumbnailURL *string             `json:"thumbnail_url"`
	Kategori     *KategoriVideoBrief `json:"kategori,omitempty"`
	DurasiDetik  int                 `json:"durasi_detik"`
	IsActive     bool                `json:"is_active"`
	ViewCount    int                 `json:"view_count"`
	PublishedAt  *time.Time          `json:"published_at"`
	CreatedAt    time.Time           `json:"created_at"`
}

type KategoriVideoBrief struct {
	ID     uuid.UUID `json:"id"`
	NamaID string    `json:"nama_id"`
	NamaEN *string   `json:"nama_en"`
	Slug   string    `json:"slug"`
}

// Video Filter Request
type VideoFilterRequest struct {
	models.PaginationRequest
	KategoriID *uuid.UUID `form:"kategori_id"`
}

func (p *VideoFilterRequest) SetDefaults() {
	p.PaginationRequest.SetDefaults()

	// Default sort by updated_at desc
	if p.SortBy == "" {
		p.SortBy = "updated_at"
	}
	if p.Order == "" {
		p.Order = "desc"
	}

	// Validate sort_by - only allow specific fields
	allowedSortBy := []string{"created_at", "updated_at", "view_count", "judul_id"}
	isValid := false
	for _, field := range allowedSortBy {
		if p.SortBy == field {
			isValid = true
			break
		}
	}
	if !isValid {
		p.SortBy = "updated_at"
	}
}

// Kategori Video DTOs
type KategoriVideoFilterRequest struct {
	models.PaginationRequest
}

func (p *KategoriVideoFilterRequest) SetDefaults() {
	p.PaginationRequest.SetDefaults()

	// Hanya sort by urutan
	p.SortBy = "urutan"
	p.Order = "asc"
}

type CreateKategoriVideoRequest struct {
	NamaID   string  `json:"nama_id" binding:"required,max=100"`
	NamaEN   *string `json:"nama_en" binding:"required,max=100"`
	Slug     string  `json:"slug" binding:"required,max=100"`
	IsActive bool    `json:"is_active"`
	Urutan   int     `json:"urutan"`
}

type UpdateKategoriVideoRequest struct {
	NamaID   *string `json:"nama_id" validate:"omitempty,max=100"`
	NamaEN   *string `json:"nama_en" validate:"omitempty,max=100"`
	Slug     *string `json:"slug" validate:"omitempty,max=100"`
	IsActive *bool   `json:"is_active"`
	Urutan   *int    `json:"urutan"`
}

// Dropdown Response DTO
type KategoriVideoDropdownResponse struct {
	ID   uuid.UUID              `json:"id"`
	Nama map[string]interface{} `json:"nama"` // {"id": "...", "en": "..."}
	Slug string                 `json:"slug"`
}
