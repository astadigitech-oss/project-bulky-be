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

	// Generate unique filename for thumbnail
	thumbnailFilename := fmt.Sprintf("%s.jpg", uuid.New().String())
	thumbnailPath := filepath.Join(thumbnailDir, thumbnailFilename)

	// Full path to video file
	fullVideoPath := filepath.Join(cfg.UploadPath, filepath.FromSlash(videoPath))

	// Extract frame at 1 second using ffmpeg
	// -i: input file
	// -ss: seek to timestamp (1 second)
	// -vframes: number of frames to extract (1)
	// -q:v: quality (2 is high quality, range 2-31)
	// -y: overwrite output file if exists
	cmd := exec.Command(
		"ffmpeg",
		"-i", fullVideoPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-q:v", "2",
		"-y",
		thumbnailPath,
	)

	// Run command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("gagal generate thumbnail dari video: %w (output: %s)", err, string(output))
	}

	// Verify thumbnail was created
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		return "", fmt.Errorf("thumbnail file tidak terbuat")
	}

	// Return relative path for URL
	relativePath := filepath.Join("/uploads", directory, thumbnailFilename)
	// Convert backslashes to forward slashes for URL
	relativePath = filepath.ToSlash(relativePath)

	return relativePath, nil
}
