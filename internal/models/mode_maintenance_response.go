package models

import "time"

type MaintenanceListResponse struct {
	ID              string    `json:"id"`
	Judul           string    `json:"judul"`
	JudulEn         string    `json:"judul_en"`
	TipeMaintenance string    `json:"tipe_maintenance"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

type MaintenanceDetailResponse struct {
	ID              string    `json:"id"`
	Judul           string    `json:"judul"`
	JudulEn         string    `json:"judul_en"`
	TipeMaintenance string    `json:"tipe_maintenance"`
	Deskripsi       string    `json:"deskripsi"`
	DeskripsiEn     string    `json:"deskripsi_en"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CheckMaintenanceResponse struct {
	IsMaintenance   bool    `json:"is_maintenance"`
	Judul           *string `json:"judul"`
	TipeMaintenance *string `json:"tipe_maintenance"`
	Deskripsi       *string `json:"deskripsi"`
}

type AppStatusResponse struct {
	Maintenance CheckMaintenanceResponse `json:"maintenance"`
	Version     CheckVersionResponse     `json:"version"`
}
