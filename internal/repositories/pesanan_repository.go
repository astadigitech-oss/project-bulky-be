package repositories

import (
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PesananRepository interface {
	FindByID(id uuid.UUID) (*models.Pesanan, error)
	FindByBuyerID(buyerID uuid.UUID) ([]models.Pesanan, error)
}

type pesananRepository struct {
	db *gorm.DB
}

func NewPesananRepository(db *gorm.DB) PesananRepository {
	return &pesananRepository{db: db}
}

func (r *pesananRepository) FindByID(id uuid.UUID) (*models.Pesanan, error) {
	var pesanan models.Pesanan
	if err := r.db.First(&pesanan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &pesanan, nil
}

func (r *pesananRepository) FindByBuyerID(buyerID uuid.UUID) ([]models.Pesanan, error) {
	var pesanan []models.Pesanan
	if err := r.db.Where("buyer_id = ?", buyerID).Find(&pesanan).Error; err != nil {
		return nil, err
	}
	return pesanan, nil
}
