package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type AlamatBuyerService interface {
	Create(ctx context.Context, req *models.CreateAlamatBuyerRequest) (*models.AlamatBuyerResponse, error)
	FindByID(ctx context.Context, id string) (*models.AlamatBuyerResponse, error)
	FindAll(ctx context.Context, params *models.AlamatBuyerFilterRequest) ([]models.AlamatBuyerResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateAlamatBuyerRequest) (*models.AlamatBuyerResponse, error)
	Delete(ctx context.Context, id string) error
	SetDefault(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
}

type alamatBuyerService struct {
	repo          repositories.AlamatBuyerRepository
	buyerRepo     repositories.BuyerRepository
	kelurahanRepo repositories.KelurahanRepository
}

func NewAlamatBuyerService(
	repo repositories.AlamatBuyerRepository,
	buyerRepo repositories.BuyerRepository,
	kelurahanRepo repositories.KelurahanRepository,
) AlamatBuyerService {
	return &alamatBuyerService{repo: repo, buyerRepo: buyerRepo, kelurahanRepo: kelurahanRepo}
}

func (s *alamatBuyerService) Create(ctx context.Context, req *models.CreateAlamatBuyerRequest) (*models.AlamatBuyerResponse, error) {
	_, err := s.buyerRepo.FindByID(ctx, req.BuyerID)
	if err != nil {
		return nil, errors.New("buyer tidak ditemukan")
	}

	_, err = s.kelurahanRepo.FindByID(ctx, req.KelurahanID)
	if err != nil {
		return nil, errors.New("kelurahan tidak ditemukan")
	}

	buyerID, _ := uuid.Parse(req.BuyerID)
	kelurahanID, _ := uuid.Parse(req.KelurahanID)

	// Check if first address
	count, _ := s.repo.CountByBuyerID(ctx, req.BuyerID)
	isDefault := req.IsDefault || count == 0

	// If setting as default, unset others
	if isDefault {
		s.repo.UnsetDefaultByBuyerID(ctx, req.BuyerID, nil)
	}

	alamat := &models.AlamatBuyer{
		ID:              uuid.New(),
		BuyerID:         buyerID,
		KelurahanID:     kelurahanID,
		Label:           req.Label,
		NamaPenerima:    req.NamaPenerima,
		TeleponPenerima: req.TeleponPenerima,
		KodePos:         req.KodePos,
		AlamatLengkap:   req.AlamatLengkap,
		Catatan:         req.Catatan,
		IsDefault:       isDefault,
	}

	if err := s.repo.Create(ctx, alamat); err != nil {
		return nil, err
	}

	alamat, _ = s.repo.FindByIDWithWilayah(ctx, alamat.ID.String())
	return s.toResponse(alamat), nil
}

func (s *alamatBuyerService) FindByID(ctx context.Context, id string) (*models.AlamatBuyerResponse, error) {
	alamat, err := s.repo.FindByIDWithWilayah(ctx, id)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan")
	}
	return s.toResponse(alamat), nil
}

func (s *alamatBuyerService) FindAll(ctx context.Context, params *models.AlamatBuyerFilterRequest) ([]models.AlamatBuyerResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	alamats, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.AlamatBuyerResponse
	for _, a := range alamats {
		items = append(items, *s.toResponse(&a))
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

func (s *alamatBuyerService) Update(ctx context.Context, id string, req *models.UpdateAlamatBuyerRequest) (*models.AlamatBuyerResponse, error) {
	alamat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan")
	}

	if req.KelurahanID != nil {
		_, err = s.kelurahanRepo.FindByID(ctx, *req.KelurahanID)
		if err != nil {
			return nil, errors.New("kelurahan tidak ditemukan")
		}
		kelurahanID, _ := uuid.Parse(*req.KelurahanID)
		alamat.KelurahanID = kelurahanID
	}

	if req.Label != nil {
		alamat.Label = *req.Label
	}
	if req.NamaPenerima != nil {
		alamat.NamaPenerima = *req.NamaPenerima
	}
	if req.TeleponPenerima != nil {
		alamat.TeleponPenerima = *req.TeleponPenerima
	}
	if req.KodePos != nil {
		alamat.KodePos = *req.KodePos
	}
	if req.AlamatLengkap != nil {
		alamat.AlamatLengkap = *req.AlamatLengkap
	}
	if req.Catatan != nil {
		alamat.Catatan = req.Catatan
	}
	if req.IsDefault != nil && *req.IsDefault {
		s.repo.UnsetDefaultByBuyerID(ctx, alamat.BuyerID.String(), &id)
		alamat.IsDefault = true
	}

	if err := s.repo.Update(ctx, alamat); err != nil {
		return nil, err
	}

	alamat, _ = s.repo.FindByIDWithWilayah(ctx, id)
	return s.toResponse(alamat), nil
}

func (s *alamatBuyerService) Delete(ctx context.Context, id string) error {
	alamat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("alamat tidak ditemukan")
	}

	if alamat.IsDefault {
		hasOther, _ := s.repo.HasOtherAddresses(ctx, alamat.BuyerID.String(), id)
		if hasOther {
			return errors.New("tidak dapat menghapus alamat default. Silakan set alamat lain sebagai default terlebih dahulu")
		}
	}

	return s.repo.Delete(ctx, id)
}

func (s *alamatBuyerService) SetDefault(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	alamat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan")
	}

	s.repo.UnsetDefaultByBuyerID(ctx, alamat.BuyerID.String(), &id)
	alamat.IsDefault = true

	if err := s.repo.Update(ctx, alamat); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       alamat.ID.String(),
		IsActive: true,
	}, nil
}

func (s *alamatBuyerService) toResponse(a *models.AlamatBuyer) *models.AlamatBuyerResponse {
	kode := func(k *string) string {
		if k == nil {
			return ""
		}
		return *k
	}

	wilayah := models.WilayahResponse{
		Kelurahan: models.WilayahItemResponse{
			ID:   a.Kelurahan.ID.String(),
			Nama: a.Kelurahan.Nama,
			Kode: kode(a.Kelurahan.Kode),
		},
		Kecamatan: models.WilayahItemResponse{
			ID:   a.Kelurahan.Kecamatan.ID.String(),
			Nama: a.Kelurahan.Kecamatan.Nama,
			Kode: kode(a.Kelurahan.Kecamatan.Kode),
		},
		Kota: models.WilayahItemResponse{
			ID:   a.Kelurahan.Kecamatan.Kota.ID.String(),
			Nama: a.Kelurahan.Kecamatan.Kota.Nama,
			Kode: kode(a.Kelurahan.Kecamatan.Kota.Kode),
		},
		Provinsi: models.WilayahItemResponse{
			ID:   a.Kelurahan.Kecamatan.Kota.Provinsi.ID.String(),
			Nama: a.Kelurahan.Kecamatan.Kota.Provinsi.Nama,
			Kode: kode(a.Kelurahan.Kecamatan.Kota.Provinsi.Kode),
		},
	}

	formatted := a.AlamatLengkap + ", " + a.Kelurahan.Nama + ", " +
		a.Kelurahan.Kecamatan.Nama + ", " + a.Kelurahan.Kecamatan.Kota.Nama + ", " +
		a.Kelurahan.Kecamatan.Kota.Provinsi.Nama + " " + a.KodePos

	return &models.AlamatBuyerResponse{
		ID:              a.ID.String(),
		BuyerID:         a.BuyerID.String(),
		Label:           a.Label,
		NamaPenerima:    a.NamaPenerima,
		TeleponPenerima: a.TeleponPenerima,
		Wilayah:         wilayah,
		KodePos:         a.KodePos,
		AlamatLengkap:   a.AlamatLengkap,
		AlamatFormatted: formatted,
		Catatan:         a.Catatan,
		IsDefault:       a.IsDefault,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}
