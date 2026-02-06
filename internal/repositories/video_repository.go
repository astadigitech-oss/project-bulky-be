package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoRepository interface {
	Create(ctx context.Context, video *models.Video) error
	Update(ctx context.Context, video *models.Video) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Video, error)
	FindBySlug(ctx context.Context, slug string) (*models.Video, error)
	FindAll(ctx context.Context, isActive *bool, kategoriID *uuid.UUID, limit, offset int) ([]models.Video, int64, error)
	Search(ctx context.Context, keyword string, isActive *bool, limit, offset int) ([]models.Video, int64, error)
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
	FindPopular(ctx context.Context, limit int) ([]models.Video, error)
	FindRelated(ctx context.Context, videoID uuid.UUID, kategoriID uuid.UUID, limit int) ([]models.Video, error)
	GetStatistics(ctx context.Context) (map[string]interface{}, error)
	ToggleStatus(ctx context.Context, id uuid.UUID) error
}

type videoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) Create(ctx context.Context, video *models.Video) error {
	return r.db.WithContext(ctx).Create(video).Error
}

func (r *videoRepository) Update(ctx context.Context, video *models.Video) error {
	return r.db.WithContext(ctx).Save(video).Error
}

func (r *videoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Video{}, id).Error
}

func (r *videoRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Video, error) {
	var video models.Video
	err := r.db.WithContext(ctx).
		Preload("Kategori").
		First(&video, id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) FindBySlug(ctx context.Context, slug string) (*models.Video, error) {
	var video models.Video
	err := r.db.WithContext(ctx).
		Preload("Kategori").
		Where("slug = ?", slug).
		First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) FindAll(ctx context.Context, isActive *bool, kategoriID *uuid.UUID, limit, offset int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Video{})

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
		Order("published_at DESC NULLS LAST, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error

	return videos, total, err
}

func (r *videoRepository) Search(ctx context.Context, keyword string, isActive *bool, limit, offset int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Video{}).
		Where("to_tsvector('indonesian', judul_id || ' ' || deskripsi_id) @@ plainto_tsquery('indonesian', ?)", keyword)

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Preload("Kategori").
		Order("view_count DESC").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error

	return videos, total, err
}

func (r *videoRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).
		Error
}

func (r *videoRepository) FindPopular(ctx context.Context, limit int) ([]models.Video, error) {
	var videos []models.Video
	err := r.db.WithContext(ctx).
		Preload("Kategori").
		Where("is_active = ?", true).
		Order("view_count DESC").
		Limit(limit).
		Find(&videos).Error
	return videos, err
}

func (r *videoRepository) FindRelated(ctx context.Context, videoID uuid.UUID, kategoriID uuid.UUID, limit int) ([]models.Video, error) {
	var videos []models.Video
	err := r.db.WithContext(ctx).
		Preload("Kategori").
		Where("kategori_id = ? AND id != ? AND is_active = ? AND deleted_at IS NULL", kategoriID, videoID, true).
		Order("view_count DESC").
		Limit(limit).
		Find(&videos).Error
	return videos, err
}

func (r *videoRepository) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total video
	var totalVideo int64
	r.db.WithContext(ctx).Model(&models.Video{}).Count(&totalVideo)
	stats["total_video"] = totalVideo

	// Total published
	var totalPublished int64
	r.db.WithContext(ctx).Model(&models.Video{}).Where("is_active = ?", true).Count(&totalPublished)
	stats["total_published"] = totalPublished

	// Total draft
	stats["total_draft"] = totalVideo - totalPublished

	// Total views
	var totalViews int64
	r.db.WithContext(ctx).Model(&models.Video{}).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)
	stats["total_views"] = totalViews

	// Total durasi
	var totalDurasi int64
	r.db.WithContext(ctx).Model(&models.Video{}).Select("COALESCE(SUM(durasi_detik), 0)").Scan(&totalDurasi)
	stats["total_durasi_menit"] = totalDurasi / 60

	return stats, nil
}

func (r *videoRepository) ToggleStatus(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Video{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active")).
		Error
}
