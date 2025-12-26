package services

import (
	"context"
	"errors"
	"fmt"

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
	SetDefault(ctx context.Context, id string) (*models.AlamatBuyerResponse, error)
}

type alamatBuyerService struct {
	repo      repositories.AlamatBuyerRepository
	buyerRepo repositories.BuyerRepository
}

func NewAlamatBuyerService(repo repositories.AlamatBuyerRepository, buyerRepo repositories.BuyerRepository) AlamatBuyerService {
	return &alamatBuyerService{repo: repo, buyerRepo: buyerRepo}
}

func (s *alamatBuyerService) Create(ctx context.Context, req *models.CreateAlamatBuyerRequest) (*models.AlamatBuyerResponse, error) {
	// Validate buyer exists
	_, err := s.buyerRepo.FindByID(ctx, req.BuyerID)
	if err != nil {
		return nil, errors.New("buyer tidak ditemukan")
	}

	buyerUUID, _ := uuid.Parse(req.BuyerID)

	alamat := &models.AlamatBuyer{
		ID:              uuid.New(),
		BuyerID:         buyerUUID,
		Label:           req.Label,
		NamaPenerima:    req.NamaPenerima,
		TeleponPenerima: req.TeleponPenerima,
		Provinsi:        req.Provinsi,
		Kota:            req.Kota,
		Kecamatan:       req.Kecamatan,
		Kelurahan:       req.Kelurahan,
		KodePos:         req.KodePos,
		AlamatLengkap:   req.AlamatLengkap,
		Catatan:         req.Catatan,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		GooglePlaceID:   req.GooglePlaceID,
		IsDefault:       req.IsDefault,
	}

	if err := s.repo.Create(ctx, alamat); err != nil {
		return nil, err
	}

	return s.toResponse(alamat), nil
}

func (s *alamatBuyerService) FindByID(ctx context.Context, id string) (*models.AlamatBuyerResponse, error) {
	alamat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan")
	}
	return s.toResponse(alamat), nil
}

func (s *alamatBuyerService) FindAll(ctx context.Context, params *models.AlamatBuyerFilterRequest) ([]models.AlamatBuyerResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	items, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var responses []models.AlamatBuyerResponse
	for _, item := range items {
		responses = append(responses, *s.toResponse(&item))
	}

	totalHalaman := (total + int64(params.PerHalaman) - 1) / int64(params.PerHalaman)

	meta := &models.PaginationMeta{
		Halaman:      params.Halaman,
		PerHalaman:   params.PerHalaman,
		TotalData:    total,
		TotalHalaman: totalHalaman,
	}

	return responses, meta, nil
}

func (s *alamatBuyerService) Update(ctx context.Context, id string, req *models.UpdateAlamatBuyerRequest) (*models.AlamatBuyerResponse, error) {
	alamat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan")
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
	if req.Provinsi != nil {
		alamat.Provinsi = *req.Provinsi
	}
	if req.Kota != nil {
		alamat.Kota = *req.Kota
	}
	if req.Kecamatan != nil {
		alamat.Kecamatan = req.Kecamatan
	}
	if req.Kelurahan != nil {
		alamat.Kelurahan = req.Kelurahan
	}
	if req.KodePos != nil {
		alamat.KodePos = req.KodePos
	}
	if req.AlamatLengkap != nil {
		alamat.AlamatLengkap = *req.AlamatLengkap
	}
	if req.Catatan != nil {
		alamat.Catatan = req.Catatan
	}
	if req.Latitude != nil {
		alamat.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		alamat.Longitude = req.Longitude
	}
	if req.GooglePlaceID != nil {
		alamat.GooglePlaceID = req.GooglePlaceID
	}
	if req.IsDefault != nil {
		alamat.IsDefault = *req.IsDefault
	}

	if err := s.repo.Update(ctx, alamat); err != nil {
		return nil, err
	}

	return s.toResponse(alamat), nil
}

func (s *alamatBuyerService) Delete(ctx context.Context, id string) error {
	alamat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("alamat tidak ditemukan")
	}

	// Check if trying to delete default address
	if alamat.IsDefault {
		count, _ := s.repo.CountByBuyerID(ctx, alamat.BuyerID.String())
		if count > 1 {
			return errors.New("tidak dapat menghapus alamat default. Set alamat lain sebagai default terlebih dahulu")
		}
	}

	return s.repo.Delete(ctx, id)
}

func (s *alamatBuyerService) SetDefault(ctx context.Context, id string) (*models.AlamatBuyerResponse, error) {
	alamat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan")
	}

	if err := s.repo.SetDefault(ctx, id, alamat.BuyerID.String()); err != nil {
		return nil, err
	}

	// Refresh data
	alamat, _ = s.repo.FindByID(ctx, id)
	return s.toResponse(alamat), nil
}

func (s *alamatBuyerService) toResponse(a *models.AlamatBuyer) *models.AlamatBuyerResponse {
	// Build formatted address
	alamatFormatted := a.AlamatLengkap
	if a.Kelurahan != nil && *a.Kelurahan != "" {
		alamatFormatted = fmt.Sprintf("%s, %s", alamatFormatted, *a.Kelurahan)
	}
	if a.Kecamatan != nil && *a.Kecamatan != "" {
		alamatFormatted = fmt.Sprintf("%s, %s", alamatFormatted, *a.Kecamatan)
	}
	alamatFormatted = fmt.Sprintf("%s, %s, %s", alamatFormatted, a.Kota, a.Provinsi)
	if a.KodePos != nil && *a.KodePos != "" {
		alamatFormatted = fmt.Sprintf("%s %s", alamatFormatted, *a.KodePos)
	}

	return &models.AlamatBuyerResponse{
		ID:              a.ID.String(),
		BuyerID:         a.BuyerID.String(),
		Label:           a.Label,
		NamaPenerima:    a.NamaPenerima,
		TeleponPenerima: a.TeleponPenerima,
		Provinsi:        a.Provinsi,
		Kota:            a.Kota,
		Kecamatan:       a.Kecamatan,
		Kelurahan:       a.Kelurahan,
		KodePos:         a.KodePos,
		AlamatLengkap:   a.AlamatLengkap,
		AlamatFormatted: alamatFormatted,
		Catatan:         a.Catatan,
		Latitude:        a.Latitude,
		Longitude:       a.Longitude,
		GooglePlaceID:   a.GooglePlaceID,
		IsDefault:       a.IsDefault,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}
