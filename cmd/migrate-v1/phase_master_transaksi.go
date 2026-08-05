package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// phaseMasterTransaksi memigrasi master data yang dibutuhkan fase transaksi
// (dokumen §5.4): disclaimer, lookup metode pembayaran, dan audit settings.
func (a *App) phaseMasterTransaksi() error {
	if err := a.masterDisclaimer(); err != nil {
		return fmt.Errorf("disclaimer: %w", err)
	}
	if err := a.loadPaymentMethodCodes(); err != nil {
		return fmt.Errorf("payment_methods: %w", err)
	}
	return a.auditSettings()
}

// masterDisclaimer memigrasi disclaimers v1 -> disclaimer v2.
//
// Dua hal yang menentukan bentuk kode di bawah:
//  1. judul_en dan konten_en NOT NULL di v2, sementara v1 menyimpannya di kolom
//     JSON title_trans/content_trans. Bila terjemahan tidak ada, fallback ke ID.
//  2. Trigger fn_ensure_single_active_disclaimer menonaktifkan baris lain setiap
//     kali ada baris aktif masuk. Baris diurutkan is_active ASC supaya yang aktif
//     ter-insert paling akhir dan status akhirnya benar.
func (a *App) masterDisclaimer() error {
	a.disclaimerKnown = map[string]bool{}
	a.disclaimerBySlug = map[string]string{}

	rows, err := a.my.Query(`SELECT id, title, title_trans, slug, content, content_trans,
			is_active, deleted_at, created_at, updated_at
		FROM disclaimers ORDER BY is_active ASC, created_at ASC`)
	if err != nil {
		return err
	}
	defer rows.Close()

	slugs := newSlugSpace(a.tgt.DisclaimerSlugs)

	for rows.Next() {
		var id string
		var title, titleTrans, slug, content, contentTrans sql.NullString
		var isActive sql.NullInt64
		var deletedAt, createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&id, &title, &titleTrans, &slug, &content, &contentTrans,
			&isActive, &deletedAt, &createdAt, &updatedAt); err != nil {
			return err
		}

		if a.tgt.DisclaimerIDs[id] {
			a.disclaimerKnown[id] = true
			a.rep.Count("disclaimer.sudah_ada_di_v2")
			continue
		}

		judulID, judulEN := parseTrans(titleTrans)
		judulID = firstNonEmpty(judulID, title.String)
		if judulID == "" {
			a.rep.Add("5", "disclaimer", id, "skip", "judul kosong")
			continue
		}
		if judulEN == "" {
			a.rep.Add("5", "disclaimer", id, "judul_en_fallback", "title_trans.en kosong -> pakai judul id")
		}
		judulEN = firstNonEmpty(judulEN, judulID)

		kontenID, kontenEN := parseTrans(contentTrans)
		kontenID = firstNonEmpty(kontenID, content.String)
		if kontenID == "" {
			a.rep.Add("5", "disclaimer", id, "skip", "konten kosong")
			continue
		}
		if kontenEN == "" {
			a.rep.Add("5", "disclaimer", id, "konten_en_fallback", "content_trans.en kosong -> pakai konten id")
		}
		kontenEN = firstNonEmpty(kontenEN, kontenID)

		base := firstNonEmpty(slug.String, slugify(judulID))
		finalSlug := slugs.alloc(base, deletedAt.Valid, id)
		if finalSlug != base {
			a.rep.Add("5", "disclaimer", id, "slug_diubah", base+" -> "+finalSlug)
		}
		slugEN := slugs.alloc(slugify(judulEN), deletedAt.Valid, id)

		// Baris terhapus di v1 tidak boleh ikut aktif di v2.
		aktif := boolFromTiny(isActive) && !deletedAt.Valid

		if err := a.exec(`INSERT INTO disclaimer
				(id, judul, judul_en, slug, slug_id, slug_en, konten, konten_en,
				 is_active, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			id, truncate(judulID, 200), truncate(judulEN, 200),
			finalSlug, finalSlug, slugEN,
			kontenID, kontenEN, aktif,
			nullTimePtr(createdAt), nullTimePtr(updatedAt), nullTimePtr(deletedAt)); err != nil {
			return err
		}
		a.disclaimerKnown[id] = true
		a.disclaimerBySlug[finalSlug] = id
		a.rep.Count("disclaimer.insert")
	}
	return rows.Err()
}

// loadPaymentMethodCodes membangun peta payment_methods.id -> code untuk dipakai
// Fase 7 mengisi pesanan_pembayaran.xendit_payment_method. Master metode
// pembayaran v1 sendiri tidak dimigrasi: v2 memakai tabel metode_pembayaran
// dengan ID sendiri yang sudah di-seed, sehingga hanya kodenya yang relevan.
func (a *App) loadPaymentMethodCodes() error {
	a.paymentMethodCode = map[string]string{}

	rows, err := a.my.Query(`SELECT id, code FROM payment_methods`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var code sql.NullString
		if err := rows.Scan(&id, &code); err != nil {
			return err
		}
		c := strings.ToUpper(strings.TrimSpace(code.String))
		if c == "" {
			continue
		}
		a.paymentMethodCode[id] = truncate(c, 50)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	a.rep.Count(fmt.Sprintf("payment_methods.kode_dimuat=%d", len(a.paymentMethodCode)))
	return nil
}

// auditSettings mencatat isi tabel settings v1 ke report. Tidak ada yang ditulis
// ke v2: tarif pajak sudah identik dengan seed v2 (11%, aktif), sedangkan sisanya
// (kontak WhatsApp, info pickup, email form wholesale) adalah konfigurasi yang di
// v2 diatur lewat env/admin panel. Report dipakai operator untuk menyalinnya manual.
func (a *App) auditSettings() error {
	rows, err := a.my.Query("SELECT `group`, `name`, payload FROM settings ORDER BY `group`, `name`")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var group, name, payload sql.NullString
		if err := rows.Scan(&group, &name, &payload); err != nil {
			return err
		}
		key := group.String + "." + name.String

		// payload disimpan sebagai JSON-encoded value; unwrap agar report terbaca.
		val := strings.TrimSpace(payload.String)
		var decoded interface{}
		if json.Unmarshal([]byte(val), &decoded) == nil {
			if s, ok := decoded.(string); ok {
				val = s
			}
		}

		if key == "tax.rate" || key == "tax.enabled" {
			a.rep.Add("5", "settings", key, "sudah_sesuai_seed_v2", truncate(val, 200))
			continue
		}
		a.rep.Add("5", "settings", key, "manual_admin_panel", truncate(val, 200))
	}
	return rows.Err()
}
