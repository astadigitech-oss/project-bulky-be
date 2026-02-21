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
	GetAllForDropdown(ctx context.Context) ([]map[string]interface{}, error)
}

type kondisiPaketService struct {
	repo           repositories.KondisiPaketRepository
	reorderService *ReorderService
}

func NewKondisiPaketService(repo repositories.KondisiPaketRepository, reorderService *ReorderService) KondisiPaketService {
	return &kondisiPaketService{
		repo:           repo,
		reorderService: reorderService,
	}
}

func (s *kondisiPaketService) Create(ctx context.Context, req *models.CreateKondisiPaketRequest) (*models.KondisiPaketResponse, error) {
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

	kondisi := &models.KondisiPaket{
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

func (s *kondisiPaketService) Delete(ctx context.Context, id string) error {
	kondisi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("kondisi paket tidak ditemukan")
	}

	// Check if kondisi paket has products
	hasProducts, err := s.repo.HasProducts(ctx, id)
	if err != nil {
		return errors.New("gagal memeriksa relasi kondisi paket")
	}
	if hasProducts {
		return errors.New("kondisi paket tidak dapat dihapus karena masih digunakan oleh produk")
	}

	deletedUrutan := kondisi.Urutan

	// Soft delete
	if err := s.repo.Delete(ctx, kondisi); err != nil {
		return err
	}

	// Reorder remaining items to fill gap
	return s.reorderService.ReorderAfterDelete(
		ctx,
		"kondisi_paket",
		deletedUrutan,
		"", // No scope
		nil,
	)
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
		ID:        kondisi.ID.String(),
		IsActive:  kondisi.IsActive,
		UpdatedAt: kondisi.UpdatedAt,
	}, nil
}

func (s *kondisiPaketService) Reorder(ctx context.Context, req *models.ReorderRequest) error {
	return s.repo.UpdateOrder(ctx, req.Items)
}

func (s *kondisiPaketService) GetAllForDropdown(ctx context.Context) ([]map[string]interface{}, error) {
	pakets, err := s.repo.GetAllForDropdown(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(pakets))
	for i, p := range pakets {
		result[i] = map[string]interface{}{
			"id":      p.ID.String(),
			"nama_id": p.NamaID,
			"nama_en": p.NamaEN,
		}
	}
	return result, nil
}

func (s *kondisiPaketService) toResponse(k *models.KondisiPaket) *models.KondisiPaketResponse {
	return &models.KondisiPaketResponse{
		ID:        k.ID.String(),
		Nama:      k.GetNama(),
		SlugID:    k.SlugID,
		SlugEN:    k.SlugEN,
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
