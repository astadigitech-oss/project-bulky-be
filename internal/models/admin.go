package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama        string         `gorm:"type:varchar(100);not null" json:"nama"`
	Email       string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password    string         `gorm:"type:varchar(255);not null" json:"-"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Admin) TableName() string {
	return "admin"
}
