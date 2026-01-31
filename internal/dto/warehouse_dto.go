package dto

// JadwalGudangResponse - jadwal operasional per hari
type JadwalGudangResponse struct {
	ID       string  `json:"id,omitempty"`
	Hari     int     `json:"hari"`
	NamaHari string  `json:"nama_hari"`
	JamBuka  *string `json:"jam_buka"`
	JamTutup *string `json:"jam_tutup"`
	IsBuka   bool    `json:"is_buka"`
}

// WarehouseResponse - untuk admin panel
type WarehouseResponse struct {
	ID             string   `json:"id"`
	Nama           string   `json:"nama"`
	Slug           string   `json:"slug"`
	Alamat         *string  `json:"alamat"`
	Kota           *string  `json:"kota"`
	KodePos        *string  `json:"kode_pos"`
	Telepon        *string  `json:"telepon"`
	Latitude       *float64 `json:"latitude"`
	Longitude      *float64 `json:"longitude"`
	JamOperasional *string  `json:"jam_operasional"`
	IsActive       bool     `json:"is_active"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// WarehousePublicResponse - untuk public (simplified)
type WarehousePublicResponse struct {
	Nama      string   `json:"nama"`
	Alamat    *string  `json:"alamat"`
	Kota      *string  `json:"kota"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

// InformasiPickupResponse - untuk public (warehouse + jadwal)
type InformasiPickupResponse struct {
	Alamat         *string  `json:"alamat"`
	JamOperasional *string  `json:"jam_operasional"`
	Telepon        *string  `json:"telepon"`
	WhatsappURL    string   `json:"whatsapp_url"`
	Latitude       *float64 `json:"latitude"`
	Longitude      *float64 `json:"longitude"`
	// GoogleMapsURL  string                 `json:"google_maps_url"`
	IsOpenNow bool `json:"is_open_now"`
	// StatusText     string                 `json:"status_text"`
	JadwalHariIni *JadwalGudangResponse  `json:"jadwal_hari_ini"`
	Jadwal        []JadwalGudangResponse `json:"jadwal"`
}

// WarehouseUpdateRequest - request untuk update warehouse
type WarehouseUpdateRequest struct {
	Nama           string   `json:"nama" binding:"required,min=3,max=100"`
	Alamat         *string  `json:"alamat" binding:"omitempty,min=10"`
	Kota           *string  `json:"kota" binding:"omitempty,max=100"`
	KodePos        *string  `json:"kode_pos" binding:"omitempty,max=10,numeric"`
	Telepon        *string  `json:"telepon" binding:"omitempty,max=20"`
	Latitude       *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude      *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	JamOperasional *string  `json:"jam_operasional" binding:"omitempty,max=100"`
}

// UpdateJadwalRequest - request untuk update jadwal
type UpdateJadwalRequest struct {
	Jadwal []JadwalItem `json:"jadwal" binding:"required,len=7"`
}

type JadwalItem struct {
	Hari     int     `json:"hari" binding:"required,min=0,max=6"`
	JamBuka  *string `json:"jam_buka"`
	JamTutup *string `json:"jam_tutup"`
	IsBuka   bool    `json:"is_buka"`
}
