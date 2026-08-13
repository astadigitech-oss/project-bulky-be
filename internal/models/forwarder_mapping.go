package models

import "time"

// ========================================
// Request DTOs
// ========================================

// ForwarderCityFilterRequest filter untuk list city mapping (panel admin).
type ForwarderCityFilterRequest struct {
	PaginationRequest
}

// ForwarderSubdistrictFilterRequest filter untuk list subdistrict mapping (panel admin).
type ForwarderSubdistrictFilterRequest struct {
	PaginationRequest
	// ForwarderCityID filter opsional berdasarkan city_id Forwarder.
	ForwarderCityID *int `query:"forwarder_city_id"`
}

// SyncForwarderMappingRequest memungkinkan sync sebagian (opsional):
// - sync_city=true (default) → tarik citylist dari API Forwarder
// - sync_subdistrict=true (default) → tarik subdistrictlist dari API Forwarder
// Kalau keduanya false → 400.
type SyncForwarderMappingRequest struct {
	SyncCity        *bool `json:"sync_city"`
	SyncSubdistrict *bool `json:"sync_subdistrict"`
}

// ========================================
// Response DTOs
// ========================================

type ForwarderCityMappingResponse struct {
	ID                int       `json:"id"`
	KotaPattern       string    `json:"kota_pattern"`
	ForwarderCityID   int       `json:"forwarder_city_id"`
	ForwarderCityName string    `json:"forwarder_city_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ForwarderSubdistrictMappingResponse struct {
	ID                       int       `json:"id"`
	KecamatanPattern         string    `json:"kecamatan_pattern"`
	ForwarderCityID          int       `json:"forwarder_city_id"`
	ForwarderSubdistrictID   int       `json:"forwarder_subdistrict_id"`
	ForwarderSubdistrictName string    `json:"forwarder_subdistrict_name"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

// SyncForwarderMappingResponse merangkum hasil satu kali proses Sync dari API Forwarder.
type SyncForwarderMappingResponse struct {
	CityCreated             int       `json:"city_created"`
	CityUpdated             int       `json:"city_updated"`
	CityTotalFromAPI        int       `json:"city_total_from_api"`
	SubdistrictCreated      int       `json:"subdistrict_created"`
	SubdistrictUpdated      int       `json:"subdistrict_updated"`
	SubdistrictTotalFromAPI int       `json:"subdistrict_total_from_api"`
	SyncedAt                time.Time `json:"synced_at"`
}

type ForwarderCityMapping struct {
	ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
	KotaPattern       string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"kota_pattern"`
	ForwarderCityID   int       `gorm:"not null" json:"forwarder_city_id"`
	ForwarderCityName string    `gorm:"type:varchar(100);not null" json:"forwarder_city_name"`
	CreatedAt         time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
}

func (ForwarderCityMapping) TableName() string {
	return "forwarder_city_mapping"
}

type ForwarderSubdistrictMapping struct {
	ID                       int       `gorm:"primaryKey;autoIncrement" json:"id"`
	KecamatanPattern         string    `gorm:"type:varchar(100);not null" json:"kecamatan_pattern"`
	ForwarderCityID          int       `gorm:"not null" json:"forwarder_city_id"`
	ForwarderSubdistrictID   int       `gorm:"not null" json:"forwarder_subdistrict_id"`
	ForwarderSubdistrictName string    `gorm:"type:varchar(100);not null" json:"forwarder_subdistrict_name"`
	CreatedAt                time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt                time.Time `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
}

func (ForwarderSubdistrictMapping) TableName() string {
	return "forwarder_subdistrict_mapping"
}
