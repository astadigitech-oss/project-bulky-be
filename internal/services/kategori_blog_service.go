package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
)

type KategoriBlogService interface {
	Create(ctx context.Context, req *dto.CreateKategoriBlogRequest) (*models.KategoriBlog, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateKategoriBlogRequest) (*models.KategoriBlog, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.KategoriBlogDetailResponse, error)
	GetBySlug(ctx context.Context, slug string) (*models.KategoriBlog, error)
	GetAll(ctx context.Context, isActive *bool) ([]models.KategoriBlog, error)
	GetAllPaginated(ctx context.Context, params *dto.KategoriBlogFilterRequest) ([]dto.KategoriBlogListResponse, models.PaginationMeta, error)
	GetAllActive(ctx context.Context) ([]dto.KategoriBlogDropdownResponse, error)
	ToggleStatus(ctx context.Context, id uuid.UUID) error
	GetAllPublicWithCount(ctx context.Context) ([]models.KategoriBlog, error)
}

type kategoriBlogService struct {
	repo repositories.KategoriBlogRepository
}

func NewKategoriBlogService(repo repositories.KategoriBlogRepository) KategoriBlogService {
	return &kategoriBlogService{repo: repo}
}

func (s *kategoriBlogService) Create(ctx context.Context, req *dto.CreateKategoriBlogRequest) (*models.KategoriBlog, error) {
	// Generate/use slug_id
	var slugIDVal string
	if req.SlugID != nil && *req.SlugID != "" {
		slugIDVal = *req.SlugID
	} else {
		slugIDVal = utils.GenerateSlug(req.NamaID)
	}
	slugIDPtr := &slugIDVal

	// Generate/use slug_en
	var slugEN *string
	if req.SlugEN != nil && *req.SlugEN != "" {
		slugEN = req.SlugEN
	} else if req.NamaEN != nil && *req.NamaEN != "" {
		s := utils.GenerateSlug(*req.NamaEN)
		slugEN = &s
	}

	// Auto-assign urutan jika tidak disertakan
	urutan := req.Urutan
	if urutan == 0 {
		maxUrutan, err := s.repo.GetMaxUrutan(ctx)
		if err != nil {
			return nil, err
		}
		urutan = maxUrutan + 1
	}

	kategori := &models.KategoriBlog{
		NamaID:   req.NamaID,
		NamaEN:   req.NamaEN,
		Slug:     slugIDVal,
		SlugID:   slugIDPtr,
		SlugEN:   slugEN,
		IsActive: req.IsActive,
		Urutan:   urutan,
	}

	if err := s.repo.Create(ctx, kategori); err != nil {
		return nil, err
	}

	return kategori, nil
}

func (s *kategoriBlogService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateKategoriBlogRequest) (*models.KategoriBlog, error) {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.NamaID != nil {
		kategori.NamaID = *req.NamaID
	}
	if req.NamaEN != nil {
		kategori.NamaEN = req.NamaEN
	}

	// SlugID: explicit > auto dari NamaID > keep existing
	if req.SlugID != nil && *req.SlugID != "" {
		kategori.SlugID = req.SlugID
		kategori.Slug = *req.SlugID
	} else if req.NamaID != nil {
		generated := utils.GenerateSlug(*req.NamaID)
		kategori.SlugID = &generated
		kategori.Slug = generated
	}

	// SlugEN: explicit > auto dari NamaEN > keep existing
	if req.SlugEN != nil && *req.SlugEN != "" {
		kategori.SlugEN = req.SlugEN
	} else if req.NamaEN != nil && *req.NamaEN != "" {
		generated := utils.GenerateSlug(*req.NamaEN)
		kategori.SlugEN = &generated
	}

	if req.IsActive != nil {
		kategori.IsActive = *req.IsActive
	}
	if req.Urutan != nil {
		kategori.Urutan = *req.Urutan
	}

	if err := s.repo.Update(ctx, kategori); err != nil {
		return nil, err
	}

	return kategori, nil
}

func (s *kategoriBlogService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if kategori exists
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if kategori has blogs
	count, err := s.repo.CountBlogByKategori(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("kategori tidak dapat dihapus karena masih memiliki artikel blog")
	}

	// Rename slug sebelum soft-delete agar slug lama bisa dipakai ulang
	suffix := "_deleted_" + id.String()[:8]
	var newSlugID *string
	if kategori.SlugID != nil {
		v := *kategori.SlugID + suffix
		newSlugID = &v
	}
	var newSlugEN *string
	if kategori.SlugEN != nil {
		v := *kategori.SlugEN + suffix
		newSlugEN = &v
	}
	if err := s.repo.UpdateSlugs(ctx, id, kategori.Slug+suffix, newSlugID, newSlugEN); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Reorder ulang setelah delete
	allKategoris, err := s.repo.FindAllOrdered(ctx)
	if err != nil {
		return err
	}
	for i, k := range allKategoris {
		if err := s.repo.UpdateUrutan(ctx, k.ID, i+1); err != nil {
			return err
		}
	}

	return nil
}

func (s *kategoriBlogService) GetByID(ctx context.Context, id uuid.UUID) (*dto.KategoriBlogDetailResponse, error) {
	k, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.KategoriBlogDetailResponse{
		ID:        k.ID.String(),
		NamaID:    k.NamaID,
		NamaEN:    k.NamaEN,
		SlugID:    k.SlugID,
		SlugEN:    k.SlugEN,
		Urutan:    k.Urutan,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}, nil
}

func (s *kategoriBlogService) GetBySlug(ctx context.Context, slug string) (*models.KategoriBlog, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *kategoriBlogService) GetAll(ctx context.Context, isActive *bool) ([]models.KategoriBlog, error) {
	return s.repo.FindAll(ctx, isActive)
}

func (s *kategoriBlogService) GetAllPaginated(ctx context.Context, params *dto.KategoriBlogFilterRequest) ([]dto.KategoriBlogListResponse, models.PaginationMeta, error) {
	params.SetDefaults()

	kategoris, total, err := s.repo.FindAllPaginated(ctx, params)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	result := make([]dto.KategoriBlogListResponse, len(kategoris))
	for i, k := range kategoris {
		result[i] = dto.KategoriBlogListResponse{
			ID:        k.ID.String(),
			NamaID:    k.NamaID,
			NamaEN:    k.NamaEN,
			Urutan:    k.Urutan,
			UpdatedAt: k.UpdatedAt,
		}
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return result, meta, nil
}

func (s *kategoriBlogService) ToggleStatus(ctx context.Context, id uuid.UUID) error {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	kategori.IsActive = !kategori.IsActive
	return s.repo.Update(ctx, kategori)
}

func (s *kategoriBlogService) GetAllActive(ctx context.Context) ([]dto.KategoriBlogDropdownResponse, error) {
	isActive := true
	kategoriList, err := s.repo.FindAll(ctx, &isActive)
	if err != nil {
		return nil, err
	}

	var result []dto.KategoriBlogDropdownResponse
	for _, k := range kategoriList {
		nama := map[string]interface{}{
			"id": k.NamaID,
		}
		if k.NamaEN != nil {
			nama["en"] = *k.NamaEN
		} else {
			nama["en"] = k.NamaID
		}

		result = append(result, dto.KategoriBlogDropdownResponse{
			ID:     k.ID,
			Nama:   nama,
			SlugID: k.SlugID,
			SlugEN: k.SlugEN,
		})
	}

	return result, nil
}

func (s *kategoriBlogService) GetAllPublicWithCount(ctx context.Context) ([]models.KategoriBlog, error) {
	return s.repo.FindAllPublicWithCount(ctx)
}
