package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type KondisiPaketService interface {
	Create(ctx context.Context, req *models.CreateKondisiPaketRequest) (*models.KondisiPaketResponse, error)
	FindByID(ctx context.Context, id string) (*models.KondisiPaketResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.KondisiPaketResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.KondisiPaketSimpleResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateKondisiPaketRequest) (*models.KondisiPaketResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	Reorder(ctx context.Context, req *models.ReorderRequest) error
}

type kondisiPaketService struct {
	repo repositories.KondisiPaketRepository
}

func NewKondisiPaketService(repo repositories.KondisiPaketRepository) KondisiPaketService {
	return &kondisiPaketService{repo: repo}
}

func (s *kondisiPaketService) Create(ctx context.Context, req *models.CreateKondisiPaketRequest) (*models.KondisiPaketResponse, error) {
	slug := utils.GenerateSlug(req.NamaID)

	exists, _ := s.repo.ExistsBySlug(ctx, slug, nil)
	if exists {
		return nil, errors.New("kondisi paket dengan nama tersebut sudah ada")
	}

	kondisi := &models.KondisiPaket{
		NamaID:    req.NamaID,
		NamaEN:    req.NamaEN,
		Slug:      slug,
		Deskripsi: req.Deskripsi,
		IsActive:  true,
	}

	if req.Urutan != nil {
		kondisi.Urutan = *req.Urutan
	}

	if err := s.repo.Create(ctx, kondisi); err != nil {
		return nil, err
	}

	return s.toResponse(kondisi), nil
}

func (s *kondisiPaketService) FindByID(ctx context.Context, id string) (*models.KondisiPaketResponse, error) {
	kondisi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kondisi paket tidak ditemukan")
	}
	return s.toResponse(kondisi), nil
}

func (s *kondisiPaketService) FindBySlug(ctx context.Context, slug string) (*models.KondisiPaketResponse, error) {
	kondisi, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("kondisi paket tidak ditemukan")
	}
	return s.toResponse(kondisi), nil
}

func (s *kondisiPaketService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.KondisiPaketSimpleResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	kondisis, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	// Ensure empty array instead of null
	if kondisis == nil {
		kondisis = []models.KondisiPaketSimpleResponse{}
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return kondisis, &meta, nil
}

func (s *kondisiPaketService) Update(ctx context.Context, id string, req *models.UpdateKondisiPaketRequest) (*models.KondisiPaketResponse, error) {
	kondisi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kondisi paket tidak ditemukan")
	}

	if req.NamaID != nil {
		newSlug := utils.GenerateSlug(*req.NamaID)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("kondisi paket dengan nama tersebut sudah ada")
		}
		kondisi.NamaID = *req.NamaID
		kondisi.Slug = newSlug
	}
	if req.NamaEN != nil {
		kondisi.NamaEN = req.NamaEN
	}
	if req.Deskripsi != nil {
		kondisi.Deskripsi = req.Deskripsi
	}
	if req.Urutan != nil {
		kondisi.Urutan = *req.Urutan
	}
	if req.IsActive != nil {
		kondisi.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, kondisi); err != nil {
		return nil, err
	}

	return s.toResponse(kondisi), nil
}

func (s *kondisiPaketService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("kondisi paket tidak ditemukan")
	}

	// TODO: Check if kondisi paket has products

	return s.repo.Delete(ctx, id)
}

func (s *kondisiPaketService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	kondisi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kondisi paket tidak ditemukan")
	}

	kondisi.IsActive = !kondisi.IsActive
	if err := s.repo.Update(ctx, kondisi); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       kondisi.ID.String(),
		IsActive: kondisi.IsActive,
	}, nil
}

func (s *kondisiPaketService) Reorder(ctx context.Context, req *models.ReorderRequest) error {
	return s.repo.UpdateOrder(ctx, req.Items)
}

func (s *kondisiPaketService) toResponse(k *models.KondisiPaket) *models.KondisiPaketResponse {
	return &models.KondisiPaketResponse{
		ID:        k.ID.String(),
		Nama:      k.GetNama(),
		Slug:      k.Slug,
		Deskripsi: k.Deskripsi,
		Urutan:    k.Urutan,
		IsActive:  k.IsActive,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func (s *kondisiPaketService) toSimpleResponse(k *models.KondisiPaket) *models.KondisiPaketSimpleResponse {
	return &models.KondisiPaketSimpleResponse{
		ID:   k.ID.String(),
		Nama: k.GetNama(),
		// Slug:         k.Slug,
		// Deskripsi:    k.Deskripsi,
		Urutan:   k.Urutan,
		IsActive: k.IsActive,
		// CreatedAt:    k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}
