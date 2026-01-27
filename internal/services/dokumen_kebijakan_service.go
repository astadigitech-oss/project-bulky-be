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
	FindByID(ctx context.Context, id string) (*models.DokumenKebijakanDetailResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.DokumenKebijakanDetailResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateDokumenKebijakanRequest) (*models.DokumenKebijakanDetailResponse, error)
	UpdateBySlug(ctx context.Context, slug string, req *models.UpdateDokumenKebijakanRequest) (*models.DokumenKebijakanDetailResponse, error)
	GetByIDPublic(ctx context.Context, id string, lang string) (*models.DokumenKebijakanPublicResponse, error)
	GetBySlugPublic(ctx context.Context, slug string, lang string) (*models.DokumenKebijakanPublicResponse, error)
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
			ID:      d.ID.String(),
			Judul:   d.Judul,
			JudulEN: d.JudulEN,
			Slug:    d.Slug,
			Urutan:  d.Urutan,
			// IsActive:  d.IsActive,
			UpdatedAt: d.UpdatedAt,
		}
	}

	return items, nil
}

func (s *dokumenKebijakanService) FindByID(ctx context.Context, id string) (*models.DokumenKebijakanDetailResponse, error) {
	dokumen, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dokumen kebijakan tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data dokumen kebijakan")
	}

	return &models.DokumenKebijakanDetailResponse{
		ID:       dokumen.ID.String(),
		Judul:    dokumen.Judul,
		JudulEN:  dokumen.JudulEN,
		Slug:     dokumen.Slug,
		Konten:   dokumen.Konten,
		KontenEN: dokumen.KontenEN,
		Urutan:   dokumen.Urutan,
		// IsActive:  dokumen.IsActive,
		CreatedAt: dokumen.CreatedAt,
		UpdatedAt: dokumen.UpdatedAt,
	}, nil
}

func (s *dokumenKebijakanService) FindBySlug(ctx context.Context, slug string) (*models.DokumenKebijakanDetailResponse, error) {
	dokumen, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dokumen kebijakan tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data dokumen kebijakan")
	}

	return &models.DokumenKebijakanDetailResponse{
		ID:       dokumen.ID.String(),
		Judul:    dokumen.Judul,
		JudulEN:  dokumen.JudulEN,
		Slug:     dokumen.Slug,
		Konten:   dokumen.Konten,
		KontenEN: dokumen.KontenEN,
		Urutan:   dokumen.Urutan,
		// IsActive:  dokumen.IsActive,
		CreatedAt: dokumen.CreatedAt,
		UpdatedAt: dokumen.UpdatedAt,
	}, nil
}

func (s *dokumenKebijakanService) Update(ctx context.Context, id string, req *models.UpdateDokumenKebijakanRequest) (*models.DokumenKebijakanDetailResponse, error) {
	dokumen, err := s.repo.FindByID(ctx, id)
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

	// Update judul_en if provided
	if req.JudulEN != nil {
		dokumen.JudulEN = *req.JudulEN
	}

	// Update konten if provided (sanitize HTML)
	if req.Konten != nil {
		sanitized := s.sanitizer.Sanitize(*req.Konten)
		dokumen.Konten = sanitized
	}

	// Update konten_en if provided (sanitize HTML)
	if req.KontenEN != nil {
		sanitized := s.sanitizer.Sanitize(*req.KontenEN)
		dokumen.KontenEN = sanitized
	}

	// Update is_active if provided
	if req.IsActive != nil {
		dokumen.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, dokumen); err != nil {
		return nil, errors.New("gagal mengupdate dokumen kebijakan")
	}

	return &models.DokumenKebijakanDetailResponse{
		ID:       dokumen.ID.String(),
		Judul:    dokumen.Judul,
		JudulEN:  dokumen.JudulEN,
		Slug:     dokumen.Slug,
		Konten:   dokumen.Konten,
		KontenEN: dokumen.KontenEN,
		Urutan:   dokumen.Urutan,
		// IsActive:  dokumen.IsActive,
		CreatedAt: dokumen.CreatedAt,
		UpdatedAt: dokumen.UpdatedAt,
	}, nil
}

func (s *dokumenKebijakanService) UpdateBySlug(ctx context.Context, slug string, req *models.UpdateDokumenKebijakanRequest) (*models.DokumenKebijakanDetailResponse, error) {
	dokumen, err := s.repo.FindBySlug(ctx, slug)
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

	// Update judul_en if provided
	if req.JudulEN != nil {
		dokumen.JudulEN = *req.JudulEN
	}

	// Update konten if provided (sanitize HTML)
	if req.Konten != nil {
		sanitized := s.sanitizer.Sanitize(*req.Konten)
		dokumen.Konten = sanitized
	}

	// Update konten_en if provided (sanitize HTML)
	if req.KontenEN != nil {
		sanitized := s.sanitizer.Sanitize(*req.KontenEN)
		dokumen.KontenEN = sanitized
	}

	// Update is_active if provided
	if req.IsActive != nil {
		dokumen.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, dokumen); err != nil {
		return nil, errors.New("gagal mengupdate dokumen kebijakan")
	}

	return &models.DokumenKebijakanDetailResponse{
		ID:       dokumen.ID.String(),
		Judul:    dokumen.Judul,
		JudulEN:  dokumen.JudulEN,
		Slug:     dokumen.Slug,
		Konten:   dokumen.Konten,
		KontenEN: dokumen.KontenEN,
		Urutan:   dokumen.Urutan,
		// IsActive:  dokumen.IsActive,
		CreatedAt: dokumen.CreatedAt,
		UpdatedAt: dokumen.UpdatedAt,
	}, nil
}

func (s *dokumenKebijakanService) GetByIDPublic(ctx context.Context, id string, lang string) (*models.DokumenKebijakanPublicResponse, error) {
	dokumen, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dokumen kebijakan tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data dokumen kebijakan")
	}

	if !dokumen.IsActive {
		return nil, errors.New("dokumen kebijakan tidak aktif")
	}

	// Return based on language
	judul := dokumen.Judul
	konten := dokumen.Konten
	if lang == "en" {
		judul = dokumen.JudulEN
		konten = dokumen.KontenEN
	}

	return &models.DokumenKebijakanPublicResponse{
		ID:     dokumen.ID.String(),
		Judul:  judul,
		Slug:   dokumen.Slug,
		Konten: konten,
		Urutan: dokumen.Urutan,
	}, nil
}

func (s *dokumenKebijakanService) GetBySlugPublic(ctx context.Context, slug string, lang string) (*models.DokumenKebijakanPublicResponse, error) {
	dokumen, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dokumen kebijakan tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data dokumen kebijakan")
	}

	if !dokumen.IsActive {
		return nil, errors.New("dokumen kebijakan tidak aktif")
	}

	// Return based on language
	judul := dokumen.Judul
	konten := dokumen.Konten
	if lang == "en" {
		judul = dokumen.JudulEN
		konten = dokumen.KontenEN
	}

	return &models.DokumenKebijakanPublicResponse{
		ID:     dokumen.ID.String(),
		Judul:  judul,
		Slug:   dokumen.Slug,
		Konten: konten,
		Urutan: dokumen.Urutan,
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
			ID:     d.ID.String(),
			Judul:  d.Judul,
			Urutan: d.Urutan,
		}
	}

	return items, nil
}
