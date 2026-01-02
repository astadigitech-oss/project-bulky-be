package repositories

import (
	"context"
	"time"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type BuyerRepository interface {
	FindByID(ctx context.Context, id string) (*models.Buyer, error)
	FindByIDWithAlamat(ctx context.Context, id string) (*models.Buyer, error)
	FindAll(ctx context.Context, params *models.BuyerFilterRequest) ([]models.Buyer, int64, error)
	Update(ctx context.Context, buyer *models.Buyer) error
	Delete(ctx context.Context, id string) error
	ExistsByUsername(ctx context.Context, username string, excludeID *string) (bool, error)
	ExistsByEmail(ctx context.Context, email string, excludeID *string) (bool, error)
	CountAlamat(ctx context.Context, buyerID string) (int64, error)
	GetStatistik(ctx context.Context) (*models.BuyerStatistikResponse, error)
}

type buyerRepository struct {
	db *gorm.DB
}

func NewBuyerRepository(db *gorm.DB) BuyerRepository {
	return &buyerRepository{db: db}
}

func (r *buyerRepository) FindByID(ctx context.Context, id string) (*models.Buyer, error) {
	var buyer models.Buyer
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&buyer).Error
	return &buyer, err
}

func (r *buyerRepository) FindByIDWithAlamat(ctx context.Context, id string) (*models.Buyer, error) {
	var buyer models.Buyer
	err := r.db.WithContext(ctx).
		Preload("Alamat", "deleted_at IS NULL").
		Where("id = ?", id).First(&buyer).Error
	return &buyer, err
}

func (r *buyerRepository) FindAll(ctx context.Context, params *models.BuyerFilterRequest) ([]models.Buyer, int64, error) {
	var buyers []models.Buyer
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Buyer{})

	if params.Cari != "" {
		search := "%" + params.Cari + "%"
		query = query.Where("nama ILIKE ? OR username ILIKE ? OR email ILIKE ? OR telepon ILIKE ?", search, search, search, search)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if params.IsVerified != nil {
		query = query.Where("is_verified = ?", *params.IsVerified)
	}

	query.Count(&total)

	orderClause := params.UrutBerdasarkan + " " + params.Urutan
	err := query.Order(orderClause).
		Offset(params.GetOffset()).
		Limit(params.PerHalaman).
		Find(&buyers).Error

	return buyers, total, err
}

func (r *buyerRepository) Update(ctx context.Context, buyer *models.Buyer) error {
	return r.db.WithContext(ctx).Save(buyer).Error
}

func (r *buyerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Buyer{}, "id = ?", id).Error
}

func (r *buyerRepository) ExistsByUsername(ctx context.Context, username string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Buyer{}).Where("username = ?", username)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *buyerRepository) ExistsByEmail(ctx context.Context, email string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Buyer{}).Where("email = ?", email)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *buyerRepository) CountAlamat(ctx context.Context, buyerID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AlamatBuyer{}).Where("buyer_id = ?", buyerID).Count(&count).Error
	return count, err
}

func (r *buyerRepository) GetStatistik(ctx context.Context) (*models.BuyerStatistikResponse, error) {
	var stats models.BuyerStatistikResponse

	// Total buyer
	r.db.WithContext(ctx).Model(&models.Buyer{}).Count(&stats.TotalBuyer)

	// Buyer aktif
	r.db.WithContext(ctx).Model(&models.Buyer{}).Where("is_active = ?", true).Count(&stats.BuyerAktif)

	// Buyer nonaktif
	stats.BuyerNonaktif = stats.TotalBuyer - stats.BuyerAktif

	// Buyer verified
	r.db.WithContext(ctx).Model(&models.Buyer{}).Where("is_verified = ?", true).Count(&stats.BuyerVerified)

	// Buyer unverified
	stats.BuyerUnverified = stats.TotalBuyer - stats.BuyerVerified

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := startOfDay.AddDate(0, 0, -int(now.Weekday()))
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Registrasi hari ini
	r.db.WithContext(ctx).Model(&models.Buyer{}).Where("created_at >= ?", startOfDay).Count(&stats.RegistrasiHariIni)

	// Registrasi minggu ini
	r.db.WithContext(ctx).Model(&models.Buyer{}).Where("created_at >= ?", startOfWeek).Count(&stats.RegistrasiMingguIni)

	// Registrasi bulan ini
	r.db.WithContext(ctx).Model(&models.Buyer{}).Where("created_at >= ?", startOfMonth).Count(&stats.RegistrasiBulanIni)

	return &stats, nil
}
