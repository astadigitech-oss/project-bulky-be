package repositories

import (
	"context"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KategoriVideoRepository interface {
	Create(ctx context.Context, kategori *models.KategoriVideo) error
	Update(ctx context.Context, kategori *models.KategoriVideo) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.KategoriVideo, error)
	FindBySlug(ctx context.Context, slug string) (*models.KategoriVideo, error)
	FindAll(ctx context.Context, isActive *bool) ([]models.KategoriVideo, error)
	FindAllPaginated(ctx context.Context, params *dto.KategoriVideoFilterRequest) ([]models.KategoriVideo, int64, error)
	GetMaxUrutan(ctx context.Context) (int, error)
	UpdateUrutan(ctx context.Context, id uuid.UUID, urutan int) error
	FindAllPublicWithCount(ctx context.Context) ([]models.KategoriVideo, error)
}

type kategoriVideoRepository struct {
	db *gorm.DB
}

func NewKategoriVideoRepository(db *gorm.DB) KategoriVideoRepository {
	return &kategoriVideoRepository{db: db}
}

func (r *kategoriVideoRepository) Create(ctx context.Context, kategori *models.KategoriVideo) error {
	return r.db.WithContext(ctx).Create(kategori).Error
}

func (r *kategoriVideoRepository) Update(ctx context.Context, kategori *models.KategoriVideo) error {
	return r.db.WithContext(ctx).Save(kategori).Error
}

func (r *kategoriVideoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.KategoriVideo{}, id).Error
}

func (r *kategoriVideoRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.KategoriVideo, error) {
	var kategori models.KategoriVideo
	err := r.db.WithContext(ctx).First(&kategori, id).Error
	if err != nil {
		return nil, err
	}
	return &kategori, nil
}

func (r *kategoriVideoRepository) FindBySlug(ctx context.Context, slug string) (*models.KategoriVideo, error) {
	var kategori models.KategoriVideo
	err := r.db.WithContext(ctx).Where("slug_id = ? OR slug_en = ?", slug, slug).First(&kategori).Error
	if err != nil {
		return nil, err
	}
	return &kategori, nil
}

func (r *kategoriVideoRepository) FindAll(ctx context.Context, isActive *bool) ([]models.KategoriVideo, error) {
	var kategoris []models.KategoriVideo
	query := r.db.WithContext(ctx).Model(&models.KategoriVideo{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("urutan ASC, nama_id ASC").Find(&kategoris).Error
	return kategoris, err
}

func (r *kategoriVideoRepository) FindAllPaginated(ctx context.Context, params *dto.KategoriVideoFilterRequest) ([]models.KategoriVideo, int64, error) {
	var kategoris []models.KategoriVideo
	var total int64

	query := r.db.WithContext(ctx).Model(&models.KategoriVideo{})

	// Filter by is_active
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Search with ILIKE pada nama_id dan nama_en
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("nama_id ILIKE ? OR nama_en ILIKE ?", search, search)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting: urutan atau updated_at
	sortColumn := "urutan"
	if params.SortBy == "updated_at" {
		sortColumn = "updated_at"
	}
	orderDir := "ASC"
	if params.Order == "desc" {
		orderDir = "DESC"
	}
	query = query.Order(sortColumn + " " + orderDir)

	// Apply pagination
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Find(&kategoris).Error; err != nil {
		return nil, 0, err
	}

	return kategoris, total, nil
}

func (r *kategoriVideoRepository) UpdateUrutan(ctx context.Context, id uuid.UUID, urutan int) error {
	return r.db.WithContext(ctx).Model(&models.KategoriVideo{}).
		Where("id = ?", id).
		Update("urutan", urutan).Error
}

func (r *kategoriVideoRepository) GetMaxUrutan(ctx context.Context) (int, error) {
	var maxUrutan int
	result := r.db.WithContext(ctx).Model(&models.KategoriVideo{}).
		Select("COALESCE(MAX(urutan), 0)").
		Scan(&maxUrutan)
	return maxUrutan, result.Error
}

func (r *kategoriVideoRepository) FindAllPublicWithCount(ctx context.Context) ([]models.KategoriVideo, error) {
	var kategoris []models.KategoriVideo
	err := r.db.WithContext(ctx).
		Select("kategori_video.*, COUNT(video.id) as video_count").
		Joins("LEFT JOIN video ON video.kategori_id = kategori_video.id AND video.deleted_at IS NULL AND video.is_active = true").
		Where("kategori_video.is_active = ? AND kategori_video.deleted_at IS NULL", true).
		Group("kategori_video.id").
		Having("COUNT(video.id) > 0").
		Order("kategori_video.urutan ASC").
		Find(&kategoris).Error
	return kategoris, err
}
