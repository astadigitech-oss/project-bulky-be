package services

import (
	"context"
	"errors"
	"os"
	"time"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UlasanAdminService interface {
	GetAll(ctx context.Context, params *dto.UlasanAdminQueryParams) ([]dto.UlasanAdminListResponse, *models.PaginationMeta, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.UlasanAdminDetailResponse, error)
	Approve(ctx context.Context, id uuid.UUID, adminID uuid.UUID) error
	Reject(ctx context.Context, id uuid.UUID) error
	BulkApprove(ctx context.Context, ids []uuid.UUID, adminID uuid.UUID) (*dto.BulkApproveUlasanResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ulasanAdminService struct {
	ulasanRepo repositories.UlasanRepository
}

func NewUlasanAdminService(ulasanRepo repositories.UlasanRepository) UlasanAdminService {
	return &ulasanAdminService{
		ulasanRepo: ulasanRepo,
	}
}

func (s *ulasanAdminService) GetAll(ctx context.Context, params *dto.UlasanAdminQueryParams) ([]dto.UlasanAdminListResponse, *models.PaginationMeta, error) {
	// Build filters
	filters := make(map[string]interface{})
	if params.Cari != "" {
		filters["cari"] = params.Cari
	}
	if params.Rating != nil {
		filters["rating"] = *params.Rating
	}
	if params.IsApproved != nil {
		filters["is_approved"] = *params.IsApproved
	}

	// Get data from repository
	ulasan, total, err := s.ulasanRepo.AdminFindAll(filters, params.Page, params.PerPage, params.SortBy, params.SortOrder)
	if err != nil {
		return nil, nil, err
	}

	// Map to response DTO
	response := make([]dto.UlasanAdminListResponse, len(ulasan))
	for i, u := range ulasan {
		response[i] = s.mapToListResponse(&u)
	}

	// Build meta
	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return response, &meta, nil
}

func (s *ulasanAdminService) GetByID(ctx context.Context, id uuid.UUID) (*dto.UlasanAdminDetailResponse, error) {
	ulasan, err := s.ulasanRepo.AdminFindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ulasan tidak ditemukan")
		}
		return nil, err
	}

	return s.mapToDetailResponse(ulasan), nil
}

func (s *ulasanAdminService) Approve(ctx context.Context, id uuid.UUID, adminID uuid.UUID) error {
	// Check if ulasan exists
	ulasan, err := s.ulasanRepo.AdminFindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ulasan tidak ditemukan")
		}
		return err
	}

	// Check if already approved
	if ulasan.IsApproved {
		return errors.New("ulasan sudah di-approve")
	}

	return s.ulasanRepo.Approve(id, true, adminID)
}

func (s *ulasanAdminService) Reject(ctx context.Context, id uuid.UUID) error {
	// Check if ulasan exists
	ulasan, err := s.ulasanRepo.AdminFindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ulasan tidak ditemukan")
		}
		return err
	}

	// Check if already rejected
	if !ulasan.IsApproved {
		return errors.New("ulasan sudah di-reject")
	}

	// Use zero UUID for reject (no approver)
	zeroUUID := uuid.UUID{}
	return s.ulasanRepo.Approve(id, false, zeroUUID)
}

func (s *ulasanAdminService) BulkApprove(ctx context.Context, ids []uuid.UUID, adminID uuid.UUID) (*dto.BulkApproveUlasanResponse, error) {
	approvedIDs, err := s.ulasanRepo.BulkApprove(ids, true, adminID)
	if err != nil {
		return nil, err
	}

	return &dto.BulkApproveUlasanResponse{
		ApprovedCount: len(approvedIDs),
		ApprovedIDs:   approvedIDs,
	}, nil
}

func (s *ulasanAdminService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if ulasan exists
	_, err := s.ulasanRepo.AdminFindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ulasan tidak ditemukan")
		}
		return err
	}

	return s.ulasanRepo.Delete(id)
}

// Helper functions to map models to DTOs

func (s *ulasanAdminService) mapToListResponse(u *models.Ulasan) dto.UlasanAdminListResponse {
	var approvedBy *string
	if u.ApprovedBy != nil && u.Approver != nil {
		nama := u.Approver.Nama
		approvedBy = &nama
	}

	var approvedAt *time.Time
	if u.ApprovedAt != nil {
		t := u.ApprovedAt.UTC()
		approvedAt = &t
	}

	return dto.UlasanAdminListResponse{
		ID:     u.ID,
		Rating: u.Rating,
		Buyer: dto.UlasanAdminBuyerResponse{
			ID:   u.BuyerID,
			Nama: u.Buyer.Nama,
		},
		Approved: dto.UlasanAdminApprovedInfo{
			At:     approvedAt,
			By:     approvedBy,
			Status: u.IsApproved,
		},
		Gambar:    u.Gambar != nil,
		CreatedAt: u.CreatedAt.UTC(),
		Pesanan: dto.UlasanAdminPesananResponse{
			ID:   u.PesananID,
			Kode: u.Pesanan.Kode,
		},
	}
}

func (s *ulasanAdminService) mapToDetailResponse(u *models.Ulasan) *dto.UlasanAdminDetailResponse {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	var gambarURL *string
	if u.Gambar != nil {
		url := baseURL + "/uploads/" + *u.Gambar
		gambarURL = &url
	}

	var produkGambarURL *string
	if len(u.Produk.Gambar) > 0 {
		url := baseURL + "/uploads/" + u.Produk.Gambar[0].GambarURL
		produkGambarURL = &url
	}

	response := &dto.UlasanAdminDetailResponse{
		ID: u.ID,
		Buyer: dto.UlasanAdminBuyerDetailResponse{
			ID:    u.BuyerID,
			Nama:  u.Buyer.Nama,
			Email: u.Buyer.Email,
			Telepon: func() string {
				if u.Buyer.Telepon != nil {
					return *u.Buyer.Telepon
				}
				return ""
			}(),
		},
		Pesanan: dto.UlasanAdminPesananDetailResponse{
			ID:          u.PesananID,
			Kode:        u.Pesanan.Kode,
			OrderStatus: string(u.Pesanan.OrderStatus),
			CreatedAt:   u.Pesanan.CreatedAt.UTC(),
		},
		Produk: dto.UlasanAdminProdukDetailResponse{
			ID:        u.ProdukID,
			Nama:      u.Produk.NamaID,
			Slug:      u.Produk.Slug,
			GambarURL: produkGambarURL,
		},
		Rating:     u.Rating,
		Komentar:   u.Komentar,
		GambarURL:  gambarURL,
		IsApproved: u.IsApproved,
		ApprovedAt: func() *time.Time {
			if u.ApprovedAt == nil {
				return nil
			}
			t := u.ApprovedAt.UTC()
			return &t
		}(),
		CreatedAt: u.CreatedAt.UTC(),
		UpdatedAt: u.UpdatedAt.UTC(),
	}

	if u.ApprovedBy != nil && u.Approver != nil {
		response.ApprovedBy = &dto.UlasanAdminApproverResponse{
			ID:   *u.ApprovedBy,
			Nama: u.Approver.Nama,
		}
	}

	return response
}
