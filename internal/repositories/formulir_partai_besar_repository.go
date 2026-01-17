package repositories

import (
	"context"
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FormulirPartaiBesarRepository interface {
	// Config
	GetConfig(ctx context.Context) (*models.FormulirPartaiBesarConfig, error)
	UpdateConfig(ctx context.Context, config *models.FormulirPartaiBesarConfig) error

	// Anggaran
	CreateAnggaran(ctx context.Context, anggaran *models.FormulirPartaiBesarAnggaran) error
	FindAllAnggaran(ctx context.Context) ([]models.FormulirPartaiBesarAnggaran, error)
	FindAnggaranByID(ctx context.Context, id uuid.UUID) (*models.FormulirPartaiBesarAnggaran, error)
	UpdateAnggaran(ctx context.Context, anggaran *models.FormulirPartaiBesarAnggaran) error
	DeleteAnggaran(ctx context.Context, id uuid.UUID) error
	ReorderAnggaran(ctx context.Context, items []models.ReorderItem) error

	// Submission
	CreateSubmission(ctx context.Context, submission *models.FormulirPartaiBesarSubmission) error
	FindAllSubmission(ctx context.Context, params *models.FormulirSubmissionFilterRequest) ([]models.FormulirPartaiBesarSubmission, int64, error)
	FindSubmissionByID(ctx context.Context, id uuid.UUID) (*models.FormulirPartaiBesarSubmission, error)
	UpdateSubmission(ctx context.Context, submission *models.FormulirPartaiBesarSubmission) error
}

type formulirPartaiBesarRepository struct {
	db *gorm.DB
}

func NewFormulirPartaiBesarRepository(db *gorm.DB) FormulirPartaiBesarRepository {
	return &formulirPartaiBesarRepository{db: db}
}

// ========================================
// Config
// ========================================

func (r *formulirPartaiBesarRepository) GetConfig(ctx context.Context) (*models.FormulirPartaiBesarConfig, error) {
	var config models.FormulirPartaiBesarConfig
	if err := r.db.WithContext(ctx).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *formulirPartaiBesarRepository) UpdateConfig(ctx context.Context, config *models.FormulirPartaiBesarConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

// ========================================
// Anggaran
// ========================================

func (r *formulirPartaiBesarRepository) CreateAnggaran(ctx context.Context, anggaran *models.FormulirPartaiBesarAnggaran) error {
	return r.db.WithContext(ctx).Create(anggaran).Error
}

func (r *formulirPartaiBesarRepository) FindAllAnggaran(ctx context.Context) ([]models.FormulirPartaiBesarAnggaran, error) {
	var items []models.FormulirPartaiBesarAnggaran

	query := r.db.WithContext(ctx).Model(&models.FormulirPartaiBesarAnggaran{})

	// Sorting by urutan ASC
	query = query.Order("urutan ASC, created_at ASC")

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *formulirPartaiBesarRepository) FindAnggaranByID(ctx context.Context, id uuid.UUID) (*models.FormulirPartaiBesarAnggaran, error) {
	var anggaran models.FormulirPartaiBesarAnggaran
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&anggaran).Error; err != nil {
		return nil, err
	}
	return &anggaran, nil
}

func (r *formulirPartaiBesarRepository) UpdateAnggaran(ctx context.Context, anggaran *models.FormulirPartaiBesarAnggaran) error {
	return r.db.WithContext(ctx).Save(anggaran).Error
}

func (r *formulirPartaiBesarRepository) DeleteAnggaran(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.FormulirPartaiBesarAnggaran{}).Error
}

func (r *formulirPartaiBesarRepository) ReorderAnggaran(ctx context.Context, items []models.ReorderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			id, _ := uuid.Parse(item.ID)
			if err := tx.Model(&models.FormulirPartaiBesarAnggaran{}).
				Where("id = ?", id).
				Update("urutan", item.Urutan).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ========================================
// Submission
// ========================================

func (r *formulirPartaiBesarRepository) CreateSubmission(ctx context.Context, submission *models.FormulirPartaiBesarSubmission) error {
	return r.db.WithContext(ctx).Create(submission).Error
}

func (r *formulirPartaiBesarRepository) FindAllSubmission(ctx context.Context, params *models.FormulirSubmissionFilterRequest) ([]models.FormulirPartaiBesarSubmission, int64, error) {
	var items []models.FormulirPartaiBesarSubmission
	var total int64

	query := r.db.WithContext(ctx).Model(&models.FormulirPartaiBesarSubmission{}).
		Preload("Buyer").
		Preload("Anggaran")

	// Search
	if params.Search != "" {
		query = query.Where("nama ILIKE ? OR telepon ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Filter by email sent
	if params.EmailSent != nil {
		query = query.Where("email_sent = ?", *params.EmailSent)
	}

	// Filter by date range
	if params.TanggalDari != nil && *params.TanggalDari != "" {
		query = query.Where("created_at >= ?", *params.TanggalDari)
	}
	if params.TanggalSampai != nil && *params.TanggalSampai != "" {
		query = query.Where("created_at <= ?", *params.TanggalSampai)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Validate sort fields
	validSortFields := map[string]bool{
		"nama":       true,
		"telepon":    true,
		"email_sent": true,
		"created_at": true,
	}

	sortBy := params.SortBy
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	order := params.Order
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// Sorting
	orderClause := sortBy + " " + order
	query = query.Order(orderClause)

	// Pagination
	if params.PerPage > 0 {
		query = query.Offset(params.GetOffset()).Limit(params.PerPage)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *formulirPartaiBesarRepository) FindSubmissionByID(ctx context.Context, id uuid.UUID) (*models.FormulirPartaiBesarSubmission, error) {
	var submission models.FormulirPartaiBesarSubmission
	if err := r.db.WithContext(ctx).
		Preload("Buyer").
		Preload("Anggaran").
		Where("id = ?", id).
		First(&submission).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *formulirPartaiBesarRepository) UpdateSubmission(ctx context.Context, submission *models.FormulirPartaiBesarSubmission) error {
	return r.db.WithContext(ctx).Save(submission).Error
}
