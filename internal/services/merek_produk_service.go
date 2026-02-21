package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type MerekProdukService interface {
	Create(ctx context.Context, req *models.CreateMerekProdukRequest) (*models.MerekProdukResponse, error)
	CreateWithLogo(ctx context.Context, req *models.CreateMerekProdukRequest, logoURL *string) (*models.MerekProdukResponse, error)
	FindByID(ctx context.Context, id string) (*models.MerekProdukResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.MerekProdukResponse, error)
	FindAll(ctx context.Context, params *models.MerekProdukFilterRequest) ([]models.MerekProdukSimpleResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateMerekProdukRequest) (*models.MerekProdukResponse, error)
	UpdateWithLogo(ctx context.Context, id string, req *models.UpdateMerekProdukRequest, logoURL *string) (*models.MerekProdukResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	GetAllForDropdown(ctx context.Context) ([]map[string]interface{}, error)
}

type merekProdukService struct {
	repo repositories.MerekProdukRepository
	cfg  *config.Config
}

func NewMerekProdukService(repo repositories.MerekProdukRepository, cfg *config.Config) MerekProdukService {
	return &merekProdukService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *merekProdukService) Create(ctx context.Context, req *models.CreateMerekProdukRequest) (*models.MerekProdukResponse, error) {
	// Generate slug_id
	var slugID *string
	if req.SlugID != nil && *req.SlugID != "" {
		s := *req.SlugID
		slugID = &s
	} else {
		s := utils.GenerateSlug(req.NamaID)
		slugID = &s
	}

	// Generate slug_en
	var slugEN *string
	if req.SlugEN != nil && *req.SlugEN != "" {
		s := *req.SlugEN
		slugEN = &s
	} else if req.NamaEN != nil && *req.NamaEN != "" {
		s := utils.GenerateSlug(*req.NamaEN)
		slugEN = &s
	}

	merek := &models.MerekProduk{
		NamaID:   req.NamaID,
		NamaEN:   req.NamaEN,
		Slug:     *slugID,
		SlugID:   slugID,
		SlugEN:   slugEN,
		LogoURL:  req.Logo,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, merek); err != nil {
		return nil, err
	}

	return s.toResponse(merek), nil
}

func (s *merekProdukService) CreateWithLogo(ctx context.Context, req *models.CreateMerekProdukRequest, logoURL *string) (*models.MerekProdukResponse, error) {
	// Generate slug_id
	var slugID *string
	if req.SlugID != nil && *req.SlugID != "" {
		s := *req.SlugID
		slugID = &s
	} else {
		s := utils.GenerateSlug(req.NamaID)
		slugID = &s
	}

	// Generate slug_en
	var slugEN *string
	if req.SlugEN != nil && *req.SlugEN != "" {
		s := *req.SlugEN
		slugEN = &s
	} else if req.NamaEN != nil && *req.NamaEN != "" {
		s := utils.GenerateSlug(*req.NamaEN)
		slugEN = &s
	}

	merek := &models.MerekProduk{
		NamaID:   req.NamaID,
		NamaEN:   req.NamaEN,
		Slug:     *slugID,
		SlugID:   slugID,
		SlugEN:   slugEN,
		LogoURL:  logoURL,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, merek); err != nil {
		return nil, err
	}

	return s.toResponse(merek), nil
}

func (s *merekProdukService) FindByID(ctx context.Context, id string) (*models.MerekProdukResponse, error) {
	merek, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("merek produk tidak ditemukan")
	}
	return s.toResponse(merek), nil
}

func (s *merekProdukService) FindBySlug(ctx context.Context, slug string) (*models.MerekProdukResponse, error) {
	merek, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("merek produk tidak ditemukan")
	}
	return s.toResponse(merek), nil
}

func (s *merekProdukService) FindAll(ctx context.Context, params *models.MerekProdukFilterRequest) ([]models.MerekProdukSimpleResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	mereks, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	// Add BASE_URL to LogoURL
	for i := range mereks {
		mereks[i].LogoURL = utils.GetFileURLPtr(mereks[i].LogoURL, s.cfg)
	}

	// Ensure empty array instead of null
	if mereks == nil {
		mereks = []models.MerekProdukSimpleResponse{}
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return mereks, &meta, nil
}

func (s *merekProdukService) Update(ctx context.Context, id string, req *models.UpdateMerekProdukRequest) (*models.MerekProdukResponse, error) {
	merek, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("merek produk tidak ditemukan")
	}

	if req.NamaID != nil {
		merek.NamaID = *req.NamaID
		// Regenerate slug_id from new nama_id (unless manually provided)
		if req.SlugID == nil || *req.SlugID == "" {
			s := utils.GenerateSlug(*req.NamaID)
			merek.SlugID = &s
			merek.Slug = s // backward compat
		}
	}
	if req.SlugID != nil && *req.SlugID != "" {
		merek.SlugID = req.SlugID
		merek.Slug = *req.SlugID // backward compat
	}
	if req.NamaEN != nil {
		merek.NamaEN = req.NamaEN
		// Regenerate slug_en from new nama_en (unless manually provided)
		if (req.SlugEN == nil || *req.SlugEN == "") && *req.NamaEN != "" {
			s := utils.GenerateSlug(*req.NamaEN)
			merek.SlugEN = &s
		}
	}
	if req.SlugEN != nil && *req.SlugEN != "" {
		merek.SlugEN = req.SlugEN
	}
	if req.Logo != nil {
		merek.LogoURL = req.Logo
	}
	if req.IsActive != nil {
		merek.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, merek); err != nil {
		return nil, err
	}

	return s.toResponse(merek), nil
}

func (s *merekProdukService) UpdateWithLogo(ctx context.Context, id string, req *models.UpdateMerekProdukRequest, logoURL *string) (*models.MerekProdukResponse, error) {
	merek, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("merek produk tidak ditemukan")
	}

	// Update text fields
	if req.NamaID != nil {
		merek.NamaID = *req.NamaID
		// Regenerate slug_id from new nama_id (unless manually provided)
		if req.SlugID == nil || *req.SlugID == "" {
			s := utils.GenerateSlug(*req.NamaID)
			merek.SlugID = &s
			merek.Slug = s // backward compat
		}
	}
	if req.SlugID != nil && *req.SlugID != "" {
		merek.SlugID = req.SlugID
		merek.Slug = *req.SlugID // backward compat
	}
	if req.NamaEN != nil {
		merek.NamaEN = req.NamaEN
		// Regenerate slug_en from new nama_en (unless manually provided)
		if (req.SlugEN == nil || *req.SlugEN == "") && *req.NamaEN != "" {
			s := utils.GenerateSlug(*req.NamaEN)
			merek.SlugEN = &s
		}
	}
	if req.SlugEN != nil && *req.SlugEN != "" {
		merek.SlugEN = req.SlugEN
	}
	if req.IsActive != nil {
		merek.IsActive = *req.IsActive
	}

	// Update logo if provided
	if logoURL != nil {
		// Delete old logo if exists
		if merek.LogoURL != nil && *merek.LogoURL != "" {
			if err := utils.DeleteFile(*merek.LogoURL, s.cfg); err == nil {
				// Log error but don't fail the update
			}
		}
		merek.LogoURL = logoURL
	}

	if err := s.repo.Update(ctx, merek); err != nil {
		return nil, err
	}

	return s.toResponse(merek), nil
}

func (s *merekProdukService) Delete(ctx context.Context, id string) error {
	merek, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("merek produk tidak ditemukan")
	}

	// Rename slug dengan suffix _deleted_{8-char-id} agar tidak conflict unique constraint
	suffix := "_deleted_" + id[:8]
	merek.Slug = merek.Slug + suffix
	if merek.SlugID != nil {
		v := *merek.SlugID + suffix
		merek.SlugID = &v
	}
	if merek.SlugEN != nil {
		v := *merek.SlugEN + suffix
		merek.SlugEN = &v
	}
	if err := s.repo.Update(ctx, merek); err != nil {
		return err
	}

	return s.repo.Delete(ctx, merek)
}

func (s *merekProdukService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	merek, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("merek produk tidak ditemukan")
	}

	merek.IsActive = !merek.IsActive
	if err := s.repo.Update(ctx, merek); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:        merek.ID.String(),
		IsActive:  merek.IsActive,
		UpdatedAt: merek.UpdatedAt,
	}, nil
}

func (s *merekProdukService) GetAllForDropdown(ctx context.Context) ([]map[string]interface{}, error) {
	mereks, err := s.repo.GetAllForDropdown(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(mereks))
	for i, m := range mereks {
		result[i] = map[string]interface{}{
			"id":      m.ID.String(),
			"nama_id": m.NamaID,
			"nama_en": m.NamaEN,
		}
	}
	return result, nil
}

func (s *merekProdukService) toResponse(m *models.MerekProduk) *models.MerekProdukResponse {
	return &models.MerekProdukResponse{
		ID:        m.ID.String(),
		Nama:      m.GetNama(),
		SlugID:    m.SlugID,
		SlugEN:    m.SlugEN,
		LogoURL:   utils.GetFileURLPtr(m.LogoURL, s.cfg),
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
