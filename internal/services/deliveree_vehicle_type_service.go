package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

// DelivereeVehicleTypeService mengelola master data kendaraan Deliveree (Sync dari
// API vehicle_types) dan menyediakan pemilihan kendaraan berdasarkan kubikasi & berat
// untuk menggantikan logic lama yang berbasis jumlah qty/palet.
type DelivereeVehicleTypeService interface {
	FindByID(ctx context.Context, id string) (*models.DelivereeVehicleTypeResponse, error)
	FindAll(ctx context.Context, params *models.DelivereeVehicleTypeFilterRequest) ([]models.DelivereeVehicleTypeResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateDelivereeVehicleTypeRequest) (*models.DelivereeVehicleTypeResponse, error)
	// Sync menarik data terbaru dari GET /public_api/v10/vehicle_types memakai
	// DELIVEREE_BASE_URL/DELIVEREE_API_KEY yang sedang aktif (mengikuti environment
	// deployment saat ini — sandbox atau production, tanpa perlu credential ganda),
	// upsert ke DB, lalu menghitung ulang threshold_kubikasi & threshold_berat
	// berdasarkan urutan kubikasi_max.
	Sync(ctx context.Context) (*models.SyncDelivereeVehicleTypeResponse, error)
	// SelectVehicle memilih kendaraan aktif TERKECIL pada environment tertentu yang
	// kapasitasnya (kubikasi_max & berat_max) mencukupi totalKubikasi & totalBerat.
	// Mengembalikan error jika tidak ada kendaraan yang cukup besar.
	SelectVehicle(ctx context.Context, environment string, totalKubikasi, totalBerat float64) (*models.DelivereeVehicleType, error)
	// FindActiveByIDDeliveree mencari satu kendaraan AKTIF berdasarkan id_deliveree+environment.
	// Dipakai saat create booking untuk memakai deliveree_vehicle_type_id yang disimpan
	// storefront saat checkout; return nil jika tidak ada (kendaraan dinonaktifkan).
	FindActiveByIDDeliveree(ctx context.Context, idDeliveree int, environment string) (*models.DelivereeVehicleType, error)
	// NextLargerVehicle mengembalikan kendaraan aktif dengan kubikasi_max satu tingkat
	// di atas currentIDDeliveree pada environment yang sama. Dipakai untuk retry otomatis
	// saat booking ke kendaraan yang dipilih gagal (mis. tidak ada driver tersedia).
	NextLargerVehicle(ctx context.Context, environment string, currentIDDeliveree int) (*models.DelivereeVehicleType, error)
}

type delivereeVehicleTypeService struct {
	repo          repositories.DelivereeVehicleTypeRepository
	warehouseRepo repositories.WarehouseRepository
	activityLog   ActivityLogService
}

func NewDelivereeVehicleTypeService(repo repositories.DelivereeVehicleTypeRepository, warehouseRepo repositories.WarehouseRepository, activityLog ActivityLogService) DelivereeVehicleTypeService {
	return &delivereeVehicleTypeService{repo: repo, warehouseRepo: warehouseRepo, activityLog: activityLog}
}

func toDelivereeVehicleTypeResponse(v *models.DelivereeVehicleType) models.DelivereeVehicleTypeResponse {
	return models.DelivereeVehicleTypeResponse{
		ID:                v.ID.String(),
		Nama:              v.Nama,
		IDDeliveree:       v.IDDeliveree,
		Environment:       v.Environment,
		KubikasiMax:       v.KubikasiMax,
		BeratMax:          v.BeratMax,
		ThresholdKubikasi: v.ThresholdKubikasi,
		ThresholdBerat:    v.ThresholdBerat,
		CargoLength:       v.CargoLength,
		CargoWidth:        v.CargoWidth,
		CargoHeight:       v.CargoHeight,
		IsActive:          v.IsActive,
		LastSyncedAt:      v.LastSyncedAt,
		CreatedAt:         v.CreatedAt,
		UpdatedAt:         v.UpdatedAt,
	}
}

func (s *delivereeVehicleTypeService) FindByID(ctx context.Context, id string) (*models.DelivereeVehicleTypeResponse, error) {
	vehicle, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toDelivereeVehicleTypeResponse(vehicle)
	return &resp, nil
}

func (s *delivereeVehicleTypeService) FindAll(ctx context.Context, params *models.DelivereeVehicleTypeFilterRequest) ([]models.DelivereeVehicleTypeResponse, *models.PaginationMeta, error) {
	vehicles, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	result := make([]models.DelivereeVehicleTypeResponse, 0, len(vehicles))
	for i := range vehicles {
		result = append(result, toDelivereeVehicleTypeResponse(&vehicles[i]))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return result, &meta, nil
}

func (s *delivereeVehicleTypeService) Update(ctx context.Context, id string, req *models.UpdateDelivereeVehicleTypeRequest) (*models.DelivereeVehicleTypeResponse, error) {
	vehicle, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ThresholdKubikasi != nil {
		vehicle.ThresholdKubikasi = *req.ThresholdKubikasi
	}
	if req.ThresholdBerat != nil {
		vehicle.ThresholdBerat = *req.ThresholdBerat
	}
	if req.IsActive != nil {
		vehicle.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, vehicle); err != nil {
		return nil, err
	}

	resp := toDelivereeVehicleTypeResponse(vehicle)
	return &resp, nil
}

func (s *delivereeVehicleTypeService) Sync(ctx context.Context) (*models.SyncDelivereeVehicleTypeResponse, error) {
	baseURL := os.Getenv("DELIVEREE_BASE_URL")
	apiKey := os.Getenv("DELIVEREE_API_KEY")
	if baseURL == "" || apiKey == "" {
		return nil, fmt.Errorf("konfigurasi Deliveree tidak lengkap")
	}

	// Environment dideteksi otomatis dari DELIVEREE_BASE_URL yang sedang aktif —
	// sama seperti logic di bookDeliveree(), sehingga tidak perlu credential ganda.
	// Untuk sync data dari environment lain, cukup ganti DELIVEREE_BASE_URL/API_KEY
	// di .env sesuai environment yang dituju (sandbox/production itu sendiri sudah
	// merupakan deployment terpisah).
	environment := string(models.DelivereeEnvProduction)
	if strings.Contains(baseURL, "sandbox") {
		environment = string(models.DelivereeEnvSandbox)
	}

	// pickup_location wajib diisi oleh Deliveree API — pakai lokasi warehouse aktif
	// sebagai referensi (sekadar untuk menentukan area layanan, tidak memengaruhi
	// kapasitas kendaraan yang dikembalikan).
	warehouse, err := s.warehouseRepo.FindFirstActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan data warehouse untuk Sync: %w", err)
	}
	lat, lng := 0.0, 0.0
	if warehouse.Latitude != nil {
		lat = *warehouse.Latitude
	}
	if warehouse.Longitude != nil {
		lng = *warehouse.Longitude
	}

	url := fmt.Sprintf("%s/vehicle_types?pickup_location[latitude]=%f&pickup_location[longitude]=%f", baseURL, lat, lng)
	log.Printf("[deliveree-vehicle-sync] --> GET /vehicle_types environment=%s", environment)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi Deliveree (%s): %w", environment, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[deliveree-vehicle-sync] <-- GET /vehicle_types environment=%s status=%d body=%s", environment, resp.StatusCode, string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Deliveree API error saat Sync (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Response Deliveree membungkus array vehicle di dalam key "data",
	// mis. {"data":[{...},{...}]} — jadi wajib parse ke wrapper dulu.
	var apiResp struct {
		Data []DelivereeVehicleTypeInfo `json:"data"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("gagal parse response Deliveree vehicle_types: %w", err)
	}
	apiVehicles := apiResp.Data
	if apiVehicles == nil {
		apiVehicles = []DelivereeVehicleTypeInfo{}
	}

	now := time.Now().UTC()
	activeDelivereeIDs := make([]int, 0, len(apiVehicles))
	created, updated := 0, 0

	for _, av := range apiVehicles {
		activeDelivereeIDs = append(activeDelivereeIDs, av.ID)

		existing, findErr := s.findExisting(ctx, av.ID, environment)
		cargoLength := av.CargoLength
		cargoWidth := av.CargoWidth
		cargoHeight := av.CargoHeight

		vehicle := models.DelivereeVehicleType{
			Nama:         av.Name,
			IDDeliveree:  av.ID,
			Environment:  environment,
			KubikasiMax:  av.CargoCubicMeter,
			BeratMax:     av.CargoWeight,
			CargoLength:  &cargoLength,
			CargoWidth:   &cargoWidth,
			CargoHeight:  &cargoHeight,
			IsActive:     true,
			LastSyncedAt: &now,
		}
		if findErr == nil && existing != nil {
			vehicle.ID = existing.ID
			// threshold tidak disentuh di sini, direcompute setelah semua vehicle di-upsert
			vehicle.ThresholdKubikasi = existing.ThresholdKubikasi
			vehicle.ThresholdBerat = existing.ThresholdBerat
			updated++
		} else {
			vehicle.ID = uuid.New()
			created++
		}

		if err := s.repo.Upsert(ctx, &vehicle); err != nil {
			return nil, fmt.Errorf("gagal menyimpan kendaraan %q (id_deliveree=%d): %w", av.Name, av.ID, err)
		}
	}

	deactivated, err := s.repo.DeactivateMissing(ctx, environment, activeDelivereeIDs)
	if err != nil {
		return nil, fmt.Errorf("gagal menonaktifkan kendaraan yang sudah tidak ada di API: %w", err)
	}

	if err := s.recomputeThresholds(ctx, environment); err != nil {
		return nil, fmt.Errorf("gagal menghitung ulang threshold: %w", err)
	}

	return &models.SyncDelivereeVehicleTypeResponse{
		Environment:  environment,
		TotalFromAPI: len(apiVehicles),
		Created:      created,
		Updated:      updated,
		Deactivated:  int(deactivated),
		SyncedAt:     now,
	}, nil
}

// recomputeThresholds mengurutkan kendaraan aktif berdasarkan kubikasi_max ascending,
// lalu men-set threshold masing-masing = kapasitas kendaraan satu tingkat di bawahnya
// (0 untuk kendaraan terkecil). Ini memberi buffer/slack sebelum sistem "naik kelas"
// ke kendaraan yang lebih besar.
func (s *delivereeVehicleTypeService) recomputeThresholds(ctx context.Context, environment string) error {
	vehicles, err := s.repo.FindActiveByEnvironment(ctx, environment)
	if err != nil {
		return err
	}

	sort.Slice(vehicles, func(i, j int) bool {
		return vehicles[i].KubikasiMax < vehicles[j].KubikasiMax
	})

	prevKubikasi, prevBerat := 0.0, 0.0
	for i := range vehicles {
		if err := s.repo.UpdateThresholds(ctx, vehicles[i].ID, prevKubikasi, prevBerat); err != nil {
			return err
		}
		prevKubikasi = vehicles[i].KubikasiMax
		prevBerat = vehicles[i].BeratMax
	}
	return nil
}

func (s *delivereeVehicleTypeService) SelectVehicle(ctx context.Context, environment string, totalKubikasi, totalBerat float64) (*models.DelivereeVehicleType, error) {
	vehicles, err := s.repo.FindActiveByEnvironment(ctx, environment)
	if err != nil {
		return nil, err
	}

	sort.Slice(vehicles, func(i, j int) bool {
		return vehicles[i].KubikasiMax < vehicles[j].KubikasiMax
	})

	for i := range vehicles {
		if vehicles[i].KubikasiMax >= totalKubikasi && vehicles[i].BeratMax >= totalBerat {
			return &vehicles[i], nil
		}
	}

	return nil, fmt.Errorf("tidak ada kendaraan Deliveree (%s) yang cukup untuk kubikasi %.3f m3 / berat %.2f kg", environment, totalKubikasi, totalBerat)
}

func (s *delivereeVehicleTypeService) FindActiveByIDDeliveree(ctx context.Context, idDeliveree int, environment string) (*models.DelivereeVehicleType, error) {
	return s.repo.FindActiveByIDDeliveree(ctx, idDeliveree, environment)
}

func (s *delivereeVehicleTypeService) NextLargerVehicle(ctx context.Context, environment string, currentIDDeliveree int) (*models.DelivereeVehicleType, error) {
	vehicles, err := s.repo.FindActiveByEnvironment(ctx, environment)
	if err != nil {
		return nil, err
	}

	sort.Slice(vehicles, func(i, j int) bool {
		return vehicles[i].KubikasiMax < vehicles[j].KubikasiMax
	})

	foundCurrent := false
	for i := range vehicles {
		if foundCurrent {
			return &vehicles[i], nil
		}
		if vehicles[i].IDDeliveree == currentIDDeliveree {
			foundCurrent = true
		}
	}

	return nil, fmt.Errorf("tidak ada kendaraan Deliveree (%s) yang lebih besar dari id_deliveree=%d", environment, currentIDDeliveree)
}

// findExisting mencari kendaraan existing berdasarkan id_deliveree+environment.
// Termasuk yang is_active=false agar saat Sync kendaraan yang sebelumnya
// dinonaktifkan tetap dihitung "updated" (bukan "created") dan ID-nya dipertahankan.
func (s *delivereeVehicleTypeService) findExisting(ctx context.Context, idDeliveree int, environment string) (*models.DelivereeVehicleType, error) {
	return s.repo.FindByIDDeliveree(ctx, idDeliveree, environment)
}
