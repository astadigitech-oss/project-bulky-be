package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FAQRepository interface {
	GetAll(ctx context.Context) ([]models.FAQ, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.FAQ, error)
	GetActive(ctx context.Context, lang string) ([]models.FAQ, error)
	Create(ctx context.Context, faq *models.FAQ) error
	Update(ctx context.Context, faq *models.FAQ) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetMaxUrutan(ctx context.Context) (int, error)
}

type faqRepository struct {
	db *gorm.DB
}

func NewFAQRepository(db *gorm.DB) FAQRepository {
	return &faqRepository{db: db}
}

func (r *faqRepository) GetAll(ctx context.Context) ([]models.FAQ, error) {
	var faqs []models.FAQ
	err := r.db.WithContext(ctx).
		Order("urutan ASC").
		Find(&faqs).Error
	return faqs, err
}

func (r *faqRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.FAQ, error) {
	var faq models.FAQ
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&faq).Error
	if err != nil {
		return nil, err
	}
	return &faq, nil
}

func (r *faqRepository) GetActive(ctx context.Context, lang string) ([]models.FAQ, error) {
	var faqs []models.FAQ
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("urutan ASC").
		Find(&faqs).Error
	return faqs, err
}

func (r *faqRepository) Create(ctx context.Context, faq *models.FAQ) error {
	return r.db.WithContext(ctx).Create(faq).Error
}

func (r *faqRepository) Update(ctx context.Context, faq *models.FAQ) error {
	return r.db.WithContext(ctx).Save(faq).Error
}

func (r *faqRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.FAQ{}).Error
}

func (r *faqRepository) GetMaxUrutan(ctx context.Context) (int, error) {
	var maxUrutan int
	err := r.db.WithContext(ctx).
		Model(&models.FAQ{}).
		Select("COALESCE(MAX(urutan), 0)").
		Scan(&maxUrutan).Error
	return maxUrutan, err
}
