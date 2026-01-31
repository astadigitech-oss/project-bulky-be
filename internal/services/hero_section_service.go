package services

import (
	"context"
	"errors"
	"time"

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
	ToggleStatus(ctx context.Context, id string) (*models.ToggleDefaultResponse, error)
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

func (s *heroSectionService) parseDateString(dateStr *string) (*time.Time, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}

	// Try parsing with different formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, *dateStr); err == nil {
			return &t, nil
		}
	}

	return nil, errors.New("format tanggal tidak valid")
}

func (s *heroSectionService) validateDateRange(tanggalMulai, tanggalSelesai *time.Time, excludeID *string) error {
	// Skip validation if no date range
	if tanggalMulai == nil || tanggalSelesai == nil {
		return nil
	}

	// Check if end date is after start date
	if tanggalSelesai.Before(*tanggalMulai) {
		return errors.New("tanggal selesai harus setelah tanggal mulai")
	}

	// Check for overlapping records
	hasOverlap, err := s.repo.CheckDateRangeOverlap(context.Background(), tanggalMulai, tanggalSelesai, excludeID)
	if err != nil {
		return err
	}

	if hasOverlap {
		return errors.New("sudah ada hero section dengan rentang tanggal yang overlap")
	}

	return nil
}

func (s *heroSectionService) isVisible(hero *models.HeroSection) bool {
	// Visible if is_default = true
	if hero.IsDefault {
		return true
	}

	// Visible if NOW() is within date range
	if hero.TanggalMulai != nil && hero.TanggalSelesai != nil {
		now := time.Now()
		return !now.Before(*hero.TanggalMulai) && !now.After(*hero.TanggalSelesai)
	}

	return false
}

func (s *heroSectionService) Create(ctx context.Context, req *models.CreateHeroSectionRequest) (*models.HeroSectionResponse, error) {
	// Parse dates
	tanggalMulai, err := s.parseDateString(req.TanggalMulai)
	if err != nil {
		return nil, errors.New("format tanggal_mulai tidak valid")
	}

	tanggalSelesai, err := s.parseDateString(req.TanggalSelesai)
	if err != nil {
		return nil, errors.New("format tanggal_selesai tidak valid")
	}

	// Validate date range overlap
	if err := s.validateDateRange(tanggalMulai, tanggalSelesai, nil); err != nil {
		return nil, err
	}

	hero := &models.HeroSection{
		ID:             uuid.New(),
		Nama:           req.Nama,
		GambarURLID:    req.GambarID,
		GambarURLEN:    req.GambarEN,
		IsDefault:      req.IsDefault,
		TanggalMulai:   tanggalMulai,
		TanggalSelesai: tanggalSelesai,
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

	// Parse dates if provided
	var tanggalMulai, tanggalSelesai *time.Time
	if req.TanggalMulai != nil {
		tanggalMulai, err = s.parseDateString(req.TanggalMulai)
		if err != nil {
			return nil, errors.New("format tanggal_mulai tidak valid")
		}
	} else {
		tanggalMulai = hero.TanggalMulai
	}

	if req.TanggalSelesai != nil {
		tanggalSelesai, err = s.parseDateString(req.TanggalSelesai)
		if err != nil {
			return nil, errors.New("format tanggal_selesai tidak valid")
		}
	} else {
		tanggalSelesai = hero.TanggalSelesai
	}

	// Validate date range overlap (exclude current record)
	idStr := id
	if err := s.validateDateRange(tanggalMulai, tanggalSelesai, &idStr); err != nil {
		return nil, err
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
	if req.IsDefault != nil {
		hero.IsDefault = *req.IsDefault
	}
	if req.TanggalMulai != nil {
		hero.TanggalMulai = tanggalMulai
	}
	if req.TanggalSelesai != nil {
		hero.TanggalSelesai = tanggalSelesai
	}

	if err := s.repo.Update(ctx, hero); err != nil {
		return nil, err
	}

	return s.toResponse(hero), nil
}

func (s *heroSectionService) Delete(ctx context.Context, id string) error {
	hero, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("hero section tidak ditemukan")
	}

	// Soft delete
	return s.repo.Delete(ctx, hero.ID.String())
}

func (s *heroSectionService) ToggleStatus(ctx context.Context, id string) (*models.ToggleDefaultResponse, error) {
	hero, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("hero section tidak ditemukan")
	}

	hero.IsDefault = !hero.IsDefault
	if err := s.repo.Update(ctx, hero); err != nil {
		return nil, err
	}

	return &models.ToggleDefaultResponse{
		ID:        hero.ID.String(),
		IsDefault: hero.IsDefault,
	}, nil
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
	}, nil
}

func (s *heroSectionService) toResponse(h *models.HeroSection) *models.HeroSectionResponse {
	return &models.HeroSectionResponse{
		ID:             h.ID.String(),
		Nama:           h.Nama,
		GambarURL:      h.GetGambarURL().GetFullURL(s.cfg.BaseURL),
		IsDefault:      h.IsDefault,
		IsVisible:      s.isVisible(h),
		TanggalMulai:   h.TanggalMulai,
		TanggalSelesai: h.TanggalSelesai,
		CreatedAt:      h.CreatedAt,
		UpdatedAt:      h.UpdatedAt,
	}
}

func (s *heroSectionService) toSimpleResponse(h *models.HeroSection) *models.HeroSectionSimpleResponse {
	return &models.HeroSectionSimpleResponse{
		ID:             h.ID.String(),
		Nama:           h.Nama,
		GambarURL:      h.GetGambarURL().GetFullURL(s.cfg.BaseURL),
		IsDefault:      h.IsDefault,
		IsVisible:      s.isVisible(h),
		TanggalMulai:   h.TanggalMulai,
		TanggalSelesai: h.TanggalSelesai,
		UpdatedAt:      h.UpdatedAt,
	}
}
