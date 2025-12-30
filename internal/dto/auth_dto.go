package dto

// AdminUpdateProfileRequest untuk update profile admin
type AdminUpdateProfileRequest struct {
	Nama  string `json:"nama" binding:"required,min=3,max=100"`
	Email string `json:"email" binding:"required,email,max=255"`
}

// BuyerUpdateProfileRequest untuk update profile buyer
type BuyerUpdateProfileRequest struct {
	Nama     string `json:"nama" binding:"required,min=3,max=100"`
	Username string `json:"username" binding:"required,min=3,max=50,alphanumund"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Telepon  string `json:"telepon" binding:"required,min=10,max=15"`
}
