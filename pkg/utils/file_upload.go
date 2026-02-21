package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"project-bulky-be/internal/config"

	"github.com/google/uuid"
)

var allowedImageTypes = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
	"image/svg+xml",
}

var allowedDocumentTypes = []string{
	"application/pdf",
}

// allowedVideoTypes - Limited to MP4-based formats supported by go-mp4 library
// Supported: MP4, MOV, M4V (for auto-duration detection)
var allowedVideoTypes = []string{
	"video/mp4",
	"video/quicktime", // MOV
}

var allowedVideoExtensions = []string{
	".mp4",
	".mov",
	".m4v",
}

// IsValidImageType validates if the uploaded file is a valid image type
func IsValidImageType(file *multipart.FileHeader) bool {
	contentType := file.Header.Get("Content-Type")
	for _, allowed := range allowedImageTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
}

// IsValidDocumentType validates if the uploaded file is a valid document type
func IsValidDocumentType(file *multipart.FileHeader) bool {
	contentType := file.Header.Get("Content-Type")
	for _, allowed := range allowedDocumentTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
}

// IsValidVideoType validates if the uploaded file is a valid video type
// Only supports MP4-based formats (MP4, MOV, M4V) for auto-duration detection
func IsValidVideoType(file *multipart.FileHeader) bool {
	// Check Content-Type
	contentType := file.Header.Get("Content-Type")
	validContentType := false
	for _, allowed := range allowedVideoTypes {
		if contentType == allowed {
			validContentType = true
			break
		}
	}

	// Check file extension as additional validation
	filename := strings.ToLower(file.Filename)
	validExtension := false
	for _, ext := range allowedVideoExtensions {
		if strings.HasSuffix(filename, ext) {
			validExtension = true
			break
		}
	}

	// Both content type and extension must be valid
	return validContentType && validExtension
}

// SaveUploadedFile saves an uploaded file to the specified directory
// Returns the relative path for URL generation (e.g., "product-categories/uuid.png")
// Supports images, documents (PDF), and videos (MP4, MOV, M4V)
func SaveUploadedFile(file *multipart.FileHeader, directory string, cfg *config.Config) (string, error) {
	// Validate file type (image, document, or video)
	if !IsValidImageType(file) && !IsValidDocumentType(file) && !IsValidVideoType(file) {
		return "", errors.New("tipe file tidak didukung. Hanya jpg, png, webp, svg, pdf, dan video (mp4, mov, m4v) yang diperbolehkan")
	}

	// Create directory if not exists (use config upload path)
	uploadPath := filepath.Join(cfg.UploadPath, directory)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori upload: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	fullPath := filepath.Join(uploadPath, filename)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file: %w", err)
	}

	// Copy file
	if _, err = dst.ReadFrom(src); err != nil {
		dst.Close()
		return "", fmt.Errorf("gagal menyimpan file: %w", err)
	}

	// Ensure file is fully written to disk before returning
	if err := dst.Sync(); err != nil {
		dst.Close()
		return "", fmt.Errorf("gagal sync file: %w", err)
	}

	// Close file explicitly before returning path
	if err := dst.Close(); err != nil {
		return "", fmt.Errorf("gagal close file: %w", err)
	}

	// Return relative path for URL (without upload path prefix)
	relativePath := filepath.Join(directory, filename)
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	return relativePath, nil
}

// DeleteFile deletes a file from the filesystem
func DeleteFile(filePath string, cfg *config.Config) error {
	if filePath == "" {
		return nil
	}

	// Build full path (normalize path separator for OS)
	fullPath := filepath.Join(cfg.UploadPath, filepath.FromSlash(filePath))

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // File doesn't exist, no need to delete
	}

	// Delete file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("gagal menghapus file: %w", err)
	}

	return nil
}

// SaveUploadedFileWithCustomName saves an uploaded file with a custom filename
// Returns the relative path for URL generation (e.g., "product-categories/custom-name.png")
func SaveUploadedFileWithCustomName(file *multipart.FileHeader, directory, customName string, cfg *config.Config) (string, error) {
	// Validate image type
	if !IsValidImageType(file) {
		return "", errors.New("tipe file tidak didukung. Hanya jpg, png, webp, dan svg yang diperbolehkan")
	}

	// Create directory if not exists (use config upload path)
	uploadPath := filepath.Join(cfg.UploadPath, directory)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori upload: %w", err)
	}

	// Use custom name with original extension
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", customName, ext)
	fullPath := filepath.Join(uploadPath, filename)

	// Delete old file if exists
	relativePath := filepath.Join(directory, filename)
	if err := DeleteFile(relativePath, cfg); err != nil {
		// Log but don't fail if old file can't be deleted
		fmt.Printf("Warning: %v\n", err)
	}

	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file: %w", err)
	}

	// Copy file
	if _, err = dst.ReadFrom(src); err != nil {
		dst.Close()
		return "", fmt.Errorf("gagal menyimpan file: %w", err)
	}

	// Ensure file is fully written to disk before returning
	if err := dst.Sync(); err != nil {
		dst.Close()
		return "", fmt.Errorf("gagal sync file: %w", err)
	}

	// Close file explicitly before returning path
	if err := dst.Close(); err != nil {
		return "", fmt.Errorf("gagal close file: %w", err)
	}

	// Return relative path for URL (without upload path prefix)
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	return relativePath, nil
}

// GetFileURL returns the full URL for a file path
func GetFileURL(filePath interface{}, cfg *config.Config) string {
	base := strings.TrimRight(cfg.BaseURL, "/")
	base = strings.TrimSuffix(base, "/uploads")
	normalize := func(p string) string {
		if p == "" {
			return ""
		}
		if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
			return p
		}
		p = strings.TrimPrefix(p, "uploads/")
		return fmt.Sprintf("%s/uploads/%s", base, p)
	}
	switch v := filePath.(type) {
	case string:
		return normalize(v)
	case *string:
		if v == nil {
			return ""
		}
		return normalize(*v)
	default:
		return ""
	}
}

// GetFileURLPtr returns the full URL for a file path pointer
func GetFileURLPtr(filePath *string, cfg *config.Config) *string {
	if filePath == nil || *filePath == "" {
		return nil
	}
	p := *filePath
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return &p
	}
	p = strings.TrimPrefix(p, "uploads/")
	base := strings.TrimRight(cfg.BaseURL, "/")
	base = strings.TrimSuffix(base, "/uploads")
	url := fmt.Sprintf("%s/uploads/%s", base, p)
	return &url
}
