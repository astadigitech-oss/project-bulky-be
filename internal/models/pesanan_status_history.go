package models

import (
	"time"

	"github.com/google/uuid"
)

type StatusHistoryType string

const (
	StatusHistoryTypeOrder   StatusHistoryType = "ORDER"
	StatusHistoryTypePayment StatusHistoryType = "PAYMENT"
)

type PesananStatusHistory struct {
	ID         uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PesananID  uuid.UUID         `gorm:"type:uuid;not null" json:"pesanan_id"`
	StatusFrom *string           `gorm:"type:varchar(20)" json:"status_from"`
	StatusTo   string            `gorm:"type:varchar(20);not null" json:"status_to"`
	StatusType StatusHistoryType `gorm:"type:status_history_type;not null" json:"status_type"`
	ChangedBy  *uuid.UUID        `gorm:"type:uuid" json:"changed_by"`
	Note       *string           `gorm:"type:text" json:"note"`
	CreatedAt  time.Time         `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Admin *Admin `gorm:"foreignKey:ChangedBy" json:"admin,omitempty"`
}

func (PesananStatusHistory) TableName() string {
	return "pesanan_status_history"
}

// Response DTOs
type PesananStatusHistoryResponse struct {
	ID         string         `json:"id"`
	PesananID  string         `json:"pesanan_id"`
	StatusFrom *string        `json:"status_from"`
	StatusTo   string         `json:"status_to"`
	StatusType string         `json:"status_type"`
	ChangedBy  *string        `json:"changed_by"`
	Note       *string        `json:"note"`
	CreatedAt  time.Time      `json:"created_at"`
	Admin      *AdminResponse `json:"admin,omitempty"`
}
