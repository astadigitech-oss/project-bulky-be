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
	log.Print("  validasi selesai (mode dry-run: angka target belum berubah)")
}
