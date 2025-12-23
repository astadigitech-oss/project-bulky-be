package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type WarehouseService interface {
	Create(ctx context.Context, req *models.CreateWarehouseRequest) (*models.WarehouseResponse, error)
	FindByID(ctx context.Context, id string) (*models.WarehouseResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest, kota string) ([]models.WarehouseResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateWarehouseRequest) (*models.WarehouseResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
}

type warehouseService struct {
	repo repositories.WarehouseRepository
}

func NewWarehouseService(repo repositories.WarehouseRepository) WarehouseService {
	return &warehouseService{repo: repo}
}

func (s *warehouseService) Create(ctx context.Context, req *models.CreateWarehouseRequest) (*models.WarehouseResponse, error) {
	slug := utils.GenerateSlug(req.Nama)

	exists, _ := s.repo.ExistsBySlug(ctx, slug, nil)
	if exists {
		return nil, errors.New("warehouse dengan nama tersebut sudah ada")
	}

	warehouse := &models.Warehouse{
		Nama:     req.Nama,
		Slug:     slug,
		Alamat:   req.Alamat,
		Kota:     req.Kota,
		KodePos:  req.KodePos,
		Telepon:  req.Telepon,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, warehouse); err != nil {
		return nil, err
	}

	return s.toResponse(warehouse), nil
}


func (s *warehouseService) FindByID(ctx context.Context, id string) (*models.WarehouseResponse, error) {
	warehouse, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}
	return s.toResponse(warehouse), nil
}

func (s *warehouseService) FindAll(ctx context.Context, params *models.PaginationRequest, kota string) ([]models.WarehouseResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	warehouses, total, err := s.repo.FindAll(ctx, params, kota)
	if err != nil {
		return nil, nil, err
	}

	var items []models.WarehouseResponse
	for _, w := range warehouses {
		items = append(items, *s.toResponse(&w))
	}

	totalHalaman := (total + int64(params.PerHalaman) - 1) / int64(params.PerHalaman)

	meta := &models.PaginationMeta{
		Halaman:      params.Halaman,
		PerHalaman:   params.PerHalaman,
		TotalData:    total,
		TotalHalaman: totalHalaman,
	}

	return items, meta, nil
}

func (s *warehouseService) Update(ctx context.Context, id string, req *models.UpdateWarehouseRequest) (*models.WarehouseResponse, error) {
	warehouse, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}

	if req.Nama != nil {
		newSlug := utils.GenerateSlug(*req.Nama)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("warehouse dengan nama tersebut sudah ada")
		}
		warehouse.Nama = *req.Nama
		warehouse.Slug = newSlug
	}
	if req.Alamat != nil {
		warehouse.Alamat = req.Alamat
	}
	if req.Kota != nil {
		warehouse.Kota = req.Kota
	}
	if req.KodePos != nil {
		warehouse.KodePos = req.KodePos
	}
	if req.Telepon != nil {
		warehouse.Telepon = req.Telepon
	}
	if req.IsActive != nil {
		warehouse.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, warehouse); err != nil {
		return nil, err
	}

	return s.toResponse(warehouse), nil
}

func (s *warehouseService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("warehouse tidak ditemukan")
	}
	// TODO: Check if warehouse has products
	return s.repo.Delete(ctx, id)
}

func (s *warehouseService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	warehouse, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}

	warehouse.IsActive = !warehouse.IsActive
	if err := s.repo.Update(ctx, warehouse); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       warehouse.ID.String(),
		IsActive: warehouse.IsActive,
	}, nil
}

func (s *warehouseService) toResponse(w *models.Warehouse) *models.WarehouseResponse {
	return &models.WarehouseResponse{
		ID:           w.ID.String(),
		Nama:         w.Nama,
		Slug:         w.Slug,
		Alamat:       w.Alamat,
		Kota:         w.Kota,
		KodePos:      w.KodePos,
		Telepon:      w.Telepon,
		IsActive:     w.IsActive,
		JumlahProduk: 0, // TODO: Count from produk table
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
	}
}
