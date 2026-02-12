package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AdminID   uuid.UUID `gorm:"type:uuid;not null" json:"admin_id"`
	Token     string    `gorm:"type:varchar(500);uniqueIndex;not null" json:"-"`
	IPAddress *string   `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent *string   `gorm:"type:varchar(500)" json:"user_agent"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`

	// Relations
	Admin Admin `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
}

func (AdminSession) TableName() string {
	return "admin_session"
}
