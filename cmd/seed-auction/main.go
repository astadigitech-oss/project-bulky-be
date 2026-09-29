package main

import (
	"fmt"
	"log"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/pkg/database"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Seeder fixture QA modul Lelang (Auction) sesuai PRD 05A.
// Membuat batch DRAFT, OPEN (kosong), OPEN (multiple bid), dan SOLD,
// beserta bid, reservasi, winner, audit, event view, dan idempotency.
//
// Idempotent: jika kode batch "QA-..." sudah ada, dilewati.

const codePrefix = "QA-"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	cfg := config.LoadConfig()
	database.InitDB(cfg)
	db := database.GetDB()

	fmt.Println("=== Seeding QA Lelang (Auction) ===")
	fmt.Println()

	// Lookup master data.
	admin := findAdmin(db)
	if admin == nil {
		log.Fatal("Tidak ada admin aktif. Seed dibatalkan.")
	}
	warehouse := findWarehouse(db)
	products := findProducts(db, 4)
	buyers := findBuyers(db, 2)
	kategori := findKategori(db)
	kondisi := findKondisi(db)
	kondisiPaket := findKondisiPaket(db)
	sumber := findSumber(db)
	mereks := findMereks(db, 2)

	if len(products) < 2 {
		log.Fatal("Butuh minimal 2 produk aktif untuk seed. Dibatalkan.")
	}
	if len(buyers) < 2 {
		log.Fatal("Butuh minimal 2 buyer aktif untuk seed. Dibatalkan.")
	}

	fmt.Printf("Admin: %s\n", admin.Nama)
	fmt.Printf("Warehouse: %s\n", safeStr(warehouse))
	fmt.Printf("Produk: %d, Buyer: %d\n", len(products), len(buyers))

	// 1) DRAFT batch (tanpa item/aset).
	seedDraft(db, admin.ID)

	// 2) OPEN tanpa bid.
	seedOpenEmpty(db, admin.ID, warehouse, products, kategori, kondisi, kondisiPaket, sumber, mereks)

	// 3) OPEN dengan multiple bid (D02: A 100k & 200k, B 150k).
	seedOpenMultiBid(db, admin.ID, warehouse, products, buyers, kategori, kondisi, kondisiPaket, sumber, mereks)

	// 4) SOLD: winner B (150k), opened 10:00, sold 11:30 (D05).
	seedSold(db, admin.ID, warehouse, products, buyers, kategori, kondisi, kondisiPaket, sumber, mereks)

	fmt.Println()
	fmt.Println("=== Seeding selesai ===")
}

// ============================================================
// Lookup master
// ============================================================

func findAdmin(db *gorm.DB) *models.Admin {
	var admin models.Admin
	if err := db.Where("deleted_at IS NULL AND is_active = true").
		Order("created_at ASC").First(&admin).Error; err != nil {
		return nil
	}
	return &admin
}

func findWarehouse(db *gorm.DB) *uuid.UUID {
	var w models.Warehouse
	if err := db.Where("deleted_at IS NULL AND is_active = true").
		Order("created_at ASC").First(&w).Error; err != nil {
		return nil
	}
	return &w.ID
}

func findProducts(db *gorm.DB, n int) []models.Produk {
	var ps []models.Produk
	db.Where("deleted_at IS NULL AND is_active = true AND is_sold = false AND quantity > 0").
		Order("quantity DESC").Limit(n).Find(&ps)
	return ps
}

func findBuyers(db *gorm.DB, n int) []models.Buyer {
	var bs []models.Buyer
	db.Where("deleted_at IS NULL AND is_active = true").
		Order("created_at ASC").Limit(n).Find(&bs)
	return bs
}

func findKategori(db *gorm.DB) *uuid.UUID {
	var m models.KategoriProduk
	if err := db.Where("deleted_at IS NULL").Order("created_at ASC").First(&m).Error; err != nil {
		return nil
	}
	return &m.ID
}

func findKondisi(db *gorm.DB) *uuid.UUID {
	var m models.KondisiProduk
	if err := db.Where("deleted_at IS NULL").Order("created_at ASC").First(&m).Error; err != nil {
		return nil
	}
	return &m.ID
}

func findKondisiPaket(db *gorm.DB) *uuid.UUID {
	var m models.KondisiPaket
	if err := db.Where("deleted_at IS NULL").Order("created_at ASC").First(&m).Error; err != nil {
		return nil
	}
	return &m.ID
}

func findSumber(db *gorm.DB) *uuid.UUID {
	var m models.SumberProduk
	if err := db.Where("deleted_at IS NULL").Order("created_at ASC").First(&m).Error; err != nil {
		return nil
	}
	return &m.ID
}

func findMereks(db *gorm.DB, n int) []uuid.UUID {
	var ms []models.MerekProduk
	db.Where("deleted_at IS NULL").Order("created_at ASC").Limit(n).Find(&ms)
	ids := make([]uuid.UUID, 0, len(ms))
	for _, m := range ms {
		ids = append(ids, m.ID)
	}
	return ids
}

// ============================================================
// Fixture helpers
// ============================================================

// ensureBatchCreated checks if a QA batch with given code exists.
func batchExists(db *gorm.DB, code string) bool {
	var count int64
	db.Model(&models.AuctionBatch{}).Where("code = ?", code).Count(&count)
	return count > 0
}

func makeCode(suffix string) string {
	return codePrefix + suffix + "-" + shortID()
}

func shortID() string {
	return stringsUpper(uuid.NewString()[:6])
}

func stringsUpper(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'a' && b[i] <= 'z' {
			b[i] -= 32
		}
	}
	return string(b)
}

func ceilDecimal(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f).Ceil()
}

func safeStr(p *uuid.UUID) string {
	if p == nil {
		return "-"
	}
	return p.String()
}

// buildBatchItems membuat snapshot item dari produk.
func buildBatchItems(products []models.Produk, quantities []int) ([]models.AuctionBatchItem, decimal.Decimal, int) {
	items := make([]models.AuctionBatchItem, 0, len(products))
	var grandTotal decimal.Decimal = decimal.Zero
	totalQty := 0
	for i, p := range products {
		qty := quantities[i]
		unit := ceilDecimal(p.HargaSesudahDiskon)
		subtotal := unit.Mul(decimal.NewFromInt(int64(qty)))
		grandTotal = grandTotal.Add(subtotal)
		totalQty += qty
		items = append(items, models.AuctionBatchItem{
			ProdukID:          &p.ID,
			SourceType:        "CATALOG",
			Quantity:          qty,
			NamaSnapshot:      p.NamaID,
			UnitPriceSnapshot: unit,
			SubtotalSnapshot:  subtotal,
		})
	}
	return items, grandTotal, totalQty
}

// publishBatch sets batch to OPEN + reserves stock.
func publishBatch(tx *gorm.DB, batch *models.AuctionBatch, adminID uuid.UUID, items []models.AuctionBatchItem, warehouseID *uuid.UUID, kategoriID, kondisiID, kondisiPaketID, sumberID *uuid.UUID, merekIDs []uuid.UUID) error {
	now := time.Now().UTC()
	batch.WarehouseID = warehouseID
	batch.KategoriID = kategoriID
	batch.KondisiID = kondisiID
	batch.KondisiPaketID = kondisiPaketID
	batch.SumberID = sumberID
	batch.Status = models.AuctionBatchStatusOPEN
	// Hanya set opened_at jika belum diisi (agar fixture D05 bisa menetapkan
	// waktu buka yang disengaja).
	if batch.OpenedAt == nil {
		batch.OpenedAt = &now
	}
	batch.UpdatedBy = adminID
	batch.Version++
	if err := tx.Save(batch).Error; err != nil {
		return err
	}
	// Insert items.
	for i := range items {
		items[i].BatchID = batch.ID
	}
	if len(items) > 0 {
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
	}
	// Brands pivot.
	if len(merekIDs) > 0 {
		brands := make([]models.AuctionBatchBrand, 0, len(merekIDs))
		for _, m := range merekIDs {
			brands = append(brands, models.AuctionBatchBrand{BatchID: batch.ID, MerekID: m})
		}
		if err := tx.Create(&brands).Error; err != nil {
			return err
		}
	}
	// Reservasi stok.
	for _, it := range items {
		if it.ProdukID == nil {
			continue
		}
		res := &models.AuctionStockReservation{
			BatchID:  batch.ID,
			ProdukID: *it.ProdukID,
			Quantity: it.Quantity,
			Status:   "ACTIVE",
		}
		if err := tx.Create(res).Error; err != nil {
			return err
		}
	}
	// Audit log.
	if err := tx.Create(&models.AuctionAuditLog{
		BatchID:      batch.ID,
		ActorAdminID: adminID,
		Action:       "publish",
		AfterData:    models.JSONMap{"status": "OPEN"},
	}).Error; err != nil {
		return err
	}
	return nil
}

// ============================================================
// Fixtures
// ============================================================

func seedDraft(db *gorm.DB, adminID uuid.UUID) {
	code := codePrefix + "DRAFT-1"
	if batchExists(db, code) {
		fmt.Printf("⊘ Skip DRAFT (sudah ada)\n")
		return
	}
	batch := &models.AuctionBatch{
		Code:          code,
		NamaID:        "QA DRAFT Batch",
		NamaEN:        strPtr("QA Draft Batch"),
		Status:        models.AuctionBatchStatusDRAFT,
		MinBidPercent: decimal.NewFromFloat(0.1),
		Version:       1,
		CreatedBy:     adminID,
		UpdatedBy:     adminID,
	}
	if err := db.Create(batch).Error; err != nil {
		log.Printf("✗ DRAFT gagal: %v\n", err)
		return
	}
	fmt.Printf("✓ DRAFT: %s (%s)\n", batch.NamaID, batch.ID)
}

func seedOpenEmpty(db *gorm.DB, adminID uuid.UUID, warehouseID *uuid.UUID, products []models.Produk, kategoriID, kondisiID, kondisiPaketID, sumberID *uuid.UUID, merekIDs []uuid.UUID) {
	code := codePrefix + "OPEN-EMPTY-1"
	if batchExists(db, code) {
		fmt.Printf("⊘ Skip OPEN empty (sudah ada)\n")
		return
	}
	batch := &models.AuctionBatch{
		Code:          code,
		NamaID:        "QA OPEN Kosong",
		NamaEN:        strPtr("QA Open Empty"),
		Status:        models.AuctionBatchStatusDRAFT,
		MinBidPercent: decimal.NewFromFloat(0.1),
		Version:       1,
		CreatedBy:     adminID,
		UpdatedBy:     adminID,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(batch).Error; err != nil {
			return err
		}
		items, grandTotal, totalQty := buildBatchItems(products[:1], []int{3})
		batch.GrandTotal = grandTotal
		batch.TotalQuantity = totalQty
		batch.PanjangCm = decimal.NewFromInt(100)
		batch.LebarCm = decimal.NewFromInt(80)
		batch.TinggiCm = decimal.NewFromInt(60)
		batch.BeratKg = decimal.NewFromInt(50)
		batch.VolumeM3 = decimal.NewFromInt(100).Mul(decimal.NewFromInt(80)).Mul(decimal.NewFromInt(60)).Div(decimal.NewFromInt(1000000))
		return publishBatch(tx, batch, adminID, items, warehouseID, kategoriID, kondisiID, kondisiPaketID, sumberID, merekIDs)
	})
	if err != nil {
		log.Printf("✗ OPEN empty gagal: %v\n", err)
		return
	}
	fmt.Printf("✓ OPEN Kosong: %s (%s)\n", batch.NamaID, batch.ID)
}

func seedOpenMultiBid(db *gorm.DB, adminID uuid.UUID, warehouseID *uuid.UUID, products []models.Produk, buyers []models.Buyer, kategoriID, kondisiID, kondisiPaketID, sumberID *uuid.UUID, merekIDs []uuid.UUID) {
	code := codePrefix + "OPEN-MULTI-1"
	if batchExists(db, code) {
		fmt.Printf("⊘ Skip OPEN multi (sudah ada)\n")
		return
	}
	batch := &models.AuctionBatch{
		Code:          code,
		NamaID:        "QA OPEN Multi Bid",
		NamaEN:        strPtr("QA Open Multi Bid"),
		Status:        models.AuctionBatchStatusDRAFT,
		MinBidPercent: decimal.NewFromFloat(0.1),
		Version:       1,
		CreatedBy:     adminID,
		UpdatedBy:     adminID,
	}
	var batchID uuid.UUID
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(batch).Error; err != nil {
			return err
		}
		items, grandTotal, totalQty := buildBatchItems(products[:2], []int{2, 2})
		batch.GrandTotal = grandTotal
		batch.TotalQuantity = totalQty
		batch.PanjangCm = decimal.NewFromInt(100)
		batch.LebarCm = decimal.NewFromInt(80)
		batch.TinggiCm = decimal.NewFromInt(60)
		batch.BeratKg = decimal.NewFromInt(50)
		batch.VolumeM3 = decimal.NewFromInt(100).Mul(decimal.NewFromInt(80)).Mul(decimal.NewFromInt(60)).Div(decimal.NewFromInt(1000000))
		if err := publishBatch(tx, batch, adminID, items, warehouseID, kategoriID, kondisiID, kondisiPaketID, sumberID, merekIDs); err != nil {
			return err
		}
		batchID = batch.ID
		// Bid: buyerA seq1=100k, buyerA seq2=200k, buyerB seq1=150k.
		bids := []models.AuctionBid{
			{BatchID: batchID, BuyerID: buyers[0].ID, Sequence: 1, InputMode: "AMOUNT", Amount: decimal.NewFromInt(100000), GrandTotalSnapshot: batch.GrandTotal},
			{BatchID: batchID, BuyerID: buyers[0].ID, Sequence: 2, InputMode: "AMOUNT", Amount: decimal.NewFromInt(200000), GrandTotalSnapshot: batch.GrandTotal},
			{BatchID: batchID, BuyerID: buyers[1].ID, Sequence: 1, InputMode: "AMOUNT", Amount: decimal.NewFromInt(150000), GrandTotalSnapshot: batch.GrandTotal},
		}
		return tx.Create(&bids).Error
	})
	if err != nil {
		log.Printf("✗ OPEN multi gagal: %v\n", err)
		return
	}
	fmt.Printf("✓ OPEN Multi Bid: %s (%s) grand_total=%s\n", batch.NamaID, batch.ID, batch.GrandTotal.String())
}

func seedSold(db *gorm.DB, adminID uuid.UUID, warehouseID *uuid.UUID, products []models.Produk, buyers []models.Buyer, kategoriID, kondisiID, kondisiPaketID, sumberID *uuid.UUID, merekIDs []uuid.UUID) {
	code := codePrefix + "SOLD-1"
	if batchExists(db, code) {
		fmt.Printf("⊘ Skip SOLD (sudah ada)\n")
		return
	}
	batch := &models.AuctionBatch{
		Code:          code,
		NamaID:        "QA SOLD Batch",
		NamaEN:        strPtr("QA Sold Batch"),
		Status:        models.AuctionBatchStatusDRAFT,
		MinBidPercent: decimal.NewFromFloat(0.1),
		Version:       1,
		CreatedBy:     adminID,
		UpdatedBy:     adminID,
	}
	var batchID uuid.UUID
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(batch).Error; err != nil {
			return err
		}
		items, grandTotal, totalQty := buildBatchItems(products[:1], []int{2})
		batch.GrandTotal = grandTotal
		batch.TotalQuantity = totalQty
		batch.PanjangCm = decimal.NewFromInt(100)
		batch.LebarCm = decimal.NewFromInt(80)
		batch.TinggiCm = decimal.NewFromInt(60)
		batch.BeratKg = decimal.NewFromInt(50)
		batch.VolumeM3 = decimal.NewFromInt(100).Mul(decimal.NewFromInt(80)).Mul(decimal.NewFromInt(60)).Div(decimal.NewFromInt(1000000))
		// Publish: opened_at = now - 1.5h (agar time_to_sold = 5400s).
		openedAt := time.Now().UTC().Add(-90 * time.Minute)
		batch.OpenedAt = &openedAt
		if err := publishBatch(tx, batch, adminID, items, warehouseID, kategoriID, kondisiID, kondisiPaketID, sumberID, merekIDs); err != nil {
			return err
		}
		batchID = batch.ID

		// Bids: A=100k (seq1), B=150k (seq1). Pilih B sebagai winner.
		bidA := models.AuctionBid{BatchID: batchID, BuyerID: buyers[0].ID, Sequence: 1, InputMode: "AMOUNT", Amount: decimal.NewFromInt(100000), GrandTotalSnapshot: batch.GrandTotal}
		bidB := models.AuctionBid{BatchID: batchID, BuyerID: buyers[1].ID, Sequence: 1, InputMode: "AMOUNT", Amount: decimal.NewFromInt(150000), GrandTotalSnapshot: batch.GrandTotal}
		if err := tx.Create(&bidA).Error; err != nil {
			return err
		}
		if err := tx.Create(&bidB).Error; err != nil {
			return err
		}

		// Winner = bid B. sold_at = now, opened_at = now-90min => time_to_sold 5400s.
		soldAt := time.Now().UTC()
		winner := &models.AuctionWinner{
			BatchID:           batchID,
			BidID:             bidB.ID,
			DealAmount:        bidB.Amount,
			SelectedBy:        adminID,
			SelectedAt:        soldAt,
			PaymentStatus:     "PAID",
			PaidAt:            &soldAt,
			PaymentNote:       strPtr("Transfer QA"),
			FulfillmentStatus: "COMPLETED",
			CompletedAt:       &soldAt,
			FulfillmentNote:   strPtr("Serah terima QA"),
		}
		if err := tx.Create(winner).Error; err != nil {
			return err
		}

		// SOLD status + sold_at.
		batch.Status = models.AuctionBatchStatusSOLD
		batch.SoldAt = &soldAt
		batch.UpdatedBy = adminID
		batch.Version++
		if err := tx.Save(batch).Error; err != nil {
			return err
		}

		// Konsumsi reservasi & kurangi stok.
		var reservations []models.AuctionStockReservation
		if err := tx.Where("batch_id = ? AND status = ?", batchID, "ACTIVE").Find(&reservations).Error; err != nil {
			return err
		}
		for _, r := range reservations {
			if err := tx.Model(&models.AuctionStockReservation{}).Where("id = ?", r.ID).
				Updates(map[string]interface{}{"status": "CONSUMED", "consumed_at": soldAt}).Error; err != nil {
				return err
			}
			var p models.Produk
			if err := tx.First(&p, "id = ?", r.ProdukID).Error; err != nil {
				return err
			}
			newQty := p.Quantity - r.Quantity
			if newQty < 0 {
				newQty = 0
			}
			updates := map[string]interface{}{"quantity": newQty}
			if newQty == 0 {
				updates["is_sold"] = true
			}
			if err := tx.Model(&models.Produk{}).Where("id = ?", r.ProdukID).Updates(updates).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("✗ SOLD gagal: %v\n", err)
		return
	}
	fmt.Printf("✓ SOLD: %s (%s) deal=150000 time_to_sold=5400s\n", batch.NamaID, batch.ID)
}

func strPtr(s string) *string { return &s }
