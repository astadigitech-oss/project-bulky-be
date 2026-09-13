package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"

	"project-bulky-be/internal/config"

	"github.com/google/uuid"
)

const (
	SeasonalWebLogoWidth      = 1368
	SeasonalWebLogoHeight     = 288
	SeasonalMobileLogoWidth   = 822
	SeasonalMobileLogoHeight  = 294
	SeasonalWebLogoMaxSize    = 150 * 1024
	SeasonalMobileLogoMaxSize = 100 * 1024

	SeasonalWebNavbarDecorationWidth   = 1920
	SeasonalWebNavbarDecorationHeight  = 160
	SeasonalWebNavbarDecorationMaxSize = 120 * 1024
	SeasonalMobileTopAppBarWidth       = 1080
	SeasonalMobileTopAppBarHeight      = 240
	SeasonalMobileTopAppBarMaxSize     = 100 * 1024
)

type SeasonalLogoVariants struct {
	WebLogoPath           string
	MobileLoadingLogoPath string
}

// SeasonalAssetTarget defines the exact canvas and maximum transfer size that
// Store FE is allowed to receive for a seasonal asset.
type SeasonalAssetTarget struct {
	Width       int
	Height      int
	MaxFileSize int64
	Filename    string
}

var (
	SeasonalWebNavbarDecorationTarget = SeasonalAssetTarget{
		Width: SeasonalWebNavbarDecorationWidth, Height: SeasonalWebNavbarDecorationHeight,
		MaxFileSize: SeasonalWebNavbarDecorationMaxSize, Filename: "web-navbar-decoration",
	}
	SeasonalMobileTopAppBarTarget = SeasonalAssetTarget{
		Width: SeasonalMobileTopAppBarWidth, Height: SeasonalMobileTopAppBarHeight,
		MaxFileSize: SeasonalMobileTopAppBarMaxSize, Filename: "mobile-top-app-bar-ornament",
	}
)

// SaveSeasonalLogoVariants creates platform canvases from one framed source.
// scale=contain preserves the logo's aspect ratio and transparent padding avoids distortion.
func SaveSeasonalLogoVariants(file *multipart.FileHeader, cfg *config.Config) (*SeasonalLogoVariants, error) {
	if !IsValidBannerImageType(file) {
		return nil, errors.New("tipe file tidak didukung")
	}
	if file.Size > MaxInputUploadSize {
		return nil, errors.New("ukuran file gambar maksimal 10MB")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, fmt.Errorf("tool ffmpeg tidak ditemukan: %w", err)
	}
	if _, err := exec.LookPath("cwebp"); err != nil {
		return nil, fmt.Errorf("tool cwebp tidak ditemukan: %w", err)
	}

	source, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer source.Close()
	temp, err := os.CreateTemp("", "seasonal-logo-*")
	if err != nil {
		return nil, err
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)
	if _, err := io.Copy(temp, source); err != nil {
		temp.Close()
		return nil, err
	}
	if err := temp.Close(); err != nil {
		return nil, err
	}

	baseDir := filepath.Join(cfg.UploadPath, "seasonal-campaign")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	webName, mobileName := fmt.Sprintf("web-logo-%s.webp", uuid.NewString()), fmt.Sprintf("mobile-loading-logo-%s.webp", uuid.NewString())
	webPath, mobilePath := filepath.Join(baseDir, webName), filepath.Join(baseDir, mobileName)
	if err := renderSeasonalLogo(tempPath, webPath, SeasonalWebLogoWidth, SeasonalWebLogoHeight, SeasonalWebLogoMaxSize); err != nil {
		return nil, err
	}
	if err := renderSeasonalLogo(tempPath, mobilePath, SeasonalMobileLogoWidth, SeasonalMobileLogoHeight, SeasonalMobileLogoMaxSize); err != nil {
		_ = os.Remove(webPath)
		return nil, err
	}
	return &SeasonalLogoVariants{WebLogoPath: filepath.ToSlash(filepath.Join("seasonal-campaign", webName)), MobileLoadingLogoPath: filepath.ToSlash(filepath.Join("seasonal-campaign", mobileName))}, nil
}

// SaveSeasonalAsset re-renders an uploaded ornament into its platform canvas.
// This keeps the Store FE payload predictable even when a caller bypasses the
// admin panel's browser-side framing step.
func SaveSeasonalAsset(file *multipart.FileHeader, cfg *config.Config, target SeasonalAssetTarget) (string, error) {
	if !IsValidBannerImageType(file) {
		return "", errors.New("tipe file tidak didukung")
	}
	if file.Size > MaxInputUploadSize {
		return "", errors.New("ukuran file gambar maksimal 10MB")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return "", fmt.Errorf("tool ffmpeg tidak ditemukan: %w", err)
	}
	if _, err := exec.LookPath("cwebp"); err != nil {
		return "", fmt.Errorf("tool cwebp tidak ditemukan: %w", err)
	}

	source, err := file.Open()
	if err != nil {
		return "", err
	}
	defer source.Close()
	temp, err := os.CreateTemp("", "seasonal-asset-*")
	if err != nil {
		return "", err
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)
	if _, err := io.Copy(temp, source); err != nil {
		temp.Close()
		return "", err
	}
	if err := temp.Close(); err != nil {
		return "", err
	}

	baseDir := filepath.Join(cfg.UploadPath, "seasonal-campaign")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%s-%s.webp", target.Filename, uuid.NewString())
	outputPath := filepath.Join(baseDir, name)
	if err := renderSeasonalLogo(tempPath, outputPath, target.Width, target.Height, target.MaxFileSize); err != nil {
		return "", err
	}
	return filepath.ToSlash(filepath.Join("seasonal-campaign", name)), nil
}

func renderSeasonalLogo(input, output string, width, height int, maxOutputSize int64) error {
	filter := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2:color=black@0", width, height, width, height)
	intermediate := output + ".png"
	defer os.Remove(intermediate)
	cmd := exec.Command("ffmpeg", "-y", "-i", input, "-frames:v", "1", "-vf", filter, "-c:v", "png", "-pix_fmt", "rgba", intermediate)
	if data, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("gagal membuat varian logo: %w (%s)", err, data)
	}
	for _, quality := range []int{82, 75, 68, 60} {
		_ = os.Remove(output)
		cmd = exec.Command("cwebp", "-q", fmt.Sprintf("%d", quality), "-size", fmt.Sprintf("%d", maxOutputSize), "-m", "6", "-sharp_yuv", "-metadata", "none", intermediate, "-o", output)
		if data, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("gagal mengonversi varian logo ke WebP: %w (%s)", err, data)
		}
		info, err := os.Stat(output)
		if err != nil {
			return fmt.Errorf("gagal memeriksa ukuran varian logo: %w", err)
		}
		if info.Size() <= maxOutputSize {
			return nil
		}
	}
	_ = os.Remove(output)
	return fmt.Errorf("varian logo tidak dapat dikompresi hingga di bawah %dKB", maxOutputSize/1024)
}
