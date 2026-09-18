package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"project-bulky-be/internal/models"
)

type AuctionEducationBannerRepository interface {
	Create(context.Context, *models.AuctionEducationBanner) error
	FindByID(context.Context, uuid.UUID) (*models.AuctionEducationBanner, error)
	List(context.Context, string, int, int) ([]models.AuctionEducationBanner, int64, error)
	Save(context.Context, *models.AuctionEducationBanner) error
	Reorder(context.Context, []uuid.UUID) ([]models.AuctionEducationBanner, error)
	Delete(context.Context, *models.AuctionEducationBanner) error
}
type auctionEducationBannerRepository struct{ db *gorm.DB }

func NewAuctionEducationBannerRepository(db *gorm.DB) AuctionEducationBannerRepository {
	return &auctionEducationBannerRepository{db}
}
func (r *auctionEducationBannerRepository) Create(c context.Context, b *models.AuctionEducationBanner) error {
	return r.db.WithContext(c).Create(b).Error
}
func (r *auctionEducationBannerRepository) FindByID(c context.Context, id uuid.UUID) (*models.AuctionEducationBanner, error) {
	var b models.AuctionEducationBanner
	err := r.db.WithContext(c).First(&b, "id = ?", id).Error
	return &b, err
}
func (r *auctionEducationBannerRepository) List(c context.Context, search string, page, perPage int) ([]models.AuctionEducationBanner, int64, error) {
	var b []models.AuctionEducationBanner
	q := r.db.WithContext(c).Model(&models.AuctionEducationBanner{})
	if search != "" {
		q = q.Where("nama ILIKE ?", "%"+search+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("urutan ASC, id ASC").Offset((page - 1) * perPage).Limit(perPage).Find(&b).Error
	return b, total, err
}
func (r *auctionEducationBannerRepository) Save(c context.Context, b *models.AuctionEducationBanner) error {
	return r.db.WithContext(c).Save(b).Error
}
func (r *auctionEducationBannerRepository) Reorder(c context.Context, ids []uuid.UUID) ([]models.AuctionEducationBanner, error) {
	var banners []models.AuctionEducationBanner
	err := r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Order("urutan ASC, id ASC").Find(&banners).Error; err != nil {
			return err
		}
		if len(banners) != len(ids) {
			return gorm.ErrRecordNotFound
		}
		existing := make(map[uuid.UUID]struct{}, len(banners))
		for _, banner := range banners {
			existing[banner.ID] = struct{}{}
		}
		for _, bannerID := range ids {
			if _, ok := existing[bannerID]; !ok {
				return gorm.ErrRecordNotFound
			}
		}
		for index, bannerID := range ids {
			if err := tx.Model(&models.AuctionEducationBanner{}).Where("id = ?", bannerID).Update("urutan", index).Error; err != nil {
				return err
			}
		}
		return tx.Order("urutan ASC, id ASC").Find(&banners).Error
	})
	return banners, err
}
func (r *auctionEducationBannerRepository) Delete(c context.Context, b *models.AuctionEducationBanner) error {
	return r.db.WithContext(c).Delete(b).Error
}
