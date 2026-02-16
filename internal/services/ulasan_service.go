package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
	"strings"

	"github.com/google/uuid"
)

type UlasanService interface {
	// Admin
	AdminFindAll(filters map[string]interface{}, page, perPage int, sortBy, sortOrder string) ([]models.UlasanAdminListResponse, int64, map[string]int64, error)
	AdminFindByID(id string) (*models.UlasanAdminDetailResponse, error)
	Approve(id string, isApproved bool, adminID uuid.UUID) error
	BulkApprove(ids []string, isApproved bool, adminID uuid.UUID) (int64, error)
	Delete(id string) error

	// Buyer
	GetPendingReviews(buyerID uuid.UUID) ([]models.BuyerPendingReviewItem, error)
	BuyerFindAll(buyerID uuid.UUID, page, perPage int) ([]models.BuyerUlasanResponse, int64, error)
	Create(req models.CreateUlasanRequest, buyerID uuid.UUID, gambarFile *multipart.FileHeader) (*models.Ulasan, error)

	// Public
	GetProdukUlasan(produkID string, filters map[string]interface{}, page, perPage int, sortBy, sortOrder string) (*models.ProdukUlasanWithStats, int64, error)
	GetProdukRatingStats(produkID string) (*models.ProdukRatingStats, error)
}

type ulasanService struct {
	repo            repositories.UlasanRepository
	pesananItemRepo repositories.PesananItemRepository
	pesananRepo     repositories.PesananRepository
	uploadPath      string
	baseURL         string
}

func NewUlasanService(
	repo repositories.UlasanRepository,
	pesananItemRepo repositories.PesananItemRepository,
	pesananRepo repositories.PesananRepository,
	uploadPath, baseURL string,
) UlasanService {
	return &ulasanService{
		repo:            repo,
		pesananItemRepo: pesananItemRepo,
		pesananRepo:     pesananRepo,
		uploadPath:      uploadPath,
		baseURL:         baseURL,
	}
}

// ========================================
// Admin Methods
// ========================================

func (s *ulasanService) AdminFindAll(filters map[string]interface{}, page, perPage int, sortBy, sortOrder string) ([]models.UlasanAdminListResponse, int64, map[string]int64, error) {
	ulasan, total, err := s.repo.AdminFindAll(filters, page, perPage, sortBy, sortOrder)
	if err != nil {
		return nil, 0, nil, err
	}

	// Get summary
	summary, err := s.repo.GetSummary()
	if err != nil {
		return nil, 0, nil, err
	}

	// Map to response
	response := make([]models.UlasanAdminListResponse, len(ulasan))
	for i, u := range ulasan {
		response[i] = models.UlasanAdminListResponse{
			ID:          u.ID.String(),
			NamaBuyer:   u.Buyer.Nama,
			EmailBuyer:  u.Buyer.Email,
			KodePesanan: u.Pesanan.Kode,
			NamaProduk:  u.Produk.NamaID,
			Rating:      u.Rating,
			Komentar:    u.Komentar,
			Gambar:      u.Gambar,
			GambarURL:   s.getImageURL(u.Gambar),
			IsApproved:  u.IsApproved,
			ApprovedAt:  u.ApprovedAt,
			CreatedAt:   u.CreatedAt,
		}
	}

	return response, total, summary, nil
}

func (s *ulasanService) AdminFindByID(id string) (*models.UlasanAdminDetailResponse, error) {
	ulasanID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid ulasan ID")
	}

	ulasan, err := s.repo.AdminFindByID(ulasanID)
	if err != nil {
		return nil, err
	}

	response := &models.UlasanAdminDetailResponse{
		ID: ulasan.ID.String(),
		Buyer: models.BuyerSimpleResponse{
			ID:      ulasan.Buyer.ID.String(),
			Nama:    ulasan.Buyer.Nama,
			Email:   ulasan.Buyer.Email,
			Telepon: ulasan.Buyer.Telepon,
		},
		Pesanan: struct {
			ID   string `json:"id"`
			Kode string `json:"kode"`
		}{
			ID:   ulasan.Pesanan.ID.String(),
			Kode: ulasan.Pesanan.Kode,
		},
		Produk: struct {
			ID   string `json:"id"`
			Nama string `json:"nama"`
			SKU  string `json:"sku"`
		}{
			ID:   ulasan.Produk.ID.String(),
			Nama: ulasan.Produk.NamaID,
			SKU: func() string {
				if ulasan.Produk.IDCargo != nil {
					return *ulasan.Produk.IDCargo
				}
				return ulasan.Produk.ID.String()
			}(),
		},
		Rating:     ulasan.Rating,
		Komentar:   ulasan.Komentar,
		Gambar:     ulasan.Gambar,
		GambarURL:  s.getImageURL(ulasan.Gambar),
		IsApproved: ulasan.IsApproved,
		ApprovedAt: ulasan.ApprovedAt,
		CreatedAt:  ulasan.CreatedAt,
	}

	if ulasan.ApprovedBy != nil {
		approvedByStr := ulasan.ApprovedBy.String()
		response.ApprovedBy = &approvedByStr
	}

	return response, nil
}

func (s *ulasanService) Approve(id string, isApproved bool, adminID uuid.UUID) error {
	ulasanID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid ulasan ID")
	}

	return s.repo.Approve(ulasanID, isApproved, adminID)
}

func (s *ulasanService) BulkApprove(ids []string, isApproved bool, adminID uuid.UUID) (int64, error) {
	uuids := make([]uuid.UUID, len(ids))
	for i, id := range ids {
		parsed, err := uuid.Parse(id)
		if err != nil {
			return 0, fmt.Errorf("invalid ID at index %d: %s", i, id)
		}
		uuids[i] = parsed
	}

	approvedIDs, err := s.repo.BulkApprove(uuids, isApproved, adminID)
	if err != nil {
		return 0, err
	}

	return int64(len(approvedIDs)), nil
}

func (s *ulasanService) Delete(id string) error {
	ulasanID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid ulasan ID")
	}

	return s.repo.Delete(ulasanID)
}

// ========================================
// Buyer Methods
// ========================================

func (s *ulasanService) GetPendingReviews(buyerID uuid.UUID) ([]models.BuyerPendingReviewItem, error) {
	items, err := s.repo.GetPendingReviews(buyerID)
	if err != nil {
		return nil, err
	}

	response := make([]models.BuyerPendingReviewItem, len(items))
	for i, item := range items {
		var gambarProduk *string
		if len(item.Produk.Gambar) > 0 {
			gambarProduk = &item.Produk.Gambar[0].GambarURL
		}

		response[i] = models.BuyerPendingReviewItem{
			PesananItemID: item.ID.String(),
			PesananKode:   item.Pesanan.Kode,
			ProdukID:      item.ProdukID.String(),
			NamaProduk:    item.Produk.NamaID,
			GambarProduk:  gambarProduk,
			Qty:           item.Qty,
			CompletedAt:   item.Pesanan.UpdatedAt, // Use pesanan updated_at as completed_at
		}
	}

	return response, nil
}

func (s *ulasanService) BuyerFindAll(buyerID uuid.UUID, page, perPage int) ([]models.BuyerUlasanResponse, int64, error) {
	ulasan, total, err := s.repo.BuyerFindAll(buyerID, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	response := make([]models.BuyerUlasanResponse, len(ulasan))
	for i, u := range ulasan {
		var gambarProduk *string
		if len(u.Produk.Gambar) > 0 {
			gambarProduk = &u.Produk.Gambar[0].GambarURL
		}

		response[i] = models.BuyerUlasanResponse{
			ID:           u.ID.String(),
			NamaProduk:   u.Produk.NamaID,
			GambarProduk: gambarProduk,
			Rating:       u.Rating,
			Komentar:     u.Komentar,
			GambarURL:    s.getImageURL(u.Gambar),
			IsApproved:   u.IsApproved,
			CreatedAt:    u.CreatedAt,
		}
	}

	return response, total, nil
}

func (s *ulasanService) Create(req models.CreateUlasanRequest, buyerID uuid.UUID, gambarFile *multipart.FileHeader) (*models.Ulasan, error) {
	// Parse pesanan item ID
	pesananItemID, err := uuid.Parse(req.PesananItemID)
	if err != nil {
		return nil, errors.New("invalid pesanan_item_id")
	}

	// Check if already reviewed
	exists, err := s.repo.CheckExistingReview(pesananItemID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("item sudah pernah di-review")
	}

	// Get pesanan item to validate
	pesananItem, err := s.pesananItemRepo.FindByID(pesananItemID)
	if err != nil {
		return nil, errors.New("pesanan item tidak ditemukan")
	}

	// Get pesanan to validate buyer and status
	pesanan, err := s.pesananRepo.FindByID(pesananItem.PesananID)
	if err != nil {
		return nil, errors.New("pesanan tidak ditemukan")
	}

	// Validate buyer ownership
	if pesanan.BuyerID != buyerID {
		return nil, errors.New("pesanan bukan milik buyer ini")
	}

	// Validate order status
	if pesanan.OrderStatus != "COMPLETED" {
		return nil, errors.New("ulasan hanya dapat diberikan untuk pesanan yang sudah selesai (COMPLETED)")
	}

	// Handle image upload
	var gambarPath *string
	if gambarFile != nil {
		path, err := s.uploadImage(gambarFile)
		if err != nil {
			return nil, fmt.Errorf("gagal upload gambar: %v", err)
		}
		gambarPath = &path
	}

	// Create ulasan
	ulasan := &models.Ulasan{
		PesananID:     pesanan.ID,
		PesananItemID: pesananItem.ID,
		BuyerID:       buyerID,
		ProdukID:      pesananItem.ProdukID,
		Rating:        req.Rating,
		Komentar:      req.Komentar,
		Gambar:        gambarPath,
		IsApproved:    false,
	}

	if err := s.repo.Create(ulasan); err != nil {
		// Cleanup uploaded image if create fails
		if gambarPath != nil {
			os.Remove(filepath.Join(s.uploadPath, filepath.FromSlash(*gambarPath)))
		}
		return nil, err
	}

	return ulasan, nil
}

// ========================================
// Public Methods
// ========================================

func (s *ulasanService) GetProdukUlasan(produkID string, filters map[string]interface{}, page, perPage int, sortBy, sortOrder string) (*models.ProdukUlasanWithStats, int64, error) {
	produkUUID, err := uuid.Parse(produkID)
	if err != nil {
		return nil, 0, errors.New("invalid produk ID")
	}

	// Get rating stats
	stats, err := s.repo.GetProdukRatingStats(produkUUID)
	if err != nil {
		return nil, 0, err
	}

	// Get ulasan list
	ulasan, total, err := s.repo.GetProdukUlasan(produkUUID, filters, page, perPage, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	// Map to public response
	ulasanResponse := make([]models.UlasanPublicResponse, len(ulasan))
	for i, u := range ulasan {
		ulasanResponse[i] = models.UlasanPublicResponse{
			ID:        u.ID.String(),
			NamaBuyer: utils.MaskName(u.Buyer.Nama),
			Rating:    u.Rating,
			Komentar:  u.Komentar,
			GambarURL: s.getImageURL(u.Gambar),
			CreatedAt: u.CreatedAt,
		}
	}

	result := &models.ProdukUlasanWithStats{
		Stats:  *stats,
		Ulasan: ulasanResponse,
	}

	return result, total, nil
}

func (s *ulasanService) GetProdukRatingStats(produkID string) (*models.ProdukRatingStats, error) {
	produkUUID, err := uuid.Parse(produkID)
	if err != nil {
		return nil, errors.New("invalid produk ID")
	}

	return s.repo.GetProdukRatingStats(produkUUID)
}

// ========================================
// Helper Methods
// ========================================

func (s *ulasanService) getImageURL(path *string) *string {
	if path == nil || *path == "" {
		return nil
	}
	fullURL := fmt.Sprintf("%s/%s", s.baseURL, *path)
	return &fullURL
}

func (s *ulasanService) uploadImage(file *multipart.FileHeader) (string, error) {
	// Validate file type
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", errors.New("format file tidak didukung. Gunakan jpg, jpeg, atau png")
	}

	// Validate file size (max 2MB)
	if file.Size > 2*1024*1024 {
		return "", errors.New("ukuran file maksimal 2MB")
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	relativePath := filepath.Join("ulasan", filename)
	// Normalize to forward slash for URL
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	fullPath := filepath.Join(s.uploadPath, filepath.FromSlash(relativePath))

	// Create directory if not exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Save file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	if _, err := dst.ReadFrom(src); err != nil {
		dst.Close()
		return "", err
	}

	// Ensure file is fully written to disk
	if err := dst.Sync(); err != nil {
		dst.Close()
		return "", err
	}

	// Close file explicitly
	if err := dst.Close(); err != nil {
		return "", err
	}

	return relativePath, nil
}
