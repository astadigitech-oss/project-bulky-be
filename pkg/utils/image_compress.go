package utils

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"

	"project-bulky-be/internal/config"

	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
)

const (
	// TargetMaxImageSize adalah batas ukuran gambar akhir (< 200KB)
	TargetMaxImageSize int64 = 200 * 1024 // 204,800 bytes

	// MaxInputUploadSize batas input file dari user (10MB)
	MaxInputUploadSize int64 = 10 * 1024 * 1024

	// MaxImageResolution sisi terpanjang maksimal (1920px) untuk menjaga ketajaman tanpa bloat
	MaxImageResolution = 1920

	// SecondPassResolution resolusi clamp cadangan jika file awal masih > 200KB
	SecondPassResolution = 1440
)

// CompressAndSaveImageWebP mengonversi file gambar ke WebP dan mengompres ukurannya di bawah 200KB.
// Aspek rasio gambar asli selalu dipertahankan (anti-gepeng / proportional scaling).
// Mengembalikan relative path gambar (contoh: "products/uuid.webp").
func CompressAndSaveImageWebP(file *multipart.FileHeader, directory string, cfg *config.Config) (string, error) {
	// 1. Validasi tipe file gambar
	if !IsValidImageType(file) {
		return "", errors.New("tipe file tidak didukung. Hanya JPG, PNG, dan WebP yang diperbolehkan")
	}

	// 2. Validasi batas ukuran input (maks 10MB)
	if file.Size > MaxInputUploadSize {
		return "", errors.New("ukuran file gambar maksimal 10MB")
	}

	// 3. Pastikan tool cwebp tersedia
	cwebpPath, err := exec.LookPath("cwebp")
	if err != nil {
		return "", fmt.Errorf("tool cwebp tidak ditemukan. Pastikan libwebp-tools terinstall: %w", err)
	}

	// 4. Siapkan direktori penyimpanan
	uploadPath := filepath.Join(cfg.UploadPath, directory)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori upload: %w", err)
	}

	// 5. Simpan file upload ke temporary file
	srcFile, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membaca file upload: %w", err)
	}
	defer srcFile.Close()

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".tmp"
	}
	tempPattern := fmt.Sprintf("upload-img-*%s", ext)
	tempFile, err := os.CreateTemp("", tempPattern)
	if err != nil {
		return "", fmt.Errorf("gagal membuat temporary file: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath) // Bersihkan temp file saat selesai

	if _, err := io.Copy(tempFile, srcFile); err != nil {
		tempFile.Close()
		return "", fmt.Errorf("gagal menyalin ke temporary file: %w", err)
	}
	tempFile.Close()

	// 6. Deteksi dimensi asli gambar
	origWidth, origHeight, _ := getImageDimensions(tempPath)

	// 7. Siapkan path tujuan (.webp)
	filename := fmt.Sprintf("%s.webp", uuid.New().String())
	outputPath := filepath.Join(uploadPath, filename)

	// 8. Coba kompresi bertahap sampai ukuran output memenuhi batas.
	if err := compressToTargetWebP(cwebpPath, tempPath, outputPath, origWidth, origHeight); err != nil {
		return "", err
	}

	// 11. Format relative path untuk database / URL
	relativePath := filepath.ToSlash(filepath.Join(directory, filename))
	return relativePath, nil
}

// CompressExistingFileWebP mengonversi dan mengompresi file gambar yang sudah ada di disk ke file WebP baru (< 200KB).
// File sumber asli tidak diubah atau dihapus.
func CompressExistingFileWebP(srcFullPath, dstFullPath string) error {
	cwebpPath, err := exec.LookPath("cwebp")
	if err != nil {
		return fmt.Errorf("tool cwebp tidak ditemukan. Pastikan libwebp-tools terinstall: %w", err)
	}

	// Pastikan direktori tujuan tersedia
	if err := os.MkdirAll(filepath.Dir(dstFullPath), 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori tujuan: %w", err)
	}

	// Deteksi dimensi asli
	origWidth, origHeight, _ := getImageDimensions(srcFullPath)

	return compressToTargetWebP(cwebpPath, srcFullPath, dstFullPath, origWidth, origHeight)
}

// compressToTargetWebP menurunkan resolusi dan kualitas secara bertahap sampai
// file hasil benar-benar berada di bawah batas ukuran.
func compressToTargetWebP(cwebpPath, inputPath, outputPath string, width, height int) error {
	passes := [...]struct {
		maxDimension int
		quality      int
	}{
		{MaxImageResolution, 82},
		{SecondPassResolution, 75},
		{1280, 68},
		{1024, 58},
		{800, 48},
		{640, 38},
	}

	for index, pass := range passes {
		if index > 0 {
			if err := os.Remove(outputPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("gagal menyiapkan kompresi ulang gambar: %w", err)
			}
		}

		if err := executeCwebp(cwebpPath, inputPath, outputPath, width, height, pass.maxDimension, pass.quality); err != nil {
			_ = os.Remove(outputPath)
			return fmt.Errorf("gagal mengonversi gambar ke WebP: %w", err)
		}

		outputInfo, err := os.Stat(outputPath)
		if err != nil {
			_ = os.Remove(outputPath)
			return fmt.Errorf("gagal memeriksa file hasil kompresi: %w", err)
		}
		if outputInfo.Size() <= TargetMaxImageSize {
			return nil
		}
	}

	_ = os.Remove(outputPath)
	return fmt.Errorf("gambar tidak dapat dikompresi hingga di bawah %dKB", TargetMaxImageSize/1024)
}

// getImageDimensions membaca lebar dan tinggi gambar tanpa memuat seluruh piksel ke memori
func getImageDimensions(filePath string) (int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}

// executeCwebp menjalankan perintah cwebp dengan parameter kompresi dan kalkulasi proporsional anti-gepeng
func executeCwebp(cwebpPath, inputPath, outputPath string, width, height, maxDimension, quality int) error {
	args := []string{
		"-q", fmt.Sprintf("%d", quality),
		"-size", fmt.Sprintf("%d", TargetMaxImageSize), // Target kompresi < 200KB
		"-m", "6", // Best compression algorithm
		"-sharp_yuv",        // Menjaga ketajaman warna RGB->YUV
		"-metadata", "none", // Hapus EXIF/metadata agar ukuran efisien
	}

	// Kalkulasi resize proporsional (salah satu sisi diset 0 agar cwebp mengunci aspek rasio asli)
	if width > 0 && height > 0 {
		if width > maxDimension || height > maxDimension {
			if width >= height {
				// Landscape / Square: clamp lebar, tinggi dihitung proporsional oleh cwebp
				args = append(args, "-resize", fmt.Sprintf("%d", maxDimension), "0")
			} else {
				// Portrait: clamp tinggi, lebar dihitung proporsional oleh cwebp
				args = append(args, "-resize", "0", fmt.Sprintf("%d", maxDimension))
			}
		}
	}

	args = append(args, inputPath, "-o", outputPath)

	cmd := exec.Command(cwebpPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w (output: %s)", err, string(output))
	}

	return nil
}
