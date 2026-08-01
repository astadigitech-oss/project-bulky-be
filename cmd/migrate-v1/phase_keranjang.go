package main

import (
	"database/sql"
	"fmt"
)

// phaseKeranjang memigrasi carts + cart_items v1 ke keranjang + keranjang_item v2.
//
// Bentuk datanya berbeda jauh: v1 membuat baris carts untuk hampir setiap
// pengunjung (39.735 baris, mayoritas kosong) dan membolehkan satu buyer punya
// beberapa cart, sedangkan v2 membatasi satu keranjang per buyer
// (UNIQUE buyer_id). Karena itu hanya keranjang yang benar-benar berisi yang
// dimigrasi, dan isi dari beberapa cart milik buyer yang sama digabung.
// Kolom v1 seperti coupon_code, shipping_cost, dan tax_amount tidak ikut:
// v2 menghitungnya ulang saat checkout, bukan menyimpannya di keranjang.
func (a *App) phaseKeranjang() error {
	if err := a.keranjang(); err != nil {
		return fmt.Errorf("keranjang: %w", err)
	}
	return a.keranjangItem()
}

func (a *App) keranjang() error {
	a.keranjangByBuyer = map[string]string{}

	// Hanya buyer dengan isi keranjang yang diambil; cart kosong tidak membawa
	// informasi apa pun. Cart terbaru dipakai sebagai pemilik id keranjang v2.
	rows, err := a.my.Query(`SELECT c.user_id, MAX(c.id) AS cart_id,
			MIN(c.created_at) AS dibuat, MAX(c.updated_at) AS diubah
		FROM carts c
		JOIN cart_items ci ON ci.cart_id = c.id AND ci.deleted_at IS NULL
		WHERE c.deleted_at IS NULL
		GROUP BY c.user_id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var userID, cartID string
		var dibuat, diubah sql.NullTime
		if err := rows.Scan(&userID, &cartID, &dibuat, &diubah); err != nil {
			return err
		}

		if id, ada := a.tgt.KeranjangByBuyer[userID]; ada {
			a.keranjangByBuyer[userID] = id
			a.rep.Count("keranjang.sudah_ada_di_v2")
			continue
		}
		if !a.buyerKnown[userID] && !a.tgt.BuyerIDs[userID] {
			a.rep.Add("9", "keranjang", cartID, "skip", "buyer tidak termigrasi: "+id8(userID))
			continue
		}

		if err := a.exec(`INSERT INTO keranjang (id, buyer_id, created_at, updated_at)
			VALUES (?, ?, ?, ?) ON CONFLICT (buyer_id) DO NOTHING`,
			cartID, userID, nullTimePtr(dibuat), nullTimePtr(diubah)); err != nil {
			return err
		}
		a.keranjangByBuyer[userID] = cartID
		a.rep.Count("keranjang.insert")
	}
	return rows.Err()
}

// keranjangItem menulis isi keranjang. Kuantitas dijumlahkan per pasangan
// (buyer, produk) karena UNIQUE (keranjang_id, produk_id) di v2 tidak
// mengizinkan baris kembar, sementara v1 punya 7 kasus produk yang sama muncul
// dua kali dalam satu cart, ditambah buyer yang isinya tersebar di lebih dari
// satu cart. harga tidak disalin: v2 selalu membaca harga terkini dari produk.
func (a *App) keranjangItem() error {
	rows, err := a.my.Query(`SELECT c.user_id, ci.product_id,
			SUM(ci.quantity) AS qty, MAX(ci.is_selected) AS dipilih,
			MIN(ci.created_at) AS dibuat, MAX(ci.updated_at) AS diubah
		FROM cart_items ci
		JOIN carts c ON c.id = ci.cart_id
		WHERE ci.deleted_at IS NULL AND c.deleted_at IS NULL
		GROUP BY c.user_id, ci.product_id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var userID, produkID string
		var qty, dipilih sql.NullInt64
		var dibuat, diubah sql.NullTime
		if err := rows.Scan(&userID, &produkID, &qty, &dipilih, &dibuat, &diubah); err != nil {
			return err
		}

		keranjangID, ada := a.keranjangByBuyer[userID]
		if !ada {
			continue // keranjang induk sudah dilaporkan saat di-skip
		}
		if !a.produkKnown[produkID] {
			a.rep.Add("9", "keranjang_item", produkID, "skip",
				"produk tidak termigrasi (buyer "+id8(userID)+")")
			continue
		}
		// CHECK (quantity > 0)
		if qty.Int64 <= 0 {
			a.rep.Add("9", "keranjang_item", produkID, "skip", "quantity tidak valid")
			continue
		}

		if err := a.exec(`INSERT INTO keranjang_item
				(id, keranjang_id, produk_id, quantity, is_selected, created_at, updated_at)
			VALUES (uuid_generate_v4(), ?, ?, ?, ?, ?, ?)
			ON CONFLICT (keranjang_id, produk_id) DO NOTHING`,
			keranjangID, produkID, qty.Int64, dipilih.Int64 == 1,
			nullTimePtr(dibuat), nullTimePtr(diubah)); err != nil {
			return err
		}
		a.rep.Count("keranjang_item.insert")
	}
	return rows.Err()
}
