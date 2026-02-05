package services

import (
	"context"
	"errors"
	"time"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FAQService interface {
	GetAll(ctx context.Context, params *models.FAQFilterRequest) ([]models.FAQResponse, *models.PaginationMeta, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.FAQResponse, error)
	GetPublic(ctx context.Context, lang string) ([]models.FAQPublicResponse, error)
	Create(ctx context.Context, req *models.FAQCreateRequest) (*models.FAQResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *models.FAQUpdateRequest) (*models.FAQResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Reorder(ctx context.Context, id uuid.UUID, direction string) (*ReorderResult, error)
}

type faqService struct {
	repo           repositories.FAQRepository
	reorderService *ReorderService
}

func NewFAQService(repo repositories.FAQRepository, reorderService *ReorderService) FAQService {
	return &faqService{
		repo:           repo,
		reorderService: reorderService,
	}
}

func (s *faqService) GetAll(ctx context.Context, params *models.FAQFilterRequest) ([]models.FAQResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	faqs, total, err := s.repo.GetAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]models.FAQResponse, len(faqs))
	for i, faq := range faqs {
		responses[i] = models.FAQResponse{
			ID:         faq.ID.String(),
			Question:   faq.Question,
			QuestionEN: faq.QuestionEN,
			Answer:     faq.Answer,
			AnswerEN:   faq.AnswerEN,
			Urutan:     faq.Urutan,
			IsActive:   faq.IsActive,
			CreatedAt:  faq.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  faq.UpdatedAt.Format(time.RFC3339),
		}
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return responses, &meta, nil
}

func (s *faqService) GetByID(ctx context.Context, id uuid.UUID) (*models.FAQResponse, error) {
	faq, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("FAQ tidak ditemukan")
		}
		return nil, err
	}

	return &models.FAQResponse{
		ID:         faq.ID.String(),
		Question:   faq.Question,
		QuestionEN: faq.QuestionEN,
		Answer:     faq.Answer,
		AnswerEN:   faq.AnswerEN,
		Urutan:     faq.Urutan,
		IsActive:   faq.IsActive,
		CreatedAt:  faq.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  faq.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *faqService) GetPublic(ctx context.Context, lang string) ([]models.FAQPublicResponse, error) {
	faqs, err := s.repo.GetActive(ctx, lang)
	if err != nil {
		return nil, err
	}

	responses := make([]models.FAQPublicResponse, len(faqs))
	for i, faq := range faqs {
		question := faq.Question
		answer := faq.Answer
		if lang == "en" {
			question = faq.QuestionEN
			answer = faq.AnswerEN
		}

		responses[i] = models.FAQPublicResponse{
			ID:       faq.ID.String(),
			Question: question,
			Answer:   answer,
		}
	}

	return responses, nil
}

func (s *faqService) Create(ctx context.Context, req *models.FAQCreateRequest) (*models.FAQResponse, error) {
	// Get max urutan
	maxUrutan, err := s.repo.GetMaxUrutan(ctx)
	if err != nil {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	faq := &models.FAQ{
		Question:   req.Question,
		QuestionEN: req.QuestionEN,
		Answer:     req.Answer,
		AnswerEN:   req.AnswerEN,
		Urutan:     maxUrutan + 1,
		IsActive:   isActive,
	}

	if err := s.repo.Create(ctx, faq); err != nil {
		return nil, err
	}

	return &models.FAQResponse{
		ID:         faq.ID.String(),
		Question:   faq.Question,
		QuestionEN: faq.QuestionEN,
		Answer:     faq.Answer,
		AnswerEN:   faq.AnswerEN,
		Urutan:     faq.Urutan,
		IsActive:   faq.IsActive,
		CreatedAt:  faq.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  faq.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *faqService) Update(ctx context.Context, id uuid.UUID, req *models.FAQUpdateRequest) (*models.FAQResponse, error) {
	faq, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("FAQ tidak ditemukan")
		}
		return nil, err
	}

	faq.Question = req.Question
	faq.QuestionEN = req.QuestionEN
	faq.Answer = req.Answer
	faq.AnswerEN = req.AnswerEN

	if req.IsActive != nil {
		faq.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, faq); err != nil {
		return nil, err
	}

	return &models.FAQResponse{
		ID:         faq.ID.String(),
		Question:   faq.Question,
		QuestionEN: faq.QuestionEN,
		Answer:     faq.Answer,
		AnswerEN:   faq.AnswerEN,
		Urutan:     faq.Urutan,
		IsActive:   faq.IsActive,
		CreatedAt:  faq.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  faq.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *faqService) Delete(ctx context.Context, id uuid.UUID) error {
	faq, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("FAQ tidak ditemukan")
		}
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Reorder after delete
	return s.reorderService.ReorderAfterDelete(ctx, "faq", faq.Urutan, "", nil)
}

func (s *faqService) Reorder(ctx context.Context, id uuid.UUID, direction string) (*ReorderResult, error) {
	return s.reorderService.Reorder(ctx, "faq", id, direction, "", nil)
}
