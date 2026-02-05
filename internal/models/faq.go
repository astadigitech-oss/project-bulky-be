package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FAQ Model - Tabel terpisah untuk FAQ
type FAQ struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Question   string         `gorm:"type:varchar(500);not null" json:"question"`
	QuestionEN string         `gorm:"column:question_en;type:varchar(500);not null" json:"question_en"`
	Answer     string         `gorm:"type:text;not null" json:"answer"`
	AnswerEN   string         `gorm:"column:answer_en;type:text;not null" json:"answer_en"`
	Urutan     int            `gorm:"type:int;not null;default:0" json:"urutan"`
	IsActive   bool           `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (FAQ) TableName() string {
	return "faq"
}

// FAQ Request DTOs
type FAQCreateRequest struct {
	Question   string `json:"question" binding:"required,min=5,max=500"`
	QuestionEN string `json:"question_en" binding:"required,min=5,max=500"`
	Answer     string `json:"answer" binding:"required,min=10"`
	AnswerEN   string `json:"answer_en" binding:"required,min=10"`
	IsActive   *bool  `json:"is_active"`
}

type FAQUpdateRequest struct {
	Question   string `json:"question" binding:"required,min=5,max=500"`
	QuestionEN string `json:"question_en" binding:"required,min=5,max=500"`
	Answer     string `json:"answer" binding:"required,min=10"`
	AnswerEN   string `json:"answer_en" binding:"required,min=10"`
	IsActive   *bool  `json:"is_active"`
}

type FAQFilterRequest struct {
	PaginationRequest
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
	SortBy   string `form:"sort_by"`
}

func (p *FAQFilterRequest) SetDefaults() {
	p.PaginationRequest.SetDefaults()

	// Whitelist sort_by
	allowedSortBy := map[string]bool{
		"urutan": true, "question": true, "created_at": true, "updated_at": true,
	}
	if p.SortBy == "" || !allowedSortBy[p.SortBy] {
		p.SortBy = "urutan"
	}
}

// FAQ Response DTOs
type FAQResponse struct {
	ID         string `json:"id"`
	Question   string `json:"question"`
	QuestionEN string `json:"question_en"`
	Answer     string `json:"answer"`
	AnswerEN   string `json:"answer_en"`
	Urutan     int    `json:"urutan"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// Public FAQ Response (simplified)
type FAQPublicResponse struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
