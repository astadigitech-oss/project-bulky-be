package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabelBlog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	NamaID    string         `gorm:"type:varchar(100);not null" json:"nama_id"`
	NamaEN    string         `gorm:"type:varchar(100);not null" json:"nama_en"`
	Slug      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	SlugID    *string        `gorm:"type:varchar(100);uniqueIndex" json:"slug_id"`
	SlugEN    *string        `gorm:"type:varchar(100);uniqueIndex" json:"slug_en"`
	Urutan    int            `gorm:"default:0" json:"urutan"`
	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"deleted_at,omitempty"`
}

func (LabelBlog) TableName() string {
	return "label_blog"
}
