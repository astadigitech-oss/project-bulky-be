package dto

// PaymentMethodResponse - individual method
type PaymentMethodResponse struct {
	ID        string  `json:"id"`
	Nama      string  `json:"nama"`
	Kode      string  `json:"kode"`
	LogoValue *string `json:"logo_value"`
	Urutan    int     `json:"urutan,omitempty"`    // Only for admin
	IsActive  bool    `json:"is_active,omitempty"` // Only for admin
}

// PaymentMethodGroupResponse - grouped response
type PaymentMethodGroupResponse struct {
	Group    string                  `json:"group"`
	Urutan   int                     `json:"urutan,omitempty"`    // Only for admin
	IsActive bool                    `json:"is_active,omitempty"` // Only for admin
	Methods  []PaymentMethodResponse `json:"methods"`
}

// ToggleMethodStatusResponse - response for toggle method status
type ToggleMethodStatusResponse struct {
	ID       string `json:"id"`
	Nama     string `json:"nama"`
	Kode     string `json:"kode"`
	IsActive bool   `json:"is_active"`
}

// ToggleGroupStatusResponse - response for toggle group status
type ToggleGroupStatusResponse struct {
	Group    string `json:"group"`
	Urutan   int    `json:"urutan"`
	IsActive bool   `json:"is_active"`
}
