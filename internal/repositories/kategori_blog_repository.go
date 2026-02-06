package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KategoriBlogRepository interface {
	Create(ctx context.Context, kategori *models.KategoriBlog) error
	Update(ctx context.Context, kategori *models.KategoriBlog) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.KategoriBlog, error)
	FindBySlug(ctx context.Context, slug string) (*models.KategoriBlog, error)
	FindAll(ctx context.Context, isActive *bool) ([]models.KategoriBlog, error)
	CountBlogByKategori(ctx context.Context, kategoriID uuid.UUID) (int64, error)
	UpdateUrutan(ctx context.Context, id uuid.UUID, urutan int) error
	FindAllPublicWithCount(ctx context.Context) ([]models.KategoriBlog, error)
}

type kategoriBlogRepository struct {
	db *gorm.DB
}

func NewKategoriBlogRepository(db *gorm.DB) KategoriBlogRepository {
	return &kategoriBlogRepository{db: db}
}

func (r *kategoriBlogRepository) Create(ctx context.Context, kategori *models.KategoriBlog) error {
	return r.db.WithContext(ctx).Create(kategori).Error
}

func (r *kategoriBlogRepository) Update(ctx context.Context, kategori *models.KategoriBlog) error {
	return r.db.WithContext(ctx).Save(kategori).Error
}

func (r *kategoriBlogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.KategoriBlog{}, id).Error
}

func (r *kategoriBlogRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.KategoriBlog, error) {
	var kategori models.KategoriBlog
	err := r.db.WithContext(ctx).First(&kategori, id).Error
	if err != nil {
		return nil, err
	}
	return &kategori, nil
}

func (r *kategoriBlogRepository) FindBySlug(ctx context.Context, slug string) (*models.KategoriBlog, error) {
	var kategori models.KategoriBlog
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&kategori).Error
	if err != nil {
		return nil, err
	}
	return &kategori, nil
}

func (r *kategoriBlogRepository) FindAll(ctx context.Context, isActive *bool) ([]models.KategoriBlog, error) {
	var kategoris []models.KategoriBlog
	query := r.db.WithContext(ctx).Model(&models.KategoriBlog{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("urutan ASC, nama->>'id' ASC").Find(&kategoris).Error
	return kategoris, err
}

func (r *kategoriBlogRepository) CountBlogByKategori(ctx context.Context, kategoriID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Blog{}).
		Where("kategori_id = ? AND deleted_at IS NULL", kategoriID).
		Count(&count).Error
	return count, err
}

func (r *kategoriBlogRepository) UpdateUrutan(ctx context.Context, id uuid.UUID, urutan int) error {
	return r.db.WithContext(ctx).Model(&models.KategoriBlog{}).
		Where("id = ?", id).
		Update("urutan", urutan).Error
}

func (r *kategoriBlogRepository) FindAllPublicWithCount(ctx context.Context) ([]models.KategoriBlog, error) {
	var kategoris []models.KategoriBlog
	err := r.db.WithContext(ctx).
		Select("kategori_blog.*, COUNT(blog.id) as blog_count").
		Joins("LEFT JOIN blog ON blog.kategori_id = kategori_blog.id AND blog.deleted_at IS NULL AND blog.is_active = true").
		Where("kategori_blog.is_active = ? AND kategori_blog.deleted_at IS NULL", true).
		Group("kategori_blog.id").
		Having("COUNT(blog.id) > 0").
		Order("kategori_blog.urutan ASC").
		Find(&kategoris).Error
	return kategoris, err
}
