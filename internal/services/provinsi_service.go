package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type ProvinsiService interface {
	Create(ctx context.Context, req *models.CreateProvinsiRequest) (*models.ProvinsiDetailResponse, error)
	FindByID(ctx context.Context, id string) (*models.ProvinsiDetailResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.ProvinsiListResponse, *models.PaginationMeta, error)
	FindAllDropdown(ctx context.Context) ([]models.WilayahDropdownItem, error)
	Update(ctx context.Context, id string, req *models.UpdateProvinsiRequest) (*models.ProvinsiDetailResponse, error)
	Delete(ctx context.Context, id string) error
}

type provinsiService struct {
	repo repositories.ProvinsiRepository
}

func NewProvinsiService(repo repositories.ProvinsiRepository) ProvinsiService {
	return &provinsiService{repo: repo}
}

func (s *provinsiService) Create(ctx context.Context, req *models.CreateProvinsiRequest) (*models.ProvinsiDetailResponse, error) {
	if req.Kode != nil && *req.Kode != "" {
		exists, _ := s.repo.ExistsByKode(ctx, *req.Kode, nil)
		if exists {
			return nil, errors.New("kode provinsi sudah digunakan")
		}
	}

	provinsi := &models.Provinsi{
		ID:   uuid.New(),
		Nama: req.Nama,
		Kode: req.Kode,
	}

	if err := s.repo.Create(ctx, provinsi); err != nil {
		return nil, err
	}

	return s.toDetailResponse(provinsi), nil
}

func (s *provinsiService) FindByID(ctx context.Context, id string) (*models.ProvinsiDetailResponse, error) {
	provinsi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("provinsi tidak ditemukan")
	}
	return s.toDetailResponse(provinsi), nil
}

func (s *provinsiService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.ProvinsiListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	provinsis, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.ProvinsiListResponse
	for _, p := range provinsis {
		count, _ := s.repo.CountKota(ctx, p.ID.String())
		items = append(items, models.ProvinsiListResponse{
			ID:         p.ID.String(),
			Nama:       p.Nama,
			Kode:       s.getKode(p.Kode),
			JumlahKota: int(count),
		})
	}

	totalHalaman := (total + int64(params.PerHalaman) - 1) / int64(params.PerHalaman)

	meta := &models.PaginationMeta{
		Halaman:      params.Halaman,
		PerHalaman:   params.PerHalaman,
		TotalData:    total,
		TotalHalaman: totalHalaman,
	}

	return items, meta, nil
}

func (s *provinsiService) FindAllDropdown(ctx context.Context) ([]models.WilayahDropdownItem, error) {
	provinsis, err := s.repo.FindAllDropdown(ctx)
	if err != nil {
		return nil, err
	}

	var items []models.WilayahDropdownItem
	for _, p := range provinsis {
		items = append(items, models.WilayahDropdownItem{
			ID:   p.ID.String(),
			Nama: p.Nama,
			Kode: s.getKode(p.Kode),
		})
	}

	return items, nil
}

func (s *provinsiService) Update(ctx context.Context, id string, req *models.UpdateProvinsiRequest) (*models.ProvinsiDetailResponse, error) {
	provinsi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("provinsi tidak ditemukan")
	}

	if req.Kode != nil && *req.Kode != "" {
		exists, _ := s.repo.ExistsByKode(ctx, *req.Kode, &id)
		if exists {
			return nil, errors.New("kode provinsi sudah digunakan")
		}
		provinsi.Kode = req.Kode
	}

	if req.Nama != nil {
		provinsi.Nama = *req.Nama
	}

	if err := s.repo.Update(ctx, provinsi); err != nil {
		return nil, err
	}

	return s.toDetailResponse(provinsi), nil
}

func (s *provinsiService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("provinsi tidak ditemukan")
	}

	count, _ := s.repo.CountKota(ctx, id)
	if count > 0 {
		return errors.New("provinsi tidak dapat dihapus karena masih memiliki kota/kabupaten terkait")
	}

	return s.repo.Delete(ctx, id)
}

func (s *provinsiService) toDetailResponse(p *models.Provinsi) *models.ProvinsiDetailResponse {
	return &models.ProvinsiDetailResponse{
		ID:        p.ID.String(),
		Nama:      p.Nama,
		Kode:      s.getKode(p.Kode),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (s *provinsiService) getKode(kode *string) string {
	if kode == nil {
		return ""
	}
	return *kode
}
