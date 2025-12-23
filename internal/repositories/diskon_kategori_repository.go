package repositories

import (
	"context"
	"time"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type DiskonKategoriRepository interface {
	Create(ctx context.Context, diskon *models.DiskonKategori) error
	FindByID(ctx context.Context, id string) (*models.DiskonKategori, error)
	FindAll(ctx context.Context, params *models.PaginationRequest, kategoriID string, berlakuHariIni bool) ([]models.DiskonKategori, int64, error)
	FindActiveByKategoriID(ctx context.Context, kategoriID string) (*models.DiskonKategori, error)
	Update(ctx context.Context, diskon *models.DiskonKategori) error
	Delete(ctx context.Context, id string) error
}

type diskonKategoriRepository struct {
	db *gorm.DB
}

func NewDiskonKategoriRepository(db *gorm.DB) DiskonKategoriRepository {
	return &diskonKategoriRepository{db: db}
}

func (r *diskonKategoriRepository) Create(ctx context.Context, diskon *models.DiskonKategori) error {
	return r.db.WithContext(ctx).Create(diskon).Error
}

func (r *diskonKategoriRepository) FindByID(ctx context.Context, id string) (*models.DiskonKategori, error) {
	var diskon models.DiskonKategori
	err := r.db.WithContext(ctx).Preload("Kategori").Where("id = ?", id).First(&diskon).Error
	if err != nil {
		return nil, err
	}
	return &diskon, nil
}


func (r *diskonKategoriRepository) FindAll(ctx context.Context, params *models.PaginationRequest, kategoriID string, berlakuHariIni bool) ([]models.DiskonKategori, int64, error) {
	var diskons []models.DiskonKategori
	var total int64

	query := r.db.WithContext(ctx).Model(&models.DiskonKategori{}).Preload("Kategori")

	if kategoriID != "" {
		query = query.Where("kategori_id = ?", kategoriID)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if berlakuHariIni {
		today := time.Now().Format("2006-01-02")
		query = query.Where("(tanggal_mulai IS NULL OR tanggal_mulai <= ?) AND (tanggal_selesai IS NULL OR tanggal_selesai >= ?)", today, today)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := params.UrutBerdasarkan + " " + params.Urutan
	query = query.Order(orderClause)
	query = query.Offset(params.GetOffset()).Limit(params.PerHalaman)

	if err := query.Find(&diskons).Error; err != nil {
		return nil, 0, err
	}

	return diskons, total, nil
}

func (r *diskonKategoriRepository) FindActiveByKategoriID(ctx context.Context, kategoriID string) (*models.DiskonKategori, error) {
	var diskon models.DiskonKategori
	today := time.Now().Format("2006-01-02")

	err := r.db.WithContext(ctx).
		Where("kategori_id = ?", kategoriID).
		Where("is_active = ?", true).
		Where("(tanggal_mulai IS NULL OR tanggal_mulai <= ?)", today).
		Where("(tanggal_selesai IS NULL OR tanggal_selesai >= ?)", today).
		Order("created_at DESC").
		First(&diskon).Error

	if err != nil {
		return nil, err
	}
	return &diskon, nil
}

func (r *diskonKategoriRepository) Update(ctx context.Context, diskon *models.DiskonKategori) error {
	return r.db.WithContext(ctx).Save(diskon).Error
}

func (r *diskonKategoriRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.DiskonKategori{}).Error
}
