package controllers

import (
	"fmt"
	"net/http"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

type DasborController struct {
	svc services.DasborService
}

func NewDasborController(svc services.DasborService) *DasborController {
	return &DasborController{svc: svc}
}

// GetChartTransaksi handles GET /api/panel/dasbor/chart-transaksi
func (c *DasborController) GetChartTransaksi(ctx *fiber.Ctx) error {
	var q dto.DasborPeriodeQuery
	if err := ctx.QueryParser(&q); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}
	q.SetDefault()

	result, err := c.svc.GetChartTransaksi(ctx.UserContext(), q.Periode)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data chart transaksi", err.Error())
	}
	return utils.SuccessResponse(ctx, "Data chart transaksi berhasil diambil", result)
}

// GetChartRevenue handles GET /api/panel/dasbor/chart-revenue
func (c *DasborController) GetChartRevenue(ctx *fiber.Ctx) error {
	var q dto.DasborPeriodeQuery
	if err := ctx.QueryParser(&q); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}
	q.SetDefault()

	result, err := c.svc.GetChartRevenue(ctx.UserContext(), q.Periode)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data chart revenue", err.Error())
	}
	return utils.SuccessResponse(ctx, "Data chart revenue berhasil diambil", result)
}

// GetChartTransaksiPerKategori handles GET /api/panel/dasbor/chart-transaksi-per-kategori
func (c *DasborController) GetChartTransaksiPerKategori(ctx *fiber.Ctx) error {
	var q dto.DasborPeriodeQuery
	if err := ctx.QueryParser(&q); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}
	q.SetDefault()

	result, err := c.svc.GetChartTransaksiPerKategori(ctx.UserContext(), q.Periode)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data chart transaksi per kategori", err.Error())
	}
	return utils.SuccessResponse(ctx, "Data chart transaksi per kategori berhasil diambil", result)
}

// GetKPI handles GET /api/panel/dasbor/kpi
func (c *DasborController) GetKPI(ctx *fiber.Ctx) error {
	var q dto.DasborPeriodeQuery
	if err := ctx.QueryParser(&q); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}
	q.SetDefault()

	result, err := c.svc.GetKPI(ctx.UserContext(), q.Periode)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data KPI", err.Error())
	}
	return utils.SuccessResponse(ctx, "Data KPI berhasil diambil", result)
}

// GetStokPerKategori handles GET /api/panel/dasbor/stok-per-kategori
func (c *DasborController) GetStokPerKategori(ctx *fiber.Ctx) error {
	result, err := c.svc.GetStokPerKategori(ctx.UserContext())
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data stok per kategori", err.Error())
	}
	return utils.SuccessResponse(ctx, "Data stok per kategori berhasil diambil", result)
}

// GetPenjualanPerBuyer handles GET /api/panel/dasbor/penjualan-per-buyer
func (c *DasborController) GetPenjualanPerBuyer(ctx *fiber.Ctx) error {
	var q dto.DasborPenjualanPerBuyerQuery
	if err := ctx.QueryParser(&q); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}
	q.SetDefaults()

	result, err := c.svc.GetPenjualanPerBuyer(ctx.UserContext(), q.Periode, q.Limit)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data penjualan per buyer", err.Error())
	}
	return utils.SuccessResponse(ctx, "Data penjualan per buyer berhasil diambil", result)
}

// GetTabelTransaksi handles GET /api/panel/dasbor/tabel-transaksi
func (c *DasborController) GetTabelTransaksi(ctx *fiber.Ctx) error {
	var q dto.DasborTabelTransaksiQuery
	if err := ctx.QueryParser(&q); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}
	q.SetDefaults()

	items, meta, err := c.svc.GetTabelTransaksi(ctx.UserContext(), q.Periode, q.Halaman, q.PerHalaman)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data tabel transaksi", err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data tabel transaksi berhasil diambil",
		"data":    items,
		"meta":    meta,
	})
}

// EksporTransaksi handles GET /api/panel/dasbor/ekspor-transaksi
func (c *DasborController) EksporTransaksi(ctx *fiber.Ctx) error {
	var q dto.DasborPeriodeQuery
	if err := ctx.QueryParser(&q); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}
	q.SetDefault()

	items, err := c.svc.GetAllTransaksiForExport(ctx.UserContext(), q.Periode)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data ekspor transaksi", err.Error())
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Transaksi"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{
		"No", "Kode Pesanan", "Nama Pembeli", "Nama Palet", "Kategori",
		"Harga Produk", "Ongkos Kirim", "Diskon", "PPN", "Total",
		"Delivery Type", "Payment Type", "Jenis Pembayaran", "Status Pesanan", "Tanggal Pesanan",
	}
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, item := range items {
		row := i + 2
		ppn := item.Total - item.Harga - item.OngkosKirim + item.Diskon
		f.SetCellValue(sheet, cellName(1, row), i+1)
		f.SetCellValue(sheet, cellName(2, row), item.Kode)
		f.SetCellValue(sheet, cellName(3, row), item.NamaPembeli)
		f.SetCellValue(sheet, cellName(4, row), item.Palet)
		f.SetCellValue(sheet, cellName(5, row), item.Kategori)
		f.SetCellValue(sheet, cellName(6, row), item.Harga)
		f.SetCellValue(sheet, cellName(7, row), item.OngkosKirim)
		f.SetCellValue(sheet, cellName(8, row), item.Diskon)
		f.SetCellValue(sheet, cellName(9, row), ppn)
		f.SetCellValue(sheet, cellName(10, row), item.Total)
		f.SetCellValue(sheet, cellName(11, row), item.DeliveryType)
		f.SetCellValue(sheet, cellName(12, row), item.PaymentType)
		f.SetCellValue(sheet, cellName(13, row), strings.Join(item.JenisPembayaran, ", "))
		f.SetCellValue(sheet, cellName(14, row), item.OrderStatus)
		f.SetCellValue(sheet, cellName(15, row), item.TanggalPesanan)
	}

	tanggal := time.Now().Format("20060102")
	filename := fmt.Sprintf("transaksi-bulky-%s-%s.xlsx", q.Periode, tanggal)

	ctx.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	if err := f.Write(ctx.Response().BodyWriter()); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat file Excel", err.Error())
	}
	return nil
}

// GetUserTransaction handles GET /api/panel/dasbor/user-transaction
func (c *DasborController) GetUserTransaction(ctx *fiber.Ctx) error {
	var q dto.DasborPeriodeQuery
	if err := ctx.QueryParser(&q); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}
	q.SetDefault()

	items, err := c.svc.GetUserTransaction(ctx.UserContext(), q.Periode)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data user transaction", err.Error())
	}
	return utils.SuccessResponse(ctx, "Data user transaction berhasil diambil", items)
}

func cellName(col, row int) string {
	name, _ := excelize.CoordinatesToCellName(col, row)
	return name
}
