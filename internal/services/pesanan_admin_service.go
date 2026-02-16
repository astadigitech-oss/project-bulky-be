package services

import (
	"context"
	"errors"
	"os"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PesananAdminService interface {
	GetAll(ctx context.Context, params *dto.PesananAdminQueryParams) ([]dto.PesananAdminListResponse, *models.PaginationMeta, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.PesananAdminDetailResponse, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, req *dto.UpdatePesananStatusRequest, adminID uuid.UUID) (*dto.UpdatePesananStatusResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetStatistics(ctx context.Context, tanggalDari, tanggalSampai *string) (*dto.PesananStatisticsResponse, error)
}

type pesananAdminService struct {
	pesananRepo repositories.PesananRepository
	db          *gorm.DB
}

func NewPesananAdminService(pesananRepo repositories.PesananRepository, db *gorm.DB) PesananAdminService {
	return &pesananAdminService{
		pesananRepo: pesananRepo,
		db:          db,
	}
}

func (s *pesananAdminService) GetAll(ctx context.Context, params *dto.PesananAdminQueryParams) ([]dto.PesananAdminListResponse, *models.PaginationMeta, error) {
	// Build filters
	filters := make(map[string]interface{})
	if params.Cari != "" {
		filters["cari"] = params.Cari
	}
	if params.OrderStatus != "" {
		filters["order_status"] = params.OrderStatus
	}
	if params.PaymentStatus != "" {
		filters["payment_status"] = params.PaymentStatus
	}
	if params.DeliveryType != "" {
		filters["delivery_type"] = params.DeliveryType
	}
	if params.TanggalDari != "" {
		tanggal, err := time.Parse("2006-01-02", params.TanggalDari)
		if err == nil {
			filters["tanggal_dari"] = tanggal
		}
	}
	if params.TanggalSampai != "" {
		tanggal, err := time.Parse("2006-01-02", params.TanggalSampai)
		if err == nil {
			filters["tanggal_sampai"] = tanggal
		}
	}

	// Get data from repository
	pesanan, total, err := s.pesananRepo.AdminFindAll(filters, params.Page, params.PerPage, params.SortBy, params.SortOrder)
	if err != nil {
		return nil, nil, err
	}

	// Map to response DTO
	response := make([]dto.PesananAdminListResponse, len(pesanan))
	for i, p := range pesanan {
		response[i] = s.mapToListResponse(&p)
	}

	// Build meta
	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return response, &meta, nil
}

func (s *pesananAdminService) GetByID(ctx context.Context, id uuid.UUID) (*dto.PesananAdminDetailResponse, error) {
	pesanan, err := s.pesananRepo.AdminFindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pesanan tidak ditemukan")
		}
		return nil, err
	}

	// Get status history
	var statusHistory []models.PesananStatusHistory
	if err := s.db.Where("pesanan_id = ?", id).Order("created_at ASC").Find(&statusHistory).Error; err != nil {
		return nil, err
	}

	return s.mapToDetailResponse(pesanan, statusHistory), nil
}

func (s *pesananAdminService) UpdateStatus(ctx context.Context, id uuid.UUID, req *dto.UpdatePesananStatusRequest, adminID uuid.UUID) (*dto.UpdatePesananStatusResponse, error) {
	// Get current pesanan
	pesanan, err := s.pesananRepo.AdminFindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pesanan tidak ditemukan")
		}
		return nil, err
	}

	previousStatus := string(pesanan.OrderStatus)

	// Update status
	orderStatus := models.OrderStatus(req.OrderStatus)
	if err := s.pesananRepo.UpdateStatus(id, orderStatus, req.Note, adminID); err != nil {
		return nil, err
	}

	return &dto.UpdatePesananStatusResponse{
		ID:             id,
		Kode:           pesanan.Kode,
		OrderStatus:    req.OrderStatus,
		PreviousStatus: previousStatus,
		UpdatedAt:      time.Now(),
		UpdatedBy:      adminID,
	}, nil
}

func (s *pesananAdminService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if pesanan exists
	_, err := s.pesananRepo.AdminFindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("pesanan tidak ditemukan")
		}
		return err
	}

	return s.pesananRepo.Delete(id)
}

func (s *pesananAdminService) GetStatistics(ctx context.Context, tanggalDari, tanggalSampai *string) (*dto.PesananStatisticsResponse, error) {
	var dari, sampai *time.Time

	// Parse dates
	if tanggalDari != nil && *tanggalDari != "" {
		t, err := time.Parse("2006-01-02", *tanggalDari)
		if err != nil {
			return nil, errors.New("format tanggal_dari tidak valid")
		}
		dari = &t
	} else {
		// Default: 30 days ago
		t := time.Now().AddDate(0, 0, -30)
		dari = &t
	}

	if tanggalSampai != nil && *tanggalSampai != "" {
		t, err := time.Parse("2006-01-02", *tanggalSampai)
		if err != nil {
			return nil, errors.New("format tanggal_sampai tidak valid")
		}
		sampai = &t
	} else {
		// Default: today
		t := time.Now()
		sampai = &t
	}

	// Get statistics
	stats, err := s.pesananRepo.GetStatistics(dari, sampai)
	if err != nil {
		return nil, err
	}

	// Map to response DTO
	response := &dto.PesananStatisticsResponse{
		TotalPesanan:     stats["total_pesanan"].(int64),
		TotalRevenue:     stats["total_revenue"].(decimal.Decimal),
		PerStatus:        stats["per_status"].(map[string]int64),
		PerDeliveryType:  stats["per_delivery_type"].(map[string]int64),
		PerPaymentStatus: stats["per_payment_status"].(map[string]int64),
	}

	return response, nil
}

// Helper functions to map models to DTOs

func (s *pesananAdminService) mapToListResponse(p *models.Pesanan) dto.PesananAdminListResponse {
	totalItem := len(p.Items)

	return dto.PesananAdminListResponse{
		ID:   p.ID,
		Kode: p.Kode,
		Buyer: dto.PesananAdminBuyerResponse{
			ID:    p.BuyerID,
			Nama:  p.Buyer.Nama,
			Email: p.Buyer.Email,
		},
		DeliveryType:    string(p.DeliveryType),
		PaymentType:     string(p.PaymentType),
		PaymentStatus:   string(p.PaymentStatus),
		OrderStatus:     string(p.OrderStatus),
		TotalItem:       totalItem,
		BiayaProduk:     p.BiayaProduk,
		BiayaPengiriman: p.BiayaPengiriman,
		BiayaPPN:        p.BiayaPPN,
		TotalBayar:      p.Total,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func (s *pesananAdminService) mapToDetailResponse(p *models.Pesanan, statusHistory []models.PesananStatusHistory) *dto.PesananAdminDetailResponse {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Map items
	items := make([]dto.PesananAdminItemResponse, len(p.Items))
	for i, item := range p.Items {
		var gambarURL *string
		if len(item.Produk.Gambar) > 0 {
			url := baseURL + "/uploads/" + item.Produk.Gambar[0].GambarURL
			gambarURL = &url
		}

		items[i] = dto.PesananAdminItemResponse{
			ID: item.ID,
			Produk: dto.PesananAdminItemProdukResponse{
				ID:        item.ProdukID,
				Nama:      item.Produk.NamaID,
				Slug:      item.Produk.Slug,
				GambarURL: gambarURL,
			},
			NamaProduk:   item.NamaProduk,
			Qty:          item.Qty,
			HargaSatuan:  item.HargaSatuan,
			DiskonSatuan: item.DiskonSatuan,
			Subtotal:     item.Subtotal,
		}
	}

	// Map pembayaran
	pembayaran := make([]dto.PesananAdminPembayaranResponse, len(p.Pembayaran))
	for i, bayar := range p.Pembayaran {
		pembayaran[i] = dto.PesananAdminPembayaranResponse{
			ID:              bayar.ID,
			BuyerID:         bayar.BuyerID,
			NamaPembayar:    bayar.Buyer.Nama,
			Jumlah:          bayar.Jumlah,
			Status:          string(bayar.Status),
			PaidAt:          bayar.PaidAt,
			XenditInvoiceID: bayar.XenditInvoiceID,
		}

		if bayar.MetodePembayaran != nil {
			pembayaran[i].MetodePembayaran = dto.PesananAdminMetodePembayaranResponse{
				ID:   bayar.MetodePembayaran.ID,
				Nama: bayar.MetodePembayaran.Nama,
				Kode: bayar.MetodePembayaran.Kode,
			}
		}
	}

	// Map status history
	history := make([]dto.PesananAdminStatusHistoryResponse, len(statusHistory))
	for i, h := range statusHistory {
		history[i] = dto.PesananAdminStatusHistoryResponse{
			StatusFrom: h.StatusFrom,
			StatusTo:   h.StatusTo,
			StatusType: string(h.StatusType),
			Note:       h.Note,
			CreatedAt:  h.CreatedAt,
		}
	}

	response := &dto.PesananAdminDetailResponse{
		ID:   p.ID,
		Kode: p.Kode,
		Buyer: dto.PesananAdminBuyerDetailResponse{
			ID:    p.BuyerID,
			Nama:  p.Buyer.Nama,
			Email: p.Buyer.Email,
			Telepon: func() string {
				if p.Buyer.Telepon != nil {
					return *p.Buyer.Telepon
				}
				return ""
			}(),
		},
		DeliveryType:    string(p.DeliveryType),
		PaymentType:     string(p.PaymentType),
		PaymentStatus:   string(p.PaymentStatus),
		OrderStatus:     string(p.OrderStatus),
		Items:           items,
		Pembayaran:      pembayaran,
		StatusHistory:   history,
		BiayaProduk:     p.BiayaProduk,
		BiayaPengiriman: p.BiayaPengiriman,
		BiayaPPN:        p.BiayaPPN,
		PotonganKupon:   decimal.Zero, // TODO: implement kupon if needed
		TotalBayar:      p.Total,
		CatatanBuyer:    p.Catatan,
		CatatanAdmin:    p.CatatanAdmin,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}

	// Map alamat if exists
	if p.AlamatBuyer != nil {
		response.AlamatPengiriman = &dto.PesananAdminAlamatResponse{
			ID:            p.AlamatBuyer.ID,
			Label:         p.AlamatBuyer.Label,
			NamaPenerima:  p.AlamatBuyer.NamaPenerima,
			Telepon:       p.AlamatBuyer.TeleponPenerima,
			AlamatLengkap: p.AlamatBuyer.AlamatLengkap,
			Kota:          p.AlamatBuyer.Kota,
			Provinsi:      p.AlamatBuyer.Provinsi,
			KodePos: func() string {
				if p.AlamatBuyer.KodePos != nil {
					return *p.AlamatBuyer.KodePos
				}
				return ""
			}(),
		}
	}

	return response
}
