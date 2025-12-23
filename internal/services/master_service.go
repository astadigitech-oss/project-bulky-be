package services

import (
	"context"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
)

type MasterService interface {
	GetDropdown(ctx context.Context) (*models.MasterDropdownResponse, error)
}

type masterService struct {
	kategoriRepo     repositories.KategoriProdukRepository
	merekRepo        repositories.MerekProdukRepository
	kondisiRepo      repositories.KondisiProdukRepository
	kondisiPaketRepo repositories.KondisiPaketRepository
	sumberRepo       repositories.SumberProdukRepository
}

func NewMasterService(
	kategoriRepo repositories.KategoriProdukRepository,
	merekRepo repositories.MerekProdukRepository,
	kondisiRepo repositories.KondisiProdukRepository,
	kondisiPaketRepo repositories.KondisiPaketRepository,
	sumberRepo repositories.SumberProdukRepository,
) MasterService {
	return &masterService{
		kategoriRepo:     kategoriRepo,
		merekRepo:        merekRepo,
		kondisiRepo:      kondisiRepo,
		kondisiPaketRepo: kondisiPaketRepo,
		sumberRepo:       sumberRepo,
	}
}

func (s *masterService) GetDropdown(ctx context.Context) (*models.MasterDropdownResponse, error) {
	response := &models.MasterDropdownResponse{
		KategoriProduk: []models.DropdownItem{},
		MerekProduk:    []models.DropdownItem{},
		KondisiProduk:  []models.DropdownItem{},
		KondisiPaket:   []models.DropdownItem{},
		SumberProduk:   []models.DropdownItem{},
	}

	// Kategori Produk
	kategoris, err := s.kategoriRepo.GetAllForDropdown(ctx)
	if err == nil {
		for _, k := range kategoris {
			response.KategoriProduk = append(response.KategoriProduk, models.DropdownItem{
				ID:   k.ID.String(),
				Nama: k.Nama,
				Slug: k.Slug,
			})
		}
	}

	// Merek Produk
	mereks, err := s.merekRepo.GetAllForDropdown(ctx)
	if err == nil {
		for _, m := range mereks {
			response.MerekProduk = append(response.MerekProduk, models.DropdownItem{
				ID:   m.ID.String(),
				Nama: m.Nama,
				Slug: m.Slug,
			})
		}
	}

	// Kondisi Produk
	kondisis, err := s.kondisiRepo.GetAllForDropdown(ctx)
	if err == nil {
		for _, k := range kondisis {
			response.KondisiProduk = append(response.KondisiProduk, models.DropdownItem{
				ID:   k.ID.String(),
				Nama: k.Nama,
				Slug: k.Slug,
			})
		}
	}

	// Kondisi Paket
	pakets, err := s.kondisiPaketRepo.GetAllForDropdown(ctx)
	if err == nil {
		for _, p := range pakets {
			response.KondisiPaket = append(response.KondisiPaket, models.DropdownItem{
				ID:   p.ID.String(),
				Nama: p.Nama,
				Slug: p.Slug,
			})
		}
	}

	// Sumber Produk
	sumbers, err := s.sumberRepo.GetAllForDropdown(ctx)
	if err == nil {
		for _, sb := range sumbers {
			response.SumberProduk = append(response.SumberProduk, models.DropdownItem{
				ID:   sb.ID.String(),
				Nama: sb.Nama,
				Slug: sb.Slug,
			})
		}
	}

	return response, nil
}
