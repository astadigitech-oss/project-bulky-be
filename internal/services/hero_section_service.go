package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type HeroSectionService interface {
	Create(ctx context.Context, req *models.CreateHeroSectionRequest) (*models.HeroSectionResponse, error)
	FindByID(ctx context.Context, id string) (*models.HeroSectionResponse, error)
	FindAll(ctx context.Context, params *models.HeroSectionFilterRequest) ([]models.HeroSectionSimpleResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateHeroSectionRequest) (*models.HeroSectionResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	Reorder(ctx context.Context, req *models.ReorderRequest) error
	GetVisibleHero(ctx context.Context) (*models.HeroSectionPublicResponse, error)
}

type heroSectionService struct {
	repo repositories.HeroSectionRepository
	cfg  *config.Config
}

func NewHeroSectionService(repo repositories.HeroSectionRepository, cfg *config.Config) HeroSectionService {
	return &heroSectionService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *heroSectionService) Create(ctx context.Context, req *models.CreateHeroSectionRequest) (*models.HeroSectionResponse, error) {
	// Auto-increment urutan
	maxUrutan, err := s.repo.GetMaxUrutan(ctx)
	if err != nil {
		return nil, err
	}

	hero := &models.HeroSection{
		ID:          uuid.New(),
		Nama:        req.Nama,
		GambarURLID: req.GambarID,
		GambarURLEN: req.GambarEN,
		Urutan:      maxUrutan + 1,
		IsActive:    req.IsActive,
	}

	if err := s.repo.Create(ctx, hero); err != nil {
		return nil, err
	}

	return s.toResponse(hero), nil
}

func (s *heroSectionService) FindByID(ctx context.Context, id string) (*models.HeroSectionResponse, error) {
	hero, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("hero section tidak ditemukan")
	}
	return s.toResponse(hero), nil
}

func (s *heroSectionService) FindAll(ctx context.Context, params *models.HeroSectionFilterRequest) ([]models.HeroSectionSimpleResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	heroes, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	items := []models.HeroSectionSimpleResponse{}
	for _, h := range heroes {
		items = append(items, *s.toSimpleResponse(&h))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *heroSectionService) Update(ctx context.Context, id string, req *models.UpdateHeroSectionRequest) (*models.HeroSectionResponse, error) {
	hero, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("hero section tidak ditemukan")
	}

	if req.Nama != nil {
		hero.Nama = *req.Nama
	}
	if req.GambarID != nil {
		hero.GambarURLID = *req.GambarID
	}
	if req.GambarEN != nil {
		hero.GambarURLEN = req.GambarEN
	}
	if req.IsActive != nil {
		hero.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, hero); err != nil {
		return nil, err
	}

	return s.toResponse(hero), nil
}

func (s *heroSectionService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("hero section tidak ditemukan")
	}
	return s.repo.Delete(ctx, id)
}

func (s *heroSectionService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	hero, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("hero section tidak ditemukan")
	}

	hero.IsActive = !hero.IsActive
	if err := s.repo.Update(ctx, hero); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       hero.ID.String(),
		IsActive: hero.IsActive,
	}, nil
}

func (s *heroSectionService) Reorder(ctx context.Context, req *models.ReorderRequest) error {
	return s.repo.UpdateOrder(ctx, req.Items)
}

func (s *heroSectionService) GetVisibleHero(ctx context.Context) (*models.HeroSectionPublicResponse, error) {
	hero, err := s.repo.GetVisibleHero(ctx)
	if err != nil {
		return nil, nil // Return nil if no visible hero
	}

	return &models.HeroSectionPublicResponse{
		ID:        hero.ID.String(),
		Nama:      hero.Nama,
		GambarURL: hero.GetGambarURL().GetFullURL(s.cfg.BaseURL),
		// LinkURL:   hero.LinkURL,
	}, nil
}

func (s *heroSectionService) toResponse(h *models.HeroSection) *models.HeroSectionResponse {
	return &models.HeroSectionResponse{
		ID:        h.ID.String(),
		Nama:      h.Nama,
		GambarURL: h.GetGambarURL().GetFullURL(s.cfg.BaseURL),
		// LinkURL:   h.LinkURL,
		Urutan:    h.Urutan,
		IsActive:  h.IsActive,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}

func (s *heroSectionService) toSimpleResponse(h *models.HeroSection) *models.HeroSectionSimpleResponse {
	return &models.HeroSectionSimpleResponse{
		ID:        h.ID.String(),
		Nama:      h.Nama,
		GambarURL: h.GetGambarURL().GetFullURL(s.cfg.BaseURL),
		// LinkURL:   h.LinkURL,
		Urutan:    h.Urutan,
		IsActive:  h.IsActive,
		UpdatedAt: h.UpdatedAt,
	}
}
