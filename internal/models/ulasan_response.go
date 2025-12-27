package models

import "time"

// ========================================
// Response untuk Admin
// ========================================

type UlasanAdminListResponse struct {
	ID          string     `json:"id"`
	NamaBuyer   string     `json:"nama_buyer"`
	EmailBuyer  string     `json:"email_buyer"`
	KodePesanan string     `json:"kode_pesanan"`
	NamaProduk  string     `json:"nama_produk"`
	Rating      int        `json:"rating"`
	Komentar    *string    `json:"komentar"`
	Gambar      *string    `json:"gambar"`
	GambarURL   *string    `json:"gambar_url"`
	IsApproved  bool       `json:"is_approved"`
	ApprovedAt  *time.Time `json:"approved_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type UlasanAdminDetailResponse struct {
	ID      string              `json:"id"`
	Buyer   BuyerSimpleResponse `json:"buyer"`
	Pesanan struct {
		ID   string `json:"id"`
		Kode string `json:"kode"`
	} `json:"pesanan"`
	Produk struct {
		ID   string `json:"id"`
		Nama string `json:"nama"`
		SKU  string `json:"sku"`
	} `json:"produk"`
	Rating     int        `json:"rating"`
	Komentar   *string    `json:"komentar"`
	Gambar     *string    `json:"gambar"`
	GambarURL  *string    `json:"gambar_url"`
	IsApproved bool       `json:"is_approved"`
	ApprovedAt *time.Time `json:"approved_at"`
	ApprovedBy *string    `json:"approved_by"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ========================================
// Response untuk Buyer
// ========================================

type BuyerPendingReviewItem struct {
	PesananItemID string    `json:"pesanan_item_id"`
	PesananKode   string    `json:"pesanan_kode"`
	ProdukID      string    `json:"produk_id"`
	NamaProduk    string    `json:"nama_produk"`
	GambarProduk  *string   `json:"gambar_produk"`
	Qty           int       `json:"qty"`
	CompletedAt   time.Time `json:"completed_at"`
}

type BuyerUlasanResponse struct {
	ID           string    `json:"id"`
	NamaProduk   string    `json:"nama_produk"`
	GambarProduk *string   `json:"gambar_produk"`
	Rating       int       `json:"rating"`
	Komentar     *string   `json:"komentar"`
	GambarURL    *string   `json:"gambar_url"`
	IsApproved   bool      `json:"is_approved"`
	CreatedAt    time.Time `json:"created_at"`
}

// ========================================
// Response untuk Public (Customer App)
// ========================================

type UlasanPublicResponse struct {
	ID        string    `json:"id"`
	NamaBuyer string    `json:"nama_buyer"` // Bisa di-mask: "J***n"
	Rating    int       `json:"rating"`
	Komentar  *string   `json:"komentar"`
	GambarURL *string   `json:"gambar_url"`
	CreatedAt time.Time `json:"created_at"`
}

// Statistik rating produk
type ProdukRatingStats struct {
	ProdukID    string  `json:"produk_id"`
	TotalUlasan int     `json:"total_ulasan"`
	RataRating  float64 `json:"rata_rating"`
	Rating5     int     `json:"rating_5"`
	Rating4     int     `json:"rating_4"`
	Rating3     int     `json:"rating_3"`
	Rating2     int     `json:"rating_2"`
	Rating1     int     `json:"rating_1"`
}

type ProdukUlasanWithStats struct {
	Stats  ProdukRatingStats      `json:"stats"`
	Ulasan []UlasanPublicResponse `json:"ulasan"`
}
