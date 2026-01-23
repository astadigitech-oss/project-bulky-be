package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/shopspring/decimal"
)

type InformasiPickupService interface {
	Get(ctx context.Context) (*models.InformasiPickupResponse, error)
	Update(ctx context.Context, req *models.UpdateInformasiPickupRequest) (*models.InformasiPickupResponse, error)
	UpdateJadwal(ctx context.Context, req *models.UpdateJadwalGudangRequest) (*models.InformasiPickupResponse, error)
	GetPublic(ctx context.Context) (*models.PublicInformasiPickupResponse, error)
}

type informasiPickupService struct {
	repo repositories.InformasiPickupRepository
}

func NewInformasiPickupService(repo repositories.InformasiPickupRepository) InformasiPickupService {
	return &informasiPickupService{repo: repo}
}

func (s *informasiPickupService) Get(ctx context.Context) (*models.InformasiPickupResponse, error) {
	pickup, err := s.repo.Get(ctx)
	if err != nil {
		return nil, errors.New("informasi pickup tidak ditemukan")
	}
	return s.toResponse(pickup), nil
}

func (s *informasiPickupService) Update(ctx context.Context, req *models.UpdateInformasiPickupRequest) (*models.InformasiPickupResponse, error) {
	pickup, err := s.repo.Get(ctx)
	if err != nil {
		return nil, errors.New("informasi pickup tidak ditemukan")
	}

	// Validate nomor whatsapp
	if len(req.NomorWhatsapp) < 10 || req.NomorWhatsapp[:2] != "62" {
		return nil, errors.New("nomor whatsapp harus diawali dengan 62")
	}

	pickup.Alamat = req.Alamat
	pickup.JamOperasional = req.JamOperasional
	pickup.NomorWhatsapp = req.NomorWhatsapp

	if req.Latitude != nil {
		lat := decimal.NewFromFloat(*req.Latitude)
		pickup.Latitude = &lat
	} else {
		pickup.Latitude = nil
	}

	if req.Longitude != nil {
		lon := decimal.NewFromFloat(*req.Longitude)
		pickup.Longitude = &lon
	} else {
		pickup.Longitude = nil
	}

	pickup.GoogleMapsURL = req.GoogleMapsURL

	if err := s.repo.Update(ctx, pickup); err != nil {
		return nil, err
	}

	// Reload with jadwal
	pickup, err = s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return s.toResponse(pickup), nil
}

func (s *informasiPickupService) UpdateJadwal(ctx context.Context, req *models.UpdateJadwalGudangRequest) (*models.InformasiPickupResponse, error) {
	pickup, err := s.repo.Get(ctx)
	if err != nil {
		return nil, errors.New("informasi pickup tidak ditemukan")
	}

	// Validate jadwal
	for _, j := range req.Jadwal {
		if j.IsBuka {
			if j.JamBuka == nil || j.JamTutup == nil {
				return nil, fmt.Errorf("jam buka dan jam tutup wajib diisi untuk hari %d", j.Hari)
			}
			// Validate time format and logic
			if *j.JamBuka >= *j.JamTutup {
				return nil, fmt.Errorf("jam tutup harus lebih besar dari jam buka untuk hari %d", j.Hari)
			}
		}
	}

	// Convert to JadwalGudang models
	jadwalModels := []models.JadwalGudang{}
	for _, j := range req.Jadwal {
		jadwalModels = append(jadwalModels, models.JadwalGudang{
			InformasiPickupID: pickup.ID,
			Hari:              j.Hari,
			JamBuka:           j.JamBuka,
			JamTutup:          j.JamTutup,
			IsBuka:            j.IsBuka,
		})
	}

	if err := s.repo.UpdateJadwal(ctx, pickup.ID.String(), jadwalModels); err != nil {
		return nil, err
	}

	// Reload with updated jadwal
	pickup, err = s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return s.toResponse(pickup), nil
}

func (s *informasiPickupService) GetPublic(ctx context.Context) (*models.PublicInformasiPickupResponse, error) {
	pickup, err := s.repo.Get(ctx)
	if err != nil {
		return nil, errors.New("informasi pickup tidak ditemukan")
	}

	isOpenNow := s.isOpenNow(pickup.JadwalGudang)

	return s.toPublicResponse(pickup, isOpenNow), nil
}

// Helper: Check if warehouse is currently open
func (s *informasiPickupService) isOpenNow(jadwal []models.JadwalGudang) bool {
	now := time.Now()
	currentDay := int(now.Weekday()) // 0=Sunday, 1=Monday, ...
	currentTime := now.Format("15:04")

	for _, j := range jadwal {
		if j.Hari == currentDay && j.IsBuka {
			if j.JamBuka != nil && j.JamTutup != nil {
				if currentTime >= *j.JamBuka && currentTime <= *j.JamTutup {
					return true
				}
			}
		}
	}

	return false
}

func (s *informasiPickupService) toResponse(p *models.InformasiPickup) *models.InformasiPickupResponse {
	var lat, lon *float64
	if p.Latitude != nil {
		f, _ := p.Latitude.Float64()
		lat = &f
	}
	if p.Longitude != nil {
		f, _ := p.Longitude.Float64()
		lon = &f
	}

	jadwalResp := []models.JadwalGudangResponse{}
	for _, j := range p.JadwalGudang {
		jadwalResp = append(jadwalResp, models.JadwalGudangResponse{
			ID:       j.ID.String(),
			Hari:     j.Hari,
			HariNama: j.GetHariNama(),
			JamBuka:  j.JamBuka,
			JamTutup: j.JamTutup,
			IsBuka:   j.IsBuka,
		})
	}

	return &models.InformasiPickupResponse{
		ID:             p.ID.String(),
		Alamat:         p.Alamat,
		JamOperasional: p.JamOperasional,
		NomorWhatsapp:  p.NomorWhatsapp,
		Latitude:       lat,
		Longitude:      lon,
		GoogleMapsURL:  p.GoogleMapsURL,
		JadwalGudang:   jadwalResp,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func (s *informasiPickupService) toPublicResponse(p *models.InformasiPickup, isOpenNow bool) *models.PublicInformasiPickupResponse {
	var lat, lon *float64
	if p.Latitude != nil {
		f, _ := p.Latitude.Float64()
		lat = &f
	}
	if p.Longitude != nil {
		f, _ := p.Longitude.Float64()
		lon = &f
	}

	jadwalResp := []models.JadwalGudangResponse{}
	for _, j := range p.JadwalGudang {
		jadwalResp = append(jadwalResp, models.JadwalGudangResponse{
			ID:       j.ID.String(),
			Hari:     j.Hari,
			HariNama: j.GetHariNama(),
			JamBuka:  j.JamBuka,
			JamTutup: j.JamTutup,
			IsBuka:   j.IsBuka,
		})
	}

	return &models.PublicInformasiPickupResponse{
		Alamat:         p.Alamat,
		JamOperasional: p.JamOperasional,
		NomorWhatsapp:  p.NomorWhatsapp,
		Latitude:       lat,
		Longitude:      lon,
		GoogleMapsURL:  p.GoogleMapsURL,
		IsOpenNow:      isOpenNow,
		JadwalGudang:   jadwalResp,
	}
}
