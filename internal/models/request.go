package models

import "time"

// ========================================
// Pagination Request
// ========================================

type PaginationRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PerPage  int    `form:"per_page" binding:"min=1,max=100"`
	Search   string `form:"search"`
	SortBy   string `form:"sort_by"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
	IsActive *bool  `form:"is_active"`
}

func (p *PaginationRequest) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = 10
	}
	if p.PerPage > 100 {
		p.PerPage = 100
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.Order == "" {
		p.Order = "desc"
	}
}

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

// ========================================
// Kategori Produk Request
// ========================================

type CreateKategoriProdukRequest struct {
	NamaID              string  `json:"nama_id" binding:"required,min=2,max=100"`
	NamaEN              *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Deskripsi           *string `json:"deskripsi" binding:"omitempty,max=1000"`
	Icon                *string `json:"icon"`
	TipeKondisiTambahan *string `json:"tipe_kondisi_tambahan" binding:"omitempty,oneof=TEKS GAMBAR"`
	GambarKondisi       *string `json:"gambar_kondisi"`
	TeksKondisi         *string `json:"teks_kondisi" binding:"omitempty,max=500"`
}

type UpdateKategoriProdukRequest struct {
	NamaID              *string `json:"nama_id" binding:"omitempty,min=2,max=100"`
	NamaEN              *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Deskripsi           *string `json:"deskripsi" binding:"omitempty,max=1000"`
	Icon                *string `json:"icon"`
	TipeKondisiTambahan *string `json:"tipe_kondisi_tambahan" binding:"omitempty,oneof=TEKS GAMBAR"`
	GambarKondisi       *string `json:"gambar_kondisi"`
	TeksKondisi         *string `json:"teks_kondisi" binding:"omitempty,max=500"`
	IsActive            *bool   `json:"is_active"`
}

type KategoriProdukFilterRequest struct {
	PaginationRequest
	IsActive  *bool   `form:"is_active"`
	UpdatedAt *string `form:"updated_at"`
}

// ========================================
// Merek Produk Request
// ========================================

type CreateMerekProdukRequest struct {
	NamaID string  `json:"nama_id" binding:"required,min=2,max=100"`
	NamaEN *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Logo   *string `json:"logo"`
}

type UpdateMerekProdukRequest struct {
	NamaID   *string `json:"nama_id" binding:"omitempty,min=2,max=100"`
	NamaEN   *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Logo     *string `json:"logo"`
	IsActive *bool   `json:"is_active"`
}

type MerekProdukFilterRequest struct {
	PaginationRequest
	IsActive  *bool   `form:"is_active"`
	UpdatedAt *string `form:"updated_at"`
}

// ========================================
// Kondisi Produk Request
// ========================================

type CreateKondisiProdukRequest struct {
	NamaID    string  `json:"nama_id" binding:"required,min=2,max=100"`
	NamaEN    *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
}

type UpdateKondisiProdukRequest struct {
	NamaID    *string `json:"nama_id" binding:"omitempty,min=2,max=100"`
	NamaEN    *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
	IsActive  *bool   `json:"is_active"`
}

type KondisiProdukFilterRequest struct {
	PaginationRequest
	IsActive  *bool   `form:"is_active"`
	UpdatedAt *string `form:"updated_at"`
}

// ========================================
// Kondisi Paket Request
// ========================================

type CreateKondisiPaketRequest struct {
	NamaID    string  `json:"nama_id" binding:"required,min=2,max=100"`
	NamaEN    *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
}

type UpdateKondisiPaketRequest struct {
	NamaID    *string `json:"nama_id" binding:"omitempty,min=2,max=100"`
	NamaEN    *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
	IsActive  *bool   `json:"is_active"`
}

// ========================================
// Sumber Produk Request
// ========================================

type CreateSumberProdukRequest struct {
	NamaID    string  `json:"nama_id" binding:"required,min=2,max=100"`
	NamaEN    *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
}

type UpdateSumberProdukRequest struct {
	NamaID    *string `json:"nama_id" binding:"omitempty,min=2,max=100"`
	NamaEN    *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
	IsActive  *bool   `json:"is_active"`
}

type SumberProdukFilterRequest struct {
	PaginationRequest
	IsActive  *bool   `form:"is_active"`
	UpdatedAt *string `form:"updated_at"`
}

// ========================================
// Reorder Request
// ========================================

type ReorderItem struct {
	ID     string `json:"id" binding:"required,uuid"`
	Urutan int    `json:"urutan" binding:"min=0"`
}

type ReorderRequest struct {
	Items []ReorderItem `json:"items" binding:"required,min=1,dive"`
}

// ReorderByDirectionRequest - For single item reorder (up/down)
type ReorderByDirectionRequest struct {
	Direction string `json:"direction" binding:"required,oneof=up down"`
}

// ========================================
// Warehouse Request
// ========================================

type CreateWarehouseRequest struct {
	Nama    string  `json:"nama" binding:"required,min=2,max=100"`
	Alamat  *string `json:"alamat" binding:"omitempty,max=500"`
	Kota    *string `json:"kota" binding:"omitempty,max=100"`
	KodePos *string `json:"kode_pos" binding:"omitempty,max=10"`
	Telepon *string `json:"telepon" binding:"omitempty,max=20"`
}

type UpdateWarehouseRequest struct {
	Nama     *string `json:"nama" binding:"omitempty,min=2,max=100"`
	Alamat   *string `json:"alamat" binding:"omitempty,max=500"`
	Kota     *string `json:"kota" binding:"omitempty,max=100"`
	KodePos  *string `json:"kode_pos" binding:"omitempty,max=10"`
	Telepon  *string `json:"telepon" binding:"omitempty,max=20"`
	IsActive *bool   `json:"is_active"`
}

// ========================================
// Tipe Produk Request
// ========================================
// Tipe Produk Request (DEPRECATED - Read Only)
// ========================================
// Note: Tipe produk is read-only (Paletbox, Container, Truckload)
// Data managed via migration only

type TipeProdukFilterRequest struct {
	PaginationRequest
	JumlahProduk *int `form:"jumlah_produk"`
}

// Deprecated: CreateTipeProdukRequest is no longer used - tipe produk is read-only
type CreateTipeProdukRequest struct {
	Nama      string  `json:"nama" binding:"required,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
	Urutan    *int    `json:"urutan"`
}

// Deprecated: UpdateTipeProdukRequest is no longer used - tipe produk is read-only
type UpdateTipeProdukRequest struct {
	Nama      *string `json:"nama" binding:"omitempty,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
	Urutan    *int    `json:"urutan"`
	IsActive  *bool   `json:"is_active"`
}

// ========================================
// Diskon Kategori Request
// ========================================

type CreateDiskonKategoriRequest struct {
	KategoriID       string  `json:"kategori_id" binding:"required,uuid"`
	PersentaseDiskon float64 `json:"persentase_diskon" binding:"required,min=0,max=100"`
	NominalDiskon    float64 `json:"nominal_diskon" binding:"min=0"`
	TanggalMulai     *string `json:"tanggal_mulai"`
	TanggalSelesai   *string `json:"tanggal_selesai"`
}

type UpdateDiskonKategoriRequest struct {
	KategoriID       *string  `json:"kategori_id" binding:"omitempty,uuid"`
	PersentaseDiskon *float64 `json:"persentase_diskon" binding:"omitempty,min=0,max=100"`
	NominalDiskon    *float64 `json:"nominal_diskon" binding:"omitempty,min=0"`
	TanggalMulai     *string  `json:"tanggal_mulai"`
	TanggalSelesai   *string  `json:"tanggal_selesai"`
	IsActive         *bool    `json:"is_active"`
}

// ========================================
// Banner Tipe Produk Request
// ========================================

type CreateBannerTipeProdukRequest struct {
	TipeProdukID string `json:"tipe_produk_id" binding:"required,uuid"`
	Nama         string `json:"nama" binding:"required,min=2,max=100"`
	GambarURL    string `json:"gambar_url" binding:"required,max=500"`
}

type UpdateBannerTipeProdukRequest struct {
	TipeProdukID *string `json:"tipe_produk_id" binding:"omitempty,uuid"`
	Nama         *string `json:"nama" binding:"omitempty,min=2,max=100"`
	GambarURL    *string `json:"gambar_url" binding:"omitempty,max=500"`
	IsActive     *bool   `json:"is_active"`
}

type BannerTipeProdukFilterRequest struct {
	PaginationRequest
	// TipeProdukID *string `form:"tipe_produk_id"`
	IsActive  *bool   `form:"is_active"`
	UpdatedAt *string `form:"updated_at"`
}

// ========================================
// Produk Request
// ========================================

type CreateProdukRequest struct {
	Nama               string  `form:"nama" binding:"required,min=2,max=255"`
	IDCargo            *string `form:"id_cargo" binding:"omitempty,max=50"`
	KategoriID         string  `form:"kategori_id" binding:"required,uuid"`
	MerekID            *string `form:"merek_id" binding:"omitempty,uuid"`
	KondisiID          string  `form:"kondisi_id" binding:"required,uuid"`
	KondisiPaketID     string  `form:"kondisi_paket_id" binding:"required,uuid"`
	SumberID           *string `form:"sumber_id" binding:"omitempty,uuid"`
	WarehouseID        string  `form:"warehouse_id" binding:"required,uuid"`
	TipeProdukID       string  `form:"tipe_produk_id" binding:"required,uuid"`
	HargaSebelumDiskon float64 `form:"harga_sebelum_diskon" binding:"required,gt=0"`
	PersentaseDiskon   float64 `form:"persentase_diskon" binding:"min=0,max=100"`
	HargaSesudahDiskon float64 `form:"harga_sesudah_diskon" binding:"required,min=0"`
	Quantity           int     `form:"quantity" binding:"min=0"`
	Discrepancy        *string `form:"discrepancy" binding:"omitempty,max=1000"`
	Panjang            float64 `form:"panjang" binding:"required,gt=0"` // cm
	Lebar              float64 `form:"lebar" binding:"required,gt=0"`   // cm
	Tinggi             float64 `form:"tinggi" binding:"required,gt=0"`  // cm
	Berat              float64 `form:"berat" binding:"required,gt=0"`   // kg
	GambarUtamaIndex   int     `form:"gambar_utama_index"`
}

type UpdateProdukRequest struct {
	Nama               *string  `form:"nama" binding:"omitempty,min=2,max=255"`
	IDCargo            *string  `form:"id_cargo" binding:"omitempty,max=50"`
	KategoriID         *string  `form:"kategori_id" binding:"omitempty,uuid"`
	MerekID            *string  `form:"merek_id" binding:"omitempty,uuid"`
	KondisiID          *string  `form:"kondisi_id" binding:"omitempty,uuid"`
	KondisiPaketID     *string  `form:"kondisi_paket_id" binding:"omitempty,uuid"`
	SumberID           *string  `form:"sumber_id" binding:"omitempty,uuid"`
	WarehouseID        *string  `form:"warehouse_id" binding:"omitempty,uuid"`
	TipeProdukID       *string  `form:"tipe_produk_id" binding:"omitempty,uuid"`
	HargaSebelumDiskon *float64 `form:"harga_sebelum_diskon" binding:"omitempty,gt=0"`
	PersentaseDiskon   *float64 `form:"persentase_diskon" binding:"omitempty,min=0,max=100"`
	HargaSesudahDiskon *float64 `form:"harga_sesudah_diskon" binding:"omitempty,min=0"`
	Quantity           *int     `form:"quantity" binding:"omitempty,min=0"`
	Discrepancy        *string  `form:"discrepancy" binding:"omitempty,max=1000"`
	Panjang            *float64 `form:"panjang" binding:"omitempty,gt=0"` // cm
	Lebar              *float64 `form:"lebar" binding:"omitempty,gt=0"`   // cm
	Tinggi             *float64 `form:"tinggi" binding:"omitempty,gt=0"`  // cm
	Berat              *float64 `form:"berat" binding:"omitempty,gt=0"`   // kg
	IsActive           *bool    `form:"is_active"`
}

type UpdateStockRequest struct {
	Quantity int     `json:"quantity" binding:"required,min=0"`
	Catatan  *string `json:"catatan"`
}

type ProdukFilterRequest struct {
	PaginationRequest
	KategoriID     string  `form:"kategori_id"`
	MerekID        string  `form:"merek_id"`
	KondisiID      string  `form:"kondisi_id"`
	KondisiPaketID string  `form:"kondisi_paket_id"`
	SumberID       string  `form:"sumber_id"`
	WarehouseID    string  `form:"warehouse_id"`
	TipeProdukID   string  `form:"tipe_produk_id"`
	HargaMin       float64 `form:"harga_min"`
	HargaMax       float64 `form:"harga_max"`
}

// ========================================
// Produk Gambar Request
// ========================================

type CreateProdukGambarRequest struct {
	IsPrimary bool `form:"is_primary"`
}

type UpdateProdukGambarRequest struct {
	IsPrimary *bool `json:"is_primary"`
}

type ReorderGambarRequest struct {
	Direction string `json:"direction" binding:"required,oneof=up down"`
}

// ========================================
// Produk Dokumen Request
// ========================================

type CreateProdukDokumenRequest struct {
	NamaDokumen string `form:"nama_dokumen" binding:"required,max=255"`
}

// ========================================
// Auth Request
// ========================================

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// ========================================
// Admin CRUD Request
// ========================================

type AdminFilterRequest struct {
	PaginationRequest
	IsActive  *bool   `form:"is_active"`
	CreatedAt *string `form:"created_at"`
}

type CreateAdminRequest struct {
	Nama            string `json:"nama" binding:"required,min=2,max=100"`
	Email           string `json:"email" binding:"required,email,max=255"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	RoleID          string `json:"role_id" binding:"required,uuid"`
}

type UpdateAdminRequest struct {
	Nama     *string `json:"nama" binding:"omitempty,min=2,max=100"`
	Email    *string `json:"email" binding:"omitempty,email,max=255"`
	IsActive *bool   `json:"is_active"`
}

type UpdateProfileRequest struct {
	Nama  *string `json:"nama" binding:"omitempty,min=2,max=100"`
	Email *string `json:"email" binding:"omitempty,email,max=255"`
}

type ResetPasswordRequest struct {
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// ========================================
// Buyer Request (Admin: RUD only, no Create)
// ========================================

type UpdateBuyerRequest struct {
	Nama     *string `json:"nama" binding:"omitempty,min=2,max=100"`
	Username *string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    *string `json:"email" binding:"omitempty,email,max=255"`
	Telepon  *string `json:"telepon" binding:"omitempty,max=20"`
	IsActive *bool   `json:"is_active"`
}

type ResetBuyerPasswordRequest struct {
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

type BuyerFilterRequest struct {
	PaginationRequest
	CreatedAt *bool `form:"created_at"`
}

type ChartParams struct {
	Filter        string    `form:"filter"`                                  // year, month, week, custom
	Tahun         int       `form:"tahun"`                                   // Year
	Bulan         int       `form:"bulan"`                                   // Month (1-12)
	Minggu        int       `form:"minggu"`                                  // Week (1-5)
	TanggalDari   time.Time `form:"tanggal_dari" time_format:"2006-01-02"`   // Start date for custom
	TanggalSampai time.Time `form:"tanggal_sampai" time_format:"2006-01-02"` // End date for custom
}

// ========================================
// Alamat Buyer Request (Google Maps API)
// ========================================

type CreateAlamatBuyerRequest struct {
	BuyerID         string   `json:"buyer_id" binding:"required,uuid"`
	Label           string   `json:"label" binding:"required,min=1,max=50"`
	NamaPenerima    string   `json:"nama_penerima" binding:"required,min=2,max=100"`
	TeleponPenerima string   `json:"telepon_penerima" binding:"required,max=20"`
	Provinsi        string   `json:"provinsi" binding:"required,max=100"`
	Kota            string   `json:"kota" binding:"required,max=100"`
	Kecamatan       *string  `json:"kecamatan" binding:"omitempty,max=100"`
	Kelurahan       *string  `json:"kelurahan" binding:"omitempty,max=100"`
	KodePos         *string  `json:"kode_pos" binding:"omitempty,max=10"`
	AlamatLengkap   string   `json:"alamat_lengkap" binding:"required,max=500"`
	Catatan         *string  `json:"catatan" binding:"omitempty,max=500"`
	Latitude        *float64 `json:"latitude"`
	Longitude       *float64 `json:"longitude"`
	GooglePlaceID   *string  `json:"google_place_id" binding:"omitempty,max=255"`
	IsDefault       bool     `json:"is_default"`
}

type UpdateAlamatBuyerRequest struct {
	Label           *string  `json:"label" binding:"omitempty,min=1,max=50"`
	NamaPenerima    *string  `json:"nama_penerima" binding:"omitempty,min=2,max=100"`
	TeleponPenerima *string  `json:"telepon_penerima" binding:"omitempty,max=20"`
	Provinsi        *string  `json:"provinsi" binding:"omitempty,max=100"`
	Kota            *string  `json:"kota" binding:"omitempty,max=100"`
	Kecamatan       *string  `json:"kecamatan" binding:"omitempty,max=100"`
	Kelurahan       *string  `json:"kelurahan" binding:"omitempty,max=100"`
	KodePos         *string  `json:"kode_pos" binding:"omitempty,max=10"`
	AlamatLengkap   *string  `json:"alamat_lengkap" binding:"omitempty,max=500"`
	Catatan         *string  `json:"catatan" binding:"omitempty,max=500"`
	Latitude        *float64 `json:"latitude"`
	Longitude       *float64 `json:"longitude"`
	GooglePlaceID   *string  `json:"google_place_id" binding:"omitempty,max=255"`
	IsDefault       *bool    `json:"is_default"`
}

type AlamatBuyerFilterRequest struct {
	PaginationRequest
	BuyerID string `form:"buyer_id" binding:"required,uuid"`
}

// ========================================
// Hero Section Request
// ========================================

type CreateHeroSectionRequest struct {
	Nama           string  `json:"nama" binding:"required,min=1,max=100"`
	GambarID       string  `json:"gambar_id" binding:"required,max=500"`
	GambarEN       *string `json:"gambar_en" binding:"omitempty,max=500"`
	IsDefault      bool    `json:"is_default"`
	TanggalMulai   *string `json:"tanggal_mulai"`
	TanggalSelesai *string `json:"tanggal_selesai"`
}

type UpdateHeroSectionRequest struct {
	Nama           *string `json:"nama" binding:"omitempty,min=1,max=100"`
	GambarID       *string `json:"gambar_id" binding:"omitempty,max=500"`
	GambarEN       *string `json:"gambar_en" binding:"omitempty,max=500"`
	TanggalMulai   *string `json:"tanggal_mulai"`
	TanggalSelesai *string `json:"tanggal_selesai"`
}

type HeroSectionFilterRequest struct {
	PaginationRequest
	IsDefault *bool   `form:"is_default"`
	UpdatedAt *string `form:"updated_at"`
}

// ========================================
// Banner Event Promo Request
// ========================================

type CreateBannerEventPromoRequest struct {
	Nama           string   `json:"nama" binding:"required,min=1,max=100"`
	GambarID       string   `json:"gambar_id" binding:"required,max=255"`
	GambarEN       string   `json:"gambar_en" binding:"required,max=255"`
	Tujuan         []string `json:"tujuan"` // Array of kategori IDs
	TanggalMulai   *string  `json:"tanggal_mulai"`
	TanggalSelesai *string  `json:"tanggal_selesai"`
}

type UpdateBannerEventPromoRequest struct {
	Nama           *string   `json:"nama" binding:"omitempty,min=1,max=100"`
	GambarID       *string   `json:"gambar_id" binding:"omitempty,max=255"`
	GambarEN       *string   `json:"gambar_en" binding:"omitempty,max=255"`
	Tujuan         *[]string `json:"tujuan"` // Array of kategori IDs
	TanggalMulai   *string   `json:"tanggal_mulai"`
	TanggalSelesai *string   `json:"tanggal_selesai"`
}

type BannerEventPromoFilterRequest struct {
	PaginationRequest
	UpdatedAt *string `form:"updated_at"`
}

// ========================================
// Formulir Partai Besar - Config Request
// ========================================

type UpdateFormulirConfigRequest struct {
	DaftarEmail []string `json:"daftar_email" binding:"required,min=1,dive,email"`
}

// ========================================
// Formulir Partai Besar - Anggaran Request
// ========================================

type CreateAnggaranRequest struct {
	Label string `json:"label" binding:"required,min=1,max=100"`
}

type UpdateAnggaranRequest struct {
	Label *string `json:"label" binding:"omitempty,min=1,max=100"`
}

type AnggaranFilterRequest struct {
	PaginationRequest
}

// ========================================
// Formulir Partai Besar - Submission Request (Buyer)
// ========================================

type CreateFormulirSubmissionRequest struct {
	Nama        string   `json:"nama" binding:"required,min=2,max=100"`
	Telepon     string   `json:"telepon" binding:"required,max=20"`
	Alamat      string   `json:"alamat" binding:"required"`
	AnggaranID  string   `json:"anggaran_id" binding:"required,uuid"`
	KategoriIDs []string `json:"kategori_ids" binding:"required,min=1,dive,uuid"`
}

type FormulirSubmissionFilterRequest struct {
	PaginationRequest
	TanggalDari   *string `form:"tanggal_dari"`
	TanggalSampai *string `form:"tanggal_sampai"`
	EmailSent     *bool   `form:"email_sent"`
}

// ========================================
// WhatsApp Handler Request
// ========================================

type CreateWhatsAppHandlerRequest struct {
	NomorWA   string `json:"nomor_wa" binding:"required,max=20"`
	PesanAwal string `json:"pesan_awal" binding:"required"`
	IsActive  bool   `json:"is_active"`
}

type UpdateWhatsAppHandlerRequest struct {
	NomorWA   *string `json:"nomor_wa" binding:"omitempty,max=20"`
	PesanAwal *string `json:"pesan_awal"`
	IsActive  *bool   `json:"is_active"`
}

type WhatsAppHandlerFilterRequest struct {
	PaginationRequest
	IsActive  *bool   `form:"is_active"`
	UpdatedAt *string `form:"updated_at"`
}
