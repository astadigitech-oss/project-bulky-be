package transcoder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type TranscodeResult struct {
	// OutputPath adalah absolute path ke file hasil transcode.
	OutputPath  string
	DurasiDetik int
}

// Transcode mengkonversi video input ke MP4 streamable (faststart, H.264, height 1280).
// Cocok untuk konten vertikal 9:16 Bulky TV.
// inputPath harus berupa absolute path.
// Output disimpan di direktori yang sama dengan prefix "stream_".
func Transcode(inputPath string) (*TranscodeResult, error) {
	dir := filepath.Dir(inputPath)
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	output := filepath.Join(dir, "stream_"+base+".mp4")

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-threads", "2",
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "28",
		"-vf", "scale=-2:1280",
		"-c:a", "aac",
		"-b:a", "96k",
		"-movflags", "+faststart",
		"-f", "mp4",
		"-y",
		output,
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg gagal: %w\noutput: %s", err, string(out))
	}

	durasi, err := extractDuration(output)
	if err != nil {
		durasi = 0 // fallback, tidak fatal
	}

	return &TranscodeResult{OutputPath: output, DurasiDetik: durasi}, nil
}

// extractDuration mengambil durasi video dalam detik via ffprobe.
func extractDuration(path string) (int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return 0, err
	}
	return int(f), nil
}

// Cleanup menghapus file di path yang diberikan (biasanya file raw upload sementara).
func Cleanup(path string) {
	_ = os.Remove(path)
}
