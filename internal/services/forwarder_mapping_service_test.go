package services

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// TestForwarderMappingFetchListLive memverifikasi bahwa layer HTTP service
// (getToken + fetchList) berhasil menarik data dari API Forwarder dengan body
// yang benar — persis alur yang dipakai Sync(). Dijalankan otomatis ketika
// env FORWARDER_* tersedia (mis. setelah `go test ./internal/services/`),
// dan di-skip jika tidak ada network/env.
func TestForwarderMappingFetchListLive(t *testing.T) {
	_ = godotenv.Load()

	apiURL := os.Getenv("FORWARDER_API_URL")
	clientName := os.Getenv("FORWARDER_CLIENT_NAME")
	username := os.Getenv("FORWARDER_USERNAME")
	password := os.Getenv("FORWARDER_PASSWORD")
	if apiURL == "" || clientName == "" || username == "" || password == "" {
		t.Skip("FORWARDER_* env tidak tersedia — skip live test")
	}

	svc := &forwarderMappingService{}
	ctx := context.Background()

	// ── Citylist ────────────────────────────────────────────────
	token, err := svc.getToken(ctx, apiURL, clientName, username, password, "CITYLIST")
	if err != nil {
		t.Fatalf("getToken CITYLIST: %v", err)
	}
	if token == "" {
		t.Fatal("token CITYLIST kosong")
	}

	cities, err := svc.fetchList(ctx, apiURL, clientName, token, "/citylist", map[string]string{"city_name": ""})
	if err != nil {
		t.Fatalf("fetchList /citylist: %v", err)
	}
	if len(cities) == 0 {
		t.Fatal("/citylist mengembalikan 0 data")
	}
	t.Logf("/citylist OK: %d kota, contoh: %d %s", len(cities), cities[0].ItemID, cities[0].ItemName)

	// ── Subdistrictlist ─────────────────────────────────────────
	token, err = svc.getToken(ctx, apiURL, clientName, username, password, "SUBDISTRICTLIST")
	if err != nil {
		t.Fatalf("getToken SUBDISTRICTLIST: %v", err)
	}
	if token == "" {
		t.Fatal("token SUBDISTRICTLIST kosong")
	}

	subdistricts, err := svc.fetchList(ctx, apiURL, clientName, token, "/subdistrictlist", map[string]string{"subdistrict_name": "", "city_id": ""})
	if err != nil {
		t.Fatalf("fetchList /subdistrictlist: %v", err)
	}
	if len(subdistricts) == 0 {
		t.Fatal("/subdistrictlist mengembalikan 0 data")
	}
	t.Logf("/subdistrictlist OK: %d kecamatan, contoh: %d %s", len(subdistricts), subdistricts[0].ItemID, subdistricts[0].ItemName)
}
