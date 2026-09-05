package models

type CreateMaintenanceRequest struct {
	Judul           string `json:"judul" binding:"required,max=100"`
	JudulEn         string `json:"judul_en" binding:"required,max=100"`
	TipeMaintenance string `json:"tipe_maintenance" binding:"required,oneof=BUG ERROR BIG_UPDATE OTHER"`
	Deskripsi       string `json:"deskripsi" binding:"required"`
	DeskripsiEn     string `json:"deskripsi_en" binding:"required"`
	IsActive        bool   `json:"is_active"`
}

type UpdateMaintenanceRequest struct {
	Judul           *string `json:"judul" binding:"omitempty,max=100"`
	JudulEn         *string `json:"judul_en" binding:"omitempty,max=100"`
	TipeMaintenance *string `json:"tipe_maintenance" binding:"omitempty,oneof=BUG ERROR BIG_UPDATE OTHER"`
	Deskripsi       *string `json:"deskripsi"`
	DeskripsiEn     *string `json:"deskripsi_en"`
	IsActive        *bool   `json:"is_active"`
}
