package models

import (
	"time"

	"github.com/google/uuid"
)

type BuyerDisclaimerConsent struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BuyerID      uuid.UUID `gorm:"type:uuid;not null"                               json:"buyer_id"`
	PesananID    uuid.UUID `gorm:"type:uuid;not null"                               json:"pesanan_id"`
	DisclaimerID uuid.UUID `gorm:"type:uuid;not null"                               json:"disclaimer_id"`
	DisetujuiAt  time.Time `gorm:"type:timestamptz;not null;default:now()"          json:"disetujui_at"`
	IPAddress    *string   `gorm:"type:varchar(45)"                                 json:"ip_address"`
	UserAgent    *string   `gorm:"type:text"                                        json:"user_agent"`
	CreatedAt    time.Time `gorm:"type:timestamptz;autoCreateTime"                  json:"created_at"`

	// Relations (optional preload)
	Buyer      *Buyer      `gorm:"foreignKey:BuyerID"      json:"buyer,omitempty"`
	Pesanan    *Pesanan    `gorm:"foreignKey:PesananID"    json:"pesanan,omitempty"`
	Disclaimer *Disclaimer `gorm:"foreignKey:DisclaimerID" json:"disclaimer,omitempty"`
}

func (BuyerDisclaimerConsent) TableName() string {
	return "buyer_disclaimer_consent"
}

// --- Response DTOs ---

type DisclaimerConsentAdminResponse struct {
	ID           string    `json:"id"`
	BuyerID      string    `json:"buyer_id"`
	BuyerNama    string    `json:"buyer_nama"`
	BuyerEmail   string    `json:"buyer_email"`
	PesananID    string    `json:"pesanan_id"`
	PesananKode  string    `json:"pesanan_kode"`
	DisclaimerID string    `json:"disclaimer_id"`
	DisetujuiAt  time.Time `json:"disetujui_at"`
	IPAddress    *string   `json:"ip_address"`
	UserAgent    *string   `json:"user_agent"`
	CreatedAt    time.Time `json:"created_at"`
}
