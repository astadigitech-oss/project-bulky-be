package dto

// DasborPeriodeQuery is a common query param for dasbor endpoints
type DasborPeriodeQuery struct {
	Periode string `query:"periode"`
}

func (q *DasborPeriodeQuery) SetDefault() {
	if q.Periode == "" {
		q.Periode = "bulan_ini"
	}
}

// DasborPenjualanPerBuyerQuery query params for penjualan-per-buyer
type DasborPenjualanPerBuyerQuery struct {
	Periode string `query:"periode"`
	Limit   int    `query:"limit"`
}

func (q *DasborPenjualanPerBuyerQuery) SetDefaults() {
	if q.Periode == "" {
		q.Periode = "bulan_ini"
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
}

// DasborTabelTransaksiQuery query params for tabel-transaksi
type DasborTabelTransaksiQuery struct {
	Periode    string `query:"periode"`
	Halaman    int    `query:"halaman"`
	PerHalaman int    `query:"per_halaman"`
}

func (q *DasborTabelTransaksiQuery) SetDefaults() {
	if q.Periode == "" {
		q.Periode = "bulan_ini"
	}
	if q.Halaman < 1 {
		q.Halaman = 1
	}
	if q.PerHalaman < 1 {
		q.PerHalaman = 10
	}
	if q.PerHalaman > 100 {
		q.PerHalaman = 100
	}
}

// ---- Response types ----

// DasborChartTransaksiResponse for GET /chart-transaksi
type DasborChartTransaksiResponse struct {
	Periode string                           `json:"periode"`
	Labels  []string                         `json:"labels"`
	Series  DasborChartTransaksiSeriesData   `json:"series"`
}

type DasborChartTransaksiSeriesData struct {
	Success []int64 `json:"success"`
	Cancel  []int64 `json:"cancel"`
}

// DasborChartRevenueResponse for GET /chart-revenue
type DasborChartRevenueResponse struct {
	Periode         string                         `json:"periode"`
	Labels          []string                       `json:"labels"`
	Series          DasborChartRevenueSeries       `json:"series"`
	TotalKeseluruhan float64                       `json:"total_keseluruhan"`
}

type DasborChartRevenueSeries struct {
	TotalPenjualan []float64 `json:"total_penjualan"`
}

// DasborChartKategoriResponse for GET /chart-transaksi-per-kategori
type DasborChartKategoriResponse struct {
	Periode string                        `json:"periode"`
	Labels  []string                      `json:"labels"`
	Series  []DasborChartKategoriSerie    `json:"series"`
}

type DasborChartKategoriSerie struct {
	Kategori   string  `json:"kategori"`
	KategoriID string  `json:"kategori_id"`
	Data       []int64 `json:"data"`
}

// DasborKPIResponse for GET /kpi
type DasborKPIResponse struct {
	Periode           string  `json:"periode"`
	StokPaletbox   int64   `json:"stok_paletbox"`
	PaletboxTerjual int64   `json:"paletbox_terjual"`
	Revenue           float64 `json:"revenue"`
}

// DasborStokPerKategoriResponse for GET /stok-per-kategori
type DasborStokPerKategoriResponse struct {
	Labels []string                      `json:"labels"`
	Series DasborStokPerKategoriSeries   `json:"series"`
}

type DasborStokPerKategoriSeries struct {
	Stok []int64 `json:"stok"`
}

// DasborPenjualanPerBuyerResponse for GET /penjualan-per-buyer
type DasborPenjualanPerBuyerResponse struct {
	Periode string                         `json:"periode"`
	Labels  []string                       `json:"labels"`
	Series  DasborPenjualanPerBuyerSeries  `json:"series"`
	Buyers  []DasborBuyerDetail            `json:"buyers"`
}

type DasborPenjualanPerBuyerSeries struct {
	TotalPembelian []float64 `json:"total_pembelian"`
}

type DasborBuyerDetail struct {
	BuyerID        string  `json:"buyer_id"`
	Nama           string  `json:"nama"`
	TotalPembelian float64 `json:"total_pembelian"`
}

// DasborTabelTransaksiItem single row in tabel-transaksi
type DasborTabelTransaksiItem struct {
	PesananID      string   `json:"pesanan_id"`
	Kode           string   `json:"kode"`
	NamaPembeli    string   `json:"nama_pembeli"`
	Palet          string   `json:"palet"`
	Kategori       string   `json:"kategori"`
	Harga          float64  `json:"harga"`
	OngkosKirim    float64  `json:"ongkos_kirim"`
	Diskon         float64  `json:"diskon"`
	Total          float64  `json:"total"`
	TanggalPesanan string   `json:"tanggal_pesanan"`
	DeliveryType   string   `json:"delivery_type"`
	PaymentType    string   `json:"payment_type"`
	OrderStatus    string   `json:"order_status"`
	JenisPembayaran []string `json:"jenis_pembayaran"`
}

// DasborTabelTransaksiMeta pagination meta for tabel-transaksi
type DasborTabelTransaksiMeta struct {
	Halaman      int   `json:"halaman"`
	PerHalaman   int   `json:"per_halaman"`
	TotalData    int64 `json:"total_data"`
	TotalHalaman int   `json:"total_halaman"`
}

// DasborUserTransaksiItem single row in user-transaction
type DasborUserTransaksiItem struct {
	BuyerID        string  `json:"buyer_id"`
	Nama           string  `json:"nama"`
	TotalTransaksi int64   `json:"total_transaksi"`
	TotalBelanja   float64 `json:"total_belanja"`
}

// ---- Raw DB scan types ----

// DasborDateCount used to scan date + count rows
type DasborDateCount struct {
	Tanggal string `gorm:"column:tanggal"`
	Jumlah  int64  `gorm:"column:jumlah"`
}

// DasborDateAmount used to scan date + amount rows
type DasborDateAmount struct {
	Tanggal        string  `gorm:"column:tanggal"`
	TotalPenjualan float64 `gorm:"column:total_penjualan"`
}

// DasborKategoriDateCount used for kategori chart raw data
type DasborKategoriDateCount struct {
	Tanggal    string `gorm:"column:tanggal"`
	Kategori   string `gorm:"column:kategori"`
	KategoriID string `gorm:"column:kategori_id"`
	Jumlah     int64  `gorm:"column:jumlah_transaksi"`
}

// DasborKategoriStok stok per kategori raw data
type DasborKategoriStok struct {
	Kategori  string `gorm:"column:kategori"`
	TotalStok int64  `gorm:"column:total_stok"`
}

// DasborBuyerPenjualan penjualan per buyer raw data
type DasborBuyerPenjualan struct {
	BuyerID        string  `gorm:"column:buyer_id"`
	Nama           string  `gorm:"column:nama"`
	TotalPembelian float64 `gorm:"column:total_pembelian"`
}

// DasborUserTransaksiRaw raw data for user transaction ranking
type DasborUserTransaksiRaw struct {
	BuyerID        string  `gorm:"column:buyer_id"`
	Nama           string  `gorm:"column:nama"`
	TotalTransaksi int64   `gorm:"column:total_transaksi"`
	TotalBelanja   float64 `gorm:"column:total_belanja"`
}

// DasborTabelRow raw row from tabel transaksi query
type DasborTabelRow struct {
	PesananID   string  `gorm:"column:pesanan_id"`
	Kode        string  `gorm:"column:kode"`
	NamaPembeli string  `gorm:"column:nama_pembeli"`
	BiayaProduk float64 `gorm:"column:biaya_produk"`
	OngkosKirim float64 `gorm:"column:ongkos_kirim"`
	BiayaPPN    float64 `gorm:"column:biaya_ppn"`
	Total       float64 `gorm:"column:total"`
	CreatedAt   string  `gorm:"column:tanggal_pesanan"`
	DeliveryType string `gorm:"column:delivery_type"`
	PaymentType string  `gorm:"column:payment_type"`
	OrderStatus string  `gorm:"column:order_status"`
}

// DasborTabelItemRow raw item row for a pesanan
type DasborTabelItemRow struct {
	PesananID    string  `gorm:"column:pesanan_id"`
	NamaProduk   string  `gorm:"column:nama_produk"`
	DiskonSatuan float64 `gorm:"column:diskon_satuan"`
	Qty          int     `gorm:"column:qty"`
	KategoriNama string  `gorm:"column:kategori_nama"`
	RowNum       int     `gorm:"column:row_num"`
	TotalItems   int     `gorm:"column:total_items"`
}

// DasborTabelPembayaranRow raw payment row for a pesanan
type DasborTabelPembayaranRow struct {
	PesananID    string `gorm:"column:pesanan_id"`
	NamaMetode   string `gorm:"column:nama_metode"`
}
