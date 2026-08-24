package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
)

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{5242880, "5.00 MB"},
		{1073741824, "1.00 GB"},
	}

	for _, tt := range tests {
		actual := models.FormatFileSize(tt.bytes)
		if actual != tt.expected {
			t.Errorf("FormatFileSize(%d) = %s, expected %s", tt.bytes, actual, tt.expected)
		}
	}
}

func TestBackupFilenameRegex(t *testing.T) {
	validNames := []string{
		"backup_bulky_2026-08-24_170000_auto.sql.gz",
		"backup_bulky_2026-08-24_170000_manual.sql.gz",
		"backup_test_123.sql.gz",
		"backup_custom-name.sql.gz",
	}

	for _, name := range validNames {
		if !backupFilenameRegex.MatchString(name) {
			t.Errorf("Expected filename %s to be valid, but regex rejected it", name)
		}
	}

	invalidNames := []string{
		"../backup_123.sql.gz",
		"backup/test.sql.gz",
		"backup_test.sql",
		"backup_test.tar.gz",
		"script.sh",
		"backup_$(rm -rf).sql.gz",
		"backup_with space.sql.gz",
	}

	for _, name := range invalidNames {
		if backupFilenameRegex.MatchString(name) {
			t.Errorf("Expected filename %s to be invalid, but regex accepted it", name)
		}
	}
}

func TestGetBackupFilePath_PathTraversal(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "bulky-backup-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := &config.Config{
		BackupPath: tempDir,
	}

	svc := NewBackupService(cfg, nil)

	// Test invalid filenames
	_, err = svc.GetBackupFilePath("../etc/passwd")
	if err != ErrInvalidFilename {
		t.Errorf("Expected ErrInvalidFilename, got %v", err)
	}

	_, err = svc.GetBackupFilePath("backup_non_existent.sql.gz")
	if err != ErrBackupNotFound {
		t.Errorf("Expected ErrBackupNotFound, got %v", err)
	}

	// Create a legitimate test backup file
	validFile := "backup_bulky_2026-08-24_120000_auto.sql.gz"
	fullPath := filepath.Join(tempDir, validFile)
	if err := os.WriteFile(fullPath, []byte("test-content"), 0644); err != nil {
		t.Fatalf("Failed to write test backup file: %v", err)
	}

	resolvedPath, err := svc.GetBackupFilePath(validFile)
	if err != nil {
		t.Errorf("Expected successful resolution, got error: %v", err)
	}
	if resolvedPath != fullPath {
		t.Errorf("Resolved path mismatch: got %s, expected %s", resolvedPath, fullPath)
	}
}

func TestListAndPruneBackups(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "bulky-backup-list-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := &config.Config{
		BackupPath:          tempDir,
		BackupRetentionDays: 7,
	}

	svc := NewBackupService(cfg, nil)

	// Create 3 files: 1 recent auto, 1 recent manual, 1 old (expired)
	recentAuto := "backup_bulky_2026-08-24_120000_auto.sql.gz"
	recentManual := "backup_bulky_2026-08-24_140000_manual.sql.gz"
	oldFile := "backup_bulky_2026-08-01_100000_auto.sql.gz"

	os.WriteFile(filepath.Join(tempDir, recentAuto), []byte("auto-content"), 0644)
	os.WriteFile(filepath.Join(tempDir, recentManual), []byte("manual-content"), 0644)
	os.WriteFile(filepath.Join(tempDir, oldFile), []byte("old-content"), 0644)

	// Manually set old timestamp on oldFile (10 days ago)
	tenDaysAgo := time.Now().UTC().AddDate(0, 0, -10)
	os.Chtimes(filepath.Join(tempDir, oldFile), tenDaysAgo, tenDaysAgo)

	// Test ListBackups
	res, err := svc.ListBackups(context.Background())
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}

	if len(res.Items) != 3 {
		t.Errorf("Expected 3 backup items, got %d", len(res.Items))
	}
	if res.Stats.TotalFiles != 3 {
		t.Errorf("Expected stats.TotalFiles = 3, got %d", res.Stats.TotalFiles)
	}

	// Test PruneOldBackups (should delete oldFile)
	pruned, err := svc.PruneOldBackups(context.Background())
	if err != nil {
		t.Fatalf("PruneOldBackups failed: %v", err)
	}
	if pruned != 1 {
		t.Errorf("Expected 1 file pruned, got %d", pruned)
	}

	// Verify remaining files
	resAfter, err := svc.ListBackups(context.Background())
	if err != nil {
		t.Fatalf("ListBackups after prune failed: %v", err)
	}
	if len(resAfter.Items) != 2 {
		t.Errorf("Expected 2 backup items remaining, got %d", len(resAfter.Items))
	}
}

func TestCalculateNextRunTime(t *testing.T) {
	cfg := &config.Config{
		BackupScheduleHour:   17,
		BackupScheduleMinute: 0,
	}

	svc := &backupService{cfg: cfg}

	// Case 1: Current time is 10:00 UTC (before 17:00 UTC) -> should run today at 17:00 UTC
	nowBefore := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)
	nextBefore := svc.calculateNextRunTime(nowBefore)
	expectedBefore := time.Date(2026, 8, 24, 17, 0, 0, 0, time.UTC)

	if !nextBefore.Equal(expectedBefore) {
		t.Errorf("Expected next run %v, got %v", expectedBefore, nextBefore)
	}

	// Case 2: Current time is 18:00 UTC (after 17:00 UTC) -> should run tomorrow at 17:00 UTC
	nowAfter := time.Date(2026, 8, 24, 18, 0, 0, 0, time.UTC)
	nextAfter := svc.calculateNextRunTime(nowAfter)
	expectedAfter := time.Date(2026, 8, 25, 17, 0, 0, 0, time.UTC)

	if !nextAfter.Equal(expectedAfter) {
		t.Errorf("Expected next run %v, got %v", expectedAfter, nextAfter)
	}
}
