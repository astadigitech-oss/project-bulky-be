package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
)

// ProdukAutoArchiveService menjalankan job berkala untuk mengarsipkan (is_active=false)
// produk yang sudah terjual (is_sold=true) lebih dari 1 hari, dihitung dari order-nya
// mencapai status SHIPPED atau COMPLETED.
type ProdukAutoArchiveService interface {
	// Run mengeksekusi satu kali proses pengecekan & pengarsipan produk.
	Run(ctx context.Context)
	// StartScheduler menjalankan Run secara berkala setiap `interval` sampai ctx dibatalkan.
	StartScheduler(ctx context.Context, interval time.Duration)
}

type produkAutoArchiveService struct {
	produkRepo      repositories.ProdukRepository
	activityLogRepo repositories.ActivityLogRepository
	// delay adalah ambang waktu "1 hari setelah terjual" sebelum produk diarsipkan otomatis.
	delay time.Duration
}

func NewProdukAutoArchiveService(
	produkRepo repositories.ProdukRepository,
	activityLogRepo repositories.ActivityLogRepository,
	delay time.Duration,
) ProdukAutoArchiveService {
	return &produkAutoArchiveService{
		produkRepo:      produkRepo,
		activityLogRepo: activityLogRepo,
		delay:           delay,
	}
}

func (s *produkAutoArchiveService) Run(ctx context.Context) {
	threshold := time.Now().Add(-s.delay)

	produkList, err := s.produkRepo.FindSoldProdukToArchive(ctx, threshold)
	if err != nil {
		log.Printf("[produk-auto-archive] gagal mengambil daftar produk: %v", err)
		return
	}

	if len(produkList) == 0 {
		return
	}

	archived := 0
	for _, produk := range produkList {
		if err := s.produkRepo.ArchiveProduk(ctx, produk.ID); err != nil {
			log.Printf("[produk-auto-archive] gagal mengarsipkan produk %s: %v", produk.ID, err)
			continue
		}
		archived++
		s.logArchive(produk)
	}

	log.Printf("[produk-auto-archive] %d/%d produk berhasil diarsipkan otomatis", archived, len(produkList))
}

func (s *produkAutoArchiveService) logArchive(produk models.Produk) {
	oldData, _ := json.Marshal(map[string]interface{}{"is_active": true})
	newData, _ := json.Marshal(map[string]interface{}{"is_active": false})

	entityType := "produk"
	entityID := produk.ID
	deskripsi := "Produk otomatis diarsipkan (non-aktif) karena sudah terjual lebih dari 1 hari: " + produk.NamaID

	entry := &models.ActivityLog{
		UserType:   "SYSTEM",
		Action:     models.ActionToggleStatus,
		Modul:      "produk",
		EntityType: &entityType,
		EntityID:   &entityID,
		Deskripsi:  deskripsi,
		OldData:    oldData,
		NewData:    newData,
	}

	if err := s.activityLogRepo.Create(entry); err != nil {
		log.Printf("[produk-auto-archive] gagal mencatat activity log untuk produk %s: %v", produk.ID, err)
	}
}

func (s *produkAutoArchiveService) StartScheduler(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[produk-auto-archive] scheduler dihentikan")
			return
		case <-ticker.C:
			s.Run(ctx)
		}
	}
}
