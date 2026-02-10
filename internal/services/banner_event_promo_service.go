package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type BannerEventPromoService interface {
	Create(ctx context.Context, req *models.CreateBannerEventPromoRequest) (*models.BannerEventPromoResponse, error)
	FindByID(ctx context.Context, id string) (*models.BannerEventPromoResponse, error)
	FindAll(ctx context.Context, params *models.BannerEventPromoFilterRequest) ([]models.BannerEventPromoSimpleResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateBannerEventPromoRequest) (*models.BannerEventPromoResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.BannerEventPromoResponse, bool, error) // Returns: response, wasActivated, error
	Reorder(ctx context.Context, req *models.ReorderRequest) error
	GetVisibleBanners(ctx context.Context) ([]models.BannerEventPromoPublicResponse, error)
}

type bannerEventPromoService struct {
	repo            repositories.BannerEventPromoRepository
	reorderService  *ReorderService
	kategoriService KategoriProdukService
	cfg             *config.Config
}

func NewBannerEventPromoService(repo repositories.BannerEventPromoRepository, reorderService *ReorderService, kategoriService KategoriProdukService, cfg *config.Config) BannerEventPromoService {
	return &bannerEventPromoService{
		repo:            repo,
		reorderService:  reorderService,
		kategoriService: kategoriService,
		cfg:             cfg,
	}
}

func (s *bannerEventPromoService) Create(ctx context.Context, req *models.CreateBannerEventPromoRequest) (*models.BannerEventPromoResponse, error) {
	// Validate and parse tujuan IDs
	kategoriIDs, err := s.validateAndParseKategoriIDs(req.Tujuan)
	if err != nil {
		return nil, err
	}

	// Auto-increment urutan
	maxUrutan, err := s.repo.GetMaxUrutan(ctx)
	if err != nil {
		return nil, err
	}

	banner := &models.BannerEventPromo{
		ID:          uuid.New(),
		Nama:        req.Nama,
		GambarURLID: req.GambarID,
		GambarURLEN: req.GambarEN,
		Urutan:      maxUrutan + 1,
	}

	if req.TanggalMulai != nil {
		if t, err := parseFlexibleDate(*req.TanggalMulai); err == nil {
			banner.TanggalMulai = &t
		}
	}

	if req.TanggalSelesai != nil {
		if t, err := parseFlexibleDate(*req.TanggalSelesai); err == nil {
			banner.TanggalSelesai = &t
		}
	}

	// Create banner and pivot relations in transaction
	if err := s.repo.CreateWithKategori(ctx, banner, kategoriIDs); err != nil {
		return nil, err
	}

	// Reload with relations
	banner, err = s.repo.FindByID(ctx, banner.ID.String())
	if err != nil {
		return nil, err
	}

	return s.toResponse(banner), nil
}

func (s *bannerEventPromoService) FindByID(ctx context.Context, id string) (*models.BannerEventPromoResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}
	return s.toResponse(banner), nil
}

func (s *bannerEventPromoService) FindAll(ctx context.Context, params *models.BannerEventPromoFilterRequest) ([]models.BannerEventPromoSimpleResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	banners, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	items := []models.BannerEventPromoSimpleResponse{}
	for _, b := range banners {
		items = append(items, *s.toSimpleResponse(&b))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *bannerEventPromoService) Update(ctx context.Context, id string, req *models.UpdateBannerEventPromoRequest) (*models.BannerEventPromoResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}

	if req.Nama != nil {
		banner.Nama = *req.Nama
	}
	if req.GambarID != nil {
		banner.GambarURLID = *req.GambarID
	}
	if req.GambarEN != nil {
		banner.GambarURLEN = *req.GambarEN
	}

	// Handle tujuan update
	var kategoriIDs []uuid.UUID
	if req.Tujuan != nil {
		// Validate and parse kategori IDs
		kategoriIDs, err = s.validateAndParseKategoriIDs(*req.Tujuan)
		if err != nil {
			return nil, err
		}
	}

	if req.TanggalMulai != nil {
		if t, err := parseFlexibleDate(*req.TanggalMulai); err == nil {
			banner.TanggalMulai = &t
		}
	}
	if req.TanggalSelesai != nil {
		if t, err := parseFlexibleDate(*req.TanggalSelesai); err == nil {
			banner.TanggalSelesai = &t
		}
	}

	// Update banner and replace kategori relations in transaction
	if req.Tujuan != nil {
		if err := s.repo.UpdateWithKategori(ctx, banner, kategoriIDs); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.Update(ctx, banner); err != nil {
			return nil, err
		}
	}

	// Reload with relations
	banner, err = s.repo.FindByID(ctx, banner.ID.String())
	if err != nil {
		return nil, err
	}

	return s.toResponse(banner), nil
}

func (s *bannerEventPromoService) Delete(ctx context.Context, id string) error {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("banner tidak ditemukan")
	}

	deletedUrutan := banner.Urutan

	// Soft delete banner
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Reorder remaining items to fill gap
	return s.reorderService.ReorderAfterDelete(
		ctx,
		"banner_event_promo",
		deletedUrutan,
		"", // No scope
		nil,
	)
}

func (s *bannerEventPromoService) ToggleStatus(ctx context.Context, id string) (*models.BannerEventPromoResponse, bool, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, false, errors.New("banner tidak ditemukan")
	}

	now := time.Now()
	past := now.Add(-1 * time.Second)
	wasActivated := false

	// SIMPLE LOGIC: Cek tanggal_selesai NULL atau engga
	if banner.TanggalSelesai == nil {
		// tanggal_selesai NULL = Banner AKTIF -> DEACTIVATE
		banner.TanggalSelesai = &past
		if err := s.repo.UpdateToggleStatus(ctx, id, nil, &past); err != nil {
			return nil, false, err
		}
		wasActivated = false // We DEACTIVATED
	} else {
		// tanggal_selesai ada value = Banner NONAKTIF -> ACTIVATE
		banner.TanggalMulai = &past
		banner.TanggalSelesai = nil
		if err := s.repo.UpdateToggleStatus(ctx, id, &past, nil); err != nil {
			return nil, false, err
		}
		wasActivated = true // We ACTIVATED
	}

	// Reload with relations
	banner, err = s.repo.FindByID(ctx, banner.ID.String())
	if err != nil {
		return nil, wasActivated, err
	}

	return s.toResponse(banner), wasActivated, nil
}

func (s *bannerEventPromoService) Reorder(ctx context.Context, req *models.ReorderRequest) error {
	return s.repo.UpdateOrder(ctx, req.Items)
}

func (s *bannerEventPromoService) GetVisibleBanners(ctx context.Context) ([]models.BannerEventPromoPublicResponse, error) {
	banners, err := s.repo.GetVisibleBanners(ctx)
	if err != nil {
		return nil, err
	}

	items := []models.BannerEventPromoPublicResponse{}
	for _, b := range banners {
		items = append(items, models.BannerEventPromoPublicResponse{
			ID:        b.ID.String(),
			Nama:      b.Nama,
			GambarURL: b.GetGambarURL().GetFullURL(s.cfg.BaseURL),
			Tujuan:    b.GetKategoriIDStrings(),
		})
	}

	return items, nil
}

func (s *bannerEventPromoService) toResponse(b *models.BannerEventPromo) *models.BannerEventPromoResponse {
	return &models.BannerEventPromoResponse{
		ID:             b.ID.String(),
		Nama:           b.Nama,
		GambarURL:      b.GetGambarURL().GetFullURL(s.cfg.BaseURL),
		Tujuan:         b.GetKategoriIDStrings(),
		Urutan:         b.Urutan,
		IsVisible:      b.IsCurrentlyVisible(),
		TanggalMulai:   b.TanggalMulai,
		TanggalSelesai: b.TanggalSelesai,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}

func (s *bannerEventPromoService) toSimpleResponse(b *models.BannerEventPromo) *models.BannerEventPromoSimpleResponse {
	return &models.BannerEventPromoSimpleResponse{
		ID:        b.ID.String(),
		Nama:      b.Nama,
		GambarURL: b.GetGambarURL().GetFullURL(s.cfg.BaseURL),
		Tujuan:    b.GetKategoriIDStrings(),
		Urutan:    b.Urutan,
		IsVisible: b.IsCurrentlyVisible(),
		UpdatedAt: b.UpdatedAt,
	}
}

// parseFlexibleDate parses date string in multiple formats
// Supports: "2026-01-10" (date only) or "2026-01-10T00:00:00Z" (RFC3339)
func parseFlexibleDate(dateStr string) (time.Time, error) {
	// Try RFC3339 first (full datetime with timezone)
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t, nil
	}

	// Try date only format (yyyy-mm-dd)
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}

	// Try datetime without timezone
	if t, err := time.Parse("2006-01-02T15:04:05", dateStr); err == nil {
		return t, nil
	}

	return time.Time{}, errors.New("format tanggal tidak valid, gunakan yyyy-mm-dd atau RFC3339")
}

// validateAndParseKategoriIDs validates and converts string IDs to UUIDs
func (s *bannerEventPromoService) validateAndParseKategoriIDs(ids []string) ([]uuid.UUID, error) {
	if len(ids) == 0 {
		return []uuid.UUID{}, nil
	}

	result := make([]uuid.UUID, 0, len(ids))

	for _, idStr := range ids {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		// Parse UUID
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID format: %s", idStr)
		}

		// Check if kategori exists
		existingKategori, err := s.kategoriService.FindActiveByIDs(context.Background(), []string{id.String()})
		if err != nil {
			return nil, err
		}
		if len(existingKategori) == 0 {
			return nil, fmt.Errorf("kategori produk tidak ditemukan: %s", idStr)
		}

		result = append(result, id)
	}

	return result, nil
}
