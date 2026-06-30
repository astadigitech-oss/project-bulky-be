package services

import (
	"context"
	"errors"
	"os"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PesananAdminService interface {
	GetAll(ctx context.Context, params *dto.PesananAdminQueryParams) ([]dto.PesananAdminListResponse, *models.PaginationMeta, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.PesananAdminDetailResponse, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, req *dto.UpdatePesananStatusRequest, adminID uuid.UUID) (*dto.UpdatePesananStatusResponse, error)
	RetryBooking(ctx context.Context, id uuid.UUID) (*dto.RetryBookingResponse, error)
	TrackDelivery(ctx context.Context, id uuid.UUID) (*TrackingResult, error)
	GetForwarderInvoice(ctx context.Context, id uuid.UUID) ([]ForwarderInvoice, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetStatistics(ctx context.Context, params *dto.StatisticsQueryParams) (*dto.PesananStatisticsResponse, error)
}

type pesananAdminService struct {
	pesananRepo     repositories.PesananRepository
	shippingService ShippingService
	db              *gorm.DB
}

func NewPesananAdminService(pesananRepo repositories.PesananRepository, shippingService ShippingService, db *gorm.DB) PesananAdminService {
	return &pesananAdminService{
		pesananRepo:     pesananRepo,
		shippingService: shippingService,
		db:              db,
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

	// PICKUP tidak punya status SHIPPED — buyer ambil sendiri ke gudang
	if pesanan.DeliveryType == models.DeliveryTypePickup && orderStatus == models.OrderStatusShipped {
		return nil, errors.New("pesanan tipe PICKUP tidak memiliki status SHIPPED")
	}

	if err := s.pesananRepo.UpdateStatus(id, orderStatus, req.Note, adminID); err != nil {
		return nil, err
	}

	// Trigger booking async when status → READY for DELIVEREE/FORWARDER
	if orderStatus == models.OrderStatusReady &&
		(pesanan.DeliveryType == models.DeliveryTypeDeliveree || pesanan.DeliveryType == models.DeliveryTypeForwarder) {
		s.shippingService.TriggerBookingAsync(pesanan)
	}

	return &dto.UpdatePesananStatusResponse{
		ID:             id,
		Kode:           pesanan.Kode,
		OrderStatus:    req.OrderStatus,
		PreviousStatus: previousStatus,
		UpdatedAt:      time.Now().UTC(),
		UpdatedBy:      adminID,
	}, nil
}

func (s *pesananAdminService) RetryBooking(ctx context.Context, id uuid.UUID) (*dto.RetryBookingResponse, error) {
	pesanan, err := s.pesananRepo.AdminFindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pesanan tidak ditemukan")
		}
		return nil, err
	}

	if pesanan.DeliveryType == models.DeliveryTypePickup {
		return nil, errors.New("retry:bad_request:Retry tidak diperlukan. Pesanan tipe PICKUP tidak memerlukan booking.")
	}
	if pesanan.OrderStatus != models.OrderStatusProcessing && pesanan.OrderStatus != models.OrderStatusShipped {
		return nil, errors.New("retry:bad_request:Retry hanya bisa dilakukan pada pesanan berstatus PROCESSING atau SHIPPED.")
	}

	// Already booked
	if pesanan.DeliveryType == models.DeliveryTypeDeliveree && pesanan.DelivereeBookingID != nil {
		return nil, errors.New("retry:already_booked:" + *pesanan.DelivereeBookingID)
	}
	if pesanan.DeliveryType == models.DeliveryTypeForwarder && pesanan.ForwarderTrackingNo != nil {
		return nil, errors.New("retry:already_booked:" + *pesanan.ForwarderTrackingNo)
	}

	// Run synchronous booking
	delivereeID, trackingNo, bookErr := s.shippingService.BookDelivery(ctx, pesanan)
	if bookErr != nil {
		errMsg := bookErr.Error()
		_ = s.pesananRepo.UpdateBookingResult(id, nil, nil, &errMsg)

		// City not mapped — distinguish 422
		if strings.Contains(errMsg, "tidak ditemukan di Forwarder mapping") {
			return nil, errors.New("retry:city_not_mapped:" + errMsg)
		}
		return nil, errors.New("retry:provider_error:" + string(pesanan.DeliveryType) + ":" + errMsg)
	}

	_ = s.pesananRepo.UpdateBookingResult(id, delivereeID, trackingNo, nil)

	return &dto.RetryBookingResponse{
		PesananID:    id.String(),
		DeliveryType: string(pesanan.DeliveryType),
		BookingID:    delivereeID,
		TrackingNo:   trackingNo,
	}, nil
}


func (s *pesananAdminService) TrackDelivery(ctx context.Context, id uuid.UUID) (*TrackingResult, error) {
	pesanan, err := s.pesananRepo.AdminFindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pesanan tidak ditemukan")
		}
		return nil, err
	}

	if pesanan.DeliveryType == models.DeliveryTypePickup {
		return nil, errors.New("tracking:not_applicable:Pesanan tipe PICKUP tidak memiliki tracking pengiriman")
	}

	return s.shippingService.TrackDelivery(ctx, pesanan)
}

func (s *pesananAdminService) GetForwarderInvoice(ctx context.Context, id uuid.UUID) ([]ForwarderInvoice, error) {
	pesanan, err := s.pesananRepo.AdminFindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pesanan tidak ditemukan")
		}
		return nil, err
	}

	if pesanan.DeliveryType != models.DeliveryTypeForwarder && pesanan.DeliveryType != models.DeliveryTypeForwarderLCL {
		return nil, errors.New("invoice:not_applicable:Pesanan ini tidak menggunakan layanan Forwarder")
	}

	if pesanan.ForwarderTrackingNo == nil || *pesanan.ForwarderTrackingNo == "" {
		return nil, errors.New("invoice:not_applicable:Pesanan belum memiliki booking number Forwarder")
	}

	return s.shippingService.GetForwarderInvoice(ctx, *pesanan.ForwarderTrackingNo)
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

func (s *pesananAdminService) GetStatistics(ctx context.Context, params *dto.StatisticsQueryParams) (*dto.PesananStatisticsResponse, error) {
	now := time.Now().UTC()
	var chartDari, chartSampai time.Time
	var summaryDari, summarySampai *time.Time
	var groupBy string
	isWeekFilter := false

	switch {
	case params.Tahun != nil && params.Bulan != nil:
		// Bulan tertentu → group by hari
		year, month := *params.Tahun, time.Month(*params.Bulan)
		chartDari = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		chartSampai = chartDari.AddDate(0, 1, 0).Add(-time.Nanosecond)
		groupBy = "day"

	case params.Tahun != nil && params.Minggu != nil:
		// Minggu tertentu → group by hari (7 data point)
		chartDari = isoWeekToDate(*params.Tahun, *params.Minggu)
		chartSampai = chartDari.AddDate(0, 0, 7).Add(-time.Nanosecond)
		groupBy = "day"
		isWeekFilter = true

	case params.Tahun != nil:
		// Tahun tertentu → group by bulan
		year := *params.Tahun
		chartDari = time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		chartSampai = time.Date(year, time.December, 31, 23, 59, 59, 999999999, time.UTC)
		groupBy = "month"

	case params.TanggalDari != "" || params.TanggalSampai != "":
		// Custom range
		if params.TanggalDari != "" {
			t, err := time.Parse("2006-01-02", params.TanggalDari)
			if err != nil {
				return nil, errors.New("format tanggal_dari tidak valid")
			}
			chartDari = t
		} else {
			chartDari = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		}
		if params.TanggalSampai != "" {
			t, err := time.Parse("2006-01-02", params.TanggalSampai)
			if err != nil {
				return nil, errors.New("format tanggal_sampai tidak valid")
			}
			chartSampai = t.Add(24*time.Hour - time.Nanosecond)
		} else {
			chartSampai = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.UTC)
		}
		// Auto group: <=90 hari → per hari, >90 hari → per bulan
		if int(chartSampai.Sub(chartDari).Hours()/24) > 90 {
			groupBy = "month"
		} else {
			groupBy = "day"
		}

	default:
		// Default: tahun berjalan, group by bulan
		chartDari = time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
		chartSampai = time.Date(now.Year(), time.December, 31, 23, 59, 59, 999999999, time.UTC)
		groupBy = "month"
	}

	// Summary stats ikut filter chart jika filter aktif, otherwise all-time
	if params.Tahun != nil || params.TanggalDari != "" || params.TanggalSampai != "" {
		summaryDari = &chartDari
		summarySampai = &chartSampai
	}

	stats, err := s.pesananRepo.GetStatistics(summaryDari, summarySampai)
	if err != nil {
		return nil, err
	}

	rawPoints, err := s.pesananRepo.GetChartData(&chartDari, &chartSampai, groupBy)
	if err != nil {
		return nil, err
	}

	return &dto.PesananStatisticsResponse{
		TotalPesanan:     stats["total_pesanan"].(int64),
		TotalRevenue:     stats["total_revenue"].(decimal.Decimal),
		PerStatus:        stats["per_status"].(map[string]int64),
		PerDeliveryType:  stats["per_delivery_type"].(map[string]int64),
		PerPaymentStatus: stats["per_payment_status"].(map[string]int64),
		ChartData:        buildChartData(rawPoints, chartDari, chartSampai, groupBy, isWeekFilter),
	}, nil
}

func buildShippingInfo(p *models.Pesanan) dto.PesananShippingInfo {
	info := dto.PesananShippingInfo{
		DeliveryType: string(p.DeliveryType),
		BookingID:    p.DelivereeBookingID,
		TrackingNo:   p.ForwarderTrackingNo,
		BookingError: p.BookingError,
	}

	switch {
	case p.DeliveryType == models.DeliveryTypePickup:
		info.BookingStatus = "NOT_APPLICABLE"
	case p.BookingError != nil:
		info.BookingStatus = "FAILED"
	case p.OrderStatus == models.OrderStatusReady && p.DelivereeBookingID == nil && p.ForwarderTrackingNo == nil:
		info.BookingStatus = "IN_PROGRESS"
	case p.DelivereeBookingID != nil || p.ForwarderTrackingNo != nil:
		info.BookingStatus = "BOOKED"
	default:
		info.BookingStatus = "PENDING"
	}

	return info
}

// buildChartData fills all periods in range with data (zero for missing periods)
func buildChartData(rawPoints []models.ChartRawPoint, dari, sampai time.Time, groupBy string, isWeekFilter bool) []dto.ChartDataPoint {
	dataMap := make(map[string]int64)
	for _, p := range rawPoints {
		dataMap[p.Period] = p.TotalPesanan
	}

	var result []dto.ChartDataPoint

	if groupBy == "month" {
		current := time.Date(dari.Year(), dari.Month(), 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(sampai.Year(), sampai.Month(), 1, 0, 0, 0, 0, time.UTC)
		for !current.After(end) {
			period := current.Format("2006-01")
			result = append(result, dto.ChartDataPoint{
				Label:        indonesianMonthShort(current.Month()),
				Period:       period,
				TotalPesanan: dataMap[period],
			})
			current = current.AddDate(0, 1, 0)
		}
	} else {
		current := time.Date(dari.Year(), dari.Month(), dari.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(sampai.Year(), sampai.Month(), sampai.Day(), 0, 0, 0, 0, time.UTC)
		for !current.After(end) {
			period := current.Format("2006-01-02")
			var label string
			if isWeekFilter {
				label = indonesianDayShort(current.Weekday())
			} else {
				label = current.Format("02")
			}
			result = append(result, dto.ChartDataPoint{
				Label:        label,
				Period:       period,
				TotalPesanan: dataMap[period],
			})
			current = current.AddDate(0, 0, 1)
		}
	}

	if result == nil {
		return []dto.ChartDataPoint{}
	}
	return result
}

func indonesianMonthShort(m time.Month) string {
	names := [13]string{"", "Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}
	return names[int(m)]
}

func indonesianDayShort(d time.Weekday) string {
	names := map[time.Weekday]string{
		time.Sunday: "Min", time.Monday: "Sen", time.Tuesday: "Sel",
		time.Wednesday: "Rab", time.Thursday: "Kam", time.Friday: "Jum", time.Saturday: "Sab",
	}
	return names[d]
}

// isoWeekToDate returns Monday of the given ISO week in the given year
func isoWeekToDate(year, week int) time.Time {
	jan4 := time.Date(year, time.January, 4, 0, 0, 0, 0, time.UTC)
	weekday := int(jan4.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := jan4.AddDate(0, 0, 1-weekday)
	return monday.AddDate(0, 0, (week-1)*7)
}

// Helper functions to map models to DTOs

func (s *pesananAdminService) mapToListResponse(p *models.Pesanan) dto.PesananAdminListResponse {
	totalItem := len(p.Items)

	createdAt := p.CreatedAt.UTC()

	return dto.PesananAdminListResponse{
		ID: p.ID,
		Buyer: dto.PesananAdminBuyerResponse{
			ID:   p.BuyerID,
			Nama: p.Buyer.Nama,
		},
		Kode:        p.Kode,
		TotalBayar:  p.Total,
		TotalItem:   totalItem,
		PaymentType: string(p.PaymentType),
		Status:      string(p.OrderStatus),
		OrderAt:     createdAt,
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
			ID:           bayar.ID,
			BuyerID:      bayar.BuyerID,
			NamaPembayar: bayar.Buyer.Nama,
			Jumlah:       bayar.Jumlah,
			Status:       string(bayar.Status),
			PaidAt: func() *time.Time {
				if bayar.PaidAt == nil {
					return nil
				}
				t := bayar.PaidAt.UTC()
				return &t
			}(),
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
			CreatedAt:  h.CreatedAt.UTC(),
		}
	}

	response := &dto.PesananAdminDetailResponse{
		ID:   p.ID,
		Kode: p.Kode,
		Buyer: dto.PesananAdminBuyerDetailResponse{
			ID:   p.BuyerID,
			Nama: p.Buyer.Nama,
			Email: func() string {
				if p.Buyer.Email != nil {
					return *p.Buyer.Email
				}
				return ""
			}(),
			Telepon: p.Buyer.Telepon,
		},
		DeliveryType:    string(p.DeliveryType),
		ShippingInfo:    buildShippingInfo(p),
		PaymentType:     string(p.PaymentType),
		PaymentStatus:   string(p.PaymentStatus),
		OrderStatus:     string(p.OrderStatus),
		Items:           items,
		Pembayaran:      pembayaran,
		StatusHistory:   history,
		BiayaProduk:     p.BiayaProduk,
		BiayaPengiriman: p.BiayaPengiriman,
		BiayaPPN:        p.BiayaPPN,
		BiayaLainnya:    p.BiayaLainnya,
		PotonganKupon:   decimal.Zero, // TODO: implement kupon if needed
		TotalBayar:      p.Total,
		CatatanBuyer:    p.Catatan,
		CatatanAdmin:    p.CatatanAdmin,
		CreatedAt:       p.CreatedAt.UTC(),
		UpdatedAt:       p.UpdatedAt.UTC(),
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
