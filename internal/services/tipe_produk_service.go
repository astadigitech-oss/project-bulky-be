package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type TipeProdukService interface {
	Create(ctx context.Context, req *models.CreateTipeProdukRequest) (*models.TipeProdukResponse, error)
	FindByID(ctx context.Context, id string) (*models.TipeProdukResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.TipeProdukResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.TipeProdukResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateTipeProdukRequest) (*models.TipeProdukResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	Reorder(ctx context.Context, req *models.ReorderRequest) error
}

type tipeProdukService struct {
	repo repositories.TipeProdukRepository
}

func NewTipeProdukService(repo repositories.TipeProdukRepository) TipeProdukService {
	return &tipeProdukService{repo: repo}
}

func (s *tipeProdukService) Create(ctx context.Context, req *models.CreateTipeProdukRequest) (*models.TipeProdukResponse, error) {
	slug := utils.GenerateSlug(req.Nama)

	exists, _ := s.repo.ExistsBySlug(ctx, slug, nil)
	if exists {
		return nil, errors.New("tipe produk dengan nama tersebut sudah ada")
	}

	tipe := &models.TipeProduk{
		Nama:      req.Nama,
		Slug:      slug,
		Deskripsi: req.Deskripsi,
		IsActive:  true,
	}

	if req.Urutan != nil {
		tipe.Urutan = *req.Urutan
	}

	if err := s.repo.Create(ctx, tipe); err != nil {
		return nil, err
	}

	return s.toResponse(tipe), nil
}

func (s *tipeProdukService) FindByID(ctx context.Context, id string) (*models.TipeProdukResponse, error) {
	tipe, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("tipe produk tidak ditemukan")
	}
	return s.toResponse(tipe), nil
}

func (s *tipeProdukService) FindBySlug(ctx context.Context, slug string) (*models.TipeProdukResponse, error) {
	tipe, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("tipe produk tidak ditemukan")
	}
	return s.toResponse(tipe), nil
}

func (s *tipeProdukService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.TipeProdukResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	tipes, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.TipeProdukResponse
	for _, t := range tipes {
		items = append(items, *s.toResponse(&t))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *tipeProdukService) Update(ctx context.Context, id string, req *models.UpdateTipeProdukRequest) (*models.TipeProdukResponse, error) {
	tipe, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("tipe produk tidak ditemukan")
	}

	if req.Nama != nil {
		newSlug := utils.GenerateSlug(*req.Nama)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("tipe produk dengan nama tersebut sudah ada")
		}
		tipe.Nama = *req.Nama
		tipe.Slug = newSlug
	}
	if req.Deskripsi != nil {
		tipe.Deskripsi = req.Deskripsi
	}
	if req.Urutan != nil {
		tipe.Urutan = *req.Urutan
	}
	if req.IsActive != nil {
		tipe.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, tipe); err != nil {
		return nil, err
	}

	return s.toResponse(tipe), nil
}

func (s *tipeProdukService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("tipe produk tidak ditemukan")
	}
	// TODO: Check if tipe has products
	return s.repo.Delete(ctx, id)
}

func (s *tipeProdukService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	tipe, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("tipe produk tidak ditemukan")
	}

	tipe.IsActive = !tipe.IsActive
	if err := s.repo.Update(ctx, tipe); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       tipe.ID.String(),
		IsActive: tipe.IsActive,
	}, nil
}

func (s *tipeProdukService) Reorder(ctx context.Context, req *models.ReorderRequest) error {
	return s.repo.UpdateOrder(ctx, req.Items)
}

func (s *tipeProdukService) toResponse(t *models.TipeProduk) *models.TipeProdukResponse {
	return &models.TipeProdukResponse{
		ID:           t.ID.String(),
		Nama:         t.Nama,
		Slug:         t.Slug,
		Deskripsi:    t.Deskripsi,
		Urutan:       t.Urutan,
		IsActive:     t.IsActive,
		JumlahProduk: 0, // TODO: Count from produk table
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}
