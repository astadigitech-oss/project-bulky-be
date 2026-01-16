package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

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
	// Check slug uniqueness if provided
	if req.Slug != nil && *req.Slug != "" {
		exists, err := s.repo.ExistsBySlug(ctx, *req.Slug, nil)
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
		Slug:     req.Slug,
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

	// Check slug uniqueness if provided
	if req.Slug != nil && *req.Slug != "" {
		exists, err := s.repo.ExistsBySlug(ctx, *req.Slug, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("slug sudah digunakan")
		}
	}

	if req.Judul != nil {
		disclaimer.Judul = *req.Judul
	}
	if req.JudulEn != nil {
		disclaimer.JudulEn = *req.JudulEn
	}
	if req.Slug != nil {
		disclaimer.Slug = req.Slug
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
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("disclaimer tidak ditemukan")
	}
	return s.repo.Delete(ctx, id)
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
		ID:       disclaimer.ID.String(),
		IsActive: disclaimer.IsActive,
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
		Slug:   *disclaimer.Slug,
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
	slug := ""
	if d.Slug != nil {
		slug = *d.Slug
	}

	return &models.DisclaimerDetailResponse{
		ID: d.ID.String(),
		Judul: models.TranslatableString{
			ID: d.Judul,
			EN: &d.JudulEn,
		},
		Slug: slug,
		Konten: models.TranslatableString{
			ID: d.Konten,
			EN: &d.KontenEn,
		},
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
