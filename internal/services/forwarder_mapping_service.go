package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

// ForwarderMappingService mengelola master data mapping kota & kecamatan Forwarder
// (tabel forwarder_city_mapping & forwarder_subdistrict_mapping). Data ditarik dari
// API Forwarder (endpoint /citylist dan /subdistrictlist) lewat tombol Sync di
// admin panel — menggantikan skrip seed one-time cmd/seed-forwarder-cities.
type ForwarderMappingService interface {
	FindCities(ctx context.Context, params *models.ForwarderCityFilterRequest) ([]models.ForwarderCityMappingResponse, *models.PaginationMeta, error)
	FindSubdistricts(ctx context.Context, params *models.ForwarderSubdistrictFilterRequest) ([]models.ForwarderSubdistrictMappingResponse, *models.PaginationMeta, error)
	// Sync menarik data terbaru dari API Forwarder (citylist &/atau subdistrictlist)
	// dan upsert ke tabel mapping. Body opsional untuk memilih bagian yang di-sync.
	Sync(ctx context.Context, req *models.SyncForwarderMappingRequest) (*models.SyncForwarderMappingResponse, error)
}

type forwarderMappingService struct {
	repo repositories.ForwarderMappingRepository
}

func NewForwarderMappingService(repo repositories.ForwarderMappingRepository) ForwarderMappingService {
	return &forwarderMappingService{repo: repo}
}

// ─── Response API Forwarder ───────────────────────────────────────────────────

// forwarderListResponse adalah bentuk respons umum endpoint citylist/subdistrictlist:
// {"msg":"Success","data":[{"item_id":1,"item_name":"AMBON"}, ...],"isSuccess":"ok"}
type forwarderListResponse struct {
	Msg       string `json:"msg"`
	IsSuccess string `json:"isSuccess"`
	Data      []struct {
		ItemID   int    `json:"item_id"`
		ItemName string `json:"item_name"`
	} `json:"data"`
}

func (s *forwarderMappingService) FindCities(ctx context.Context, params *models.ForwarderCityFilterRequest) ([]models.ForwarderCityMappingResponse, *models.PaginationMeta, error) {
	items, total, err := s.repo.FindCities(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	result := make([]models.ForwarderCityMappingResponse, 0, len(items))
	for _, it := range items {
		result = append(result, models.ForwarderCityMappingResponse{
			ID:                it.ID,
			KotaPattern:       it.KotaPattern,
			ForwarderCityID:   it.ForwarderCityID,
			ForwarderCityName: it.ForwarderCityName,
			CreatedAt:         it.CreatedAt,
			UpdatedAt:         it.UpdatedAt,
		})
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return result, &meta, nil
}

func (s *forwarderMappingService) FindSubdistricts(ctx context.Context, params *models.ForwarderSubdistrictFilterRequest) ([]models.ForwarderSubdistrictMappingResponse, *models.PaginationMeta, error) {
	items, total, err := s.repo.FindSubdistricts(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	result := make([]models.ForwarderSubdistrictMappingResponse, 0, len(items))
	for _, it := range items {
		result = append(result, models.ForwarderSubdistrictMappingResponse{
			ID:                       it.ID,
			KecamatanPattern:         it.KecamatanPattern,
			ForwarderCityID:          it.ForwarderCityID,
			ForwarderSubdistrictID:   it.ForwarderSubdistrictID,
			ForwarderSubdistrictName: it.ForwarderSubdistrictName,
			CreatedAt:                it.CreatedAt,
			UpdatedAt:                it.UpdatedAt,
		})
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return result, &meta, nil
}

func (s *forwarderMappingService) Sync(ctx context.Context, req *models.SyncForwarderMappingRequest) (*models.SyncForwarderMappingResponse, error) {
	apiURL := os.Getenv("FORWARDER_API_URL")
	clientName := os.Getenv("FORWARDER_CLIENT_NAME")
	username := os.Getenv("FORWARDER_USERNAME")
	password := os.Getenv("FORWARDER_PASSWORD")
	if apiURL == "" || clientName == "" || username == "" || password == "" {
		return nil, fmt.Errorf("konfigurasi Forwarder tidak lengkap (FORWARDER_API_URL/CLIENT_NAME/USERNAME/PASSWORD)")
	}

	syncCity := true
	if req != nil && req.SyncCity != nil {
		syncCity = *req.SyncCity
	}
	syncSubdistrict := true
	if req != nil && req.SyncSubdistrict != nil {
		syncSubdistrict = *req.SyncSubdistrict
	}
	if !syncCity && !syncSubdistrict {
		return nil, fmt.Errorf("minimal salah satu dari sync_city atau sync_subdistrict harus true")
	}

	resp := &models.SyncForwarderMappingResponse{SyncedAt: time.Now().UTC()}

	if syncCity {
		token, err := s.getToken(ctx, apiURL, clientName, username, password, "CITYLIST")
		if err != nil {
			return nil, fmt.Errorf("gagal mendapatkan token Forwarder (CITYLIST): %w", err)
		}
		list, err := s.fetchList(ctx, apiURL, clientName, token, "/citylist", map[string]string{"city_name": ""})
		if err != nil {
			return nil, err
		}
		resp.CityTotalFromAPI = len(list)

		items := make([]models.ForwarderCityMapping, 0, len(list))
		for _, it := range list {
			items = append(items, models.ForwarderCityMapping{
				KotaPattern:       utils.NormalizeKota(it.ItemName),
				ForwarderCityID:   it.ItemID,
				ForwarderCityName: it.ItemName,
			})
		}
		created, updated, err := s.repo.UpsertCities(ctx, items)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan mapping kota: %w", err)
		}
		resp.CityCreated = created
		resp.CityUpdated = updated
	}

	if syncSubdistrict {
		token, err := s.getToken(ctx, apiURL, clientName, username, password, "SUBDISTRICTLIST")
		if err != nil {
			return nil, fmt.Errorf("gagal mendapatkan token Forwarder (SUBDISTRICTLIST): %w", err)
		}
		list, err := s.fetchList(ctx, apiURL, clientName, token, "/subdistrictlist", map[string]string{"subdistrict_name": "", "city_id": ""})
		if err != nil {
			return nil, err
		}
		resp.SubdistrictTotalFromAPI = len(list)

		items := make([]models.ForwarderSubdistrictMapping, 0, len(list))
		for _, it := range list {
			items = append(items, models.ForwarderSubdistrictMapping{
				KecamatanPattern:         utils.NormalizeKecamatan(it.ItemName),
				ForwarderCityID:          0, // list global tanpa city_id per kecamatan
				ForwarderSubdistrictID:   it.ItemID,
				ForwarderSubdistrictName: it.ItemName,
			})
		}
		created, updated, err := s.repo.UpsertSubdistricts(ctx, items)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan mapping kecamatan: %w", err)
		}
		resp.SubdistrictCreated = created
		resp.SubdistrictUpdated = updated
	}

	return resp, nil
}

// getToken mengambil access token Forwarder untuk scope tertentu (CITYLIST/SUBDISTRICTLIST).
func (s *forwarderMappingService) getToken(ctx context.Context, apiURL, clientName, username, password, scope string) (string, error) {
	body, _ := json.Marshal(map[string]string{"scope": scope})

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL+"/accesstoken", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("client_name", clientName)
	httpReq.Header.Set("username", username)
	httpReq.Header.Set("password", password)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("connection timeout saat minta token: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[forwarder-mapping] <-- POST /accesstoken scope=%s status=%d body=%s", scope, resp.StatusCode, string(respBody))

	var result struct {
		AccessToken string `json:"access_token"`
		Status      string `json:"status"`
		Message     string `json:"message"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("gagal parse response token: %w", err)
	}
	if result.Status != "ok" || result.AccessToken == "" {
		return "", fmt.Errorf("gagal mendapatkan token (%s): %s", scope, result.Message)
	}
	return result.AccessToken, nil
}

// fetchList memanggil endpoint list Forwarder (POST /citylist atau /subdistrictlist)
// dengan body opsional dan mem-parsing respons {"msg","data":[...]}.
func (s *forwarderMappingService) fetchList(ctx context.Context, apiURL, clientName, token, path string, bodyMap map[string]string) ([]struct {
	ItemID   int    `json:"item_id"`
	ItemName string `json:"item_name"`
}, error) {
	bodyBytes := []byte("{}")
	if bodyMap != nil {
		b, err := json.Marshal(bodyMap)
		if err != nil {
			return nil, err
		}
		bodyBytes = b
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL+path, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Client_name", clientName)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("connection timeout saat memanggil %s: %w", path, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[forwarder-mapping] <-- POST %s status=%d body=%s", path, resp.StatusCode, string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Forwarder API error saat Sync %s (status %d): %s", path, resp.StatusCode, string(respBody))
	}

	var result forwarderListResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("gagal parse response %s: %w", path, err)
	}
	if result.IsSuccess != "ok" {
		return nil, fmt.Errorf("Forwarder API error saat Sync %s: %s", path, result.Msg)
	}
	if result.Data == nil {
		result.Data = []struct {
			ItemID   int    `json:"item_id"`
			ItemName string `json:"item_name"`
		}{}
	}
	return result.Data, nil
}
