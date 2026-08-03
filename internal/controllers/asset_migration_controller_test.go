package controllers

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"project-bulky-be/internal/config"

	"gorm.io/gorm"
)

// TestProcessV1ZipMapping membuat ZIP sampel berisi file dengan struktur
// storage Laravel v1, lalu memastikan processV1Zip memetakannya ke folder
// tujuan v2 yang benar dan men-skip file yang tidak dikenali.
func TestProcessV1ZipMapping(t *testing.T) {
	uploadDir := t.TempDir()

	// Buat ZIP sampel
	zipPath := filepath.Join(t.TempDir(), "sample-v1.zip")
	zf, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(zf)
	entries := []struct{ name, content string }{
		{"storage/app/public/products/produk-a.jpg", "image-a"},
		{"storage/app/public/reviews/review-1.jpg", "review-1"},
		{"storage/app/public/public/profile/user-1.png", "profile-1"},
		{"storage/app/public/br001.pdf", "pdf-1"},
		{"storage/app/public/br001.jpg", "jpg-root"},
		{"storage/app/private/secret.txt", "skip"},
		{"storage/app/public/unknown-folder/foo.txt", "skip"},
	}
	for _, e := range entries {
		w, err := zw.Create(e.name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(e.content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	zf.Close()

	ctrl := &AssetMigrationController{db: &gorm.DB{}, cfg: &config.Config{UploadPath: uploadDir}}
	res, err := ctrl.processV1Zip(zipPath)
	if err != nil {
		t.Fatalf("processV1Zip error: %v", err)
	}

	if res.Imported != 4 {
		t.Errorf("imported = %d, want 4", res.Imported)
	}
	if res.Skipped != 3 {
		t.Errorf("skipped = %d, want 3", res.Skipped)
	}

	// Cek file hasil ekstrak
	want := map[string]bool{
		"product-images/produk-a.jpg": true,
		"reviews/review-1.jpg":        true,
		"profile/user-1.png":          true,
		"product-documents/br001.pdf": true,
	}
	for rel := range want {
		if _, err := os.Stat(filepath.Join(uploadDir, filepath.FromSlash(rel))); err != nil {
			t.Errorf("file %s tidak diekstrak: %v", rel, err)
		}
	}

	// File non-PDF di root harus di-skip
	if _, err := os.Stat(filepath.Join(uploadDir, "product-documents/br001.jpg")); err == nil {
		t.Error("file jpg di root seharusnya di-skip, tapi ada di hasil")
	}
}
