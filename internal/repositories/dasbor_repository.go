package repositories

import (
	"fmt"
	"project-bulky-be/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type DasborRepository interface {
	GetChartTransaksiSuccess(periode string) ([]dto.DasborDateCount, error)
	GetChartTransaksiCancel(periode string) ([]dto.DasborDateCount, error)
	GetChartRevenue(periode string) ([]dto.DasborDateAmount, float64, error)
	GetChartTransaksiPerKategori(periode string) ([]dto.DasborKategoriDateCount, error)
	GetKPIPaletboxAvailable() (int64, error)
	GetKPIPaletboxSold(periode string) (int64, error)
	GetKPIRevenue(periode string) (float64, error)
	GetStokPerKategori() ([]dto.DasborKategoriStok, error)
	GetPenjualanPerBuyer(periode string, limit int) ([]dto.DasborBuyerPenjualan, error)
	GetTabelTransaksi(periode string, page, perPage int) ([]dto.DasborTabelRow, int64, error)
	GetTabelTransaksiItems(pesananIDs []string) ([]dto.DasborTabelItemRow, error)
	GetTabelTransaksiPembayaran(pesananIDs []string) ([]dto.DasborTabelPembayaranRow, error)
	GetAllTransaksi(periode string) ([]dto.DasborTabelRow, error)
	GetUserTransaction(periode string) ([]dto.DasborUserTransaksiRaw, error)
}

type dasborRepository struct {
	db *gorm.DB
}

func NewDasborRepository(db *gorm.DB) DasborRepository {
	return &dasborRepository{db: db}
}

const jakartaTZ = "Asia/Jakarta"

// buildPeriodeWhereClause returns a SQL WHERE fragment for the periode filter.
// tableAlias.column should be the timestamptz column reference.
func buildPeriodeWhereClause(col, periode string) string {
	switch periode {
	case "tahun_ini":
		return fmt.Sprintf(
			"DATE_TRUNC('year', %s AT TIME ZONE '%s') = DATE_TRUNC('year', NOW() AT TIME ZONE '%s')",
			col, jakartaTZ, jakartaTZ,
		)
	case "semua":
		return ""
	default: // bulan_ini
		return fmt.Sprintf(
			"DATE_TRUNC('month', %s AT TIME ZONE '%s') = DATE_TRUNC('month', NOW() AT TIME ZONE '%s')",
			col, jakartaTZ, jakartaTZ,
		)
	}
}

// GetChartTransaksiSuccess returns daily count of COMPLETED orders.
func (r *dasborRepository) GetChartTransaksiSuccess(periode string) ([]dto.DasborDateCount, error) {
	periodeCond := buildPeriodeWhereClause("created_at", periode)
	query := `
		SELECT
			DATE(created_at AT TIME ZONE $1)::text AS tanggal,
			COUNT(*) AS jumlah
		FROM pesanan
		WHERE order_status = 'COMPLETED'
		  AND deleted_at IS NULL`
	if periodeCond != "" {
		query += " AND " + periodeCond
	}
	query += " GROUP BY tanggal ORDER BY tanggal"

	var rows []dto.DasborDateCount
	if err := r.db.Raw(query, jakartaTZ).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetChartTransaksiCancel returns daily count of CANCELLED orders.
func (r *dasborRepository) GetChartTransaksiCancel(periode string) ([]dto.DasborDateCount, error) {
	periodeCond := buildPeriodeWhereClause("created_at", periode)
	query := `
		SELECT
			DATE(created_at AT TIME ZONE $1)::text AS tanggal,
			COUNT(*) AS jumlah
		FROM pesanan
		WHERE order_status = 'CANCELLED'
		  AND deleted_at IS NULL`
	if periodeCond != "" {
		query += " AND " + periodeCond
	}
	query += " GROUP BY tanggal ORDER BY tanggal"

	var rows []dto.DasborDateCount
	if err := r.db.Raw(query, jakartaTZ).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetChartRevenue returns daily revenue from PAID orders.
func (r *dasborRepository) GetChartRevenue(periode string) ([]dto.DasborDateAmount, float64, error) {
	periodeCond := buildPeriodeWhereClause("paid_at", periode)
	query := `
		SELECT
			DATE(paid_at AT TIME ZONE $1)::text AS tanggal,
			SUM(total) AS total_penjualan
		FROM pesanan
		WHERE payment_status = 'PAID'
		  AND deleted_at IS NULL
		  AND paid_at IS NOT NULL`
	if periodeCond != "" {
		query += " AND " + periodeCond
	}
	query += " GROUP BY tanggal ORDER BY tanggal"

	var rows []dto.DasborDateAmount
	if err := r.db.Raw(query, jakartaTZ).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	var total float64
	for _, row := range rows {
		total += row.TotalPenjualan
	}
	return rows, total, nil
}

// GetChartTransaksiPerKategori returns transaction count per day per kategori.
func (r *dasborRepository) GetChartTransaksiPerKategori(periode string) ([]dto.DasborKategoriDateCount, error) {
	periodeCond := buildPeriodeWhereClause("p.created_at", periode)
	query := `
		SELECT
			DATE(p.created_at AT TIME ZONE $1)::text AS tanggal,
			kp.nama_id AS kategori,
			kp.id::text AS kategori_id,
			COUNT(DISTINCT p.id) AS jumlah_transaksi
		FROM pesanan p
		JOIN pesanan_item pi ON pi.pesanan_id = p.id
		JOIN produk pr ON pr.id = pi.produk_id
		JOIN kategori_produk kp ON kp.id = pr.kategori_id
		WHERE p.deleted_at IS NULL
		  AND p.order_status != 'CANCELLED'
		  AND kp.deleted_at IS NULL`
	if periodeCond != "" {
		query += " AND " + periodeCond
	}
	query += " GROUP BY tanggal, kp.id, kp.nama_id ORDER BY tanggal, kp.nama_id"

	var rows []dto.DasborKategoriDateCount
	if err := r.db.Raw(query, jakartaTZ).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetKPIPaletboxAvailable returns current count of available (unsold) active products.
func (r *dasborRepository) GetKPIPaletboxAvailable() (int64, error) {
	var total int64
	err := r.db.Raw(`
		SELECT COUNT(*) AS paletbox_available
		FROM produk
		WHERE is_sold = false
		  AND is_active = true AND deleted_at IS NULL
	`).Scan(&total).Error
	return total, err
}

// GetKPIPaletboxSold returns total qty sold in the given periode.
func (r *dasborRepository) GetKPIPaletboxSold(periode string) (int64, error) {
	periodeCond := buildPeriodeWhereClause("p.created_at", periode)
	query := `
		SELECT COALESCE(SUM(pi.qty), 0)
		FROM pesanan_item pi
		JOIN pesanan p ON p.id = pi.pesanan_id
		WHERE p.order_status = 'COMPLETED'
		  AND p.deleted_at IS NULL`
	if periodeCond != "" {
		query += " AND " + periodeCond
	}

	var total int64
	if err := r.db.Raw(query).Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// GetKPIRevenue returns total revenue from PAID orders in the given periode.
func (r *dasborRepository) GetKPIRevenue(periode string) (float64, error) {
	periodeCond := buildPeriodeWhereClause("p.paid_at", periode)
	query := `
		SELECT COALESCE(SUM(p.total), 0)
		FROM pesanan p
		WHERE p.payment_status = 'PAID'
		  AND p.deleted_at IS NULL
		  AND p.paid_at IS NOT NULL`
	if periodeCond != "" {
		query += " AND " + periodeCond
	}

	var total float64
	if err := r.db.Raw(query).Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// GetStokPerKategori returns current available palet count per kategori (all categories, zero if none).
func (r *dasborRepository) GetStokPerKategori() ([]dto.DasborKategoriStok, error) {
	var rows []dto.DasborKategoriStok
	err := r.db.Raw(`
		SELECT
			kp.nama_id AS kategori,
			COUNT(p.id) AS total_stok
		FROM kategori_produk kp
		LEFT JOIN produk p ON p.kategori_id = kp.id
			AND p.is_sold = false
			AND p.is_active = true
			AND p.deleted_at IS NULL
		WHERE kp.deleted_at IS NULL
		GROUP BY kp.id, kp.nama_id
		ORDER BY kp.nama_id ASC
	`).Scan(&rows).Error
	return rows, err
}

// GetPenjualanPerBuyer returns top N buyers by total revenue.
func (r *dasborRepository) GetPenjualanPerBuyer(periode string, limit int) ([]dto.DasborBuyerPenjualan, error) {
	periodeCond := buildPeriodeWhereClause("p.paid_at", periode)
	query := `
		SELECT
			b.id::text AS buyer_id,
			b.nama AS nama,
			SUM(p.total) AS total_pembelian
		FROM pesanan p
		JOIN buyer b ON b.id = p.buyer_id
		WHERE p.payment_status = 'PAID'
		  AND p.deleted_at IS NULL
		  AND p.paid_at IS NOT NULL
		  AND b.deleted_at IS NULL`
	if periodeCond != "" {
		query += " AND " + periodeCond
	}
	query += fmt.Sprintf(" GROUP BY b.id, b.nama ORDER BY total_pembelian DESC LIMIT %d", limit)

	var rows []dto.DasborBuyerPenjualan
	if err := r.db.Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetTabelTransaksi returns paginated transaction rows.
func (r *dasborRepository) GetTabelTransaksi(periode string, page, perPage int) ([]dto.DasborTabelRow, int64, error) {
	periodeCond := buildPeriodeWhereClause("p.created_at", periode)
	whereClause := "p.deleted_at IS NULL"
	if periodeCond != "" {
		whereClause += " AND " + periodeCond
	}

	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM pesanan p WHERE %s`, whereClause)
	if err := r.db.Raw(countQuery).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	dataQuery := fmt.Sprintf(`
		SELECT
			p.id::text AS pesanan_id,
			p.kode AS kode,
			b.nama AS nama_pembeli,
			p.biaya_produk AS biaya_produk,
			p.biaya_pengiriman AS ongkos_kirim,
			p.biaya_ppn AS biaya_ppn,
			p.total AS total,
			DATE(p.created_at AT TIME ZONE '%s')::text AS tanggal_pesanan,
			p.delivery_type::text AS delivery_type,
			p.payment_type::text AS payment_type,
			p.order_status::text AS order_status
		FROM pesanan p
		JOIN buyer b ON b.id = p.buyer_id
		WHERE %s
		ORDER BY p.created_at DESC
		LIMIT %d OFFSET %d
	`, jakartaTZ, whereClause, perPage, offset)

	var rows []dto.DasborTabelRow
	if err := r.db.Raw(dataQuery).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// GetTabelTransaksiItems returns item rows for given pesanan IDs (first item + count).
func (r *dasborRepository) GetTabelTransaksiItems(pesananIDs []string) ([]dto.DasborTabelItemRow, error) {
	if len(pesananIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(pesananIDs))
	args := make([]interface{}, len(pesananIDs))
	for i, id := range pesananIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT
			pi.pesanan_id::text AS pesanan_id,
			pi.nama_produk AS nama_produk,
			pi.diskon_satuan AS diskon_satuan,
			pi.qty AS qty,
			kp.nama_id AS kategori_nama,
			ROW_NUMBER() OVER (PARTITION BY pi.pesanan_id ORDER BY pi.created_at ASC) AS row_num,
			COUNT(*) OVER (PARTITION BY pi.pesanan_id) AS total_items
		FROM pesanan_item pi
		JOIN produk pr ON pr.id = pi.produk_id
		JOIN kategori_produk kp ON kp.id = pr.kategori_id
		WHERE pi.pesanan_id::text IN (%s)
	`, strings.Join(placeholders, ","))

	var rows []dto.DasborTabelItemRow
	if err := r.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetTabelTransaksiPembayaran returns payment method names for given pesanan IDs.
func (r *dasborRepository) GetTabelTransaksiPembayaran(pesananIDs []string) ([]dto.DasborTabelPembayaranRow, error) {
	if len(pesananIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(pesananIDs))
	args := make([]interface{}, len(pesananIDs))
	for i, id := range pesananIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT
			pp.pesanan_id::text AS pesanan_id,
			mp.nama AS nama_metode
		FROM pesanan_pembayaran pp
		JOIN metode_pembayaran mp ON mp.id = pp.metode_pembayaran_id
		WHERE pp.pesanan_id::text IN (%s)
		  AND pp.metode_pembayaran_id IS NOT NULL
	`, strings.Join(placeholders, ","))

	var rows []dto.DasborTabelPembayaranRow
	if err := r.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetAllTransaksi returns all transactions for export (no pagination).
func (r *dasborRepository) GetAllTransaksi(periode string) ([]dto.DasborTabelRow, error) {
	periodeCond := buildPeriodeWhereClause("p.created_at", periode)
	whereClause := "p.deleted_at IS NULL"
	if periodeCond != "" {
		whereClause += " AND " + periodeCond
	}

	dataQuery := fmt.Sprintf(`
		SELECT
			p.id::text AS pesanan_id,
			p.kode AS kode,
			b.nama AS nama_pembeli,
			p.biaya_produk AS biaya_produk,
			p.biaya_pengiriman AS ongkos_kirim,
			p.biaya_ppn AS biaya_ppn,
			p.total AS total,
			DATE(p.created_at AT TIME ZONE '%s')::text AS tanggal_pesanan,
			p.delivery_type::text AS delivery_type,
			p.payment_type::text AS payment_type,
			p.order_status::text AS order_status
		FROM pesanan p
		JOIN buyer b ON b.id = p.buyer_id
		WHERE %s
		ORDER BY p.created_at DESC
	`, jakartaTZ, whereClause)

	var rows []dto.DasborTabelRow
	if err := r.db.Raw(dataQuery).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetUserTransaction returns buyer ranking by total transactions.
func (r *dasborRepository) GetUserTransaction(periode string) ([]dto.DasborUserTransaksiRaw, error) {
	periodeCond := buildPeriodeWhereClause("p.created_at", periode)
	joinCond := "p.deleted_at IS NULL AND p.order_status != 'CANCELLED'"
	if periodeCond != "" {
		joinCond += " AND " + periodeCond
	}

	query := fmt.Sprintf(`
		SELECT
			b.id::text AS buyer_id,
			b.nama AS nama,
			COUNT(p.id) AS total_transaksi,
			COALESCE(SUM(p.total), 0) AS total_belanja
		FROM buyer b
		JOIN pesanan p ON p.buyer_id = b.id AND %s
		WHERE b.deleted_at IS NULL
		GROUP BY b.id, b.nama
		ORDER BY total_transaksi DESC
	`, joinCond)

	var rows []dto.DasborUserTransaksiRaw
	if err := r.db.Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
