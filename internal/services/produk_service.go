package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"strings"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/constants"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProdukService interface {
	Create(ctx context.Context, req *models.CreateProdukRequest) (*models.ProdukDetailResponse, error)
	CreateWithFiles(ctx context.Context, req *models.CreateProdukRequest, isActive bool, gambarFiles, dokumenFiles []*multipart.FileHeader, dokumenNama []string) (*models.ProdukDetailResponse, error)
	FindByID(ctx context.Context, id string) (*models.ProdukDetailResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.ProdukDetailResponse, error)
	FindAll(ctx context.Context, params *models.ProdukFilterRequest) ([]models.ProdukPanelListResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateProdukRequest) (*models.ProdukDetailResponse, error)
	UpdateWithFiles(ctx context.Context, id string, req *models.UpdateProdukRequest, dokumenFiles []*multipart.FileHeader, dokumenNama []string) (*models.ProdukDetailResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	UpdateStock(ctx context.Context, id string, req *models.UpdateStockRequest) (*models.ProdukDetailResponse, error)
}

type produkService struct {
	repo           repositories.ProdukRepository
	gambarRepo     repositories.ProdukGambarRepository
	dokumenRepo    repositories.ProdukDokumenRepository
	warehouseRepo  repositories.WarehouseRepository
	tipeProdukRepo repositories.TipeProdukRepository
	cfg            *config.Config
	db             *gorm.DB
}

func NewProdukService(
	repo repositories.ProdukRepository,
	gambarRepo repositories.ProdukGambarRepository,
	dokumenRepo repositories.ProdukDokumenRepository,
	warehouseRepo repositories.WarehouseRepository,
	tipeProdukRepo repositories.TipeProdukRepository,
	cfg *config.Config,
	db *gorm.DB,
) ProdukService {
	return &produkService{
		repo:           repo,
		gambarRepo:     gambarRepo,
		dokumenRepo:    dokumenRepo,
		warehouseRepo:  warehouseRepo,
		tipeProdukRepo: tipeProdukRepo,
		cfg:            cfg,
		db:             db,
	}
}

func (s *produkService) Create(ctx context.Context, req *models.CreateProdukRequest) (*models.ProdukDetailResponse, error) {
	// This method is deprecated - use CreateWithFiles instead
	return nil, errors.New("use CreateWithFiles method for multipart form-data")
}

func (s *produkService) CreateWithFiles(
	ctx context.Context,
	req *models.CreateProdukRequest,
	isActive bool,
	gambarFiles, dokumenFiles []*multipart.FileHeader,
	dokumenNama []string,
) (*models.ProdukDetailResponse, error) {
	// Generate slug_id
	var slugID *string
	if req.SlugID != nil && *req.SlugID != "" {
		s := *req.SlugID
		slugID = &s
	} else {
		s := utils.GenerateSlug(req.NamaID)
		slugID = &s
	}

	// Generate slug_en
	var slugEN *string
	if req.SlugEN != nil && *req.SlugEN != "" {
		s := *req.SlugEN
		slugEN = &s
	} else if req.NamaEN != "" {
		s := utils.GenerateSlug(req.NamaEN)
		slugEN = &s
	}

	if req.IDCargo != nil && *req.IDCargo != "" {
		exists, _ := s.repo.ExistsByIDCargo(ctx, *req.IDCargo, nil)
		if exists {
			return nil, errors.New("ID Cargo sudah digunakan oleh produk lain")
		}
	}

	kategoriID, _ := uuid.Parse(req.KategoriID)
	kondisiID, _ := uuid.Parse(req.KondisiID)
	kondisiPaketID, _ := uuid.Parse(req.KondisiPaketID)

	// Auto-set warehouse_id by querying slug
	warehouse, err := s.warehouseRepo.FindBySlug(ctx, constants.DefaultWarehouseSlug)
	if err != nil {
		return nil, fmt.Errorf("warehouse default '%s' tidak ditemukan", constants.DefaultWarehouseSlug)
	}
	warehouseID := warehouse.ID

	// Auto-set tipe_produk_id by querying slug
	tipeProduk, err := s.tipeProdukRepo.FindBySlug(ctx, constants.DefaultTipeProdukSlug)
	if err != nil {
		return nil, fmt.Errorf("tipe produk default '%s' tidak ditemukan", constants.DefaultTipeProdukSlug)
	}
	tipeProdukID := tipeProduk.ID

	produk := &models.Produk{
		NamaID:             req.NamaID,
		NamaEN:             req.NamaEN,
		Slug:               *slugID,
		SlugID:             slugID,
		SlugEN:             slugEN,
		IDCargo:            req.IDCargo,
		ReferenceID:        req.ReferenceID,
		KategoriID:         kategoriID,
		KondisiID:          kondisiID,
		KondisiPaketID:     kondisiPaketID,
		WarehouseID:        warehouseID,
		TipeProdukID:       tipeProdukID,
		HargaSebelumDiskon: req.HargaSebelumDiskon,
		HargaSesudahDiskon: req.HargaSesudahDiskon,
		Quantity:           req.Quantity,
		Discrepancy:        req.Discrepancy,
		Panjang:            req.Panjang,
		Lebar:              req.Lebar,
		Tinggi:             req.Tinggi,
		Berat:              req.Berat,
		IsActive:           isActive,
	}

	if req.SumberID != nil {
		sumberID, _ := uuid.Parse(*req.SumberID)
		produk.SumberID = &sumberID
	}

	// Begin transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create produk record
	if err := tx.Create(produk).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. Create produk_merek relations (many-to-many)
	merekIDs := req.GetMerekIDs()
	if len(merekIDs) > 0 {
		produkMereks := make([]models.ProdukMerek, len(merekIDs))
		for i, merekID := range merekIDs {
			produkMereks[i] = models.ProdukMerek{
				ProdukID: produk.ID,
				MerekID:  merekID,
			}
		}

		// Batch insert relation records
		if err := tx.Create(&produkMereks).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("gagal membuat relasi merek: %w", err)
		}
	}
	// 3. Upload and create gambar records
	produkDir := fmt.Sprintf("products/%s", produk.ID.String())
	for i, file := range gambarFiles {
		// Upload to storage
		relativePath, err := utils.SaveUploadedFile(file, produkDir, s.cfg)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("gagal upload gambar: %w", err)
		}

		gambar := &models.ProdukGambar{
			ProdukID:  produk.ID,
			GambarURL: relativePath,
			Urutan:    i + 1,
			IsPrimary: i == 0, // Gambar pertama otomatis jadi primary
		}

		if err := tx.Create(gambar).Error; err != nil {
			tx.Rollback()
			// Cleanup uploaded file
			utils.DeleteFile(relativePath, s.cfg)
			return nil, err
		}
	}

	// 4. Upload and create dokumen records
	if len(dokumenFiles) > 0 {
		dokumenDir := fmt.Sprintf("documents/%s", produk.ID.String())
		for i, file := range dokumenFiles {
			// Upload to storage
			relativePath, err := utils.SaveUploadedFile(file, dokumenDir, s.cfg)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("gagal upload dokumen: %w", err)
			}

			// Get document name
			nama := file.Filename
			if i < len(dokumenNama) && dokumenNama[i] != "" {
				nama = dokumenNama[i]
			}

			dokumen := &models.ProdukDokumen{
				ProdukID:    produk.ID,
				NamaDokumen: nama,
				FileURL:     relativePath,
				TipeFile:    "pdf",
				UkuranFile:  intPtr(int(file.Size)),
			}

			if err := tx.Create(dokumen).Error; err != nil {
				tx.Rollback()
				// Cleanup uploaded file
				utils.DeleteFile(relativePath, s.cfg)
				return nil, err
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload with relations
	return s.FindByID(ctx, produk.ID.String())
}

func (s *produkService) FindByID(ctx context.Context, id string) (*models.ProdukDetailResponse, error) {
	produk, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	return s.toDetailResponse(produk), nil
}

func (s *produkService) FindBySlug(ctx context.Context, slug string) (*models.ProdukDetailResponse, error) {
	produk, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	return s.toDetailResponse(produk), nil
}

func (s *produkService) FindAll(ctx context.Context, params *models.ProdukFilterRequest) ([]models.ProdukPanelListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	produks, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	items := []models.ProdukPanelListResponse{}
	for _, p := range produks {
		items = append(items, *s.toPanelListResponse(&p))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *produkService) Update(ctx context.Context, id string, req *models.UpdateProdukRequest) (*models.ProdukDetailResponse, error) {
	produk, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	if req.NamaID != nil {
		produk.NamaID = *req.NamaID
		// Regenerate slug_id from new nama_id (unless manually provided)
		if req.SlugID == nil {
			s := utils.GenerateSlug(*req.NamaID)
			produk.SlugID = &s
			produk.Slug = s // backward compat
		}
	}
	if req.SlugID != nil && *req.SlugID != "" {
		produk.SlugID = req.SlugID
		produk.Slug = *req.SlugID // backward compat
	}

	if req.NamaEN != nil {
		produk.NamaEN = *req.NamaEN
		// Regenerate slug_en from new nama_en (unless manually provided)
		if req.SlugEN == nil && *req.NamaEN != "" {
			s := utils.GenerateSlug(*req.NamaEN)
			produk.SlugEN = &s
		}
	}
	if req.SlugEN != nil && *req.SlugEN != "" {
		produk.SlugEN = req.SlugEN
	}

	if req.IDCargo != nil {
		if *req.IDCargo != "" {
			exists, _ := s.repo.ExistsByIDCargo(ctx, *req.IDCargo, &id)
			if exists {
				return nil, errors.New("id_cargo sudah digunakan")
			}
		}
		produk.IDCargo = req.IDCargo
	}

	if req.ReferenceID != nil {
		produk.ReferenceID = req.ReferenceID
	}

	if req.KategoriID != nil {
		kategoriID, _ := uuid.Parse(*req.KategoriID)
		produk.KategoriID = kategoriID
	}
	if req.KondisiID != nil {
		kondisiID, _ := uuid.Parse(*req.KondisiID)
		produk.KondisiID = kondisiID
	}
	if req.KondisiPaketID != nil {
		kondisiPaketID, _ := uuid.Parse(*req.KondisiPaketID)
		produk.KondisiPaketID = kondisiPaketID
	}
	if req.SumberID != nil {
		sumberID, _ := uuid.Parse(*req.SumberID)
		produk.SumberID = &sumberID
	}
	// Note: warehouse_id, tipe_produk_id, persentase_diskon are auto-managed
	// Use dedicated endpoints if these need to be changed
	if req.HargaSebelumDiskon != nil {
		produk.HargaSebelumDiskon = *req.HargaSebelumDiskon
	}
	if req.HargaSesudahDiskon != nil {
		produk.HargaSesudahDiskon = *req.HargaSesudahDiskon
	}
	if req.Quantity != nil {
		produk.Quantity = *req.Quantity
	}
	if req.Discrepancy != nil {
		produk.Discrepancy = req.Discrepancy
	}
	if req.Panjang != nil {
		produk.Panjang = *req.Panjang
	}
	if req.Lebar != nil {
		produk.Lebar = *req.Lebar
	}
	if req.Tinggi != nil {
		produk.Tinggi = *req.Tinggi
	}
	if req.Berat != nil {
		produk.Berat = *req.Berat
	}
	if req.IsActive != nil {
		// Parse string to bool: "true"/"1" = true, "false"/"0" = false
		val := strings.ToLower(strings.TrimSpace(*req.IsActive))
		produk.IsActive = (val == "true" || val == "1")
	}

	// Begin transaction for update
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update produk record
	if err := tx.Save(produk).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update merek relations if merek_ids provided
	if merekIDs, shouldUpdate := req.GetMerekIDs(); shouldUpdate {
		produkID := produk.ID

		// Delete existing relations
		if err := tx.Where("produk_id = ?", produkID).Delete(&models.ProdukMerek{}).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("gagal menghapus relasi merek lama: %w", err)
		}

		// Create new relations if any
		if len(merekIDs) > 0 {
			produkMereks := make([]models.ProdukMerek, len(merekIDs))
			for i, merekID := range merekIDs {
				produkMereks[i] = models.ProdukMerek{
					ProdukID: produkID,
					MerekID:  merekID,
				}
			}

			if err := tx.Create(&produkMereks).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("gagal membuat relasi merek baru: %w", err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return s.FindByID(ctx, id)
}

func (s *produkService) UpdateWithFiles(ctx context.Context, id string, req *models.UpdateProdukRequest, dokumenFiles []*multipart.FileHeader, dokumenNama []string) (*models.ProdukDetailResponse, error) {
	// First do the regular product update
	result, err := s.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// If no new dokumen files, return the result as-is
	if len(dokumenFiles) == 0 {
		return result, nil
	}

	// Validate max 5 dokumen
	if len(dokumenFiles) > 5 {
		return nil, errors.New("maksimal 5 dokumen")
	}

	// Begin transaction for dokumen replacement
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Fetch existing dokumen to delete files from storage
	existingDokumens, err := s.dokumenRepo.FindByProdukID(ctx, id)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal mengambil dokumen lama: %w", err)
	}

	// Delete all existing dokumen files from storage
	for _, dok := range existingDokumens {
		if err := utils.DeleteFile(dok.FileURL, s.cfg); err != nil {
			fmt.Printf("Warning: gagal menghapus file dokumen %s: %v\n", dok.FileURL, err)
		}
	}

	// Delete all existing dokumen records from DB
	produkUUID, _ := uuid.Parse(id)
	if err := tx.Where("produk_id = ?", produkUUID).Delete(&models.ProdukDokumen{}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal menghapus dokumen lama: %w", err)
	}

	// Upload and insert new dokumen
	dokumenDir := fmt.Sprintf("documents/%s", id)
	for i, file := range dokumenFiles {
		relativePath, err := utils.SaveUploadedFile(file, dokumenDir, s.cfg)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("gagal upload dokumen: %w", err)
		}

		nama := file.Filename
		if i < len(dokumenNama) && dokumenNama[i] != "" {
			nama = dokumenNama[i]
		}

		dokumen := &models.ProdukDokumen{
			ProdukID:    produkUUID,
			NamaDokumen: nama,
			FileURL:     relativePath,
			TipeFile:    "pdf",
			UkuranFile:  intPtr(int(file.Size)),
		}

		if err := tx.Create(dokumen).Error; err != nil {
			tx.Rollback()
			utils.DeleteFile(relativePath, s.cfg)
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return s.FindByID(ctx, id)
}

func (s *produkService) Delete(ctx context.Context, id string) error {
	produk, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}
	return s.repo.Delete(ctx, produk)
}

func (s *produkService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	produk, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	produk.IsActive = !produk.IsActive
	if err := s.repo.Update(ctx, produk); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:        produk.ID.String(),
		IsActive:  produk.IsActive,
		UpdatedAt: produk.UpdatedAt,
	}, nil
}

func (s *produkService) UpdateStock(ctx context.Context, id string, req *models.UpdateStockRequest) (*models.ProdukDetailResponse, error) {
	produk, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	produk.Quantity = req.Quantity
	if err := s.repo.Update(ctx, produk); err != nil {
		return nil, err
	}

	return s.FindByID(ctx, id)
}

func (s *produkService) toListResponse(p *models.Produk) *models.ProdukListResponse {
	resp := &models.ProdukListResponse{
		ID:     p.ID.String(),
		NamaID: p.NamaID,
		NamaEN: p.NamaEN,
		// Slug:    p.Slug,
		// IDCargo: p.IDCargo,
		Kategori: models.SimpleProdukRelationInfo{
			// ID:   p.Kategori.ID.String(),
			Nama: p.Kategori.GetNama().ID,
			// Slug: p.Kategori.Slug,
		},
		Kondisi: models.SimpleProdukRelationInfo{
			// ID:   p.Kondisi.ID.String(),
			Nama: p.Kondisi.GetNama().ID,
			// Slug: p.Kondisi.Slug,
		},
		// KondisiPaket: models.SimpleProdukRelationInfo{
		// 	// ID:   p.KondisiPaket.ID.String(),
		// 	Nama: p.KondisiPaket.GetNama().ID,
		// 	// Slug: p.KondisiPaket.Slug,
		// },
		Warehouse: models.SimpleProdukWarehouseInfo{
			// ID:   p.Warehouse.ID.String(),
			Nama: p.Warehouse.Nama,
			// Slug: p.Warehouse.Slug,
		},
		TipeProduk: models.SimpleProdukRelationInfo{
			// ID:   p.TipeProduk.ID.String(),
			Nama: p.TipeProduk.Nama,
			// Slug: p.TipeProduk.Slug,
		},
		// HargaSebelumDiskon: p.HargaSebelumDiskon,
		// PersentaseDiskon:   p.PersentaseDiskon,
		HargaSesudahDiskon: p.HargaSesudahDiskon,
		Quantity:           p.Quantity,
		QuantityTerjual:    p.QuantityTerjual,
		Berat:              p.Berat,
		IsActive:           p.IsActive,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}

	// Map multiple mereks
	resp.Mereks = make([]models.SimpleProdukRelationInfo, 0, len(p.Mereks))
	for _, merek := range p.Mereks {
		resp.Mereks = append(resp.Mereks, models.SimpleProdukRelationInfo{
			// ID:   merek.ID.String(),
			Nama: merek.GetNama().ID,
			// Slug: merek.Slug,
		})
	}

	// if p.Sumber != nil {
	// 	resp.Sumber = &models.SimpleProdukRelationInfo{
	// 		// ID:   p.Sumber.ID.String(),
	// 		Nama: p.Sumber.GetNama().ID,
	// 		// Slug: p.Sumber.Slug,
	// 	}
	// }

	// Get primary image
	if len(p.Gambar) > 0 {
		fullURL := utils.GetFileURL(p.Gambar[0].GambarURL, s.cfg)
		resp.GambarUtama = &fullURL
	}

	return resp
}

// toPanelListResponse converts Produk to simplified ProdukPanelListResponse for admin panel
func (s *produkService) toPanelListResponse(p *models.Produk) *models.ProdukPanelListResponse {
	resp := &models.ProdukPanelListResponse{
		ID:      p.ID.String(),
		NamaID:  p.NamaID,
		NamaEN:  p.NamaEN,
		IDCargo: p.IDCargo,
		Status:  p.IsActive,
	}

	// Get primary/first image - prioritize is_primary, fallback to first by urutan
	if len(p.Gambar) > 0 {
		var selectedGambar *models.ProdukGambar
		// Find primary image first
		for i := range p.Gambar {
			if p.Gambar[i].IsPrimary {
				selectedGambar = &p.Gambar[i]
				break
			}
		}
		// If no primary, use first one
		if selectedGambar == nil {
			selectedGambar = &p.Gambar[0]
		}
		fullURL := utils.GetFileURL(selectedGambar.GambarURL, s.cfg)
		resp.GambarUtama = &fullURL
	}

	// Get first PDF document
	for _, dok := range p.Dokumen {
		if strings.HasSuffix(strings.ToLower(dok.FileURL), ".pdf") {
			fullURL := utils.GetFileURL(dok.FileURL, s.cfg)
			resp.FilePDF = &fullURL
			break
		}
	}

	return resp
}

func (s *produkService) toDetailResponse(p *models.Produk) *models.ProdukDetailResponse {
	resp := &models.ProdukDetailResponse{
		ID:          p.ID.String(),
		NamaID:      p.NamaID,
		NamaEN:      p.NamaEN,
		SlugID:      p.SlugID,
		SlugEN:      p.SlugEN,
		IDCargo:     p.IDCargo,
		ReferenceID: p.ReferenceID,
		Kategori: models.SimpleProdukRelationInfo{
			ID:   p.Kategori.ID.String(),
			Nama: p.Kategori.GetNama().ID,
			// Slug: p.Kategori.Slug,
		},
		Kondisi: models.SimpleProdukRelationInfo{
			ID:   p.Kondisi.ID.String(),
			Nama: p.Kondisi.GetNama().ID,
			// Slug: p.Kondisi.Slug,
		},
		KondisiPaket: models.SimpleProdukRelationInfo{
			ID:   p.KondisiPaket.ID.String(),
			Nama: p.KondisiPaket.GetNama().ID,
			// Slug: p.KondisiPaket.Slug,
		},
		Warehouse: models.ProdukWarehouseInfo{
			ID:   p.Warehouse.ID.String(),
			Nama: p.Warehouse.Nama,
			// Slug: p.Warehouse.Slug,
			Kota: p.Warehouse.Kota,
		},
		TipeProduk: models.SimpleProdukRelationInfo{
			ID:   p.TipeProduk.ID.String(),
			Nama: p.TipeProduk.Nama,
			// Slug: p.TipeProduk.Slug,
		},
		HargaSebelumDiskon: p.HargaSebelumDiskon,
		HargaSesudahDiskon: p.HargaSesudahDiskon,
		Quantity:           p.Quantity,
		QuantityTerjual:    p.QuantityTerjual,
		Discrepancy:        p.Discrepancy,
		Panjang:            p.Panjang,
		Lebar:              p.Lebar,
		Tinggi:             p.Tinggi,
		Berat:              p.Berat,
		BeratVolumetrik:    s.calculateBeratVolumetrik(p.Panjang, p.Lebar, p.Tinggi),
		IsActive:           p.IsActive,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}

	// Map multiple mereks
	resp.Mereks = make([]models.SimpleProdukRelationInfo, 0, len(p.Mereks))
	for _, merek := range p.Mereks {
		resp.Mereks = append(resp.Mereks, models.SimpleProdukRelationInfo{
			ID:   merek.ID.String(),
			Nama: merek.GetNama().ID,
			// Slug: merek.Slug,
		})
	}

	if p.Sumber != nil {
		resp.Sumber = &models.SimpleProdukRelationInfo{
			ID:   p.Sumber.ID.String(),
			Nama: p.Sumber.GetNama().ID,
			// Slug: p.Sumber.Slug,
		}
	}

	// Map gambar
	for _, g := range p.Gambar {
		resp.Gambar = append(resp.Gambar, models.ProdukGambarResponse{
			ID:        g.ID.String(),
			GambarURL: utils.GetFileURL(g.GambarURL, s.cfg),
			Urutan:    g.Urutan,
			IsPrimary: g.IsPrimary,
		})
	}

	// Map dokumen
	for _, d := range p.Dokumen {
		resp.Dokumen = append(resp.Dokumen, models.ProdukDokumenResponse{
			ID:          d.ID.String(),
			NamaDokumen: d.NamaDokumen,
			FileURL:     utils.GetFileURL(d.FileURL, s.cfg),
			TipeFile:    d.TipeFile,
			UkuranFile:  d.UkuranFile,
		})
	}

	return resp
}

// calculateBeratVolumetrik menghitung berat volumetrik berdasarkan dimensi produk
// Formula: (Panjang × Lebar × Tinggi) / 6000
// Hasil dalam kg, dibulatkan 2 desimal
func (s *produkService) calculateBeratVolumetrik(panjang, lebar, tinggi float64) float64 {
	if panjang == 0 || lebar == 0 || tinggi == 0 {
		return 0
	}
	// Formula standar: (P × L × T) / 6000
	// Divisor 6000 untuk konversi cm³ ke kg
	volumetrik := (panjang * lebar * tinggi) / 6000
	return math.Round(volumetrik*100) / 100
}

// Helper function
func intPtr(i int) *int {
	return &i
}
