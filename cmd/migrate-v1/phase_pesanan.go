package main

import (
	"database/sql"
	"fmt"
	"strings"
)

// phasePesanan memigrasi orders + order_shippings -> pesanan, lalu
// order_items -> pesanan_item (dokumen §5.1 dan §5.2).
func (a *App) phasePesanan() error {
	if err := a.pesananHeader(); err != nil {
		return fmt.Errorf("pesanan: %w", err)
	}
	return a.pesananItem()
}

// mapDeliveryType memetakan shipping_method + shipping_provider v1 ke enum delivery_type v2.
//
// v1 hanya mengenal self_pickup / courier_pickup, sedangkan penyedia kurir
// disimpan di order_shippings.shipping_provider. Mayoritas order courier_pickup
// (322 dari 436) tidak punya baris order_shippings sama sekali, sehingga
// penyedianya tidak diketahui; kasus itu dipetakan ke FORWARDER sebagai default
// dan ditandai di catatan_admin agar operator tahu asal datanya.
func mapDeliveryType(method, provider string) (dtype string, catatan string) {
	m := strings.ToLower(strings.TrimSpace(method))
	p := strings.ToLower(strings.TrimSpace(provider))

	if m != "courier_pickup" {
		// self_pickup dan NULL sama-sama ambil di gudang
		return "PICKUP", ""
	}
	switch {
	case strings.Contains(p, "deliveree"):
		return "DELIVEREE", ""
	case strings.Contains(p, "forwarder"):
		return "FORWARDER", ""
	default:
		return "FORWARDER", "Migrasi v1: kurir tidak tercatat di data lama, diasumsikan FORWARDER."
	}
}

// mapPaymentStatus memetakan payment_status v1 ke enum payment_status v2.
// Keputusan #17: CANCELED v1 tidak punya padanan langsung; pesanan batal karena
// pembayaran tidak diselesaikan, jadi dipetakan ke EXPIRED.
func mapPaymentStatus(v1 string) string {
	switch strings.ToUpper(strings.TrimSpace(v1)) {
	case "PAID":
		return "PAID"
	case "CANCELED", "CANCELLED", "EXPIRED":
		return "EXPIRED"
	case "FAILED":
		return "FAILED"
	case "REFUNDED":
		return "REFUNDED"
	default:
		return "PENDING"
	}
}

// mapOrderStatus memetakan order_status v1 ke enum order_status v2.
func mapOrderStatus(v1 string) string {
	switch strings.ToLower(strings.TrimSpace(v1)) {
	case "processing":
		return "PROCESSING"
	case "ready_to_pickup", "ready":
		return "READY"
	case "shipped", "delivering":
		return "SHIPPED"
	case "delivered", "completed":
		return "COMPLETED"
	case "canceled", "cancelled":
		return "CANCELLED"
	default:
		return "PENDING"
	}
}

func (a *App) pesananHeader() error {
	a.pesananKnown = map[string]bool{}
	a.orderExists = map[string]bool{}

	rows, err := a.my.Query(`SELECT o.id, o.user_id, o.order_number, o.order_date,
			o.total_price, o.discount_amount, o.tax_amount,
			o.shipping_address, o.name, o.phone_number, o.latitude, o.longitude,
			o.shipping_method, o.payment_status, o.order_status,
			o.payment_expired_at, o.paid_off_at, o.cancel_reason, o.notes,
			o.tracking_number, o.deleted_at, o.created_at, o.updated_at,
			os.shipping_provider, os.shipping_cost, os.insurance_amount,
			os.booking_id, os.tracking_url
		FROM orders o
		LEFT JOIN order_shippings os ON os.order_id = o.id
		ORDER BY o.created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, orderNumber string
		var userID, shippingAddress, name, phone sql.NullString
		var shippingMethod, paymentStatus, orderStatus sql.NullString
		var cancelReason, notes, trackingNumber sql.NullString
		var shippingProvider, bookingID, trackingURL sql.NullString
		var lat, lng, totalPrice, discountAmount, taxAmount sql.NullFloat64
		var shippingCost, insuranceAmount sql.NullFloat64
		var orderDate, expiredAt, paidAt, deletedAt, createdAt, updatedAt sql.NullTime

		if err := rows.Scan(&id, &userID, &orderNumber, &orderDate,
			&totalPrice, &discountAmount, &taxAmount,
			&shippingAddress, &name, &phone, &lat, &lng,
			&shippingMethod, &paymentStatus, &orderStatus,
			&expiredAt, &paidAt, &cancelReason, &notes,
			&trackingNumber, &deletedAt, &createdAt, &updatedAt,
			&shippingProvider, &shippingCost, &insuranceAmount,
			&bookingID, &trackingURL); err != nil {
			return err
		}

		a.orderExists[id] = true

		if a.tgt.PesananIDs[id] {
			a.pesananKnown[id] = true
			a.rep.Count("pesanan.sudah_ada_di_v2")
			continue
		}

		// Keputusan #18: pesanan tanpa pemilik tidak bisa dimigrasi
		// (buyer_id NOT NULL) dan tidak berguna bagi siapa pun.
		if !userID.Valid || userID.String == "" {
			a.rep.Add("6", "pesanan", id, "skip", "user_id NULL di v1")
			continue
		}
		if !a.buyerKnown[userID.String] && !a.tgt.BuyerIDs[userID.String] {
			a.rep.Add("6", "pesanan", id, "skip", "buyer tidak termigrasi: "+id8(userID.String))
			continue
		}

		deliveryType, catatanAdmin := mapDeliveryType(shippingMethod.String, shippingProvider.String)
		payStatus := mapPaymentStatus(paymentStatus.String)
		ordStatus := mapOrderStatus(orderStatus.String)

		alamat := strings.TrimSpace(shippingAddress.String)
		// chk_alamat_required: non-PICKUP wajib punya alamat_buyer_id atau snapshot.
		// Pesanan hasil migrasi tidak punya relasi alamat (keputusan #16), jadi
		// yang alamatnya kosong diturunkan ke PICKUP agar constraint terpenuhi.
		if deliveryType != "PICKUP" && alamat == "" {
			deliveryType = "PICKUP"
			catatanAdmin = strings.TrimSpace(catatanAdmin +
				" Migrasi v1: alamat pengiriman tidak tercatat, pesanan ditandai PICKUP.")
			a.rep.Add("6", "pesanan", id, "delivery_turun_pickup", "alamat kosong pada "+shippingMethod.String)
		}

		// Rincian biaya mengikuti aturan bisnis: diskon langsung memotong harga
		// barang, asuransi masuk biaya lainnya, ppn dan ongkir apa adanya.
		// biaya_produk diturunkan dari total (bukan dari SUM item) karena
		// total_price v1 terbukti otoritatif: 280 dari 280 pesanan PAID cocok
		// persis dengan jumlah invoice yang benar-benar dibayar, sementara
		// rincian item pada 109 pesanan tidak konsisten (diskon terhitung ganda).
		// Dengan cara ini total selalu utuh dan biaya_produk tidak pernah negatif.
		total := floatOrZero(totalPrice)
		ongkir := floatOrZero(shippingCost)
		ppn := floatOrZero(taxAmount)
		lainnya := floatOrZero(insuranceAmount)
		biayaProduk := total - ongkir - ppn - lainnya
		if biayaProduk < 0 {
			a.rep.Add("6", "pesanan", id, "biaya_produk_negatif",
				fmt.Sprintf("total=%.2f ongkir=%.2f ppn=%.2f lainnya=%.2f", total, ongkir, ppn, lainnya))
			biayaProduk = 0
		}

		// Timestamp status: v1 hanya menyimpan paid_off_at, sisanya diperkirakan
		// dari updated_at agar riwayat tetap punya penanda waktu yang masuk akal.
		stamp := nullTimePtr(updatedAt)
		var tProcessed, tReady, tShipped, tCompleted, tCancelled interface{}
		switch ordStatus {
		case "PROCESSING":
			tProcessed = stamp
		case "READY":
			tProcessed, tReady = stamp, stamp
		case "SHIPPED":
			tProcessed, tShipped = stamp, stamp
		case "COMPLETED":
			tProcessed, tCompleted = stamp, stamp
		case "CANCELLED":
			tCancelled = stamp
		}

		catatan := nullTrunc(notes, 1000)
		var catatanAdminPtr *string
		if catatanAdmin != "" {
			catatanAdminPtr = strPtr(strings.TrimSpace(catatanAdmin))
		}

		// forwarder_tracking_no: v1 punya dua sumber (orders.tracking_number dan
		// order_shippings.booking_id); ambil yang terisi.
		var tracking *string
		if t := firstNonEmpty(strings.TrimSpace(trackingNumber.String), strings.TrimSpace(bookingID.String)); t != "" {
			tracking = strPtr(truncate(t, 100))
		}
		var deliveree *string
		if deliveryType == "DELIVEREE" && strings.TrimSpace(bookingID.String) != "" {
			deliveree = strPtr(truncate(strings.TrimSpace(bookingID.String), 100))
			tracking = nil
		}

		// created_at memakai order_date bila ada: itulah waktu transaksi
		// sebenarnya di v1, sedangkan created_at bisa berbeda karena backfill.
		created := nullTimePtr(createdAt)
		if orderDate.Valid {
			created = nullTimePtr(orderDate)
		}

		if err := a.exec(`INSERT INTO pesanan
				(id, kode, buyer_id, delivery_type, alamat_buyer_id, payment_type,
				 payment_status, order_status, biaya_produk, biaya_pengiriman,
				 biaya_ppn, biaya_lainnya, total, catatan, catatan_admin,
				 expired_at, paid_at, processed_at, ready_at, shipped_at,
				 completed_at, cancelled_at, cancelled_reason,
				 deliveree_booking_id, forwarder_tracking_no,
				 alamat_snapshot, nama_penerima_snapshot, telepon_penerima_snapshot,
				 latitude_snapshot, longitude_snapshot,
				 created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?::delivery_type, NULL, 'REGULAR'::payment_type,
				 ?::payment_status, ?::order_status, ?, ?, ?, ?, ?, ?, ?,
				 ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			id, truncate(orderNumber, 20), userID.String, deliveryType,
			payStatus, ordStatus, biayaProduk, ongkir, ppn, lainnya, total,
			catatan, catatanAdminPtr,
			nullTimePtr(expiredAt), nullTimePtr(paidAt),
			tProcessed, tReady, tShipped, tCompleted, tCancelled,
			nullTrunc(cancelReason, 1000),
			deliveree, tracking,
			nullStrPtr(shippingAddress), nullTrunc(name, 100), nullTrunc(phone, 20),
			nullFloatPtr(lat), nullFloatPtr(lng),
			created, nullTimePtr(updatedAt), nullTimePtr(deletedAt)); err != nil {
			return err
		}

		a.pesananKnown[id] = true
		a.rep.Count("pesanan.insert")
		if trackingURL.Valid && strings.TrimSpace(trackingURL.String) != "" {
			a.rep.Count("pesanan.tracking_url_tidak_dimigrasi")
		}
	}
	return rows.Err()
}

// pesananItem memigrasi order_items -> pesanan_item.
//
// nama_produk NOT NULL di v2 tetapi tidak ada di order_items v1, jadi di-JOIN
// dari products (name, fallback name_trans.id untuk 372 produk yang name-nya
// kosong). Nilainya adalah nama produk saat migrasi, bukan saat transaksi —
// keterbatasan yang tak terhindarkan karena v1 tidak menyimpan snapshot.
// subtotal sengaja dikirim 0: trigger trg_calculate_subtotal berjalan
// BEFORE INSERT tanpa klausa WHEN sehingga nilainya selalu dihitung ulang.
func (a *App) pesananItem() error {
	a.pesananItemInserted = map[string]bool{}

	rows, err := a.my.Query(`SELECT oi.id, oi.order_id, oi.product_id, oi.quantity,
			oi.price, oi.discount_amount, oi.created_at, oi.updated_at,
			p.name, p.name_trans
		FROM order_items oi
		LEFT JOIN products p ON p.id = oi.product_id
		ORDER BY oi.created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var orderID, productID, pname, ptrans sql.NullString
		var qty sql.NullInt64
		var price, disc sql.NullFloat64
		var createdAt, updatedAt sql.NullTime

		if err := rows.Scan(&id, &orderID, &productID, &qty, &price, &disc,
			&createdAt, &updatedAt, &pname, &ptrans); err != nil {
			return err
		}

		if a.tgt.PesananItemIDs[id] {
			a.rep.Count("pesanan_item.sudah_ada_di_v2")
			continue
		}
		if !orderID.Valid || !a.orderExists[orderID.String] {
			// order induk sudah tidak ada di tabel orders v1 (data yatim)
			a.rep.Add("6", "pesanan_item", id, "skip_yatim", "order_id tidak ada di orders v1")
			continue
		}
		if !a.pesananKnown[orderID.String] {
			a.rep.Add("6", "pesanan_item", id, "skip", "pesanan induk tidak termigrasi")
			continue
		}
		if !productID.Valid || !a.produkKnown[productID.String] {
			a.rep.Add("6", "pesanan_item", id, "skip", "produk tidak termigrasi")
			continue
		}
		// CHECK (qty > 0)
		if qty.Int64 <= 0 {
			a.rep.Add("6", "pesanan_item", id, "skip", fmt.Sprintf("qty tidak valid: %d", qty.Int64))
			continue
		}

		namaID, _ := parseTrans(ptrans)
		nama := firstNonEmpty(strings.TrimSpace(pname.String), namaID)
		if nama == "" {
			nama = "Produk " + id8(productID.String)
			a.rep.Add("6", "pesanan_item", id, "nama_produk_fallback", "nama produk v1 kosong")
		}

		if err := a.exec(`INSERT INTO pesanan_item
				(id, pesanan_id, produk_id, nama_produk, sku, qty,
				 harga_satuan, diskon_satuan, subtotal, created_at, updated_at)
			VALUES (?, ?, ?, ?, NULL, ?, ?, ?, 0, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			id, orderID.String, productID.String, truncate(nama, 200),
			qty.Int64, floatOrZero(price), floatOrZero(disc),
			nullTimePtr(createdAt), nullTimePtr(updatedAt)); err != nil {
			return err
		}
		a.pesananItemInserted[id] = true
		a.rep.Count("pesanan_item.insert")
	}
	return rows.Err()
}
