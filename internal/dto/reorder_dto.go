package dto

// ReorderDirectionRequest - Request untuk reorder dengan direction (up/down)
// Digunakan untuk semua endpoint reorder yang memindahkan item satu posisi
type ReorderDirectionRequest struct {
	Direction string `json:"direction" binding:"required,oneof=up down"`
}
