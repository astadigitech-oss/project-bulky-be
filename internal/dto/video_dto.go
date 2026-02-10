package dto

import (
	"time"

	"github.com/google/uuid"
)

// Video DTOs
type CreateVideoRequest struct {
	JudulID           string    `json:"judul_id" validate:"required,max=200"`
	JudulEN           *string   `json:"judul_en" validate:"omitempty,max=200"`
	Slug              string    `json:"slug" validate:"required,max=250"`
	DeskripsiID       string    `json:"deskripsi_id" validate:"required"`
	DeskripsiEN       *string   `json:"deskripsi_en"`
	VideoURL          string    `json:"video_url" validate:"required,max=500"`
	ThumbnailURL      *string   `json:"thumbnail_url" validate:"omitempty,max=500"`
	KategoriID        uuid.UUID `json:"kategori_id" validate:"required"`
	DurasiDetik       int       `json:"durasi_detik" validate:"required,gt=0"`
	MetaTitleID       *string   `json:"meta_title_id" validate:"omitempty,max=200"`
	MetaTitleEN       *string   `json:"meta_title_en" validate:"omitempty,max=200"`
	MetaDescriptionID *string   `json:"meta_description_id"`
	MetaDescriptionEN *string   `json:"meta_description_en"`
	MetaKeywords      *string   `json:"meta_keywords"`
	IsActive          bool      `json:"is_active"`
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
	IsActive          *bool      `json:"is_active"`
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

// Kategori Video DTOs
type CreateKategoriVideoRequest struct {
	NamaID   string  `json:"nama_id" validate:"required,max=100"`
	NamaEN   *string `json:"nama_en" validate:"omitempty,max=100"`
	Slug     string  `json:"slug" validate:"required,max=100"`
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
