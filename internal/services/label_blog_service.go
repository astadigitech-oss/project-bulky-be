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

type LabelBlogService interface {
	Create(ctx context.Context, req *dto.CreateLabelBlogRequest) (*models.LabelBlog, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateLabelBlogRequest) (*models.LabelBlog, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.LabelBlogDetailResponse, error)
	GetBySlug(ctx context.Context, slug string) (*models.LabelBlog, error)
	GetAll(ctx context.Context, params *dto.LabelBlogFilterRequest) ([]dto.LabelBlogListResponse, models.PaginationMeta, error)
	GetAllActive(ctx context.Context) ([]dto.LabelBlogDropdownResponse, error)
	GetAllPublicWithCount(ctx context.Context) ([]models.LabelBlog, error)
}

type labelBlogService struct {
	repo repositories.LabelBlogRepository
}

func NewLabelBlogService(repo repositories.LabelBlogRepository) LabelBlogService {
	return &labelBlogService{repo: repo}
}

func (s *labelBlogService) Create(ctx context.Context, req *dto.CreateLabelBlogRequest) (*models.LabelBlog, error) {
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
	} else if req.NamaEN != "" {
		s := utils.GenerateSlug(req.NamaEN)
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

	label := &models.LabelBlog{
		NamaID: req.NamaID,
		NamaEN: req.NamaEN,
		Slug:   slugIDVal,
		SlugID: slugIDPtr,
		SlugEN: slugEN,
		Urutan: urutan,
	}

	if err := s.repo.Create(ctx, label); err != nil {
		return nil, err
	}

	return label, nil
}

func (s *labelBlogService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateLabelBlogRequest) (*models.LabelBlog, error) {
	label, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.NamaID != nil {
		label.NamaID = *req.NamaID
	}
	if req.NamaEN != nil {
		label.NamaEN = *req.NamaEN
	}

	// SlugID: explicit > auto dari NamaID > keep existing
	if req.SlugID != nil && *req.SlugID != "" {
		label.SlugID = req.SlugID
		label.Slug = *req.SlugID
	} else if req.NamaID != nil {
		generated := utils.GenerateSlug(*req.NamaID)
		label.SlugID = &generated
		label.Slug = generated
	}

	// SlugEN: explicit > auto dari NamaEN > keep existing
	if req.SlugEN != nil && *req.SlugEN != "" {
		label.SlugEN = req.SlugEN
	} else if req.NamaEN != nil && label.NamaEN != "" {
		generated := utils.GenerateSlug(label.NamaEN)
		label.SlugEN = &generated
	}

	if req.Urutan != nil {
		label.Urutan = *req.Urutan
	}

	if err := s.repo.Update(ctx, label); err != nil {
		return nil, err
	}

	return label, nil
}

func (s *labelBlogService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if label exists
	label, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if label is used in any blog posts
	count, err := s.repo.CountBlogByLabel(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("label tidak dapat dihapus karena masih memiliki artikel blog")
	}

	// Rename slug sebelum soft-delete agar slug lama bisa dipakai ulang
	suffix := "_deleted_" + id.String()[:8]
	var newSlugID *string
	if label.SlugID != nil {
		v := *label.SlugID + suffix
		newSlugID = &v
	}
	var newSlugEN *string
	if label.SlugEN != nil {
		v := *label.SlugEN + suffix
		newSlugEN = &v
	}
	if err := s.repo.UpdateSlugs(ctx, id, label.Slug+suffix, newSlugID, newSlugEN); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Reorder ulang setelah delete
	allLabels, err := s.repo.FindAllOrdered(ctx)
	if err != nil {
		return err
	}
	for i, l := range allLabels {
		if err := s.repo.UpdateUrutan(ctx, l.ID, i+1); err != nil {
			return err
		}
	}

	return nil
}

func (s *labelBlogService) GetByID(ctx context.Context, id uuid.UUID) (*dto.LabelBlogDetailResponse, error) {
	l, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.LabelBlogDetailResponse{
		ID:        l.ID.String(),
		NamaID:    l.NamaID,
		NamaEN:    l.NamaEN,
		SlugID:    l.SlugID,
		SlugEN:    l.SlugEN,
		Urutan:    l.Urutan,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}, nil
}

func (s *labelBlogService) GetBySlug(ctx context.Context, slug string) (*models.LabelBlog, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *labelBlogService) GetAll(ctx context.Context, params *dto.LabelBlogFilterRequest) ([]dto.LabelBlogListResponse, models.PaginationMeta, error) {
	params.SetDefaults()

	labels, meta, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	result := make([]dto.LabelBlogListResponse, len(labels))
	for i, l := range labels {
		result[i] = dto.LabelBlogListResponse{
			ID:        l.ID.String(),
			NamaID:    l.NamaID,
			NamaEN:    l.NamaEN,
			Urutan:    l.Urutan,
			UpdatedAt: l.UpdatedAt,
		}
	}

	return result, meta, nil
}

func (s *labelBlogService) GetAllActive(ctx context.Context) ([]dto.LabelBlogDropdownResponse, error) {
	params := &dto.LabelBlogFilterRequest{}
	params.SetDefaults()
	params.PerPage = 1000
	labelList, _, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []dto.LabelBlogDropdownResponse
	for _, l := range labelList {
		nama := map[string]interface{}{
			"id": l.NamaID,
		}
		if l.NamaEN != "" {
			nama["en"] = l.NamaEN
		} else {
			nama["en"] = l.NamaID
		}

		result = append(result, dto.LabelBlogDropdownResponse{
			ID:     l.ID,
			Nama:   nama,
			SlugID: l.SlugID,
			SlugEN: l.SlugEN,
		})
	}

	return result, nil
}

func (s *labelBlogService) GetAllPublicWithCount(ctx context.Context) ([]models.LabelBlog, error) {
	return s.repo.FindAllPublicWithCount(ctx)
}
