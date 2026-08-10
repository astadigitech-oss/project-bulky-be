package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DelivereeEnvironment membedakan credential/vehicle_type_id sandbox vs production Deliveree.
type DelivereeEnvironment string

const (
	DelivereeEnvSandbox    DelivereeEnvironment = "sandbox"
	DelivereeEnvProduction DelivereeEnvironment = "production"
)

// DelivereeVehicleType adalah master data kendaraan Deliveree (hasil Sync dari
// GET /public_api/v10/vehicle_types), dipakai sebagai basis keputusan pemilihan
// kendaraan saat create booking Deliveree berdasarkan kubikasi & berat.
type DelivereeVehicleType struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama              string         `gorm:"type:varchar(100);not null" json:"nama"`
	IDDeliveree       int            `gorm:"column:id_deliveree;not null" json:"id_deliveree"`
	Environment       string         `gorm:"type:varchar(20);not null" json:"environment"`
	KubikasiMax       float64        `gorm:"column:kubikasi_max;type:decimal(10,3);not null;default:0" json:"kubikasi_max"`
	BeratMax          float64        `gorm:"column:berat_max;type:decimal(10,2);not null;default:0" json:"berat_max"`
	ThresholdKubikasi float64        `gorm:"column:threshold_kubikasi;type:decimal(10,3);not null;default:0" json:"threshold_kubikasi"`
	ThresholdBerat    float64        `gorm:"column:threshold_berat;type:decimal(10,2);not null;default:0" json:"threshold_berat"`
	CargoLength       *float64       `gorm:"column:cargo_length;type:decimal(10,2)" json:"cargo_length"`
	CargoWidth        *float64       `gorm:"column:cargo_width;type:decimal(10,2)" json:"cargo_width"`
	CargoHeight       *float64       `gorm:"column:cargo_height;type:decimal(10,2)" json:"cargo_height"`
	IsActive          bool           `gorm:"column:is_active;default:true" json:"is_active"`
	LastSyncedAt      *time.Time     `gorm:"column:last_synced_at;type:timestamptz" json:"last_synced_at"`
	CreatedAt         time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"type:timestamptz;index" json:"-"`
}

func (DelivereeVehicleType) TableName() string {
	return "deliveree_vehicle_type"
}

// ========================================
// Request DTOs
// ========================================

// UpdateDelivereeVehicleTypeRequest dipakai untuk override manual field
// yang dihitung otomatis oleh Sync (mis. threshold) atau mengubah status aktif.
type UpdateDelivereeVehicleTypeRequest struct {
	ThresholdKubikasi *float64 `json:"threshold_kubikasi" binding:"omitempty,min=0"`
	ThresholdBerat    *float64 `json:"threshold_berat" binding:"omitempty,min=0"`
	IsActive          *bool    `json:"is_active"`
}

type DelivereeVehicleTypeFilterRequest struct {
	PaginationRequest
	Environment *string `query:"environment"`
}

// ========================================
// Response DTOs
// ========================================

type DelivereeVehicleTypeResponse struct {
	ID                string     `json:"id"`
	Nama              string     `json:"nama"`
	IDDeliveree       int        `json:"id_deliveree"`
	Environment       string     `json:"environment"`
	KubikasiMax       float64    `json:"kubikasi_max"`
	BeratMax          float64    `json:"berat_max"`
	ThresholdKubikasi float64    `json:"threshold_kubikasi"`
	ThresholdBerat    float64    `json:"threshold_berat"`
	CargoLength       *float64   `json:"cargo_length"`
	CargoWidth        *float64   `json:"cargo_width"`
	CargoHeight       *float64   `json:"cargo_height"`
	IsActive          bool       `json:"is_active"`
	LastSyncedAt      *time.Time `json:"last_synced_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// SyncDelivereeVehicleTypeResponse merangkum hasil satu kali proses Sync.
type SyncDelivereeVehicleTypeResponse struct {
	Environment  string    `json:"environment"`
	TotalFromAPI int       `json:"total_from_api"`
	Created      int       `json:"created"`
	Updated      int       `json:"updated"`
	Deactivated  int       `json:"deactivated"`
	SyncedAt     time.Time `json:"synced_at"`
}
