package models

import (
	"time"

	"github.com/google/uuid"
)

type UserType string

const (
	UserTypeAdmin  UserType = "ADMIN"
	UserTypeBuyer  UserType = "BUYER"
	UserTypeSystem UserType = "SYSTEM"
)

// RefreshToken untuk token revocation & multi-device support
type RefreshToken struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserType   UserType   `gorm:"type:varchar(20);not null" json:"user_type"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Token      string     `gorm:"type:varchar(500);not null;unique" json:"token"`
	DeviceInfo *string    `gorm:"type:varchar(255)" json:"device_info"`
	IPAddress  *string    `gorm:"type:varchar(50)" json:"ip_address"`
	ExpiredAt  time.Time  `gorm:"not null" json:"expired_at"`
	IsRevoked  bool       `gorm:"default:false" json:"is_revoked"`
	RevokedAt  *time.Time `json:"revoked_at"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_token"
}
