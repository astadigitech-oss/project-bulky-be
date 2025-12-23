package models

// ========================================
// Pagination Request
// ========================================

type PaginationRequest struct {
	Halaman         int    `form:"halaman"`
	PerHalaman      int    `form:"per_halaman"`
	Cari            string `form:"cari"`
	UrutBerdasarkan string `form:"urut_berdasarkan"`
	Urutan          string `form:"urutan"`
	IsActive        *bool  `form:"is_active"`
}

func (p *PaginationRequest) SetDefaults() {
	if p.Halaman <= 0 {
		p.Halaman = 1
	}
	if p.PerHalaman <= 0 {
		p.PerHalaman = 10
	}
	if p.PerHalaman > 100 {
		p.PerHalaman = 100
	}
	if p.UrutBerdasarkan == "" {
		p.UrutBerdasarkan = "created_at"
	}
	if p.Urutan == "" {
		p.Urutan = "desc"
	}
}

func (p *PaginationRequest) GetOffset() int {
	return (p.Halaman - 1) * p.PerHalaman
}

// ========================================
// Kategori Produk Request
// ========================================

type CreateKategoriProdukRequest struct {
	Nama                    string  `json:"nama" binding:"required,min=2,max=100"`
	Deskripsi               *string `json:"deskripsi" binding:"omitempty,max=1000"`
	Icon                    *string `json:"icon"`
	MemilikiKondisiTambahan bool    `json:"memiliki_kondisi_tambahan"`
	TipeKondisiTambahan     *string `json:"tipe_kondisi_tambahan" binding:"omitempty,oneof=gambar teks"`
	GambarKondisi           *string `json:"gambar_kondisi"`
	TeksKondisi             *string `json:"teks_kondisi" binding:"omitempty,max=500"`
}

type UpdateKategoriProdukRequest struct {
	Nama                    *string `json:"nama" binding:"omitempty,min=2,max=100"`
	Deskripsi               *string `json:"deskripsi" binding:"omitempty,max=1000"`
	Icon                    *string `json:"icon"`
	MemilikiKondisiTambahan *bool   `json:"memiliki_kondisi_tambahan"`
	TipeKondisiTambahan     *string `json:"tipe_kondisi_tambahan" binding:"omitempty,oneof=gambar teks"`
	GambarKondisi           *string `json:"gambar_kondisi"`
	TeksKondisi             *string `json:"teks_kondisi" binding:"omitempty,max=500"`
	IsActive                *bool   `json:"is_active"`
}

// ========================================
// Merek Produk Request
// ========================================

type CreateMerekProdukRequest struct {
	Nama string  `json:"nama" binding:"required,min=2,max=100"`
	Logo *string `json:"logo"`
}

type UpdateMerekProdukRequest struct {
	Nama     *string `json:"nama" binding:"omitempty,min=2,max=100"`
	Logo     *string `json:"logo"`
	IsActive *bool   `json:"is_active"`
}

// ========================================
// Kondisi Produk Request
// ========================================

type CreateKondisiProdukRequest struct {
	Nama      string  `json:"nama" binding:"required,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
	Urutan    *int    `json:"urutan"`
}

type UpdateKondisiProdukRequest struct {
	Nama      *string `json:"nama" binding:"omitempty,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
	Urutan    *int    `json:"urutan"`
	IsActive  *bool   `json:"is_active"`
}

// ========================================
// Kondisi Paket Request
// ========================================

type CreateKondisiPaketRequest struct {
	Nama      string  `json:"nama" binding:"required,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
	Urutan    *int    `json:"urutan"`
}

type UpdateKondisiPaketRequest struct {
	Nama      *string `json:"nama" binding:"omitempty,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
	Urutan    *int    `json:"urutan"`
	IsActive  *bool   `json:"is_active"`
}

// ========================================
// Sumber Produk Request
// ========================================

type CreateSumberProdukRequest struct {
	Nama      string  `json:"nama" binding:"required,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
}

type UpdateSumberProdukRequest struct {
	Nama      *string `json:"nama" binding:"omitempty,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
	IsActive  *bool   `json:"is_active"`
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

type CreateTipeProdukRequest struct {
	Nama      string  `json:"nama" binding:"required,min=2,max=100"`
	Deskripsi *string `json:"deskripsi" binding:"omitempty,max=500"`
	Urutan    *int    `json:"urutan"`
}

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
	Urutan       *int   `json:"urutan"`
}

type UpdateBannerTipeProdukRequest struct {
	TipeProdukID *string `json:"tipe_produk_id" binding:"omitempty,uuid"`
	Nama         *string `json:"nama" binding:"omitempty,min=2,max=100"`
	GambarURL    *string `json:"gambar_url" binding:"omitempty,max=500"`
	Urutan       *int    `json:"urutan"`
	IsActive     *bool   `json:"is_active"`
}


// ========================================
// Produk Request
// ========================================

type CreateProdukRequest struct {
	Nama               string   `json:"nama" binding:"required,min=2,max=255"`
	IDCargo            *string  `json:"id_cargo" binding:"omitempty,max=50"`
	KategoriID         string   `json:"kategori_id" binding:"required,uuid"`
	MerekID            *string  `json:"merek_id" binding:"omitempty,uuid"`
	KondisiID          string   `json:"kondisi_id" binding:"required,uuid"`
	KondisiPaketID     string   `json:"kondisi_paket_id" binding:"required,uuid"`
	SumberID           *string  `json:"sumber_id" binding:"omitempty,uuid"`
	WarehouseID        string   `json:"warehouse_id" binding:"required,uuid"`
	TipeProdukID       string   `json:"tipe_produk_id" binding:"required,uuid"`
	HargaSebelumDiskon float64  `json:"harga_sebelum_diskon" binding:"required,gt=0"`
	PersentaseDiskon   float64  `json:"persentase_diskon" binding:"min=0,max=100"`
	HargaSesudahDiskon float64  `json:"harga_sesudah_diskon" binding:"required,min=0"`
	Quantity           int      `json:"quantity" binding:"min=0"`
	Discrepancy        *string  `json:"discrepancy" binding:"omitempty,max=1000"`
	GambarURLs         []string `json:"gambar_urls" binding:"required,min=1"`
	GambarUtamaIndex   int      `json:"gambar_utama_index"`
}

type UpdateProdukRequest struct {
	Nama               *string  `json:"nama" binding:"omitempty,min=2,max=255"`
	IDCargo            *string  `json:"id_cargo" binding:"omitempty,max=50"`
	KategoriID         *string  `json:"kategori_id" binding:"omitempty,uuid"`
	MerekID            *string  `json:"merek_id" binding:"omitempty,uuid"`
	KondisiID          *string  `json:"kondisi_id" binding:"omitempty,uuid"`
	KondisiPaketID     *string  `json:"kondisi_paket_id" binding:"omitempty,uuid"`
	SumberID           *string  `json:"sumber_id" binding:"omitempty,uuid"`
	WarehouseID        *string  `json:"warehouse_id" binding:"omitempty,uuid"`
	TipeProdukID       *string  `json:"tipe_produk_id" binding:"omitempty,uuid"`
	HargaSebelumDiskon *float64 `json:"harga_sebelum_diskon" binding:"omitempty,gt=0"`
	PersentaseDiskon   *float64 `json:"persentase_diskon" binding:"omitempty,min=0,max=100"`
	HargaSesudahDiskon *float64 `json:"harga_sesudah_diskon" binding:"omitempty,min=0"`
	Quantity           *int     `json:"quantity" binding:"omitempty,min=0"`
	Discrepancy        *string  `json:"discrepancy" binding:"omitempty,max=1000"`
	IsActive           *bool    `json:"is_active"`
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
	GambarURL string `json:"gambar_url" binding:"required,max=500"`
	Urutan    *int   `json:"urutan"`
	IsPrimary bool   `json:"is_primary"`
}

type UpdateProdukGambarRequest struct {
	Urutan    *int  `json:"urutan"`
	IsPrimary *bool `json:"is_primary"`
}

// ========================================
// Produk Dokumen Request
// ========================================

type CreateProdukDokumenRequest struct {
	NamaDokumen string `json:"nama_dokumen" binding:"required,max=255"`
	FileURL     string `json:"file_url" binding:"required,max=500"`
	TipeFile    string `json:"tipe_file" binding:"required,max=50"`
	UkuranFile  *int   `json:"ukuran_file"`
}
