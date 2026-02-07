package models

import (
	"database/sql/driver"
	"encoding/json"
)

// DualLanguage untuk field JSONB dual bahasa
type DualLanguage map[string]interface{}

// NewDualLanguage creates DualLanguage from map[string]string
func NewDualLanguage(m map[string]string) DualLanguage {
	dl := make(DualLanguage)
	for k, v := range m {
		dl[k] = v
	}
	return dl
}

// Scan implements sql.Scanner interface
func (dl *DualLanguage) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, dl)
}

// Value implements driver.Valuer interface
func (dl DualLanguage) Value() (driver.Value, error) {
	return json.Marshal(dl)
}

// TranslatableString untuk field dual bahasa (nama)
type TranslatableString struct {
	ID string  `json:"id"`
	EN *string `json:"en,omitempty"`
}

// TranslatableImage untuk field dual bahasa (gambar)
type TranslatableImage struct {
	ID string  `json:"id"`
	EN *string `json:"en,omitempty"`
}

// GetFullURL returns TranslatableImage with full URLs
func (t TranslatableImage) GetFullURL(baseURL string) TranslatableImage {
	if baseURL == "" {
		return t
	}

	result := TranslatableImage{
		ID: t.ID,
		EN: t.EN,
	}

	if t.ID != "" {
		result.ID = baseURL + "/uploads/" + t.ID
	}

	if t.EN != nil && *t.EN != "" {
		fullEN := baseURL + "/uploads/" + *t.EN
		result.EN = &fullEN
	}

	return result
}

// SafeString converts *string to string (empty if nil)
func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
