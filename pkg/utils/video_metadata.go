package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/abema/go-mp4"
)

// VideoMetadata represents video file metadata
type VideoMetadata struct {
	Duration float64 `json:"duration"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
}

// ExtractVideoMetadata extracts metadata from a video file using pure Go parser
// Supports MP4, MOV, M4V formats (MP4-based container formats)
// Does NOT support: AVI, MKV, WebM, FLV, MPEG
// Returns duration in seconds, width and height in pixels
func ExtractVideoMetadata(filePath string) (*VideoMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open video file: %w", err)
	}
	defer file.Close()

	var duration float64
	var width, height int

	// Parse MP4 structure
	_, err = mp4.ReadBoxStructure(file, func(h *mp4.ReadHandle) (interface{}, error) {
		// Read movie header (mvhd) to get duration and timescale
		if h.BoxInfo.Type == mp4.BoxTypeMvhd() {
			box, _, err := h.ReadPayload()
			if err != nil {
				return nil, err
			}
			mvhd := box.(*mp4.Mvhd)
			if mvhd.Timescale > 0 {
				duration = float64(mvhd.DurationV1) / float64(mvhd.Timescale)
			}
		}

		// Read track header (tkhd) to get video dimensions
		if h.BoxInfo.Type == mp4.BoxTypeTkhd() {
			box, _, err := h.ReadPayload()
			if err != nil {
				return nil, err
			}
			tkhd := box.(*mp4.Tkhd)
			width = int(tkhd.Width >> 16)   // Fixed-point 16.16 format
			height = int(tkhd.Height >> 16) // Fixed-point 16.16 format
		}

		// Expand boxes to continue parsing
		return h.Expand()
	})

	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to parse video metadata: %w", err)
	}

	if duration == 0 {
		return nil, fmt.Errorf("could not extract duration from video file")
	}

	return &VideoMetadata{
		Duration: duration,
		Width:    width,
		Height:   height,
	}, nil
}

// GetVideoDurationInSeconds returns video duration in seconds (rounded)
func GetVideoDurationInSeconds(filePath string) (int, error) {
	// Try go-mp4 parser first (pure Go, faster)
	metadata, err := ExtractVideoMetadata(filePath)
	if err == nil {
		return int(metadata.Duration), nil
	}

	// Fallback to ffprobe if go-mp4 fails
	return getVideoDurationWithFFProbe(filePath)
}

// getVideoDurationWithFFProbe uses ffprobe to get video duration
func getVideoDurationWithFFProbe(filePath string) (int, error) {
	// Check if ffprobe is available
	_, err := exec.LookPath("ffprobe")
	if err != nil {
		return 0, fmt.Errorf("ffprobe tidak ditemukan dan go-mp4 parser gagal. Install ffmpeg untuk fallback detection")
	}

	// Run ffprobe to get duration
	// ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 file.mp4
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe gagal: %w", err)
	}

	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("gagal parse duration dari ffprobe: %w", err)
	}

	return int(duration), nil
}
