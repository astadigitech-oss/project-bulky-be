package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type WarehouseService interface {
	Create(ctx context.Context, req *models.CreateWarehouseRequest) (*models.WarehouseResponse, error)
	FindByID(ctx context.Context, id string) (*models.WarehouseResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest, kota string) ([]models.WarehouseResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateWarehouseRequest) (*models.WarehouseResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
	// New methods for singleton pattern
	Get(ctx context.Context) (*dto.WarehouseResponse, error)
	UpdateSingleton(ctx context.Context, req *dto.WarehouseUpdateRequest) (*dto.WarehouseResponse, error)
	GetPublic(ctx context.Context) (*dto.WarehousePublicResponse, error)
	// Informasi Pickup methods (warehouse + jadwal)
	GetInformasiPickup(ctx context.Context) (*dto.InformasiPickupResponse, error)
	GetJadwal(ctx context.Context) ([]dto.JadwalGudangResponse, error)
	UpdateJadwal(ctx context.Context, req *dto.UpdateJadwalRequest) ([]dto.JadwalGudangResponse, error)
}

type warehouseService struct {
	repo       repositories.WarehouseRepository
	jadwalRepo repositories.JadwalGudangRepository
}

func NewWarehouseService(repo repositories.WarehouseRepository, jadwalRepo repositories.JadwalGudangRepository) WarehouseService {
	return &warehouseService{
		repo:       repo,
		jadwalRepo: jadwalRepo,
	}
}

func (s *warehouseService) Create(ctx context.Context, req *models.CreateWarehouseRequest) (*models.WarehouseResponse, error) {
	slug := utils.GenerateSlug(req.Nama)

	exists, _ := s.repo.ExistsBySlug(ctx, slug, nil)
	if exists {
		return nil, errors.New("warehouse dengan nama tersebut sudah ada")
	}

	warehouse := &models.Warehouse{
		Nama:     req.Nama,
		Slug:     slug,
		Alamat:   req.Alamat,
		Kota:     req.Kota,
		KodePos:  req.KodePos,
		Telepon:  req.Telepon,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, warehouse); err != nil {
		return nil, err
	}

	return s.toResponse(warehouse), nil
}

func (s *warehouseService) FindByID(ctx context.Context, id string) (*models.WarehouseResponse, error) {
	warehouse, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}
	return s.toResponse(warehouse), nil
}

func (s *warehouseService) FindAll(ctx context.Context, params *models.PaginationRequest, kota string) ([]models.WarehouseResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	warehouses, total, err := s.repo.FindAll(ctx, params, kota)
	if err != nil {
		return nil, nil, err
	}

	items := []models.WarehouseResponse{}
	for _, w := range warehouses {
		items = append(items, *s.toResponse(&w))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *warehouseService) Update(ctx context.Context, id string, req *models.UpdateWarehouseRequest) (*models.WarehouseResponse, error) {
	warehouse, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}

	if req.Nama != nil {
		newSlug := utils.GenerateSlug(*req.Nama)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("warehouse dengan nama tersebut sudah ada")
		}
		warehouse.Nama = *req.Nama
		warehouse.Slug = newSlug
	}
	if req.Alamat != nil {
		warehouse.Alamat = req.Alamat
	}
	if req.Kota != nil {
		warehouse.Kota = req.Kota
	}
	if req.KodePos != nil {
		warehouse.KodePos = req.KodePos
	}
	if req.Telepon != nil {
		warehouse.Telepon = req.Telepon
	}
	if req.IsActive != nil {
		warehouse.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, warehouse); err != nil {
		return nil, err
	}

	return s.toResponse(warehouse), nil
}

func (s *warehouseService) Delete(ctx context.Context, id string) error {
	warehouse, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("warehouse tidak ditemukan")
	}
	// TODO: Check if warehouse has products
	return s.repo.Delete(ctx, warehouse)
}

func (s *warehouseService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	warehouse, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}

	warehouse.IsActive = !warehouse.IsActive
	if err := s.repo.Update(ctx, warehouse); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:        warehouse.ID.String(),
		IsActive:  warehouse.IsActive,
		UpdatedAt: warehouse.UpdatedAt,
	}, nil
}

func (s *warehouseService) toResponse(w *models.Warehouse) *models.WarehouseResponse {
	return &models.WarehouseResponse{
		ID:           w.ID.String(),
		Nama:         w.Nama,
		Slug:         w.Slug,
		Alamat:       w.Alamat,
		Kota:         w.Kota,
		KodePos:      w.KodePos,
		Telepon:      w.Telepon,
		IsActive:     w.IsActive,
		JumlahProduk: 0, // TODO: Count from produk table
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
	}
}

// Get returns the first active warehouse (singleton pattern)
func (s *warehouseService) Get(ctx context.Context) (*dto.WarehouseResponse, error) {
	warehouse, err := s.repo.FindFirstActive(ctx)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}
	return s.toWarehouseResponse(warehouse), nil
}

// UpdateSingleton updates the first active warehouse (singleton pattern)
func (s *warehouseService) UpdateSingleton(ctx context.Context, req *dto.WarehouseUpdateRequest) (*dto.WarehouseResponse, error) {
	warehouse, err := s.repo.FindFirstActive(ctx)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}

	// Update fields
	warehouse.Nama = req.Nama
	warehouse.Alamat = req.Alamat
	warehouse.Kota = req.Kota
	warehouse.KodePos = req.KodePos
	warehouse.Latitude = req.Latitude
	warehouse.Longitude = req.Longitude
	warehouse.JamOperasional = req.JamOperasional

	// Regenerate slug if nama changed
	newSlug := utils.GenerateSlug(req.Nama)
	if warehouse.Slug != newSlug {
		warehouse.Slug = newSlug
	}

	if err := s.repo.Update(ctx, warehouse); err != nil {
		return nil, err
	}

	return s.toWarehouseResponse(warehouse), nil
}

// GetPublic returns simplified warehouse data for public
func (s *warehouseService) GetPublic(ctx context.Context) (*dto.WarehousePublicResponse, error) {
	warehouse, err := s.repo.FindFirstActive(ctx)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}
	return &dto.WarehousePublicResponse{
		Nama:      warehouse.Nama,
		Alamat:    warehouse.Alamat,
		Kota:      warehouse.Kota,
		Latitude:  warehouse.Latitude,
		Longitude: warehouse.Longitude,
	}, nil
}

func (s *warehouseService) toWarehouseResponse(w *models.Warehouse) *dto.WarehouseResponse {
	return &dto.WarehouseResponse{
		ID:             w.ID.String(),
		Nama:           w.Nama,
		Alamat:         w.Alamat,
		Kota:           w.Kota,
		KodePos:        w.KodePos,
		Latitude:       w.Latitude,
		Longitude:      w.Longitude,
		JamOperasional: w.JamOperasional,
		CreatedAt:      w.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      w.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// GetInformasiPickup returns warehouse + jadwal for public informasi pickup endpoint
func (s *warehouseService) GetInformasiPickup(ctx context.Context) (*dto.InformasiPickupResponse, error) {
	warehouse, err := s.repo.FindFirstActive(ctx)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}

	// Get jadwal
	jadwal, err := s.jadwalRepo.FindByWarehouseID(ctx, warehouse.ID)
	if err != nil {
		return nil, err
	}

	// Convert jadwal to response
	jadwalResp := []dto.JadwalGudangResponse{}
	for _, j := range jadwal {
		jadwalResp = append(jadwalResp, dto.JadwalGudangResponse{
			Hari:     j.Hari,
			NamaHari: j.GetHariNama(),
			JamBuka:  j.JamBuka,
			JamTutup: j.JamTutup,
			IsBuka:   j.IsBuka,
		})
	}

	// Calculate is_open_now
	isOpenNow := s.isOpenNow(jadwal)
	// statusText := s.getStatusText(jadwal)
	jadwalHariIni := s.getJadwalHariIni(jadwal)

	// Generate URLs
	whatsappURL := ""
	if warehouse.Telepon != nil {
		whatsappURL = fmt.Sprintf("https://wa.me/%s", *warehouse.Telepon)
	}

	// googleMapsURL := ""
	// if warehouse.Latitude != nil && warehouse.Longitude != nil {
	// 	googleMapsURL = fmt.Sprintf("https://maps.google.com/?q=%f,%f", *warehouse.Latitude, *warehouse.Longitude)
	// }

	return &dto.InformasiPickupResponse{
		Alamat:         warehouse.Alamat,
		JamOperasional: warehouse.JamOperasional,
		Telepon:        warehouse.Telepon,
		WhatsappURL:    whatsappURL,
		Latitude:       warehouse.Latitude,
		Longitude:      warehouse.Longitude,
		// GoogleMapsURL:  googleMapsURL,
		IsOpenNow: isOpenNow,
		// StatusText:     statusText,
		JadwalHariIni: jadwalHariIni,
		Jadwal:        jadwalResp,
	}, nil
}

// GetJadwal returns jadwal gudang as array
func (s *warehouseService) GetJadwal(ctx context.Context) ([]dto.JadwalGudangResponse, error) {
	warehouse, err := s.repo.FindFirstActive(ctx)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}

	// Get jadwal
	jadwal, err := s.jadwalRepo.FindByWarehouseID(ctx, warehouse.ID)
	if err != nil {
		return nil, err
	}

	// Convert jadwal to response
	jadwalResp := []dto.JadwalGudangResponse{}
	for _, j := range jadwal {
		jadwalResp = append(jadwalResp, dto.JadwalGudangResponse{
			ID:       j.ID.String(),
			Hari:     j.Hari,
			NamaHari: j.GetHariNama(),
			JamBuka:  j.JamBuka,
			JamTutup: j.JamTutup,
			IsBuka:   j.IsBuka,
		})
	}

	return jadwalResp, nil
}

// UpdateJadwal updates jadwal gudang and returns updated jadwal array
func (s *warehouseService) UpdateJadwal(ctx context.Context, req *dto.UpdateJadwalRequest) ([]dto.JadwalGudangResponse, error) {
	warehouse, err := s.repo.FindFirstActive(ctx)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("warehouse tidak ditemukan")
	}

	// Validate jadwal
	for _, j := range req.Jadwal {
		if j.IsBuka {
			if j.JamBuka == nil || j.JamTutup == nil {
				return nil, fmt.Errorf("jam buka dan jam tutup wajib diisi untuk hari %d", j.Hari)
			}
			if *j.JamBuka >= *j.JamTutup {
				return nil, fmt.Errorf("jam tutup harus lebih besar dari jam buka untuk hari %d", j.Hari)
			}
		}
	}

	// Convert to models
	jadwalModels := []models.JadwalGudang{}
	for _, j := range req.Jadwal {
		jadwalModels = append(jadwalModels, models.JadwalGudang{
			WarehouseID: warehouse.ID,
			Hari:        j.Hari,
			JamBuka:     j.JamBuka,
			JamTutup:    j.JamTutup,
			IsBuka:      j.IsBuka,
		})
	}

	if err := s.jadwalRepo.UpdateBatch(ctx, warehouse.ID, jadwalModels); err != nil {
		return nil, err
	}

	// Return updated jadwal
	return s.GetJadwal(ctx)
}

// Helper functions
func (s *warehouseService) isOpenNow(jadwal []models.JadwalGudang) bool {
	now := time.Now()
	currentDay := int(now.Weekday())
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

func (s *warehouseService) getStatusText(jadwal []models.JadwalGudang) string {
	now := time.Now()
	currentDay := int(now.Weekday())
	currentTime := now.Format("15:04")

	for _, j := range jadwal {
		if j.Hari == currentDay {
			if !j.IsBuka {
				return "Tutup"
			}
			if j.JamBuka != nil && j.JamTutup != nil {
				if currentTime >= *j.JamBuka && currentTime <= *j.JamTutup {
					return "Buka"
				}
				if currentTime < *j.JamBuka {
					return fmt.Sprintf("Buka pukul %s", *j.JamBuka)
				}
			}
			return "Tutup"
		}
	}
	return "Tutup"
}

func (s *warehouseService) getJadwalHariIni(jadwal []models.JadwalGudang) *dto.JadwalGudangResponse {
	now := time.Now()
	currentDay := int(now.Weekday())

	for _, j := range jadwal {
		if j.Hari == currentDay {
			return &dto.JadwalGudangResponse{
				Hari:     j.Hari,
				NamaHari: j.GetHariNama(),
				JamBuka:  j.JamBuka,
				JamTutup: j.JamTutup,
				IsBuka:   j.IsBuka,
			}
		}
	}
	return nil
}
