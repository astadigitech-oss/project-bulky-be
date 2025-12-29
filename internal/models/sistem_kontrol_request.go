package models

// ========================================
// Force Update Request
// ========================================

type CreateForceUpdateRequest struct {
	KodeVersi       string `json:"kode_versi" binding:"required,max=20"`
	UpdateType      string `json:"update_type" binding:"required,oneof=OPTIONAL MANDATORY"`
	InformasiUpdate string `json:"informasi_update" binding:"required"`
	IsActive        bool   `json:"is_active"`
}

type UpdateForceUpdateRequest struct {
	KodeVersi       *string `json:"kode_versi" binding:"omitempty,max=20"`
	UpdateType      *string `json:"update_type" binding:"omitempty,oneof=OPTIONAL MANDATORY"`
	InformasiUpdate *string `json:"informasi_update"`
	IsActive        *bool   `json:"is_active"`
}

// ========================================
// Mode Maintenance Request
// ========================================

type CreateMaintenanceRequest struct {
	Judul           string `json:"judul" binding:"required,max=100"`
	TipeMaintenance string `json:"tipe_maintenance" binding:"required,oneof=BUG ERROR BIG_UPDATE OTHER"`
	Deskripsi       string `json:"deskripsi" binding:"required"`
	IsActive        bool   `json:"is_active"`
}

type UpdateMaintenanceRequest struct {
	Judul           *string `json:"judul" binding:"omitempty,max=100"`
	TipeMaintenance *string `json:"tipe_maintenance" binding:"omitempty,oneof=BUG ERROR BIG_UPDATE OTHER"`
	Deskripsi       *string `json:"deskripsi"`
	IsActive        *bool   `json:"is_active"`
}
