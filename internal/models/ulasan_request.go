package models

// ========================================
// Request dari Buyer (Mobile App)
// ========================================

type CreateUlasanRequest struct {
	PesananItemID string  `json:"pesanan_item_id" binding:"required,uuid"`
	Rating        int     `json:"rating" binding:"required,min=1,max=5"`
	Komentar      *string `json:"komentar" binding:"omitempty,max=1000"`
	// Gambar di-upload via multipart/form-data
}

// ========================================
// Request dari Admin (Approval)
// ========================================

type ApproveUlasanRequest struct {
	IsApproved bool `json:"is_approved"`
}

// Bulk approve/reject
type BulkApproveUlasanRequest struct {
	IDs        []string `json:"ids" binding:"required,min=1"`
	IsApproved bool     `json:"is_approved"`
}
