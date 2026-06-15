package models

import "time"

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
