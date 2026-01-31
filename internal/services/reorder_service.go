package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReorderService struct {
	db *gorm.DB
}

func NewReorderService(db *gorm.DB) *ReorderService {
	return &ReorderService{db: db}
}

type ReorderResult struct {
	ItemID        uuid.UUID `json:"item_id"`
	ItemUrutan    int       `json:"item_urutan"`
	SwappedID     uuid.UUID `json:"swapped_id"`
	SwappedUrutan int       `json:"swapped_urutan"`
}

// Reorder - Generic reorder function for any table with 'urutan' field
// tableName: nama tabel (e.g., "hero_section")
// id: ID item yang mau di-reorder
// direction: "up" atau "down"
// scopeColumn: optional WHERE condition column untuk scoping (e.g., "produk_id")
// scopeValue: optional WHERE condition value untuk scoping
func (s *ReorderService) Reorder(
	ctx context.Context,
	tableName string,
	id uuid.UUID,
	direction string,
	scopeColumn string,
	scopeValue interface{},
) (*ReorderResult, error) {

	// Validate direction
	if direction != "up" && direction != "down" {
		return nil, errors.New("direction harus 'up' atau 'down'")
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get current item
	var currentUrutan int
	query := tx.Table(tableName).
		Select("urutan").
		Where("id = ? AND deleted_at IS NULL", id)

	if scopeColumn != "" && scopeValue != nil {
		query = query.Where(scopeColumn+" = ?", scopeValue)
	}

	queryResult := query.Scan(&currentUrutan)
	if queryResult.Error != nil {
		tx.Rollback()
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, queryResult.Error
	}

	// Check if record actually found
	if queryResult.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.New("data tidak ditemukan")
	}

	// Find adjacent item to swap with
	var adjacentID uuid.UUID
	var adjacentUrutan int

	adjacentQuery := tx.Table(tableName).
		Select("id, urutan").
		Where("deleted_at IS NULL")

	if scopeColumn != "" && scopeValue != nil {
		adjacentQuery = adjacentQuery.Where(scopeColumn+" = ?", scopeValue)
	}

	if direction == "up" {
		// Find item with urutan < current, closest one (ORDER BY urutan DESC)
		adjacentQuery = adjacentQuery.
			Where("urutan < ?", currentUrutan).
			Order("urutan DESC").
			Limit(1)
	} else {
		// Find item with urutan > current, closest one (ORDER BY urutan ASC)
		adjacentQuery = adjacentQuery.
			Where("urutan > ?", currentUrutan).
			Order("urutan ASC").
			Limit(1)
	}

	var adjacentResult struct {
		ID     uuid.UUID
		Urutan int
	}

	if err := adjacentQuery.Scan(&adjacentResult).Error; err != nil {
		tx.Rollback()
		if direction == "up" {
			return nil, errors.New("item sudah berada di urutan paling atas")
		}
		return nil, errors.New("item sudah berada di urutan paling bawah")
	}

	if adjacentResult.ID == uuid.Nil {
		tx.Rollback()
		if direction == "up" {
			return nil, errors.New("item sudah berada di urutan paling atas")
		}
		return nil, errors.New("item sudah berada di urutan paling bawah")
	}

	adjacentID = adjacentResult.ID
	adjacentUrutan = adjacentResult.Urutan

	// Swap urutan values
	// Update current item → adjacent's urutan
	if err := tx.Table(tableName).
		Where("id = ?", id).
		Update("urutan", adjacentUrutan).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update adjacent item → current's urutan
	if err := tx.Table(tableName).
		Where("id = ?", adjacentID).
		Update("urutan", currentUrutan).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &ReorderResult{
		ItemID:        id,
		ItemUrutan:    adjacentUrutan,
		SwappedID:     adjacentID,
		SwappedUrutan: currentUrutan,
	}, nil
}

// ReorderAfterDelete - Reorder items after deletion to fill gaps
// tableName: nama tabel (e.g., "banner_event_promo")
// deletedUrutan: urutan dari item yang dihapus
// scopeColumn: optional WHERE condition column untuk scoping (e.g., "tipe_produk_id")
// scopeValue: optional WHERE condition value untuk scoping
func (s *ReorderService) ReorderAfterDelete(
	ctx context.Context,
	tableName string,
	deletedUrutan int,
	scopeColumn string,
	scopeValue interface{},
) error {
	// Build base query
	query := s.db.WithContext(ctx).Table(tableName).
		Where("urutan > ? AND deleted_at IS NULL", deletedUrutan)

	// Add scope condition if provided
	if scopeColumn != "" && scopeValue != nil {
		query = query.Where(scopeColumn+" = ?", scopeValue)
	}

	// Decrement urutan for all items after deleted item
	return query.UpdateColumn("urutan", gorm.Expr("urutan - 1")).Error
}
