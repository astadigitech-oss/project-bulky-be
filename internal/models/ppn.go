package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PPN struct {
	ID         uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Persentase decimal.Decimal `gorm:"type:decimal(5,2);not null" json:"persentase"`
	IsActive   bool            `gorm:"default:false" json:"is_active"`
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (PPN) TableName() string {
	return "ppn"
}

// Request DTOs
type CreatePPNRequest struct {
	Persentase float64 `json:"persentase" binding:"required,min=0,max=100"`
	IsActive   bool    `json:"is_active"`
}

type UpdatePPNRequest struct {
	Persentase *float64 `json:"persentase" binding:"omitempty,min=0,max=100"`
	IsActive   *bool    `json:"is_active"`
}

// Response DTOs
type PPNResponse struct {
	ID         string    `json:"id"`
	Persentase float64   `json:"persentase"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
