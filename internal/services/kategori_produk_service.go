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
	CreateWithIcon(ctx context.Context, req *models.CreateKategoriProdukRequest, iconURL, gambarKondisiURL *string) (*models.KategoriProdukResponse, error)
	FindByID(ctx context.Context, id string) (*models.KategoriProdukResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.KategoriProdukResponse, error)
	FindAll(ctx context.Context, params *models.KategoriProdukFilterRequest) ([]models.KategoriProdukSimpleResponse, *models.PaginationMeta, error)
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
	slug := utils.GenerateSlug(req.NamaID)

	exists, _ := s.repo.ExistsBySlug(ctx, slug, nil)
	if exists {
		return nil, errors.New("kategori dengan nama tersebut sudah ada")
	}

	kategori := &models.KategoriProduk{
		NamaID:              req.NamaID,
		NamaEN:              req.NamaEN,
		Slug:                slug,
		Deskripsi:           req.Deskripsi,
		TipeKondisiTambahan: req.TipeKondisiTambahan,
		TeksKondisi:         req.TeksKondisi,
		IsActive:            true,
	}

	if err := s.repo.Create(ctx, kategori); err != nil {
		return nil, err
	}

	return s.toResponse(kategori), nil
}

func (s *kategoriProdukService) CreateWithIcon(ctx context.Context, req *models.CreateKategoriProdukRequest, iconURL, gambarKondisiURL *string) (*models.KategoriProdukResponse, error) {
	slug := utils.GenerateSlug(req.NamaID)

	exists, _ := s.repo.ExistsBySlug(ctx, slug, nil)
	if exists {
		return nil, errors.New("kategori dengan nama tersebut sudah ada")
	}

	kategori := &models.KategoriProduk{
		NamaID:              req.NamaID,
		NamaEN:              req.NamaEN,
		Slug:                slug,
		Deskripsi:           req.Deskripsi,
		IconURL:             iconURL,
		TipeKondisiTambahan: req.TipeKondisiTambahan,
		GambarKondisiURL:    gambarKondisiURL,
		TeksKondisi:         req.TeksKondisi,
		IsActive:            true,
	}

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

func (s *kategoriProdukService) FindAll(ctx context.Context, params *models.KategoriProdukFilterRequest) ([]models.KategoriProdukSimpleResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	kategoris, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	// Map entities to simple response
	items := []models.KategoriProdukSimpleResponse{}
	for _, k := range kategoris {
		items = append(items, models.KategoriProdukSimpleResponse{
			ID:      k.ID.String(),
			Nama:    k.GetNama(),
			IconURL: utils.GetFileURL(k.IconURL, s.cfg),
			// TipeKondisiTambahan: k.TipeKondisiTambahan,
			IsActive: k.IsActive,
		})
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *kategoriProdukService) Update(ctx context.Context, id string, req *models.UpdateKategoriProdukRequest) (*models.KategoriProdukResponse, error) {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kategori produk tidak ditemukan")
	}

	if req.NamaID != nil {
		newSlug := utils.GenerateSlug(*req.NamaID)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("kategori dengan nama tersebut sudah ada")
		}
		kategori.NamaID = *req.NamaID
		kategori.Slug = newSlug
	}
	if req.NamaEN != nil {
		kategori.NamaEN = req.NamaEN
	}
	if req.Deskripsi != nil {
		kategori.Deskripsi = req.Deskripsi
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
	if req.NamaID != nil {
		newSlug := utils.GenerateSlug(*req.NamaID)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("kategori dengan nama tersebut sudah ada")
		}
		kategori.NamaID = *req.NamaID
		kategori.Slug = newSlug
	}
	if req.NamaEN != nil {
		kategori.NamaEN = req.NamaEN
	}
	if req.Deskripsi != nil {
		kategori.Deskripsi = req.Deskripsi
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
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("kategori produk tidak ditemukan")
	}

	return s.repo.Delete(ctx, kategori)
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

// validateKondisiTambahan validates kondisi tambahan fields based on tipe
func (s *kategoriProdukService) validateKondisiTambahan(tipe *string, gambar *string, teks *string) error {
	if tipe == nil {
		return nil
	}

	switch *tipe {
	case models.TipeKondisiGambar:
		if gambar == nil || *gambar == "" {
			return errors.New("gambar_kondisi wajib diisi jika tipe_kondisi_tambahan = GAMBAR")
		}
	case models.TipeKondisiTeks:
		if teks == nil || *teks == "" {
			return errors.New("teks_kondisi wajib diisi jika tipe_kondisi_tambahan = TEKS")
		}
	}

	return nil
}

func (s *kategoriProdukService) toResponse(k *models.KategoriProduk) *models.KategoriProdukResponse {
	return &models.KategoriProdukResponse{
		ID:                  k.ID.String(),
		Nama:                k.GetNama(),
		Slug:                k.Slug,
		Deskripsi:           models.SafeString(k.Deskripsi),
		IconURL:             utils.GetFileURL(k.IconURL, s.cfg),
		TipeKondisiTambahan: k.TipeKondisiTambahan,
		GambarKondisiURL:    utils.GetFileURL(k.GambarKondisiURL, s.cfg),
		TeksKondisi:         models.SafeString(k.TeksKondisi),
		IsActive:            k.IsActive,
		CreatedAt:           k.CreatedAt,
		UpdatedAt:           k.UpdatedAt,
	}
}
