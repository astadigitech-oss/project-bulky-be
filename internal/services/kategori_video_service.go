package services

import (
	"context"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
)

type KategoriVideoService interface {
	Create(ctx context.Context, req *dto.CreateKategoriVideoRequest) (*models.KategoriVideo, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateKategoriVideoRequest) (*models.KategoriVideo, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.KategoriVideoDetailResponse, error)
	GetBySlug(ctx context.Context, slug string) (*models.KategoriVideo, error)
	GetAll(ctx context.Context, isActive *bool) ([]models.KategoriVideo, error)
	GetAllPaginated(ctx context.Context, params *dto.KategoriVideoFilterRequest) ([]dto.KategoriVideoListResponse, models.PaginationMeta, error)
	GetAllActive(ctx context.Context) ([]dto.KategoriVideoDropdownResponse, error)
	ToggleStatus(ctx context.Context, id uuid.UUID) error
	GetAllPublicWithCount(ctx context.Context) ([]models.KategoriVideo, error)
}

type kategoriVideoService struct {
	repo repositories.KategoriVideoRepository
}

func NewKategoriVideoService(repo repositories.KategoriVideoRepository) KategoriVideoService {
	return &kategoriVideoService{repo: repo}
}

func (s *kategoriVideoService) Create(ctx context.Context, req *dto.CreateKategoriVideoRequest) (*models.KategoriVideo, error) {
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

	// Auto-increment urutan: ambil nilai max dari DB
	maxUrutan, err := s.repo.GetMaxUrutan(ctx)
	if err != nil {
		return nil, err
	}

	kategori := &models.KategoriVideo{
		NamaID:   req.NamaID,
		NamaEN:   req.NamaEN,
		Slug:     slugIDVal,
		SlugID:   slugIDPtr,
		SlugEN:   slugEN,
		IsActive: req.IsActive,
		Urutan:   maxUrutan + 1,
	}

	if err := s.repo.Create(ctx, kategori); err != nil {
		return nil, err
	}

	return kategori, nil
}

func (s *kategoriVideoService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateKategoriVideoRequest) (*models.KategoriVideo, error) {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.NamaID != nil {
		kategori.NamaID = *req.NamaID
		if req.SlugID == nil || *req.SlugID == "" {
			s := utils.GenerateSlug(*req.NamaID)
			kategori.SlugID = &s
			kategori.Slug = s
		}
	}
	if req.NamaEN != nil {
		kategori.NamaEN = req.NamaEN
		// Auto-generate slug_en dari nama_en jika slug_en tidak disertakan
		if (req.SlugEN == nil || *req.SlugEN == "") && *req.NamaEN != "" {
			s := utils.GenerateSlug(*req.NamaEN)
			kategori.SlugEN = &s
		}
	}
	if req.SlugID != nil && *req.SlugID != "" {
		kategori.SlugID = req.SlugID
		kategori.Slug = *req.SlugID
	}
	if req.SlugEN != nil && *req.SlugEN != "" {
		kategori.SlugEN = req.SlugEN
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

func (s *kategoriVideoService) Delete(ctx context.Context, id uuid.UUID) error {
	// Ambil data yang akan dihapus
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Rename slug dengan suffix _deleted_{8-char-id} agar tidak conflict unique constraint
	suffix := "_deleted_" + id.String()[:8]
	kategori.Slug = kategori.Slug + suffix
	if kategori.SlugID != nil {
		s := *kategori.SlugID + suffix
		kategori.SlugID = &s
	}
	if kategori.SlugEN != nil {
		s := *kategori.SlugEN + suffix
		kategori.SlugEN = &s
	}
	if err := s.repo.Update(ctx, kategori); err != nil {
		return err
	}

	// Soft delete
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Reorder sisa data: compact urutan menjadi 1, 2, 3, ...
	remaining, err := s.repo.FindAll(ctx, nil)
	if err != nil {
		return err
	}
	for i, item := range remaining {
		if err := s.repo.UpdateUrutan(ctx, item.ID, i+1); err != nil {
			return err
		}
	}

	return nil
}

func (s *kategoriVideoService) GetByID(ctx context.Context, id uuid.UUID) (*dto.KategoriVideoDetailResponse, error) {
	k, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.KategoriVideoDetailResponse{
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

func (s *kategoriVideoService) GetBySlug(ctx context.Context, slug string) (*models.KategoriVideo, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *kategoriVideoService) GetAll(ctx context.Context, isActive *bool) ([]models.KategoriVideo, error) {
	return s.repo.FindAll(ctx, isActive)
}

func (s *kategoriVideoService) GetAllPaginated(ctx context.Context, params *dto.KategoriVideoFilterRequest) ([]dto.KategoriVideoListResponse, models.PaginationMeta, error) {
	params.SetDefaults()

	kategoris, total, err := s.repo.FindAllPaginated(ctx, params)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	result := make([]dto.KategoriVideoListResponse, len(kategoris))
	for i, k := range kategoris {
		result[i] = dto.KategoriVideoListResponse{
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

func (s *kategoriVideoService) ToggleStatus(ctx context.Context, id uuid.UUID) error {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	kategori.IsActive = !kategori.IsActive
	return s.repo.Update(ctx, kategori)
}

func (s *kategoriVideoService) GetAllActive(ctx context.Context) ([]dto.KategoriVideoDropdownResponse, error) {
	isActive := true
	kategoriList, err := s.repo.FindAll(ctx, &isActive)
	if err != nil {
		return nil, err
	}

	var result []dto.KategoriVideoDropdownResponse
	for _, k := range kategoriList {
		nama := map[string]interface{}{
			"id": k.NamaID,
		}
		if k.NamaEN != nil {
			nama["en"] = *k.NamaEN
		} else {
			nama["en"] = k.NamaID
		}

		result = append(result, dto.KategoriVideoDropdownResponse{
			ID:     k.ID,
			Nama:   nama,
			SlugID: k.SlugID,
			SlugEN: k.SlugEN,
		})
	}

	return result, nil
}

func (s *kategoriVideoService) GetAllPublicWithCount(ctx context.Context) ([]models.KategoriVideo, error) {
	return s.repo.FindAllPublicWithCount(ctx)
}
