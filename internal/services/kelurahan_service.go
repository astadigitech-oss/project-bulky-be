package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type KelurahanService interface {
	Create(ctx context.Context, req *models.CreateKelurahanRequest) (*models.KelurahanDetailResponse, error)
	FindByID(ctx context.Context, id string) (*models.KelurahanDetailResponse, error)
	FindAll(ctx context.Context, params *models.KelurahanFilterRequest) ([]models.KelurahanListResponse, *models.PaginationMeta, error)
	FindByKecamatanID(ctx context.Context, kecamatanID string) ([]models.WilayahDropdownItem, error)
	Update(ctx context.Context, id string, req *models.UpdateKelurahanRequest) (*models.KelurahanDetailResponse, error)
	Delete(ctx context.Context, id string) error
}

type kelurahanService struct {
	repo          repositories.KelurahanRepository
	kecamatanRepo repositories.KecamatanRepository
}

func NewKelurahanService(repo repositories.KelurahanRepository, kecamatanRepo repositories.KecamatanRepository) KelurahanService {
	return &kelurahanService{repo: repo, kecamatanRepo: kecamatanRepo}
}

func (s *kelurahanService) Create(ctx context.Context, req *models.CreateKelurahanRequest) (*models.KelurahanDetailResponse, error) {
	_, err := s.kecamatanRepo.FindByID(ctx, req.KecamatanID)
	if err != nil {
		return nil, errors.New("kecamatan tidak ditemukan")
	}

	if req.Kode != nil && *req.Kode != "" {
		exists, _ := s.repo.ExistsByKode(ctx, *req.Kode, nil)
		if exists {
			return nil, errors.New("kode kelurahan sudah digunakan")
		}
	}

	kecamatanID, _ := uuid.Parse(req.KecamatanID)
	kelurahan := &models.Kelurahan{
		ID:          uuid.New(),
		KecamatanID: kecamatanID,
		Nama:        req.Nama,
		Kode:        req.Kode,
	}

	if err := s.repo.Create(ctx, kelurahan); err != nil {
		return nil, err
	}

	kelurahan, _ = s.repo.FindByID(ctx, kelurahan.ID.String())
	return s.toDetailResponse(kelurahan), nil
}

func (s *kelurahanService) FindByID(ctx context.Context, id string) (*models.KelurahanDetailResponse, error) {
	kelurahan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kelurahan tidak ditemukan")
	}
	return s.toDetailResponse(kelurahan), nil
}

func (s *kelurahanService) FindAll(ctx context.Context, params *models.KelurahanFilterRequest) ([]models.KelurahanListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	kelurahans, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.KelurahanListResponse
	for _, k := range kelurahans {
		items = append(items, s.toListResponse(&k))
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

func (s *kelurahanService) FindByKecamatanID(ctx context.Context, kecamatanID string) ([]models.WilayahDropdownItem, error) {
	kelurahans, err := s.repo.FindByKecamatanID(ctx, kecamatanID)
	if err != nil {
		return nil, err
	}

	var items []models.WilayahDropdownItem
	for _, k := range kelurahans {
		items = append(items, models.WilayahDropdownItem{
			ID:   k.ID.String(),
			Nama: k.Nama,
			Kode: s.getKode(k.Kode),
		})
	}

	return items, nil
}

func (s *kelurahanService) Update(ctx context.Context, id string, req *models.UpdateKelurahanRequest) (*models.KelurahanDetailResponse, error) {
	kelurahan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kelurahan tidak ditemukan")
	}

	if req.KecamatanID != nil {
		_, err := s.kecamatanRepo.FindByID(ctx, *req.KecamatanID)
		if err != nil {
			return nil, errors.New("kecamatan tidak ditemukan")
		}
		kecamatanID, _ := uuid.Parse(*req.KecamatanID)
		kelurahan.KecamatanID = kecamatanID
	}

	if req.Kode != nil && *req.Kode != "" {
		exists, _ := s.repo.ExistsByKode(ctx, *req.Kode, &id)
		if exists {
			return nil, errors.New("kode kelurahan sudah digunakan")
		}
		kelurahan.Kode = req.Kode
	}

	if req.Nama != nil {
		kelurahan.Nama = *req.Nama
	}

	if err := s.repo.Update(ctx, kelurahan); err != nil {
		return nil, err
	}

	kelurahan, _ = s.repo.FindByID(ctx, id)
	return s.toDetailResponse(kelurahan), nil
}

func (s *kelurahanService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("kelurahan tidak ditemukan")
	}

	count, _ := s.repo.CountAlamatBuyer(ctx, id)
	if count > 0 {
		return errors.New("kelurahan tidak dapat dihapus karena masih digunakan di alamat buyer")
	}

	return s.repo.Delete(ctx, id)
}

func (s *kelurahanService) toDetailResponse(k *models.Kelurahan) *models.KelurahanDetailResponse {
	return &models.KelurahanDetailResponse{
		ID:   k.ID.String(),
		Nama: k.Nama,
		Kode: s.getKode(k.Kode),
		Kecamatan: models.KecamatanSimple{
			ID:   k.Kecamatan.ID.String(),
			Nama: k.Kecamatan.Nama,
			Kode: s.getKode(k.Kecamatan.Kode),
			Kota: models.KotaSimple{
				ID:   k.Kecamatan.Kota.ID.String(),
				Nama: k.Kecamatan.Kota.Nama,
				Kode: s.getKode(k.Kecamatan.Kota.Kode),
				Provinsi: models.ProvinsiSimple{
					ID:   k.Kecamatan.Kota.Provinsi.ID.String(),
					Nama: k.Kecamatan.Kota.Provinsi.Nama,
					Kode: s.getKode(k.Kecamatan.Kota.Provinsi.Kode),
				},
			},
		},
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func (s *kelurahanService) toListResponse(k *models.Kelurahan) models.KelurahanListResponse {
	return models.KelurahanListResponse{
		ID:   k.ID.String(),
		Nama: k.Nama,
		Kode: s.getKode(k.Kode),
		Kecamatan: models.KecamatanSimple{
			ID:   k.Kecamatan.ID.String(),
			Nama: k.Kecamatan.Nama,
			Kode: s.getKode(k.Kecamatan.Kode),
			Kota: models.KotaSimple{
				ID:   k.Kecamatan.Kota.ID.String(),
				Nama: k.Kecamatan.Kota.Nama,
				Kode: s.getKode(k.Kecamatan.Kota.Kode),
				Provinsi: models.ProvinsiSimple{
					ID:   k.Kecamatan.Kota.Provinsi.ID.String(),
					Nama: k.Kecamatan.Kota.Provinsi.Nama,
					Kode: s.getKode(k.Kecamatan.Kota.Provinsi.Kode),
				},
			},
		},
	}
}

func (s *kelurahanService) getKode(kode *string) string {
	if kode == nil {
		return ""
	}
	return *kode
}
