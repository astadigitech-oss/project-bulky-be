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

// ========================================
// Kategori Produk Response
// ========================================

type KategoriProdukResponse struct {
	ID                      string             `json:"id"`
	Nama                    TranslatableString `json:"nama"`
	Slug                    string             `json:"slug"`
	Deskripsi               *string            `json:"deskripsi,omitempty"`
	IconURL                 *string            `json:"icon_url,omitempty"`
	MemilikiKondisiTambahan bool               `json:"memiliki_kondisi_tambahan"`
	TipeKondisiTambahan     *string            `json:"tipe_kondisi_tambahan,omitempty"`
	GambarKondisiURL        *string            `json:"gambar_kondisi_url,omitempty"`
	TeksKondisi             *string            `json:"teks_kondisi,omitempty"`
	IsActive                bool               `json:"is_active"`
	// JumlahProduk            int64     `json:"jumlah_produk"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KategoriProdukSimpleResponse struct {
	ID   string             `json:"id"`
	Nama TranslatableString `json:"nama"`
	// Slug                    string             `json:"slug"`
	// Deskripsi               *string            `json:"deskripsi,omitempty"`
	IconURL                 *string `json:"icon_url,omitempty"`
	MemilikiKondisiTambahan bool    `json:"memiliki_kondisi_tambahan"`
	// TipeKondisiTambahan     *string            `json:"tipe_kondisi_tambahan,omitempty"`
	// GambarKondisiURL        *string            `json:"gambar_kondisi_url,omitempty"`
	// TeksKondisi             *string            `json:"teks_kondisi,omitempty"`
	IsActive bool `json:"is_active"`
	// JumlahProduk            int64     `json:"jumlah_produk"`
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

// ========================================
// Produk Response
// ========================================

type ProdukRelationInfo struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
	Slug string `json:"slug"`
}

type ProdukWarehouseInfo struct {
	ID   string  `json:"id"`
	Nama string  `json:"nama"`
	Slug string  `json:"slug"`
	Kota *string `json:"kota"`
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
	ID                 string              `json:"id"`
	Nama               string              `json:"nama"`
	Slug               string              `json:"slug"`
	IDCargo            *string             `json:"id_cargo"`
	Kategori           ProdukRelationInfo  `json:"kategori"`
	Merek              *ProdukRelationInfo `json:"merek"`
	Kondisi            ProdukRelationInfo  `json:"kondisi"`
	KondisiPaket       ProdukRelationInfo  `json:"kondisi_paket"`
	Sumber             *ProdukRelationInfo `json:"sumber"`
	Warehouse          ProdukRelationInfo  `json:"warehouse"`
	TipeProduk         ProdukRelationInfo  `json:"tipe_produk"`
	HargaSebelumDiskon float64             `json:"harga_sebelum_diskon"`
	PersentaseDiskon   float64             `json:"persentase_diskon"`
	HargaSesudahDiskon float64             `json:"harga_sesudah_diskon"`
	Quantity           int                 `json:"quantity"`
	QuantityTerjual    int                 `json:"quantity_terjual"`
	GambarUtama        *string             `json:"gambar_utama"`
	IsActive           bool                `json:"is_active"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}

type ProdukDetailResponse struct {
	ID                 string                  `json:"id"`
	Nama               string                  `json:"nama"`
	Slug               string                  `json:"slug"`
	IDCargo            *string                 `json:"id_cargo"`
	Kategori           ProdukRelationInfo      `json:"kategori"`
	Merek              *ProdukRelationInfo     `json:"merek"`
	Kondisi            ProdukRelationInfo      `json:"kondisi"`
	KondisiPaket       ProdukRelationInfo      `json:"kondisi_paket"`
	Sumber             *ProdukRelationInfo     `json:"sumber"`
	Warehouse          ProdukWarehouseInfo     `json:"warehouse"`
	TipeProduk         ProdukRelationInfo      `json:"tipe_produk"`
	HargaSebelumDiskon float64                 `json:"harga_sebelum_diskon"`
	PersentaseDiskon   float64                 `json:"persentase_diskon"`
	HargaSesudahDiskon float64                 `json:"harga_sesudah_diskon"`
	Quantity           int                     `json:"quantity"`
	QuantityTerjual    int                     `json:"quantity_terjual"`
	Discrepancy        *string                 `json:"discrepancy"`
	Gambar             []ProdukGambarResponse  `json:"gambar"`
	Dokumen            []ProdukDokumenResponse `json:"dokumen"`
	IsActive           bool                    `json:"is_active"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
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
	ID        string            `json:"id"`
	Nama      string            `json:"nama"`
	GambarURL TranslatableImage `json:"gambar_url"`
	LinkURL   *string           `json:"link_url,omitempty"`
	Urutan    int               `json:"urutan"`
	IsActive  bool              `json:"is_active"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type HeroSectionSimpleResponse struct {
	ID        string            `json:"id"`
	Nama      string            `json:"nama"`
	GambarURL TranslatableImage `json:"gambar_url"`
	LinkURL   *string           `json:"link_url,omitempty"`
	Urutan    int               `json:"urutan"`
	IsActive  bool              `json:"is_active"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Public response (minimal data)
type HeroSectionPublicResponse struct {
	ID        string            `json:"id"`
	Nama      string            `json:"nama"`
	GambarURL TranslatableImage `json:"gambar_url"`
	LinkURL   *string           `json:"link_url,omitempty"`
}

// ========================================
// Banner Event Promo Response
// ========================================

type BannerEventPromoResponse struct {
	ID           string            `json:"id"`
	Nama         string            `json:"nama"`
	GambarURL    TranslatableImage `json:"gambar_url"`
	LinkURL      *string           `json:"link_url,omitempty"`
	Urutan       int               `json:"urutan"`
	IsActive     bool              `json:"is_active"`
	TanggalMulai *time.Time        `json:"tanggal_mulai,omitempty"`
	TanggalAkhir *time.Time        `json:"tanggal_akhir,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type BannerEventPromoSimpleResponse struct {
	ID        string            `json:"id"`
	Nama      string            `json:"nama"`
	GambarURL TranslatableImage `json:"gambar_url"`
	LinkURL   *string           `json:"link_url,omitempty"`
	Urutan    int               `json:"urutan"`
	IsActive  bool              `json:"is_active"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Public response (minimal data)
type BannerEventPromoPublicResponse struct {
	ID        string            `json:"id"`
	Nama      string            `json:"nama"`
	GambarURL TranslatableImage `json:"gambar_url"`
	LinkURL   *string           `json:"link_url,omitempty"`
}
