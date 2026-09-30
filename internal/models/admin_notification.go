package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminPushDevice struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AdminID    uuid.UUID `gorm:"type:uuid;not null" json:"admin_id"`
	FID        string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"`
	Platform   string    `gorm:"type:varchar(20);not null" json:"platform"`
	LastSeenAt time.Time `gorm:"type:timestamptz;not null" json:"last_seen_at"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (AdminPushDevice) TableName() string { return "admin_push_devices" }

type AdminNotification struct {
	ID        uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	AdminID   uuid.UUID         `gorm:"type:uuid;not null" json:"-"`
	EventKey  string            `gorm:"type:varchar(160);not null" json:"-"`
	Type      string            `gorm:"type:varchar(50);not null" json:"type"`
	Title     string            `gorm:"type:varchar(160);not null" json:"title"`
	Body      string            `gorm:"type:text;not null" json:"body"`
	Data      map[string]string `gorm:"type:jsonb;serializer:json;not null" json:"data"`
	IsRead    bool              `gorm:"not null;default:false" json:"is_read"`
	ReadAt    *time.Time        `gorm:"type:timestamptz" json:"read_at,omitempty"`
	CreatedAt time.Time         `gorm:"type:timestamptz;not null" json:"created_at"`
}

func (AdminNotification) TableName() string { return "admin_notifications" }

type AdminNotificationListResponse struct {
	Notifications []AdminNotification `json:"notifications"`
	UnreadCount   int64               `json:"unread_count"`
}
