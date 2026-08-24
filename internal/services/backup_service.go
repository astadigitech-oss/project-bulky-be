package services

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

var (
	ErrBackupInProgress = errors.New("proses backup database sedang berjalan, silakan coba beberapa saat lagi")
	ErrInvalidFilename  = errors.New("nama file backup tidak valid atau mengandung karakter terlarang")
	ErrBackupNotFound   = errors.New("file backup tidak ditemukan")
)

// Regex untuk validasi nama file backup agar aman dari path traversal
var backupFilenameRegex = regexp.MustCompile(`^backup_[a-zA-Z0-9_\-\.]+\.sql\.gz$`)

type BackupService interface {
	// CreateBackup mengeksekusi pg_dump + gzip untuk membuat file dump terkompresi.
	CreateBackup(ctx context.Context, triggerType models.BackupTriggerType, adminID *uuid.UUID, adminEmail string) (*models.BackupItem, error)
	// ListBackups mengambil daftar semua file backup di storage beserta statistik storage.
	ListBackups(ctx context.Context) (*models.BackupListResponse, error)
	// GetBackupFilePath mengembalikan path file backup yang sudah divalidasi keamanannya.
	GetBackupFilePath(filename string) (string, error)
	// DeleteBackup menghapus file backup tertentu secara permanen.
	DeleteBackup(ctx context.Context, filename string, adminID *uuid.UUID, adminEmail string) error
	// PruneOldBackups menghapus file backup yang telah melewati batas hari retensi.
	PruneOldBackups(ctx context.Context) (int, error)
	// StartDailyScheduler menjalankan cron scheduler backup harian di background.
	StartDailyScheduler(ctx context.Context)
	// GetNextScheduledRun mengembalikan waktu eksekusi backup otomatis berikutnya.
	GetNextScheduledRun() time.Time
}

type backupService struct {
	cfg             *config.Config
	activityLogRepo repositories.ActivityLogRepository
	isBackingUp     atomic.Bool
	mu              sync.Mutex
	nextRun         time.Time
	nextRunMu       sync.RWMutex
}

func NewBackupService(cfg *config.Config, activityLogRepo repositories.ActivityLogRepository) BackupService {
	return &backupService{
		cfg:             cfg,
		activityLogRepo: activityLogRepo,
	}
}

// ensureBackupDir memastikan direktori backup tersedia
func (s *backupService) ensureBackupDir() error {
	return os.MkdirAll(s.cfg.BackupPath, 0755)
}

// CreateBackup mengeksekusi pg_dump secara streaming ke file .sql.gz
func (s *backupService) CreateBackup(
	ctx context.Context,
	triggerType models.BackupTriggerType,
	adminID *uuid.UUID,
	adminEmail string,
) (*models.BackupItem, error) {
	// Mutex & atomic check untuk mencegah backup ganda bersamaan
	if !s.isBackingUp.CompareAndSwap(false, true) {
		return nil, ErrBackupInProgress
	}
	defer s.isBackingUp.Store(false)

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.ensureBackupDir(); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori backup: %w", err)
	}

	now := time.Now().UTC()
	suffix := "auto"
	if triggerType == models.BackupTriggerManual {
		suffix = "manual"
	}

	filename := fmt.Sprintf("backup_bulky_%s_%s.sql.gz", now.Format("2006-01-02_150405"), suffix)
	fullPath := filepath.Join(s.cfg.BackupPath, filename)

	// Persiapkan file tujuan
	dstFile, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat file backup di disk: %w", err)
	}

	// Gzip writer dengan kompresi standar (hemat CPU & ukuran optimal)
	gzipWriter, err := gzip.NewWriterLevel(dstFile, gzip.DefaultCompression)
	if err != nil {
		dstFile.Close()
		os.Remove(fullPath)
		return nil, fmt.Errorf("gagal menginisialisasi gzip: %w", err)
	}

	// Siapkan perintah pg_dump
	// Gunakan argumen non-blocking dan bersih
	args := []string{
		"-h", s.cfg.DBHost,
		"-p", s.cfg.DBPort,
		"-U", s.cfg.DBUser,
		"-d", s.cfg.DBName,
		"--no-owner",
		"--no-acl",
		"--clean",
		"--if-exists",
	}

	cmd := exec.CommandContext(ctx, "pg_dump", args...)
	// Password dikirim aman via environment variable PGPASSWORD
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", s.cfg.DBPassword))

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		gzipWriter.Close()
		dstFile.Close()
		os.Remove(fullPath)
		return nil, fmt.Errorf("gagal menghubungkan stdout pg_dump: %w", err)
	}

	var stderrBuf strings.Builder
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		gzipWriter.Close()
		dstFile.Close()
		os.Remove(fullPath)
		if strings.Contains(err.Error(), "executable file not found") {
			return nil, fmt.Errorf("tool pg_dump tidak ditemukan di sistem host. Di Mac, jalankan 'brew install libpq && brew link --force libpq' (di Docker server sudah terinstall otomatis)")
		}
		return nil, fmt.Errorf("gagal menjalankan pg_dump: %w (stderr: %s)", err, stderrBuf.String())
	}

	// Stream dari stdout pg_dump langsung ke gzip writer
	written, copyErr := io.Copy(gzipWriter, stdoutPipe)
	if copyErr != nil {
		gzipWriter.Close()
		dstFile.Close()
		os.Remove(fullPath)
		return nil, fmt.Errorf("gagal menulis stream kompresi backup: %w", copyErr)
	}

	// Tutup gzip dan file sebelum menunggu process exit
	if err := gzipWriter.Close(); err != nil {
		dstFile.Close()
		os.Remove(fullPath)
		return nil, fmt.Errorf("gagal menyelesaikan kompresi gzip: %w", err)
	}

	if err := dstFile.Sync(); err != nil {
		dstFile.Close()
		os.Remove(fullPath)
		return nil, fmt.Errorf("gagal flush file backup ke disk: %w", err)
	}

	if err := dstFile.Close(); err != nil {
		os.Remove(fullPath)
		return nil, fmt.Errorf("gagal menutup file backup: %w", err)
	}

	// Tunggu proses pg_dump selesai
	if err := cmd.Wait(); err != nil {
		os.Remove(fullPath)
		return nil, fmt.Errorf("pg_dump gagal dieksekusi: %w (stderr: %s)", err, stderrBuf.String())
	}

	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca metadata file backup: %w", err)
	}

	item := &models.BackupItem{
		Filename:      filename,
		SizeBytes:     fileInfo.Size(),
		SizeFormatted: models.FormatFileSize(fileInfo.Size()),
		CreatedAt:     fileInfo.ModTime().UTC(),
		TriggerType:   triggerType,
		DownloadURL:   fmt.Sprintf("/api/panel/backups/%s/download", filename),
	}

	log.Printf("[backup-service] Berhasil membuat backup database %s (Ukuran: %s, RAW Stream: %s)",
		filename, item.SizeFormatted, models.FormatFileSize(written))

	// Catat audit trail ke activity log
	s.logActivity(triggerType, models.ActionCreate, filename, item.SizeFormatted, adminID, adminEmail)

	return item, nil
}

// ListBackups memindai direktori backup dan mengembalikan daftar file beserta statistik
func (s *backupService) ListBackups(ctx context.Context) (*models.BackupListResponse, error) {
	if err := s.ensureBackupDir(); err != nil {
		return nil, fmt.Errorf("gagal membaca direktori backup: %w", err)
	}

	entries, err := os.ReadDir(s.cfg.BackupPath)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca isi direktori backup: %w", err)
	}

	var items []models.BackupItem
	var totalSizeBytes int64

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !backupFilenameRegex.MatchString(filename) {
			continue // Lewati file di luar format backup
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		triggerType := models.BackupTriggerAuto
		if strings.Contains(filename, "_manual.sql.gz") {
			triggerType = models.BackupTriggerManual
		}

		totalSizeBytes += info.Size()

		items = append(items, models.BackupItem{
			Filename:      filename,
			SizeBytes:     info.Size(),
			SizeFormatted: models.FormatFileSize(info.Size()),
			CreatedAt:     info.ModTime().UTC(),
			TriggerType:   triggerType,
			DownloadURL:   fmt.Sprintf("/api/panel/backups/%s/download", filename),
		})
	}

	// Urutkan dari yang paling baru ke paling lama
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	nextRun := s.GetNextScheduledRun()
	var nextRunPtr *time.Time
	if !nextRun.IsZero() {
		nextRunPtr = &nextRun
	}

	stats := models.BackupStorageStats{
		TotalFiles:         len(items),
		TotalSizeBytes:     totalSizeBytes,
		TotalSizeFormatted: models.FormatFileSize(totalSizeBytes),
		RetentionDays:      s.cfg.BackupRetentionDays,
		BackupPath:         s.cfg.BackupPath,
		NextScheduledRun:   nextRunPtr,
	}

	return &models.BackupListResponse{
		Items: items,
		Stats: stats,
	}, nil
}

// GetBackupFilePath memvalidasi nama file dari path traversal dan memastikan file ada
func (s *backupService) GetBackupFilePath(filename string) (string, error) {
	if !backupFilenameRegex.MatchString(filename) {
		return "", ErrInvalidFilename
	}

	cleanPath := filepath.Clean(filepath.Join(s.cfg.BackupPath, filename))
	cleanBackupDir := filepath.Clean(s.cfg.BackupPath)

	// Path Traversal Defense
	if filepath.Dir(cleanPath) != cleanBackupDir {
		return "", ErrInvalidFilename
	}

	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return "", ErrBackupNotFound
	}

	return cleanPath, nil
}

// DeleteBackup menghapus file backup tertentu secara permanen
func (s *backupService) DeleteBackup(
	ctx context.Context,
	filename string,
	adminID *uuid.UUID,
	adminEmail string,
) error {
	filePath, err := s.GetBackupFilePath(filename)
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("gagal menghapus file backup: %w", err)
	}

	log.Printf("[backup-service] File backup %s berhasil dihapus oleh %s", filename, adminEmail)

	// Catat audit trail
	s.logActivity(models.BackupTriggerManual, models.ActionDelete, filename, "", adminID, adminEmail)

	return nil
}

// PruneOldBackups menghapus file backup yang lebih tua dari batas retensi
func (s *backupService) PruneOldBackups(ctx context.Context) (int, error) {
	if s.cfg.BackupRetentionDays <= 0 {
		return 0, nil
	}

	if err := s.ensureBackupDir(); err != nil {
		return 0, err
	}

	entries, err := os.ReadDir(s.cfg.BackupPath)
	if err != nil {
		return 0, err
	}

	cutoff := time.Now().UTC().AddDate(0, 0, -s.cfg.BackupRetentionDays)
	prunedCount := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !backupFilenameRegex.MatchString(filename) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().UTC().Before(cutoff) {
			fullPath := filepath.Join(s.cfg.BackupPath, filename)
			if err := os.Remove(fullPath); err == nil {
				prunedCount++
				log.Printf("[backup-service] Auto-prune menghapus backup kadaluarsa: %s (Dibuat: %s)",
					filename, info.ModTime().UTC().Format(time.RFC3339))
			}
		}
	}

	if prunedCount > 0 {
		log.Printf("[backup-service] Auto-prune selesai: %d file backup kadaluarsa berhasil dibersihkan", prunedCount)
	}

	return prunedCount, nil
}

// StartDailyScheduler menjalankan timer harian di background
func (s *backupService) StartDailyScheduler(ctx context.Context) {
	for {
		nextRun := s.calculateNextRunTime(time.Now().UTC())
		s.setNextScheduledRun(nextRun)

		waitDuration := time.Until(nextRun)
		log.Printf("[backup-scheduler] Backup otomatis berikutnya dijadwalkan pada %s UTC (dalam %v)",
			nextRun.Format("2006-01-02 15:04:05"), waitDuration.Round(time.Minute))

		timer := time.NewTimer(waitDuration)

		select {
		case <-ctx.Done():
			timer.Stop()
			log.Println("[backup-scheduler] Scheduler backup database dihentikan")
			return
		case <-timer.C:
			log.Println("[backup-scheduler] Menjalankan auto backup database harian...")
			if _, err := s.CreateBackup(ctx, models.BackupTriggerAuto, nil, "system"); err != nil {
				log.Printf("[backup-scheduler] Gagal membuat auto backup: %v", err)
			}

			// Bersihkan file yang sudah melewati masa retensi
			if _, err := s.PruneOldBackups(ctx); err != nil {
				log.Printf("[backup-scheduler] Gagal menjalankan auto-prune: %v", err)
			}
		}
	}
}

// calculateNextRunTime menghitung target waktu berikutnya berdasarkan jam/menit UTC yang dikonfigurasi
func (s *backupService) calculateNextRunTime(now time.Time) time.Time {
	targetHour := s.cfg.BackupScheduleHour
	targetMinute := s.cfg.BackupScheduleMinute

	target := time.Date(
		now.Year(), now.Month(), now.Day(),
		targetHour, targetMinute, 0, 0,
		time.UTC,
	)

	// Jika target waktu hari ini sudah lewat, jadwalkan untuk besok
	if !target.After(now) {
		target = target.Add(24 * time.Hour)
	}

	return target
}

func (s *backupService) setNextScheduledRun(t time.Time) {
	s.nextRunMu.Lock()
	defer s.nextRunMu.Unlock()
	s.nextRun = t
}

func (s *backupService) GetNextScheduledRun() time.Time {
	s.nextRunMu.RLock()
	defer s.nextRunMu.RUnlock()
	return s.nextRun
}

// logActivity mencatat aktivitas ke tabel audit activity_log
func (s *backupService) logActivity(
	triggerType models.BackupTriggerType,
	action models.ActivityAction,
	filename string,
	sizeFormatted string,
	adminID *uuid.UUID,
	adminEmail string,
) {
	if s.activityLogRepo == nil {
		return
	}

	userType := "ADMIN"
	if triggerType == models.BackupTriggerAuto || adminID == nil {
		userType = "SYSTEM"
	}

	var deskripsi string
	if action == models.ActionCreate {
		if userType == "SYSTEM" {
			deskripsi = fmt.Sprintf("Auto-backup database berhasil dibuat: %s (Ukuran: %s)", filename, sizeFormatted)
		} else {
			deskripsi = fmt.Sprintf("Manual backup database dipicu oleh %s: %s (Ukuran: %s)", adminEmail, filename, sizeFormatted)
		}
	} else if action == models.ActionDelete {
		deskripsi = fmt.Sprintf("File backup database dihapus oleh %s: %s", adminEmail, filename)
	}

	entityType := "backup"
	newData, _ := json.Marshal(map[string]interface{}{
		"filename":       filename,
		"size_formatted": sizeFormatted,
		"trigger_type":   triggerType,
	})

	entry := &models.ActivityLog{
		UserType:   userType,
		UserID:     adminID,
		Action:     action,
		Modul:      "backup",
		EntityType: &entityType,
		Deskripsi:  deskripsi,
		NewData:    newData,
	}

	if err := s.activityLogRepo.Create(entry); err != nil {
		log.Printf("[backup-service] Gagal mencatat activity log backup: %v", err)
	}
}
