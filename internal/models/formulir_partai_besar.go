package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ========================================
// Config (Singleton)
// ========================================

type FormulirPartaiBesarConfig struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	DaftarEmail string    `gorm:"type:text;not null;default:'[]'" json:"daftar_email"` // JSON array
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (FormulirPartaiBesarConfig) TableName() string {
	return "formulir_partai_besar_config"
}

// GetEmails parses JSON array to []string
func (c *FormulirPartaiBesarConfig) GetEmails() []string {
	var emails []string
	json.Unmarshal([]byte(c.DaftarEmail), &emails)
	return emails
}

// SetEmails converts []string to JSON array
func (c *FormulirPartaiBesarConfig) SetEmails(emails []string) {
	data, _ := json.Marshal(emails)
	c.DaftarEmail = string(data)
}

// ========================================
// Anggaran
// ========================================

type FormulirPartaiBesarAnggaran struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Label     string         `gorm:"type:varchar(100);not null" json:"label"`
	Urutan    int            `gorm:"default:0" json:"urutan"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FormulirPartaiBesarAnggaran) TableName() string {
	return "formulir_partai_besar_anggaran"
}

// ========================================
// Submission
// ========================================

type FormulirPartaiBesarSubmission struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BuyerID     *uuid.UUID `gorm:"type:uuid" json:"buyer_id"`
	Nama        string     `gorm:"type:varchar(100);not null" json:"nama"`
	Telepon     string     `gorm:"type:varchar(20);not null" json:"telepon"`
	Alamat      string     `gorm:"type:text;not null" json:"alamat"`
	AnggaranID  *uuid.UUID `gorm:"type:uuid" json:"anggaran_id"`
	KategoriIDs string     `gorm:"type:text;not null;default:'[]'" json:"kategori_ids"` // JSON array
	EmailSent   bool       `gorm:"default:false" json:"email_sent"`
	EmailSentAt *time.Time `json:"email_sent_at"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Buyer    *Buyer                       `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Anggaran *FormulirPartaiBesarAnggaran `gorm:"foreignKey:AnggaranID" json:"anggaran,omitempty"`
}

func (FormulirPartaiBesarSubmission) TableName() string {
	return "formulir_partai_besar_submission"
}

// GetKategoriIDs parses JSON array to []uuid.UUID
func (s *FormulirPartaiBesarSubmission) GetKategoriIDs() []uuid.UUID {
	var ids []uuid.UUID
	var strIDs []string
	json.Unmarshal([]byte(s.KategoriIDs), &strIDs)
	for _, str := range strIDs {
		if id, err := uuid.Parse(str); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// SetKategoriIDs converts []string to JSON array
func (s *FormulirPartaiBesarSubmission) SetKategoriIDs(ids []string) {
	data, _ := json.Marshal(ids)
	s.KategoriIDs = string(data)
}
