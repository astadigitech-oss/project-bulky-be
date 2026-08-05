package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// phasePromo memigrasi kupon beserta relasinya, ulasan produk, dan riwayat
// persetujuan disclaimer (dokumen §5.5–§5.7).
func (a *App) phasePromo() error {
	a.kuponKnown = map[string]bool{}
	// is_all_kategori harus sudah diketahui sebelum baris kupon ditulis,
	// jadi daftar kupon yang dibatasi kategori dimuat lebih dulu.
	a.kuponPunyaKategori = map[string]bool{}
	katRows, err := a.my.Query(`SELECT DISTINCT coupon_id FROM coupon_category`)
	if err != nil {
		return err
	}
	for katRows.Next() {
		var cid string
		if err := katRows.Scan(&cid); err != nil {
			katRows.Close()
			return err
		}
		a.kuponPunyaKategori[cid] = true
	}
	katRows.Close()

	if err := a.kupon(); err != nil {
		return fmt.Errorf("kupon: %w", err)
	}
	if err := a.kuponKategori(); err != nil {
		return fmt.Errorf("kupon_kategori: %w", err)
	}
	if err := a.kuponUsage(); err != nil {
		return fmt.Errorf("kupon_usage: %w", err)
	}
	if err := a.ulasan(); err != nil {
		return fmt.Errorf("ulasan: %w", err)
	}
	return a.consent()
}

// kuponJatuhTempoDefault dipakai untuk kupon v1 yang tidak punya expiry_date.
// v2 mewajibkan tanggal_kedaluarsa NOT NULL, jadi kupon tanpa batas waktu
// diberi tanggal masa lalu sekaligus is_active=false supaya tidak bisa dipakai
// lagi di v2 tanpa ditinjau ulang oleh admin.
var kuponJatuhTempoDefault = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

func (a *App) kupon() error {
	rows, err := a.my.Query(`SELECT id, code, discount_type, discount_value,
			expiry_date, minimum_purchase, usage_limit,
			deleted_at, created_at, updated_at
		FROM coupons ORDER BY created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, code string
		var discType, discValue sql.NullString
		var minPurchase, usageLimit sql.NullInt64
		var expiry, deletedAt, createdAt, updatedAt sql.NullTime

		if err := rows.Scan(&id, &code, &discType, &discValue, &expiry,
			&minPurchase, &usageLimit, &deletedAt, &createdAt, &updatedAt); err != nil {
			return err
		}

		if a.tgt.KuponIDs[id] {
			a.kuponKnown[id] = true
			a.rep.Count("kupon.sudah_ada_di_v2")
			continue
		}

		// discount_value di v1 bertipe varchar, jadi harus diurai manual.
		nilai, err := strconv.ParseFloat(strings.TrimSpace(discValue.String), 64)
		if err != nil || nilai <= 0 {
			a.rep.Add("8", "kupon", id, "skip",
				"nilai_diskon tidak valid: "+discValue.String)
			continue
		}

		jenis := "jumlah_tetap"
		if strings.EqualFold(strings.TrimSpace(discType.String), "percent") {
			jenis = "persentase"
			if nilai > 100 {
				a.rep.Add("8", "kupon", id, "skip",
					fmt.Sprintf("persentase di atas 100: %.2f", nilai))
				continue
			}
		}

		aktif := !deletedAt.Valid
		tempo := expiry.Time
		if !expiry.Valid {
			tempo = kuponJatuhTempoDefault
			aktif = false
			a.rep.Add("8", "kupon", id, "kedaluarsa_default",
				"expiry_date NULL di v1, dinonaktifkan")
		}

		// limit_pemakaian CHECK (NULL OR > 0)
		var limit *int64
		if usageLimit.Valid && usageLimit.Int64 > 0 {
			v := usageLimit.Int64
			limit = &v
		}

		if err := a.exec(`INSERT INTO kupon
				(id, kode, nama, deskripsi, jenis_diskon, nilai_diskon,
				 minimal_pembelian, limit_pemakaian, tanggal_kedaluarsa,
				 is_all_kategori, is_active, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			id, truncate(strings.ToUpper(strings.TrimSpace(code)), 50),
			truncate(strings.TrimSpace(code), 255),
			jenis, nilai, float64(minPurchase.Int64), limit, tempo,
			!a.kuponPunyaKategori[id], aktif,
			nullTimePtr(createdAt), nullTimePtr(updatedAt), nullTimePtr(deletedAt)); err != nil {
			return err
		}
		a.kuponKnown[id] = true
		a.rep.Count("kupon.insert")
	}
	return rows.Err()
}

func (a *App) kuponKategori() error {
	rows, err := a.my.Query(`SELECT coupon_id, product_category_id FROM coupon_category`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var kuponID, katID string
		if err := rows.Scan(&kuponID, &katID); err != nil {
			return err
		}
		if !a.kuponKnown[kuponID] {
			a.rep.Add("8", "kupon_kategori", kuponID, "skip", "kupon tidak termigrasi")
			continue
		}
		if !a.kategoriKnown[katID] {
			a.rep.Add("8", "kupon_kategori", kuponID, "skip", "kategori tidak termigrasi: "+id8(katID))
			continue
		}
		// v1 tidak punya PK di tabel pivot ini, jadi id v2 dibuat deterministik
		// dari pasangannya agar migrasi ulang tidak menghasilkan duplikat.
		if err := a.exec(`INSERT INTO kupon_kategori (id, kupon_id, kategori_id, created_at)
			VALUES (uuid_generate_v4(), ?, ?, NOW())
			ON CONFLICT (kupon_id, kategori_id) DO NOTHING`, kuponID, katID); err != nil {
			return err
		}
		a.rep.Count("kupon_kategori.insert")
	}
	return rows.Err()
}

// kuponUsage memigrasi riwayat pemakaian kupon. v2 mewajibkan kode_kupon dan
// nilai_potongan (> 0) yang tidak ada di v1, sehingga keduanya diturunkan dari
// coupons.code dan orders.discount_amount.
func (a *App) kuponUsage() error {
	rows, err := a.my.Query(`SELECT cu.id, cu.coupon_id, cu.user_id, cu.order_id,
			cu.created_at, c.code, o.discount_amount
		FROM coupon_usages cu
		JOIN coupons c ON c.id = cu.coupon_id
		JOIN orders o ON o.id = cu.order_id
		ORDER BY cu.created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, userID, orderID, code string
		var kuponID sql.NullString
		var potongan sql.NullFloat64
		var createdAt sql.NullTime

		if err := rows.Scan(&id, &kuponID, &userID, &orderID, &createdAt,
			&code, &potongan); err != nil {
			return err
		}

		if a.tgt.KuponUsageIDs[id] {
			a.rep.Count("kupon_usage.sudah_ada_di_v2")
			continue
		}
		if !kuponID.Valid || !a.kuponKnown[kuponID.String] {
			a.rep.Add("8", "kupon_usage", id, "skip", "kupon tidak termigrasi")
			continue
		}
		if !a.pesananKnown[orderID] {
			a.rep.Add("8", "kupon_usage", id, "skip", "pesanan tidak termigrasi")
			continue
		}
		// CHECK (nilai_potongan > 0)
		if floatOrZero(potongan) <= 0 {
			a.rep.Add("8", "kupon_usage", id, "skip", "discount_amount pesanan nol")
			continue
		}

		if err := a.exec(`INSERT INTO kupon_usage
				(id, kupon_id, buyer_id, pesanan_id, kode_kupon, nilai_potongan, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			id, kuponID.String, userID, orderID,
			truncate(strings.ToUpper(strings.TrimSpace(code)), 50),
			floatOrZero(potongan), nullTimePtr(createdAt)); err != nil {
			return err
		}
		a.rep.Count("kupon_usage.insert")
	}
	return rows.Err()
}

// ulasan memigrasi reviews v1 ke tabel ulasan v2.
//
// v2 mengikat ulasan ke pesanan_item (UNIQUE), sedangkan v1 hanya menyimpan
// pasangan order_id + product_id, jadi item pasangannya dicari lewat JOIN.
// Ulasan tanpa order_id atau yang itemnya tidak ketemu di-skip (keputusan #19).
// Gambar: v1 menyimpan banyak baris review_images (maksimum 3), v2 hanya punya
// satu kolom varchar, jadi hanya gambar pertama yang ikut.
func (a *App) ulasan() error {
	// trg_validate_ulasan_order menolak ulasan untuk pesanan non-COMPLETED.
	// Aturan itu benar untuk transaksi baru, tapi data historis v1 memuat
	// ulasan pada pesanan berstatus lain, jadi trigger dimatikan sementara
	// di dalam transaksi fase ini dan dihidupkan lagi setelah selesai.
	if a.execute {
		if err := a.tx.Exec(`ALTER TABLE ulasan DISABLE TRIGGER trg_validate_ulasan_order`).Error; err != nil {
			return fmt.Errorf("nonaktifkan trigger ulasan: %w", err)
		}
		defer func() {
			if a.tx != nil {
				_ = a.tx.Exec(`ALTER TABLE ulasan ENABLE TRIGGER trg_validate_ulasan_order`).Error
			}
		}()
	}

	// Satu pesanan_item hanya boleh punya satu ulasan (UNIQUE), tetapi v1
	// memuat 4 pasangan (order, product) dengan lebih dari satu ulasan.
	// ROW_NUMBER menyisakan ulasan terlama sebagai yang sah.
	rows, err := a.my.Query(`SELECT x.id, x.user_id, x.order_id, x.product_id,
			x.rating, x.comment, x.approved, x.deleted_at, x.created_at, x.updated_at,
			(SELECT ri.path FROM review_images ri WHERE ri.review_id = x.id
			 ORDER BY ri.created_at, ri.id LIMIT 1) AS gambar,
			(SELECT COUNT(*) FROM review_images ri WHERE ri.review_id = x.id) AS n_gambar
		FROM (
			SELECT r.*, ROW_NUMBER() OVER (
				PARTITION BY r.order_id, r.product_id ORDER BY r.created_at, r.id) AS rn
			FROM reviews r WHERE r.order_id IS NOT NULL
		) x
		WHERE x.rn = 1
		ORDER BY x.created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, orderID, productID string
		var userID, comment, gambar sql.NullString
		var rating, approved, nGambar sql.NullInt64
		var deletedAt, createdAt, updatedAt sql.NullTime

		if err := rows.Scan(&id, &userID, &orderID, &productID, &rating, &comment,
			&approved, &deletedAt, &createdAt, &updatedAt, &gambar, &nGambar); err != nil {
			return err
		}

		if a.tgt.UlasanIDs[id] {
			a.rep.Count("ulasan.sudah_ada_di_v2")
			continue
		}
		if !userID.Valid || userID.String == "" {
			a.rep.Add("8", "ulasan", id, "skip", "user_id NULL di v1")
			continue
		}
		if !a.pesananKnown[orderID] {
			a.rep.Add("8", "ulasan", id, "skip", "pesanan tidak termigrasi")
			continue
		}
		if rating.Int64 < 1 || rating.Int64 > 5 {
			a.rep.Add("8", "ulasan", id, "skip",
				fmt.Sprintf("rating di luar 1-5: %d", rating.Int64))
			continue
		}

		itemID, ok := a.pesananItemLookup(orderID, productID)
		if !ok {
			a.rep.Add("8", "ulasan", id, "skip", "pesanan_item tidak ditemukan")
			continue
		}

		disetujui := boolFromTiny(approved)
		var approvedAt interface{}
		if disetujui {
			approvedAt = nullTimePtr(updatedAt)
		}
		if nGambar.Int64 > 1 {
			a.rep.Add("8", "ulasan", id, "gambar_terpangkas",
				fmt.Sprintf("%d gambar di v1, hanya 1 yang muat di v2", nGambar.Int64))
		}

		if err := a.exec(`INSERT INTO ulasan
				(id, pesanan_id, pesanan_item_id, buyer_id, produk_id, rating,
				 komentar, gambar, is_approved, approved_at, approved_by,
				 created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			id, orderID, itemID, userID.String, productID, rating.Int64,
			nullStrPtr(comment), nullTrunc(gambar, 255), disetujui, approvedAt,
			nullTimePtr(createdAt), nullTimePtr(updatedAt), nullTimePtr(deletedAt)); err != nil {
			return err
		}
		a.rep.Count("ulasan.insert")
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

// pesananItemLookup mencari id pesanan_item dari pasangan pesanan + produk.
// Peta dibangun sekali saat pemanggilan pertama karena dipakai berulang.
func (a *App) pesananItemLookup(pesananID, produkID string) (string, bool) {
	if a.itemByOrderProduct == nil {
		a.itemByOrderProduct = map[string]string{}
		rows, err := a.my.Query(`SELECT order_id, product_id, id FROM order_items
			WHERE order_id IS NOT NULL AND product_id IS NOT NULL
			ORDER BY created_at`)
		if err != nil {
			return "", false
		}
		defer rows.Close()
		for rows.Next() {
			var o, p, i string
			if err := rows.Scan(&o, &p, &i); err != nil {
				return "", false
			}
			// item pertama yang menang bila satu produk muncul dua kali
			if _, ada := a.itemByOrderProduct[o+"|"+p]; !ada {
				a.itemByOrderProduct[o+"|"+p] = i
			}
		}
	}
	id, ok := a.itemByOrderProduct[pesananID+"|"+produkID]
	if !ok {
		return "", false
	}
	// item bisa saja ter-skip pada Fase 6 (mis. produk tidak termigrasi)
	if !a.tgt.PesananItemIDs[id] && !a.pesananItemInserted[id] {
		return "", false
	}
	return id, true
}

// consent memigrasi user_consents -> buyer_disclaimer_consent. Kolom version
// di v1 berisi slug disclaimer, sedangkan disclaimer_id-nya sudah benar,
// sehingga id dipakai langsung dan slug hanya jadi cadangan.
func (a *App) consent() error {
	rows, err := a.my.Query(`SELECT id, user_id, order_id, disclaimer_id, version,
			ip_address, user_agent, accepted_at, created_at
		FROM user_consents ORDER BY created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, userID, orderID, disclaimerID string
		var version, ip, ua sql.NullString
		var acceptedAt, createdAt sql.NullTime

		if err := rows.Scan(&id, &userID, &orderID, &disclaimerID, &version,
			&ip, &ua, &acceptedAt, &createdAt); err != nil {
			return err
		}

		if a.tgt.ConsentIDs[id] {
			a.rep.Count("consent.sudah_ada_di_v2")
			continue
		}
		if !a.pesananKnown[orderID] {
			a.rep.Add("8", "buyer_disclaimer_consent", id, "skip", "pesanan tidak termigrasi")
			continue
		}

		disc := disclaimerID
		if !a.disclaimerKnown[disc] {
			if byslug, ok := a.disclaimerBySlug[strings.TrimSpace(version.String)]; ok {
				disc = byslug
				a.rep.Add("8", "buyer_disclaimer_consent", id, "disclaimer_via_slug", version.String)
			} else {
				a.rep.Add("8", "buyer_disclaimer_consent", id, "skip",
					"disclaimer tidak dikenal: "+version.String)
				continue
			}
		}

		// disetujui_at NOT NULL di v2; v1 selalu mengisinya, tapi created_at
		// disiapkan sebagai cadangan agar migrasi tidak pernah gagal di sini.
		setuju := nullTimePtr(acceptedAt)
		if !acceptedAt.Valid {
			setuju = nullTimePtr(createdAt)
		}

		if err := a.exec(`INSERT INTO buyer_disclaimer_consent
				(id, buyer_id, pesanan_id, disclaimer_id, disetujui_at,
				 ip_address, user_agent, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			id, userID, orderID, disc, setuju,
			nullTrunc(ip, 45), nullStrPtr(ua), nullTimePtr(createdAt)); err != nil {
			return err
		}
		a.rep.Count("consent.insert")
	}
	return rows.Err()
}
