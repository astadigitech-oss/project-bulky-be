package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlogRepository interface {
	Create(ctx context.Context, blog *models.Blog) error
	Update(ctx context.Context, blog *models.Blog) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Blog, error)
	FindBySlug(ctx context.Context, slug string) (*models.Blog, error)
	FindAll(ctx context.Context, isActive *bool, kategoriID *uuid.UUID, limit, offset int) ([]models.Blog, int64, error)
	Search(ctx context.Context, keyword string, isActive *bool, limit, offset int) ([]models.Blog, int64, error)
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
	AddLabels(ctx context.Context, blogID uuid.UUID, labelIDs []uuid.UUID) error
	RemoveAllLabels(ctx context.Context, blogID uuid.UUID) error
	FindRelated(ctx context.Context, blogID uuid.UUID, kategoriID uuid.UUID, limit int) ([]models.Blog, error)
	FindPopular(ctx context.Context, limit int) ([]models.Blog, error)
	GetStatistics(ctx context.Context) (map[string]interface{}, error)
	ToggleStatus(ctx context.Context, id uuid.UUID) error
}

type blogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) Create(ctx context.Context, blog *models.Blog) error {
	return r.db.WithContext(ctx).Create(blog).Error
}

func (r *blogRepository) Update(ctx context.Context, blog *models.Blog) error {
	return r.db.WithContext(ctx).Save(blog).Error
}

func (r *blogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Blog{}, id).Error
}

func (r *blogRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Blog, error) {
	var blog models.Blog
	err := r.db.WithContext(ctx).
		Preload("Kategori").
		Preload("Labels").
		First(&blog, id).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *blogRepository) FindBySlug(ctx context.Context, slug string) (*models.Blog, error) {
	var blog models.Blog
	err := r.db.WithContext(ctx).
		Preload("Kategori").
		Preload("Labels").
		Where("slug = ?", slug).
		First(&blog).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *blogRepository) FindAll(ctx context.Context, isActive *bool, kategoriID *uuid.UUID, limit, offset int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Blog{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if kategoriID != nil {
		query = query.Where("kategori_id = ?", *kategoriID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Preload("Kategori").
		Preload("Labels").
		Order("published_at DESC NULLS LAST, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&blogs).Error

	return blogs, total, err
}

func (r *blogRepository) Search(ctx context.Context, keyword string, isActive *bool, limit, offset int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Blog{}).
		Where("to_tsvector('indonesian', judul_id || ' ' || deskripsi_singkat_id) @@ plainto_tsquery('indonesian', ?)", keyword)

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Preload("Kategori").
		Preload("Labels").
		Order("view_count DESC").
		Limit(limit).
		Offset(offset).
		Find(&blogs).Error

	return blogs, total, err
}

func (r *blogRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Blog{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).
		Error
}

func (r *blogRepository) AddLabels(ctx context.Context, blogID uuid.UUID, labelIDs []uuid.UUID) error {
	if len(labelIDs) == 0 {
		return nil
	}

	var blogLabels []models.BlogLabel
	for _, labelID := range labelIDs {
		blogLabels = append(blogLabels, models.BlogLabel{
			BlogID:  blogID,
			LabelID: labelID,
		})
	}

	return r.db.WithContext(ctx).Create(&blogLabels).Error
}

func (r *blogRepository) RemoveAllLabels(ctx context.Context, blogID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("blog_id = ?", blogID).
		Delete(&models.BlogLabel{}).
		Error
}

func (r *blogRepository) FindRelated(ctx context.Context, blogID uuid.UUID, kategoriID uuid.UUID, limit int) ([]models.Blog, error) {
	var blogs []models.Blog
	err := r.db.WithContext(ctx).
		Preload("Kategori").
		Preload("Labels").
		Where("kategori_id = ? AND id != ? AND is_active = ? AND deleted_at IS NULL", kategoriID, blogID, true).
		Order("view_count DESC").
		Limit(limit).
		Find(&blogs).Error
	return blogs, err
}

func (r *blogRepository) FindPopular(ctx context.Context, limit int) ([]models.Blog, error) {
	var blogs []models.Blog
	err := r.db.WithContext(ctx).
		Preload("Kategori").
		Where("is_active = ? AND deleted_at IS NULL", true).
		Order("view_count DESC").
		Limit(limit).
		Find(&blogs).Error
	return blogs, err
}

func (r *blogRepository) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total blog
	var totalBlog int64
	r.db.WithContext(ctx).Model(&models.Blog{}).Count(&totalBlog)
	stats["total_blog"] = totalBlog

	// Total published
	var totalPublished int64
	r.db.WithContext(ctx).Model(&models.Blog{}).Where("is_active = ?", true).Count(&totalPublished)
	stats["total_published"] = totalPublished

	// Total draft
	stats["total_draft"] = totalBlog - totalPublished

	// Total views
	var totalViews int64
	r.db.WithContext(ctx).Model(&models.Blog{}).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)
	stats["total_views"] = totalViews

	return stats, nil
}

func (r *blogRepository) ToggleStatus(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Blog{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active")).
		Error
}
