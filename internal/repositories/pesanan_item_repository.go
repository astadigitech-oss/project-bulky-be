package repositories

import (
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PesananItemRepository interface {
	FindByID(id uuid.UUID) (*models.PesananItem, error)
	FindByPesananID(pesananID uuid.UUID) ([]models.PesananItem, error)
}

type pesananItemRepository struct {
	db *gorm.DB
}

func NewPesananItemRepository(db *gorm.DB) PesananItemRepository {
	return &pesananItemRepository{db: db}
}

func (r *pesananItemRepository) FindByID(id uuid.UUID) (*models.PesananItem, error) {
	var item models.PesananItem
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *pesananItemRepository) FindByPesananID(pesananID uuid.UUID) ([]models.PesananItem, error) {
	var items []models.PesananItem
	if err := r.db.Where("pesanan_id = ?", pesananID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
