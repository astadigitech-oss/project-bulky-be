package models

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WhatsAppHandler struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	NomorWA   string         `gorm:"type:varchar(20);not null" json:"nomor_wa"`
	PesanAwal string         `gorm:"type:text;not null" json:"pesan_awal"`
	IsActive  bool           `gorm:"default:false" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WhatsAppHandler) TableName() string {
	return "whatsapp_handler"
}

// GetWhatsAppURL generates wa.me URL with encoded message
func (w *WhatsAppHandler) GetWhatsAppURL() string {
	return fmt.Sprintf("https://wa.me/%s?text=%s", w.NomorWA, url.QueryEscape(w.PesanAwal))
}
