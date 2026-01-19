package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"gorm.io/gorm"
)

type DokumenKebijakanService interface {
	FindAll(ctx context.Context) ([]models.DokumenKebijakanListResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.DokumenKebijakanDetailResponse, error)
	Update(ctx context.Context, slug string, req *models.UpdateDokumenKebijakanRequest) (*models.DokumenKebijakanDetailResponse, error)
	GetBySlugPublic(ctx context.Context, slug string) (*models.DokumenKebijakanPublicResponse, error)
	GetActiveListPublic(ctx context.Context) ([]models.DokumenKebijakanPublicResponse, error)
}

type dokumenKebijakanService struct {
	repo      repositories.DokumenKebijakanRepository
	sanitizer *HTMLSanitizer
}

func NewDokumenKebijakanService(repo repositories.DokumenKebijakanRepository) DokumenKebijakanService {
	return &dokumenKebijakanService{
		repo:      repo,
		sanitizer: NewHTMLSanitizer(),
	}
}

func (s *dokumenKebijakanService) FindAll(ctx context.Context) ([]models.DokumenKebijakanListResponse, error) {
	dokumens, err := s.repo.FindAllSimple(ctx)
	if err != nil {
		return nil, errors.New("gagal mengambil data dokumen kebijakan")
	}

	items := make([]models.DokumenKebijakanListResponse, len(dokumens))
	for i, d := range dokumens {
		items[i] = models.DokumenKebijakanListResponse{
			ID:        d.ID.String(),
			Judul:     d.Judul,
			Slug:      d.Slug,
			Urutan:    d.Urutan,
			IsActive:  d.IsActive,
			UpdatedAt: d.UpdatedAt,
		}
	}

	return items, nil
}

func (s *dokumenKebijakanService) FindBySlug(ctx context.Context, slug string) (*models.DokumenKebijakanDetailResponse, error) {
	dokumen, err := s.repo.FindBySlugForEdit(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dokumen kebijakan tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data dokumen kebijakan")
	}

	return &models.DokumenKebijakanDetailResponse{
		ID:        dokumen.ID.String(),
		Judul:     dokumen.Judul,
		Slug:      dokumen.Slug,
		Konten:    dokumen.Konten,
		Urutan:    dokumen.Urutan,
		IsActive:  dokumen.IsActive,
		CreatedAt: dokumen.CreatedAt,
		UpdatedAt: dokumen.UpdatedAt,
	}, nil
}

func (s *dokumenKebijakanService) Update(ctx context.Context, slug string, req *models.UpdateDokumenKebijakanRequest) (*models.DokumenKebijakanDetailResponse, error) {
	dokumen, err := s.repo.FindBySlugForEdit(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dokumen kebijakan tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data dokumen kebijakan")
	}

	// Update judul if provided
	if req.Judul != nil {
		dokumen.Judul = *req.Judul
	}

	// Update konten if provided (sanitize HTML)
	if req.Konten != nil {
		sanitized := s.sanitizer.Sanitize(*req.Konten)
		dokumen.Konten = sanitized
	}

	// Update is_active if provided
	if req.IsActive != nil {
		dokumen.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, dokumen); err != nil {
		return nil, errors.New("gagal mengupdate dokumen kebijakan")
	}

	return &models.DokumenKebijakanDetailResponse{
		ID:        dokumen.ID.String(),
		Judul:     dokumen.Judul,
		Slug:      dokumen.Slug,
		Konten:    dokumen.Konten,
		Urutan:    dokumen.Urutan,
		IsActive:  dokumen.IsActive,
		CreatedAt: dokumen.CreatedAt,
		UpdatedAt: dokumen.UpdatedAt,
	}, nil
}

func (s *dokumenKebijakanService) GetBySlugPublic(ctx context.Context, slug string) (*models.DokumenKebijakanPublicResponse, error) {
	dokumen, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dokumen kebijakan tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data dokumen kebijakan")
	}

	return &models.DokumenKebijakanPublicResponse{
		Judul:  dokumen.Judul,
		Slug:   dokumen.Slug,
		Konten: dokumen.Konten,
	}, nil
}

func (s *dokumenKebijakanService) GetActiveListPublic(ctx context.Context) ([]models.DokumenKebijakanPublicResponse, error) {
	dokumens, err := s.repo.GetActiveList(ctx)
	if err != nil {
		return nil, errors.New("gagal mengambil data dokumen kebijakan")
	}

	items := make([]models.DokumenKebijakanPublicResponse, len(dokumens))
	for i, d := range dokumens {
		items[i] = models.DokumenKebijakanPublicResponse{
			Judul: d.Judul,
			Slug:  d.Slug,
		}
	}

	return items, nil
}
