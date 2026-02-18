package services

import (
	"context"
	"errors"
	"time"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KuponService interface {
	Create(ctx context.Context, req *dto.CreateKuponRequest) (*dto.KuponDetailResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateKuponRequest) (*dto.KuponDetailResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.KuponDetailResponse, error)
	GetAll(ctx context.Context, params *dto.KuponQueryParams) ([]dto.KuponListResponse, *models.PaginationMeta, error)
	ToggleStatus(ctx context.Context, id uuid.UUID) (*dto.ToggleStatusResponse, error)
	GenerateKode(ctx context.Context, req *dto.GenerateKodeRequest) (*dto.GeneratedKodeResponse, error)
	GetUsages(ctx context.Context, kuponID uuid.UUID, page, limit int) (*dto.KuponUsageListResponse, *models.PaginationMeta, error)
	GetKategoriDropdown(ctx context.Context) ([]dto.KategoriDropdownResponse, error)
}

type kuponService struct {
	kuponRepo    repositories.KuponRepository
	kategoriRepo repositories.KategoriProdukRepository
	db           *gorm.DB
}

func NewKuponService(
	kuponRepo repositories.KuponRepository,
	kategoriRepo repositories.KategoriProdukRepository,
	db *gorm.DB,
) KuponService {
	return &kuponService{
		kuponRepo:    kuponRepo,
		kategoriRepo: kategoriRepo,
		db:           db,
	}
}

func (s *kuponService) Create(ctx context.Context, req *dto.CreateKuponRequest) (*dto.KuponDetailResponse, error) {
	// Custom validation
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if kode already exists
	exists, err := s.kuponRepo.IsKodeExists(ctx, req.Kode, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("kode kupon sudah digunakan")
	}

	// Validate kategori if not all kategori
	if req.IsAllKategori != nil && !*req.IsAllKategori && len(req.KategoriIDs) > 0 {
		for _, kategoriID := range req.KategoriIDs {
			_, err := s.kategoriRepo.FindByID(ctx, kategoriID.String())
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("kategori tidak ditemukan")
				}
				return nil, err
			}
		}
	}

	// Create kupon
	isActiveTrue := true
	kupon := &models.Kupon{
		Kode:              req.Kode,
		Nama:              req.Nama,
		Deskripsi:         req.Deskripsi,
		JenisDiskon:       models.JenisDiskon(req.JenisDiskon),
		NilaiDiskon:       req.NilaiDiskon,
		MinimalPembelian:  req.MinimalPembelian,
		LimitPemakaian:    req.LimitPemakaian,
		TanggalKedaluarsa: req.TanggalKedaluarsa.UTC(),
		IsAllKategori:     req.IsAllKategori,
		IsActive:          &isActiveTrue,
	}

	if err := s.kuponRepo.Create(ctx, kupon); err != nil {
		return nil, err
	}

	// Add kategori if not all kategori
	if req.IsAllKategori != nil && !*req.IsAllKategori && len(req.KategoriIDs) > 0 {
		if err := s.kuponRepo.AddKategori(ctx, kupon.ID, req.KategoriIDs); err != nil {
			return nil, err
		}
	}

	// Get created kupon with relations
	created, err := s.kuponRepo.FindByID(ctx, kupon.ID)
	if err != nil {
		return nil, err
	}

	return s.toDetailResponse(created), nil
}

func (s *kuponService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateKuponRequest) (*dto.KuponDetailResponse, error) {
	// Custom validation
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Find existing kupon
	kupon, err := s.kuponRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kupon tidak ditemukan")
		}
		return nil, err
	}

	// Check if kode already exists (exclude current)
	exists, err := s.kuponRepo.IsKodeExists(ctx, req.Kode, &id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("kode kupon sudah digunakan")
	}

	// Validate kategori if not all kategori
	if req.IsAllKategori != nil && !*req.IsAllKategori && len(req.KategoriIDs) > 0 {
		for _, kategoriID := range req.KategoriIDs {
			_, err := s.kategoriRepo.FindByID(ctx, kategoriID.String())
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("kategori tidak ditemukan")
				}
				return nil, err
			}
		}
	}

	// Update kupon fields
	kupon.Kode = req.Kode
	kupon.Nama = req.Nama
	kupon.Deskripsi = req.Deskripsi
	kupon.JenisDiskon = models.JenisDiskon(req.JenisDiskon)
	kupon.NilaiDiskon = req.NilaiDiskon
	kupon.MinimalPembelian = req.MinimalPembelian
	kupon.LimitPemakaian = req.LimitPemakaian
	kupon.TanggalKedaluarsa = req.TanggalKedaluarsa.UTC()
	kupon.IsAllKategori = req.IsAllKategori

	if err := s.kuponRepo.Update(ctx, kupon); err != nil {
		return nil, err
	}

	// Update kategori relations
	if err := s.kuponRepo.RemoveAllKategori(ctx, id); err != nil {
		return nil, err
	}

	if req.IsAllKategori != nil && !*req.IsAllKategori && len(req.KategoriIDs) > 0 {
		if err := s.kuponRepo.AddKategori(ctx, id, req.KategoriIDs); err != nil {
			return nil, err
		}
	}

	// Get updated kupon with relations
	updated, err := s.kuponRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toDetailResponse(updated), nil
}

func (s *kuponService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if kupon exists
	_, err := s.kuponRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("kupon tidak ditemukan")
		}
		return err
	}

	return s.kuponRepo.Delete(ctx, id)
}

func (s *kuponService) GetByID(ctx context.Context, id uuid.UUID) (*dto.KuponDetailResponse, error) {
	kupon, err := s.kuponRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kupon tidak ditemukan")
		}
		return nil, err
	}

	return s.toDetailResponse(kupon), nil
}

func (s *kuponService) GetAll(ctx context.Context, params *dto.KuponQueryParams) ([]dto.KuponListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	offset := (params.Page - 1) * params.PerPage

	kupons, total, err := s.kuponRepo.FindAll(
		ctx,
		params.JenisDiskon,
		params.IsActive,
		params.IsExpired,
		params.Search,
		params.SortBy,
		params.Order,
		params.PerPage,
		offset,
	)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]dto.KuponListResponse, len(kupons))
	for i, kupon := range kupons {
		responses[i] = s.toListResponse(&kupon)
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return responses, &meta, nil
}

func (s *kuponService) ToggleStatus(ctx context.Context, id uuid.UUID) (*dto.ToggleStatusResponse, error) {
	// Check if kupon exists
	_, err := s.kuponRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kupon tidak ditemukan")
		}
		return nil, err
	}

	if err := s.kuponRepo.ToggleStatus(ctx, id); err != nil {
		return nil, err
	}

	// Get updated kupon
	updated, err := s.kuponRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.ToggleStatusResponse{
		ID:        updated.ID,
		Kode:      updated.Kode,
		IsActive:  updated.IsActive != nil && *updated.IsActive,
		UpdatedAt: updated.UpdatedAt.UTC(),
	}, nil
}

func (s *kuponService) GenerateKode(ctx context.Context, req *dto.GenerateKodeRequest) (*dto.GeneratedKodeResponse, error) {
	req.SetDefaults()

	kode := generateRandomKode(req.Prefix, req.Length)

	return &dto.GeneratedKodeResponse{
		Kode: kode,
	}, nil
}

func (s *kuponService) GetUsages(ctx context.Context, kuponID uuid.UUID, page, limit int) (*dto.KuponUsageListResponse, *models.PaginationMeta, error) {
	// Check if kupon exists
	kupon, err := s.kuponRepo.FindByID(ctx, kuponID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("kupon tidak ditemukan")
		}
		return nil, nil, err
	}

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	usages, total, err := s.kuponRepo.FindUsagesByKuponID(ctx, kuponID, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	usageResponses := make([]dto.KuponUsageItemResponse, len(usages))
	for i, usage := range usages {
		usageResponses[i] = dto.KuponUsageItemResponse{
			ID: usage.ID,
			Buyer: dto.KuponUsageBuyerInfo{
				ID:    usage.Buyer.ID,
				Nama:  usage.Buyer.Nama,
				Email: usage.Buyer.Email,
			},
			Pesanan: dto.KuponUsagePesananInfo{
				ID:   usage.Pesanan.ID,
				Kode: usage.Pesanan.Kode,
			},
			NilaiPotongan: usage.NilaiPotongan,
			CreatedAt:     usage.CreatedAt.UTC(),
		}
	}

	response := &dto.KuponUsageListResponse{
		Kupon: dto.KuponUsageSummary{
			ID:         kupon.ID,
			Kode:       kupon.Kode,
			TotalUsage: int(total),
		},
		Usages: usageResponses,
	}

	meta := models.NewPaginationMeta(page, limit, total)

	return response, &meta, nil
}

func (s *kuponService) GetKategoriDropdown(ctx context.Context) ([]dto.KategoriDropdownResponse, error) {
	kategoris, err := s.kategoriRepo.FindAllActiveForDropdown(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.KategoriDropdownResponse, len(kategoris))
	for i, kategori := range kategoris {
		responses[i] = dto.KategoriDropdownResponse{
			ID:   kategori.ID,
			Nama: kategori.GetNama(),
		}
	}

	return responses, nil
}

// Helper methods

func (s *kuponService) toListResponse(kupon *models.Kupon) dto.KuponListResponse {
	return dto.KuponListResponse{
		ID:          kupon.ID,
		Status:      kupon.IsActive != nil && *kupon.IsActive,
		JenisDiskon: string(kupon.JenisDiskon),
		Kode:        kupon.Kode,
		Nama:        kupon.Nama,
		NilaiDiskon: kupon.NilaiDiskon,
		UpdatedAt:   kupon.UpdatedAt.UTC(),
	}
}

func (s *kuponService) toDetailResponse(kupon *models.Kupon) *dto.KuponDetailResponse {
	totalUsage := int(kupon.GetUsageCount(s.db))
	remainingUsage := kupon.GetRemainingUsage(s.db)

	kategoriResponses := make([]dto.KuponKategoriResponse, 0)
	if kupon.IsAllKategori != nil && !*kupon.IsAllKategori && len(kupon.Kategori) > 0 {
		for _, kk := range kupon.Kategori {
			if kk.Kategori != nil {
				kategoriResponses = append(kategoriResponses, dto.KuponKategoriResponse{
					ID:   kk.Kategori.ID,
					Nama: kk.Kategori.GetNama(),
					Slug: kk.Kategori.Slug,
				})
			}
		}
	}

	return &dto.KuponDetailResponse{
		ID:                kupon.ID,
		Kode:              kupon.Kode,
		Nama:              kupon.Nama,
		Deskripsi:         kupon.Deskripsi,
		JenisDiskon:       string(kupon.JenisDiskon),
		NilaiDiskon:       kupon.NilaiDiskon,
		MinimalPembelian:  kupon.MinimalPembelian,
		LimitPemakaian:    kupon.LimitPemakaian,
		TotalUsage:        totalUsage,
		RemainingUsage:    remainingUsage,
		TanggalKedaluarsa: kupon.TanggalKedaluarsa.UTC(),
		IsAllKategori:     kupon.IsAllKategori != nil && *kupon.IsAllKategori,
		Kategori:          kategoriResponses,
		IsActive:          kupon.IsActive != nil && *kupon.IsActive,
		IsExpired:         kupon.IsExpired(),
		CreatedAt:         kupon.CreatedAt.UTC(),
		UpdatedAt:         kupon.UpdatedAt.UTC(),
	}
}

// generateRandomKode generates a random kupon code
func generateRandomKode(prefix string, length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // Ensure different values
	}

	if prefix != "" {
		return prefix + string(b)
	}
	return string(b)
}
