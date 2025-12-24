package models

import "time"

// ========================================
// Pagination Meta
// ========================================

type PaginationMeta struct {
	Halaman      int   `json:"halaman"`
	PerHalaman   int   `json:"per_halaman"`
	TotalData    int64 `json:"total_data"`
	TotalHalaman int64 `json:"total_halaman"`
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
	ID                      string    `json:"id"`
	Nama                    string    `json:"nama"`
	Slug                    string    `json:"slug"`
	Deskripsi               *string   `json:"deskripsi"`
	IconURL                 *string   `json:"icon_url"`
	MemilikiKondisiTambahan bool      `json:"memiliki_kondisi_tambahan"`
	TipeKondisiTambahan     *string   `json:"tipe_kondisi_tambahan"`
	GambarKondisiURL        *string   `json:"gambar_kondisi_url"`
	TeksKondisi             *string   `json:"teks_kondisi"`
	IsActive                bool      `json:"is_active"`
	JumlahProduk            int64     `json:"jumlah_produk"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// ========================================
// Merek Produk Response
// ========================================

type MerekProdukResponse struct {
	ID           string    `json:"id"`
	Nama         string    `json:"nama"`
	Slug         string    `json:"slug"`
	LogoURL      *string   `json:"logo_url"`
	IsActive     bool      `json:"is_active"`
	JumlahProduk int64     `json:"jumlah_produk"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ========================================
// Kondisi Produk Response
// ========================================

type KondisiProdukResponse struct {
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

// ========================================
// Kondisi Paket Response
// ========================================

type KondisiPaketResponse struct {
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

// ========================================
// Sumber Produk Response
// ========================================

type SumberProdukResponse struct {
	ID           string    `json:"id"`
	Nama         string    `json:"nama"`
	Slug         string    `json:"slug"`
	Deskripsi    *string   `json:"deskripsi"`
	IsActive     bool      `json:"is_active"`
	JumlahProduk int64     `json:"jumlah_produk"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
	ID                   string              `json:"id"`
	Nama                 string              `json:"nama"`
	Slug                 string              `json:"slug"`
	IDCargo              *string             `json:"id_cargo"`
	Kategori             ProdukRelationInfo  `json:"kategori"`
	Merek                *ProdukRelationInfo `json:"merek"`
	Kondisi              ProdukRelationInfo  `json:"kondisi"`
	KondisiPaket         ProdukRelationInfo  `json:"kondisi_paket"`
	Sumber               *ProdukRelationInfo `json:"sumber"`
	Warehouse            ProdukRelationInfo  `json:"warehouse"`
	TipeProduk           ProdukRelationInfo  `json:"tipe_produk"`
	HargaSebelumDiskon   float64             `json:"harga_sebelum_diskon"`
	PersentaseDiskon     float64             `json:"persentase_diskon"`
	HargaSesudahDiskon   float64             `json:"harga_sesudah_diskon"`
	Quantity             int                 `json:"quantity"`
	QuantityTerjual      int                 `json:"quantity_terjual"`
	GambarUtama          *string             `json:"gambar_utama"`
	IsActive             bool                `json:"is_active"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
}

type ProdukDetailResponse struct {
	ID                   string                  `json:"id"`
	Nama                 string                  `json:"nama"`
	Slug                 string                  `json:"slug"`
	IDCargo              *string                 `json:"id_cargo"`
	Kategori             ProdukRelationInfo      `json:"kategori"`
	Merek                *ProdukRelationInfo     `json:"merek"`
	Kondisi              ProdukRelationInfo      `json:"kondisi"`
	KondisiPaket         ProdukRelationInfo      `json:"kondisi_paket"`
	Sumber               *ProdukRelationInfo     `json:"sumber"`
	Warehouse            ProdukWarehouseInfo     `json:"warehouse"`
	TipeProduk           ProdukRelationInfo      `json:"tipe_produk"`
	HargaSebelumDiskon   float64                 `json:"harga_sebelum_diskon"`
	PersentaseDiskon     float64                 `json:"persentase_diskon"`
	HargaSesudahDiskon   float64                 `json:"harga_sesudah_diskon"`
	Quantity             int                     `json:"quantity"`
	QuantityTerjual      int                     `json:"quantity_terjual"`
	Discrepancy          *string                 `json:"discrepancy"`
	Gambar               []ProdukGambarResponse  `json:"gambar"`
	Dokumen              []ProdukDokumenResponse `json:"dokumen"`
	IsActive             bool                    `json:"is_active"`
	CreatedAt            time.Time               `json:"created_at"`
	UpdatedAt            time.Time               `json:"updated_at"`
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
	ID           string     `json:"id"`
	Nama         string     `json:"nama"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Telepon      *string    `json:"telepon"`
	IsActive     bool       `json:"is_active"`
	IsVerified   bool       `json:"is_verified"`
	JumlahAlamat int        `json:"jumlah_alamat"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type BuyerDetailResponse struct {
	ID              string                `json:"id"`
	Nama            string                `json:"nama"`
	Username        string                `json:"username"`
	Email           string                `json:"email"`
	Telepon         *string               `json:"telepon"`
	IsActive        bool                  `json:"is_active"`
	IsVerified      bool                  `json:"is_verified"`
	EmailVerifiedAt *time.Time            `json:"email_verified_at"`
	LastLoginAt     *time.Time            `json:"last_login_at"`
	Alamat          []AlamatBuyerResponse `json:"alamat"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

type BuyerStatistikResponse struct {
	TotalBuyer          int64 `json:"total_buyer"`
	BuyerAktif          int64 `json:"buyer_aktif"`
	BuyerNonaktif       int64 `json:"buyer_nonaktif"`
	BuyerVerified       int64 `json:"buyer_verified"`
	BuyerUnverified     int64 `json:"buyer_unverified"`
	RegistrasiHariIni   int64 `json:"registrasi_hari_ini"`
	RegistrasiMingguIni int64 `json:"registrasi_minggu_ini"`
	RegistrasiBulanIni  int64 `json:"registrasi_bulan_ini"`
}

// ========================================
// Alamat Buyer Response
// ========================================

type AlamatBuyerResponse struct {
	ID              string          `json:"id"`
	BuyerID         string          `json:"buyer_id"`
	Label           string          `json:"label"`
	NamaPenerima    string          `json:"nama_penerima"`
	TeleponPenerima string          `json:"telepon_penerima"`
	Wilayah         WilayahResponse `json:"wilayah"`
	KodePos         string          `json:"kode_pos"`
	AlamatLengkap   string          `json:"alamat_lengkap"`
	AlamatFormatted string          `json:"alamat_formatted"`
	Catatan         *string         `json:"catatan"`
	IsDefault       bool            `json:"is_default"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type WilayahResponse struct {
	Kelurahan WilayahItemResponse `json:"kelurahan"`
	Kecamatan WilayahItemResponse `json:"kecamatan"`
	Kota      WilayahItemResponse `json:"kota"`
	Provinsi  WilayahItemResponse `json:"provinsi"`
}

type WilayahItemResponse struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
	Kode string `json:"kode"`
}

// ========================================
// Provinsi Response
// ========================================

type ProvinsiListResponse struct {
	ID         string `json:"id"`
	Nama       string `json:"nama"`
	Kode       string `json:"kode"`
	JumlahKota int    `json:"jumlah_kota"`
}

type ProvinsiDetailResponse struct {
	ID        string    `json:"id"`
	Nama      string    `json:"nama"`
	Kode      string    `json:"kode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ========================================
// Kota Response
// ========================================

type KotaListResponse struct {
	ID              string         `json:"id"`
	Nama            string         `json:"nama"`
	Kode            string         `json:"kode"`
	Provinsi        ProvinsiSimple `json:"provinsi"`
	JumlahKecamatan int            `json:"jumlah_kecamatan"`
}

type KotaDetailResponse struct {
	ID        string         `json:"id"`
	Nama      string         `json:"nama"`
	Kode      string         `json:"kode"`
	Provinsi  ProvinsiSimple `json:"provinsi"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// ========================================
// Kecamatan Response
// ========================================

type KecamatanListResponse struct {
	ID              string     `json:"id"`
	Nama            string     `json:"nama"`
	Kode            string     `json:"kode"`
	Kota            KotaSimple `json:"kota"`
	JumlahKelurahan int        `json:"jumlah_kelurahan"`
}

type KecamatanDetailResponse struct {
	ID        string     `json:"id"`
	Nama      string     `json:"nama"`
	Kode      string     `json:"kode"`
	Kota      KotaSimple `json:"kota"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ========================================
// Kelurahan Response
// ========================================

type KelurahanListResponse struct {
	ID        string          `json:"id"`
	Nama      string          `json:"nama"`
	Kode      string          `json:"kode"`
	Kecamatan KecamatanSimple `json:"kecamatan"`
}

type KelurahanDetailResponse struct {
	ID        string          `json:"id"`
	Nama      string          `json:"nama"`
	Kode      string          `json:"kode"`
	Kecamatan KecamatanSimple `json:"kecamatan"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// ========================================
// Simple Response (for nested)
// ========================================

type ProvinsiSimple struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
	Kode string `json:"kode"`
}

type KotaSimple struct {
	ID       string         `json:"id"`
	Nama     string         `json:"nama"`
	Kode     string         `json:"kode"`
	Provinsi ProvinsiSimple `json:"provinsi"`
}

type KecamatanSimple struct {
	ID   string     `json:"id"`
	Nama string     `json:"nama"`
	Kode string     `json:"kode"`
	Kota KotaSimple `json:"kota"`
}

// ========================================
// Wilayah Dropdown Response
// ========================================

type WilayahDropdownItem struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
	Kode string `json:"kode"`
}
