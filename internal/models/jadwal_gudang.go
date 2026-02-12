package models

import (
	"time"

	"github.com/google/uuid"
)

type JadwalGudang struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	WarehouseID uuid.UUID `gorm:"type:uuid;not null" json:"warehouse_id"`
	Hari        int       `gorm:"not null" json:"hari"` // 0=Minggu, 1=Senin, ..., 6=Sabtu
	JamBuka     *string   `gorm:"type:time" json:"jam_buka"`
	JamTutup    *string   `gorm:"type:time" json:"jam_tutup"`
	IsBuka      bool      `gorm:"default:false" json:"is_buka"`
	CreatedAt   time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
}

func (JadwalGudang) TableName() string {
	return "jadwal_gudang"
}

// GetHariNama returns Indonesian day name
func (j *JadwalGudang) GetHariNama() string {
	days := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	if j.Hari >= 0 && j.Hari <= 6 {
		return days[j.Hari]
	}
	return ""
}

// Request DTOs
type JadwalGudangItem struct {
	Hari     int     `json:"hari" binding:"min=0,max=6"`
	JamBuka  *string `json:"jam_buka"`  // Format: "09:00"
	JamTutup *string `json:"jam_tutup"` // Format: "18:00"
	IsBuka   bool    `json:"is_buka"`
}

type UpdateJadwalGudangRequest struct {
	Jadwal []JadwalGudangItem `json:"jadwal" binding:"required,dive"`
}

// Response DTOs
type JadwalGudangResponse struct {
	ID       string  `json:"id"`
	Hari     int     `json:"hari"`
	HariNama string  `json:"hari_nama"`
	JamBuka  *string `json:"jam_buka"`
	JamTutup *string `json:"jam_tutup"`
	IsBuka   bool    `json:"is_buka"`
}
