package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"project-bulky-be/internal/config"
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
	GetChart(ctx context.Context, params *models.ChartParams) (*models.ChartResponse, error)
	UpdateProfile(ctx context.Context, id, nama, username, email, telepon string) (*models.Buyer, error)
	IsEmailExistExcludeID(ctx context.Context, email, excludeID string) (bool, error)
	IsUsernameExistExcludeID(ctx context.Context, username, excludeID string) (bool, error)
}

type buyerService struct {
	repo       repositories.BuyerRepository
	alamatRepo repositories.AlamatBuyerRepository
	cfg        *config.Config
}

func NewBuyerService(repo repositories.BuyerRepository, alamatRepo repositories.AlamatBuyerRepository) BuyerService {
	return &buyerService{
		repo:       repo,
		alamatRepo: alamatRepo,
		cfg:        config.LoadConfig(),
	}
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

	items := []models.BuyerListResponse{}
	for _, b := range buyers {
		items = append(items, models.BuyerListResponse{
			ID:        b.ID.String(),
			Nama:      b.Nama,
			Username:  b.Username,
			Email:     b.Email,
			Telepon:   b.Telepon,
			CreatedAt: b.CreatedAt,
		})
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
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

	// Use configured bcrypt cost
	hashedPassword, err := utils.HashPasswordWithCost(req.NewPassword, s.cfg.BcryptCost)
	if err != nil {
		return err
	}

	buyer.Password = hashedPassword
	return s.repo.Update(ctx, buyer)
}

func (s *buyerService) GetStatistik(ctx context.Context) (*models.BuyerStatistikResponse, error) {
	return s.repo.GetStatistik(ctx)
}

func (s *buyerService) GetChart(ctx context.Context, params *models.ChartParams) (*models.ChartResponse, error) {
	now := time.Now()
	var startDate, endDate time.Time
	var mode string

	// Default filter to "year" if empty
	if params.Filter == "" {
		params.Filter = "year"
	}

	switch params.Filter {
	case "year":
		mode = "year"
		year := params.Tahun
		if year == 0 {
			year = now.Year()
		}
		startDate = time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
		endDate = time.Date(year, 12, 31, 23, 59, 59, 0, now.Location())

	case "month":
		mode = "month"
		year := params.Tahun
		month := params.Bulan
		if year == 0 {
			year = now.Year()
		}
		if month == 0 {
			month = int(now.Month())
		}
		startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
		// Get last day of month
		endDate = startDate.AddDate(0, 1, -1)
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, now.Location())

	case "week":
		mode = "month"
		year := params.Tahun
		month := params.Bulan
		week := params.Minggu
		if year == 0 {
			year = now.Year()
		}
		if month == 0 {
			month = int(now.Month())
		}
		if week == 0 {
			week = 1
		}
		startDate = time.Date(year, time.Month(month), (week-1)*7+1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 0, 6)
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, now.Location())

	case "custom":
		mode = "month"
		if !params.TanggalDari.IsZero() && !params.TanggalSampai.IsZero() {
			startDate = params.TanggalDari
			endDate = params.TanggalSampai
		} else {
			// Default to current month if dates not provided
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			endDate = now
		}
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, now.Location())

	default:
		// Default: year
		mode = "year"
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
	}

	// Get data from repository
	var chart []models.ChartData
	var total int64
	var err error

	if mode == "year" {
		// Group by month
		chart, total, err = s.repo.GetRegistrationByMonth(ctx, startDate, endDate)
	} else {
		// Group by day
		chart, total, err = s.repo.GetRegistrationByDay(ctx, startDate, endDate)
	}

	if err != nil {
		return nil, err
	}

	return &models.ChartResponse{
		Mode:  mode,
		Chart: chart,
		Total: total,
	}, nil
}

func (s *buyerService) toDetailResponse(b *models.Buyer) *models.BuyerDetailResponse {
	var alamatResponses []models.AlamatBuyerResponse
	for _, a := range b.Alamat {
		alamatResponses = append(alamatResponses, s.toAlamatResponse(&a))
	}

	return &models.BuyerDetailResponse{
		ID:        b.ID.String(),
		Nama:      b.Nama,
		Username:  b.Username,
		Email:     b.Email,
		Telepon:   b.Telepon,
		FotoURL:   b.FotoURL,
		Alamat:    alamatResponses,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
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

// UpdateProfile updates buyer profile (nama, username, email, telepon)
func (s *buyerService) UpdateProfile(ctx context.Context, id, nama, username, email, telepon string) (*models.Buyer, error) {
	buyer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("buyer tidak ditemukan")
	}

	buyer.Nama = nama
	buyer.Username = username
	buyer.Email = email
	buyer.Telepon = &telepon

	if err := s.repo.Update(ctx, buyer); err != nil {
		return nil, err
	}

	return buyer, nil
}

// IsEmailExistExcludeID checks if email exists excluding specific ID
func (s *buyerService) IsEmailExistExcludeID(ctx context.Context, email, excludeID string) (bool, error) {
	return s.repo.ExistsByEmail(ctx, email, &excludeID)
}

// IsUsernameExistExcludeID checks if username exists excluding specific ID
func (s *buyerService) IsUsernameExistExcludeID(ctx context.Context, username, excludeID string) (bool, error) {
	return s.repo.ExistsByUsername(ctx, username, &excludeID)
}
