package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type KotaService interface {
	Create(ctx context.Context, req *models.CreateKotaRequest) (*models.KotaDetailResponse, error)
	FindByID(ctx context.Context, id string) (*models.KotaDetailResponse, error)
	FindAll(ctx context.Context, params *models.KotaFilterRequest) ([]models.KotaListResponse, *models.PaginationMeta, error)
	FindByProvinsiID(ctx context.Context, provinsiID string) ([]models.WilayahDropdownItem, error)
	Update(ctx context.Context, id string, req *models.UpdateKotaRequest) (*models.KotaDetailResponse, error)
	Delete(ctx context.Context, id string) error
}

type kotaService struct {
	repo         repositories.KotaRepository
	provinsiRepo repositories.ProvinsiRepository
}

func NewKotaService(repo repositories.KotaRepository, provinsiRepo repositories.ProvinsiRepository) KotaService {
	return &kotaService{repo: repo, provinsiRepo: provinsiRepo}
}

func (s *kotaService) Create(ctx context.Context, req *models.CreateKotaRequest) (*models.KotaDetailResponse, error) {
	// Validate provinsi exists
	_, err := s.provinsiRepo.FindByID(ctx, req.ProvinsiID)
	if err != nil {
		return nil, errors.New("provinsi tidak ditemukan")
	}

	if req.Kode != nil && *req.Kode != "" {
		exists, _ := s.repo.ExistsByKode(ctx, *req.Kode, nil)
		if exists {
			return nil, errors.New("kode kota sudah digunakan")
		}
	}

	provinsiID, _ := uuid.Parse(req.ProvinsiID)
	kota := &models.Kota{
		ID:         uuid.New(),
		ProvinsiID: provinsiID,
		Nama:       req.Nama,
		Kode:       req.Kode,
	}

	if err := s.repo.Create(ctx, kota); err != nil {
		return nil, err
	}

	kota, _ = s.repo.FindByID(ctx, kota.ID.String())
	return s.toDetailResponse(kota), nil
}

func (s *kotaService) FindByID(ctx context.Context, id string) (*models.KotaDetailResponse, error) {
	kota, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kota tidak ditemukan")
	}
	return s.toDetailResponse(kota), nil
}

func (s *kotaService) FindAll(ctx context.Context, params *models.KotaFilterRequest) ([]models.KotaListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	kotas, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.KotaListResponse
	for _, k := range kotas {
		count, _ := s.repo.CountKecamatan(ctx, k.ID.String())
		items = append(items, models.KotaListResponse{
			ID:   k.ID.String(),
			Nama: k.Nama,
			Kode: s.getKode(k.Kode),
			Provinsi: models.ProvinsiSimple{
				ID:   k.Provinsi.ID.String(),
				Nama: k.Provinsi.Nama,
				Kode: s.getKode(k.Provinsi.Kode),
			},
			JumlahKecamatan: int(count),
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

func (s *kotaService) FindByProvinsiID(ctx context.Context, provinsiID string) ([]models.WilayahDropdownItem, error) {
	kotas, err := s.repo.FindByProvinsiID(ctx, provinsiID)
	if err != nil {
		return nil, err
	}

	var items []models.WilayahDropdownItem
	for _, k := range kotas {
		items = append(items, models.WilayahDropdownItem{
			ID:   k.ID.String(),
			Nama: k.Nama,
			Kode: s.getKode(k.Kode),
		})
	}

	return items, nil
}

func (s *kotaService) Update(ctx context.Context, id string, req *models.UpdateKotaRequest) (*models.KotaDetailResponse, error) {
	kota, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kota tidak ditemukan")
	}

	if req.ProvinsiID != nil {
		_, err := s.provinsiRepo.FindByID(ctx, *req.ProvinsiID)
		if err != nil {
			return nil, errors.New("provinsi tidak ditemukan")
		}
		provinsiID, _ := uuid.Parse(*req.ProvinsiID)
		kota.ProvinsiID = provinsiID
	}

	if req.Kode != nil && *req.Kode != "" {
		exists, _ := s.repo.ExistsByKode(ctx, *req.Kode, &id)
		if exists {
			return nil, errors.New("kode kota sudah digunakan")
		}
		kota.Kode = req.Kode
	}

	if req.Nama != nil {
		kota.Nama = *req.Nama
	}

	if err := s.repo.Update(ctx, kota); err != nil {
		return nil, err
	}

	kota, _ = s.repo.FindByID(ctx, id)
	return s.toDetailResponse(kota), nil
}

func (s *kotaService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("kota tidak ditemukan")
	}

	count, _ := s.repo.CountKecamatan(ctx, id)
	if count > 0 {
		return errors.New("kota tidak dapat dihapus karena masih memiliki kecamatan terkait")
	}

	return s.repo.Delete(ctx, id)
}

func (s *kotaService) toDetailResponse(k *models.Kota) *models.KotaDetailResponse {
	return &models.KotaDetailResponse{
		ID:   k.ID.String(),
		Nama: k.Nama,
		Kode: s.getKode(k.Kode),
		Provinsi: models.ProvinsiSimple{
			ID:   k.Provinsi.ID.String(),
			Nama: k.Provinsi.Nama,
			Kode: s.getKode(k.Provinsi.Kode),
		},
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func (s *kotaService) getKode(kode *string) string {
	if kode == nil {
		return ""
	}
	return *kode
}
