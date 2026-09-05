package models

type CreateForceUpdateRequest struct {
	KodeVersi          string `json:"kode_versi" binding:"required,max=20"`
	MinimumBuildNumber int64  `json:"minimum_build_number" binding:"required,min=1"`
	UpdateType         string `json:"update_type" binding:"required,oneof=OPTIONAL MANDATORY"`
	InformasiUpdate    string `json:"informasi_update" binding:"required"`
	InformasiUpdateEn  string `json:"informasi_update_en" binding:"required"`
	Platform           string `json:"platform" binding:"required,oneof=ANDROID IOS"`
	IsActive           bool   `json:"is_active"`
}

type UpdateForceUpdateRequest struct {
	KodeVersi          *string `json:"kode_versi" binding:"omitempty,max=20"`
	MinimumBuildNumber *int64  `json:"minimum_build_number" binding:"omitempty,min=1"`
	UpdateType         *string `json:"update_type" binding:"omitempty,oneof=OPTIONAL MANDATORY"`
	InformasiUpdate    *string `json:"informasi_update"`
	InformasiUpdateEn  *string `json:"informasi_update_en"`
	Platform           *string `json:"platform" binding:"omitempty,oneof=ANDROID IOS"`
	IsActive           *bool   `json:"is_active"`
}
