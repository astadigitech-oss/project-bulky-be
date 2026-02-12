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
	RoleID      uuid.UUID      `gorm:"type:uuid;not null" json:"role_id"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt *time.Time     `gorm:"type:timestamptz" json:"last_login_at"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`

	// Relations
	Role *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (Admin) TableName() string {
	return "admin"
}
