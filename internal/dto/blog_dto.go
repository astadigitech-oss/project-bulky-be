package dto

import (
	"project-bulky-be/internal/models"
	"time"

	"github.com/google/uuid"
)

// Blog DTOs
type CreateBlogRequest struct {
	JudulID           string      `json:"judul_id" validate:"required,max=200"`
	JudulEN           *string     `json:"judul_en" validate:"omitempty,max=200"`
	Slug              string      `json:"slug" validate:"required,max=250"`
	KontenID          string      `json:"konten_id" validate:"required"`
	KontenEN          *string     `json:"konten_en"`
	FeaturedImageURL  *string     `json:"featured_image_url" validate:"omitempty,max=500"`
	KategoriID        uuid.UUID   `json:"kategori_id" validate:"required"`
	MetaTitleID       *string     `json:"meta_title_id" validate:"omitempty,max=200"`
	MetaTitleEN       *string     `json:"meta_title_en" validate:"omitempty,max=200"`
	MetaDescriptionID *string     `json:"meta_description_id"`
	MetaDescriptionEN *string     `json:"meta_description_en"`
	MetaKeywords      *string     `json:"meta_keywords"`
	IsActive          bool        `json:"is_active"`
	LabelIDs          []uuid.UUID `json:"label_ids"`
}

type UpdateBlogRequest struct {
	JudulID           *string     `json:"judul_id" validate:"omitempty,max=200"`
	JudulEN           *string     `json:"judul_en" validate:"omitempty,max=200"`
	Slug              *string     `json:"slug" validate:"omitempty,max=250"`
	KontenID          *string     `json:"konten_id"`
	KontenEN          *string     `json:"konten_en"`
	FeaturedImageURL  *string     `json:"featured_image_url" validate:"omitempty,max=500"`
	KategoriID        *uuid.UUID  `json:"kategori_id"`
	MetaTitleID       *string     `json:"meta_title_id" validate:"omitempty,max=200"`
	MetaTitleEN       *string     `json:"meta_title_en" validate:"omitempty,max=200"`
	MetaDescriptionID *string     `json:"meta_description_id"`
	MetaDescriptionEN *string     `json:"meta_description_en"`
	MetaKeywords      *string     `json:"meta_keywords"`
	IsActive          *bool       `json:"is_active"`
	LabelIDs          []uuid.UUID `json:"label_ids"`
}

type BlogResponse struct {
	ID                uuid.UUID          `json:"id"`
	JudulID           string             `json:"judul_id"`
	JudulEN           *string            `json:"judul_en"`
	Slug              string             `json:"slug"`
	KontenID          string             `json:"konten_id"`
	KontenEN          *string            `json:"konten_en"`
	FeaturedImageURL  *string            `json:"featured_image_url"`
	KategoriID        uuid.UUID          `json:"kategori_id"`
	Kategori          *KategoriBlogBrief `json:"kategori,omitempty"`
	MetaTitleID       *string            `json:"meta_title_id"`
	MetaTitleEN       *string            `json:"meta_title_en"`
	MetaDescriptionID *string            `json:"meta_description_id"`
	MetaDescriptionEN *string            `json:"meta_description_en"`
	MetaKeywords      *string            `json:"meta_keywords"`
	IsActive          bool               `json:"is_active"`
	ViewCount         int                `json:"view_count"`
	PublishedAt       *time.Time         `json:"published_at"`
	Labels            []LabelBlogBrief   `json:"labels,omitempty"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

type BlogListResponse struct {
	ID               uuid.UUID          `json:"id"`
	JudulID          string             `json:"judul_id"`
	JudulEN          *string            `json:"judul_en"`
	Slug             string             `json:"slug"`
	FeaturedImageURL *string            `json:"featured_image_url"`
	Kategori         *KategoriBlogBrief `json:"kategori,omitempty"`
	IsActive         bool               `json:"is_active"`
	ViewCount        int                `json:"view_count"`
	PublishedAt      *time.Time         `json:"published_at"`
	Labels           []LabelBlogBrief   `json:"labels,omitempty"`
	CreatedAt        time.Time          `json:"created_at"`
}

type KategoriBlogBrief struct {
	ID     uuid.UUID `json:"id"`
	NamaID string    `json:"nama_id"`
	NamaEN *string   `json:"nama_en"`
	Slug   string    `json:"slug"`
}

type LabelBlogBrief struct {
	ID     uuid.UUID `json:"id"`
	NamaID string    `json:"nama_id"`
	NamaEN string    `json:"nama_en"`
	Slug   string    `json:"slug"`
}

// Blog Filter Request
type BlogFilterRequest struct {
	models.PaginationRequest
	KategoriID *uuid.UUID `form:"kategori_id"`
	LabelID    *uuid.UUID `form:"label_id"`
}

func (p *BlogFilterRequest) SetDefaults() {
	p.PaginationRequest.SetDefaults()

	// Default sort by updated_at desc
	if p.SortBy == "" {
		p.SortBy = "updated_at"
	}
	if p.Order == "" {
		p.Order = "desc"
	}

	// Validate sort_by - only allow is_active and updated_at
	if p.SortBy != "is_active" && p.SortBy != "updated_at" {
		p.SortBy = "updated_at"
	}
}

// Kategori Blog DTOs
type KategoriBlogFilterRequest struct {
	models.PaginationRequest
}

func (p *KategoriBlogFilterRequest) SetDefaults() {
	p.PaginationRequest.SetDefaults()

	// Hanya sort by urutan
	p.SortBy = "urutan"
	p.Order = "asc"
}

// Label Blog DTOs
type LabelBlogFilterRequest struct {
	models.PaginationRequest
}

func (p *LabelBlogFilterRequest) SetDefaults() {
	p.PaginationRequest.SetDefaults()

	// Hanya sort by urutan
	p.SortBy = "urutan"
	p.Order = "asc"
}

type CreateKategoriBlogRequest struct {
	NamaID   string  `json:"nama_id" validate:"required,max=100"`
	NamaEN   *string `json:"nama_en" validate:"required,max=100"`
	Slug     string  `json:"slug" validate:"required,max=100"`
	IsActive bool    `json:"is_active"`
	Urutan   int     `json:"urutan"`
}

type UpdateKategoriBlogRequest struct {
	NamaID   *string `json:"nama_id" validate:"omitempty,max=100"`
	NamaEN   *string `json:"nama_en" validate:"omitempty,max=100"`
	Slug     *string `json:"slug" validate:"omitempty,max=100"`
	IsActive *bool   `json:"is_active"`
	Urutan   *int    `json:"urutan"`
}

// Label Blog DTOs
type CreateLabelBlogRequest struct {
	NamaID string  `json:"nama_id" binding:"required,max=100" validate:"required,max=100"`
	NamaEN string  `json:"nama_en" binding:"required,max=100" validate:"required,max=100"`
	Slug   *string `json:"slug" binding:"omitempty,max=100" validate:"omitempty,max=100"`
	Urutan int     `json:"urutan"`
}

type UpdateLabelBlogRequest struct {
	NamaID *string `json:"nama_id" validate:"omitempty,max=100"`
	NamaEN *string `json:"nama_en" validate:"omitempty,max=100"`
	Slug   *string `json:"slug" validate:"omitempty,max=100"`
	Urutan *int    `json:"urutan"`
}

// Dropdown Response DTOs
type KategoriBlogDropdownResponse struct {
	ID   uuid.UUID              `json:"id"`
	Nama map[string]interface{} `json:"nama"` // {"id": "...", "en": "..."}
	Slug string                 `json:"slug"`
}

type LabelBlogDropdownResponse struct {
	ID   uuid.UUID              `json:"id"`
	Nama map[string]interface{} `json:"nama"` // {"id": "...", "en": "..."}
	Slug string                 `json:"slug"`
}

// Reorder DTO
type ReorderItem struct {
	ID     uuid.UUID `json:"id" validate:"required"`
	Urutan int       `json:"urutan" validate:"required"`
}

type ReorderRequest struct {
	Items []ReorderItem `json:"items" validate:"required,min=1"`
}
