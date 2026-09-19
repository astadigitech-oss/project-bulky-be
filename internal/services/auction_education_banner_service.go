package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
)

type AuctionEducationBannerService interface {
	Create(context.Context, string, string, string, int, *time.Time, *time.Time) (*models.AuctionEducationBannerResponse, error)
	Get(context.Context, string) (*models.AuctionEducationBannerResponse, error)
	List(context.Context, string, int, int) ([]models.AuctionEducationBannerResponse, *models.PaginationMeta, error)
	Update(context.Context, string, *string, *string, *string, *int, *time.Time, *time.Time) (*models.AuctionEducationBannerResponse, error)
	Reorder(context.Context, []string) ([]models.AuctionEducationBannerResponse, error)
	Delete(context.Context, string) error
}
type auctionEducationBannerService struct {
	repo repositories.AuctionEducationBannerRepository
	cfg  *config.Config
}

func NewAuctionEducationBannerService(r repositories.AuctionEducationBannerRepository, cfg *config.Config) AuctionEducationBannerService {
	return &auctionEducationBannerService{r, cfg}
}
func (s *auctionEducationBannerService) Create(ctx context.Context, nama, id, en string, urutan int, start, end *time.Time) (*models.AuctionEducationBannerResponse, error) {
	if nama == "" || id == "" || en == "" {
		return nil, errors.New("nama dan kedua gambar wajib diisi")
	}
	if len(nama) > 100 {
		return nil, errors.New("nama maksimal 100 karakter")
	}
	if urutan < 0 {
		return nil, errors.New("urutan tidak boleh negatif")
	}
	if err := validAuctionBannerSchedule(start, end); err != nil {
		return nil, err
	}
	b := &models.AuctionEducationBanner{ID: uuid.New(), Nama: nama, GambarURLID: id, GambarURLEN: en, Urutan: urutan, TanggalMulai: start, TanggalSelesai: end}
	if err := s.repo.Create(ctx, b); err != nil {
		return nil, err
	}
	return s.response(b), nil
}
func (s *auctionEducationBannerService) Get(ctx context.Context, id string) (*models.AuctionEducationBannerResponse, error) {
	b, err := s.repo.FindByID(ctx, parseAuctionBannerID(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("banner edukasi lelang tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	return s.response(b), nil
}
func (s *auctionEducationBannerService) List(ctx context.Context, search string, page, perPage int) ([]models.AuctionEducationBannerResponse, *models.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	b, total, err := s.repo.List(ctx, search, page, perPage)
	if err != nil {
		return nil, nil, err
	}
	out := make([]models.AuctionEducationBannerResponse, 0, len(b))
	for i := range b {
		out = append(out, *s.response(&b[i]))
	}
	meta := models.NewPaginationMeta(page, perPage, total)
	return out, &meta, nil
}
func (s *auctionEducationBannerService) Update(ctx context.Context, id string, nama, gambarID, gambarEN *string, urutan *int, start, end *time.Time) (*models.AuctionEducationBannerResponse, error) {
	b, err := s.repo.FindByID(ctx, parseAuctionBannerID(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("banner edukasi lelang tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	if nama != nil {
		if *nama == "" || len(*nama) > 100 {
			return nil, errors.New("nama wajib dan maksimal 100 karakter")
		}
		b.Nama = *nama
	}
	if gambarID != nil {
		b.GambarURLID = *gambarID
	}
	if gambarEN != nil {
		b.GambarURLEN = *gambarEN
	}
	// Position changes are handled by Reorder so form edits cannot introduce
	// duplicate or colliding positions. Keep the parameter for API compatibility.
	_ = urutan
	b.TanggalMulai = start
	b.TanggalSelesai = end
	if err := validAuctionBannerSchedule(start, end); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, b); err != nil {
		return nil, err
	}
	return s.response(b), nil
}
func (s *auctionEducationBannerService) Delete(ctx context.Context, id string) error {
	b, err := s.repo.FindByID(ctx, parseAuctionBannerID(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("banner edukasi lelang tidak ditemukan")
	}
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, b)
}
func (s *auctionEducationBannerService) Reorder(ctx context.Context, ids []string) ([]models.AuctionEducationBannerResponse, error) {
	if len(ids) == 0 {
		return nil, errors.New("ids wajib diisi")
	}
	parsed := make([]uuid.UUID, 0, len(ids))
	seen := make(map[uuid.UUID]struct{}, len(ids))
	for _, value := range ids {
		parsedID, err := uuid.Parse(value)
		if err != nil {
			return nil, errors.New("id banner tidak valid")
		}
		if _, duplicate := seen[parsedID]; duplicate {
			return nil, errors.New("ids tidak boleh duplikat")
		}
		seen[parsedID] = struct{}{}
		parsed = append(parsed, parsedID)
	}
	banners, err := s.repo.Reorder(ctx, parsed)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("ids harus memuat seluruh banner aktif")
	}
	if err != nil {
		return nil, err
	}
	responses := make([]models.AuctionEducationBannerResponse, 0, len(banners))
	for index := range banners {
		responses = append(responses, *s.response(&banners[index]))
	}
	return responses, nil
}
func (s *auctionEducationBannerService) response(b *models.AuctionEducationBanner) *models.AuctionEducationBannerResponse {
	return &models.AuctionEducationBannerResponse{ID: b.ID.String(), Nama: b.Nama, GambarURL: models.TranslatableImage{ID: fullAuctionBannerURL(s.cfg.BaseURL, b.GambarURLID), EN: stringPtr(fullAuctionBannerURL(s.cfg.BaseURL, b.GambarURLEN))}, Urutan: b.Urutan, TanggalMulai: b.TanggalMulai, TanggalSelesai: b.TanggalSelesai, CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt}
}
func parseAuctionBannerID(v string) uuid.UUID {
	id, err := uuid.Parse(v)
	if err != nil {
		return uuid.Nil
	}
	return id
}
func validAuctionBannerSchedule(a, b *time.Time) error {
	if a != nil && b != nil && !b.After(*a) {
		return errors.New("tanggal selesai harus setelah tanggal mulai")
	}
	return nil
}
func fullAuctionBannerURL(base, path string) string {
	if path == "" {
		return ""
	}
	return fmt.Sprintf("%s/uploads/%s", base, path)
}
func stringPtr(v string) *string { return &v }
