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
	// Validate tujuan
	if err := s.validateTujuan(req.Tujuan); err != nil {
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

	// Set tujuan (trim whitespace, remove empty)
	if req.Tujuan != "" {
		tujuan := s.cleanTujuanString(req.Tujuan)
		banner.Tujuan = &tujuan
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

	if err := s.repo.Create(ctx, banner); err != nil {
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
	if req.Tujuan != nil {
		if *req.Tujuan == "" {
			// Empty string = clear tujuan
			banner.Tujuan = nil
		} else {
			// Validate and clean tujuan
			if err := s.validateTujuan(*req.Tujuan); err != nil {
				return nil, err
			}
			tujuan := s.cleanTujuanString(*req.Tujuan)
			banner.Tujuan = &tujuan
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

	if err := s.repo.Update(ctx, banner); err != nil {
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
			Tujuan:    b.Tujuan,
		})
	}

	return items, nil
}

func (s *bannerEventPromoService) toResponse(b *models.BannerEventPromo) *models.BannerEventPromoResponse {
	return &models.BannerEventPromoResponse{
		ID:             b.ID.String(),
		Nama:           b.Nama,
		GambarURL:      b.GetGambarURL().GetFullURL(s.cfg.BaseURL),
		Tujuan:         b.Tujuan,
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
		Tujuan:    b.Tujuan,
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

// validateTujuan validates comma-separated kategori IDs
func (s *bannerEventPromoService) validateTujuan(tujuanStr string) error {
	if tujuanStr == "" {
		return nil // Empty is allowed (no redirect)
	}

	ids := strings.Split(tujuanStr, ",")

	for _, idStr := range ids {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		// Validate UUID format
		id, err := uuid.Parse(idStr)
		if err != nil {
			return fmt.Errorf("invalid UUID format: %s", idStr)
		}

		// Check if kategori exists
		existingKategori, err := s.kategoriService.FindActiveByIDs(context.Background(), []string{id.String()})
		if err != nil {
			return err
		}
		if len(existingKategori) == 0 {
			return fmt.Errorf("kategori produk tidak ditemukan: %s", idStr)
		}
	}

	return nil
}

// cleanTujuanString removes whitespace and empty entries
func (s *bannerEventPromoService) cleanTujuanString(tujuan string) string {
	ids := strings.Split(tujuan, ",")
	cleaned := make([]string, 0, len(ids))

	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id != "" {
			cleaned = append(cleaned, id)
		}
	}

	return strings.Join(cleaned, ",")
}
