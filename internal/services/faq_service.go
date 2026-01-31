package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type FAQService interface {
	Get(ctx context.Context) (*models.FAQAdminResponse, error)
	Update(ctx context.Context, req *models.FAQUpdateRequest) (*models.FAQAdminResponse, error)
	AddItem(ctx context.Context, req *models.FAQItemRequest) (*models.FAQItemResponse, error)
	UpdateItem(ctx context.Context, index int, req *models.FAQItemRequest) (*models.FAQItemResponse, error)
	DeleteItem(ctx context.Context, index int) error
	ReorderItem(ctx context.Context, req *models.FAQReorderRequest) (*models.FAQAdminResponse, error)
}

type faqService struct {
	db *gorm.DB
}

func NewFAQService(db *gorm.DB) FAQService {
	return &faqService{db: db}
}

// Internal struct for JSON parsing
type faqItem struct {
	Question   string `json:"question"`
	QuestionEN string `json:"question_en"`
	Answer     string `json:"answer"`
	AnswerEN   string `json:"answer_en"`
}

func (s *faqService) Get(ctx context.Context) (*models.FAQAdminResponse, error) {
	var doc models.DokumenKebijakan

	err := s.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", "faq").First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("FAQ tidak ditemukan")
		}
		return nil, err
	}

	// Parse JSON content
	var items []faqItem
	if doc.Konten != "" {
		if err := json.Unmarshal([]byte(doc.Konten), &items); err != nil {
			items = []faqItem{}
		}
	}

	// Build response
	response := &models.FAQAdminResponse{
		ID:        doc.ID.String(),
		Judul:     doc.Judul,
		JudulEN:   doc.JudulEN,
		Slug:      doc.Slug,
		IsActive:  doc.IsActive,
		Items:     make([]models.FAQItemResponse, len(items)),
		CreatedAt: doc.CreatedAt.Format(time.RFC3339),
		UpdatedAt: doc.UpdatedAt.Format(time.RFC3339),
	}

	for i, item := range items {
		response.Items[i] = models.FAQItemResponse{
			Index:      i,
			Question:   item.Question,
			QuestionEN: item.QuestionEN,
			Answer:     item.Answer,
			AnswerEN:   item.AnswerEN,
		}
	}

	return response, nil
}

func (s *faqService) Update(ctx context.Context, req *models.FAQUpdateRequest) (*models.FAQAdminResponse, error) {
	var doc models.DokumenKebijakan

	err := s.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", "faq").First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("FAQ tidak ditemukan")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Judul != nil {
		doc.Judul = *req.Judul
	}
	if req.JudulEN != nil {
		doc.JudulEN = *req.JudulEN
	}
	if req.IsActive != nil {
		doc.IsActive = *req.IsActive
	}

	// Update items if provided
	if req.Items != nil {
		items := make([]faqItem, len(req.Items))
		for i, item := range req.Items {
			items[i] = faqItem{
				Question:   item.Question,
				QuestionEN: item.QuestionEN,
				Answer:     item.Answer,
				AnswerEN:   item.AnswerEN,
			}
		}

		jsonBytes, err := json.Marshal(items)
		if err != nil {
			return nil, err
		}
		doc.Konten = string(jsonBytes)
		// For FAQ, konten_en is same as konten (bilingual in one JSON)
		doc.KontenEN = string(jsonBytes)
	}

	if err := s.db.WithContext(ctx).Save(&doc).Error; err != nil {
		return nil, err
	}

	return s.Get(ctx)
}

func (s *faqService) AddItem(ctx context.Context, req *models.FAQItemRequest) (*models.FAQItemResponse, error) {
	var doc models.DokumenKebijakan

	err := s.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", "faq").First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("FAQ tidak ditemukan")
		}
		return nil, err
	}

	// Parse existing items
	var items []faqItem
	if doc.Konten != "" {
		json.Unmarshal([]byte(doc.Konten), &items)
	}

	// Add new item
	newItem := faqItem{
		Question:   req.Question,
		QuestionEN: req.QuestionEN,
		Answer:     req.Answer,
		AnswerEN:   req.AnswerEN,
	}
	items = append(items, newItem)

	// Save back
	jsonBytes, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	doc.Konten = string(jsonBytes)
	doc.KontenEN = string(jsonBytes)

	if err := s.db.WithContext(ctx).Save(&doc).Error; err != nil {
		return nil, err
	}

	return &models.FAQItemResponse{
		Index:      len(items) - 1,
		Question:   newItem.Question,
		QuestionEN: newItem.QuestionEN,
		Answer:     newItem.Answer,
		AnswerEN:   newItem.AnswerEN,
	}, nil
}

func (s *faqService) UpdateItem(ctx context.Context, index int, req *models.FAQItemRequest) (*models.FAQItemResponse, error) {
	var doc models.DokumenKebijakan

	err := s.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", "faq").First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("FAQ tidak ditemukan")
		}
		return nil, err
	}

	// Parse existing items
	var items []faqItem
	if doc.Konten != "" {
		json.Unmarshal([]byte(doc.Konten), &items)
	}

	// Check bounds
	if index < 0 || index >= len(items) {
		return nil, fmt.Errorf("FAQ item dengan index %d tidak ditemukan", index)
	}

	// Update item (partial update supported)
	if req.Question != "" {
		items[index].Question = req.Question
	}
	if req.QuestionEN != "" {
		items[index].QuestionEN = req.QuestionEN
	}
	if req.Answer != "" {
		items[index].Answer = req.Answer
	}
	if req.AnswerEN != "" {
		items[index].AnswerEN = req.AnswerEN
	}

	// Save back
	jsonBytes, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	doc.Konten = string(jsonBytes)
	doc.KontenEN = string(jsonBytes)

	if err := s.db.WithContext(ctx).Save(&doc).Error; err != nil {
		return nil, err
	}

	return &models.FAQItemResponse{
		Index:      index,
		Question:   items[index].Question,
		QuestionEN: items[index].QuestionEN,
		Answer:     items[index].Answer,
		AnswerEN:   items[index].AnswerEN,
	}, nil
}

func (s *faqService) DeleteItem(ctx context.Context, index int) error {
	var doc models.DokumenKebijakan

	err := s.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", "faq").First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("FAQ tidak ditemukan")
		}
		return err
	}

	// Parse existing items
	var items []faqItem
	if doc.Konten != "" {
		json.Unmarshal([]byte(doc.Konten), &items)
	}

	// Check bounds
	if index < 0 || index >= len(items) {
		return fmt.Errorf("FAQ item dengan index %d tidak ditemukan", index)
	}

	// Remove item
	items = append(items[:index], items[index+1:]...)

	// Save back
	jsonBytes, err := json.Marshal(items)
	if err != nil {
		return err
	}
	doc.Konten = string(jsonBytes)
	doc.KontenEN = string(jsonBytes)

	return s.db.WithContext(ctx).Save(&doc).Error
}

func (s *faqService) ReorderItem(ctx context.Context, req *models.FAQReorderRequest) (*models.FAQAdminResponse, error) {
	var doc models.DokumenKebijakan

	err := s.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", "faq").First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("FAQ tidak ditemukan")
		}
		return nil, err
	}

	// Parse existing items
	var items []faqItem
	if doc.Konten != "" {
		json.Unmarshal([]byte(doc.Konten), &items)
	}

	// Check bounds
	if req.Index < 0 || req.Index >= len(items) {
		return nil, fmt.Errorf("FAQ item dengan index %d tidak ditemukan", req.Index)
	}

	// Check edge cases
	if req.Direction == "up" && req.Index == 0 {
		return nil, errors.New("item sudah berada di posisi paling atas")
	}
	if req.Direction == "down" && req.Index == len(items)-1 {
		return nil, errors.New("item sudah berada di posisi paling bawah")
	}

	// Swap items
	targetIndex := req.Index - 1
	if req.Direction == "down" {
		targetIndex = req.Index + 1
	}
	items[req.Index], items[targetIndex] = items[targetIndex], items[req.Index]

	// Save back
	jsonBytes, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	doc.Konten = string(jsonBytes)
	doc.KontenEN = string(jsonBytes)

	if err := s.db.WithContext(ctx).Save(&doc).Error; err != nil {
		return nil, err
	}

	return s.Get(ctx)
}
