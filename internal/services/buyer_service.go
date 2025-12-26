package services

import (
	"context"
	"errors"
	"fmt"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type BuyerService interface {
	FindByID(ctx context.Context, id string) (*models.BuyerDetailResponse, error)
	FindAll(ctx context.Context, params *models.BuyerFilterRequest) ([]models.BuyerListResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateBuyerRequest) (*models.BuyerDetailResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	ResetPassword(ctx context.Context, id string, req *models.ResetBuyerPasswordRequest) error
	GetStatistik(ctx context.Context) (*models.BuyerStatistikResponse, error)
}

type buyerService struct {
	repo       repositories.BuyerRepository
	alamatRepo repositories.AlamatBuyerRepository
}

func NewBuyerService(repo repositories.BuyerRepository, alamatRepo repositories.AlamatBuyerRepository) BuyerService {
	return &buyerService{repo: repo, alamatRepo: alamatRepo}
}

func (s *buyerService) FindByID(ctx context.Context, id string) (*models.BuyerDetailResponse, error) {
	buyer, err := s.repo.FindByIDWithAlamat(ctx, id)
	if err != nil {
		return nil, errors.New("buyer tidak ditemukan")
	}
	return s.toDetailResponse(buyer), nil
}

func (s *buyerService) FindAll(ctx context.Context, params *models.BuyerFilterRequest) ([]models.BuyerListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	buyers, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.BuyerListResponse
	for _, b := range buyers {
		count, _ := s.repo.CountAlamat(ctx, b.ID.String())
		items = append(items, models.BuyerListResponse{
			ID:           b.ID.String(),
			Nama:         b.Nama,
			Username:     b.Username,
			Email:        b.Email,
			Telepon:      b.Telepon,
			IsActive:     b.IsActive,
			IsVerified:   b.IsVerified,
			JumlahAlamat: int(count),
			LastLoginAt:  b.LastLoginAt,
			CreatedAt:    b.CreatedAt,
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

func (s *buyerService) Update(ctx context.Context, id string, req *models.UpdateBuyerRequest) (*models.BuyerDetailResponse, error) {
	buyer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("buyer tidak ditemukan")
	}

	if req.Username != nil && *req.Username != buyer.Username {
		exists, _ := s.repo.ExistsByUsername(ctx, *req.Username, &id)
		if exists {
			return nil, errors.New("username sudah digunakan oleh buyer lain")
		}
		buyer.Username = *req.Username
	}

	if req.Email != nil && *req.Email != buyer.Email {
		exists, _ := s.repo.ExistsByEmail(ctx, *req.Email, &id)
		if exists {
			return nil, errors.New("email sudah digunakan oleh buyer lain")
		}
		buyer.Email = *req.Email
	}

	if req.Nama != nil {
		buyer.Nama = *req.Nama
	}
	if req.Telepon != nil {
		buyer.Telepon = req.Telepon
	}
	if req.IsActive != nil {
		buyer.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, buyer); err != nil {
		return nil, err
	}

	buyer, _ = s.repo.FindByIDWithAlamat(ctx, id)
	return s.toDetailResponse(buyer), nil
}

func (s *buyerService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("buyer tidak ditemukan")
	}
	return s.repo.Delete(ctx, id)
}

func (s *buyerService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	buyer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("buyer tidak ditemukan")
	}

	buyer.IsActive = !buyer.IsActive
	if err := s.repo.Update(ctx, buyer); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       buyer.ID.String(),
		IsActive: buyer.IsActive,
	}, nil
}

func (s *buyerService) ResetPassword(ctx context.Context, id string, req *models.ResetBuyerPasswordRequest) error {
	buyer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("buyer tidak ditemukan")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	buyer.Password = hashedPassword
	return s.repo.Update(ctx, buyer)
}

func (s *buyerService) GetStatistik(ctx context.Context) (*models.BuyerStatistikResponse, error) {
	return s.repo.GetStatistik(ctx)
}

func (s *buyerService) toDetailResponse(b *models.Buyer) *models.BuyerDetailResponse {
	var alamatResponses []models.AlamatBuyerResponse
	for _, a := range b.Alamat {
		alamatResponses = append(alamatResponses, s.toAlamatResponse(&a))
	}

	return &models.BuyerDetailResponse{
		ID:              b.ID.String(),
		Nama:            b.Nama,
		Username:        b.Username,
		Email:           b.Email,
		Telepon:         b.Telepon,
		IsActive:        b.IsActive,
		IsVerified:      b.IsVerified,
		EmailVerifiedAt: b.EmailVerifiedAt,
		LastLoginAt:     b.LastLoginAt,
		Alamat:          alamatResponses,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
	}
}

func (s *buyerService) toAlamatResponse(a *models.AlamatBuyer) models.AlamatBuyerResponse {
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

	return models.AlamatBuyerResponse{
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
