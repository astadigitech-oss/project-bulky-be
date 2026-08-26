package models

import (
	"fmt"
	"time"
)

// BackupTriggerType mendefinisikan sumber pemicu pembuatan backup
type BackupTriggerType string

const (
	BackupTriggerAuto   BackupTriggerType = "AUTO"
	BackupTriggerManual BackupTriggerType = "MANUAL"
)

// BackupItem merepresentasikan informasi satu file backup database
type BackupItem struct {
	Filename      string            `json:"filename"`
	SizeBytes     int64             `json:"size_bytes"`
	SizeFormatted string            `json:"size_formatted"`
	CreatedAt     time.Time         `json:"created_at"`
	TriggerType   BackupTriggerType `json:"trigger_type"`
	DownloadURL   string            `json:"download_url"`
}

// BackupStorageStats merangkum statistik penyimpanan backup database
type BackupStorageStats struct {
	TotalFiles         int        `json:"total_files"`
	TotalSizeBytes     int64      `json:"total_size_bytes"`
	TotalSizeFormatted string     `json:"total_size_formatted"`
	RetentionDays      int        `json:"retention_days"`
	BackupPath         string     `json:"backup_path"`
	NextScheduledRun   *time.Time `json:"next_scheduled_run,omitempty"`
}

// BackupListResponse merangkum respon endpoint GET /api/panel/backups
type BackupListResponse struct {
	Items []BackupItem       `json:"items"`
	Stats BackupStorageStats `json:"stats"`
}

// FormatFileSize mengonversi byte ke representasi yang mudah dibaca (KB, MB, GB)
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
