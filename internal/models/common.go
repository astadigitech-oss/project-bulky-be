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
