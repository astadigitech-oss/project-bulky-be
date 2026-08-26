package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"project-bulky-be/internal/config"

	"github.com/google/uuid"
)

// GenerateThumbnailFromVideo extracts a thumbnail image from a video file using ffmpeg
// Returns the relative path for URL generation (e.g., "/uploads/video/thumbnail/uuid.jpg")
func GenerateThumbnailFromVideo(videoPath string, directory string, cfg *config.Config) (string, error) {
	// Check if ffmpeg is available
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("ffmpeg tidak ditemukan. Install ffmpeg terlebih dahulu untuk auto-generate thumbnail")
	}

	// Create thumbnail directory if not exists
	thumbnailDir := filepath.Join(cfg.UploadPath, directory)
	if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori thumbnail: %w", err)
	}

	// Temporary jpg filename for frame extraction
	tmpFilename := fmt.Sprintf("%s_tmp.jpg", uuid.New().String())
	tmpPath := filepath.Join(thumbnailDir, tmpFilename)
	defer os.Remove(tmpPath)

	// Final WebP filename
	finalFilename := fmt.Sprintf("%s.webp", uuid.New().String())
	finalPath := filepath.Join(thumbnailDir, finalFilename)

	// Full path to video file
	fullVideoPath := filepath.Join(cfg.UploadPath, filepath.FromSlash(videoPath))

	// Extract frame at 1 second using ffmpeg
	cmd := exec.Command(
		"ffmpeg",
		"-i", fullVideoPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-q:v", "2",
		"-y",
		tmpPath,
	)

	// Run command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("gagal generate thumbnail dari video: %w (output: %s)", err, string(output))
	}

	// Verify temporary thumbnail was created
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		return "", fmt.Errorf("thumbnail temporary file tidak terbuat")
	}

	// Compress extracted frame to WebP (< 200KB)
	if err := CompressExistingFileWebP(tmpPath, finalPath); err != nil {
		return "", fmt.Errorf("gagal kompresi thumbnail video ke webp: %w", err)
	}

	// Return relative path for URL (tanpa prefix /uploads/, konsisten dengan SaveUploadedFile)
	relativePath := filepath.ToSlash(filepath.Join(directory, finalFilename))

	return relativePath, nil
}
