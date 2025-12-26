package services

import (
	"context"
	"errors"
	"time"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type BannerEventPromoService interface {
	Create(ctx context.Context, req *models.CreateBannerEventPromoRequest) (*models.BannerEventPromoResponse, error)
	FindByID(ctx context.Context, id string) (*models.BannerEventPromoResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.BannerEventPromoResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateBannerEventPromoRequest) (*models.BannerEventPromoResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	Reorder(ctx context.Context, req *models.ReorderRequest) error
	GetVisibleBanners(ctx context.Context) ([]models.BannerEventPromoPublicResponse, error)
}

type bannerEventPromoService struct {
	repo repositories.BannerEventPromoRepository
}

func NewBannerEventPromoService(repo repositories.BannerEventPromoRepository) BannerEventPromoService {
	return &bannerEventPromoService{repo: repo}
}

func (s *bannerEventPromoService) Create(ctx context.Context, req *models.CreateBannerEventPromoRequest) (*models.BannerEventPromoResponse, error) {
	banner := &models.BannerEventPromo{
		ID:        uuid.New(),
		Nama:      req.Nama,
		Gambar:    req.Gambar,
		UrlTujuan: req.UrlTujuan,
		Urutan:    req.Urutan,
		IsActive:  req.IsActive,
	}

	if req.TanggalMulai != nil {
		t, err := time.Parse(time.RFC3339, *req.TanggalMulai)
		if err == nil {
			banner.TanggalMulai = &t
		}
	}

	if req.TanggalSelesai != nil {
		t, err := time.Parse(time.RFC3339, *req.TanggalSelesai)
		if err == nil {
			banner.TanggalSelesai = &t
		}
	}

	if err := s.repo.Create(ctx, banner); err != nil {
		return nil, err
	}

	return s.toResponse(banner), nil
}

func (s *bannerEventPromoService) FindByID(ctx context.Context, id string) (*models.BannerEventPromoResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}
	return s.toResponse(banner), nil
}

func (s *bannerEventPromoService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.BannerEventPromoResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	banners, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.BannerEventPromoResponse
	for _, b := range banners {
		items = append(items, *s.toResponse(&b))
	}

	totalHalaman := (total + int64(params.PerHalaman) - 1) / int64(params.PerHalaman)

	meta := &models.PaginationMeta{
		Halaman:      params.Halaman,
		PerHalaman:   params.PerHalaman,
		TotalData:    total,
		TotalHalaman: totalHalaman,
	}

	return items, meta, nil
}

func (s *bannerEventPromoService) Update(ctx context.Context, id string, req *models.UpdateBannerEventPromoRequest) (*models.BannerEventPromoResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}

	if req.Nama != nil {
		banner.Nama = *req.Nama
	}
	if req.Gambar != nil {
		banner.Gambar = *req.Gambar
	}
	if req.UrlTujuan != nil {
		banner.UrlTujuan = req.UrlTujuan
	}
	if req.Urutan != nil {
		banner.Urutan = *req.Urutan
	}
	if req.IsActive != nil {
		banner.IsActive = *req.IsActive
	}
	if req.TanggalMulai != nil {
		t, err := time.Parse(time.RFC3339, *req.TanggalMulai)
		if err == nil {
			banner.TanggalMulai = &t
		}
	}
	if req.TanggalSelesai != nil {
		t, err := time.Parse(time.RFC3339, *req.TanggalSelesai)
		if err == nil {
			banner.TanggalSelesai = &t
		}
	}

	if err := s.repo.Update(ctx, banner); err != nil {
		return nil, err
	}

	return s.toResponse(banner), nil
}

func (s *bannerEventPromoService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("banner tidak ditemukan")
	}
	return s.repo.Delete(ctx, id)
}

func (s *bannerEventPromoService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}

	banner.IsActive = !banner.IsActive
	if err := s.repo.Update(ctx, banner); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       banner.ID.String(),
		IsActive: banner.IsActive,
	}, nil
}

func (s *bannerEventPromoService) Reorder(ctx context.Context, req *models.ReorderRequest) error {
	return s.repo.UpdateOrder(ctx, req.Items)
}

func (s *bannerEventPromoService) GetVisibleBanners(ctx context.Context) ([]models.BannerEventPromoPublicResponse, error) {
	banners, err := s.repo.GetVisibleBanners(ctx)
	if err != nil {
		return nil, err
	}

	var items []models.BannerEventPromoPublicResponse
	for _, b := range banners {
		items = append(items, models.BannerEventPromoPublicResponse{
			ID:        b.ID.String(),
			Nama:      b.Nama,
			Gambar:    b.Gambar,
			UrlTujuan: b.UrlTujuan,
		})
	}

	return items, nil
}

func (s *bannerEventPromoService) toResponse(b *models.BannerEventPromo) *models.BannerEventPromoResponse {
	return &models.BannerEventPromoResponse{
		ID:             b.ID.String(),
		Nama:           b.Nama,
		Gambar:         b.Gambar,
		UrlTujuan:      b.UrlTujuan,
		Urutan:         b.Urutan,
		IsActive:       b.IsActive,
		IsVisible:      b.IsCurrentlyVisible(),
		TanggalMulai:   b.TanggalMulai,
		TanggalSelesai: b.TanggalSelesai,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}
