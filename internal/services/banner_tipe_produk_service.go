package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type BannerTipeProdukService interface {
	Create(ctx context.Context, req *models.CreateBannerTipeProdukRequest) (*models.BannerTipeProdukResponse, error)
	FindByID(ctx context.Context, id string) (*models.BannerTipeProdukResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest, tipeProdukID string) ([]models.BannerTipeProdukResponse, *models.PaginationMeta, error)
	FindByTipeProdukID(ctx context.Context, tipeProdukID string) ([]models.BannerSimpleResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateBannerTipeProdukRequest) (*models.BannerTipeProdukResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	Reorder(ctx context.Context, req *models.ReorderRequest) error
}

type bannerTipeProdukService struct {
	repo repositories.BannerTipeProdukRepository
}

func NewBannerTipeProdukService(repo repositories.BannerTipeProdukRepository) BannerTipeProdukService {
	return &bannerTipeProdukService{repo: repo}
}

func (s *bannerTipeProdukService) Create(ctx context.Context, req *models.CreateBannerTipeProdukRequest) (*models.BannerTipeProdukResponse, error) {
	tipeProdukUUID, err := uuid.Parse(req.TipeProdukID)
	if err != nil {
		return nil, errors.New("tipe_produk_id tidak valid")
	}

	banner := &models.BannerTipeProduk{
		TipeProdukID: tipeProdukUUID,
		Nama:         req.Nama,
		GambarURL:    req.GambarURL,
		IsActive:     true,
	}

	if req.Urutan != nil {
		banner.Urutan = *req.Urutan
	}

	if err := s.repo.Create(ctx, banner); err != nil {
		return nil, err
	}

	return s.FindByID(ctx, banner.ID.String())
}

func (s *bannerTipeProdukService) FindByID(ctx context.Context, id string) (*models.BannerTipeProdukResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}
	return s.toResponse(banner), nil
}

func (s *bannerTipeProdukService) FindAll(ctx context.Context, params *models.PaginationRequest, tipeProdukID string) ([]models.BannerTipeProdukResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	banners, total, err := s.repo.FindAll(ctx, params, tipeProdukID)
	if err != nil {
		return nil, nil, err
	}

	var items []models.BannerTipeProdukResponse
	for _, b := range banners {
		items = append(items, *s.toResponse(&b))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *bannerTipeProdukService) FindByTipeProdukID(ctx context.Context, tipeProdukID string) ([]models.BannerSimpleResponse, error) {
	_, err := uuid.Parse(tipeProdukID)
	if err != nil {
		return nil, errors.New("tipe_produk_id tidak valid")
	}

	banners, err := s.repo.FindByTipeProdukID(ctx, tipeProdukID)
	if err != nil {
		return nil, err
	}

	var items []models.BannerSimpleResponse
	for _, b := range banners {
		items = append(items, models.BannerSimpleResponse{
			ID:        b.ID.String(),
			Nama:      b.Nama,
			GambarURL: b.GambarURL,
			Urutan:    b.Urutan,
		})
	}

	return items, nil
}

func (s *bannerTipeProdukService) Update(ctx context.Context, id string, req *models.UpdateBannerTipeProdukRequest) (*models.BannerTipeProdukResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}

	if req.TipeProdukID != nil {
		tipeProdukUUID, err := uuid.Parse(*req.TipeProdukID)
		if err != nil {
			return nil, errors.New("tipe_produk_id tidak valid")
		}
		banner.TipeProdukID = tipeProdukUUID
	}
	if req.Nama != nil {
		banner.Nama = *req.Nama
	}
	if req.GambarURL != nil {
		banner.GambarURL = *req.GambarURL
	}
	if req.Urutan != nil {
		banner.Urutan = *req.Urutan
	}
	if req.IsActive != nil {
		banner.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, banner); err != nil {
		return nil, err
	}

	return s.FindByID(ctx, id)
}

func (s *bannerTipeProdukService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("banner tidak ditemukan")
	}

	return s.repo.Delete(ctx, id)
}

func (s *bannerTipeProdukService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
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

func (s *bannerTipeProdukService) Reorder(ctx context.Context, req *models.ReorderRequest) error {
	return s.repo.UpdateOrder(ctx, req.Items)
}

func (s *bannerTipeProdukService) toResponse(b *models.BannerTipeProduk) *models.BannerTipeProdukResponse {
	return &models.BannerTipeProdukResponse{
		ID: b.ID.String(),
		TipeProduk: models.BannerTipeProdukTipeInfo{
			ID:   b.TipeProduk.ID.String(),
			Nama: b.TipeProduk.Nama,
			Slug: b.TipeProduk.Slug,
		},
		Nama:      b.Nama,
		GambarURL: b.GambarURL,
		Urutan:    b.Urutan,
		IsActive:  b.IsActive,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
