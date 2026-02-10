package models

import "time"

// ========================================
// Pagination Meta
// ========================================

type PaginationMeta struct {
	FirstPage   int   `json:"first_page"`
	LastPage    int   `json:"last_page"`
	CurrentPage int   `json:"current_page"`
	From        int   `json:"from"`
	Last        int   `json:"last"`
	Total       int64 `json:"total"`
	TotalItems  int64 `json:"total_items"` // Same as Total, for frontend convenience
	PerPage     int   `json:"per_page"`
}

// NewPaginationMeta creates a new pagination meta with calculated fields
func NewPaginationMeta(currentPage, perPage int, total int64) PaginationMeta {
	// Calculate last_page
	lastPage := int(float64(total) / float64(perPage))
	if total%int64(perPage) != 0 {
		lastPage++
	}
	if lastPage < 1 {
		lastPage = 1
	}

	// Calculate from
	from := ((currentPage - 1) * perPage) + 1
	if total == 0 {
		from = 0
	}

	// Calculate last
	last := from + perPage - 1
	if int64(last) > total {
		last = int(total)
	}
	if total == 0 {
		last = 0
	}

	return PaginationMeta{
		FirstPage:   1,
		LastPage:    lastPage,
		CurrentPage: currentPage,
		From:        from,
		Last:        last,
		Total:       total,
		TotalItems:  total, // Same as Total
		PerPage:     perPage,
	}
}

// ========================================
// Field Error (untuk validasi)
// ========================================

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ========================================
// Toggle Status Response
// ========================================

type ToggleStatusResponse struct {
	ID       string `json:"id"`
	IsActive bool   `json:"is_active"`
}

type ToggleDefaultResponse struct {
	ID        string `json:"id"`
	IsDefault bool   `json:"is_default"`
}

// ========================================
// Kategori Produk Response
// ========================================

type KategoriProdukResponse struct {
	ID                  string             `json:"id"`
	Nama                TranslatableString `json:"nama"`
	Slug                string             `json:"slug"`
	Deskripsi           string             `json:"deskripsi"`
	IconURL             string             `json:"icon_url"`
	TipeKondisiTambahan *string            `json:"tipe_kondisi_tambahan"`
	GambarKondisiURL    string             `json:"gambar_kondisi_url"`
	TeksKondisi         string             `json:"teks_kondisi"`
	IsActive            bool               `json:"is_active"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
}

type KategoriProdukSimpleResponse struct {
	ID      string             `json:"id"`
	Nama    TranslatableString `json:"nama"`
	IconURL string             `json:"icon_url"`
	// TipeKondisiTambahan *string            `json:"tipe_kondisi_tambahan"`
	IsActive bool `json:"is_active"`
	// CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Simple response untuk GET All Data (array)
//

// ========================================
// Merek Produk Response
// ========================================

type MerekProdukResponse struct {
	ID       string             `json:"id"`
	Nama     TranslatableString `json:"nama"`
	Slug     string             `json:"slug"`
	LogoURL  *string            `json:"logo_url,omitempty"`
	IsActive bool               `json:"is_active"`
	// JumlahProduk int64     `json:"jumlah_produk"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MerekProdukSimpleResponse struct {
	ID   string             `json:"id"`
	Nama TranslatableString `json:"nama"`
	// Slug     string  `json:"slug"`
	LogoURL  *string `json:"logo_url,omitempty"`
	IsActive bool    `json:"is_active"`
	// JumlahProduk int64     `json:"jumlah_produk"`
	// CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ========================================
// Kondisi Produk Response
// ========================================

type KondisiProdukResponse struct {
	ID        string             `json:"id"`
	Nama      TranslatableString `json:"nama"`
	Slug      string             `json:"slug"`
	Deskripsi *string            `json:"deskripsi,omitempty"`
	Urutan    int                `json:"urutan"`
	IsActive  bool               `json:"is_active"`
	// JumlahProduk int64     `json:"jumlah_produk"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KondisiProdukSimpleResponse struct {
	ID   string             `json:"id"`
	Nama TranslatableString `json:"nama"`
	// Slug      string  `json:"slug"`
	// Deskripsi *string `json:"deskripsi"`
	Urutan   int  `json:"urutan"`
	IsActive bool `json:"is_active"`
	// JumlahProduk int64     `json:"jumlah_produk"`
	// CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ========================================
// Kondisi Paket Response
// ========================================

type KondisiPaketResponse struct {
	ID        string             `json:"id"`
	Nama      TranslatableString `json:"nama"`
	Slug      string             `json:"slug"`
	Deskripsi *string            `json:"deskripsi,omitempty"`
	Urutan    int                `json:"urutan"`
	IsActive  bool               `json:"is_active"`
	// JumlahProduk int64     `json:"jumlah_produk"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KondisiPaketSimpleResponse struct {
	ID   string             `json:"id"`
	Nama TranslatableString `json:"nama"`
	// Slug         string    `json:"slug"`
	// Deskripsi    *string   `json:"deskripsi"`
	Urutan   int  `json:"urutan"`
	IsActive bool `json:"is_active"`
	// JumlahProduk int64     `json:"jumlah_produk"`
	// CreatedAt    time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ========================================
// Sumber Produk Response
// ========================================

type SumberProdukResponse struct {
	ID        string             `json:"id"`
	Nama      TranslatableString `json:"nama"`
	Slug      string             `json:"slug"`
	Deskripsi *string            `json:"deskripsi,omitempty"`
	IsActive  bool               `json:"is_active"`
	// JumlahProduk int64     `json:"jumlah_produk"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SumberProdukSimpleResponse struct {
	ID   string             `json:"id"`
	Nama TranslatableString `json:"nama"`
	// Slug      string  `json:"slug"`
	// Deskripsi *string `json:"deskripsi"`
	IsActive bool `json:"is_active"`
	// JumlahProduk int64     `json:"jumlah_produk"`
	// CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ========================================
// Master Dropdown Response
// ========================================

type DropdownItem struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
	Slug string `json:"slug"`
}

type MasterDropdownResponse struct {
	KategoriProduk []DropdownItem `json:"kategori_produk"`
	MerekProduk    []DropdownItem `json:"merek_produk"`
	KondisiProduk  []DropdownItem `json:"kondisi_produk"`
	KondisiPaket   []DropdownItem `json:"kondisi_paket"`
	SumberProduk   []DropdownItem `json:"sumber_produk"`
}

// ========================================
// Warehouse Response
// ========================================

type WarehouseResponse struct {
	ID           string    `json:"id"`
	Nama         string    `json:"nama"`
	Slug         string    `json:"slug"`
	Alamat       *string   `json:"alamat"`
	Kota         *string   `json:"kota"`
	KodePos      *string   `json:"kode_pos"`
	Telepon      *string   `json:"telepon"`
	IsActive     bool      `json:"is_active"`
	JumlahProduk int64     `json:"jumlah_produk"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ========================================
// Tipe Produk Response
// ========================================

type TipeProdukResponse struct {
	ID           string    `json:"id"`
	Nama         string    `json:"nama"`
	Slug         string    `json:"slug"`
	Deskripsi    *string   `json:"deskripsi"`
	Urutan       int       `json:"urutan"`
	IsActive     bool      `json:"is_active"`
	JumlahProduk int64     `json:"jumlah_produk"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TipeProdukSimpleResponse struct {
	ID           string `json:"id"`
	Nama         string `json:"nama"`
	Slug         string `json:"slug"`
	JumlahProduk int64  `json:"jumlah_produk"`
}

// ========================================
// Diskon Kategori Response
// ========================================

type DiskonKategoriKategoriInfo struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
	Slug string `json:"slug"`
}

type DiskonKategoriResponse struct {
	ID               string                     `json:"id"`
	Kategori         DiskonKategoriKategoriInfo `json:"kategori"`
	PersentaseDiskon float64                    `json:"persentase_diskon"`
	NominalDiskon    float64                    `json:"nominal_diskon"`
	TanggalMulai     *string                    `json:"tanggal_mulai"`
	TanggalSelesai   *string                    `json:"tanggal_selesai"`
	IsActive         bool                       `json:"is_active"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
}

type DiskonKategoriActiveResponse struct {
	KategoriID       string  `json:"kategori_id"`
	PersentaseDiskon float64 `json:"persentase_diskon"`
	NominalDiskon    float64 `json:"nominal_diskon"`
	TanggalMulai     *string `json:"tanggal_mulai"`
	TanggalSelesai   *string `json:"tanggal_selesai"`
}

// ========================================
// Banner Tipe Produk Response
// ========================================

type BannerTipeProdukTipeInfo struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
	Slug string `json:"slug"`
}

type BannerTipeProdukSimpleInfo struct {
	// ID   string `json:"id"`
	Nama string `json:"nama"`
	// Slug string `json:"slug"`
}

type BannerTipeProdukResponse struct {
	ID         string                   `json:"id"`
	TipeProduk BannerTipeProdukTipeInfo `json:"tipe_produk"`
	Nama       string                   `json:"nama"`
	GambarURL  string                   `json:"gambar_url"`
	Urutan     int                      `json:"urutan"`
	IsActive   bool                     `json:"is_active"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
}

type BannerTipeProdukSimpleResponse struct {
	ID         string                     `json:"id"`
	TipeProduk BannerTipeProdukSimpleInfo `json:"tipe_produk"`
	Nama       string                     `json:"nama"`
	GambarURL  string                     `json:"gambar_url"`
	Urutan     int                        `json:"urutan"`
	IsActive   bool                       `json:"is_active"`
	// CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BannerSimpleResponse struct {
	ID        string `json:"id"`
	Nama      string `json:"nama"`
	GambarURL string `json:"gambar_url"`
	Urutan    int    `json:"urutan"`
}

// Grouped response for Banner Tipe Produk
type BannerTipeProdukGroupedResponse struct {
	PaletLoad     []BannerTipeProdukSimpleResponse `json:"palet_load"`
	ContainerLoad []BannerTipeProdukSimpleResponse `json:"container_load"`
	TruckLoad     []BannerTipeProdukSimpleResponse `json:"truck_load"`
}

type BannerTipeProdukGroupedMeta struct {
	TotalByType map[string]int `json:"total_by_type"`
}

// ========================================
// Produk Response
// ========================================

type ProdukRelationInfo struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
	Slug string `json:"slug"`
}

type SimpleProdukRelationInfo struct {
	// ID   string `json:"id"`
	Nama string `json:"nama"`
	// Slug string `json:"slug"`
}

type ProdukWarehouseInfo struct {
	// ID   string  `json:"id"`
	Nama string `json:"nama"`
	// Slug string  `json:"slug"`
	Kota *string `json:"kota"`
}

type SimpleProdukWarehouseInfo struct {
	// ID   string  `json:"id"`
	Nama string `json:"nama"`
	// Slug string  `json:"slug"`
	// Kota *string `json:"kota"`
}

type ProdukGambarResponse struct {
	ID        string `json:"id"`
	GambarURL string `json:"gambar_url"`
	Urutan    int    `json:"urutan"`
	IsPrimary bool   `json:"is_primary"`
}

type ProdukDokumenResponse struct {
	ID          string `json:"id"`
	NamaDokumen string `json:"nama_dokumen"`
	FileURL     string `json:"file_url"`
	TipeFile    string `json:"tipe_file"`
	UkuranFile  *int   `json:"ukuran_file"`
}

type ProdukListResponse struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
	// Slug               string                    `json:"slug"`
	// IDCargo            *string                   `json:"id_cargo"`
	Kategori SimpleProdukRelationInfo  `json:"kategori"`
	Merek    *SimpleProdukRelationInfo `json:"merek"`
	Kondisi  SimpleProdukRelationInfo  `json:"kondisi"`
	// KondisiPaket       SimpleProdukRelationInfo  `json:"kondisi_paket"`
	// Sumber             *SimpleProdukRelationInfo `json:"sumber"`
	Warehouse  SimpleProdukWarehouseInfo `json:"warehouse"`
	TipeProduk SimpleProdukRelationInfo  `json:"tipe_produk"`
	// HargaSebelumDiskon float64                   `json:"harga_sebelum_diskon"`
	// PersentaseDiskon   float64                   `json:"persentase_diskon"`
	HargaSesudahDiskon float64   `json:"harga_sesudah_diskon"`
	Quantity           int       `json:"quantity"`
	QuantityTerjual    int       `json:"quantity_terjual"`
	Berat              float64   `json:"berat"` // kg
	GambarUtama        *string   `json:"gambar_utama"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ProdukDetailResponse struct {
	ID                 string                    `json:"id"`
	Nama               string                    `json:"nama"`
	Slug               string                    `json:"slug"`
	IDCargo            *string                   `json:"id_cargo"`
	Kategori           SimpleProdukRelationInfo  `json:"kategori"`
	Merek              *SimpleProdukRelationInfo `json:"merek"`
	Kondisi            SimpleProdukRelationInfo  `json:"kondisi"`
	KondisiPaket       SimpleProdukRelationInfo  `json:"kondisi_paket"`
	Sumber             *SimpleProdukRelationInfo `json:"sumber"`
	Warehouse          ProdukWarehouseInfo       `json:"warehouse"`
	TipeProduk         SimpleProdukRelationInfo  `json:"tipe_produk"`
	HargaSebelumDiskon float64                   `json:"harga_sebelum_diskon"`
	PersentaseDiskon   float64                   `json:"persentase_diskon"`
	HargaSesudahDiskon float64                   `json:"harga_sesudah_diskon"`
	Quantity           int                       `json:"quantity"`
	QuantityTerjual    int                       `json:"quantity_terjual"`
	Discrepancy        *string                   `json:"discrepancy"`
	Panjang            float64                   `json:"panjang"`          // cm
	Lebar              float64                   `json:"lebar"`            // cm
	Tinggi             float64                   `json:"tinggi"`           // cm
	Berat              float64                   `json:"berat"`            // kg
	BeratVolumetrik    float64                   `json:"berat_volumetrik"` // kg (calculated)
	Gambar             []ProdukGambarResponse    `json:"gambar"`
	Dokumen            []ProdukDokumenResponse   `json:"dokumen"`
	IsActive           bool                      `json:"is_active"`
	CreatedAt          time.Time                 `json:"created_at"`
	UpdatedAt          time.Time                 `json:"updated_at"`
}

// ========================================
// Auth Response
// ========================================

type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    int64         `json:"expires_in"`
	Admin        AdminResponse `json:"admin"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// ========================================
// Admin Response
// ========================================

type AdminResponse struct {
	ID          string     `json:"id"`
	Nama        string     `json:"nama"`
	Email       string     `json:"email"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type AdminListResponse struct {
	ID          string     `json:"id"`
	Nama        string     `json:"nama"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type AdminSessionResponse struct {
	ID        string    `json:"id"`
	IPAddress *string   `json:"ip_address"`
	UserAgent *string   `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	IsCurrent bool      `json:"is_current"`
}

// ========================================
// Buyer Response
// ========================================

type BuyerListResponse struct {
	ID        string    `json:"id"`
	Nama      string    `json:"nama"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Telepon   *string   `json:"telepon"`
	CreatedAt time.Time `json:"created_at"`
}

type BuyerDetailResponse struct {
	ID        string                `json:"id"`
	Nama      string                `json:"nama"`
	Username  string                `json:"username"`
	Email     string                `json:"email"`
	Telepon   *string               `json:"telepon"`
	FotoURL   *string               `json:"foto_url"`
	Alamat    []AlamatBuyerResponse `json:"alamat"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

type BuyerStatistikResponse struct {
	TotalBuyer         int64          `json:"total_buyer"`
	BuyerVerified      int64          `json:"buyer_verified"`
	PersentaseBulanIni PersentaseData `json:"persentase_bulan_ini"`
	PersentaseTahunIni PersentaseData `json:"persentase_tahun_ini"`
	RegistrasiBulanIni int64          `json:"registrasi_bulan_ini"`
	RegistrasiTahunIni int64          `json:"registrasi_tahun_ini"`
}

type PersentaseData struct {
	Value    float64 `json:"value"`
	Trend    string  `json:"trend"` // "up", "down", "stable"
	Current  int64   `json:"current"`
	Previous int64   `json:"previous"`
}

type ChartData struct {
	Date time.Time `json:"date"`
	User int       `json:"user"`
}

type ChartResponse struct {
	Mode  string      `json:"mode"`
	Chart []ChartData `json:"chart"`
	Total int64       `json:"total"`
}

// ========================================
// Alamat Buyer Response (Google Maps API)
// ========================================

type AlamatBuyerResponse struct {
	ID              string    `json:"id"`
	BuyerID         string    `json:"buyer_id"`
	Label           string    `json:"label"`
	NamaPenerima    string    `json:"nama_penerima"`
	TeleponPenerima string    `json:"telepon_penerima"`
	Provinsi        string    `json:"provinsi"`
	Kota            string    `json:"kota"`
	Kecamatan       *string   `json:"kecamatan"`
	Kelurahan       *string   `json:"kelurahan"`
	KodePos         *string   `json:"kode_pos"`
	AlamatLengkap   string    `json:"alamat_lengkap"`
	AlamatFormatted string    `json:"alamat_formatted"`
	Catatan         *string   `json:"catatan"`
	Latitude        *float64  `json:"latitude"`
	Longitude       *float64  `json:"longitude"`
	GooglePlaceID   *string   `json:"google_place_id"`
	IsDefault       bool      `json:"is_default"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ========================================
// Hero Section Response
// ========================================

type HeroSectionResponse struct {
	ID             string            `json:"id"`
	Nama           string            `json:"nama"`
	GambarURL      TranslatableImage `json:"gambar_url"`
	IsDefault      bool              `json:"is_default"`
	IsVisible      bool              `json:"is_visible"`
	TanggalMulai   *time.Time        `json:"tanggal_mulai"`
	TanggalSelesai *time.Time        `json:"tanggal_selesai"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type HeroSectionSimpleResponse struct {
	ID             string            `json:"id"`
	Nama           string            `json:"nama"`
	GambarURL      TranslatableImage `json:"gambar_url"`
	IsDefault      bool              `json:"is_default"`
	IsVisible      bool              `json:"is_visible"`
	TanggalMulai   *time.Time        `json:"tanggal_mulai"`
	TanggalSelesai *time.Time        `json:"tanggal_selesai"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// Public response (minimal data)
type HeroSectionPublicResponse struct {
	ID        string            `json:"id"`
	Nama      string            `json:"nama"`
	GambarURL TranslatableImage `json:"gambar_url"`
}

// ========================================
// Banner Event Promo Response
// ========================================

type BannerEventPromoResponse struct {
	ID             string            `json:"id"`
	Nama           string            `json:"nama"`
	GambarURL      TranslatableImage `json:"gambar_url"`
	Tujuan         []string          `json:"tujuan"` // Array of kategori ID strings
	Urutan         int               `json:"urutan"`
	IsVisible      bool              `json:"is_visible"` // Computed from tanggal
	TanggalMulai   *time.Time        `json:"tanggal_mulai"`
	TanggalSelesai *time.Time        `json:"tanggal_selesai"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type BannerEventPromoSimpleResponse struct {
	ID        string            `json:"id"`
	Nama      string            `json:"nama"`
	GambarURL TranslatableImage `json:"gambar_url"`
	Tujuan    []string          `json:"tujuan"` // Array of kategori ID strings
	Urutan    int               `json:"urutan"`
	IsVisible bool              `json:"is_visible"` // Computed from tanggal
	UpdatedAt time.Time         `json:"updated_at"`
}

// Public response (minimal data)
type BannerEventPromoPublicResponse struct {
	ID        string            `json:"id"`
	Nama      string            `json:"nama"`
	GambarURL TranslatableImage `json:"gambar_url"`
	Tujuan    []string          `json:"tujuan"` // Array of kategori ID strings
}

// ========================================
// Formulir Partai Besar - Config Response
// ========================================

type FormulirConfigResponse struct {
	ID          string    `json:"id"`
	DaftarEmail []string  `json:"daftar_email"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ========================================
// Formulir Partai Besar - Anggaran Response
// ========================================

type AnggaranResponse struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Urutan int    `json:"urutan"`
}

// ========================================
// Formulir Partai Besar - Submission Response (Admin)
// ========================================

type FormulirSubmissionListResponse struct {
	ID        string    `json:"id"`
	Nama      string    `json:"nama"`
	Telepon   string    `json:"telepon"`
	Alamat    string    `json:"alamat"`
	Anggaran  string    `json:"anggaran"`
	Kategori  []string  `json:"kategori"` // Array nama kategori
	EmailSent bool      `json:"email_sent"`
	CreatedAt time.Time `json:"created_at"`
}

type FormulirSubmissionDetailResponse struct {
	ID          string             `json:"id"`
	Buyer       *BuyerListResponse `json:"buyer,omitempty"`
	Nama        string             `json:"nama"`
	Telepon     string             `json:"telepon"`
	Alamat      string             `json:"alamat"`
	Anggaran    *AnggaranResponse  `json:"anggaran"`
	Kategori    []DropdownItem     `json:"kategori"`
	EmailSent   bool               `json:"email_sent"`
	EmailSentAt *time.Time         `json:"email_sent_at"`
	CreatedAt   time.Time          `json:"created_at"`
}

type FormulirOptionsResponse struct {
	Anggaran []AnggaranResponse `json:"anggaran"`
	Kategori []DropdownItem     `json:"kategori"`
}

// ========================================
// WhatsApp Handler Response
// ========================================

type WhatsAppHandlerResponse struct {
	ID          string    `json:"id"`
	NomorWA     string    `json:"nomor_wa"`
	PesanAwal   string    `json:"pesan_awal"`
	IsActive    bool      `json:"is_active"`
	WhatsAppURL string    `json:"whatsapp_url"` // Generated URL
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WhatsAppHandlerSimpleResponse struct {
	ID        string `json:"id"`
	NomorWA   string `json:"nomor_wa"`
	PesanAwal string `json:"pesan_awal"`
	// IsActive  bool   `json:"is_active"`
	WhatsAppURL string `json:"whatsapp_url"` // Generated URL
	// CreatedAt   time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Public response (untuk floating button)
type WhatsAppHandlerPublicResponse struct {
	NomorWA     string `json:"nomor_wa"`
	PesanAwal   string `json:"pesan_awal"`
	WhatsAppURL string `json:"whatsapp_url"`
}
