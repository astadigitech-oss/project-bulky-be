package models

import "time"

// ========================================
// Force Update Response (Admin)
// ========================================

type ForceUpdateListResponse struct {
	ID         string    `json:"id"`
	KodeVersi  string    `json:"kode_versi"`
	UpdateType string    `json:"update_type"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type ForceUpdateDetailResponse struct {
	ID              string    `json:"id"`
	KodeVersi       string    `json:"kode_versi"`
	UpdateType      string    `json:"update_type"`
	InformasiUpdate string    `json:"informasi_update"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ========================================
// Force Update Response (Public/Mobile)
// ========================================

type CheckVersionResponse struct {
	ShouldUpdate    bool    `json:"should_update"`
	UpdateType      *string `json:"update_type"`      // OPTIONAL / MANDATORY (null jika tidak perlu update)
	LatestVersion   *string `json:"latest_version"`   // Versi terbaru
	CurrentVersion  string  `json:"current_version"`  // Versi yang dikirim client
	InformasiUpdate *string `json:"informasi_update"` // Changelog
	StoreURL        *string `json:"store_url"`        // URL ke Play Store / App Store
}

// ========================================
// Mode Maintenance Response (Admin)
// ========================================

type MaintenanceListResponse struct {
	ID              string    `json:"id"`
	Judul           string    `json:"judul"`
	TipeMaintenance string    `json:"tipe_maintenance"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

type MaintenanceDetailResponse struct {
	ID              string    `json:"id"`
	Judul           string    `json:"judul"`
	TipeMaintenance string    `json:"tipe_maintenance"`
	Deskripsi       string    `json:"deskripsi"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ========================================
// Mode Maintenance Response (Public)
// ========================================

type CheckMaintenanceResponse struct {
	IsMaintenance   bool    `json:"is_maintenance"`
	Judul           *string `json:"judul"`
	TipeMaintenance *string `json:"tipe_maintenance"`
	Deskripsi       *string `json:"deskripsi"`
}

// ========================================
// App Status Response (Combined)
// ========================================

type AppStatusResponse struct {
	Maintenance CheckMaintenanceResponse `json:"maintenance"`
	Version     CheckVersionResponse     `json:"version"`
}
