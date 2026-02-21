package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type DokumenKebijakanRepository interface {
	Create(ctx context.Context, dokumen *models.DokumenKebijakan) error
	FindByID(ctx context.Context, id string) (*models.DokumenKebijakan, error)
	FindBySlug(ctx context.Context, slug string) (*models.DokumenKebijakan, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.DokumenKebijakan, int64, error)
	FindAllSimple(ctx context.Context) ([]models.DokumenKebijakan, error)
	Update(ctx context.Context, dokumen *models.DokumenKebijakan) error
	Delete(ctx context.Context, id string) error
	GetActiveList(ctx context.Context) ([]models.DokumenKebijakan, error)
}

type dokumenKebijakanRepository struct {
	db *gorm.DB
}

func NewDokumenKebijakanRepository(db *gorm.DB) DokumenKebijakanRepository {
	return &dokumenKebijakanRepository{db: db}
}

func (r *dokumenKebijakanRepository) Create(ctx context.Context, dokumen *models.DokumenKebijakan) error {
	return r.db.WithContext(ctx).Create(dokumen).Error
}

func (r *dokumenKebijakanRepository) FindByID(ctx context.Context, id string) (*models.DokumenKebijakan, error) {
	var dokumen models.DokumenKebijakan
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&dokumen).Error
	if err != nil {
		return nil, err
	}
	return &dokumen, nil
}

func (r *dokumenKebijakanRepository) FindBySlug(ctx context.Context, slug string) (*models.DokumenKebijakan, error) {
	var dokumen models.DokumenKebijakan
	err := r.db.WithContext(ctx).Where("slug_id = ? OR slug_en = ?", slug, slug).First(&dokumen).Error
	if err != nil {
		return nil, err
	}
	return &dokumen, nil
}

func (r *dokumenKebijakanRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.DokumenKebijakan, int64, error) {
	var dokumens []models.DokumenKebijakan
	var total int64

	query := r.db.WithContext(ctx).Model(&models.DokumenKebijakan{})

	if params.Search != "" {
		query = query.Where("judul ILIKE ? OR judul_en ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	validSortFields := map[string]bool{
		"judul":      true,
		"is_active":  true,
		"created_at": true,
		"updated_at": true,
		"slug":       true,
	}

	sortBy := params.SortBy
	if !validSortFields[sortBy] {
		sortBy = "slug"
	}

	order := params.Order
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	orderClause := sortBy + " " + order
	query = query.Order(orderClause)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Find(&dokumens).Error; err != nil {
		return nil, 0, err
	}

	return dokumens, total, nil
}

func (r *dokumenKebijakanRepository) Update(ctx context.Context, dokumen *models.DokumenKebijakan) error {
	return r.db.WithContext(ctx).Save(dokumen).Error
}

func (r *dokumenKebijakanRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.DokumenKebijakan{}).Error
}

func (r *dokumenKebijakanRepository) GetActiveList(ctx context.Context) ([]models.DokumenKebijakan, error) {
	var dokumens []models.DokumenKebijakan
	err := r.db.WithContext(ctx).
		Select("id", "judul", "judul_en", "slug").
		Where("is_active = ?", true).
		Order("slug ASC").
		Find(&dokumens).Error
	return dokumens, err
}

func (r *dokumenKebijakanRepository) FindAllSimple(ctx context.Context) ([]models.DokumenKebijakan, error) {
	var dokumens []models.DokumenKebijakan
	err := r.db.WithContext(ctx).
		Order("slug ASC").
		Find(&dokumens).Error
	return dokumens, err
}
