package services

import (
	"context"
	"errors"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetodePembayaranGroupService interface {
	GetAll(ctx context.Context, params *models.PaginationRequest) ([]models.MetodePembayaranGroupListResponse, *models.PaginationMeta, error)
	GetByID(ctx context.Context, id string) (*models.MetodePembayaranGroupDetailResponse, error)
	Create(ctx context.Context, req *models.CreateMetodePembayaranGroupRequest) (*models.MetodePembayaranGroupResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateMetodePembayaranGroupRequest) (*models.MetodePembayaranGroupResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.MetodePembayaranGroupResponse, error)
}

type metodePembayaranGroupService struct {
	repo repositories.MetodePembayaranGroupRepository
}

func NewMetodePembayaranGroupService(repo repositories.MetodePembayaranGroupRepository) MetodePembayaranGroupService {
	return &metodePembayaranGroupService{repo: repo}
}

func (s *metodePembayaranGroupService) GetAll(ctx context.Context, params *models.PaginationRequest) ([]models.MetodePembayaranGroupListResponse, *models.PaginationMeta, error) {
	groups, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]models.MetodePembayaranGroupListResponse, len(groups))
	for i, group := range groups {
		// Count metode for each group
		jumlahMetode, _ := s.repo.CountActiveMetode(ctx, group.ID)

		responses[i] = models.MetodePembayaranGroupListResponse{
			ID:           group.ID.String(),
			Nama:         group.Nama,
			Urutan:       group.Urutan,
			IsActive:     group.IsActive,
			JumlahMetode: int(jumlahMetode),
			UpdatedAt:    group.UpdatedAt,
		}
	}

	meta := &models.PaginationMeta{
		FirstPage:   1,
		LastPage:    int((total + int64(params.PerPage) - 1) / int64(params.PerPage)),
		CurrentPage: params.Page,
		From:        (params.Page-1)*params.PerPage + 1,
		Last:        len(groups),
		Total:       total,
		PerPage:     params.PerPage,
	}

	return responses, meta, nil
}

func (s *metodePembayaranGroupService) GetByID(ctx context.Context, id string) (*models.MetodePembayaranGroupDetailResponse, error) {
	groupID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID group tidak valid")
	}

	group, err := s.repo.FindByIDWithMetode(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Group metode pembayaran tidak ditemukan")
		}
		return nil, err
	}

	// Count metode
	jumlahMetode, _ := s.repo.CountActiveMetode(ctx, group.ID)

	// Map metode pembayaran
	metodes := make([]models.MetodePembayaranResponse, len(group.MetodePembayaran))
	for i, metode := range group.MetodePembayaran {
		metodes[i] = models.MetodePembayaranResponse{
			ID:       metode.ID.String(),
			Nama:     metode.Nama,
			Kode:     metode.Kode,
			Logo:     metode.Logo,
			Urutan:   metode.Urutan,
			IsActive: metode.IsActive,
		}
	}

	return &models.MetodePembayaranGroupDetailResponse{
		ID:               group.ID.String(),
		Nama:             group.Nama,
		Urutan:           group.Urutan,
		IsActive:         group.IsActive,
		JumlahMetode:     int(jumlahMetode),
		MetodePembayaran: metodes,
		CreatedAt:        group.CreatedAt,
		UpdatedAt:        group.UpdatedAt,
	}, nil
}

func (s *metodePembayaranGroupService) Create(ctx context.Context, req *models.CreateMetodePembayaranGroupRequest) (*models.MetodePembayaranGroupResponse, error) {
	// Check if nama already exists
	exists, err := s.repo.CheckByName(ctx, req.Nama, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("Nama group sudah digunakan")
	}

	group := &models.MetodePembayaranGroup{
		Nama:     req.Nama,
		Urutan:   req.Urutan,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, group); err != nil {
		return nil, err
	}

	return &models.MetodePembayaranGroupResponse{
		ID:        group.ID.String(),
		Nama:      group.Nama,
		Urutan:    group.Urutan,
		IsActive:  group.IsActive,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}, nil
}

func (s *metodePembayaranGroupService) Update(ctx context.Context, id string, req *models.UpdateMetodePembayaranGroupRequest) (*models.MetodePembayaranGroupResponse, error) {
	groupID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID group tidak valid")
	}

	group, err := s.repo.FindByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Group metode pembayaran tidak ditemukan")
		}
		return nil, err
	}

	// Check if nama already exists (exclude current group)
	if req.Nama != nil {
		exists, err := s.repo.CheckByName(ctx, *req.Nama, &groupID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("Nama group sudah digunakan")
		}
		group.Nama = *req.Nama
	}

	if req.Urutan != nil {
		group.Urutan = *req.Urutan
	}

	if req.IsActive != nil {
		group.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, group); err != nil {
		return nil, err
	}

	return &models.MetodePembayaranGroupResponse{
		ID:        group.ID.String(),
		Nama:      group.Nama,
		Urutan:    group.Urutan,
		IsActive:  group.IsActive,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}, nil
}

func (s *metodePembayaranGroupService) Delete(ctx context.Context, id string) error {
	groupID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("ID group tidak valid")
	}

	group, err := s.repo.FindByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Group metode pembayaran tidak ditemukan")
		}
		return err
	}

	// Check if group has active metode pembayaran
	count, err := s.repo.CountActiveMetode(ctx, group.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Tidak dapat menghapus group yang masih memiliki metode pembayaran aktif")
	}

	return s.repo.Delete(ctx, groupID)
}

func (s *metodePembayaranGroupService) ToggleStatus(ctx context.Context, id string) (*models.MetodePembayaranGroupResponse, error) {
	groupID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID group tidak valid")
	}

	group, err := s.repo.ToggleStatus(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Group metode pembayaran tidak ditemukan")
		}
		return nil, err
	}

	return &models.MetodePembayaranGroupResponse{
		ID:        group.ID.String(),
		Nama:      group.Nama,
		IsActive:  group.IsActive,
		UpdatedAt: group.UpdatedAt,
	}, nil
}
