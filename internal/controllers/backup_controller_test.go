package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/services"

	"github.com/gofiber/fiber/v2"
)

func TestBackupControllerEndpoints(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		BackupPath:          tempDir,
		BackupRetentionDays: 7,
	}

	// Create sample backup files
	sampleFile1 := "backup_bulky_2026-08-24_100000_auto.sql.gz"
	sampleFile2 := "backup_bulky_2026-08-24_110000_manual.sql.gz"
	os.WriteFile(filepath.Join(tempDir, sampleFile1), []byte("sample-data-1"), 0644)
	os.WriteFile(filepath.Join(tempDir, sampleFile2), []byte("sample-data-2"), 0644)

	svc := services.NewBackupService(cfg, nil)
	ctrl := NewBackupController(svc)

	app := fiber.New()
	app.Get("/api/panel/backups", ctrl.List)
	app.Post("/api/panel/backups", ctrl.Create)
	app.Get("/api/panel/backups/:filename/download", ctrl.Download)
	app.Delete("/api/panel/backups/:filename", ctrl.Delete)

	// 1. Test GET /api/panel/backups
	req := httptest.NewRequest(http.MethodGet, "/api/panel/backups", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /api/panel/backups failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var listResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&listResp)
	if !listResp["success"].(bool) {
		t.Errorf("Expected success = true, got %v", listResp["success"])
	}

	// 2. Test GET /api/panel/backups/:filename/download
	req = httptest.NewRequest(http.MethodGet, "/api/panel/backups/"+sampleFile1+"/download", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Download request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for download, got %d", resp.StatusCode)
	}

	// 3. Test GET Download with path traversal attempt
	req = httptest.NewRequest(http.MethodGet, "/api/panel/backups/..%2F..%2Fetc%2Fpasswd/download", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Download path traversal request failed: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected non-200 for path traversal, got %d", resp.StatusCode)
	}

	// 4. Test DELETE /api/panel/backups/:filename
	req = httptest.NewRequest(http.MethodDelete, "/api/panel/backups/"+sampleFile2, nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Delete request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for delete, got %d", resp.StatusCode)
	}

	// Verify file is gone from disk
	if _, err := os.Stat(filepath.Join(tempDir, sampleFile2)); !os.IsNotExist(err) {
		t.Errorf("Expected %s to be deleted from disk", sampleFile2)
	}
}
