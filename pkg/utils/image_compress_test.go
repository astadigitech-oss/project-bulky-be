package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"project-bulky-be/internal/config"
)

// createMultipartFileHeader membuat dummy multipart.FileHeader dari bytes gambar di memori
func createMultipartFileHeader(t *testing.T, filename, contentType string, data []byte) *multipart.FileHeader {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		t.Fatalf("failed to create multipart part: %v", err)
	}

	if _, err := part.Write(data); err != nil {
		t.Fatalf("failed to write data to part: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(int64(body.Len()))
	if err != nil {
		t.Fatalf("failed to read form: %v", err)
	}

	files := form.File["file"]
	if len(files) == 0 {
		t.Fatalf("no file found in form")
	}

	return files[0]
}

// generateTestJPEG membuat data dummy gambar JPEG dengan dimensi tertentu
func generateTestJPEG(t *testing.T, width, height int) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Gambar pola warna-warni agar tidak terlalu gampang dikompres (simulasi foto asli)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8((x * 255) / width),
				G: uint8((y * 255) / height),
				B: uint8(((x + y) * 255) / (width + height)),
				A: 255,
			})
		}
	}

	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 95}); err != nil {
		t.Fatalf("failed to encode jpeg: %v", err)
	}

	return buf.Bytes()
}

// generateTestPNG membuat data dummy gambar PNG dengan alpha channel (transparan)
func generateTestPNG(t *testing.T, width, height int) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Transparan di pojok, berwarna di tengah
	draw.Draw(img, image.Rect(width/4, height/4, (3*width)/4, (3*height)/4), &image.Uniform{color.RGBA{255, 0, 0, 200}}, image.Point{}, draw.Src)

	buf := &bytes.Buffer{}
	if err := png.Encode(buf, img); err != nil {
		t.Fatalf("failed to encode png: %v", err)
	}

	return buf.Bytes()
}

func TestCompressAndSaveImageWebP(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-uploads-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := &config.Config{
		UploadPath: tempDir,
	}

	t.Run("Compress Large JPEG to WebP < 200KB", func(t *testing.T) {
		// Buat gambar beresolusi 3000x2000 (landscape besar)
		jpegData := generateTestJPEG(t, 3000, 2000)
		fileHeader := createMultipartFileHeader(t, "sample-large.jpg", "image/jpeg", jpegData)

		relPath, err := CompressAndSaveImageWebP(fileHeader, "products", cfg)
		if err != nil {
			t.Fatalf("CompressAndSaveImageWebP error: %v", err)
		}

		if filepath.Ext(relPath) != ".webp" {
			t.Errorf("expected extension .webp, got %s", filepath.Ext(relPath))
		}

		fullPath := filepath.Join(cfg.UploadPath, relPath)
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			t.Fatalf("output file does not exist: %v", err)
		}

		t.Logf("Original JPEG size: %d bytes, Compressed WebP size: %d bytes (%.2f KB)",
			len(jpegData), fileInfo.Size(), float64(fileInfo.Size())/1024.0)

		if fileInfo.Size() > TargetMaxImageSize {
			t.Errorf("compressed file size %d bytes exceeds target %d bytes", fileInfo.Size(), TargetMaxImageSize)
		}

		// Pastikan dimensi proporsional (lebar max 1920, rasio 3000:2000 = 1.5 -> tinggi 1280)
		w, h, err := getImageDimensions(fullPath)
		if err != nil {
			t.Fatalf("failed to read output image dimensions: %v", err)
		}
		if w != 1920 || h != 1280 {
			t.Errorf("expected dimensions 1920x1280, got %dx%d", w, h)
		}
	})

	t.Run("Compress Portrait PNG with Transparency", func(t *testing.T) {
		// Buat gambar portrait 1000x2000
		pngData := generateTestPNG(t, 1000, 2000)
		fileHeader := createMultipartFileHeader(t, "sample-transparent.png", "image/png", pngData)

		relPath, err := CompressAndSaveImageWebP(fileHeader, "products", cfg)
		if err != nil {
			t.Fatalf("CompressAndSaveImageWebP error: %v", err)
		}

		fullPath := filepath.Join(cfg.UploadPath, relPath)
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			t.Fatalf("output file does not exist: %v", err)
		}

		if fileInfo.Size() > TargetMaxImageSize {
			t.Errorf("compressed file size %d bytes exceeds target %d bytes", fileInfo.Size(), TargetMaxImageSize)
		}

		// Portrait: tinggi max 1920, rasio 1000:2000 = 0.5 -> lebar 960
		w, h, err := getImageDimensions(fullPath)
		if err != nil {
			t.Fatalf("failed to read output image dimensions: %v", err)
		}
		if h != 1920 || w != 960 {
			t.Errorf("expected dimensions 960x1920, got %dx%d", w, h)
		}
	})

	t.Run("Compress Small Image Without Upscaling", func(t *testing.T) {
		// Buat gambar kecil 400x300 (harus tetap 400x300)
		jpegData := generateTestJPEG(t, 400, 300)
		fileHeader := createMultipartFileHeader(t, "small.jpg", "image/jpeg", jpegData)

		relPath, err := CompressAndSaveImageWebP(fileHeader, "products", cfg)
		if err != nil {
			t.Fatalf("CompressAndSaveImageWebP error: %v", err)
		}

		fullPath := filepath.Join(cfg.UploadPath, relPath)
		w, h, err := getImageDimensions(fullPath)
		if err != nil {
			t.Fatalf("failed to read output image dimensions: %v", err)
		}
		if w != 400 || h != 300 {
			t.Errorf("expected small image to preserve original dimensions 400x300, got %dx%d", w, h)
		}
	})

	t.Run("Reject Invalid Image Type", func(t *testing.T) {
		dummyData := []byte("%PDF-1.4 dummy pdf content")
		fileHeader := createMultipartFileHeader(t, "document.pdf", "application/pdf", dummyData)

		_, err := CompressAndSaveImageWebP(fileHeader, "products", cfg)
		if err == nil {
			t.Error("expected error for non-image file, got nil")
		}
	})
}
