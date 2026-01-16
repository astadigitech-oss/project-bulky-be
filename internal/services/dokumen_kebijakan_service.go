package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type DokumenKebijakanService interface {
	Create(ctx context.Context, req *models.CreateDokumenKebijakanRequest) (*models.DokumenKebijakanDetailResponse, error)
	FindByID(ctx context.Context, id string) (*models.DokumenKebijakanDetailResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.DokumenKebijakanListResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateDokumenKebijakanRequest) (*models.DokumenKebijakanDetailResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	GetBySlug(ctx context.Context, slug string, lang string) (*models.DokumenKebijakanPublicResponse, error)
	GetActiveList(ctx context.Context, lang string) ([]models.DokumenKebijakanPublicResponse, error)
}

type dokumenKebijakanService struct {
	repo repositories.DokumenKebijakanRepository
}

func NewDokumenKebijakanService(repo repositories.DokumenKebijakanRepository) DokumenKebijakanService {
	return &dokumenKebijakanService{repo: repo}
}

func (s *dokumenKebijakanService) Create(ctx context.Context, req *models.CreateDokumenKebijakanRequest) (*models.DokumenKebijakanDetailResponse, error) {
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

	dokumen := &models.DokumenKebijakan{
		ID:       uuid.New(),
		Judul:    req.Judul,
		JudulEn:  req.JudulEn,
		Slug:     req.Slug,
		Konten:   req.Konten,
		KontenEn: req.KontenEn,
		IsActive: req.IsActive,
	}

	if err := s.repo.Create(ctx, dokumen); err != nil {
		return nil, err
	}

	return s.toDetailResponse(dokumen), nil
}

func (s *dokumenKebijakanService) FindByID(ctx context.Context, id string) (*models.DokumenKebijakanDetailResponse, error) {
	dokumen, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("dokumen kebijakan tidak ditemukan")
	}
	return s.toDetailResponse(dokumen), nil
}

func (s *dokumenKebijakanService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.DokumenKebijakanListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	dokumens, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	items := []models.DokumenKebijakanListResponse{}
	for _, d := range dokumens {
		items = append(items, *s.toListResponse(&d))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *dokumenKebijakanService) Update(ctx context.Context, id string, req *models.UpdateDokumenKebijakanRequest) (*models.DokumenKebijakanDetailResponse, error) {
	dokumen, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("dokumen kebijakan tidak ditemukan")
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
		dokumen.Judul = *req.Judul
	}
	if req.JudulEn != nil {
		dokumen.JudulEn = *req.JudulEn
	}
	if req.Slug != nil {
		dokumen.Slug = req.Slug
	}
	if req.Konten != nil {
		dokumen.Konten = *req.Konten
	}
	if req.KontenEn != nil {
		dokumen.KontenEn = *req.KontenEn
	}
	if req.IsActive != nil {
		dokumen.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, dokumen); err != nil {
		return nil, err
	}

	return s.toDetailResponse(dokumen), nil
}

func (s *dokumenKebijakanService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("dokumen kebijakan tidak ditemukan")
	}
	return s.repo.Delete(ctx, id)
}

func (s *dokumenKebijakanService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	dokumen, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("dokumen kebijakan tidak ditemukan")
	}

	dokumen.IsActive = !dokumen.IsActive
	if err := s.repo.Update(ctx, dokumen); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       dokumen.ID.String(),
		IsActive: dokumen.IsActive,
	}, nil
}

func (s *dokumenKebijakanService) GetBySlug(ctx context.Context, slug string, lang string) (*models.DokumenKebijakanPublicResponse, error) {
	dokumen, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("dokumen kebijakan tidak ditemukan")
	}

	// Return based on language preference
	judul := dokumen.Judul
	konten := dokumen.Konten

	if lang == "en" {
		judul = dokumen.JudulEn
		konten = dokumen.KontenEn
	}

	return &models.DokumenKebijakanPublicResponse{
		Judul:  judul,
		Slug:   *dokumen.Slug,
		Konten: konten,
	}, nil
}

func (s *dokumenKebijakanService) GetActiveList(ctx context.Context, lang string) ([]models.DokumenKebijakanPublicResponse, error) {
	dokumens, err := s.repo.GetActiveList(ctx)
	if err != nil {
		return nil, err
	}

	items := []models.DokumenKebijakanPublicResponse{}
	for _, d := range dokumens {
		judul := d.Judul
		if lang == "en" {
			judul = d.JudulEn
		}

		items = append(items, models.DokumenKebijakanPublicResponse{
			Judul: judul,
			Slug:  *d.Slug,
		})
	}

	return items, nil
}

func (s *dokumenKebijakanService) toListResponse(d *models.DokumenKebijakan) *models.DokumenKebijakanListResponse {
	// slug := ""
	// if d.Slug != nil {
	// 	slug = *d.Slug
	// }

	return &models.DokumenKebijakanListResponse{
		ID: d.ID.String(),
		Judul: models.TranslatableString{
			ID: d.Judul,
			EN: &d.JudulEn,
		},
		// Slug:      slug,
		IsActive:  d.IsActive,
		UpdatedAt: d.UpdatedAt,
	}
}

func (s *dokumenKebijakanService) toDetailResponse(d *models.DokumenKebijakan) *models.DokumenKebijakanDetailResponse {
	slug := ""
	if d.Slug != nil {
		slug = *d.Slug
	}

	return &models.DokumenKebijakanDetailResponse{
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
