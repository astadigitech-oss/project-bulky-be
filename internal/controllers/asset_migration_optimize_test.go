package controllers

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"project-bulky-be/internal/config"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// generateTestJPEGBytes membuat bytes file JPEG dummy
func generateTestJPEGBytes(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 100, A: 255})
		}
	}
	buf := &bytes.Buffer{}
	_ = jpeg.Encode(buf, img, &jpeg.Options{Quality: 90})
	return buf.Bytes()
}

func TestOptimizeWebPEndpoint_DryRunValidation(t *testing.T) {
	app := fiber.New()
	uploadDir := t.TempDir()

	// Inisialisasi controller dummy (db nil/mocked handled safely)
	ctrl := &AssetMigrationController{
		db:  &gorm.DB{Config: &gorm.Config{Dialector: postgres.Dialector{}}},
		cfg: &config.Config{UploadPath: uploadDir},
	}
	app.Post("/optimize-webp", ctrl.OptimizeWebP)

	t.Run("Execute without DryRunToken should fail", func(t *testing.T) {
		reqBody := []byte(`{"dry_run": false}`)
		req := httptest.NewRequest(http.MethodPost, "/optimize-webp", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test failed: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400 for execute without token, got %d", resp.StatusCode)
		}
	})

	t.Run("Execute with invalid DryRunToken should fail", func(t *testing.T) {
		reqBody := []byte(`{"dry_run": false, "dry_run_token": "opt-invalid-token"}`)
		req := httptest.NewRequest(http.MethodPost, "/optimize-webp", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test failed: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400 for invalid token, got %d", resp.StatusCode)
		}
	})
}

func TestCompressExistingFileWebP(t *testing.T) {
	tempDir := t.TempDir()
	srcPath := filepath.Join(tempDir, "sample-old.jpg")
	dstPath := filepath.Join(tempDir, "sample-new.webp")

	// Buat file JPEG sumber di disk
	data := generateTestJPEGBytes(2500, 1800)
	if err := os.WriteFile(srcPath, data, 0644); err != nil {
		t.Fatalf("failed to write source image: %v", err)
	}

	// Tes helper utils.CompressExistingFileWebP
	err := utils.CompressExistingFileWebP(srcPath, dstPath)
	if err != nil {
		t.Fatalf("utils.CompressExistingFileWebP failed: %v", err)
	}

	// Pastikan file lama MASIH ADA di disk (tidak terhapus)
	if _, err := os.Stat(srcPath); err != nil {
		t.Errorf("source file should NOT be deleted, but got error: %v", err)
	}

	// Pastikan file baru tercipta dan ukurannya < 200KB
	dstInfo, err := os.Stat(dstPath)
	if err != nil {
		t.Fatalf("dest file does not exist: %v", err)
	}

	if dstInfo.Size() > 200*1024 {
		t.Errorf("dest file size %d exceeds 200KB limit", dstInfo.Size())
	}
}
