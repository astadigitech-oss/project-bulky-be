package models

import "time"

type ForceUpdateListResponse struct {
	ID                 string    `json:"id"`
	KodeVersi          string    `json:"kode_versi"`
	MinimumBuildNumber int64     `json:"minimum_build_number"`
	UpdateType         string    `json:"update_type"`
	Platform           string    `json:"platform"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
}

type ForceUpdateDetailResponse struct {
	ID                 string    `json:"id"`
	KodeVersi          string    `json:"kode_versi"`
	MinimumBuildNumber int64     `json:"minimum_build_number"`
	UpdateType         string    `json:"update_type"`
	InformasiUpdate    string    `json:"informasi_update"`
	InformasiUpdateEn  string    `json:"informasi_update_en"`
	Platform           string    `json:"platform"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// CheckVersionResponse is retained for consumers that share the system-control model.
type CheckVersionResponse struct {
	ShouldUpdate    bool    `json:"should_update"`
	UpdateType      *string `json:"update_type"`
	LatestVersion   *string `json:"latest_version"`
	CurrentVersion  string  `json:"current_version"`
	InformasiUpdate *string `json:"update_information"`
	Platform        *string `json:"platform"`
	StoreURL        *string `json:"store_url"`
}
