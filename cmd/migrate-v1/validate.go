package main

import "log"

// validate menjalankan pembandingan jumlah baris + cek constraint pasca-migrasi (dokumen §9).
func (a *App) validate() {
	countMy := func(q string) int64 {
		var n int64
		if err := a.my.QueryRow(q).Scan(&n); err != nil {
			log.Printf("  ! gagal hitung sumber (%s): %v", q, err)
			return -1
		}
		return n
	}
	countPg := func(q string) int64 {
		var n int64
		if err := a.pg.Raw(q).Scan(&n).Error; err != nil {
			log.Printf("  ! gagal hitung target (%s): %v", q, err)
			return -1
		}
		return n
	}

	type cmp struct{ label, myQ, pgQ string }
	for _, c := range []cmp{
		{"products vs produk", `SELECT COUNT(*) FROM products`, `SELECT COUNT(*) FROM produk`},
		{"users(is_admin=0) vs buyer", `SELECT COUNT(*) FROM users WHERE is_admin = 0`, `SELECT COUNT(*) FROM buyer`},
		{"users(is_admin=1)+admins vs admin", `SELECT (SELECT COUNT(*) FROM users WHERE is_admin = 1) + (SELECT COUNT(*) FROM admins)`, `SELECT COUNT(*) FROM admin`},
		{"addresses vs alamat_buyer", `SELECT COUNT(*) FROM addresses`, `SELECT COUNT(*) FROM alamat_buyer`},
		{"brand pivot vs produk_merek", `SELECT COUNT(*) FROM product_brand_pivot`, `SELECT COUNT(*) FROM produk_merek`},
		{"orders vs pesanan", `SELECT COUNT(*) FROM orders WHERE user_id IS NOT NULL`, `SELECT COUNT(*) FROM pesanan`},
		{"order_items vs pesanan_item", `SELECT COUNT(*) FROM order_items oi JOIN orders o ON o.id = oi.order_id WHERE o.user_id IS NOT NULL`, `SELECT COUNT(*) FROM pesanan_item`},
		{"invoices vs pesanan_pembayaran", `SELECT COUNT(*) FROM invoices`, `SELECT COUNT(*) FROM pesanan_pembayaran`},
		{"coupons vs kupon", `SELECT COUNT(*) FROM coupons`, `SELECT COUNT(*) FROM kupon`},
		{"coupon_category vs kupon_kategori", `SELECT COUNT(*) FROM coupon_category`, `SELECT COUNT(*) FROM kupon_kategori`},
		{"coupon_usages vs kupon_usage", `SELECT COUNT(*) FROM coupon_usages`, `SELECT COUNT(*) FROM kupon_usage`},
		{"reviews vs ulasan", `SELECT COUNT(*) FROM reviews WHERE order_id IS NOT NULL`, `SELECT COUNT(*) FROM ulasan`},
		{"user_consents vs buyer_disclaimer_consent", `SELECT COUNT(*) FROM user_consents`, `SELECT COUNT(*) FROM buyer_disclaimer_consent`},
		{"buyer bercart vs keranjang", `SELECT COUNT(DISTINCT c.user_id) FROM carts c JOIN cart_items ci ON ci.cart_id = c.id AND ci.deleted_at IS NULL WHERE c.deleted_at IS NULL`, `SELECT COUNT(*) FROM keranjang`},
		{"cart_items vs keranjang_item", `SELECT COUNT(*) FROM cart_items ci JOIN carts c ON c.id = ci.cart_id WHERE ci.deleted_at IS NULL AND c.deleted_at IS NULL`, `SELECT COUNT(*) FROM keranjang_item`},
	} {
		mv, pv := countMy(c.myQ), countPg(c.pgQ)
		note := ""
		if mv != pv {
			note = "  <- selisih; cocokkan dengan skip/seed di report"
		}
		log.Printf("  %-38s v1=%-7d v2=%-7d%s", c.label, mv, pv, note)
	}

	// cek constraint
	if n := countPg(`SELECT COUNT(*) FROM produk WHERE harga_sebelum_diskon <= 0`); n != 0 {
		log.Printf("  ! %d produk dengan harga_sebelum_diskon <= 0", n)
	}
	if n := countPg(`SELECT COUNT(*) FROM (SELECT telepon FROM buyer WHERE deleted_at IS NULL GROUP BY telepon HAVING COUNT(*) > 1) d`); n != 0 {
		log.Printf("  ! %d nomor telepon buyer aktif terduplikasi", n)
	}
	if n := countPg(`SELECT COUNT(*) FROM (SELECT slug FROM produk GROUP BY slug HAVING COUNT(*) > 1) d`); n != 0 {
		log.Printf("  ! %d slug produk terduplikasi", n)
	}
	if n := countPg(`SELECT COUNT(*) FROM (SELECT buyer_id FROM alamat_buyer WHERE is_default = true AND deleted_at IS NULL GROUP BY buyer_id HAVING COUNT(*) > 1) d`); n != 0 {
		log.Printf("  ! %d buyer dengan lebih dari satu alamat default", n)
	}

	// cek integritas transaksi (Fase 6-9)
	type chk struct{ label, q string }
	for _, c := range []chk{
		{"pesanan: total <> rincian biaya", `SELECT COUNT(*) FROM pesanan WHERE total <> biaya_produk + biaya_pengiriman + biaya_ppn + biaya_lainnya`},
		{"pesanan: kode terduplikasi", `SELECT COUNT(*) FROM (SELECT kode FROM pesanan GROUP BY kode HAVING COUNT(*) > 1) d`},
		{"pesanan: non-PICKUP tanpa alamat", `SELECT COUNT(*) FROM pesanan WHERE delivery_type <> 'PICKUP' AND alamat_buyer_id IS NULL AND COALESCE(alamat_snapshot, '') = ''`},
		{"pesanan_item: subtotal <> qty x (harga - diskon)", `SELECT COUNT(*) FROM pesanan_item WHERE subtotal <> qty * (harga_satuan - COALESCE(diskon_satuan, 0))`},
		{"pesanan_item: yatim (pesanan hilang)", `SELECT COUNT(*) FROM pesanan_item i WHERE NOT EXISTS (SELECT 1 FROM pesanan p WHERE p.id = i.pesanan_id)`},
		{"pembayaran: PAID tanpa paid_at", `SELECT COUNT(*) FROM pesanan_pembayaran WHERE status = 'PAID' AND paid_at IS NULL`},
		{"pembayaran: total PAID <> total pesanan", `SELECT COUNT(*) FROM (SELECT b.pesanan_id FROM pesanan_pembayaran b JOIN pesanan p ON p.id = b.pesanan_id WHERE b.status = 'PAID' GROUP BY b.pesanan_id, p.total HAVING SUM(b.jumlah) <> p.total) d`},
		{"kupon: kode terduplikasi", `SELECT COUNT(*) FROM (SELECT kode FROM kupon GROUP BY kode HAVING COUNT(*) > 1) d`},
		{"kupon: persentase > 100", `SELECT COUNT(*) FROM kupon WHERE jenis_diskon = 'persentase' AND nilai_diskon > 100`},
		{"ulasan: pesanan_item terduplikasi", `SELECT COUNT(*) FROM (SELECT pesanan_item_id FROM ulasan GROUP BY pesanan_item_id HAVING COUNT(*) > 1) d`},
		{"ulasan: rating di luar 1-5", `SELECT COUNT(*) FROM ulasan WHERE rating < 1 OR rating > 5`},
		{"keranjang: buyer terduplikasi", `SELECT COUNT(*) FROM (SELECT buyer_id FROM keranjang GROUP BY buyer_id HAVING COUNT(*) > 1) d`},
		{"keranjang: tanpa item", `SELECT COUNT(*) FROM keranjang k WHERE NOT EXISTS (SELECT 1 FROM keranjang_item i WHERE i.keranjang_id = k.id)`},
		{"keranjang_item: quantity <= 0", `SELECT COUNT(*) FROM keranjang_item WHERE quantity <= 0`},
	} {
		if n := countPg(c.q); n != 0 {
			log.Printf("  ! %s: %d baris", c.label, n)
		}
	}
	log.Print("  validasi selesai (mode dry-run: angka target belum berubah)")
}
