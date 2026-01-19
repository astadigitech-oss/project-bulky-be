package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type KondisiProdukService interface {
	Create(ctx context.Context, req *models.CreateKondisiProdukRequest) (*models.KondisiProdukResponse, error)
	FindByID(ctx context.Context, id string) (*models.KondisiProdukResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.KondisiProdukResponse, error)
	FindAll(ctx context.Context, params *models.KondisiProdukFilterRequest) ([]models.KondisiProdukSimpleResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateKondisiProdukRequest) (*models.KondisiProdukResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	Reorder(ctx context.Context, req *models.ReorderRequest) error
}

type kondisiProdukService struct {
	repo repositories.KondisiProdukRepository
}

func NewKondisiProdukService(repo repositories.KondisiProdukRepository) KondisiProdukService {
	return &kondisiProdukService{repo: repo}
}

func (s *kondisiProdukService) Create(ctx context.Context, req *models.CreateKondisiProdukRequest) (*models.KondisiProdukResponse, error) {
	slug := utils.GenerateSlug(req.NamaID)

	exists, _ := s.repo.ExistsBySlug(ctx, slug, nil)
	if exists {
		return nil, errors.New("kondisi produk dengan nama tersebut sudah ada")
	}

	// Auto-increment urutan
	maxUrutan, err := s.repo.GetMaxUrutan(ctx)
	if err != nil {
		return nil, err
	}

	kondisi := &models.KondisiProduk{
		NamaID:    req.NamaID,
		NamaEN:    req.NamaEN,
		Slug:      slug,
		Deskripsi: req.Deskripsi,
		Urutan:    maxUrutan + 1,
		IsActive:  true,
	}

	if err := s.repo.Create(ctx, kondisi); err != nil {
		return nil, err
	}

	return s.toResponse(kondisi), nil
}

func (s *kondisiProdukService) FindByID(ctx context.Context, id string) (*models.KondisiProdukResponse, error) {
	kondisi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kondisi produk tidak ditemukan")
	}
	return s.toResponse(kondisi), nil
}

func (s *kondisiProdukService) FindBySlug(ctx context.Context, slug string) (*models.KondisiProdukResponse, error) {
	kondisi, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("kondisi produk tidak ditemukan")
	}
	return s.toResponse(kondisi), nil
}

func (s *kondisiProdukService) FindAll(ctx context.Context, params *models.KondisiProdukFilterRequest) ([]models.KondisiProdukSimpleResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	kondisis, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	items := []models.KondisiProdukSimpleResponse{}
	for _, k := range kondisis {
		items = append(items, *s.toSimpleResponse(&k))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *kondisiProdukService) Update(ctx context.Context, id string, req *models.UpdateKondisiProdukRequest) (*models.KondisiProdukResponse, error) {
	kondisi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kondisi produk tidak ditemukan")
	}

	if req.NamaID != nil {
		newSlug := utils.GenerateSlug(*req.NamaID)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("kondisi produk dengan nama tersebut sudah ada")
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
	if req.IsActive != nil {
		kondisi.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, kondisi); err != nil {
		return nil, err
	}

	return s.toResponse(kondisi), nil
}

func (s *kondisiProdukService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("kondisi produk tidak ditemukan")
	}

	// TODO: Check if kondisi has products

	return s.repo.Delete(ctx, id)
}

func (s *kondisiProdukService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	kondisi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kondisi produk tidak ditemukan")
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

func (s *kondisiProdukService) Reorder(ctx context.Context, req *models.ReorderRequest) error {
	// Validasi: pastikan semua ID yang diberikan valid
	for _, item := range req.Items {
		_, err := s.repo.FindByID(ctx, item.ID)
		if err != nil {
			return errors.New("Data kondisi produk tidak ditemukan")
		}
	}

	return s.repo.UpdateOrder(ctx, req.Items)
}

func (s *kondisiProdukService) toResponse(k *models.KondisiProduk) *models.KondisiProdukResponse {
	return &models.KondisiProdukResponse{
		ID:        k.ID.String(),
		Nama:      k.GetNama(),
		Slug:      k.Slug,
		Deskripsi: k.Deskripsi,
		Urutan:    k.Urutan,
		IsActive:  k.IsActive,
		// JumlahProduk: 0, // TODO: Count from produk table
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func (s *kondisiProdukService) toSimpleResponse(k *models.KondisiProduk) *models.KondisiProdukSimpleResponse {
	return &models.KondisiProdukSimpleResponse{
		ID:   k.ID.String(),
		Nama: k.GetNama(),
		// Slug:         k.Slug,
		// Deskripsi:    k.Deskripsi,
		Urutan:   k.Urutan,
		IsActive: k.IsActive,
		// JumlahProduk: 0, // TODO: Count from produk table
		// CreatedAt:    k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}
