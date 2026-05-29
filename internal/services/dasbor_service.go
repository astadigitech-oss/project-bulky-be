package services

import (
	"context"
	"fmt"
	"math"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/repositories"
	"time"
)

type DasborService interface {
	GetChartTransaksi(ctx context.Context, periode string) (*dto.DasborChartTransaksiResponse, error)
	GetChartRevenue(ctx context.Context, periode string) (*dto.DasborChartRevenueResponse, error)
	GetChartTransaksiPerKategori(ctx context.Context, periode string) (*dto.DasborChartKategoriResponse, error)
	GetKPI(ctx context.Context, periode string) (*dto.DasborKPIResponse, error)
	GetStokPerKategori(ctx context.Context) (*dto.DasborStokPerKategoriResponse, error)
	GetPenjualanPerBuyer(ctx context.Context, periode string, limit int) (*dto.DasborPenjualanPerBuyerResponse, error)
	GetTabelTransaksi(ctx context.Context, periode string, page, perPage int) ([]dto.DasborTabelTransaksiItem, *dto.DasborTabelTransaksiMeta, error)
	GetAllTransaksiForExport(ctx context.Context, periode string) ([]dto.DasborTabelTransaksiItem, error)
	GetUserTransaction(ctx context.Context, periode string) ([]dto.DasborUserTransaksiItem, error)
}

type dasborService struct {
	repo repositories.DasborRepository
}

func NewDasborService(repo repositories.DasborRepository) DasborService {
	return &dasborService{repo: repo}
}

// jakartaLocation is the Asia/Jakarta timezone
var jakartaLocation = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.UTC
	}
	return loc
}()

// generateDailyLabels generates day labels for the given periode in Asia/Jakarta timezone.
// Returns (labels []string, dates []string) where dates are "YYYY-MM-DD" strings.
func generateDailyLabels(periode string) ([]string, []string) {
	now := time.Now().In(jakartaLocation)

	var start, end time.Time
	switch periode {
	case "tahun_ini":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, jakartaLocation)
		end = now
	default: // bulan_ini
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, jakartaLocation)
		end = now
	}

	monthNames := [...]string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}

	var labels, dates []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		label := fmt.Sprintf("%s-%02d", monthNames[d.Month()-1], d.Day())
		labels = append(labels, label)
		dates = append(dates, d.Format("2006-01-02"))
	}
	return labels, dates
}

func (s *dasborService) GetChartTransaksi(ctx context.Context, periode string) (*dto.DasborChartTransaksiResponse, error) {
	successRows, err := s.repo.GetChartTransaksiSuccess(periode)
	if err != nil {
		return nil, err
	}
	cancelRows, err := s.repo.GetChartTransaksiCancel(periode)
	if err != nil {
		return nil, err
	}

	labels, dates := generateDailyLabels(periode)

	successMap := rowsToDateMap(successRows)
	cancelMap := rowsToDateMap(cancelRows)

	successData := make([]int64, len(dates))
	cancelData := make([]int64, len(dates))
	for i, d := range dates {
		successData[i] = successMap[d]
		cancelData[i] = cancelMap[d]
	}

	return &dto.DasborChartTransaksiResponse{
		Periode: periode,
		Labels:  labels,
		Series: dto.DasborChartTransaksiSeriesData{
			Success: successData,
			Cancel:  cancelData,
		},
	}, nil
}

func (s *dasborService) GetChartRevenue(ctx context.Context, periode string) (*dto.DasborChartRevenueResponse, error) {
	rows, grandTotal, err := s.repo.GetChartRevenue(periode)
	if err != nil {
		return nil, err
	}

	labels, dates := generateDailyLabels(periode)

	revenueMap := make(map[string]float64, len(rows))
	for _, r := range rows {
		revenueMap[r.Tanggal] = r.TotalPenjualan
	}

	revenueData := make([]float64, len(dates))
	for i, d := range dates {
		revenueData[i] = revenueMap[d]
	}

	return &dto.DasborChartRevenueResponse{
		Periode: periode,
		Labels:  labels,
		Series: dto.DasborChartRevenueSeries{
			TotalPenjualan: revenueData,
		},
		TotalKeseluruhan: grandTotal,
	}, nil
}

func (s *dasborService) GetChartTransaksiPerKategori(ctx context.Context, periode string) (*dto.DasborChartKategoriResponse, error) {
	rows, err := s.repo.GetChartTransaksiPerKategori(periode)
	if err != nil {
		return nil, err
	}

	labels, dates := generateDailyLabels(periode)
	dateIndex := make(map[string]int, len(dates))
	for i, d := range dates {
		dateIndex[d] = i
	}

	// Collect unique categories (preserving order of first appearance)
	type kategoriMeta struct {
		nama string
		id   string
		idx  int
	}
	kategoriOrder := []kategoriMeta{}
	kategoriIdx := map[string]int{}

	// Build data matrix: kategori × date
	type cell struct{ count int64 }
	matrix := map[string][]cell{} // key = kategori_id

	for _, row := range rows {
		if _, exists := kategoriIdx[row.KategoriID]; !exists {
			idx := len(kategoriOrder)
			kategoriOrder = append(kategoriOrder, kategoriMeta{nama: row.Kategori, id: row.KategoriID, idx: idx})
			kategoriIdx[row.KategoriID] = idx
			matrix[row.KategoriID] = make([]cell, len(dates))
		}
		if di, ok := dateIndex[row.Tanggal]; ok {
			matrix[row.KategoriID][di].count += row.Jumlah
		}
	}

	series := make([]dto.DasborChartKategoriSerie, len(kategoriOrder))
	for i, km := range kategoriOrder {
		data := make([]int64, len(dates))
		for j, c := range matrix[km.id] {
			data[j] = c.count
		}
		series[i] = dto.DasborChartKategoriSerie{
			Kategori:   km.nama,
			KategoriID: km.id,
			Data:       data,
		}
	}

	return &dto.DasborChartKategoriResponse{
		Periode: periode,
		Labels:  labels,
		Series:  series,
	}, nil
}

func (s *dasborService) GetKPI(ctx context.Context, periode string) (*dto.DasborKPIResponse, error) {
	available, err := s.repo.GetKPIPaletboxAvailable()
	if err != nil {
		return nil, err
	}
	sold, err := s.repo.GetKPIPaletboxSold(periode)
	if err != nil {
		return nil, err
	}
	revenue, err := s.repo.GetKPIRevenue(periode)
	if err != nil {
		return nil, err
	}

	return &dto.DasborKPIResponse{
		Periode:           periode,
		PaletboxAvailable: available,
		PaletboxSold:      sold,
		Revenue:           revenue,
	}, nil
}

func (s *dasborService) GetStokPerKategori(ctx context.Context) (*dto.DasborStokPerKategoriResponse, error) {
	rows, err := s.repo.GetStokPerKategori()
	if err != nil {
		return nil, err
	}

	labels := make([]string, len(rows))
	stok := make([]int64, len(rows))
	for i, r := range rows {
		labels[i] = r.Kategori
		stok[i] = r.TotalStok
	}

	return &dto.DasborStokPerKategoriResponse{
		Labels: labels,
		Series: dto.DasborStokPerKategoriSeries{Stok: stok},
	}, nil
}

func (s *dasborService) GetPenjualanPerBuyer(ctx context.Context, periode string, limit int) (*dto.DasborPenjualanPerBuyerResponse, error) {
	rows, err := s.repo.GetPenjualanPerBuyer(periode, limit)
	if err != nil {
		return nil, err
	}

	labels := make([]string, len(rows))
	totals := make([]float64, len(rows))
	buyers := make([]dto.DasborBuyerDetail, len(rows))
	for i, r := range rows {
		labels[i] = r.Nama
		totals[i] = r.TotalPembelian
		buyers[i] = dto.DasborBuyerDetail{
			BuyerID:        r.BuyerID,
			Nama:           r.Nama,
			TotalPembelian: r.TotalPembelian,
		}
	}

	return &dto.DasborPenjualanPerBuyerResponse{
		Periode: periode,
		Labels:  labels,
		Series:  dto.DasborPenjualanPerBuyerSeries{TotalPembelian: totals},
		Buyers:  buyers,
	}, nil
}

func (s *dasborService) GetTabelTransaksi(ctx context.Context, periode string, page, perPage int) ([]dto.DasborTabelTransaksiItem, *dto.DasborTabelTransaksiMeta, error) {
	rows, total, err := s.repo.GetTabelTransaksi(periode, page, perPage)
	if err != nil {
		return nil, nil, err
	}

	items, err := s.enrichTabelRows(rows)
	if err != nil {
		return nil, nil, err
	}

	totalHalaman := int(math.Ceil(float64(total) / float64(perPage)))
	if totalHalaman < 1 {
		totalHalaman = 1
	}

	meta := &dto.DasborTabelTransaksiMeta{
		Halaman:      page,
		PerHalaman:   perPage,
		TotalData:    total,
		TotalHalaman: totalHalaman,
	}
	return items, meta, nil
}

func (s *dasborService) GetAllTransaksiForExport(ctx context.Context, periode string) ([]dto.DasborTabelTransaksiItem, error) {
	rows, err := s.repo.GetAllTransaksi(periode)
	if err != nil {
		return nil, err
	}
	return s.enrichTabelRows(rows)
}

func (s *dasborService) GetUserTransaction(ctx context.Context, periode string) ([]dto.DasborUserTransaksiItem, error) {
	rows, err := s.repo.GetUserTransaction(periode)
	if err != nil {
		return nil, err
	}

	items := make([]dto.DasborUserTransaksiItem, len(rows))
	for i, r := range rows {
		items[i] = dto.DasborUserTransaksiItem{
			BuyerID:        r.BuyerID,
			Nama:           r.Nama,
			TotalTransaksi: r.TotalTransaksi,
			TotalBelanja:   r.TotalBelanja,
		}
	}
	return items, nil
}

// enrichTabelRows fetches items and payment methods for a set of pesanan rows and assembles the response.
func (s *dasborService) enrichTabelRows(rows []dto.DasborTabelRow) ([]dto.DasborTabelTransaksiItem, error) {
	if len(rows) == 0 {
		return []dto.DasborTabelTransaksiItem{}, nil
	}

	pesananIDs := make([]string, len(rows))
	for i, r := range rows {
		pesananIDs[i] = r.PesananID
	}

	itemRows, err := s.repo.GetTabelTransaksiItems(pesananIDs)
	if err != nil {
		return nil, err
	}
	pembayaranRows, err := s.repo.GetTabelTransaksiPembayaran(pesananIDs)
	if err != nil {
		return nil, err
	}

	// Build item maps
	type itemData struct {
		namaFirst   string
		totalItems  int
		kategori    string
		totalDiskon float64
	}
	itemMap := map[string]*itemData{}
	for _, ir := range itemRows {
		data, exists := itemMap[ir.PesananID]
		if !exists {
			data = &itemData{}
			itemMap[ir.PesananID] = data
		}
		if ir.RowNum == 1 {
			data.namaFirst = ir.NamaProduk
			data.totalItems = ir.TotalItems
			data.kategori = ir.KategoriNama
		}
		data.totalDiskon += ir.DiskonSatuan * float64(ir.Qty)
	}

	// Build pembayaran map
	pembayaranMap := map[string][]string{}
	for _, pr := range pembayaranRows {
		pembayaranMap[pr.PesananID] = append(pembayaranMap[pr.PesananID], pr.NamaMetode)
	}

	items := make([]dto.DasborTabelTransaksiItem, len(rows))
	for i, r := range rows {
		palet := ""
		kategori := ""
		totalDiskon := 0.0

		if data, ok := itemMap[r.PesananID]; ok {
			palet = data.namaFirst
			if data.totalItems > 1 {
				palet = fmt.Sprintf("%s +%d lainnya", data.namaFirst, data.totalItems-1)
			}
			kategori = data.kategori
			totalDiskon = data.totalDiskon
		}

		jenisPembayaran := pembayaranMap[r.PesananID]
		if jenisPembayaran == nil {
			jenisPembayaran = []string{}
		}

		items[i] = dto.DasborTabelTransaksiItem{
			PesananID:       r.PesananID,
			Kode:            r.Kode,
			NamaPembeli:     r.NamaPembeli,
			Palet:           palet,
			Kategori:        kategori,
			Harga:           r.BiayaProduk,
			OngkosKirim:     r.OngkosKirim,
			Diskon:          totalDiskon,
			Total:           r.Total,
			TanggalPesanan:  r.CreatedAt,
			DeliveryType:    r.DeliveryType,
			PaymentType:     r.PaymentType,
			OrderStatus:     r.OrderStatus,
			JenisPembayaran: jenisPembayaran,
		}
	}
	return items, nil
}

// rowsToDateMap converts DasborDateCount slice to map[tanggal]jumlah.
func rowsToDateMap(rows []dto.DasborDateCount) map[string]int64 {
	m := make(map[string]int64, len(rows))
	for _, r := range rows {
		m[r.Tanggal] = r.Jumlah
	}
	return m
}
