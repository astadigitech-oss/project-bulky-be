package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role untuk admin users
type Role struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama      string         `gorm:"type:varchar(50);not null" json:"nama"`
	Kode      string         `gorm:"type:varchar(30);not null;unique" json:"kode"`
	Deskripsi *string        `gorm:"type:text" json:"deskripsi"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Permissions []Permission `gorm:"many2many:role_permission;" json:"permissions,omitempty"`
}

// Response format untuk role (array sederhana)
type RoleResponseFormat struct {
	ID        string  `json:"id"`
	Nama      string  `json:"nama"`
	Kode      string  `json:"kode"`
	Deskripsi *string `json:"deskripsi"`
}

// Response format untuk role detail dengan permissions
type RoleDetailResponse struct {
	ID          string                     `json:"id"`
	Nama        string                     `json:"nama"`
	Kode        string                     `json:"kode"`
	Deskripsi   *string                    `json:"deskripsi"`
	Permissions []PermissionSimpleResponse `json:"permissions"`
}

func (Role) TableName() string {
	return "role"
}

// Role codes
const (
	RoleSuperAdmin = "SUPER_ADMIN"
	RoleAdmin      = "ADMIN"
	RoleStaff      = "STAFF"
)
