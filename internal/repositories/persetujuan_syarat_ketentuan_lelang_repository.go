package repositories

import (
	"context"
	"strings"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type PersetujuanSyaratKetentuanLelangRepository interface {
	AmbilSemua(ctx context.Context, params *models.PaginationRequest) ([]models.PersetujuanSyaratKetentuanLelang, int64, error)
	AmbilBerdasarkanID(ctx context.Context, id string) (*models.PersetujuanSyaratKetentuanLelang, error)
}

type persetujuanSyaratKetentuanLelangRepository struct {
	db *gorm.DB
}

func NewPersetujuanSyaratKetentuanLelangRepository(db *gorm.DB) PersetujuanSyaratKetentuanLelangRepository {
	return &persetujuanSyaratKetentuanLelangRepository{db: db}
}

func (r *persetujuanSyaratKetentuanLelangRepository) AmbilSemua(
	ctx context.Context,
	params *models.PaginationRequest,
) ([]models.PersetujuanSyaratKetentuanLelang, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.PersetujuanSyaratKetentuanLelang{}).
		Joins("JOIN buyer ON buyer.id = buyer_auction_terms_consent.buyer_id").
		Joins("JOIN auction_bids ON auction_bids.id = buyer_auction_terms_consent.bid_id").
		Joins("JOIN auction_batches ON auction_batches.id = auction_bids.batch_id").
		Preload("Buyer").
		Preload("Bid").
		Preload("Bid.Batch")

	if search := strings.TrimSpace(params.Search); search != "" {
		pattern := "%" + search + "%"
		query = query.Where(
			"buyer.nama ILIKE ? OR COALESCE(buyer.email, '') ILIKE ? OR auction_batches.code ILIKE ? OR auction_batches.nama_id ILIKE ? OR auction_bids.id::text ILIKE ?",
			pattern, pattern, pattern, pattern, pattern,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderByColumns := map[string]string{
		"disetujui_at": "buyer_auction_terms_consent.disetujui_at",
		"buyer_nama":   "buyer.nama",
		"batch_kode":   "auction_batches.code",
		"bid_amount":   "auction_bids.amount",
		"created_at":   "buyer_auction_terms_consent.created_at",
	}
	orderBy := orderByColumns[params.SortBy]
	if orderBy == "" {
		orderBy = orderByColumns["disetujui_at"]
	}
	direction := "DESC"
	if strings.EqualFold(params.Order, "asc") {
		direction = "ASC"
	}

	var consents []models.PersetujuanSyaratKetentuanLelang
	if err := query.
		Order(orderBy + " " + direction).
		Order("buyer_auction_terms_consent.id DESC").
		Offset(params.GetOffset()).
		Limit(params.PerPage).
		Find(&consents).Error; err != nil {
		return nil, 0, err
	}
	return consents, total, nil
}

func (r *persetujuanSyaratKetentuanLelangRepository) AmbilBerdasarkanID(
	ctx context.Context,
	id string,
) (*models.PersetujuanSyaratKetentuanLelang, error) {
	var consent models.PersetujuanSyaratKetentuanLelang
	err := r.db.WithContext(ctx).
		Preload("Buyer").
		Preload("Bid").
		Preload("Bid.Batch").
		First(&consent, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &consent, nil
}
