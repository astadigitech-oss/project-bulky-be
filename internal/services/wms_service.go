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
	"strconv"
	"strings"
	"sync"
	"time"

	"project-bulky-be/internal/models"
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
	// menetapkan harga jual cargo — tipe "discount" (persentase dari
	// total_price) atau "fix" (nilai adalah harga jual/sale_price akhir
	// secara langsung, BUKAN potongan dari total_price).
	SetCargoPrice(ctx context.Context, cargoID string, req *models.SetWMSCargoPriceRequest) (*models.WMSCargoPriceResponse, error)
	// ListAlreadyPricedCargos memanggil GET /api/integration/cargos/already-priced
	// untuk mendapatkan daftar cargo yang sudah diberi harga di WMS tapi belum
	// dikonfirmasi sinkron — sumber dropdown "ID Cargo" saat create produk.
	ListAlreadyPricedCargos(ctx context.Context, search string) ([]models.WMSCargoPricedResponse, error)
	// DownloadCargoPricingPDF memanggil GET /api/integration/cargos/{id}/pricing-pdf
	// dan mengembalikan isi PDF mentah (bukan disimpan ke storage lokal) — FE
	// yang menyimpannya sebagai file dokumen produk seolah diupload manual.
	DownloadCargoPricingPDF(ctx context.Context, cargoID string) ([]byte, error)
	// MarkCargoSynced memanggil POST /api/integration/cargos/{id}/status untuk
	// menandai cargo sudah dikonfirmasi sinkron (is_sync = true) di WMS.
	// Idempotent.
	MarkCargoSynced(ctx context.Context, cargoID string) (*models.WMSCargoSyncStatusResponse, error)
}

type wmsService struct {
	baseURL      string
	clientID     string
	clientSecret string

	mu          sync.Mutex
	cachedToken string
	expiresAt   time.Time
}

func NewWMSService(baseURL, clientID, clientSecret string) WMSService {
	return &wmsService{
		baseURL:      strings.TrimRight(baseURL, "/"),
		clientID:     clientID,
		clientSecret: clientSecret,
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
// Bearer token untuk menetapkan harga jual cargo — tipe "discount"
// (persentase dari total_price) atau "fix" (harga jual/sale_price akhir
// secara langsung).
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

	return &envelope.Data, nil
}

// ListAlreadyPricedCargos memanggil GET /api/integration/cargos/already-priced
// dengan Bearer token untuk mendapatkan daftar cargo yang sudah diberi harga
// di WMS tapi belum dikonfirmasi sinkron — sumber dropdown "ID Cargo" saat
// create produk. Pure passthrough, tidak ada cache/filter lokal — WMS sendiri
// yang menentukan cargo mana yang masih tersedia (is_sync belum true).
func (s *wmsService) ListAlreadyPricedCargos(ctx context.Context, search string) ([]models.WMSCargoPricedResponse, error) {
	token, err := s.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if search != "" {
		query.Set("search", search)
	}
	query.Set("limit", "100")

	reqURL := s.baseURL + "/api/integration/cargos/already-priced?" + query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("connection timeout saat ambil daftar cargo sudah diberi harga: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[wms] <-- GET /api/integration/cargos/already-priced status=%d body=%s", resp.StatusCode, string(respBody))

	if resp.StatusCode == http.StatusUnauthorized {
		s.mu.Lock()
		s.cachedToken = ""
		s.mu.Unlock()
		return nil, fmt.Errorf("token WMS tidak ada / salah / kedaluwarsa / kredensial dicabut")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WMS API error saat ambil daftar cargo sudah diberi harga (status %d): %s", resp.StatusCode, string(respBody))
	}

	var envelope models.WMSCargoPricedListEnvelope
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("gagal parse response daftar cargo sudah diberi harga: %w", err)
	}
	if !envelope.Success {
		return nil, fmt.Errorf("WMS API mengembalikan gagal: %s", envelope.Message)
	}

	return envelope.Data, nil
}

// DownloadCargoPricingPDF memanggil GET /api/integration/cargos/{id}/pricing-pdf
// dan mengembalikan isi PDF mentah — tidak disimpan ke storage lokal, FE yang
// mengunggahnya kembali sebagai dokumen produk seolah file diupload manual.
func (s *wmsService) DownloadCargoPricingPDF(ctx context.Context, cargoID string) ([]byte, error) {
	token, err := s.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	reqURL := s.baseURL + "/api/integration/cargos/" + url.PathEscape(cargoID) + "/pricing-pdf"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("connection timeout saat download PDF harga cargo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		s.mu.Lock()
		s.cachedToken = ""
		s.mu.Unlock()
		return nil, fmt.Errorf("token WMS tidak ada / salah / kedaluwarsa / kredensial dicabut")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("WMS API error saat download PDF harga cargo (status %d): %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// MarkCargoSynced memanggil POST /api/integration/cargos/{id}/status untuk
// menandai cargo sudah dikonfirmasi sinkron (is_sync = true) di WMS —
// idempotent, aman dipanggil berkali-kali.
func (s *wmsService) MarkCargoSynced(ctx context.Context, cargoID string) (*models.WMSCargoSyncStatusResponse, error) {
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

	return &envelope.Data, nil
}
