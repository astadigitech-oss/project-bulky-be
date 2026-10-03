package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	PushPromotionStatusDraft     = "DRAFT"
	PushPromotionStatusScheduled = "SCHEDULED"
	PushPromotionStatusSending   = "SENDING"
	PushPromotionStatusSent      = "SENT"
	PushPromotionStatusFailed    = "FAILED"
)

type PushPromotionNotification struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Nama         string     `gorm:"type:varchar(120);not null" json:"nama"`
	TitleID      string     `gorm:"type:varchar(160);not null" json:"title_id"`
	BodyID       string     `gorm:"type:text;not null" json:"body_id"`
	TitleEN      string     `gorm:"type:varchar(160);not null" json:"title_en"`
	BodyEN       string     `gorm:"type:text;not null" json:"body_en"`
	DeepLink     string     `gorm:"type:varchar(500);not null" json:"deep_link"`
	ScheduledAt  *time.Time `gorm:"type:timestamptz" json:"scheduled_at"`
	Status       string     `gorm:"type:varchar(20);not null" json:"status"`
	SentAt       *time.Time `gorm:"type:timestamptz" json:"sent_at"`
	SuccessCount int        `gorm:"not null;default:0" json:"success_count"`
	FailureCount int        `gorm:"not null;default:0" json:"failure_count"`
	LastError    *string    `gorm:"type:text" json:"last_error,omitempty"`
	CreatedBy    *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy    *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (PushPromotionNotification) TableName() string {
	return "push_promotion_notifications"
}

type CreatePushPromotionRequest struct {
	Nama        string  `json:"nama" binding:"required,min=1,max=120"`
	TitleID     string  `json:"title_id" binding:"required,min=1,max=50"`
	BodyID      string  `json:"body_id" binding:"required,min=1,max=120"`
	TitleEN     string  `json:"title_en" binding:"required,min=1,max=50"`
	BodyEN      string  `json:"body_en" binding:"required,min=1,max=120"`
	ScheduledAt *string `json:"scheduled_at"`
}

type UpdatePushPromotionRequest struct {
	Nama        *string `json:"nama" binding:"omitempty,min=1,max=120"`
	TitleID     *string `json:"title_id" binding:"omitempty,min=1,max=50"`
	BodyID      *string `json:"body_id" binding:"omitempty,min=1,max=120"`
	TitleEN     *string `json:"title_en" binding:"omitempty,min=1,max=50"`
	BodyEN      *string `json:"body_en" binding:"omitempty,min=1,max=120"`
	ScheduledAt *string `json:"scheduled_at"`
}

type PushPromotionListRequest struct {
	PaginationRequest
	Status string `query:"status"`
}

type SchedulePushPromotionRequest struct {
	ScheduledAt string `json:"scheduled_at" binding:"required"`
}

type PushPromotionListResponse struct {
	Data []PushPromotionNotification `json:"data"`
	Meta PaginationMeta              `json:"meta"`
}
