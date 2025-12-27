package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Buyer struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama            string         `gorm:"type:varchar(100);not null" json:"nama"`
	Username        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email           string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password        string         `gorm:"type:varchar(255);not null" json:"-"`
	Telepon         *string        `gorm:"type:varchar(20)" json:"telepon"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	IsVerified      bool           `gorm:"default:false" json:"is_verified"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Alamat []AlamatBuyer `gorm:"foreignKey:BuyerID" json:"alamat,omitempty"`
}

func (Buyer) TableName() string {
	return "buyer"
}

// Response DTO
type BuyerResponse struct {
	ID              string     `json:"id"`
	Nama            string     `json:"nama"`
	Username        string     `json:"username"`
	Email           string     `json:"email"`
	Telepon         *string    `json:"telepon"`
	IsActive        bool       `json:"is_active"`
	IsVerified      bool       `json:"is_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Simple response for nested objects
type BuyerSimpleResponse struct {
	ID      string  `json:"id"`
	Nama    string  `json:"nama"`
	Email   string  `json:"email"`
	Telepon *string `json:"telepon"`
}
