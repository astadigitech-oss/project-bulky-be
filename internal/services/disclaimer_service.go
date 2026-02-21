package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
)

type DisclaimerService interface {
	Create(ctx context.Context, req *models.CreateDisclaimerRequest) (*models.DisclaimerDetailResponse, error)
	FindByID(ctx context.Context, id string) (*models.DisclaimerDetailResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.DisclaimerListResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateDisclaimerRequest) (*models.DisclaimerDetailResponse, error)
	Delete(ctx context.Context, id string) error
	SetActive(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	GetActive(ctx context.Context, lang string) (*models.DisclaimerPublicResponse, error)
}

type disclaimerService struct {
	repo repositories.DisclaimerRepository
}

func NewDisclaimerService(repo repositories.DisclaimerRepository) DisclaimerService {
	return &disclaimerService{repo: repo}
}

func (s *disclaimerService) Create(ctx context.Context, req *models.CreateDisclaimerRequest) (*models.DisclaimerDetailResponse, error) {
	// Generate/use slug_id
	var slugID *string
	if req.SlugID != nil && *req.SlugID != "" {
		slugID = req.SlugID
	} else {
		s := utils.GenerateSlug(req.Judul)
		slugID = &s
	}

	// Generate/use slug_en
	var slugEN *string
	if req.SlugEN != nil && *req.SlugEN != "" {
		slugEN = req.SlugEN
	} else if req.JudulEn != "" {
		s := utils.GenerateSlug(req.JudulEn)
		slugEN = &s
	}

	// Check slug uniqueness
	if slugID != nil && *slugID != "" {
		exists, err := s.repo.ExistsBySlug(ctx, *slugID, nil)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("slug sudah digunakan")
		}
	}
	if slugEN != nil && *slugEN != "" {
		exists, err := s.repo.ExistsBySlug(ctx, *slugEN, nil)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("slug sudah digunakan")
		}
	}

	disclaimer := &models.Disclaimer{
		ID:       uuid.New(),
		Judul:    req.Judul,
		JudulEn:  req.JudulEn,
		Slug:     slugID,
		SlugID:   slugID,
		SlugEN:   slugEN,
		Konten:   req.Konten,
		KontenEn: req.KontenEn,
		IsActive: req.IsActive,
	}

	if err := s.repo.Create(ctx, disclaimer); err != nil {
		return nil, err
	}

	return s.toDetailResponse(disclaimer), nil
}

func (s *disclaimerService) FindByID(ctx context.Context, id string) (*models.DisclaimerDetailResponse, error) {
	disclaimer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("disclaimer tidak ditemukan")
	}
	return s.toDetailResponse(disclaimer), nil
}

func (s *disclaimerService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.DisclaimerListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	disclaimers, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	items := []models.DisclaimerListResponse{}
	for _, d := range disclaimers {
		items = append(items, *s.toListResponse(&d))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *disclaimerService) Update(ctx context.Context, id string, req *models.UpdateDisclaimerRequest) (*models.DisclaimerDetailResponse, error) {
	disclaimer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("disclaimer tidak ditemukan")
	}

	// Update judul fields
	if req.Judul != nil {
		disclaimer.Judul = *req.Judul
	}
	if req.JudulEn != nil {
		disclaimer.JudulEn = *req.JudulEn
	}

	// Handle slug_id
	if req.SlugID != nil && *req.SlugID != "" {
		disclaimer.SlugID = req.SlugID
		disclaimer.Slug = req.SlugID // backward compat
	} else if req.Judul != nil {
		// Judul changed and no explicit SlugID provided, regenerate
		s := utils.GenerateSlug(*req.Judul)
		disclaimer.SlugID = &s
		disclaimer.Slug = &s // backward compat
	}

	// Handle slug_en
	if req.SlugEN != nil && *req.SlugEN != "" {
		disclaimer.SlugEN = req.SlugEN
	} else if req.JudulEn != nil {
		// JudulEn changed and no explicit SlugEN provided, regenerate
		s := utils.GenerateSlug(*req.JudulEn)
		disclaimer.SlugEN = &s
	}

	// Check slug uniqueness
	if disclaimer.SlugID != nil && *disclaimer.SlugID != "" {
		exists, err := s.repo.ExistsBySlug(ctx, *disclaimer.SlugID, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("slug sudah digunakan")
		}
	}
	if disclaimer.SlugEN != nil && *disclaimer.SlugEN != "" {
		exists, err := s.repo.ExistsBySlug(ctx, *disclaimer.SlugEN, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("slug sudah digunakan")
		}
	}

	if req.Konten != nil {
		disclaimer.Konten = *req.Konten
	}
	if req.KontenEn != nil {
		disclaimer.KontenEn = *req.KontenEn
	}
	if req.IsActive != nil {
		disclaimer.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, disclaimer); err != nil {
		return nil, err
	}

	return s.toDetailResponse(disclaimer), nil
}

func (s *disclaimerService) Delete(ctx context.Context, id string) error {
	disclaimer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("disclaimer tidak ditemukan")
	}
	return s.repo.Delete(ctx, disclaimer)
}

func (s *disclaimerService) SetActive(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	disclaimer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("disclaimer tidak ditemukan")
	}

	// Set as active (trigger will deactivate others)
	disclaimer.IsActive = true
	if err := s.repo.Update(ctx, disclaimer); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:        disclaimer.ID.String(),
		IsActive:  disclaimer.IsActive,
		UpdatedAt: disclaimer.UpdatedAt,
	}, nil
}

func (s *disclaimerService) GetActive(ctx context.Context, lang string) (*models.DisclaimerPublicResponse, error) {
	disclaimer, err := s.repo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	if disclaimer == nil {
		return nil, nil
	}

	// Return based on language preference
	judul := disclaimer.Judul
	konten := disclaimer.Konten

	if lang == "en" {
		judul = disclaimer.JudulEn
		konten = disclaimer.KontenEn
	}

	return &models.DisclaimerPublicResponse{
		Judul:  judul,
		SlugID: disclaimer.SlugID,
		SlugEN: disclaimer.SlugEN,
		Konten: konten,
	}, nil
}

func (s *disclaimerService) toListResponse(d *models.Disclaimer) *models.DisclaimerListResponse {
	// slug := ""
	// if d.Slug != nil {
	// 	slug = *d.Slug
	// }

	return &models.DisclaimerListResponse{
		ID: d.ID.String(),
		Judul: models.TranslatableString{
			ID: d.Judul,
			EN: &d.JudulEn,
		},
		// Slug:      slug,
		IsActive: d.IsActive,
		// CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (s *disclaimerService) toDetailResponse(d *models.Disclaimer) *models.DisclaimerDetailResponse {
	return &models.DisclaimerDetailResponse{
		ID: d.ID.String(),
		Judul: models.TranslatableString{
			ID: d.Judul,
			EN: &d.JudulEn,
		},
		SlugID: d.SlugID,
		SlugEN: d.SlugEN,
		Konten: models.TranslatableString{
			ID: d.Konten,
			EN: &d.KontenEn,
		},
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
