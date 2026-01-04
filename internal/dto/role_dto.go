package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateRoleRequest DTO untuk membuat role baru
type CreateRoleRequest struct {
	Nama          string   `json:"nama" binding:"required,min=3,max=100"`
	Kode          string   `json:"kode" binding:"required,min=3,max=50,uppercase_snake"`
	Deskripsi     string   `json:"deskripsi" binding:"max=500"`
	IsActive      *bool    `json:"is_active"`
	PermissionIDs []string `json:"permission_ids" binding:"dive,uuid"`
}

// UpdateRoleRequest DTO untuk update role
type UpdateRoleRequest struct {
	Nama          string   `json:"nama" binding:"omitempty,min=3,max=100"`
	Kode          string   `json:"kode" binding:"omitempty,min=3,max=50,uppercase_snake"`
	Deskripsi     *string  `json:"deskripsi" binding:"omitempty,max=500"`
	IsActive      *bool    `json:"is_active"`
	PermissionIDs []string `json:"permission_ids" binding:"omitempty,dive,uuid"`
}

// RoleQueryParams untuk query parameter list role
type RoleQueryParams struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PerPage  int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
	SortBy   string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at nama"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// RoleResponse DTO response untuk role
type RoleResponse struct {
	ID                 uuid.UUID                       `json:"id"`
	Nama               string                          `json:"nama"`
	Kode               string                          `json:"kode"`
	Deskripsi          string                          `json:"deskripsi,omitempty"`
	IsActive           bool                            `json:"is_active"`
	Permissions        []PermissionResponse            `json:"permissions,omitempty"`
	PermissionsGrouped map[string][]PermissionResponse `json:"permissions_grouped,omitempty"`
	Admins             []AdminMinimalResponse          `json:"admins,omitempty"`
	TotalAdmin         int                             `json:"total_admin,omitempty"`
	TotalPermission    int                             `json:"total_permission,omitempty"`
	CreatedAt          time.Time                       `json:"created_at"`
	UpdatedAt          time.Time                       `json:"updated_at"`
}

// RoleListResponse untuk response list (tanpa detail permissions)
type RoleListResponse struct {
	ID              uuid.UUID `json:"id"`
	Nama            string    `json:"nama"`
	Kode            string    `json:"kode"`
	Deskripsi       string    `json:"deskripsi,omitempty"`
	IsActive        bool      `json:"is_active"`
	TotalAdmin      int       `json:"total_admin"`
	TotalPermission int       `json:"total_permission"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PermissionResponse DTO response untuk permission
type PermissionResponse struct {
	ID        uuid.UUID `json:"id"`
	Nama      string    `json:"nama"`
	Kode      string    `json:"kode"`
	Modul     string    `json:"modul,omitempty"`
	Deskripsi string    `json:"deskripsi,omitempty"`
}

// AdminMinimalResponse untuk response minimal admin
type AdminMinimalResponse struct {
	ID    uuid.UUID `json:"id"`
	Nama  string    `json:"nama"`
	Email string    `json:"email"`
}

// PermissionQueryParams untuk query parameter list permission
type PermissionQueryParams struct {
	Grouped bool   `form:"grouped"`
	Modul   string `form:"modul"`
}

// PermissionListResponse untuk list permissions (flat)
type PermissionListResponse struct {
	ID        uuid.UUID `json:"id"`
	Nama      string    `json:"nama"`
	Kode      string    `json:"kode"`
	Modul     string    `json:"modul"`
	Deskripsi string    `json:"deskripsi,omitempty"`
}
