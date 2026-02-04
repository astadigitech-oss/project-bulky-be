package services

import (
	"context"
	"encoding/json"
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
	GetFAQ(ctx context.Context, lang string) (*models.FAQLegacyResponse, error)
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

// Fixed pages order (predefined)
var fixedPagesOrder = []string{
	"tentang-kami",
	"cara-membeli",
	"tentang-pembayaran",
	"hubungi-kami",
	"faq",
	"syarat-ketentuan",
	"kebijakan-privasi",
}

func (s *dokumenKebijakanService) FindAll(ctx context.Context) ([]models.DokumenKebijakanListResponse, error) {
	dokumens, err := s.repo.FindAllSimple(ctx)
	if err != nil {
		return nil, errors.New("gagal mengambil data dokumen kebijakan")
	}

	// Create map for quick lookup
	docMap := make(map[string]models.DokumenKebijakan)
	for _, d := range dokumens {
		docMap[d.Slug] = d
	}

	// Sort by predefined order
	items := make([]models.DokumenKebijakanListResponse, 0, len(fixedPagesOrder))
	for _, slug := range fixedPagesOrder {
		if d, ok := docMap[slug]; ok {
			items = append(items, models.DokumenKebijakanListResponse{
				ID:        d.ID.String(),
				Judul:     d.Judul,
				JudulEN:   d.JudulEN,
				Slug:      d.Slug,
				UpdatedAt: d.UpdatedAt,
			})
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
		ID:        dokumen.ID.String(),
		Judul:     dokumen.Judul,
		JudulEN:   dokumen.JudulEN,
		Slug:      dokumen.Slug,
		Konten:    dokumen.Konten,
		KontenEN:  dokumen.KontenEN,
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
		ID:        dokumen.ID.String(),
		Judul:     dokumen.Judul,
		JudulEN:   dokumen.JudulEN,
		Slug:      dokumen.Slug,
		Konten:    dokumen.Konten,
		KontenEN:  dokumen.KontenEN,
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

	// Update konten if provided (sanitize HTML for non-FAQ, keep JSON for FAQ)
	if req.Konten != nil {
		if dokumen.Slug == "faq" {
			// For FAQ, validate JSON format
			var faqItems []models.FAQContentItem
			if err := json.Unmarshal([]byte(*req.Konten), &faqItems); err != nil {
				return nil, errors.New("format FAQ tidak valid, harus berupa JSON array")
			}
			dokumen.Konten = *req.Konten
		} else {
			// For other pages, sanitize HTML
			sanitized := s.sanitizer.Sanitize(*req.Konten)
			dokumen.Konten = sanitized
		}
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
		ID:        dokumen.ID.String(),
		Judul:     dokumen.Judul,
		JudulEN:   dokumen.JudulEN,
		Slug:      dokumen.Slug,
		Konten:    dokumen.Konten,
		KontenEN:  dokumen.KontenEN,
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

	// Update konten if provided (sanitize HTML for non-FAQ, keep JSON for FAQ)
	if req.Konten != nil {
		if dokumen.Slug == "faq" {
			// For FAQ, validate JSON format
			var faqItems []models.FAQContentItem
			if err := json.Unmarshal([]byte(*req.Konten), &faqItems); err != nil {
				return nil, errors.New("format FAQ tidak valid, harus berupa JSON array")
			}
			dokumen.Konten = *req.Konten
		} else {
			// For other pages, sanitize HTML
			sanitized := s.sanitizer.Sanitize(*req.Konten)
			dokumen.Konten = sanitized
		}
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
		ID:        dokumen.ID.String(),
		Judul:     dokumen.Judul,
		JudulEN:   dokumen.JudulEN,
		Slug:      dokumen.Slug,
		Konten:    dokumen.Konten,
		KontenEN:  dokumen.KontenEN,
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
			ID:    d.ID.String(),
			Judul: d.Judul,
			Slug:  d.Slug,
		}
	}

	return items, nil
}

func (s *dokumenKebijakanService) GetFAQ(ctx context.Context, lang string) (*models.FAQLegacyResponse, error) {
	doc, err := s.repo.FindBySlug(ctx, "faq")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("FAQ tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data FAQ")
	}

	if !doc.IsActive {
		return nil, errors.New("FAQ tidak aktif")
	}

	// Parse JSON content
	var items []models.FAQContentItem
	if err := json.Unmarshal([]byte(doc.Konten), &items); err != nil {
		return nil, errors.New("format FAQ tidak valid")
	}

	// Build response based on language
	response := &models.FAQLegacyResponse{
		Judul: doc.Judul,
		Items: make([]models.FAQItem, 0, len(items)),
	}

	if lang == "en" {
		response.Judul = doc.JudulEN
	}

	for _, item := range items {
		faqItem := models.FAQItem{}
		if lang == "en" {
			faqItem.Question = item.QuestionEN
			faqItem.Answer = item.AnswerEN
		} else {
			faqItem.Question = item.Question
			faqItem.Answer = item.Answer
		}
		response.Items = append(response.Items, faqItem)
	}

	return response, nil
}
