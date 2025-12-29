package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ActivityAction string

const (
	ActionLogin        ActivityAction = "LOGIN"
	ActionLogout       ActivityAction = "LOGOUT"
	ActionLoginFailed  ActivityAction = "LOGIN_FAILED"
	ActionCreate       ActivityAction = "CREATE"
	ActionUpdate       ActivityAction = "UPDATE"
	ActionDelete       ActivityAction = "DELETE"
	ActionRestore      ActivityAction = "RESTORE"
	ActionToggleStatus ActivityAction = "TOGGLE_STATUS"
	ActionApprove      ActivityAction = "APPROVE"
	ActionReject       ActivityAction = "REJECT"
	ActionExport       ActivityAction = "EXPORT"
	ActionImport       ActivityAction = "IMPORT"
)

// ActivityLog untuk audit trail
type ActivityLog struct {
	ID         uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserType   string          `gorm:"type:varchar(20);not null" json:"user_type"`
	UserID     *uuid.UUID      `gorm:"type:uuid" json:"user_id"`
	Action     ActivityAction  `gorm:"type:varchar(50);not null" json:"action"`
	Modul      string          `gorm:"type:varchar(50);not null" json:"modul"`
	EntityType *string         `gorm:"type:varchar(50)" json:"entity_type"`
	EntityID   *uuid.UUID      `gorm:"type:uuid" json:"entity_id"`
	Deskripsi  string          `gorm:"type:text;not null" json:"deskripsi"`
	OldData    json.RawMessage `gorm:"type:jsonb" json:"old_data,omitempty"`
	NewData    json.RawMessage `gorm:"type:jsonb" json:"new_data,omitempty"`
	IPAddress  *string         `gorm:"type:varchar(50)" json:"ip_address"`
	UserAgent  *string         `gorm:"type:text" json:"user_agent"`
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

func (ActivityLog) TableName() string {
	return "activity_log"
}
