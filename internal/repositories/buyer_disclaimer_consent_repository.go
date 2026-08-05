package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type BuyerDisclaimerConsentRepository interface {
	FindByPesananID(ctx context.Context, pesananID string) (*models.BuyerDisclaimerConsent, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.BuyerDisclaimerConsent, int64, error)
}

type buyerDisclaimerConsentRepository struct {
	db *gorm.DB
}

func NewBuyerDisclaimerConsentRepository(db *gorm.DB) BuyerDisclaimerConsentRepository {
	return &buyerDisclaimerConsentRepository{db: db}
}

func (r *buyerDisclaimerConsentRepository) FindByPesananID(ctx context.Context, pesananID string) (*models.BuyerDisclaimerConsent, error) {
	var consent models.BuyerDisclaimerConsent
	err := r.db.WithContext(ctx).
		Preload("Buyer").
		Preload("Pesanan").
		Preload("Disclaimer").
		Where("pesanan_id = ?", pesananID).
		First(&consent).Error
	if err != nil {
		return nil, err
	}
	return &consent, nil
}

func (r *buyerDisclaimerConsentRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.BuyerDisclaimerConsent, int64, error) {
	var consents []models.BuyerDisclaimerConsent
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.BuyerDisclaimerConsent{}).
		Preload("Buyer").
		Preload("Pesanan").
		Preload("Disclaimer")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("disetujui_at DESC").
		Offset(params.GetOffset()).
		Limit(params.PerPage).
		Find(&consents).Error; err != nil {
		return nil, 0, err
	}

	return consents, total, nil
}
