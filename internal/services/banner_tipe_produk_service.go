package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BannerTipeProdukService interface {
	Create(ctx context.Context, req *models.CreateBannerTipeProdukRequest) (*models.BannerTipeProdukResponse, error)
	FindByID(ctx context.Context, id string) (*models.BannerTipeProdukResponse, error)
	FindAll(ctx context.Context, params *models.BannerTipeProdukFilterRequest, tipeProdukID string) ([]models.BannerTipeProdukSimpleResponse, *models.PaginationMeta, error)
	FindAllGrouped(ctx context.Context, search string) (*models.BannerTipeProdukGroupedResponse, *models.BannerTipeProdukGroupedMeta, error) // NEW
	FindByTipeProdukID(ctx context.Context, tipeProdukID string) ([]models.BannerSimpleResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateBannerTipeProdukRequest) (*models.BannerTipeProdukResponse, error)
	UpdateWithFile(ctx context.Context, id string, req *models.UpdateBannerTipeProdukRequest, newGambarURL *string) (*models.BannerTipeProdukResponse, error)
	Delete(ctx context.Context, id string) error
	DeleteWithFile(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	Reorder(ctx context.Context, req *models.ReorderRequest) error
}

type bannerTipeProdukService struct {
	repo           repositories.BannerTipeProdukRepository
	tipeProdukRepo repositories.TipeProdukRepository
	reorderService *ReorderService
	cfg            *config.Config
}

func NewBannerTipeProdukService(
	repo repositories.BannerTipeProdukRepository,
	tipeProdukRepo repositories.TipeProdukRepository,
	reorderService *ReorderService,
	cfg *config.Config,
) BannerTipeProdukService {
	return &bannerTipeProdukService{
		repo:           repo,
		tipeProdukRepo: tipeProdukRepo,
		reorderService: reorderService,
		cfg:            cfg,
	}
}

func (s *bannerTipeProdukService) Create(ctx context.Context, req *models.CreateBannerTipeProdukRequest) (*models.BannerTipeProdukResponse, error) {
	tipeProdukUUID, err := uuid.Parse(req.TipeProdukID)
	if err != nil {
		return nil, errors.New("tipe_produk_id tidak valid")
	}

	// Auto-increment urutan per tipe_produk (scoped)
	maxUrutan, err := s.repo.GetMaxUrutanByTipeProduk(ctx, req.TipeProdukID)
	if err != nil {
		return nil, err
	}

	banner := &models.BannerTipeProduk{
		TipeProdukID: tipeProdukUUID,
		Nama:         req.Nama,
		GambarURL:    req.GambarURL,
		Urutan:       maxUrutan + 1,
		IsActive:     true,
	}

	if err := s.repo.Create(ctx, banner); err != nil {
		return nil, err
	}

	return s.FindByID(ctx, banner.ID.String())
}

func (s *bannerTipeProdukService) FindByID(ctx context.Context, id string) (*models.BannerTipeProdukResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}
	return s.toResponse(banner), nil
}

func (s *bannerTipeProdukService) FindAll(ctx context.Context, params *models.BannerTipeProdukFilterRequest, tipeProdukID string) ([]models.BannerTipeProdukSimpleResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	banners, total, err := s.repo.FindAll(ctx, params, tipeProdukID)
	if err != nil {
		return nil, nil, err
	}

	items := []models.BannerTipeProdukSimpleResponse{}
	for _, b := range banners {
		items = append(items, *s.toSimpleResponse(&b))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

// FindAllGrouped - Get all banners grouped by tipe_produk (no pagination)
func (s *bannerTipeProdukService) FindAllGrouped(ctx context.Context, search string) (*models.BannerTipeProdukGroupedResponse, *models.BannerTipeProdukGroupedMeta, error) {
	banners, err := s.repo.FindAllGrouped(ctx, search)
	if err != nil {
		return nil, nil, err
	}

	// Initialize grouped response with empty arrays
	grouped := &models.BannerTipeProdukGroupedResponse{
		PaletLoad:     []models.BannerTipeProdukSimpleResponse{},
		ContainerLoad: []models.BannerTipeProdukSimpleResponse{},
		TruckLoad:     []models.BannerTipeProdukSimpleResponse{},
	}

	// Group by tipe_produk slug
	for _, banner := range banners {
		resp := s.toSimpleResponse(&banner)

		switch banner.TipeProduk.Slug {
		case "palet-load":
			grouped.PaletLoad = append(grouped.PaletLoad, *resp)
		case "container-load":
			grouped.ContainerLoad = append(grouped.ContainerLoad, *resp)
		case "truck-load":
			grouped.TruckLoad = append(grouped.TruckLoad, *resp)
		}
	}

	// Create meta with total by type
	meta := &models.BannerTipeProdukGroupedMeta{
		TotalByType: map[string]int{
			"palet_load":     len(grouped.PaletLoad),
			"container_load": len(grouped.ContainerLoad),
			"truck_load":     len(grouped.TruckLoad),
		},
	}

	return grouped, meta, nil
}

func (s *bannerTipeProdukService) FindByTipeProdukID(ctx context.Context, tipeProdukID string) ([]models.BannerSimpleResponse, error) {
	_, err := uuid.Parse(tipeProdukID)
	if err != nil {
		return nil, errors.New("tipe_produk_id tidak valid")
	}

	banners, err := s.repo.FindByTipeProdukID(ctx, tipeProdukID)
	if err != nil {
		return nil, err
	}

	items := []models.BannerSimpleResponse{}
	for _, b := range banners {
		items = append(items, models.BannerSimpleResponse{
			ID:        b.ID.String(),
			Nama:      b.Nama,
			GambarURL: utils.GetFileURL(b.GambarURL, s.cfg),
			Urutan:    b.Urutan,
		})
	}

	return items, nil
}

func (s *bannerTipeProdukService) Update(ctx context.Context, id string, req *models.UpdateBannerTipeProdukRequest) (*models.BannerTipeProdukResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}

	oldTipeProdukID := banner.TipeProdukID
	oldUrutan := banner.Urutan

	// Check if tipe_produk is being changed
	if req.TipeProdukID != nil {
		tipeProdukUUID, err := uuid.Parse(*req.TipeProdukID)
		if err != nil {
			return nil, errors.New("tipe_produk_id tidak valid")
		}

		// If tipe_produk is actually changing
		if tipeProdukUUID != banner.TipeProdukID {
			// Validate new tipe_produk exists by checking dropdown list
			allTipeProduk, err := s.tipeProdukRepo.GetAllForDropdown(ctx)
			if err != nil {
				return nil, errors.New("gagal validasi tipe produk")
			}

			// Check if the new tipe_produk_id exists
			found := false
			for _, tp := range allTipeProduk {
				if tp.ID == tipeProdukUUID {
					found = true
					break
				}
			}

			if !found {
				return nil, errors.New("tipe produk tidak ditemukan")
			}

			// Get max urutan in new tipe_produk group
			maxUrutan, err := s.repo.GetMaxUrutanByTipeProduk(ctx, tipeProdukUUID.String())
			if err != nil {
				return nil, err
			}

			// Update banner with new tipe_produk and place at end
			banner.TipeProdukID = tipeProdukUUID
			banner.Urutan = maxUrutan + 1
			// Clear preloaded relation to avoid GORM using old FK
			banner.TipeProduk = models.TipeProduk{}

			// Reorder old tipe_produk group to fill gap
			if err := s.repo.ReorderAfterDeleteScoped(ctx, oldTipeProdukID.String(), oldUrutan); err != nil {
				return nil, err
			}
		}
	}

	if req.Nama != nil {
		banner.Nama = *req.Nama
	}
	if req.GambarURL != nil {
		banner.GambarURL = *req.GambarURL
	}
	if req.IsActive != nil {
		banner.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, banner); err != nil {
		return nil, err
	}

	return s.FindByID(ctx, id)
}

func (s *bannerTipeProdukService) UpdateWithFile(ctx context.Context, id string, req *models.UpdateBannerTipeProdukRequest, newGambarURL *string) (*models.BannerTipeProdukResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}

	// Store old gambar URL for deletion if update successful
	oldGambarURL := banner.GambarURL
	oldTipeProdukID := banner.TipeProdukID
	oldUrutan := banner.Urutan

	// Check if tipe_produk is being changed
	if req.TipeProdukID != nil {
		tipeProdukUUID, err := uuid.Parse(*req.TipeProdukID)
		if err != nil {
			return nil, errors.New("tipe_produk_id tidak valid")
		}

		// If tipe_produk is actually changing
		if tipeProdukUUID != banner.TipeProdukID {
			// Validate new tipe_produk exists by checking dropdown list
			allTipeProduk, err := s.tipeProdukRepo.GetAllForDropdown(ctx)
			if err != nil {
				return nil, errors.New("gagal validasi tipe produk")
			}

			// Check if the new tipe_produk_id exists
			found := false
			for _, tp := range allTipeProduk {
				if tp.ID == tipeProdukUUID {
					found = true
					break
				}
			}

			if !found {
				return nil, errors.New("tipe produk tidak ditemukan")
			}

			// Get max urutan in new tipe_produk group
			maxUrutan, err := s.repo.GetMaxUrutanByTipeProduk(ctx, tipeProdukUUID.String())
			if err != nil {
				return nil, err
			}

			// Update banner with new tipe_produk and place at end
			banner.TipeProdukID = tipeProdukUUID
			banner.Urutan = maxUrutan + 1
			// Clear preloaded relation to avoid GORM using old FK
			banner.TipeProduk = models.TipeProduk{}

			// Reorder old tipe_produk group to fill gap
			if err := s.repo.ReorderAfterDeleteScoped(ctx, oldTipeProdukID.String(), oldUrutan); err != nil {
				return nil, err
			}
		}
	}

	if req.Nama != nil {
		banner.Nama = *req.Nama
	}
	if req.GambarURL != nil {
		banner.GambarURL = *req.GambarURL
	}
	if req.IsActive != nil {
		banner.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, banner); err != nil {
		return nil, err
	}

	// Delete old file if new file was uploaded
	if newGambarURL != nil && oldGambarURL != "" {
		utils.DeleteFile(oldGambarURL, s.cfg)
	}

	return s.FindByID(ctx, id)
}

func (s *bannerTipeProdukService) Delete(ctx context.Context, id string) error {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("banner tidak ditemukan")
	}

	deletedUrutan := banner.Urutan
	tipeProdukID := banner.TipeProdukID.String()

	// Soft delete banner
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Reorder remaining items within same tipe_produk (scoped)
	return s.repo.ReorderAfterDeleteScoped(ctx, tipeProdukID, deletedUrutan)
}

func (s *bannerTipeProdukService) DeleteWithFile(ctx context.Context, id string) error {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("banner tidak ditemukan")
	}

	deletedUrutan := banner.Urutan
	tipeProdukID := banner.TipeProdukID.String()

	// Delete from database first
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Delete file from storage
	if banner.GambarURL != "" {
		utils.DeleteFile(banner.GambarURL, s.cfg)
	}

	// Reorder remaining items within same tipe_produk (scoped)
	return s.repo.ReorderAfterDeleteScoped(ctx, tipeProdukID, deletedUrutan)
}

func (s *bannerTipeProdukService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	banner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("banner tidak ditemukan")
	}

	banner.IsActive = !banner.IsActive
	if err := s.repo.Update(ctx, banner); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       banner.ID.String(),
		IsActive: banner.IsActive,
	}, nil
}

func (s *bannerTipeProdukService) Reorder(ctx context.Context, req *models.ReorderRequest) error {
	if len(req.Items) == 0 {
		return errors.New("items tidak boleh kosong")
	}

	err := s.repo.UpdateOrder(ctx, req.Items)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("salah satu atau lebih banner tidak ditemukan")
		}
		return err
	}

	return nil
}

func (s *bannerTipeProdukService) toResponse(b *models.BannerTipeProduk) *models.BannerTipeProdukResponse {
	return &models.BannerTipeProdukResponse{
		ID: b.ID.String(),
		TipeProduk: models.BannerTipeProdukTipeInfo{
			ID:   b.TipeProduk.ID.String(),
			Nama: b.TipeProduk.Nama,
			Slug: b.TipeProduk.Slug,
		},
		Nama:      b.Nama,
		GambarURL: utils.GetFileURL(b.GambarURL, s.cfg),
		Urutan:    b.Urutan,
		IsActive:  b.IsActive,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

func (s *bannerTipeProdukService) toSimpleResponse(b *models.BannerTipeProduk) *models.BannerTipeProdukSimpleResponse {
	return &models.BannerTipeProdukSimpleResponse{
		ID: b.ID.String(),
		TipeProduk: models.BannerTipeProdukSimpleInfo{
			// ID:   b.TipeProduk.ID.String(),
			Nama: b.TipeProduk.Nama,
			// Slug: b.TipeProduk.Slug,
		},
		Nama:      b.Nama,
		GambarURL: utils.GetFileURL(b.GambarURL, s.cfg),
		Urutan:    b.Urutan,
		IsActive:  b.IsActive,
		// CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
