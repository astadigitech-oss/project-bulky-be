package main

import (
	"database/sql"
	"fmt"
	"strings"
)

// phasePembayaran memigrasi invoices v1 -> pesanan_pembayaran v2 (dokumen §5.3),
// lalu menandai pesanan yang dibayar lebih dari satu kali sebagai SPLIT.
func (a *App) phasePembayaran() error {
	if err := a.pembayaran(); err != nil {
		return fmt.Errorf("pesanan_pembayaran: %w", err)
	}
	return a.tandaiSplitPayment()
}

// mapInvoiceStatus memetakan invoices.status v1 ke enum payment_status v2.
// v1 memakai huruf kecil dan mengenal 'cancelled' yang tidak ada padanannya
// di v2; sama seperti keputusan #17 pada header pesanan, statusnya diturunkan
// ke EXPIRED karena artinya pembayaran tidak pernah diselesaikan.
func mapInvoiceStatus(v1 string) string {
	switch strings.ToLower(strings.TrimSpace(v1)) {
	case "paid", "settled", "success":
		return "PAID"
	case "expired":
		return "EXPIRED"
	case "cancelled", "canceled":
		return "EXPIRED"
	case "failed":
		return "FAILED"
	case "refunded":
		return "REFUNDED"
	default:
		return "PENDING"
	}
}

func (a *App) pembayaran() error {
	rows, err := a.my.Query(`SELECT i.id, i.user_id, i.order_id, i.payment_method_id,
			i.amount, i.status, i.xendit_id, i.xendit_invoice_url,
			i.midtrans_snap_token, i.midtrans_redirect_url,
			i.deleted_at, i.created_at, i.updated_at,
			o.payment_expired_at, o.paid_off_at
		FROM invoices i
		JOIN orders o ON o.id = i.order_id
		ORDER BY i.created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, orderID string
		var userID, methodID, status sql.NullString
		var xenditID, xenditURL, snapToken, redirectURL sql.NullString
		var amount sql.NullFloat64
		var deletedAt, createdAt, updatedAt, expiredAt, paidAt sql.NullTime

		if err := rows.Scan(&id, &userID, &orderID, &methodID, &amount, &status,
			&xenditID, &xenditURL, &snapToken, &redirectURL,
			&deletedAt, &createdAt, &updatedAt, &expiredAt, &paidAt); err != nil {
			return err
		}

		if a.tgt.PembayaranIDs[id] {
			a.rep.Count("pembayaran.sudah_ada_di_v2")
			continue
		}
		if !a.pesananKnown[orderID] {
			a.rep.Add("7", "pesanan_pembayaran", id, "skip", "pesanan induk tidak termigrasi")
			continue
		}
		if !userID.Valid || userID.String == "" {
			a.rep.Add("7", "pesanan_pembayaran", id, "skip", "user_id NULL di v1")
			continue
		}

		st := mapInvoiceStatus(status.String)

		// paid_at hanya bermakna untuk pembayaran yang benar-benar lunas.
		// v1 tidak menyimpannya per-invoice, jadi dipakai orders.paid_off_at.
		var tPaid interface{}
		if st == "PAID" {
			if paidAt.Valid {
				tPaid = nullTimePtr(paidAt)
			} else {
				tPaid = nullTimePtr(updatedAt)
			}
		}

		// xendit_external_id diisi id invoice v1: itulah kunci yang dikirim ke
		// Xendit sebagai external_id, sehingga rekonsiliasi lama tetap bisa
		// ditelusuri. xendit_invoice_id diisi xendit_id bila ada.
		var payURL *string
		if u := firstNonEmpty(strings.TrimSpace(xenditURL.String), strings.TrimSpace(redirectURL.String)); u != "" {
			payURL = strPtr(u)
		}

		// 6 invoice hanya punya jejak Midtrans (payment gateway lama). Datanya
		// tetap dimigrasi agar riwayat pembayaran buyer utuh, tapi ditandai.
		if !xenditID.Valid && (snapToken.Valid || redirectURL.Valid) {
			a.rep.Add("7", "pesanan_pembayaran", id, "gateway_midtrans_lama",
				"tidak ada xendit_id, dimigrasi apa adanya")
		}

		var metode *string
		if methodID.Valid {
			if code, ok := a.paymentMethodCode[methodID.String]; ok && code != "" {
				metode = strPtr(code)
			}
		}

		if err := a.exec(`INSERT INTO pesanan_pembayaran
				(id, pesanan_id, buyer_id, metode_pembayaran_id, jumlah, status,
				 xendit_invoice_id, xendit_external_id, xendit_payment_url,
				 xendit_payment_method, expired_at, paid_at, created_at, updated_at)
			VALUES (?, ?, ?, NULL, ?, ?::payment_status, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			id, orderID, userID.String, floatOrZero(amount), st,
			nullTrunc(xenditID, 100), truncate(id, 100), payURL, metode,
			nullTimePtr(expiredAt), tPaid,
			nullTimePtr(createdAt), nullTimePtr(updatedAt)); err != nil {
			return err
		}
		a.rep.Count("pembayaran.insert")
		if deletedAt.Valid {
			a.rep.Count("pembayaran.v1_soft_deleted_tetap_dimigrasi")
		}
	}
	return rows.Err()
}

// tandaiSplitPayment menyetel payment_type = SPLIT untuk pesanan yang punya
// lebih dari satu pembayaran lunas. Kolom itu tidak bisa ditentukan saat Fase 6
// karena jumlah invoice baru diketahui setelah tabel pembayaran terisi.
func (a *App) tandaiSplitPayment() error {
	if !a.execute {
		a.rep.Count("pembayaran.split_ditandai_saat_execute")
		return nil
	}
	res := a.tx.Exec(`UPDATE pesanan SET payment_type = 'SPLIT'
		WHERE payment_type <> 'SPLIT' AND id IN (
			SELECT pesanan_id FROM pesanan_pembayaran
			WHERE status = 'PAID' GROUP BY pesanan_id HAVING COUNT(*) > 1)`)
	if res.Error != nil {
		return fmt.Errorf("tandai split: %w", res.Error)
	}
	for i := int64(0); i < res.RowsAffected; i++ {
		a.rep.Count("pembayaran.pesanan_jadi_split")
	}
	return nil
}
