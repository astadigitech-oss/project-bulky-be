package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BuyerOAuth struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BuyerID     uuid.UUID      `gorm:"type:uuid;not null;index"                         json:"buyer_id"`
	Provider    string         `gorm:"type:varchar(50);not null"                        json:"provider"`
	ProviderUID string         `gorm:"type:varchar(255);not null"                       json:"provider_uid"`
	Email       *string        `gorm:"type:varchar(255)"                                json:"email"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"                                   json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"                                   json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                                            json:"-"`

	// Relations
	Buyer Buyer `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
}

func (BuyerOAuth) TableName() string {
	return "buyer_oauth"
}
