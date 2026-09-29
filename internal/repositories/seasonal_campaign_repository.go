package repositories

import (
	"context"
	"time"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type SeasonalCampaignRepository interface {
	Create(context.Context, *models.SeasonalCampaign) error
	FindByID(context.Context, string) (*models.SeasonalCampaign, error)
	FindAll(context.Context, *models.SeasonalCampaignFilterRequest) ([]models.SeasonalCampaign, int64, error)
	Save(context.Context, *models.SeasonalCampaign) error
	Delete(context.Context, *models.SeasonalCampaign) error
	HasPublishedOverlap(context.Context, time.Time, time.Time, string) (bool, error)
}

type seasonalCampaignRepository struct{ db *gorm.DB }

func NewSeasonalCampaignRepository(db *gorm.DB) SeasonalCampaignRepository {
	return &seasonalCampaignRepository{db: db}
}

func (r *seasonalCampaignRepository) Create(ctx context.Context, campaign *models.SeasonalCampaign) error {
	return r.db.WithContext(ctx).Create(campaign).Error
}

func (r *seasonalCampaignRepository) FindByID(ctx context.Context, id string) (*models.SeasonalCampaign, error) {
	var campaign models.SeasonalCampaign
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&campaign).Error
	return &campaign, err
}

func (r *seasonalCampaignRepository) FindAll(ctx context.Context, params *models.SeasonalCampaignFilterRequest) ([]models.SeasonalCampaign, int64, error) {
	var campaigns []models.SeasonalCampaign
	var total int64
	query := r.db.WithContext(ctx).Model(&models.SeasonalCampaign{})
	if params.Search != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Search+"%")
	}
	if params.TanggalMulai != nil {
		query = query.Where("tanggal_mulai >= ?", *params.TanggalMulai)
	}
	if params.TanggalSelesai != nil {
		query = query.Where("tanggal_selesai <= ?", *params.TanggalSelesai)
	}
	switch params.Status {
	case "draft":
		query = query.Where("is_published = false AND cancelled_at IS NULL")
	case "cancelled":
		query = query.Where("cancelled_at IS NOT NULL")
	case "scheduled":
		query = query.Where("is_published = true AND cancelled_at IS NULL AND tanggal_mulai > NOW()")
	case "active":
		query = query.Where("is_published = true AND cancelled_at IS NULL AND tanggal_mulai <= NOW() AND tanggal_selesai > NOW()")
	case "ended":
		query = query.Where("is_published = true AND cancelled_at IS NULL AND tanggal_selesai <= NOW()")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order(params.SortBy + " " + params.Order).Offset(params.GetOffset()).Limit(params.PerPage).Find(&campaigns).Error
	return campaigns, total, err
}

func (r *seasonalCampaignRepository) Save(ctx context.Context, campaign *models.SeasonalCampaign) error {
	return r.db.WithContext(ctx).Save(campaign).Error
}

func (r *seasonalCampaignRepository) Delete(ctx context.Context, campaign *models.SeasonalCampaign) error {
	return r.db.WithContext(ctx).Delete(campaign).Error
}

func (r *seasonalCampaignRepository) HasPublishedOverlap(ctx context.Context, start, end time.Time, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.SeasonalCampaign{}).
		Where("is_published = true AND cancelled_at IS NULL").
		Where("tanggal_mulai <= ? AND tanggal_selesai >= ?", end, start)
	if excludeID != "" {
		query = query.Where("id <> ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
