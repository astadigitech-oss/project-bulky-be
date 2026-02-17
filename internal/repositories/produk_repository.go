package repositories

import (
	"context"
	"fmt"
	"time"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type ProdukRepository interface {
	Create(ctx context.Context, produk *models.Produk) error
	FindByID(ctx context.Context, id string) (*models.Produk, error)
	FindBySlug(ctx context.Context, slug string) (*models.Produk, error)
	FindAll(ctx context.Context, params *models.ProdukFilterRequest) ([]models.Produk, int64, error)
	Update(ctx context.Context, produk *models.Produk) error
	Delete(ctx context.Context, produk *models.Produk) error
	ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error)
	ExistsByIDCargo(ctx context.Context, idCargo string, excludeID *string) (bool, error)
}

type produkRepository struct {
	db *gorm.DB
}

func NewProdukRepository(db *gorm.DB) ProdukRepository {
	return &produkRepository{db: db}
}

func (r *produkRepository) Create(ctx context.Context, produk *models.Produk) error {
	return r.db.WithContext(ctx).Create(produk).Error
}

func (r *produkRepository) FindByID(ctx context.Context, id string) (*models.Produk, error) {
	var produk models.Produk
	err := r.db.WithContext(ctx).
		Preload("Kategori").
		Preload("Mereks").
		Preload("Kondisi").
		Preload("KondisiPaket").
		Preload("Sumber").
		Preload("Warehouse").
		Preload("TipeProduk").
		Preload("Gambar", func(db *gorm.DB) *gorm.DB {
			return db.Order("urutan ASC")
		}).
		Preload("Dokumen").
		Where("id = ?", id).
		First(&produk).Error
	if err != nil {
		return nil, err
	}
	return &produk, nil
}

func (r *produkRepository) FindBySlug(ctx context.Context, slug string) (*models.Produk, error) {
	var produk models.Produk
	err := r.db.WithContext(ctx).
		Preload("Kategori").
		Preload("Mereks").
		Preload("Kondisi").
		Preload("KondisiPaket").
		Preload("Sumber").
		Preload("Warehouse").
		Preload("TipeProduk").
		Preload("Gambar", func(db *gorm.DB) *gorm.DB {
			return db.Order("urutan ASC")
		}).
		Preload("Dokumen").
		Where("slug = ?", slug).
		First(&produk).Error
	if err != nil {
		return nil, err
	}
	return &produk, nil
}

func (r *produkRepository) FindAll(ctx context.Context, params *models.ProdukFilterRequest) ([]models.Produk, int64, error) {
	var produks []models.Produk
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Produk{}).
		Preload("Kategori").
		Preload("Mereks").
		Preload("Kondisi").
		Preload("KondisiPaket").
		Preload("Sumber").
		Preload("Warehouse").
		Preload("TipeProduk").
		Preload("Gambar", func(db *gorm.DB) *gorm.DB {
			return db.Order("urutan ASC")
		}).
		Preload("Dokumen")

	// Apply filters
	if params.Search != "" {
		query = query.Where("nama_id ILIKE ? OR nama_en ILIKE ? OR id_cargo ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%")
	}
	if params.KategoriID != "" {
		query = query.Where("kategori_id = ?", params.KategoriID)
	}
	if params.MerekID != "" {
		// Filter by merek using many-to-many join
		query = query.Joins("JOIN produk_merek ON produk_merek.produk_id = produk.id").
			Where("produk_merek.merek_id = ?", params.MerekID)
	}
	if params.KondisiID != "" {
		query = query.Where("kondisi_id = ?", params.KondisiID)
	}
	if params.KondisiPaketID != "" {
		query = query.Where("kondisi_paket_id = ?", params.KondisiPaketID)
	}
	if params.SumberID != "" {
		query = query.Where("sumber_id = ?", params.SumberID)
	}
	if params.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", params.WarehouseID)
	}
	if params.TipeProdukID != "" {
		query = query.Where("tipe_produk_id = ?", params.TipeProdukID)
	}
	if params.HargaMin > 0 {
		query = query.Where("harga_sesudah_diskon >= ?", params.HargaMin)
	}
	if params.HargaMax > 0 {
		query = query.Where("harga_sesudah_diskon <= ?", params.HargaMax)
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := params.SortBy + " " + params.Order
	query = query.Order(orderClause)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Find(&produks).Error; err != nil {
		return nil, 0, err
	}

	return produks, total, nil
}

func (r *produkRepository) Update(ctx context.Context, produk *models.Produk) error {
	return r.db.WithContext(ctx).Save(produk).Error
}

func (r *produkRepository) Delete(ctx context.Context, produk *models.Produk) error {
	// Manual update slug untuk soft delete
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		deletedSlug := fmt.Sprintf("%s-deleted-%d%06d",
			produk.Slug,
			now.Unix(),
			now.Nanosecond()/1000,
		)

		if err := tx.Model(produk).Updates(map[string]interface{}{
			"slug":       deletedSlug,
			"deleted_at": now,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *produkRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Produk{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *produkRepository) ExistsByIDCargo(ctx context.Context, idCargo string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Produk{}).Where("id_cargo = ?", idCargo)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
