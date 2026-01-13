package services

import (
	"context"
	"errors"
	"time"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type DiskonKategoriService interface {
	Create(ctx context.Context, req *models.CreateDiskonKategoriRequest) (*models.DiskonKategoriResponse, error)
	FindByID(ctx context.Context, id string) (*models.DiskonKategoriResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest, kategoriID string, berlakuHariIni bool) ([]models.DiskonKategoriResponse, *models.PaginationMeta, error)
	FindActiveByKategoriID(ctx context.Context, kategoriID string) (*models.DiskonKategoriActiveResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateDiskonKategoriRequest) (*models.DiskonKategoriResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
}

type diskonKategoriService struct {
	repo repositories.DiskonKategoriRepository
}

func NewDiskonKategoriService(repo repositories.DiskonKategoriRepository) DiskonKategoriService {
	return &diskonKategoriService{repo: repo}
}

func (s *diskonKategoriService) Create(ctx context.Context, req *models.CreateDiskonKategoriRequest) (*models.DiskonKategoriResponse, error) {
	kategoriUUID, err := uuid.Parse(req.KategoriID)
	if err != nil {
		return nil, errors.New("kategori_id tidak valid")
	}

	diskon := &models.DiskonKategori{
		KategoriID:       kategoriUUID,
		PersentaseDiskon: req.PersentaseDiskon,
		NominalDiskon:    req.NominalDiskon,
		IsActive:         true,
	}

	if req.TanggalMulai != nil {
		t, _ := time.Parse("2006-01-02", *req.TanggalMulai)
		diskon.TanggalMulai = &t
	}
	if req.TanggalSelesai != nil {
		t, _ := time.Parse("2006-01-02", *req.TanggalSelesai)
		diskon.TanggalSelesai = &t
	}

	if err := s.repo.Create(ctx, diskon); err != nil {
		return nil, err
	}

	return s.FindByID(ctx, diskon.ID.String())
}

func (s *diskonKategoriService) FindByID(ctx context.Context, id string) (*models.DiskonKategoriResponse, error) {
	diskon, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("diskon kategori tidak ditemukan")
	}
	return s.toResponse(diskon), nil
}

func (s *diskonKategoriService) FindAll(ctx context.Context, params *models.PaginationRequest, kategoriID string, berlakuHariIni bool) ([]models.DiskonKategoriResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	diskons, total, err := s.repo.FindAll(ctx, params, kategoriID, berlakuHariIni)
	if err != nil {
		return nil, nil, err
	}

	items := []models.DiskonKategoriResponse{}
	for _, d := range diskons {
		items = append(items, *s.toResponse(&d))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *diskonKategoriService) FindActiveByKategoriID(ctx context.Context, kategoriID string) (*models.DiskonKategoriActiveResponse, error) {
	_, err := uuid.Parse(kategoriID)
	if err != nil {
		return nil, errors.New("kategori_id tidak valid")
	}

	diskon, err := s.repo.FindActiveByKategoriID(ctx, kategoriID)
	if err != nil {
		return nil, nil // No active discount found
	}

	var tanggalMulai, tanggalSelesai *string
	if diskon.TanggalMulai != nil {
		t := diskon.TanggalMulai.Format("2006-01-02")
		tanggalMulai = &t
	}
	if diskon.TanggalSelesai != nil {
		t := diskon.TanggalSelesai.Format("2006-01-02")
		tanggalSelesai = &t
	}

	return &models.DiskonKategoriActiveResponse{
		KategoriID:       diskon.KategoriID.String(),
		PersentaseDiskon: diskon.PersentaseDiskon,
		NominalDiskon:    diskon.NominalDiskon,
		TanggalMulai:     tanggalMulai,
		TanggalSelesai:   tanggalSelesai,
	}, nil
}

func (s *diskonKategoriService) Update(ctx context.Context, id string, req *models.UpdateDiskonKategoriRequest) (*models.DiskonKategoriResponse, error) {
	diskon, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("diskon kategori tidak ditemukan")
	}

	if req.KategoriID != nil {
		kategoriUUID, err := uuid.Parse(*req.KategoriID)
		if err != nil {
			return nil, errors.New("kategori_id tidak valid")
		}
		diskon.KategoriID = kategoriUUID
	}
	if req.PersentaseDiskon != nil {
		diskon.PersentaseDiskon = *req.PersentaseDiskon
	}
	if req.NominalDiskon != nil {
		diskon.NominalDiskon = *req.NominalDiskon
	}
	if req.TanggalMulai != nil {
		t, _ := time.Parse("2006-01-02", *req.TanggalMulai)
		diskon.TanggalMulai = &t
	}
	if req.TanggalSelesai != nil {
		t, _ := time.Parse("2006-01-02", *req.TanggalSelesai)
		diskon.TanggalSelesai = &t
	}
	if req.IsActive != nil {
		diskon.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, diskon); err != nil {
		return nil, err
	}

	return s.FindByID(ctx, id)
}

func (s *diskonKategoriService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("diskon kategori tidak ditemukan")
	}

	return s.repo.Delete(ctx, id)
}

func (s *diskonKategoriService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	diskon, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("diskon kategori tidak ditemukan")
	}

	diskon.IsActive = !diskon.IsActive
	if err := s.repo.Update(ctx, diskon); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       diskon.ID.String(),
		IsActive: diskon.IsActive,
	}, nil
}

func (s *diskonKategoriService) toResponse(d *models.DiskonKategori) *models.DiskonKategoriResponse {
	var tanggalMulai, tanggalSelesai *string
	if d.TanggalMulai != nil {
		t := d.TanggalMulai.Format("2006-01-02")
		tanggalMulai = &t
	}
	if d.TanggalSelesai != nil {
		t := d.TanggalSelesai.Format("2006-01-02")
		tanggalSelesai = &t
	}

	return &models.DiskonKategoriResponse{
		ID: d.ID.String(),
		Kategori: models.DiskonKategoriKategoriInfo{
			ID:   d.Kategori.ID.String(),
			Nama: d.Kategori.Nama,
			Slug: d.Kategori.Slug,
		},
		PersentaseDiskon: d.PersentaseDiskon,
		NominalDiskon:    d.NominalDiskon,
		TanggalMulai:     tanggalMulai,
		TanggalSelesai:   tanggalSelesai,
		IsActive:         d.IsActive,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}
