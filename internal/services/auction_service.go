package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ============================================================
// Typed error untuk mapping status HTTP di controller.
// ============================================================

type AuctionError struct {
	Status      int
	Message     string
	FieldErrors []models.FieldError
}

func (e *AuctionError) Error() string { return e.Message }

func auctionErr(status int, message string) *AuctionError {
	return &AuctionError{Status: status, Message: message}
}

func auctionErrFields(status int, message string, fields []models.FieldError) *AuctionError {
	return &AuctionError{Status: status, Message: message, FieldErrors: fields}
}

// ============================================================
// Service interface
// ============================================================

type AuctionService interface {
	CreateDraft(ctx context.Context, req *dto.AuctionDraftInput, adminID uuid.UUID, idempotencyKey string) (*dto.AuctionBatchDetail, error)
	UpdateDraft(ctx context.Context, id uuid.UUID, req *dto.AuctionDraftInput, adminID uuid.UUID, version int) (*dto.AuctionBatchDetail, error)
	DeleteDraft(ctx context.Context, id uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID, version int, adminID uuid.UUID, idempotencyKey string) (*dto.AuctionBatchDetail, error)
	ListBatches(ctx context.Context, params *dto.AuctionListQueryParams) ([]dto.AuctionBatchSummary, *models.PaginationMeta, error)
	GetBatchDetail(ctx context.Context, id uuid.UUID) (*dto.AuctionBatchDetail, error)
	ListBids(ctx context.Context, batchID uuid.UUID, params *dto.AuctionBidsQueryParams) ([]dto.AuctionBidDetail, *models.PaginationMeta, error)
	SelectWinner(ctx context.Context, id uuid.UUID, req *dto.AuctionWinnerRequest, adminID uuid.UUID, idempotencyKey string) (*dto.AuctionBatchDetail, error)
	UpdateOperations(ctx context.Context, id uuid.UUID, req *dto.AuctionOperationRequest, adminID uuid.UUID, idempotencyKey string) (*dto.AuctionBatchDetail, error)
	UploadAsset(ctx context.Context, file *multipart.FileHeader, kind string, adminID uuid.UUID) (*dto.AuctionAssetResponse, error)
	PreviewSupplierExcel(ctx context.Context, file *multipart.FileHeader) (*dto.AuctionSupplierExcelPreview, error)
	ImportSupplierExcel(ctx context.Context, file *multipart.FileHeader, mapping dto.AuctionSupplierExcelMapping, title string, adminID uuid.UUID) (*dto.AuctionSupplierExcelImport, error)
	ListProductOptions(ctx context.Context, params *dto.AuctionProductOptionsQueryParams) ([]dto.AuctionProductOption, *models.PaginationMeta, error)
}

// DeleteDraft menghapus batch yang belum pernah dipublikasikan. Batch dengan
// status lain harus tetap tersimpan karena dapat memiliki bid, pemenang, atau
// histori operasional yang tidak boleh dihilangkan.
func (s *auctionService) DeleteDraft(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var batch models.AuctionBatch
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&batch, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return auctionErr(404, "Batch tidak ditemukan")
			}
			return err
		}
		if batch.Status != models.AuctionBatchStatusDRAFT {
			return auctionErr(409, "Hanya batch berstatus DRAFT yang dapat dihapus")
		}

		// Seluruh relasi ini memakai ON DELETE RESTRICT. Bersihkan hanya relasi
		// draft yang aman dihapus, sedangkan file asset dibiarkan sebagai asset
		// upload terpisah agar tidak berisiko menghapus file yang masih dipakai.
		for _, relation := range []interface{}{
			&models.AuctionAuditLog{},
			&models.AuctionBatchAsset{},
			&models.AuctionBatchBrand{},
			&models.AuctionBatchItem{},
		} {
			if err := tx.Where("batch_id = ?", id).Delete(relation).Error; err != nil {
				return err
			}
		}

		return tx.Delete(&batch).Error
	})
}

type auctionService struct {
	repo repositories.AuctionRepository
	db   *gorm.DB
	cfg  *config.Config
}

func NewAuctionService(repo repositories.AuctionRepository, db *gorm.DB, cfg *config.Config) AuctionService {
	return &auctionService{repo: repo, db: db, cfg: cfg}
}

// ============================================================
// Create Draft
// ============================================================

func (s *auctionService) CreateDraft(ctx context.Context, req *dto.AuctionDraftInput, adminID uuid.UUID, idempotencyKey string) (*dto.AuctionBatchDetail, error) {
	operation := "create_batch"
	if idempotencyKey != "" {
		replay, bid, err := s.checkIdempotency(ctx, "ADMIN", adminID, operation, uuid.Nil, idempotencyKey, req)
		if err != nil {
			return nil, err
		}
		if replay {
			if bid == uuid.Nil {
				return nil, auctionErr(409, "Idempotency key tidak valid")
			}
			return s.GetBatchDetail(ctx, bid)
		}
	}

	// Snapshot item & hitung turunan.
	items, grandTotal, totalQty, err := s.buildItemsSnapshot(ctx, req.Items)
	if err != nil {
		return nil, err
	}

	bOriginType := originType(req.OriginType)
	var warehouseID *uuid.UUID
	var suppName, suppAddress, suppProvinsi, suppKota, suppKecamatan, suppKelurahan, suppKodePos *string
	var suppLat, suppLng *decimal.Decimal

	if bOriginType == "SUPPLIER" {
		warehouseID = nil
		suppName = req.SupplierName
		suppAddress = req.SupplierAddress
		suppProvinsi = req.SupplierProvinsi
		suppKota = req.SupplierKota
		suppKecamatan = req.SupplierKecamatan
		suppKelurahan = req.SupplierKelurahan
		suppKodePos = req.SupplierKodePos
		suppLat = parseDecimalPtr(req.SupplierLatitude)
		suppLng = parseDecimalPtr(req.SupplierLongitude)
	} else {
		warehouseID = parseUUIDPtr(req.WarehouseID)
	}

	batchID := uuid.New()
	batch := &models.AuctionBatch{
		ID:                    batchID,
		Code:                  s.generateCode(),
		SlugID:                auctionBatchSlug(req.NamaID, batchID),
		SlugEN:                auctionBatchSlug(auctionBatchSlugName(req.NamaID, req.NamaEN), batchID),
		NamaID:                req.NamaID,
		NamaEN:                req.NamaEN,
		Description:           req.Description,
		WarehouseID:           warehouseID,
		OriginType:            bOriginType,
		SupplierName:          suppName,
		SupplierAddress:       suppAddress,
		SupplierProvinsi:      suppProvinsi,
		SupplierKota:          suppKota,
		SupplierKecamatan:     suppKecamatan,
		SupplierKelurahan:     suppKelurahan,
		SupplierKodePos:       suppKodePos,
		SupplierLatitude:      suppLat,
		SupplierLongitude:     suppLng,
		KategoriID:            parseUUIDPtr(req.KategoriID),
		KondisiID:             parseUUIDPtr(req.KondisiID),
		KondisiPaketID:        parseUUIDPtr(req.KondisiPaketID),
		SumberID:              parseUUIDPtr(req.SumberID),
		DiscrepancyPercentage: parseDecimal(req.DiscrepancyPercentage, decimal.Zero),
		Status:                models.AuctionBatchStatusDRAFT,
		GrandTotal:            grandTotal,
		MinBidPercent:         decimal.NewFromFloat(0.1),
		TotalQuantity:         totalQty,
		PanjangCm:             parseDecimal(req.PanjangCm, decimal.Zero),
		LebarCm:               parseDecimal(req.LebarCm, decimal.Zero),
		TinggiCm:              parseDecimal(req.TinggiCm, decimal.Zero),
		BeratKg:               parseDecimal(req.BeratKg, decimal.Zero),
		VolumeM3:              computeVolume(parseDecimal(req.PanjangCm, decimal.Zero), parseDecimal(req.LebarCm, decimal.Zero), parseDecimal(req.TinggiCm, decimal.Zero)),
		Version:               1,
		CreatedBy:             adminID,
		UpdatedBy:             adminID,
	}

	batch.Items = items
	brandIDs := parseUUIDList(req.MerekIDs)
	assetIDs := parseAssetIDs(req.ImageAssetIDs, req.PDFAssetID)

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(batch).Error; err != nil {
			return err
		}
		if err := s.replaceItems(tx, batch.ID, items); err != nil {
			return err
		}
		if err := s.replaceBrands(tx, batch.ID, brandIDs); err != nil {
			return err
		}
		if err := s.replaceBatchAssets(tx, batch.ID, assetIDs); err != nil {
			return err
		}
		// Audit
		if err := tx.Create(&models.AuctionAuditLog{
			BatchID:      batch.ID,
			ActorAdminID: adminID,
			Action:       "create",
			AfterData:    toJSONMap(batch),
			Note:         &req.NamaID,
		}).Error; err != nil {
			return err
		}
		// Idempotency
		if idempotencyKey != "" {
			if err := s.saveIdempotency(tx, "ADMIN", adminID, operation, uuid.Nil, idempotencyKey, req, 201, batch); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, s.wrapDBError(err)
	}

	return s.GetBatchDetail(ctx, batch.ID)
}

// ============================================================
// Update Draft
// ============================================================

func (s *auctionService) UpdateDraft(ctx context.Context, id uuid.UUID, req *dto.AuctionDraftInput, adminID uuid.UUID, version int) (*dto.AuctionBatchDetail, error) {
	batch, err := s.repo.FindBatchByIDWithRelations(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auctionErr(404, "Batch tidak ditemukan")
		}
		return nil, err
	}
	if batch.Status != models.AuctionBatchStatusDRAFT {
		return nil, auctionErr(409, "Batch hanya dapat diedit saat status DRAFT")
	}
	if batch.Version != version {
		return nil, auctionErr(409, "Version batch tidak sesuai. Silakan muat ulang data terbaru.")
	}

	items, grandTotal, totalQty, err := s.buildItemsSnapshot(ctx, req.Items)
	if err != nil {
		return nil, err
	}

	batch.NamaID = req.NamaID
	batch.NamaEN = req.NamaEN
	batch.SlugID = auctionBatchSlug(req.NamaID, batch.ID)
	batch.SlugEN = auctionBatchSlug(auctionBatchSlugName(req.NamaID, req.NamaEN), batch.ID)
	batch.Description = req.Description
	bOriginType := originType(req.OriginType)
	batch.OriginType = bOriginType
	if bOriginType == "SUPPLIER" {
		batch.WarehouseID = nil
		batch.SupplierName = req.SupplierName
		batch.SupplierAddress = req.SupplierAddress
		batch.SupplierProvinsi = req.SupplierProvinsi
		batch.SupplierKota = req.SupplierKota
		batch.SupplierKecamatan = req.SupplierKecamatan
		batch.SupplierKelurahan = req.SupplierKelurahan
		batch.SupplierKodePos = req.SupplierKodePos
		batch.SupplierLatitude = parseDecimalPtr(req.SupplierLatitude)
		batch.SupplierLongitude = parseDecimalPtr(req.SupplierLongitude)
	} else {
		batch.WarehouseID = parseUUIDPtr(req.WarehouseID)
		batch.SupplierName = nil
		batch.SupplierAddress = nil
		batch.SupplierProvinsi = nil
		batch.SupplierKota = nil
		batch.SupplierKecamatan = nil
		batch.SupplierKelurahan = nil
		batch.SupplierKodePos = nil
		batch.SupplierLatitude = nil
		batch.SupplierLongitude = nil
	}
	batch.KategoriID = parseUUIDPtr(req.KategoriID)
	batch.KondisiID = parseUUIDPtr(req.KondisiID)
	batch.KondisiPaketID = parseUUIDPtr(req.KondisiPaketID)
	batch.SumberID = parseUUIDPtr(req.SumberID)
	batch.DiscrepancyPercentage = parseDecimal(req.DiscrepancyPercentage, decimal.Zero)
	batch.GrandTotal = grandTotal
	batch.TotalQuantity = totalQty
	batch.PanjangCm = parseDecimal(req.PanjangCm, decimal.Zero)
	batch.LebarCm = parseDecimal(req.LebarCm, decimal.Zero)
	batch.TinggiCm = parseDecimal(req.TinggiCm, decimal.Zero)
	batch.BeratKg = parseDecimal(req.BeratKg, decimal.Zero)
	batch.VolumeM3 = computeVolume(batch.PanjangCm, batch.LebarCm, batch.TinggiCm)
	batch.UpdatedBy = adminID
	batch.Version++

	before := toJSONMap(batch)

	brandIDs := parseUUIDList(req.MerekIDs)
	assetIDs := parseAssetIDs(req.ImageAssetIDs, req.PDFAssetID)

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(batch).Error; err != nil {
			return err
		}
		if err := s.replaceItems(tx, id, items); err != nil {
			return err
		}
		if err := s.replaceBrands(tx, id, brandIDs); err != nil {
			return err
		}
		if err := s.replaceBatchAssets(tx, id, assetIDs); err != nil {
			return err
		}
		if err := tx.Create(&models.AuctionAuditLog{
			BatchID:      id,
			ActorAdminID: adminID,
			Action:       "update",
			BeforeData:   before,
			AfterData:    toJSONMap(batch),
			Note:         &req.NamaID,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, s.wrapDBError(err)
	}

	return s.GetBatchDetail(ctx, id)
}

// ============================================================
// Publish
// ============================================================

func (s *auctionService) Publish(ctx context.Context, id uuid.UUID, version int, adminID uuid.UUID, idempotencyKey string) (*dto.AuctionBatchDetail, error) {
	operation := fmt.Sprintf("publish_batch:%s", id.String())
	if idempotencyKey != "" {
		replay, bid, err := s.checkIdempotency(ctx, "ADMIN", adminID, operation, id, idempotencyKey, version)
		if err != nil {
			return nil, err
		}
		if replay {
			return s.GetBatchDetail(ctx, bid)
		}
	}

	detail, err := s.publishInTransaction(ctx, id, version, adminID, idempotencyKey, operation)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

func (s *auctionService) publishInTransaction(ctx context.Context, id uuid.UUID, version int, adminID uuid.UUID, idempotencyKey, operation string) (*dto.AuctionBatchDetail, error) {
	batch, err := s.repo.FindBatchByIDWithRelations(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auctionErr(404, "Batch tidak ditemukan")
		}
		return nil, err
	}
	if batch.Status != models.AuctionBatchStatusDRAFT {
		return nil, auctionErr(409, "Batch hanya dapat dipublish saat status DRAFT")
	}
	if batch.Version != version {
		return nil, auctionErr(409, "Version batch tidak sesuai. Silakan muat ulang data terbaru.")
	}

	// Validasi kelayakan publish.
	if err := s.validatePublish(ctx, batch); err != nil {
		return nil, err
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock batch row.
		var locked models.AuctionBatch
		if err := tx.Clauses(lockClause()).First(&locked, "id = ?", id).Error; err != nil {
			return err
		}
		if locked.Status != models.AuctionBatchStatusDRAFT {
			return auctionErr(409, "Batch sudah tidak berstatus DRAFT")
		}

		// Reservasi stok atomik. Items diambil dari batch yang sudah di-preload
		// (locked tidak memuat Items, karena hanya mengunci row batch).
		if err := s.reserveStock(tx, id, batch.Items); err != nil {
			return err
		}

		now := time.Now().UTC()
		locked.Status = models.AuctionBatchStatusOPEN
		locked.OpenedAt = &now
		locked.UpdatedBy = adminID
		locked.Version++
		if err := tx.Save(&locked).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.AuctionAuditLog{
			BatchID:      id,
			ActorAdminID: adminID,
			Action:       "publish",
			AfterData:    toJSONMap(&locked),
			Note:         ptrString("Batch dibuka untuk lelang"),
		}).Error; err != nil {
			return err
		}

		if idempotencyKey != "" {
			if err := s.saveIdempotency(tx, "ADMIN", adminID, operation, id, idempotencyKey, version, 200, &locked); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, s.wrapDBError(err)
	}

	return s.GetBatchDetail(ctx, id)
}

// validatePublish memastikan batch lengkap sebelum dibuka.
func (s *auctionService) validatePublish(ctx context.Context, batch *models.AuctionBatch) error {
	var fieldErrs []models.FieldError

	if len(batch.Items) == 0 {
		fieldErrs = append(fieldErrs, models.FieldError{Field: "items", Message: "Minimal satu produk diperlukan untuk membuka lelang"})
	}
	if len(batch.BatchAssets) == 0 {
		fieldErrs = append(fieldErrs, models.FieldError{Field: "image_asset_ids", Message: "Minimal satu gambar diperlukan untuk membuka lelang"})
	}
	if batch.KategoriID == nil || batch.KondisiID == nil || batch.KondisiPaketID == nil {
		fieldErrs = append(fieldErrs, models.FieldError{Field: "master_data", Message: "Kategori, kondisi produk dan kondisi paket wajib diisi saat publish"})
	}
	if batch.OriginType == "BULKY_WAREHOUSE" && batch.WarehouseID == nil {
		fieldErrs = append(fieldErrs, models.FieldError{Field: "warehouse_id", Message: "Warehouse Bulky wajib dipilih"})
	}
	if batch.OriginType == "SUPPLIER" {
		if isBlank(batch.SupplierName) || isBlank(batch.SupplierAddress) || isBlank(batch.SupplierKota) || batch.SupplierLatitude == nil || batch.SupplierLongitude == nil {
			fieldErrs = append(fieldErrs, models.FieldError{Field: "supplier_origin", Message: "Nama gudang, alamat, kota, latitude, dan longitude titik asal pengiriman wajib diisi saat publish"})
		}
		if batch.SupplierLatitude != nil && (batch.SupplierLatitude.LessThan(decimal.NewFromInt(-90)) || batch.SupplierLatitude.GreaterThan(decimal.NewFromInt(90))) {
			fieldErrs = append(fieldErrs, models.FieldError{Field: "supplier_latitude", Message: "Latitude gudang supplier harus antara -90 dan 90"})
		}
		if batch.SupplierLongitude != nil && (batch.SupplierLongitude.LessThan(decimal.NewFromInt(-180)) || batch.SupplierLongitude.GreaterThan(decimal.NewFromInt(180))) {
			fieldErrs = append(fieldErrs, models.FieldError{Field: "supplier_longitude", Message: "Longitude gudang supplier harus antara -180 dan 180"})
		}
	}
	if batch.OriginType != "BULKY_WAREHOUSE" && batch.OriginType != "SUPPLIER" {
		fieldErrs = append(fieldErrs, models.FieldError{Field: "origin_type", Message: "Asal pengiriman tidak valid"})
	}
	if !batch.GrandTotal.IsPositive() {
		fieldErrs = append(fieldErrs, models.FieldError{Field: "grand_total", Message: "Grand total harus lebih dari nol"})
	}
	if !batch.PanjangCm.IsPositive() || !batch.LebarCm.IsPositive() || !batch.TinggiCm.IsPositive() || !batch.BeratKg.IsPositive() {
		fieldErrs = append(fieldErrs, models.FieldError{Field: "dimensi", Message: "Ukuran dan berat batch harus lebih dari nol saat publish"})
	}

	if len(fieldErrs) > 0 {
		return auctionErrFields(400, "Data batch belum lengkap", fieldErrs)
	}

	// Hanya item katalog yang terhubung ke stok Bulky. Snapshot MANUAL tidak
	// boleh dipaksa masuk ke tabel produk atau diresevasi sebagai stok Bulky.
	produkIDs := make([]uuid.UUID, 0, len(batch.Items))
	for _, item := range batch.Items {
		if item.SourceType == "CATALOG" && item.ProdukID != nil {
			produkIDs = append(produkIDs, *item.ProdukID)
		}
	}

	var produkList []models.Produk
	if err := s.db.WithContext(ctx).Where("id IN ?", produkIDs).Find(&produkList).Error; err != nil {
		return err
	}
	produkByID := make(map[uuid.UUID]models.Produk)
	for _, p := range produkList {
		produkByID[p.ID] = p
	}

	for _, item := range batch.Items {
		if batch.OriginType == "SUPPLIER" && item.SourceType == "CATALOG" {
			fieldErrs = append(fieldErrs, models.FieldError{Field: "items", Message: "Batch dari gudang supplier hanya dapat memakai item manual"})
			continue
		}
		if item.SourceType == "MANUAL" {
			if !item.UnitPriceSnapshot.IsPositive() || !isWholeNumber(item.UnitPriceSnapshot) {
				fieldErrs = append(fieldErrs, models.FieldError{Field: "items", Message: "Harga item manual harus rupiah bulat dan lebih dari nol"})
			}
			continue
		}
		if item.ProdukID == nil {
			fieldErrs = append(fieldErrs, models.FieldError{Field: "items", Message: "Produk katalog wajib dipilih"})
			continue
		}
		p, ok := produkByID[*item.ProdukID]
		if !ok {
			fieldErrs = append(fieldErrs, models.FieldError{Field: "items", Message: "Produk tidak ditemukan"})
			continue
		}
		if !p.IsActive {
			fieldErrs = append(fieldErrs, models.FieldError{Field: "items", Message: "Produk " + p.NamaID + " tidak aktif"})
		}
		if p.IsSold {
			fieldErrs = append(fieldErrs, models.FieldError{Field: "items", Message: "Produk " + p.NamaID + " sudah terjual"})
		}
		if batch.OriginType == "BULKY_WAREHOUSE" && batch.WarehouseID != nil && p.WarehouseID != *batch.WarehouseID {
			fieldErrs = append(fieldErrs, models.FieldError{Field: "items", Message: "Produk " + p.NamaID + " berada di warehouse berbeda"})
		}
		// Snapshot harga harus rupiah bulat (BR03).
		if !isWholeNumber(decimal.NewFromFloat(p.HargaSesudahDiskon)) {
			fieldErrs = append(fieldErrs, models.FieldError{Field: "items", Message: "Harga produk " + p.NamaID + " memiliki pecahan dan belum dapat dipublish"})
		}
		// Stok tersedia harus cukup (kurangi reservasi aktif).
		available := s.stockAvailable(ctx, p)
		if available < item.Quantity {
			fieldErrs = append(fieldErrs, models.FieldError{Field: "items", Message: "Stok produk " + p.NamaID + " tidak mencukupi"})
		}
	}

	if len(fieldErrs) > 0 {
		return auctionErrFields(400, "Data batch belum lengkap", fieldErrs)
	}
	return nil
}

// reserveStock membuat reservasi stok untuk setiap item batch. Row produk
// dikunci dalam urutan ID yang konsisten agar dua publish bersamaan tidak
// saling deadlock dan tidak oversell (BR08).
func (s *auctionService) reserveStock(tx *gorm.DB, batchID uuid.UUID, items []models.AuctionBatchItem) error {
	produkIDs := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		if item.SourceType == "CATALOG" && item.ProdukID != nil {
			produkIDs = append(produkIDs, *item.ProdukID)
		}
	}
	// Urutan lock produk konsisten (sorted) untuk mencegah deadlock.
	sort.Slice(produkIDs, func(i, j int) bool { return produkIDs[i].String() < produkIDs[j].String() })

	// Ambil reservasi aktif existing untuk produk tersebut.
	var existingReservations []models.AuctionStockReservation
	if len(produkIDs) > 0 {
		if err := tx.Where("produk_id IN ? AND status = ?", produkIDs, "ACTIVE").Find(&existingReservations).Error; err != nil {
			return err
		}
	}
	reservedByProduk := make(map[uuid.UUID]int)
	for _, r := range existingReservations {
		reservedByProduk[r.ProdukID] += r.Quantity
	}

	// Lock row produk satu per satu dalam urutan konsisten.
	produkByID := make(map[uuid.UUID]models.Produk)
	for _, pid := range produkIDs {
		var p models.Produk
		if err := tx.Clauses(lockClause()).First(&p, "id = ?", pid).Error; err != nil {
			return err
		}
		produkByID[pid] = p
	}

	// Validasi & insert reservasi (all-or-nothing).
	for _, item := range items {
		if item.SourceType != "CATALOG" || item.ProdukID == nil {
			continue
		}
		p, ok := produkByID[*item.ProdukID]
		if !ok {
			return auctionErr(409, "Produk tidak ditemukan saat reservasi stok")
		}
		available := p.Quantity - reservedByProduk[*item.ProdukID]
		if available < item.Quantity {
			return auctionErr(409, "Stok produk "+p.NamaID+" tidak mencukupi untuk reservasi")
		}
		res := &models.AuctionStockReservation{
			BatchID:  batchID,
			ProdukID: *item.ProdukID,
			Quantity: item.Quantity,
			Status:   "ACTIVE",
		}
		if err := tx.Create(res).Error; err != nil {
			return err
		}
		reservedByProduk[*item.ProdukID] += item.Quantity
	}
	return nil
}

// ============================================================
// List Batches
// ============================================================

func (s *auctionService) ListBatches(ctx context.Context, params *dto.AuctionListQueryParams) ([]dto.AuctionBatchSummary, *models.PaginationMeta, error) {
	batches, total, err := s.repo.ListBatches(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	// Agregasi massal untuk hindari N+1.
	batchIDs := make([]uuid.UUID, 0, len(batches))
	for _, b := range batches {
		batchIDs = append(batchIDs, b.ID)
	}
	analytics, _ := s.repo.GetBatchesAnalytics(ctx, batchIDs)
	thumbnails, _ := s.repo.GetFirstImagesForBatches(ctx, batchIDs)

	response := make([]dto.AuctionBatchSummary, 0, len(batches))
	for _, b := range batches {
		response = append(response, s.mapBatchSummary(b, analytics[b.ID], thumbnails[b.ID]))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return response, &meta, nil
}

// ============================================================
// Get Batch Detail
// ============================================================

func (s *auctionService) GetBatchDetail(ctx context.Context, id uuid.UUID) (*dto.AuctionBatchDetail, error) {
	batch, err := s.repo.FindBatchByIDWithRelations(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auctionErr(404, "Batch tidak ditemukan")
		}
		return nil, err
	}

	analytics, err := s.repo.GetBatchAnalytics(ctx, id)
	if err != nil {
		return nil, err
	}
	viewCount, err := s.repo.GetBatchViewCount(ctx, id)
	if err != nil {
		return nil, err
	}
	assets, err := s.repo.GetBatchAssetDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	// Winner buyer info.
	var winner *dto.AuctionWinnerResponse
	if batch.Winner != nil {
		bid, bidErr := s.repo.FindBidByID(ctx, batch.Winner.BidID)
		if bidErr != nil {
			return nil, bidErr
		}
		winner = mapWinner(batch.Winner, bid)
	}

	return s.mapBatchDetail(batch, analytics, assets, viewCount, winner), nil
}

// ============================================================
// List Bids
// ============================================================

func (s *auctionService) ListBids(ctx context.Context, batchID uuid.UUID, params *dto.AuctionBidsQueryParams) ([]dto.AuctionBidDetail, *models.PaginationMeta, error) {
	if _, err := s.repo.FindBatchByID(ctx, batchID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, auctionErr(404, "Batch tidak ditemukan")
		}
		return nil, nil, err
	}

	bids, total, err := s.repo.ListBids(ctx, batchID, params)
	if err != nil {
		return nil, nil, err
	}

	// Tandai bid yang terpilih menjadi winner.
	selectedBidID := uuid.Nil
	if winner, err := s.repo.FindWinnerByBatchID(ctx, batchID); err == nil {
		selectedBidID = winner.BidID
	}

	response := make([]dto.AuctionBidDetail, 0, len(bids))
	for _, bid := range bids {
		response = append(response, mapBidDetail(bid, bid.ID == selectedBidID))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return response, &meta, nil
}

// ============================================================
// Select Winner
// ============================================================

func (s *auctionService) SelectWinner(ctx context.Context, id uuid.UUID, req *dto.AuctionWinnerRequest, adminID uuid.UUID, idempotencyKey string) (*dto.AuctionBatchDetail, error) {
	operation := fmt.Sprintf("winner:%s", id.String())
	if idempotencyKey != "" {
		replay, bid, err := s.checkIdempotency(ctx, "ADMIN", adminID, operation, id, idempotencyKey, req)
		if err != nil {
			return nil, err
		}
		if replay {
			return s.GetBatchDetail(ctx, bid)
		}
	}

	batch, err := s.repo.FindBatchByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auctionErr(404, "Batch tidak ditemukan")
		}
		return nil, err
	}
	if batch.Status != models.AuctionBatchStatusOPEN {
		return nil, auctionErr(409, "Winner hanya dapat dipilih saat batch berstatus OPEN")
	}
	if batch.Version != req.Version {
		return nil, auctionErr(409, "Version batch tidak sesuai. Silakan muat ulang data terbaru.")
	}

	bidID, _ := uuid.Parse(req.BidID)
	bid, err := s.repo.FindBidByBatchAndID(ctx, id, bidID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auctionErr(404, "Bid tidak ditemukan pada batch ini")
		}
		return nil, err
	}

	// Pastikan belum ada winner.
	if existing, err := s.repo.FindWinnerByBatchID(ctx, id); err == nil && existing != nil {
		return nil, auctionErr(409, "Batch sudah memiliki pemenang")
	}

	// Lock batch + insert winner + konsumsi stok + SOLD dalam satu transaksi.
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var locked models.AuctionBatch
		if err := tx.Clauses(lockClause()).First(&locked, "id = ?", id).Error; err != nil {
			return err
		}
		if locked.Status != models.AuctionBatchStatusOPEN {
			return auctionErr(409, "Batch sudah tidak berstatus OPEN")
		}
		if _, err := s.repo.FindWinnerByBatchID(ctx, id); err == nil {
			return auctionErr(409, "Batch sudah memiliki pemenang")
		}

		now := time.Now().UTC()
		winner := &models.AuctionWinner{
			BatchID:           id,
			BidID:             bid.ID,
			DealAmount:        bid.Amount,
			SelectedBy:        adminID,
			SelectedAt:        now,
			Note:              req.Note,
			PaymentStatus:     "UNPAID",
			FulfillmentStatus: "PENDING",
		}
		if err := tx.Create(winner).Error; err != nil {
			return err
		}

		// Konsumsi reservasi & kurangi stok produk.
		if err := s.consumeReservations(tx, id); err != nil {
			return err
		}

		locked.Status = models.AuctionBatchStatusSOLD
		locked.SoldAt = &now
		locked.UpdatedBy = adminID
		locked.Version++
		if err := tx.Save(&locked).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.AuctionAuditLog{
			BatchID:      id,
			ActorAdminID: adminID,
			Action:       "winner",
			BeforeData:   toJSONMap(&locked),
			AfterData:    toJSONMap(winner),
			Note:         req.Note,
		}).Error; err != nil {
			return err
		}

		if idempotencyKey != "" {
			if err := s.saveIdempotency(tx, "ADMIN", adminID, operation, id, idempotencyKey, req, 200, &locked); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, s.wrapDBError(err)
	}

	return s.GetBatchDetail(ctx, id)
}

// consumeReservations menandai reservasi CONSUMED dan mengurangi stok produk.
// Row produk dikunci dalam urutan ID konsisten agar tidak deadlock dengan publish.
func (s *auctionService) consumeReservations(tx *gorm.DB, batchID uuid.UUID) error {
	var reservations []models.AuctionStockReservation
	if err := tx.Where("batch_id = ? AND status = ?", batchID, "ACTIVE").Find(&reservations).Error; err != nil {
		return err
	}
	now := time.Now().UTC()
	// Urutan lock produk konsisten.
	sort.Slice(reservations, func(i, j int) bool {
		return reservations[i].ProdukID.String() < reservations[j].ProdukID.String()
	})
	for _, r := range reservations {
		if err := tx.Model(&models.AuctionStockReservation{}).
			Where("id = ?", r.ID).
			Updates(map[string]interface{}{"status": "CONSUMED", "consumed_at": now}).Error; err != nil {
			return err
		}
		// Kurangi stok produk tepat satu kali (lock row produk).
		var p models.Produk
		if err := tx.Clauses(lockClause()).First(&p, "id = ?", r.ProdukID).Error; err != nil {
			return err
		}
		newQty := p.Quantity - r.Quantity
		if newQty < 0 {
			newQty = 0
		}
		updates := map[string]interface{}{"quantity": newQty}
		// is_sold mengikuti stok tersisa, bukan otomatis true saat sebagian terjual.
		if newQty == 0 {
			updates["is_sold"] = true
		}
		if err := tx.Model(&models.Produk{}).Where("id = ?", r.ProdukID).Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}

// ============================================================
// Operations (payment / fulfillment manual)
// ============================================================

func (s *auctionService) UpdateOperations(ctx context.Context, id uuid.UUID, req *dto.AuctionOperationRequest, adminID uuid.UUID, idempotencyKey string) (*dto.AuctionBatchDetail, error) {
	operation := fmt.Sprintf("operations:%s", id.String())
	if idempotencyKey != "" {
		replay, bid, err := s.checkIdempotency(ctx, "ADMIN", adminID, operation, id, idempotencyKey, req)
		if err != nil {
			return nil, err
		}
		if replay {
			return s.GetBatchDetail(ctx, bid)
		}
	}

	batch, err := s.repo.FindBatchByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auctionErr(404, "Batch tidak ditemukan")
		}
		return nil, err
	}
	if batch.Status != models.AuctionBatchStatusSOLD {
		return nil, auctionErr(409, "Operasi manual hanya dapat dilakukan saat batch berstatus SOLD")
	}
	if batch.Version != req.Version {
		return nil, auctionErr(409, "Version batch tidak sesuai. Silakan muat ulang data terbaru.")
	}
	if req.PaymentStatus == nil && req.FulfillmentStatus == nil {
		return nil, auctionErr(400, "Tepat satu target status (payment atau fulfillment) harus dikirim")
	}
	if req.PaymentStatus != nil && req.FulfillmentStatus != nil {
		return nil, auctionErr(400, "Hanya satu target status per request yang diperbolehkan")
	}

	winner, err := s.repo.FindWinnerByBatchID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, auctionErr(404, "Winner tidak ditemukan")
		}
		return nil, err
	}

	// Tentukan transisi.
	var changed bool
	if req.PaymentStatus != nil {
		target := *req.PaymentStatus
		if winner.PaymentStatus != target {
			if winner.PaymentStatus == "UNPAID" && target == "PAID" {
				changed = true
			} else {
				return nil, auctionErr(409, "Transisi payment status tidak valid")
			}
		}
	} else {
		target := *req.FulfillmentStatus
		if winner.FulfillmentStatus != target {
			switch {
			case winner.FulfillmentStatus == "PENDING" && target == "PROCESSING":
				if winner.PaymentStatus != "PAID" {
					return nil, auctionErr(409, "Pembayaran harus PAID sebelum fulfillment diproses")
				}
				changed = true
			case winner.FulfillmentStatus == "PROCESSING" && target == "COMPLETED":
				changed = true
			default:
				return nil, auctionErr(409, "Transisi fulfillment status tidak valid")
			}
		}
	}

	// No-op: status sama, tidak ada perubahan nyata.
	if !changed {
		return s.GetBatchDetail(ctx, id)
	}

	// Perubahan nyata wajib memiliki catatan.
	if req.Note == nil || strings.TrimSpace(*req.Note) == "" {
		return nil, auctionErr(400, "Catatan wajib untuk perubahan status")
	}

	now := time.Now().UTC()
	updates := map[string]interface{}{}
	if req.PaymentStatus != nil {
		updates["payment_status"] = "PAID"
		updates["paid_at"] = now
		updates["payment_note"] = req.Note
	} else {
		if *req.FulfillmentStatus == "COMPLETED" {
			updates["fulfillment_status"] = "COMPLETED"
			updates["completed_at"] = now
		} else {
			updates["fulfillment_status"] = "PROCESSING"
		}
		updates["fulfillment_note"] = req.Note
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.AuctionWinner{}).Where("id = ?", winner.ID).Updates(updates).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.AuctionBatch{}).Where("id = ?", id).
			Updates(map[string]interface{}{"version": batch.Version + 1, "updated_by": adminID}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.AuctionAuditLog{
			BatchID:      id,
			ActorAdminID: adminID,
			Action:       "operations",
			AfterData:    updates,
			Note:         req.Note,
		}).Error; err != nil {
			return err
		}
		if idempotencyKey != "" {
			if err := s.saveIdempotency(tx, "ADMIN", adminID, operation, id, idempotencyKey, req, 200, batch); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, s.wrapDBError(err)
	}

	return s.GetBatchDetail(ctx, id)
}

// ============================================================
// Upload Asset
// ============================================================

func (s *auctionService) UploadAsset(ctx context.Context, file *multipart.FileHeader, kind string, adminID uuid.UUID) (*dto.AuctionAssetResponse, error) {
	if kind == "" {
		kind = "IMAGE"
	}

	var validType bool
	switch kind {
	case "IMAGE":
		validType = utils.IsValidImageType(file)
		if !validType {
			return nil, auctionErr(415, "Format gambar tidak didukung. Gunakan jpg, png, atau webp")
		}
		if file.Size > utils.MaxImageSize {
			return nil, auctionErr(413, "Ukuran gambar maksimal 5MB")
		}
	case "PDF":
		validType = utils.IsValidDocumentType(file)
		if !validType {
			return nil, auctionErr(415, "Format dokumen tidak didukung. Gunakan PDF")
		}
		if file.Size > 10*1024*1024 {
			return nil, auctionErr(413, "Ukuran PDF maksimal 10MB")
		}
	default:
		return nil, auctionErr(400, "Kind harus IMAGE atau PDF")
	}

	// Gambar batch mengikuti pipeline gambar produk: dikompresi dan selalu
	// disimpan sebagai WebP. PDF tetap disimpan dalam format aslinya.
	directory := "auction"
	storageKey := ""
	var err error
	if kind == "IMAGE" {
		storageKey, err = utils.CompressAndSaveImageWebP(file, directory, s.cfg)
	} else {
		storageKey, err = utils.SaveUploadedFile(file, directory, s.cfg)
	}
	if err != nil {
		return nil, auctionErr(400, err.Error())
	}

	sizeBytes := file.Size
	mimeType := file.Header.Get("Content-Type")
	if kind == "IMAGE" {
		mimeType = "image/webp"
		if info, statErr := os.Stat(filepath.Join(s.cfg.UploadPath, filepath.FromSlash(storageKey))); statErr == nil {
			sizeBytes = info.Size()
		}
	}

	asset := &models.AuctionAsset{
		UploadedBy:   adminID,
		Kind:         kind,
		StorageKey:   storageKey,
		OriginalName: file.Filename,
		MimeType:     mimeType,
		SizeBytes:    sizeBytes,
	}
	if err := s.repo.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}

	return mapAsset(asset, s.cfg), nil
}

// ============================================================
// Product Options
// ============================================================

func (s *auctionService) ListProductOptions(ctx context.Context, params *dto.AuctionProductOptionsQueryParams) ([]dto.AuctionProductOption, *models.PaginationMeta, error) {
	produk, total, err := s.repo.ListProductOptions(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	produkIDs := make([]uuid.UUID, 0, len(produk))
	for _, p := range produk {
		produkIDs = append(produkIDs, p.ID)
	}
	reservations, _ := s.repo.GetActiveReservationsForProducts(ctx, produkIDs)
	reservedByProduk := make(map[uuid.UUID]int)
	for _, r := range reservations {
		reservedByProduk[r.ProdukID] += r.Quantity
	}

	response := make([]dto.AuctionProductOption, 0, len(produk))
	for _, p := range produk {
		available := p.Quantity - reservedByProduk[p.ID]
		if available < 0 {
			available = 0
		}
		response = append(response, dto.AuctionProductOption{
			ID:             p.ID.String(),
			NamaID:         p.NamaID,
			WarehouseID:    p.WarehouseID.String(),
			UnitPrice:      decimal.NewFromFloat(p.HargaSesudahDiskon).StringFixed(0),
			StockAvailable: available,
			IsActive:       p.IsActive,
			IsSold:         p.IsSold,
		})
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return response, &meta, nil
}

// ============================================================
// Helper: snapshot & turunan
// ============================================================

func (s *auctionService) buildItemsSnapshot(ctx context.Context, reqItems []dto.AuctionDraftItem) ([]models.AuctionBatchItem, decimal.Decimal, int, error) {
	items := make([]models.AuctionBatchItem, 0, len(reqItems))
	var grandTotal decimal.Decimal = decimal.Zero
	totalQty := 0
	seen := make(map[uuid.UUID]bool)

	if len(reqItems) == 0 {
		return items, grandTotal, totalQty, nil
	}

	produkIDs := make([]uuid.UUID, 0, len(reqItems))
	for _, it := range reqItems {
		sourceType := strings.ToUpper(it.SourceType)
		if sourceType == "" {
			sourceType = "CATALOG"
		}
		if sourceType == "CATALOG" {
			pid, err := uuid.Parse(it.ProdukID)
			if err != nil {
				return nil, decimal.Zero, 0, auctionErr(400, "Produk katalog tidak valid")
			}
			produkIDs = append(produkIDs, pid)
		}
	}

	var produkList []models.Produk
	if len(produkIDs) > 0 {
		if err := s.db.WithContext(ctx).Where("id IN ?", produkIDs).Find(&produkList).Error; err != nil {
			return nil, decimal.Zero, 0, err
		}
	}
	produkByID := make(map[uuid.UUID]models.Produk)
	for _, p := range produkList {
		produkByID[p.ID] = p
	}

	for _, it := range reqItems {
		sourceType := strings.ToUpper(it.SourceType)
		if sourceType == "" {
			sourceType = "CATALOG"
		}
		if sourceType != "CATALOG" && sourceType != "MANUAL" {
			return nil, decimal.Zero, 0, auctionErr(400, "Sumber item tidak valid")
		}
		if it.Quantity < 1 {
			return nil, decimal.Zero, 0, auctionErr(400, "Jumlah item minimal satu")
		}
		if sourceType == "MANUAL" {
			name := strings.TrimSpace(it.Nama)
			price := parseDecimal(it.UnitPrice, decimal.Zero)
			if name == "" || !price.IsPositive() || !isWholeNumber(price) {
				return nil, decimal.Zero, 0, auctionErr(400, "Item manual memerlukan nama dan harga rupiah bulat lebih dari nol")
			}
			subtotal := price.Mul(decimal.NewFromInt(int64(it.Quantity)))
			grandTotal = grandTotal.Add(subtotal)
			totalQty += it.Quantity
			items = append(items, models.AuctionBatchItem{SourceType: "MANUAL", Quantity: it.Quantity, NamaSnapshot: name, UnitPriceSnapshot: price, SubtotalSnapshot: subtotal})
			continue
		}

		pid, _ := uuid.Parse(it.ProdukID)
		if seen[pid] {
			return nil, decimal.Zero, 0, auctionErr(400, "Produk duplikat pada item batch")
		}
		seen[pid] = true

		p, ok := produkByID[pid]
		if !ok {
			return nil, decimal.Zero, 0, auctionErr(404, "Produk tidak ditemukan")
		}

		// Snapshot harga: rupiah bulat (pembulatan ke atas per D02).
		unitPrice := decimal.NewFromFloat(p.HargaSesudahDiskon).Ceil()
		subtotal := unitPrice.Mul(decimal.NewFromInt(int64(it.Quantity)))
		grandTotal = grandTotal.Add(subtotal)
		totalQty += it.Quantity

		items = append(items, models.AuctionBatchItem{
			ProdukID:          &pid,
			SourceType:        "CATALOG",
			Quantity:          it.Quantity,
			NamaSnapshot:      p.NamaID,
			UnitPriceSnapshot: unitPrice,
			SubtotalSnapshot:  subtotal,
		})
	}

	return items, grandTotal, totalQty, nil
}

func computeVolume(p, l, t decimal.Decimal) decimal.Decimal {
	if !p.IsPositive() || !l.IsPositive() || !t.IsPositive() {
		return decimal.Zero
	}
	vol := p.Mul(l).Mul(t).Div(decimal.NewFromInt(1000000))
	return vol.Round(3)
}

func minBidAmount(grandTotal decimal.Decimal) decimal.Decimal {
	return grandTotal.Mul(decimal.NewFromFloat(0.001)).Ceil()
}

// ============================================================
// Helper: idempotency
// ============================================================

func computeRequestHash(operation string, targetID uuid.UUID, payload interface{}) string {
	b, _ := json.Marshal(payload)
	h := sha256.Sum256([]byte(operation + "|" + targetID.String() + "|" + string(b)))
	return hex.EncodeToString(h[:])
}

// checkIdempotency memeriksa apakah key sudah pernah digunakan. Mengembalikan
// replay=true bila key+payload sama (retry), dan batchID dari hasil tersimpan
// agar service dapat mengembalikan detail terkini.
func (s *auctionService) checkIdempotency(ctx context.Context, actorType string, actorID uuid.UUID, operation string, targetID uuid.UUID, key string, payload interface{}) (replay bool, batchID uuid.UUID, err error) {
	rec, err := s.repo.FindIdempotency(ctx, actorType, actorID, operation, key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, uuid.Nil, nil // tidak ada, lanjut proses
		}
		return false, uuid.Nil, err
	}
	hash := computeRequestHash(operation, targetID, payload)
	if rec.RequestHash != hash {
		return false, uuid.Nil, auctionErr(409, "Idempotency key sudah dipakai dengan payload berbeda")
	}
	bid := targetID
	if bid == uuid.Nil {
		if idStr, ok := rec.ResponseBody["id"].(string); ok {
			bid, _ = uuid.Parse(idStr)
		}
	}
	return true, bid, nil
}

func (s *auctionService) saveIdempotency(tx *gorm.DB, actorType string, actorID uuid.UUID, operation string, targetID uuid.UUID, key string, payload interface{}, status int, result interface{}) error {
	hash := computeRequestHash(operation, targetID, payload)
	rec := &models.AuctionIdempotency{
		ActorID:        actorID,
		ActorType:      actorType,
		Operation:      operation,
		Key:            key,
		RequestHash:    hash,
		ResponseStatus: status,
		ResponseBody:   toJSONMap(result),
	}
	return tx.Create(rec).Error
}

// ============================================================
// Helper: replace pivot (dalam transaksi)
// ============================================================

func (s *auctionService) replaceItems(tx *gorm.DB, batchID uuid.UUID, items []models.AuctionBatchItem) error {
	if err := tx.Where("batch_id = ?", batchID).Delete(&models.AuctionBatchItem{}).Error; err != nil {
		return err
	}
	for i := range items {
		items[i].BatchID = batchID
	}
	if len(items) == 0 {
		return nil
	}
	return tx.Create(&items).Error
}

func (s *auctionService) replaceBrands(tx *gorm.DB, batchID uuid.UUID, merekIDs []uuid.UUID) error {
	if err := tx.Where("batch_id = ?", batchID).Delete(&models.AuctionBatchBrand{}).Error; err != nil {
		return err
	}
	if len(merekIDs) == 0 {
		return nil
	}
	brands := make([]models.AuctionBatchBrand, 0, len(merekIDs))
	for _, m := range merekIDs {
		brands = append(brands, models.AuctionBatchBrand{BatchID: batchID, MerekID: m})
	}
	return tx.Create(&brands).Error
}

func (s *auctionService) replaceBatchAssets(tx *gorm.DB, batchID uuid.UUID, assetIDs []uuid.UUID) error {
	if err := tx.Where("batch_id = ?", batchID).Delete(&models.AuctionBatchAsset{}).Error; err != nil {
		return err
	}
	if len(assetIDs) == 0 {
		return nil
	}
	assets := make([]models.AuctionBatchAsset, 0, len(assetIDs))
	for i, a := range assetIDs {
		assets = append(assets, models.AuctionBatchAsset{BatchID: batchID, AssetID: a, SortOrder: i})
	}
	return tx.Create(&assets).Error
}

// ============================================================
// Helper: mapping response
// ============================================================

func (s *auctionService) mapBatchSummary(b models.AuctionBatch, analytics *repositories.AuctionAnalytics, thumbnail *models.AuctionAsset) dto.AuctionBatchSummary {
	summary := dto.AuctionBatchSummary{
		ID:           b.ID.String(),
		Code:         b.Code,
		SlugID:       b.SlugID,
		SlugEN:       b.SlugEN,
		NamaID:       b.NamaID,
		Status:       b.Status,
		GrandTotal:   b.GrandTotal.StringFixed(0),
		MinBidAmount: minBidAmount(b.GrandTotal).StringFixed(0),
		CreatedAt:    b.CreatedAt,
		OpenedAt:     b.OpenedAt,
		SoldAt:       b.SoldAt,
		Version:      b.Version,
	}
	if thumbnail != nil {
		url := utils.GetFileURL(thumbnail.StorageKey, s.cfg)
		summary.ThumbnailURL = &url
	}
	if analytics != nil {
		summary.BidderCount = analytics.BidderCount
		summary.BidCount = analytics.BidCount
		if analytics.HighestBid != nil {
			v := decimal.NewFromFloat(*analytics.HighestBid).StringFixed(0)
			summary.HighestBid = &v
		}
	}
	return summary
}

func (s *auctionService) mapBatchDetail(b *models.AuctionBatch, analytics *repositories.AuctionAnalytics, assets []models.AuctionAsset, viewCount int64, winner *dto.AuctionWinnerResponse) *dto.AuctionBatchDetail {
	detail := &dto.AuctionBatchDetail{
		ID:                    b.ID.String(),
		Code:                  b.Code,
		SlugID:                b.SlugID,
		SlugEN:                b.SlugEN,
		NamaID:                b.NamaID,
		NamaEN:                b.NamaEN,
		Description:           b.Description,
		WarehouseID:           uuidPtrToString(b.WarehouseID),
		OriginType:            b.OriginType,
		SupplierName:          b.SupplierName,
		SupplierAddress:       b.SupplierAddress,
		SupplierProvinsi:      b.SupplierProvinsi,
		SupplierKota:          b.SupplierKota,
		SupplierKecamatan:     b.SupplierKecamatan,
		SupplierKelurahan:     b.SupplierKelurahan,
		SupplierKodePos:       b.SupplierKodePos,
		SupplierLatitude:      decimalPtrToString(b.SupplierLatitude),
		SupplierLongitude:     decimalPtrToString(b.SupplierLongitude),
		KategoriID:            uuidPtrToString(b.KategoriID),
		KondisiID:             uuidPtrToString(b.KondisiID),
		KondisiPaketID:        uuidPtrToString(b.KondisiPaketID),
		SumberID:              uuidPtrToString(b.SumberID),
		DiscrepancyPercentage: b.DiscrepancyPercentage.String(),
		Status:                b.Status,
		GrandTotal:            b.GrandTotal.StringFixed(0),
		MinBidPercent:         b.MinBidPercent.String(),
		MinBidAmount:          minBidAmount(b.GrandTotal).StringFixed(0),
		TotalQuantity:         b.TotalQuantity,
		PanjangCm:             b.PanjangCm.String(),
		LebarCm:               b.LebarCm.String(),
		TinggiCm:              b.TinggiCm.String(),
		BeratKg:               b.BeratKg.String(),
		VolumeM3:              b.VolumeM3.String(),
		Version:               b.Version,
		CreatedBy:             b.CreatedBy.String(),
		UpdatedBy:             b.UpdatedBy.String(),
		CreatedAt:             b.CreatedAt,
		UpdatedAt:             b.UpdatedAt,
		OpenedAt:              b.OpenedAt,
		SoldAt:                b.SoldAt,
		ViewCount:             viewCount,
		Winner:                winner,
		Items:                 []dto.AuctionItemSnapshot{},
		BrandIDs:              []string{},
		Images:                []dto.AuctionAssetResponse{},
	}

	for _, item := range b.Items {
		detail.Items = append(detail.Items, dto.AuctionItemSnapshot{
			ProdukID:          uuidPtrToString(item.ProdukID),
			SourceType:        item.SourceType,
			NamaSnapshot:      item.NamaSnapshot,
			Quantity:          item.Quantity,
			UnitPriceSnapshot: item.UnitPriceSnapshot.StringFixed(0),
			SubtotalSnapshot:  item.SubtotalSnapshot.StringFixed(0),
		})
	}
	for _, brand := range b.Brands {
		detail.BrandIDs = append(detail.BrandIDs, brand.MerekID.String())
	}
	for _, a := range assets {
		asset := mapAsset(&a, s.cfg)
		if a.Kind == "IMAGE" {
			detail.Images = append(detail.Images, *asset)
		} else if a.Kind == "PDF" {
			detail.PDF = asset
		}
	}

	if analytics != nil {
		detail.BidderCount = analytics.BidderCount
		detail.BidCount = analytics.BidCount
		detail.RepeatBidderCount = analytics.RepeatBidderCount
		detail.AdditionalBidCount = analytics.AdditionalBidCount
		if analytics.HighestBid != nil {
			v := decimal.NewFromFloat(*analytics.HighestBid).StringFixed(0)
			detail.HighestBid = &v
		}
		if analytics.LowestBid != nil {
			v := decimal.NewFromFloat(*analytics.LowestBid).StringFixed(0)
			detail.LowestBid = &v
		}
	}

	if b.SoldAt != nil && b.OpenedAt != nil {
		secs := int64(b.SoldAt.Sub(*b.OpenedAt).Seconds())
		detail.TimeToSoldSeconds = &secs
	}

	return detail
}

// ============================================================
// Helper: util parsing & konversi
// ============================================================

func parseUUIDPtr(s *string) *uuid.UUID {
	if s == nil || *s == "" {
		return nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil
	}
	return &id
}

func originType(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return "BULKY_WAREHOUSE"
	}
	return value
}

func auctionBatchSlugName(nameID string, nameEN *string) string {
	if nameEN != nil && strings.TrimSpace(*nameEN) != "" {
		return *nameEN
	}
	return nameID
}

func auctionBatchSlug(name string, id uuid.UUID) string {
	slug := utils.GenerateSlug(name)
	if slug == "" {
		slug = "auction"
	}
	return slug + "-" + id.String()[:8]
}

func isBlank(value *string) bool {
	return value == nil || strings.TrimSpace(*value) == ""
}

func parseUUIDList(list []string) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(list))
	for _, s := range list {
		if s == "" {
			continue
		}
		id, err := uuid.Parse(s)
		if err == nil {
			out = append(out, id)
		}
	}
	return out
}

// parseAssetIDs menggabungkan image_asset_ids dan pdf_asset_id, mengabaikan
// nilai kosong/Nil agar tidak menimbulkan FK violation.
func parseAssetIDs(imageIDs []string, pdfID *string) []uuid.UUID {
	out := parseUUIDList(imageIDs)
	if pdfID != nil && *pdfID != "" {
		if id, err := uuid.Parse(*pdfID); err == nil {
			out = append(out, id)
		}
	}
	return out
}

func parseDecimal(s string, def decimal.Decimal) decimal.Decimal {
	if strings.TrimSpace(s) == "" {
		return def
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return def
	}
	return d
}

func parseDecimalPtr(s *string) *decimal.Decimal {
	if s == nil || strings.TrimSpace(*s) == "" {
		return nil
	}
	d, err := decimal.NewFromString(strings.TrimSpace(*s))
	if err != nil {
		return nil
	}
	return &d
}

func decimalPtrToString(d *decimal.Decimal) *string {
	if d == nil {
		return nil
	}
	v := d.String()
	return &v
}

func isWholeNumber(d decimal.Decimal) bool {
	return d.Equal(d.Ceil())
}

func (s *auctionService) generateCode() string {
	now := time.Now().UTC()
	code := fmt.Sprintf("AUC-%s-%s", now.Format("20060102"), strings.ToUpper(shortUUID()))
	return code
}

func shortUUID() string {
	id := uuid.New()
	return strings.ReplaceAll(id.String(), "-", "")[:8]
}

func (s *auctionService) stockAvailable(ctx context.Context, p models.Produk) int {
	reservations, err := s.repo.GetActiveReservationsForProducts(ctx, []uuid.UUID{p.ID})
	if err != nil {
		return 0
	}
	reserved := 0
	for _, r := range reservations {
		reserved += r.Quantity
	}
	avail := p.Quantity - reserved
	if avail < 0 {
		return 0
	}
	return avail
}

func (s *auctionService) wrapDBError(err error) error {
	var ae *AuctionError
	if errors.As(err, &ae) {
		return err
	}
	return err
}

func toJSONMap(v interface{}) models.JSONMap {
	b, _ := json.Marshal(v)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return models.JSONMap(m)
}

func ptrString(s string) *string { return &s }

func uuidPtrToString(p *uuid.UUID) *string {
	if p == nil {
		return nil
	}
	s := p.String()
	return &s
}

func mapAsset(a *models.AuctionAsset, cfg *config.Config) *dto.AuctionAssetResponse {
	return &dto.AuctionAssetResponse{
		ID:           a.ID.String(),
		Kind:         a.Kind,
		URL:          utils.GetFileURL(a.StorageKey, cfg),
		OriginalName: a.OriginalName,
		MimeType:     a.MimeType,
		SizeBytes:    a.SizeBytes,
	}
}

func mapBidDetail(b models.AuctionBid, isSelected bool) dto.AuctionBidDetail {
	effective := computeEffectivePercent(b.Amount, b.GrandTotalSnapshot)
	detail := dto.AuctionBidDetail{
		ID:                 b.ID.String(),
		BatchID:            b.BatchID.String(),
		Sequence:           b.Sequence,
		InputMode:          b.InputMode,
		Amount:             b.Amount.StringFixed(0),
		GrandTotalSnapshot: b.GrandTotalSnapshot.StringFixed(0),
		CreatedAt:          b.CreatedAt,
		IsSelected:         isSelected,
		EffectivePercent:   effective,
	}
	if b.Buyer != nil {
		detail.Buyer = dto.AuctionBuyerSimple{
			ID:      b.Buyer.ID.String(),
			Nama:    b.Buyer.Nama,
			Telepon: b.Buyer.Telepon,
		}
	}
	if b.InputPercent != nil {
		v := b.InputPercent.String()
		detail.InputPercent = &v
	}
	return detail
}

func computeEffectivePercent(amount, grandTotal decimal.Decimal) string {
	if grandTotal.IsZero() {
		return "0"
	}
	percent := amount.Div(grandTotal).Mul(decimal.NewFromInt(100)).Round(4)
	return percent.String()
}

func mapWinner(w *models.AuctionWinner, bid *models.AuctionBid) *dto.AuctionWinnerResponse {
	resp := &dto.AuctionWinnerResponse{
		BidID:             w.BidID.String(),
		DealAmount:        w.DealAmount.StringFixed(0),
		SelectedBy:        w.SelectedBy.String(),
		SelectedAt:        w.SelectedAt,
		Note:              w.Note,
		PaymentStatus:     w.PaymentStatus,
		PaidAt:            w.PaidAt,
		PaymentNote:       w.PaymentNote,
		FulfillmentStatus: w.FulfillmentStatus,
		CompletedAt:       w.CompletedAt,
		FulfillmentNote:   w.FulfillmentNote,
	}
	if bid != nil && bid.Buyer != nil {
		resp.Buyer = dto.AuctionBuyerSimple{
			ID:      bid.Buyer.ID.String(),
			Nama:    bid.Buyer.Nama,
			Telepon: bid.Buyer.Telepon,
		}
	}
	return resp
}

func lockClause() clause.Locking {
	return clause.Locking{Strength: "UPDATE"}
}
