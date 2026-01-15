package dto

// Create/Update Request untuk master data dengan nama dual bahasa
type MerekProdukRequest struct {
	NamaID   string  `json:"nama_id" binding:"required,min=2,max=100"`
	NamaEN   *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Slug     string  `json:"slug,omitempty"`
	LogoURL  *string `json:"logo_url,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type KategoriProdukRequest struct {
	NamaID                  string  `json:"nama_id" binding:"required,min=2,max=100"`
	NamaEN                  *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Slug                    string  `json:"slug,omitempty"`
	Deskripsi               *string `json:"deskripsi,omitempty"`
	IconURL                 *string `json:"icon_url,omitempty"`
	MemilikiKondisiTambahan *bool   `json:"memiliki_kondisi_tambahan,omitempty"`
	TipeKondisiTambahan     *string `json:"tipe_kondisi_tambahan,omitempty"`
	GambarKondisiURL        *string `json:"gambar_kondisi_url,omitempty"`
	TeksKondisi             *string `json:"teks_kondisi,omitempty"`
	IsActive                *bool   `json:"is_active,omitempty"`
}

type KondisiProdukRequest struct {
	NamaID    string  `json:"nama_id" binding:"required,min=2,max=100"`
	NamaEN    *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Slug      string  `json:"slug,omitempty"`
	Deskripsi *string `json:"deskripsi,omitempty"`
	Urutan    *int    `json:"urutan,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

type SumberProdukRequest struct {
	NamaID    string  `json:"nama_id" binding:"required,min=2,max=100"`
	NamaEN    *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Slug      string  `json:"slug,omitempty"`
	Deskripsi *string `json:"deskripsi,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

type KondisiPaketRequest struct {
	NamaID    string  `json:"nama_id" binding:"required,min=2,max=100"`
	NamaEN    *string `json:"nama_en,omitempty" binding:"omitempty,min=2,max=100"`
	Slug      string  `json:"slug,omitempty"`
	Deskripsi *string `json:"deskripsi,omitempty"`
	Urutan    *int    `json:"urutan,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}
