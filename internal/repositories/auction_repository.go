package repositories

import (
	"context"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuctionRepository interface {
	// Batch
	ListBatches(ctx context.Context, params *dto.AuctionListQueryParams) ([]models.AuctionBatch, int64, error)
	FindBatchByID(ctx context.Context, id uuid.UUID) (*models.AuctionBatch, error)
	FindBatchByIDWithRelations(ctx context.Context, id uuid.UUID) (*models.AuctionBatch, error)
	CreateBatch(ctx context.Context, batch *models.AuctionBatch) error
	UpdateBatch(ctx context.Context, batch *models.AuctionBatch) error
	ReplaceBatchItems(ctx context.Context, batchID uuid.UUID, items []models.AuctionBatchItem) error
	ReplaceBatchBrands(ctx context.Context, batchID uuid.UUID, merekIDs []uuid.UUID) error
	ReplaceBatchAssets(ctx context.Context, batchID uuid.UUID, assetIDs []uuid.UUID) error
	GetBatchAssets(ctx context.Context, batchID uuid.UUID) ([]models.AuctionBatchAsset, error)
	GetBatchAssetDetails(ctx context.Context, batchID uuid.UUID) ([]models.AuctionAsset, error)

	// Asset
	CreateAsset(ctx context.Context, asset *models.AuctionAsset) error
	FindAssetByID(ctx context.Context, id uuid.UUID) (*models.AuctionAsset, error)

	// Bids
	ListBids(ctx context.Context, batchID uuid.UUID, params *dto.AuctionBidsQueryParams) ([]models.AuctionBid, int64, error)
	FindBidByID(ctx context.Context, id uuid.UUID) (*models.AuctionBid, error)
	FindBidByBatchAndID(ctx context.Context, batchID, bidID uuid.UUID) (*models.AuctionBid, error)

	// Winner
	FindWinnerByBatchID(ctx context.Context, batchID uuid.UUID) (*models.AuctionWinner, error)

	// Reservations
	GetActiveReservationsForProducts(ctx context.Context, produkIDs []uuid.UUID) ([]models.AuctionStockReservation, error)

	// Audit
	CreateAuditLog(ctx context.Context, log *models.AuctionAuditLog) error

	// Idempotency
	FindIdempotency(ctx context.Context, actorType string, actorID uuid.UUID, operation, key string) (*models.AuctionIdempotency, error)
	CreateIdempotency(ctx context.Context, rec *models.AuctionIdempotency) error

	// Analytics
	GetBatchAnalytics(ctx context.Context, batchID uuid.UUID) (*AuctionAnalytics, error)
	GetBatchesAnalytics(ctx context.Context, batchIDs []uuid.UUID) (map[uuid.UUID]*AuctionAnalytics, error)
	GetBatchViewCount(ctx context.Context, batchID uuid.UUID) (int64, error)
	GetFirstImagesForBatches(ctx context.Context, batchIDs []uuid.UUID) (map[uuid.UUID]*models.AuctionAsset, error)

	// Product options
	ListProductOptions(ctx context.Context, params *dto.AuctionProductOptionsQueryParams) ([]models.Produk, int64, error)
}

type auctionRepository struct {
	db *gorm.DB
}

func NewAuctionRepository(db *gorm.DB) AuctionRepository {
	return &auctionRepository{db: db}
}

// ========================================
// Batch
// ========================================

func (r *auctionRepository) ListBatches(ctx context.Context, params *dto.AuctionListQueryParams) ([]models.AuctionBatch, int64, error) {
	var batches []models.AuctionBatch
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AuctionBatch{})

	if params.Search != "" {
		query = query.Where("nama_id ILIKE ? OR code ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "created_at DESC, id DESC"
	if params.SortBy == "created_at_asc" {
		order = "created_at ASC, id ASC"
	}

	if err := query.
		Order(order).
		Offset((params.Page - 1) * params.PerPage).
		Limit(params.PerPage).
		Find(&batches).Error; err != nil {
		return nil, 0, err
	}

	return batches, total, nil
}

func (r *auctionRepository) FindBatchByID(ctx context.Context, id uuid.UUID) (*models.AuctionBatch, error) {
	var batch models.AuctionBatch
	if err := r.db.WithContext(ctx).First(&batch, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *auctionRepository) FindBatchByIDWithRelations(ctx context.Context, id uuid.UUID) (*models.AuctionBatch, error) {
	var batch models.AuctionBatch
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Brands").
		Preload("BatchAssets").
		Preload("Winner").
		First(&batch, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *auctionRepository) CreateBatch(ctx context.Context, batch *models.AuctionBatch) error {
	return r.db.WithContext(ctx).Create(batch).Error
}

func (r *auctionRepository) UpdateBatch(ctx context.Context, batch *models.AuctionBatch) error {
	return r.db.WithContext(ctx).Save(batch).Error
}

func (r *auctionRepository) ReplaceBatchItems(ctx context.Context, batchID uuid.UUID, items []models.AuctionBatchItem) error {
	if err := r.db.WithContext(ctx).Where("batch_id = ?", batchID).Delete(&models.AuctionBatchItem{}).Error; err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *auctionRepository) ReplaceBatchBrands(ctx context.Context, batchID uuid.UUID, merekIDs []uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("batch_id = ?", batchID).Delete(&models.AuctionBatchBrand{}).Error; err != nil {
		return err
	}
	if len(merekIDs) == 0 {
		return nil
	}
	brands := make([]models.AuctionBatchBrand, 0, len(merekIDs))
	for _, m := range merekIDs {
		brands = append(brands, models.AuctionBatchBrand{BatchID: batchID, MerekID: m})
	}
	return r.db.WithContext(ctx).Create(&brands).Error
}

func (r *auctionRepository) ReplaceBatchAssets(ctx context.Context, batchID uuid.UUID, assetIDs []uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("batch_id = ?", batchID).Delete(&models.AuctionBatchAsset{}).Error; err != nil {
		return err
	}
	if len(assetIDs) == 0 {
		return nil
	}
	assets := make([]models.AuctionBatchAsset, 0, len(assetIDs))
	for i, a := range assetIDs {
		assets = append(assets, models.AuctionBatchAsset{BatchID: batchID, AssetID: a, SortOrder: i})
	}
	return r.db.WithContext(ctx).Create(&assets).Error
}

func (r *auctionRepository) GetBatchAssets(ctx context.Context, batchID uuid.UUID) ([]models.AuctionBatchAsset, error) {
	var assets []models.AuctionBatchAsset
	if err := r.db.WithContext(ctx).Where("batch_id = ?", batchID).Order("sort_order ASC").Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *auctionRepository) GetBatchAssetDetails(ctx context.Context, batchID uuid.UUID) ([]models.AuctionAsset, error) {
	var assets []models.AuctionAsset
	err := r.db.WithContext(ctx).
		Joins("JOIN auction_batch_assets aba ON aba.asset_id = auction_assets.id AND aba.batch_id = ?", batchID).
		Order("aba.sort_order ASC").
		Find(&assets).Error
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// ========================================
// Asset
// ========================================

func (r *auctionRepository) CreateAsset(ctx context.Context, asset *models.AuctionAsset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *auctionRepository) FindAssetByID(ctx context.Context, id uuid.UUID) (*models.AuctionAsset, error) {
	var asset models.AuctionAsset
	if err := r.db.WithContext(ctx).First(&asset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

// ========================================
// Bids
// ========================================

func (r *auctionRepository) ListBids(ctx context.Context, batchID uuid.UUID, params *dto.AuctionBidsQueryParams) ([]models.AuctionBid, int64, error) {
	var bids []models.AuctionBid
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AuctionBid{}).Where("batch_id = ?", batchID)

	if params.Search != "" {
		query = query.Joins("JOIN buyer ON buyer.id = auction_bids.buyer_id").
			Where("buyer.nama ILIKE ? OR buyer.telepon ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}
	if params.BuyerID != "" {
		if bidID, err := uuid.Parse(params.BuyerID); err == nil {
			query = query.Where("buyer_id = ?", bidID)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "created_at DESC, id DESC"
	switch params.SortBy {
	case "amount_desc":
		order = "amount DESC, id DESC"
	case "amount_asc":
		order = "amount ASC, id ASC"
	}

	if err := query.
		Preload("Buyer").
		Order(order).
		Offset((params.Page - 1) * params.PerPage).
		Limit(params.PerPage).
		Find(&bids).Error; err != nil {
		return nil, 0, err
	}

	return bids, total, nil
}

func (r *auctionRepository) FindBidByID(ctx context.Context, id uuid.UUID) (*models.AuctionBid, error) {
	var bid models.AuctionBid
	if err := r.db.WithContext(ctx).Preload("Buyer").First(&bid, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &bid, nil
}

func (r *auctionRepository) FindBidByBatchAndID(ctx context.Context, batchID, bidID uuid.UUID) (*models.AuctionBid, error) {
	var bid models.AuctionBid
	if err := r.db.WithContext(ctx).Preload("Buyer").First(&bid, "batch_id = ? AND id = ?", batchID, bidID).Error; err != nil {
		return nil, err
	}
	return &bid, nil
}

// ========================================
// Winner
// ========================================

func (r *auctionRepository) FindWinnerByBatchID(ctx context.Context, batchID uuid.UUID) (*models.AuctionWinner, error) {
	var winner models.AuctionWinner
	if err := r.db.WithContext(ctx).First(&winner, "batch_id = ?", batchID).Error; err != nil {
		return nil, err
	}
	return &winner, nil
}

// ========================================
// Reservations
// ========================================

func (r *auctionRepository) GetActiveReservationsForProducts(ctx context.Context, produkIDs []uuid.UUID) ([]models.AuctionStockReservation, error) {
	var reservations []models.AuctionStockReservation
	if len(produkIDs) == 0 {
		return reservations, nil
	}
	if err := r.db.WithContext(ctx).
		Where("produk_id IN ? AND status = ?", produkIDs, "ACTIVE").
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

// ========================================
// Audit
// ========================================

func (r *auctionRepository) CreateAuditLog(ctx context.Context, log *models.AuctionAuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// ========================================
// Idempotency
// ========================================

func (r *auctionRepository) FindIdempotency(ctx context.Context, actorType string, actorID uuid.UUID, operation, key string) (*models.AuctionIdempotency, error) {
	var rec models.AuctionIdempotency
	if err := r.db.WithContext(ctx).
		Where("actor_type = ? AND actor_id = ? AND operation = ? AND key = ?", actorType, actorID, operation, key).
		First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *auctionRepository) CreateIdempotency(ctx context.Context, rec *models.AuctionIdempotency) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

// ========================================
// Analytics
// ========================================

// AuctionAnalytics hasil agregasi untuk satu batch.
type AuctionAnalytics struct {
	BidderCount        int64
	BidCount           int64
	RepeatBidderCount  int64
	AdditionalBidCount int64
	HighestBid         *float64
	LowestBid          *float64
	TotalQuantity      int64
}

func (r *auctionRepository) GetBatchAnalytics(ctx context.Context, batchID uuid.UUID) (*AuctionAnalytics, error) {
	var res struct {
		BidderCount int64
		BidCount    int64
		HighestBid  *float64
		LowestBid   *float64
	}
	err := r.db.WithContext(ctx).Model(&models.AuctionBid{}).
		Where("batch_id = ?", batchID).
		Select(
			"COUNT(DISTINCT buyer_id) AS bidder_count",
			"COUNT(*) AS bid_count",
			"MAX(amount) AS highest_bid",
			"MIN(amount) AS lowest_bid",
		).
		Scan(&res).Error
	if err != nil {
		return nil, err
	}

	// Repeat bidder = buyer dengan lebih dari satu bid.
	var repeatBidder int64
	err = r.db.WithContext(ctx).Model(&models.AuctionBid{}).
		Where("batch_id = ?", batchID).
		Group("buyer_id").
		Having("COUNT(*) > 1").
		Count(&repeatBidder).Error
	if err != nil {
		return nil, err
	}

	return &AuctionAnalytics{
		BidderCount:        res.BidderCount,
		BidCount:           res.BidCount,
		RepeatBidderCount:  repeatBidder,
		AdditionalBidCount: res.BidCount - res.BidderCount,
		HighestBid:         res.HighestBid,
		LowestBid:          res.LowestBid,
	}, nil
}

// GetBatchesAnalytics menghitung agregasi bid untuk banyak batch sekaligus
// agar list tidak melakukan query per baris (N+1).
func (r *auctionRepository) GetBatchesAnalytics(ctx context.Context, batchIDs []uuid.UUID) (map[uuid.UUID]*AuctionAnalytics, error) {
	result := make(map[uuid.UUID]*AuctionAnalytics)
	if len(batchIDs) == 0 {
		return result, nil
	}

	var rows []struct {
		BatchID     uuid.UUID
		BidderCount int64
		BidCount    int64
		HighestBid  *float64
		LowestBid   *float64
	}
	err := r.db.WithContext(ctx).Model(&models.AuctionBid{}).
		Where("batch_id IN ?", batchIDs).
		Select(
			"batch_id",
			"COUNT(DISTINCT buyer_id) AS bidder_count",
			"COUNT(*) AS bid_count",
			"MAX(amount) AS highest_bid",
			"MIN(amount) AS lowest_bid",
		).
		Group("batch_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.BatchID] = &AuctionAnalytics{
			BidderCount: row.BidderCount,
			BidCount:    row.BidCount,
			HighestBid:  row.HighestBid,
			LowestBid:   row.LowestBid,
		}
	}

	// Repeat bidder per batch.
	var repeatRows []struct {
		BatchID uuid.UUID
		Cnt     int64
	}
	err = r.db.WithContext(ctx).Model(&models.AuctionBid{}).
		Where("batch_id IN ?", batchIDs).
		Select("batch_id", "COUNT(*) AS cnt").
		Group("batch_id, buyer_id").
		Having("COUNT(*) > 1").
		Scan(&repeatRows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range repeatRows {
		if analytics, ok := result[row.BatchID]; ok {
			analytics.RepeatBidderCount++
		}
	}
	for _, analytics := range result {
		analytics.AdditionalBidCount = analytics.BidCount - analytics.BidderCount
	}

	return result, nil
}

func (r *auctionRepository) GetBatchViewCount(ctx context.Context, batchID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AuctionEvent{}).
		Where("batch_id = ? AND event_type = ?", batchID, "BATCH_VIEW").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetFirstImagesForBatches mengambil gambar pertama (sort_order terkecil)
// untuk setiap batch agar list dapat menampilkan thumbnail tanpa N+1.
func (r *auctionRepository) GetFirstImagesForBatches(ctx context.Context, batchIDs []uuid.UUID) (map[uuid.UUID]*models.AuctionAsset, error) {
	result := make(map[uuid.UUID]*models.AuctionAsset)
	if len(batchIDs) == 0 {
		return result, nil
	}

	// Ambil pivot (batch_id, asset_id, sort_order) diurutkan per batch.
	var pivots []models.AuctionBatchAsset
	if err := r.db.WithContext(ctx).
		Where("batch_id IN ?", batchIDs).
		Order("batch_id ASC, sort_order ASC").
		Find(&pivots).Error; err != nil {
		return nil, err
	}

	// Kumpulkan asset id pertama per batch yang bertipe IMAGE.
	firstAssetIDs := make(map[uuid.UUID]uuid.UUID)
	for _, p := range pivots {
		if _, ok := firstAssetIDs[p.BatchID]; ok {
			continue
		}
		firstAssetIDs[p.BatchID] = p.AssetID
	}

	allAssetIDs := make([]uuid.UUID, 0, len(firstAssetIDs))
	for _, id := range firstAssetIDs {
		allAssetIDs = append(allAssetIDs, id)
	}
	if len(allAssetIDs) == 0 {
		return result, nil
	}

	var assets []models.AuctionAsset
	if err := r.db.WithContext(ctx).
		Where("id IN ? AND kind = ?", allAssetIDs, "IMAGE").
		Find(&assets).Error; err != nil {
		return nil, err
	}
	assetByID := make(map[uuid.UUID]models.AuctionAsset)
	for _, a := range assets {
		assetByID[a.ID] = a
	}

	for batchID, assetID := range firstAssetIDs {
		if a, ok := assetByID[assetID]; ok {
			cp := a
			result[batchID] = &cp
		}
	}

	return result, nil
}

// ========================================
// Product options
// ========================================

func (r *auctionRepository) ListProductOptions(ctx context.Context, params *dto.AuctionProductOptionsQueryParams) ([]models.Produk, int64, error) {
	var produk []models.Produk
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Produk{}).Where("deleted_at IS NULL")

	if params.Search != "" {
		query = query.Where("nama_id ILIKE ? OR nama_en ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}
	if params.WarehouseID != "" {
		if wid, err := uuid.Parse(params.WarehouseID); err == nil {
			query = query.Where("warehouse_id = ?", wid)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "created_at DESC, id DESC"
	if params.SortBy == "nama_asc" {
		order = "nama_id ASC, id ASC"
	}

	if err := query.
		Order(order).
		Offset((params.Page - 1) * params.PerPage).
		Limit(params.PerPage).
		Find(&produk).Error; err != nil {
		return nil, 0, err
	}

	return produk, total, nil
}
