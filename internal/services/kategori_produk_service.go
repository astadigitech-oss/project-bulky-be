package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type KategoriProdukService interface {
	Create(ctx context.Context, req *models.CreateKategoriProdukRequest) (*models.KategoriProdukResponse, error)
	FindByID(ctx context.Context, id string) (*models.KategoriProdukResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.KategoriProdukResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.KategoriProdukResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateKategoriProdukRequest) (*models.KategoriProdukResponse, error)
	UpdateWithIcon(ctx context.Context, id string, req *models.UpdateKategoriProdukRequest, iconURL, gambarKondisiURL *string) (*models.KategoriProdukResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
}

type kategoriProdukService struct {
	repo repositories.KategoriProdukRepository
	cfg  *config.Config
}

func NewKategoriProdukService(repo repositories.KategoriProdukRepository, cfg *config.Config) KategoriProdukService {
	return &kategoriProdukService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *kategoriProdukService) Create(ctx context.Context, req *models.CreateKategoriProdukRequest) (*models.KategoriProdukResponse, error) {
	slug := utils.GenerateSlug(req.Nama)

	exists, _ := s.repo.ExistsBySlug(ctx, slug, nil)
	if exists {
		return nil, errors.New("kategori dengan nama tersebut sudah ada")
	}

	kategori := &models.KategoriProduk{
		Nama:                    req.Nama,
		Slug:                    slug,
		Deskripsi:               req.Deskripsi,
		MemilikiKondisiTambahan: req.MemilikiKondisiTambahan,
		TipeKondisiTambahan:     req.TipeKondisiTambahan,
		TeksKondisi:             req.TeksKondisi,
		IsActive:                true,
	}

	// TODO: Handle icon & gambar_kondisi upload

	if err := s.repo.Create(ctx, kategori); err != nil {
		return nil, err
	}

	return s.toResponse(kategori), nil
}

func (s *kategoriProdukService) FindByID(ctx context.Context, id string) (*models.KategoriProdukResponse, error) {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kategori produk tidak ditemukan")
	}
	return s.toResponse(kategori), nil
}

func (s *kategoriProdukService) FindBySlug(ctx context.Context, slug string) (*models.KategoriProdukResponse, error) {
	kategori, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("kategori produk tidak ditemukan")
	}
	return s.toResponse(kategori), nil
}

func (s *kategoriProdukService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.KategoriProdukResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	kategoris, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.KategoriProdukResponse
	for _, k := range kategoris {
		items = append(items, *s.toResponse(&k))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *kategoriProdukService) Update(ctx context.Context, id string, req *models.UpdateKategoriProdukRequest) (*models.KategoriProdukResponse, error) {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kategori produk tidak ditemukan")
	}

	if req.Nama != nil {
		newSlug := utils.GenerateSlug(*req.Nama)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("kategori dengan nama tersebut sudah ada")
		}
		kategori.Nama = *req.Nama
		kategori.Slug = newSlug
	}
	if req.Deskripsi != nil {
		kategori.Deskripsi = req.Deskripsi
	}
	if req.MemilikiKondisiTambahan != nil {
		kategori.MemilikiKondisiTambahan = *req.MemilikiKondisiTambahan
	}
	if req.TipeKondisiTambahan != nil {
		kategori.TipeKondisiTambahan = req.TipeKondisiTambahan
	}
	if req.TeksKondisi != nil {
		kategori.TeksKondisi = req.TeksKondisi
	}
	if req.IsActive != nil {
		kategori.IsActive = *req.IsActive
	}

	// TODO: Handle icon & gambar_kondisi upload

	if err := s.repo.Update(ctx, kategori); err != nil {
		return nil, err
	}

	return s.toResponse(kategori), nil
}

func (s *kategoriProdukService) UpdateWithIcon(ctx context.Context, id string, req *models.UpdateKategoriProdukRequest, iconURL, gambarKondisiURL *string) (*models.KategoriProdukResponse, error) {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kategori produk tidak ditemukan")
	}

	// Update text fields
	if req.Nama != nil {
		newSlug := utils.GenerateSlug(*req.Nama)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("kategori dengan nama tersebut sudah ada")
		}
		kategori.Nama = *req.Nama
		kategori.Slug = newSlug
	}
	if req.Deskripsi != nil {
		kategori.Deskripsi = req.Deskripsi
	}
	if req.MemilikiKondisiTambahan != nil {
		kategori.MemilikiKondisiTambahan = *req.MemilikiKondisiTambahan
	}
	if req.TipeKondisiTambahan != nil {
		kategori.TipeKondisiTambahan = req.TipeKondisiTambahan
	}
	if req.TeksKondisi != nil {
		kategori.TeksKondisi = req.TeksKondisi
	}
	if req.IsActive != nil {
		kategori.IsActive = *req.IsActive
	}

	// Update icon if uploaded
	if iconURL != nil {
		// Delete old icon if exists
		if kategori.IconURL != nil && *kategori.IconURL != "" {
			utils.DeleteFile(*kategori.IconURL, s.cfg)
		}
		kategori.IconURL = iconURL
	}

	// Update gambar kondisi if uploaded
	if gambarKondisiURL != nil {
		// Delete old gambar kondisi if exists
		if kategori.GambarKondisiURL != nil && *kategori.GambarKondisiURL != "" {
			utils.DeleteFile(*kategori.GambarKondisiURL, s.cfg)
		}
		kategori.GambarKondisiURL = gambarKondisiURL
	}

	if err := s.repo.Update(ctx, kategori); err != nil {
		return nil, err
	}

	return s.toResponse(kategori), nil
}

func (s *kategoriProdukService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("kategori produk tidak ditemukan")
	}

	// TODO: Check if kategori has products

	return s.repo.Delete(ctx, id)
}

func (s *kategoriProdukService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kategori produk tidak ditemukan")
	}

	kategori.IsActive = !kategori.IsActive
	if err := s.repo.Update(ctx, kategori); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       kategori.ID.String(),
		IsActive: kategori.IsActive,
	}, nil
}

func (s *kategoriProdukService) toResponse(k *models.KategoriProduk) *models.KategoriProdukResponse {
	return &models.KategoriProdukResponse{
		ID:                      k.ID.String(),
		Nama:                    k.Nama,
		Slug:                    k.Slug,
		Deskripsi:               k.Deskripsi,
		IconURL:                 k.IconURL,
		MemilikiKondisiTambahan: k.MemilikiKondisiTambahan,
		TipeKondisiTambahan:     k.TipeKondisiTambahan,
		GambarKondisiURL:        k.GambarKondisiURL,
		TeksKondisi:             k.TeksKondisi,
		IsActive:                k.IsActive,
		JumlahProduk:            0, // TODO: Count from produk table
		CreatedAt:               k.CreatedAt,
		UpdatedAt:               k.UpdatedAt,
	}
}
