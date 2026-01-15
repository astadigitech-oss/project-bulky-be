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
	FindAll(ctx context.Context, params *models.BannerEventPromoFilterRequest) ([]models.BannerEventPromoSimpleResponse, *models.PaginationMeta, error)
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
		ID:          uuid.New(),
		Nama:        req.Nama,
		GambarURLID: req.GambarID,
		GambarURLEN: req.GambarEN,
		LinkURL:     req.UrlTujuan,
		Urutan:      req.Urutan,
		IsActive:    req.IsActive,
	}

	if req.TanggalMulai != nil {
		if t, err := parseFlexibleDate(*req.TanggalMulai); err == nil {
			banner.TanggalMulai = &t
		}
	}

	if req.TanggalSelesai != nil {
		if t, err := parseFlexibleDate(*req.TanggalSelesai); err == nil {
			banner.TanggalAkhir = &t
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

func (s *bannerEventPromoService) FindAll(ctx context.Context, params *models.BannerEventPromoFilterRequest) ([]models.BannerEventPromoSimpleResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	banners, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	items := []models.BannerEventPromoSimpleResponse{}
	for _, b := range banners {
		items = append(items, *s.toSimpleResponse(&b))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *bannerEventPromoService) Update(ctx context.Context, id string, req *models.UpdateBannerEventPromoRequest) (*models.BannerEventPromoResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}

	if req.Nama != nil {
		banner.Nama = *req.Nama
	}
	if req.GambarID != nil {
		banner.GambarURLID = *req.GambarID
	}
	if req.GambarEN != nil {
		banner.GambarURLEN = req.GambarEN
	}
	if req.UrlTujuan != nil {
		banner.LinkURL = req.UrlTujuan
	}
	if req.Urutan != nil {
		banner.Urutan = *req.Urutan
	}
	if req.IsActive != nil {
		banner.IsActive = *req.IsActive
	}
	if req.TanggalMulai != nil {
		if t, err := parseFlexibleDate(*req.TanggalMulai); err == nil {
			banner.TanggalMulai = &t
		}
	}
	if req.TanggalSelesai != nil {
		if t, err := parseFlexibleDate(*req.TanggalSelesai); err == nil {
			banner.TanggalAkhir = &t
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

	items := []models.BannerEventPromoPublicResponse{}
	for _, b := range banners {
		items = append(items, models.BannerEventPromoPublicResponse{
			ID:        b.ID.String(),
			Nama:      b.Nama,
			GambarURL: b.GetGambarURL(),
			LinkURL:   b.LinkURL,
		})
	}

	return items, nil
}

func (s *bannerEventPromoService) toResponse(b *models.BannerEventPromo) *models.BannerEventPromoResponse {
	return &models.BannerEventPromoResponse{
		ID:           b.ID.String(),
		Nama:         b.Nama,
		GambarURL:    b.GetGambarURL(),
		LinkURL:      b.LinkURL,
		Urutan:       b.Urutan,
		IsActive:     b.IsActive,
		TanggalMulai: b.TanggalMulai,
		TanggalAkhir: b.TanggalAkhir,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
}

func (s *bannerEventPromoService) toSimpleResponse(b *models.BannerEventPromo) *models.BannerEventPromoSimpleResponse {
	return &models.BannerEventPromoSimpleResponse{
		ID:        b.ID.String(),
		Nama:      b.Nama,
		GambarURL: b.GetGambarURL(),
		LinkURL:   b.LinkURL,
		Urutan:    b.Urutan,
		IsActive:  b.IsActive,
		UpdatedAt: b.UpdatedAt,
	}
}

// parseFlexibleDate parses date string in multiple formats
// Supports: "2026-01-10" (date only) or "2026-01-10T00:00:00Z" (RFC3339)
func parseFlexibleDate(dateStr string) (time.Time, error) {
	// Try RFC3339 first (full datetime with timezone)
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t, nil
	}

	// Try date only format (yyyy-mm-dd)
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}

	// Try datetime without timezone
	if t, err := time.Parse("2006-01-02T15:04:05", dateStr); err == nil {
		return t, nil
	}

	return time.Time{}, errors.New("format tanggal tidak valid, gunakan yyyy-mm-dd atau RFC3339")
}
