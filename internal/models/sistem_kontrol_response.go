package models

import "time"

// ========================================
// Force Update Response (Admin)
// ========================================

type ForceUpdateListResponse struct {
	ID         string    `json:"id"`
	KodeVersi  string    `json:"kode_versi"`
	UpdateType string    `json:"update_type"`
	Platform   string    `json:"platform"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type ForceUpdateDetailResponse struct {
	ID                string    `json:"id"`
	KodeVersi         string    `json:"kode_versi"`
	UpdateType        string    `json:"update_type"`
	InformasiUpdate   string    `json:"informasi_update"`
	InformasiUpdateEn string    `json:"informasi_update_en"`
	Platform          string    `json:"platform"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ========================================
// Force Update Response (Public/Mobile)
// ========================================

type CheckVersionResponse struct {
	ShouldUpdate    bool    `json:"should_update"`
	UpdateType      *string `json:"update_type"`      // OPTIONAL / MANDATORY (null jika tidak perlu update)
	LatestVersion   *string `json:"latest_version"`   // Versi terbaru
	CurrentVersion  string  `json:"current_version"`  // Versi yang dikirim client
	InformasiUpdate *string `json:"informasi_update"` // Changelog (sesuai lang param)
	Platform        *string `json:"platform"`         // ALL / ANDROID / IOS
	StoreURL        *string `json:"store_url"`        // URL ke Play Store / App Store
}

// ========================================
// Mode Maintenance Response (Admin)
// ========================================

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
