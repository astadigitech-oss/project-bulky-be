package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeasonalCampaign is the admin-owned source of truth for seasonal branding.
// Its status is deliberately derived from publication, cancellation and period.
type SeasonalCampaign struct {
	ID                         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Nama                       string         `gorm:"type:varchar(100);not null" json:"nama"`
	IsPublished                bool           `gorm:"not null;default:false" json:"is_published"`
	TanggalMulai               *time.Time     `gorm:"type:timestamptz" json:"tanggal_mulai"`
	TanggalSelesai             *time.Time     `gorm:"type:timestamptz" json:"tanggal_selesai"`
	CancelledAt                *time.Time     `gorm:"type:timestamptz" json:"cancelled_at"`
	WebLogoURL                 *string        `gorm:"column:web_logo_url;type:varchar(500)" json:"-"`
	MobileLoadingLogoURL       *string        `gorm:"column:mobile_loading_logo_url;type:varchar(500)" json:"-"`
	WebNavbarDecorationURL     *string        `gorm:"column:web_navbar_decoration_url;type:varchar(500)" json:"-"`
	MobileTopAppBarOrnamentURL *string        `gorm:"column:mobile_top_app_bar_ornament_url;type:varchar(500)" json:"-"`
	CreatedBy                  *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy                  *uuid.UUID     `gorm:"type:uuid" json:"updated_by,omitempty"`
	CreatedAt                  time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt                  time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt                  gorm.DeletedAt `gorm:"type:timestamptz;index" json:"deleted_at,omitempty"`
}

func (SeasonalCampaign) TableName() string { return "seasonal_campaign" }

func (c *SeasonalCampaign) StatusAt(now time.Time) string {
	if c.CancelledAt != nil {
		return "cancelled"
	}
	if !c.IsPublished {
		return "draft"
	}
	if c.TanggalMulai != nil && now.Before(*c.TanggalMulai) {
		return "scheduled"
	}
	if c.TanggalSelesai != nil && !now.Before(*c.TanggalSelesai) {
		return "ended"
	}
	return "active"
}
