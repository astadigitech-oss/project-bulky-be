package main

import (
	"database/sql"
	"strings"
)

// phaseAlamat memigrasi addresses -> alamat_buyer (dokumen §7).
// Urutan insert penting: per buyer, alamat aktif dulu (primary paling awal) supaya
// trigger first-address-auto-default v2 menetapkan default yang sama dengan v1.
func (a *App) phaseAlamat() error {
	rows, err := a.my.Query(`SELECT a.id, a.user_id, a.label, a.name, a.phone_number, a.address,
			a.is_primary, a.latitude, a.longitude, a.deleted_at, a.created_at, a.updated_at,
			sd.name AS kelurahan, sd.postal_code, d.name AS kecamatan, c.name AS kota, p.name AS provinsi,
			u.is_admin
		FROM addresses a
		LEFT JOIN users u ON u.id = a.user_id
		LEFT JOIN sub_districts sd ON sd.id = a.sub_district_id
		LEFT JOIN districts d ON d.id = sd.district_id
		LEFT JOIN cities c ON c.id = d.city_id
		LEFT JOIN provinces p ON p.id = c.province_id
		ORDER BY a.user_id, (a.deleted_at IS NOT NULL), a.is_primary DESC, a.created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, userID string
		var label, name, phone, address sql.NullString
		var kelurahan, kodePos, kecamatan, kota, provinsi sql.NullString
		var isPrimary, isAdmin sql.NullInt64
		var lat, lng sql.NullFloat64
		var deletedAt, createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&id, &userID, &label, &name, &phone, &address,
			&isPrimary, &lat, &lng, &deletedAt, &createdAt, &updatedAt,
			&kelurahan, &kodePos, &kecamatan, &kota, &provinsi, &isAdmin); err != nil {
			return err
		}

		if a.tgt.AlamatIDs[id] {
			a.rep.Count("alamat_buyer.sudah_ada_di_v2")
			continue
		}
		if boolFromTiny(isAdmin) {
			a.rep.Add("4", "alamat_buyer", id, "skip", "pemilik adalah user is_admin=1 (masuk tabel admin, bukan buyer)")
			continue
		}
		if !a.buyerKnown[userID] {
			a.rep.Add("4", "alamat_buyer", id, "skip", "buyer '"+userID+"' tidak termigrasi")
			continue
		}

		provinsiV := wilayahText(provinsi)
		kotaV := wilayahText(kota)
		if provinsiV == "-" || kotaV == "-" {
			a.rep.Add("4", "alamat_buyer", id, "wilayah_kosong", "sub_district v1 NULL/orphan -> provinsi/kota '-'")
		}
		alamatLengkap := firstNonEmpty(address.String)
		if alamatLengkap == "" {
			alamatLengkap = strings.TrimSuffix(strings.Join(nonEmpty(
				wilayahOrEmpty(kelurahan), wilayahOrEmpty(kecamatan), kotaV, provinsiV), ", "), ", ")
			if alamatLengkap == "" {
				alamatLengkap = "-"
			}
			a.rep.Add("4", "alamat_buyer", id, "alamat_kosong", "address v1 NULL -> diisi teks wilayah")
		}

		teleponPenerima := "-"
		if p, ok := normalizePhone(phone.String); ok {
			teleponPenerima = p
		} else if v := strings.TrimSpace(phone.String); v != "" {
			teleponPenerima = truncate(v, 20)
		}

		isDefault := boolFromTiny(isPrimary) && !deletedAt.Valid

		if err := a.exec(`INSERT INTO alamat_buyer (id, buyer_id, label, nama_penerima, telepon_penerima,
				provinsi, kota, kecamatan, kelurahan, kode_pos, alamat_lengkap,
				latitude, longitude, is_default, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			id, userID,
			truncate(firstNonEmpty(label.String, "Alamat"), 50),
			truncate(firstNonEmpty(name.String, "-"), 100),
			teleponPenerima,
			truncate(provinsiV, 100), truncate(kotaV, 100),
			nullTrunc(kecamatan, 100), nullTrunc(kelurahan, 100), nullTrunc(kodePos, 10),
			alamatLengkap,
			nullFloatPtr(lat), nullFloatPtr(lng), isDefault,
			nullTimePtr(createdAt), nullTimePtr(updatedAt), nullTimePtr(deletedAt)); err != nil {
			return err
		}
		a.rep.Count("alamat_buyer.insert")
	}
	return rows.Err()
}

func wilayahText(ns sql.NullString) string {
	if v := strings.TrimSpace(ns.String); v != "" {
		return v
	}
	return "-"
}

func wilayahOrEmpty(ns sql.NullString) string {
	return strings.TrimSpace(ns.String)
}

func nonEmpty(vals ...string) []string {
	out := make([]string, 0, len(vals))
	for _, v := range vals {
		if v != "" && v != "-" {
			out = append(out, v)
		}
	}
	return out
}

func nullTrunc(ns sql.NullString, max int) *string {
	if v := strings.TrimSpace(ns.String); v != "" {
		return strPtr(truncate(v, max))
	}
	return nil
}

func nullFloatPtr(nf sql.NullFloat64) *float64 {
	if !nf.Valid {
		return nil
	}
	v := nf.Float64
	return &v
}
