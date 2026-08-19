package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
)

// WMSService mengelola integrasi OAuth client_credentials ke WMS (Warehouse
// Management System) — dipakai sebagai fondasi untuk fitur sync produk palet
// dari inventory WMS jadi cargo online.
//
// Alur:
//  1. POST /api/oauth/token (publik, tanpa Authorization header) menukar
//     client_id/client_secret jadi access token. Token berlaku singkat (±15 menit).
//  2. Endpoint bisnis lain (mis. GET /api/integration/me) dipanggil dengan
//     header Authorization: Bearer <access_token>.
//
// Token di-cache in-memory (thread-safe) dan otomatis diminta ulang saat
// hampir/sudah kedaluwarsa, supaya tidak perlu request token di setiap call.
type WMSService interface {
	// GetAccessToken mengembalikan access token WMS yang valid, meminta token
	// baru ke API kalau cache kosong/hampir kedaluwarsa.
	GetAccessToken(ctx context.Context) (string, error)
	// TestConnection memanggil GET /api/integration/me untuk memverifikasi
	// kredensial WMS aktif & terkoneksi. Dipanggil setelah dapat token, sebelum
	// mencoba endpoint bisnis lain.
	TestConnection(ctx context.Context) (*models.WMSConnectionInfo, error)
	// ListReadyToPriceCargos memanggil GET /api/integration/cargos/ready-to-price
	// untuk mendapatkan daftar cargo (ukuran fisik lengkap, belum pernah
	// dihargai, belum disinkronkan) yang siap ditetapkan harga jualnya.
	ListReadyToPriceCargos(ctx context.Context, params *models.WMSCargoListFilterRequest) ([]models.WMSCargoPricingResponse, *models.WMSPaginationMetaRaw, error)
	// SetCargoPrice memanggil POST /api/integration/cargos/{id}/price untuk
	// menetapkan harga jual cargo (diskon persentase atau potongan rupiah
	// fix dari total_price/harga inventory). Setelah sukses, PDF harga
	// (pricing_pdf_url) langsung di-download & disimpan ke storage lokal, lalu
	// cache cargo di tabel wms_cargo_priced di-upsert supaya siap dipilih di
	// dropdown "ID Cargo" saat create/edit produk.
	SetCargoPrice(ctx context.Context, cargoID string, req *models.SetWMSCargoPriceRequest) (*models.WMSCargoPriceResponse, error)
	// ListAlreadyPricedCargos mengembalikan cache lokal cargo yang sudah
	// diberi harga tapi belum dipakai di produk manapun (tabel
	// wms_cargo_priced), dipakai sebagai sumber dropdown "ID Cargo".
	ListAlreadyPricedCargos(ctx context.Context, search string) ([]models.WMSCargoPricedListResponse, error)
	// MarkCargoSynced memanggil POST /api/integration/cargos/{id}/status untuk
	// menandai cargo sudah dikonfirmasi sinkron (is_sync = true) di WMS, lalu
	// menandai cache lokal sebagai sudah dipakai di produk terkait. Idempotent.
	MarkCargoSynced(ctx context.Context, cargoID string, produkID string) (*models.WMSCargoSyncStatusResponse, error)
}

type wmsService struct {
	baseURL      string
	clientID     string
	clientSecret string
	cfg          *config.Config
	cargoRepo    repositories.WMSCargoPricedRepository

	mu          sync.Mutex
	cachedToken string
	expiresAt   time.Time
}

func NewWMSService(baseURL, clientID, clientSecret string, cfg *config.Config, cargoRepo repositories.WMSCargoPricedRepository) WMSService {
	return &wmsService{
		baseURL:      strings.TrimRight(baseURL, "/"),
		clientID:     clientID,
		clientSecret: clientSecret,
		cfg:          cfg,
		cargoRepo:    cargoRepo,
	}
}

// wmsTokenEnvelope bentuk respons POST /api/oauth/token:
// {"success":true,"message":"...","data":{"access_token":"...","token_type":"Bearer","expires_in":899,"expires_at":"..."}}
type wmsTokenEnvelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		ExpiresAt   string `json:"expires_at"`
	} `json:"data"`
}

// GetAccessToken mengembalikan token dari cache jika masih valid (dengan buffer
// 30 detik sebelum kedaluwarsa), atau meminta token baru ke WMS.
func (s *wmsService) GetAccessToken(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cachedToken != "" && time.Now().Before(s.expiresAt.Add(-30*time.Second)) {
		return s.cachedToken, nil
	}

	token, expiresAt, err := s.fetchToken(ctx)
	if err != nil {
		return "", err
	}

	s.cachedToken = token
	s.expiresAt = expiresAt
	return token, nil
}

// fetchToken menukar client_id/client_secret jadi access token via
// POST /api/oauth/token. Publik, tidak butuh Authorization header.
func (s *wmsService) fetchToken(ctx context.Context) (string, time.Time, error) {
	if s.baseURL == "" || s.clientID == "" || s.clientSecret == "" {
		return "", time.Time{}, fmt.Errorf("konfigurasi WMS tidak lengkap (WMS_BASE_URL/WMS_CLIENT_ID/WMS_CLIENT_SECRET)")
	}

	body, _ := json.Marshal(map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     s.clientID,
		"client_secret": s.clientSecret,
	})

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/api/oauth/token", bytes.NewReader(body))
	if err != nil {
		return "", time.Time{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("connection timeout saat minta token WMS: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[wms] <-- POST /api/oauth/token status=%d body=%s", resp.StatusCode, string(respBody))

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("WMS API error saat minta token (status %d): %s", resp.StatusCode, string(respBody))
	}

	var envelope wmsTokenEnvelope
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return "", time.Time{}, fmt.Errorf("gagal parse response token WMS: %w", err)
	}
	if !envelope.Success || envelope.Data.AccessToken == "" {
		return "", time.Time{}, fmt.Errorf("gagal mendapatkan token WMS: %s", envelope.Message)
	}

	expiresAt := time.Now().Add(time.Duration(envelope.Data.ExpiresIn) * time.Second)
	if envelope.Data.ExpiresAt != "" {
		if parsed, err := time.Parse(time.RFC3339, envelope.Data.ExpiresAt); err == nil {
			expiresAt = parsed
		}
	}

	return envelope.Data.AccessToken, expiresAt, nil
}

// TestConnection memanggil GET /api/integration/me dengan Bearer token untuk
// memverifikasi kredensial WMS aktif & terkoneksi.
func (s *wmsService) TestConnection(ctx context.Context) (*models.WMSConnectionInfo, error) {
	token, err := s.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/api/integration/me", nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("connection timeout saat cek koneksi WMS: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[wms] <-- GET /api/integration/me status=%d body=%s", resp.StatusCode, string(respBody))

	if resp.StatusCode == http.StatusUnauthorized {
		// Token mungkin sudah dicabut/invalid di sisi WMS meski masih dianggap
		// valid di cache lokal — bersihkan cache agar request berikutnya minta
		// token baru, bukan mengulang token basi.
		s.mu.Lock()
		s.cachedToken = ""
		s.mu.Unlock()
		return nil, fmt.Errorf("token WMS tidak ada / salah / kedaluwarsa / kredensial dicabut")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WMS API error saat cek koneksi (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result models.WMSConnectionInfo
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("gagal parse response cek koneksi WMS: %w", err)
	}
	if !result.Success {
		return nil, fmt.Errorf("WMS API mengembalikan gagal: %s", result.Message)
	}

	return &result, nil
}

// ListReadyToPriceCargos memanggil GET /api/integration/cargos/ready-to-price
// dengan Bearer token untuk mendapatkan daftar cargo yang siap ditetapkan
// harga jualnya (ukuran fisik lengkap, belum pernah dihargai, belum
// disinkronkan).
func (s *wmsService) ListReadyToPriceCargos(ctx context.Context, params *models.WMSCargoListFilterRequest) ([]models.WMSCargoPricingResponse, *models.WMSPaginationMetaRaw, error) {
	token, err := s.GetAccessToken(ctx)
	if err != nil {
		return nil, nil, err
	}

	if params == nil {
		params = &models.WMSCargoListFilterRequest{}
	}
	params.SetDefaults()

	query := url.Values{}
	query.Set("page", strconv.Itoa(params.Page))
	query.Set("limit", strconv.Itoa(params.Limit))
	if params.Search != "" {
		query.Set("search", params.Search)
	}

	reqURL := s.baseURL + "/api/integration/cargos/ready-to-price?" + query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, nil, fmt.Errorf("connection timeout saat ambil daftar cargo siap harga WMS: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[wms] <-- GET /api/integration/cargos/ready-to-price status=%d body=%s", resp.StatusCode, string(respBody))

	if resp.StatusCode == http.StatusUnauthorized {
		s.mu.Lock()
		s.cachedToken = ""
		s.mu.Unlock()
		return nil, nil, fmt.Errorf("token WMS tidak ada / salah / kedaluwarsa / kredensial dicabut")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("WMS API error saat ambil daftar cargo siap harga (status %d): %s", resp.StatusCode, string(respBody))
	}

	var envelope models.WMSCargoListEnvelope
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, nil, fmt.Errorf("gagal parse response daftar cargo siap harga WMS: %w", err)
	}
	if !envelope.Success {
		return nil, nil, fmt.Errorf("WMS API mengembalikan gagal: %s", envelope.Message)
	}

	return envelope.Data, &envelope.Meta, nil
}

// SetCargoPrice memanggil POST /api/integration/cargos/{id}/price dengan
// Bearer token untuk menetapkan harga jual cargo (diskon persentase atau
// potongan rupiah fix dari total_price).
func (s *wmsService) SetCargoPrice(ctx context.Context, cargoID string, req *models.SetWMSCargoPriceRequest) (*models.WMSCargoPriceResponse, error) {
	token, err := s.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("gagal encode request harga cargo: %w", err)
	}

	reqURL := s.baseURL + "/api/integration/cargos/" + url.PathEscape(cargoID) + "/price"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("connection timeout saat menetapkan harga cargo WMS: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[wms] <-- POST /api/integration/cargos/%s/price status=%d body=%s", cargoID, resp.StatusCode, string(respBody))

	if resp.StatusCode == http.StatusUnauthorized {
		s.mu.Lock()
		s.cachedToken = ""
		s.mu.Unlock()
		return nil, fmt.Errorf("token WMS tidak ada / salah / kedaluwarsa / kredensial dicabut")
	}
	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
		var errEnvelope models.WMSErrorEnvelope
		if err := json.Unmarshal(respBody, &errEnvelope); err == nil && errEnvelope.Message != "" {
			return nil, fmt.Errorf("WMS API error saat menetapkan harga cargo (status %d): %s", resp.StatusCode, errEnvelope.Message)
		}
		return nil, fmt.Errorf("WMS API error saat menetapkan harga cargo (status %d): %s", resp.StatusCode, string(respBody))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WMS API error saat menetapkan harga cargo (status %d): %s", resp.StatusCode, string(respBody))
	}

	var envelope models.WMSCargoPriceEnvelope
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("gagal parse response harga cargo WMS: %w", err)
	}
	if !envelope.Success {
		return nil, fmt.Errorf("WMS API mengembalikan gagal: %s", envelope.Message)
	}

	// Download PDF harga terbaru & simpan ke storage lokal, lalu upsert cache
	// cargo di tabel wms_cargo_priced. Kegagalan di langkah ini TIDAK
	// menggagalkan SetCargoPrice — harga sudah ter-set di WMS, PDF/cache bisa
	// direfresh belakangan (mis. saat cargo muncul lagi di ready-to-price
	// setelah diset ulang, atau retry manual).
	pdfPath, err := s.downloadAndSavePricingPDF(ctx, token, cargoID, envelope.Data.PricingPDFURL)
	if err != nil {
		log.Printf("[wms] WARNING: gagal download/simpan PDF harga cargo %s: %v", cargoID, err)
	}

	if s.cargoRepo != nil {
		cargoUUID, parseErr := uuid.Parse(cargoID)
		if parseErr == nil {
			pricingType := envelope.Data.PricingType
			pricingValue := envelope.Data.PricingValue
			salePrice := envelope.Data.SalePrice
			pricedAt := envelope.Data.PricedAt
			cache := &models.WMSCargoPriced{
				CargoID:      cargoUUID,
				Code:         envelope.Data.Code,
				TotalPrice:   envelope.Data.TotalPrice,
				PricingType:  &pricingType,
				PricingValue: &pricingValue,
				SalePrice:    &salePrice,
				PricedAt:     &pricedAt,
			}
			if pdfPath != "" {
				cache.PricingPDFPath = &pdfPath
			}
			if upsertErr := s.cargoRepo.Upsert(ctx, cache); upsertErr != nil {
				log.Printf("[wms] WARNING: gagal upsert cache cargo %s: %v", cargoID, upsertErr)
			}
		}
	}

	return &envelope.Data, nil
}

// downloadAndSavePricingPDF mengunduh PDF harga terbaru dari WMS
// (GET /api/integration/cargos/{id}/pricing-pdf) dan menyimpannya ke storage
// lokal (UPLOAD_PATH/wms-cargo/<cargo_id>/pricing.pdf, selalu ditimpa dengan
// nama tetap supaya URL publik konsisten walau harga diset ulang). Mengembalikan
// path relatif (untuk disimpan ke DB) atau error kalau gagal.
func (s *wmsService) downloadAndSavePricingPDF(ctx context.Context, token, cargoID, pdfURLHint string) (string, error) {
	if s.cfg == nil {
		return "", fmt.Errorf("konfigurasi upload tidak tersedia")
	}

	reqURL := s.baseURL + "/api/integration/cargos/" + url.PathEscape(cargoID) + "/pricing-pdf"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("connection timeout saat download PDF harga cargo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("WMS API error saat download PDF harga cargo (status %d): %s", resp.StatusCode, string(body))
	}

	relativeDir := filepath.Join("wms-cargo", cargoID)
	uploadDir := filepath.Join(s.cfg.UploadPath, relativeDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori upload PDF cargo: %w", err)
	}

	fullPath := filepath.Join(uploadDir, "pricing.pdf")
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file PDF cargo: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, resp.Body); err != nil {
		return "", fmt.Errorf("gagal menyimpan file PDF cargo: %w", err)
	}

	relativePath := strings.ReplaceAll(filepath.Join(relativeDir, "pricing.pdf"), "\\", "/")
	log.Printf("[wms] PDF harga cargo %s tersimpan di %s", cargoID, relativePath)
	return relativePath, nil
}

// ListAlreadyPricedCargos mengembalikan cache lokal cargo yang sudah diberi
// harga tapi belum dipakai di produk manapun (tabel wms_cargo_priced).
func (s *wmsService) ListAlreadyPricedCargos(ctx context.Context, search string) ([]models.WMSCargoPricedListResponse, error) {
	if s.cargoRepo == nil {
		return nil, fmt.Errorf("cargo repository tidak tersedia")
	}

	items, err := s.cargoRepo.ListUnused(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar cargo sudah diberi harga: %w", err)
	}

	result := make([]models.WMSCargoPricedListResponse, 0, len(items))
	for _, item := range items {
		salePrice := 0.0
		if item.SalePrice != nil {
			salePrice = *item.SalePrice
		}
		var pdfURL *string
		if item.PricingPDFPath != nil && *item.PricingPDFPath != "" {
			full := utils.GetFileURL(*item.PricingPDFPath, s.cfg)
			pdfURL = &full
		}
		result = append(result, models.WMSCargoPricedListResponse{
			CargoID:        item.CargoID.String(),
			Code:           item.Code,
			LengthCM:       item.LengthCM,
			WidthCM:        item.WidthCM,
			HeightCM:       item.HeightCM,
			WeightKG:       item.WeightKG,
			TotalPrice:     item.TotalPrice,
			SalePrice:      salePrice,
			PricingPDFURL:  pdfURL,
			IsUsedInProduk: item.IsUsedInProduk,
		})
	}
	return result, nil
}

// MarkCargoSynced memanggil POST /api/integration/cargos/{id}/status untuk
// menandai cargo sudah dikonfirmasi sinkron (is_sync = true) di WMS —
// idempotent, aman dipanggil berkali-kali. Setelah sukses, cache lokal
// ditandai sudah dipakai di produk terkait (produkID boleh kosong kalau
// hanya ingin menandai sinkron ke WMS tanpa mengaitkan ke produk tertentu).
func (s *wmsService) MarkCargoSynced(ctx context.Context, cargoID string, produkID string) (*models.WMSCargoSyncStatusResponse, error) {
	token, err := s.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	reqURL := s.baseURL + "/api/integration/cargos/" + url.PathEscape(cargoID) + "/status"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("connection timeout saat menandai cargo sinkron: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[wms] <-- POST /api/integration/cargos/%s/status status=%d body=%s", cargoID, resp.StatusCode, string(respBody))

	if resp.StatusCode == http.StatusUnauthorized {
		s.mu.Lock()
		s.cachedToken = ""
		s.mu.Unlock()
		return nil, fmt.Errorf("token WMS tidak ada / salah / kedaluwarsa / kredensial dicabut")
	}
	if resp.StatusCode == http.StatusBadRequest {
		var errEnvelope models.WMSErrorEnvelope
		if err := json.Unmarshal(respBody, &errEnvelope); err == nil && errEnvelope.Message != "" {
			return nil, fmt.Errorf("cargo belum pernah diberi harga: %s", errEnvelope.Message)
		}
		return nil, fmt.Errorf("cargo belum pernah diberi harga (status 400): %s", string(respBody))
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("cargo tidak ditemukan di WMS")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WMS API error saat menandai cargo sinkron (status %d): %s", resp.StatusCode, string(respBody))
	}

	var envelope models.WMSCargoSyncStatusEnvelope
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("gagal parse response status sinkron cargo: %w", err)
	}
	if !envelope.Success {
		return nil, fmt.Errorf("WMS API mengembalikan gagal: %s", envelope.Message)
	}

	if s.cargoRepo != nil && produkID != "" {
		if err := s.cargoRepo.MarkUsed(ctx, cargoID, produkID); err != nil {
			log.Printf("[wms] WARNING: gagal menandai cache cargo %s sudah dipakai: %v", cargoID, err)
		}
	}

	return &envelope.Data, nil
}
