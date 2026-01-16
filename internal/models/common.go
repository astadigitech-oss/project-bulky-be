package models

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
