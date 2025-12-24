package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type KecamatanService interface {
	Create(ctx context.Context, req *models.CreateKecamatanRequest) (*models.KecamatanDetailResponse, error)
	FindByID(ctx context.Context, id string) (*models.KecamatanDetailResponse, error)
	FindAll(ctx context.Context, params *models.KecamatanFilterRequest) ([]models.KecamatanListResponse, *models.PaginationMeta, error)
	FindByKotaID(ctx context.Context, kotaID string) ([]models.WilayahDropdownItem, error)
	Update(ctx context.Context, id string, req *models.UpdateKecamatanRequest) (*models.KecamatanDetailResponse, error)
	Delete(ctx context.Context, id string) error
}

type kecamatanService struct {
	repo     repositories.KecamatanRepository
	kotaRepo repositories.KotaRepository
}

func NewKecamatanService(repo repositories.KecamatanRepository, kotaRepo repositories.KotaRepository) KecamatanService {
	return &kecamatanService{repo: repo, kotaRepo: kotaRepo}
}

func (s *kecamatanService) Create(ctx context.Context, req *models.CreateKecamatanRequest) (*models.KecamatanDetailResponse, error) {
	_, err := s.kotaRepo.FindByID(ctx, req.KotaID)
	if err != nil {
		return nil, errors.New("kota tidak ditemukan")
	}

	if req.Kode != nil && *req.Kode != "" {
		exists, _ := s.repo.ExistsByKode(ctx, *req.Kode, nil)
		if exists {
			return nil, errors.New("kode kecamatan sudah digunakan")
		}
	}

	kotaID, _ := uuid.Parse(req.KotaID)
	kecamatan := &models.Kecamatan{
		ID:     uuid.New(),
		KotaID: kotaID,
		Nama:   req.Nama,
		Kode:   req.Kode,
	}

	if err := s.repo.Create(ctx, kecamatan); err != nil {
		return nil, err
	}

	kecamatan, _ = s.repo.FindByID(ctx, kecamatan.ID.String())
	return s.toDetailResponse(kecamatan), nil
}

func (s *kecamatanService) FindByID(ctx context.Context, id string) (*models.KecamatanDetailResponse, error) {
	kecamatan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kecamatan tidak ditemukan")
	}
	return s.toDetailResponse(kecamatan), nil
}

func (s *kecamatanService) FindAll(ctx context.Context, params *models.KecamatanFilterRequest) ([]models.KecamatanListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	kecamatans, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.KecamatanListResponse
	for _, k := range kecamatans {
		count, _ := s.repo.CountKelurahan(ctx, k.ID.String())
		items = append(items, models.KecamatanListResponse{
			ID:   k.ID.String(),
			Nama: k.Nama,
			Kode: s.getKode(k.Kode),
			Kota: models.KotaSimple{
				ID:   k.Kota.ID.String(),
				Nama: k.Kota.Nama,
				Kode: s.getKode(k.Kota.Kode),
				Provinsi: models.ProvinsiSimple{
					ID:   k.Kota.Provinsi.ID.String(),
					Nama: k.Kota.Provinsi.Nama,
					Kode: s.getKode(k.Kota.Provinsi.Kode),
				},
			},
			JumlahKelurahan: int(count),
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

func (s *kecamatanService) FindByKotaID(ctx context.Context, kotaID string) ([]models.WilayahDropdownItem, error) {
	kecamatans, err := s.repo.FindByKotaID(ctx, kotaID)
	if err != nil {
		return nil, err
	}

	var items []models.WilayahDropdownItem
	for _, k := range kecamatans {
		items = append(items, models.WilayahDropdownItem{
			ID:   k.ID.String(),
			Nama: k.Nama,
			Kode: s.getKode(k.Kode),
		})
	}

	return items, nil
}

func (s *kecamatanService) Update(ctx context.Context, id string, req *models.UpdateKecamatanRequest) (*models.KecamatanDetailResponse, error) {
	kecamatan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kecamatan tidak ditemukan")
	}

	if req.KotaID != nil {
		_, err := s.kotaRepo.FindByID(ctx, *req.KotaID)
		if err != nil {
			return nil, errors.New("kota tidak ditemukan")
		}
		kotaID, _ := uuid.Parse(*req.KotaID)
		kecamatan.KotaID = kotaID
	}

	if req.Kode != nil && *req.Kode != "" {
		exists, _ := s.repo.ExistsByKode(ctx, *req.Kode, &id)
		if exists {
			return nil, errors.New("kode kecamatan sudah digunakan")
		}
		kecamatan.Kode = req.Kode
	}

	if req.Nama != nil {
		kecamatan.Nama = *req.Nama
	}

	if err := s.repo.Update(ctx, kecamatan); err != nil {
		return nil, err
	}

	kecamatan, _ = s.repo.FindByID(ctx, id)
	return s.toDetailResponse(kecamatan), nil
}

func (s *kecamatanService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("kecamatan tidak ditemukan")
	}

	count, _ := s.repo.CountKelurahan(ctx, id)
	if count > 0 {
		return errors.New("kecamatan tidak dapat dihapus karena masih memiliki kelurahan terkait")
	}

	return s.repo.Delete(ctx, id)
}

func (s *kecamatanService) toDetailResponse(k *models.Kecamatan) *models.KecamatanDetailResponse {
	return &models.KecamatanDetailResponse{
		ID:   k.ID.String(),
		Nama: k.Nama,
		Kode: s.getKode(k.Kode),
		Kota: models.KotaSimple{
			ID:   k.Kota.ID.String(),
			Nama: k.Kota.Nama,
			Kode: s.getKode(k.Kota.Kode),
			Provinsi: models.ProvinsiSimple{
				ID:   k.Kota.Provinsi.ID.String(),
				Nama: k.Kota.Provinsi.Nama,
				Kode: s.getKode(k.Kota.Provinsi.Kode),
			},
		},
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func (s *kecamatanService) getKode(kode *string) string {
	if kode == nil {
		return ""
	}
	return *kode
}
