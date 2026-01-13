package models

import (
	"time"

	"github.com/google/uuid"
)

// Permission/hak akses granular
type Permission struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama      string    `gorm:"type:varchar(100);not null" json:"nama"`
	Kode      string    `gorm:"type:varchar(50);not null;unique" json:"kode"`
	Modul     string    `gorm:"type:varchar(50);not null" json:"modul"`
	Deskripsi *string   `gorm:"type:text" json:"deskripsi"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Permission) TableName() string {
	return "permission"
}

// PermissionSimpleResponse untuk response sederhana permission
type PermissionSimpleResponse struct {
	ID        string  `json:"id"`
	Nama      string  `json:"nama"`
	Deskripsi *string `json:"deskripsi"`
}

// RolePermission adalah pivot table
type RolePermission struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null" json:"permission_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (RolePermission) TableName() string {
	return "role_permission"
}
