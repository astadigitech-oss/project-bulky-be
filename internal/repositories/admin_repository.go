package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type AdminRepository interface {
	Create(ctx context.Context, admin *models.Admin) error
	FindByID(ctx context.Context, id string) (*models.Admin, error)
	FindByIDWithRole(ctx context.Context, id string) (*models.Admin, error)
	FindByEmail(ctx context.Context, email string) (*models.Admin, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.Admin, int64, error)
	Update(ctx context.Context, admin *models.Admin) error
	Delete(ctx context.Context, id string) error
	ExistsByEmail(ctx context.Context, email string, excludeID *string) (bool, error)
	Count(ctx context.Context) (int64, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) Create(ctx context.Context, admin *models.Admin) error {
	return r.db.WithContext(ctx).Create(admin).Error
}

func (r *adminRepository) FindByID(ctx context.Context, id string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) FindByIDWithRole(ctx context.Context, id string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.WithContext(ctx).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true)
		}).
		Preload("Role.Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Order("modul ASC, kode ASC")
		}).
		Where("id = ?", id).
		First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) FindByEmail(ctx context.Context, email string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.Admin, int64, error) {
	var admins []models.Admin
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Admin{})

	if params.Cari != "" {
		query = query.Where("nama ILIKE ? OR email ILIKE ?", "%"+params.Cari+"%", "%"+params.Cari+"%")
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := params.UrutBerdasarkan + " " + params.Urutan
	query = query.Order(orderClause)
	query = query.Offset(params.GetOffset()).Limit(params.PerHalaman)

	if err := query.Find(&admins).Error; err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

func (r *adminRepository) Update(ctx context.Context, admin *models.Admin) error {
	return r.db.WithContext(ctx).Save(admin).Error
}

func (r *adminRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Admin{}).Error
}

func (r *adminRepository) ExistsByEmail(ctx context.Context, email string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Admin{}).Where("email = ?", email)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *adminRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Admin{}).Count(&count).Error
	return count, err
}
