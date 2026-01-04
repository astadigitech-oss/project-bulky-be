package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
)

type ProdukService interface {
	Create(ctx context.Context, req *models.CreateProdukRequest) (*models.ProdukDetailResponse, error)
	FindByID(ctx context.Context, id string) (*models.ProdukDetailResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.ProdukDetailResponse, error)
	FindAll(ctx context.Context, params *models.ProdukFilterRequest) ([]models.ProdukListResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateProdukRequest) (*models.ProdukDetailResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	UpdateStock(ctx context.Context, id string, req *models.UpdateStockRequest) (*models.ProdukDetailResponse, error)
}

type produkService struct {
	repo       repositories.ProdukRepository
	gambarRepo repositories.ProdukGambarRepository
}

func NewProdukService(repo repositories.ProdukRepository, gambarRepo repositories.ProdukGambarRepository) ProdukService {
	return &produkService{repo: repo, gambarRepo: gambarRepo}
}

func (s *produkService) Create(ctx context.Context, req *models.CreateProdukRequest) (*models.ProdukDetailResponse, error) {
	slug := utils.GenerateSlug(req.Nama)

	exists, _ := s.repo.ExistsBySlug(ctx, slug, nil)
	if exists {
		return nil, errors.New("produk dengan nama tersebut sudah ada")
	}

	if req.IDCargo != nil && *req.IDCargo != "" {
		exists, _ := s.repo.ExistsByIDCargo(ctx, *req.IDCargo, nil)
		if exists {
			return nil, errors.New("id_cargo sudah digunakan")
		}
	}

	kategoriID, _ := uuid.Parse(req.KategoriID)
	kondisiID, _ := uuid.Parse(req.KondisiID)
	kondisiPaketID, _ := uuid.Parse(req.KondisiPaketID)
	warehouseID, _ := uuid.Parse(req.WarehouseID)
	tipeProdukID, _ := uuid.Parse(req.TipeProdukID)

	produk := &models.Produk{
		Nama:               req.Nama,
		Slug:               slug,
		IDCargo:            req.IDCargo,
		KategoriID:         kategoriID,
		KondisiID:          kondisiID,
		KondisiPaketID:     kondisiPaketID,
		WarehouseID:        warehouseID,
		TipeProdukID:       tipeProdukID,
		HargaSebelumDiskon: req.HargaSebelumDiskon,
		PersentaseDiskon:   req.PersentaseDiskon,
		HargaSesudahDiskon: req.HargaSesudahDiskon,
		Quantity:           req.Quantity,
		Discrepancy:        req.Discrepancy,
		IsActive:           true,
	}

	if req.MerekID != nil {
		merekID, _ := uuid.Parse(*req.MerekID)
		produk.MerekID = &merekID
	}
	if req.SumberID != nil {
		sumberID, _ := uuid.Parse(*req.SumberID)
		produk.SumberID = &sumberID
	}

	if err := s.repo.Create(ctx, produk); err != nil {
		return nil, err
	}

	// Create gambar
	for i, url := range req.GambarURLs {
		gambar := &models.ProdukGambar{
			ProdukID:  produk.ID,
			GambarURL: url,
			Urutan:    i,
			IsPrimary: i == req.GambarUtamaIndex,
		}
		if err := s.gambarRepo.Create(ctx, gambar); err != nil {
			return nil, err
		}
	}

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

func (s *produkService) FindAll(ctx context.Context, params *models.ProdukFilterRequest) ([]models.ProdukListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	produks, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.ProdukListResponse
	for _, p := range produks {
		items = append(items, *s.toListResponse(&p))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *produkService) Update(ctx context.Context, id string, req *models.UpdateProdukRequest) (*models.ProdukDetailResponse, error) {
	produk, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	if req.Nama != nil {
		newSlug := utils.GenerateSlug(*req.Nama)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("produk dengan nama tersebut sudah ada")
		}
		produk.Nama = *req.Nama
		produk.Slug = newSlug
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

	if req.KategoriID != nil {
		kategoriID, _ := uuid.Parse(*req.KategoriID)
		produk.KategoriID = kategoriID
	}
	if req.MerekID != nil {
		merekID, _ := uuid.Parse(*req.MerekID)
		produk.MerekID = &merekID
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
	if req.WarehouseID != nil {
		warehouseID, _ := uuid.Parse(*req.WarehouseID)
		produk.WarehouseID = warehouseID
	}
	if req.TipeProdukID != nil {
		tipeProdukID, _ := uuid.Parse(*req.TipeProdukID)
		produk.TipeProdukID = tipeProdukID
	}
	if req.HargaSebelumDiskon != nil {
		produk.HargaSebelumDiskon = *req.HargaSebelumDiskon
	}
	if req.PersentaseDiskon != nil {
		produk.PersentaseDiskon = *req.PersentaseDiskon
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
	if req.IsActive != nil {
		produk.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, produk); err != nil {
		return nil, err
	}

	return s.FindByID(ctx, id)
}

func (s *produkService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}
	return s.repo.Delete(ctx, id)
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
		ID:       produk.ID.String(),
		IsActive: produk.IsActive,
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
		ID:      p.ID.String(),
		Nama:    p.Nama,
		Slug:    p.Slug,
		IDCargo: p.IDCargo,
		Kategori: models.ProdukRelationInfo{
			ID:   p.Kategori.ID.String(),
			Nama: p.Kategori.Nama,
			Slug: p.Kategori.Slug,
		},
		Kondisi: models.ProdukRelationInfo{
			ID:   p.Kondisi.ID.String(),
			Nama: p.Kondisi.Nama,
			Slug: p.Kondisi.Slug,
		},
		KondisiPaket: models.ProdukRelationInfo{
			ID:   p.KondisiPaket.ID.String(),
			Nama: p.KondisiPaket.Nama,
			Slug: p.KondisiPaket.Slug,
		},
		Warehouse: models.ProdukRelationInfo{
			ID:   p.Warehouse.ID.String(),
			Nama: p.Warehouse.Nama,
			Slug: p.Warehouse.Slug,
		},
		TipeProduk: models.ProdukRelationInfo{
			ID:   p.TipeProduk.ID.String(),
			Nama: p.TipeProduk.Nama,
			Slug: p.TipeProduk.Slug,
		},
		HargaSebelumDiskon: p.HargaSebelumDiskon,
		PersentaseDiskon:   p.PersentaseDiskon,
		HargaSesudahDiskon: p.HargaSesudahDiskon,
		Quantity:           p.Quantity,
		QuantityTerjual:    p.QuantityTerjual,
		IsActive:           p.IsActive,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}

	if p.Merek != nil {
		resp.Merek = &models.ProdukRelationInfo{
			ID:   p.Merek.ID.String(),
			Nama: p.Merek.Nama,
			Slug: p.Merek.Slug,
		}
	}
	if p.Sumber != nil {
		resp.Sumber = &models.ProdukRelationInfo{
			ID:   p.Sumber.ID.String(),
			Nama: p.Sumber.Nama,
			Slug: p.Sumber.Slug,
		}
	}

	// Get primary image
	if len(p.Gambar) > 0 {
		resp.GambarUtama = &p.Gambar[0].GambarURL
	}

	return resp
}

func (s *produkService) toDetailResponse(p *models.Produk) *models.ProdukDetailResponse {
	resp := &models.ProdukDetailResponse{
		ID:      p.ID.String(),
		Nama:    p.Nama,
		Slug:    p.Slug,
		IDCargo: p.IDCargo,
		Kategori: models.ProdukRelationInfo{
			ID:   p.Kategori.ID.String(),
			Nama: p.Kategori.Nama,
			Slug: p.Kategori.Slug,
		},
		Kondisi: models.ProdukRelationInfo{
			ID:   p.Kondisi.ID.String(),
			Nama: p.Kondisi.Nama,
			Slug: p.Kondisi.Slug,
		},
		KondisiPaket: models.ProdukRelationInfo{
			ID:   p.KondisiPaket.ID.String(),
			Nama: p.KondisiPaket.Nama,
			Slug: p.KondisiPaket.Slug,
		},
		Warehouse: models.ProdukWarehouseInfo{
			ID:   p.Warehouse.ID.String(),
			Nama: p.Warehouse.Nama,
			Slug: p.Warehouse.Slug,
			Kota: p.Warehouse.Kota,
		},
		TipeProduk: models.ProdukRelationInfo{
			ID:   p.TipeProduk.ID.String(),
			Nama: p.TipeProduk.Nama,
			Slug: p.TipeProduk.Slug,
		},
		HargaSebelumDiskon: p.HargaSebelumDiskon,
		PersentaseDiskon:   p.PersentaseDiskon,
		HargaSesudahDiskon: p.HargaSesudahDiskon,
		Quantity:           p.Quantity,
		QuantityTerjual:    p.QuantityTerjual,
		Discrepancy:        p.Discrepancy,
		IsActive:           p.IsActive,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}

	if p.Merek != nil {
		resp.Merek = &models.ProdukRelationInfo{
			ID:   p.Merek.ID.String(),
			Nama: p.Merek.Nama,
			Slug: p.Merek.Slug,
		}
	}
	if p.Sumber != nil {
		resp.Sumber = &models.ProdukRelationInfo{
			ID:   p.Sumber.ID.String(),
			Nama: p.Sumber.Nama,
			Slug: p.Sumber.Slug,
		}
	}

	// Map gambar
	for _, g := range p.Gambar {
		resp.Gambar = append(resp.Gambar, models.ProdukGambarResponse{
			ID:        g.ID.String(),
			GambarURL: g.GambarURL,
			Urutan:    g.Urutan,
			IsPrimary: g.IsPrimary,
		})
	}

	// Map dokumen
	for _, d := range p.Dokumen {
		resp.Dokumen = append(resp.Dokumen, models.ProdukDokumenResponse{
			ID:          d.ID.String(),
			NamaDokumen: d.NamaDokumen,
			FileURL:     d.FileURL,
			TipeFile:    d.TipeFile,
			UkuranFile:  d.UkuranFile,
		})
	}

	return resp
}
