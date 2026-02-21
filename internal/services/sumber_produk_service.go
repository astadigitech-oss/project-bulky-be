package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type SumberProdukService interface {
	Create(ctx context.Context, req *models.CreateSumberProdukRequest) (*models.SumberProdukResponse, error)
	FindByID(ctx context.Context, id string) (*models.SumberProdukResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.SumberProdukResponse, error)
	FindAll(ctx context.Context, params *models.SumberProdukFilterRequest) ([]models.SumberProdukSimpleResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateSumberProdukRequest) (*models.SumberProdukResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
}

type sumberProdukService struct {
	repo repositories.SumberProdukRepository
}

func NewSumberProdukService(repo repositories.SumberProdukRepository) SumberProdukService {
	return &sumberProdukService{repo: repo}
}

func (s *sumberProdukService) Create(ctx context.Context, req *models.CreateSumberProdukRequest) (*models.SumberProdukResponse, error) {
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

	sumber := &models.SumberProduk{
		NamaID:    req.NamaID,
		NamaEN:    req.NamaEN,
		Slug:      *slugID,
		SlugID:    slugID,
		SlugEN:    slugEN,
		Deskripsi: req.Deskripsi,
		IsActive:  true,
	}

	if err := s.repo.Create(ctx, sumber); err != nil {
		return nil, err
	}

	return s.toResponse(sumber), nil
}

func (s *sumberProdukService) FindByID(ctx context.Context, id string) (*models.SumberProdukResponse, error) {
	sumber, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("sumber produk tidak ditemukan")
	}
	return s.toResponse(sumber), nil
}

func (s *sumberProdukService) FindBySlug(ctx context.Context, slug string) (*models.SumberProdukResponse, error) {
	sumber, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("sumber produk tidak ditemukan")
	}
	return s.toResponse(sumber), nil
}

func (s *sumberProdukService) FindAll(ctx context.Context, params *models.SumberProdukFilterRequest) ([]models.SumberProdukSimpleResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	sumbers, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	items := []models.SumberProdukSimpleResponse{}
	for _, sb := range sumbers {
		items = append(items, *s.toSimpleResponse(&sb))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *sumberProdukService) Update(ctx context.Context, id string, req *models.UpdateSumberProdukRequest) (*models.SumberProdukResponse, error) {
	sumber, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("sumber produk tidak ditemukan")
	}

	if req.NamaID != nil {
		sumber.NamaID = *req.NamaID
		// Regenerate slug_id from new nama_id (unless manually provided)
		if req.SlugID == nil || *req.SlugID == "" {
			s := utils.GenerateSlug(*req.NamaID)
			sumber.SlugID = &s
			sumber.Slug = s // backward compat
		}
	}
	if req.SlugID != nil && *req.SlugID != "" {
		sumber.SlugID = req.SlugID
		sumber.Slug = *req.SlugID // backward compat
	}
	if req.NamaEN != nil {
		sumber.NamaEN = req.NamaEN
		// Regenerate slug_en from new nama_en (unless manually provided)
		if (req.SlugEN == nil || *req.SlugEN == "") && *req.NamaEN != "" {
			s := utils.GenerateSlug(*req.NamaEN)
			sumber.SlugEN = &s
		}
	}
	if req.SlugEN != nil && *req.SlugEN != "" {
		sumber.SlugEN = req.SlugEN
	}
	if req.Deskripsi != nil {
		sumber.Deskripsi = req.Deskripsi
	}
	if req.IsActive != nil {
		sumber.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, sumber); err != nil {
		return nil, err
	}

	return s.toResponse(sumber), nil
}

func (s *sumberProdukService) Delete(ctx context.Context, id string) error {
	sumber, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("sumber produk tidak ditemukan")
	}

	// Rename slug dengan suffix _deleted_{8-char-id} agar tidak conflict unique constraint
	suffix := "_deleted_" + id[:8]
	sumber.Slug = sumber.Slug + suffix
	if sumber.SlugID != nil {
		v := *sumber.SlugID + suffix
		sumber.SlugID = &v
	}
	if sumber.SlugEN != nil {
		v := *sumber.SlugEN + suffix
		sumber.SlugEN = &v
	}
	if err := s.repo.Update(ctx, sumber); err != nil {
		return err
	}

	return s.repo.Delete(ctx, sumber)
}

func (s *sumberProdukService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	sumber, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("sumber produk tidak ditemukan")
	}

	sumber.IsActive = !sumber.IsActive
	if err := s.repo.Update(ctx, sumber); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:        sumber.ID.String(),
		IsActive:  sumber.IsActive,
		UpdatedAt: sumber.UpdatedAt,
	}, nil
}

func (s *sumberProdukService) toResponse(sb *models.SumberProduk) *models.SumberProdukResponse {
	return &models.SumberProdukResponse{
		ID:        sb.ID.String(),
		Nama:      sb.GetNama(),
		SlugID:    sb.SlugID,
		SlugEN:    sb.SlugEN,
		Deskripsi: sb.Deskripsi,
		IsActive:  sb.IsActive,
		CreatedAt: sb.CreatedAt,
		UpdatedAt: sb.UpdatedAt,
	}
}

func (s *sumberProdukService) toSimpleResponse(sb *models.SumberProduk) *models.SumberProdukSimpleResponse {
	return &models.SumberProdukSimpleResponse{
		ID:        sb.ID.String(),
		Nama:      sb.GetNama(),
		SlugID:    sb.SlugID,
		SlugEN:    sb.SlugEN,
		IsActive:  sb.IsActive,
		UpdatedAt: sb.UpdatedAt,
	}
}
