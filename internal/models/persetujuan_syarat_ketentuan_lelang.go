package models

import (
	"time"

	"github.com/google/uuid"
)

type PersetujuanSyaratKetentuanLelang struct {
	ID              uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	BuyerID         uuid.UUID   `gorm:"column:buyer_id;not null" json:"buyer_id"`
	BidID           uuid.UUID   `gorm:"column:bid_id;not null" json:"bid_id"`
	TermsDocumentID uuid.UUID   `gorm:"column:terms_document_id;not null" json:"terms_document_id"`
	Locale          string      `gorm:"column:locale;not null" json:"locale"`
	ContentHash     string      `gorm:"column:content_hash;not null" json:"content_hash"`
	DisetujuiAt     time.Time   `gorm:"column:disetujui_at;not null" json:"disetujui_at"`
	IPAddress       *string     `gorm:"column:ip_address" json:"ip_address"`
	UserAgent       *string     `gorm:"column:user_agent" json:"user_agent"`
	CreatedAt       time.Time   `gorm:"column:created_at;not null" json:"created_at"`
	Buyer           *Buyer      `gorm:"foreignKey:BuyerID" json:"-"`
	Bid             *AuctionBid `gorm:"foreignKey:BidID" json:"-"`
}

func (PersetujuanSyaratKetentuanLelang) TableName() string {
	return "buyer_auction_terms_consent"
}

type PersetujuanSyaratKetentuanLelangResponse struct {
	ID              string    `json:"id"`
	BuyerID         string    `json:"buyer_id"`
	BuyerNama       string    `json:"buyer_nama"`
	BuyerEmail      string    `json:"buyer_email"`
	BidID           string    `json:"bid_id"`
	BidSequence     int       `json:"bid_sequence"`
	BidAmount       string    `json:"bid_amount"`
	BatchID         string    `json:"batch_id"`
	BatchKode       string    `json:"batch_kode"`
	BatchNama       string    `json:"batch_nama"`
	TermsDocumentID string    `json:"terms_document_id"`
	TermsJudul      string    `json:"terms_judul"`
	Locale          string    `json:"locale"`
	ContentHash     string    `json:"content_hash"`
	DisetujuiAt     time.Time `json:"disetujui_at"`
	IPAddress       *string   `json:"ip_address"`
	UserAgent       *string   `json:"user_agent"`
	CreatedAt       time.Time `json:"created_at"`
}
