package services

import (
	"context"
	"errors"
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
	// Validate and build Tujuan list
	var tujuanList models.TujuanList
	if len(req.Tujuan) > 0 {
		kategoriIDs := make([]string, len(req.Tujuan))
		for i, t := range req.Tujuan {
			kategoriIDs[i] = t.ID
		}

		// Validate all kategori exist and active
		existingKategori, err := s.kategoriService.FindActiveByIDs(ctx, kategoriIDs)
		if err != nil {
			return nil, errors.New("gagal validasi kategori: " + err.Error())
		}

		if len(existingKategori) != len(kategoriIDs) {
			return nil, errors.New("beberapa kategori tidak ditemukan atau tidak aktif")
		}

		// Build tujuan list
		tujuanList = make(models.TujuanList, len(existingKategori))
		for i, k := range existingKategori {
			tujuanList[i] = models.TujuanKategori{
				ID:   k.ID.String(),
				Slug: k.Slug,
			}
		}
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
		Tujuan:      tujuanList,
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
		if len(req.Tujuan) > 0 {
			kategoriIDs := make([]string, len(req.Tujuan))
			for i, t := range req.Tujuan {
				kategoriIDs[i] = t.ID
			}

			existingKategori, err := s.kategoriService.FindActiveByIDs(ctx, kategoriIDs)
			if err != nil {
				return nil, errors.New("gagal validasi kategori: " + err.Error())
			}

			if len(existingKategori) != len(kategoriIDs) {
				return nil, errors.New("beberapa kategori tidak ditemukan atau tidak aktif")
			}

			tujuanList := make(models.TujuanList, len(existingKategori))
			for i, k := range existingKategori {
				tujuanList[i] = models.TujuanKategori{
					ID:   k.ID.String(),
					Slug: k.Slug,
				}
			}
			banner.Tujuan = tujuanList
		} else {
			// Empty array = clear tujuan
			banner.Tujuan = nil
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
		var tujuan []models.TujuanKategoriPublicResponse
		if b.Tujuan != nil {
			tujuan = make([]models.TujuanKategoriPublicResponse, len(b.Tujuan))
			for i, t := range b.Tujuan {
				tujuan[i] = models.TujuanKategoriPublicResponse{
					ID:   t.ID,
					Slug: t.Slug,
				}
			}
		}

		items = append(items, models.BannerEventPromoPublicResponse{
			ID:        b.ID.String(),
			Nama:      b.Nama,
			GambarURL: b.GetGambarURL().GetFullURL(s.cfg.BaseURL),
			Tujuan:    tujuan,
		})
	}

	return items, nil
}

func (s *bannerEventPromoService) toResponse(b *models.BannerEventPromo) *models.BannerEventPromoResponse {
	var tujuan []models.TujuanKategoriResponse
	if b.Tujuan != nil {
		tujuan = make([]models.TujuanKategoriResponse, len(b.Tujuan))
		for i, t := range b.Tujuan {
			tujuan[i] = models.TujuanKategoriResponse{
				ID:   t.ID,
				Slug: t.Slug,
			}
		}
	}

	return &models.BannerEventPromoResponse{
		ID:             b.ID.String(),
		Nama:           b.Nama,
		GambarURL:      b.GetGambarURL().GetFullURL(s.cfg.BaseURL),
		Tujuan:         tujuan,
		Urutan:         b.Urutan,
		IsVisible:      b.IsCurrentlyVisible(),
		TanggalMulai:   b.TanggalMulai,
		TanggalSelesai: b.TanggalSelesai,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}

func (s *bannerEventPromoService) toSimpleResponse(b *models.BannerEventPromo) *models.BannerEventPromoSimpleResponse {
	var tujuan []models.TujuanKategoriResponse
	if b.Tujuan != nil {
		tujuan = make([]models.TujuanKategoriResponse, len(b.Tujuan))
		for i, t := range b.Tujuan {
			tujuan[i] = models.TujuanKategoriResponse{
				ID:   t.ID,
				Slug: t.Slug,
			}
		}
	}

	return &models.BannerEventPromoSimpleResponse{
		ID:        b.ID.String(),
		Nama:      b.Nama,
		GambarURL: b.GetGambarURL().GetFullURL(s.cfg.BaseURL),
		// Tujuan:    tujuan,
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
