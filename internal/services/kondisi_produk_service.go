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
	GetAllForDropdown(ctx context.Context) ([]map[string]interface{}, error)
}

type kondisiProdukService struct {
	repo           repositories.KondisiProdukRepository
	reorderService *ReorderService
}

func NewKondisiProdukService(repo repositories.KondisiProdukRepository, reorderService *ReorderService) KondisiProdukService {
	return &kondisiProdukService{
		repo:           repo,
		reorderService: reorderService,
	}
}

func (s *kondisiProdukService) Create(ctx context.Context, req *models.CreateKondisiProdukRequest) (*models.KondisiProdukResponse, error) {
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
	} else if req.NamaEN != nil && *req.NamaEN != "" {
		s := utils.GenerateSlug(*req.NamaEN)
		slugEN = &s
	}

	// Auto-increment urutan
	maxUrutan, err := s.repo.GetMaxUrutan(ctx)
	if err != nil {
		return nil, err
	}

	kondisi := &models.KondisiProduk{
		NamaID:    req.NamaID,
		NamaEN:    req.NamaEN,
		Slug:      *slugID,
		SlugID:    slugID,
		SlugEN:    slugEN,
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
		kondisi.NamaID = *req.NamaID
		// Regenerate slug_id from new nama_id (unless manually provided)
		if req.SlugID == nil {
			s := utils.GenerateSlug(*req.NamaID)
			kondisi.SlugID = &s
			kondisi.Slug = s // backward compat
		}
	}
	if req.SlugID != nil && *req.SlugID != "" {
		kondisi.SlugID = req.SlugID
		kondisi.Slug = *req.SlugID // backward compat
	}
	if req.NamaEN != nil {
		kondisi.NamaEN = req.NamaEN
		// Regenerate slug_en from new nama_en (unless manually provided)
		if req.SlugEN == nil && *req.NamaEN != "" {
			s := utils.GenerateSlug(*req.NamaEN)
			kondisi.SlugEN = &s
		}
	}
	if req.SlugEN != nil && *req.SlugEN != "" {
		kondisi.SlugEN = req.SlugEN
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
	kondisi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("kondisi produk tidak ditemukan")
	}

	deletedUrutan := kondisi.Urutan

	// TODO: Check if kondisi has products

	// Soft delete
	if err := s.repo.Delete(ctx, kondisi); err != nil {
		return err
	}

	// Reorder remaining items to fill gap
	return s.reorderService.ReorderAfterDelete(
		ctx,
		"kondisi_produk",
		deletedUrutan,
		"", // No scope
		nil,
	)
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
		ID:        kondisi.ID.String(),
		IsActive:  kondisi.IsActive,
		UpdatedAt: kondisi.UpdatedAt,
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

func (s *kondisiProdukService) GetAllForDropdown(ctx context.Context) ([]map[string]interface{}, error) {
	kondisis, err := s.repo.GetAllForDropdown(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(kondisis))
	for i, k := range kondisis {
		result[i] = map[string]interface{}{
			"id":      k.ID.String(),
			"nama_id": k.NamaID,
			"nama_en": k.NamaEN,
		}
	}
	return result, nil
}

func (s *kondisiProdukService) toResponse(k *models.KondisiProduk) *models.KondisiProdukResponse {
	return &models.KondisiProdukResponse{
		ID:        k.ID.String(),
		Nama:      k.GetNama(),
		SlugID:    k.SlugID,
		SlugEN:    k.SlugEN,
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
