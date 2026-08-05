package main

import (
	"database/sql"
	"fmt"
	"strings"
)

// phaseMaster memigrasi/memetakan seluruh tabel referensi produk (dokumen §3).
func (a *App) phaseMaster() error {
	if err := a.masterKategori(); err != nil {
		return fmt.Errorf("kategori: %w", err)
	}
	if err := a.masterMerek(); err != nil {
		return fmt.Errorf("merek: %w", err)
	}
	if err := a.masterKondisi(); err != nil {
		return fmt.Errorf("kondisi: %w", err)
	}
	if err := a.masterKondisiPaket(); err != nil {
		return fmt.Errorf("kondisi_paket: %w", err)
	}
	if err := a.masterSumber(); err != nil {
		return fmt.Errorf("sumber: %w", err)
	}
	if err := a.masterWarehouse(); err != nil {
		return fmt.Errorf("warehouse: %w", err)
	}
	if err := a.masterFallbackRows(); err != nil {
		return fmt.Errorf("fallback: %w", err)
	}
	return a.masterTipe()
}

// §3.1 product_categories -> kategori_produk (mapping by ID; seed v2 memakai UUID v1)
func (a *App) masterKategori() error {
	a.kategoriKnown = map[string]bool{}
	for id := range a.tgt.KategoriIDs {
		a.kategoriKnown[id] = true
	}
	slugs := newSlugSpace(a.tgt.KategoriSlugs)
	slugsEN := newSlugSpace(a.tgt.KategoriSlugsEN)

	rows, err := a.my.Query(`SELECT id, name, name_trans, slug, icon, deleted_at, created_at, updated_at
		FROM product_categories ORDER BY created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var name, trans, slug, icon sql.NullString
		var deletedAt, createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&id, &name, &trans, &slug, &icon, &deletedAt, &createdAt, &updatedAt); err != nil {
			return err
		}
		if a.tgt.KategoriIDs[id] {
			a.rep.Count("kategori_produk.sudah_ada_di_v2")
			continue
		}
		tid, ten := parseTrans(trans)
		namaID := firstNonEmpty(tid, name.String)
		if namaID == "" {
			a.rep.Add("1", "kategori_produk", id, "skip", "nama kosong (junk)")
			continue
		}
		namaID = truncate(namaID, 100)
		// nama_en: fallback ke nama_id bila v1 tidak punya terjemahan (jangan pernah kosong)
		if ten == "" {
			a.rep.Add("1", "kategori_produk", id, "nama_en_fallback", "nama_en v1 kosong -> pakai nama_id")
		}
		namaEN := truncate(firstNonEmpty(ten, namaID), 100)
		base := slug.String
		if base == "" {
			base = slugify(namaID)
		}
		finalSlug := slugs.alloc(base, deletedAt.Valid, id)
		if finalSlug != base {
			a.rep.Add("1", "kategori_produk", id, "slug_diubah", base+" -> "+finalSlug)
		}
		// slug_en: selalu dialokasikan (dedup terpisah dari slug_id), jangan pernah kosong
		baseEN := slugify(namaEN)
		finalSlugEN := slugsEN.alloc(baseEN, deletedAt.Valid, id)
		if finalSlugEN != baseEN {
			a.rep.Add("1", "kategori_produk", id, "slug_en_diubah", baseEN+" -> "+finalSlugEN)
		}
		if err := a.exec(`INSERT INTO kategori_produk (id, nama_id, nama_en, slug, slug_id, slug_en, icon_url, is_active, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO NOTHING`,
			id, namaID, namaEN, finalSlug, finalSlug, finalSlugEN, nullStrPtr(icon), !deletedAt.Valid,
			nullTimePtr(createdAt), nullTimePtr(updatedAt), nullTimePtr(deletedAt)); err != nil {
			return err
		}
		a.kategoriKnown[id] = true
		a.rep.Count("kategori_produk.insert")
	}

	// fallback kategori (dokumen §3.8): baris "Lainnya"
	switch {
	case a.kategoriKnown[v1KategoriLainnyaID]:
		a.fallbackKategoriID = v1KategoriLainnyaID
	case a.tgt.KategoriBySlug["lainnya"] != "":
		a.fallbackKategoriID = a.tgt.KategoriBySlug["lainnya"]
	case a.tgt.KategoriBySlug["uncategorized"] != "":
		a.fallbackKategoriID = a.tgt.KategoriBySlug["uncategorized"]
	default:
		return fmt.Errorf("kategori fallback (Lainnya/uncategorized) tidak ditemukan di target")
	}
	return rows.Err()
}

// §3.2 product_brands -> merek_produk (mapping by ID)
func (a *App) masterMerek() error {
	a.merekKnown = map[string]bool{}
	for id := range a.tgt.MerekIDs {
		a.merekKnown[id] = true
	}
	slugs := newSlugSpace(a.tgt.MerekSlugs)
	slugsEN := newSlugSpace(a.tgt.MerekSlugsEN)

	rows, err := a.my.Query(`SELECT id, name, slug, deleted_at, created_at, updated_at
		FROM product_brands ORDER BY created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var name, slug sql.NullString
		var deletedAt, createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&id, &name, &slug, &deletedAt, &createdAt, &updatedAt); err != nil {
			return err
		}
		if a.tgt.MerekIDs[id] {
			a.rep.Count("merek_produk.sudah_ada_di_v2")
			continue
		}
		nama := truncate(strings.TrimSpace(name.String), 100)
		if nama == "" {
			a.rep.Add("1", "merek_produk", id, "skip", "nama kosong (junk)")
			continue
		}
		// v1 tidak punya terjemahan merek -> nama_en = nama_id (dokumen §3.2)
		base := firstNonEmpty(slug.String, slugify(nama))
		finalSlug := slugs.alloc(base, deletedAt.Valid, id)
		if finalSlug != base {
			a.rep.Add("1", "merek_produk", id, "slug_diubah", base+" -> "+finalSlug)
		}
		// slug_en: dedup terpisah dari slug_id, jangan pernah kosong
		finalSlugEN := slugsEN.alloc(base, deletedAt.Valid, id)
		if finalSlugEN != base {
			a.rep.Add("1", "merek_produk", id, "slug_en_diubah", base+" -> "+finalSlugEN)
		}
		if err := a.exec(`INSERT INTO merek_produk (id, nama_id, nama_en, slug, slug_id, slug_en, is_active, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO NOTHING`,
			id, nama, nama, finalSlug, finalSlug, finalSlugEN, !deletedAt.Valid,
			nullTimePtr(createdAt), nullTimePtr(updatedAt), nullTimePtr(deletedAt)); err != nil {
			return err
		}
		a.merekKnown[id] = true
		a.rep.Count("merek_produk.insert")
	}
	return rows.Err()
}

// §3.3 product_conditions -> kondisi_produk (mapping by slug; junk di-skip; sisanya insert dengan ID v1)
func (a *App) masterKondisi() error {
	a.condMap = map[string]string{}
	slugs := newSlugSpace(a.tgt.KondisiSlugs)
	slugsEN := newSlugSpace(a.tgt.KondisiSlugsEN)

	rows, err := a.my.Query(`SELECT id, title, title_trans, slug, deleted_at, created_at, updated_at
		FROM product_conditions ORDER BY created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var title, trans, slug sql.NullString
		var deletedAt, createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&id, &title, &trans, &slug, &deletedAt, &createdAt, &updatedAt); err != nil {
			return err
		}
		// junk v1 (keputusan #8): baris tes ber-slug kosong / '-1'
		if s := strings.TrimSpace(slug.String); s == "" || s == "-1" {
			a.rep.Add("1", "kondisi_produk", id, "skip", "junk (slug '"+s+"')")
			continue
		}
		if a.tgt.KondisiIDs[id] {
			a.condMap[id] = id
			a.rep.Count("kondisi_produk.sudah_ada_di_v2")
			continue
		}
		if tgtID := a.tgt.KondisiBySlug[slug.String]; tgtID != "" {
			a.condMap[id] = tgtID
			a.rep.Count("kondisi_produk.map_by_slug")
			continue
		}
		tid, ten := parseTrans(trans)
		namaID := truncate(firstNonEmpty(tid, title.String), 100)
		// nama_en: fallback ke nama_id bila v1 tidak punya terjemahan (jangan pernah kosong)
		if ten == "" {
			a.rep.Add("1", "kondisi_produk", id, "nama_en_fallback", "nama_en v1 kosong -> pakai nama_id")
		}
		namaEN := truncate(firstNonEmpty(ten, namaID), 100)
		finalSlug := slugs.alloc(slug.String, deletedAt.Valid, id)
		// slug_en: selalu dialokasikan (dedup terpisah dari slug_id), jangan pernah kosong
		baseEN := slugify(namaEN)
		finalSlugEN := slugsEN.alloc(baseEN, deletedAt.Valid, id)
		if finalSlugEN != baseEN {
			a.rep.Add("1", "kondisi_produk", id, "slug_en_diubah", baseEN+" -> "+finalSlugEN)
		}
		if err := a.exec(`INSERT INTO kondisi_produk (id, nama_id, nama_en, slug, slug_id, slug_en, urutan, is_active, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, 0, ?, ?, ?, ?) ON CONFLICT (id) DO NOTHING`,
			id, namaID, namaEN, finalSlug, finalSlug, finalSlugEN, !deletedAt.Valid,
			nullTimePtr(createdAt), nullTimePtr(updatedAt), nullTimePtr(deletedAt)); err != nil {
			return err
		}
		a.condMap[id] = id
		a.rep.Count("kondisi_produk.insert")
	}
	return rows.Err()
}

// §3.4 status_packages -> kondisi_paket (insert semua dengan ID v1)
func (a *App) masterKondisiPaket() error {
	a.pakMap = map[string]string{}
	slugs := newSlugSpace(a.tgt.PaketSlugs)
	slugsEN := newSlugSpace(a.tgt.PaketSlugsEN)

	rows, err := a.my.Query(`SELECT id, status, status_trans, deleted_at, created_at, updated_at
		FROM status_packages ORDER BY created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var status, trans sql.NullString
		var deletedAt, createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&id, &status, &trans, &deletedAt, &createdAt, &updatedAt); err != nil {
			return err
		}
		if a.tgt.PaketIDs[id] {
			a.pakMap[id] = id
			a.rep.Count("kondisi_paket.sudah_ada_di_v2")
			continue
		}
		tid, ten := parseTrans(trans)
		namaID := truncate(firstNonEmpty(tid, status.String), 100)
		if namaID == "" {
			a.rep.Add("1", "kondisi_paket", id, "skip", "nama kosong (junk)")
			continue
		}
		// nama_en: fallback ke nama_id bila v1 tidak punya terjemahan (jangan pernah kosong)
		if ten == "" {
			a.rep.Add("1", "kondisi_paket", id, "nama_en_fallback", "nama_en v1 kosong -> pakai nama_id")
		}
		namaEN := truncate(firstNonEmpty(ten, namaID), 100)
		finalSlug := slugs.alloc(slugify(namaID), deletedAt.Valid, id)
		// slug_en: selalu dialokasikan (dedup terpisah dari slug_id), jangan pernah kosong
		baseEN := slugify(namaEN)
		finalSlugEN := slugsEN.alloc(baseEN, deletedAt.Valid, id)
		if finalSlugEN != baseEN {
			a.rep.Add("1", "kondisi_paket", id, "slug_en_diubah", baseEN+" -> "+finalSlugEN)
		}
		if err := a.exec(`INSERT INTO kondisi_paket (id, nama_id, nama_en, slug, slug_id, slug_en, urutan, is_active, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, 0, ?, ?, ?, ?) ON CONFLICT (id) DO NOTHING`,
			id, namaID, namaEN, finalSlug, finalSlug, finalSlugEN, !deletedAt.Valid,
			nullTimePtr(createdAt), nullTimePtr(updatedAt), nullTimePtr(deletedAt)); err != nil {
			return err
		}
		a.pakMap[id] = id
		a.rep.Count("kondisi_paket.insert")
	}
	return rows.Err()
}

// §3.5 product_statuses -> sumber_produk (mapping semantik manual; unmapped -> NULL)
func (a *App) masterSumber() error {
	a.sumMap = map[string]string{}
	semantik := map[string]string{ // nama v1 -> slug sumber v2
		"overstock produsen":    "overstock",
		"failed delivery items": "retur",
	}

	rows, err := a.my.Query(`SELECT id, status FROM product_statuses`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var status sql.NullString
		if err := rows.Scan(&id, &status); err != nil {
			return err
		}
		slug, ok := semantik[strings.ToLower(strings.TrimSpace(status.String))]
		if !ok {
			a.rep.Add("1", "sumber_produk", id, "unmapped", "status v1 '"+status.String+"' -> sumber NULL")
			continue
		}
		tgtID := a.tgt.SumberBySlug[slug]
		if tgtID == "" {
			return fmt.Errorf("sumber_produk slug '%s' tidak ada di target — pastikan seed 000064 terpasang", slug)
		}
		a.sumMap[id] = tgtID
		a.rep.Count("sumber_produk.map")
	}
	return rows.Err()
}

// §3.6 warehouses -> warehouse (mapping by slug; junk terhapus di-skip; backfill koordinat)
func (a *App) masterWarehouse() error {
	a.whMap = map[string]string{}
	cibinong := a.tgt.WarehouseBySlug["warehouse-cibinong"]
	if cibinong == "" {
		return fmt.Errorf("warehouse 'warehouse-cibinong' tidak ada di target — pastikan seed 000091 terpasang")
	}
	a.fallbackWarehouseID = cibinong

	rows, err := a.my.Query(`SELECT id, name, latitude, longitude, deleted_at FROM warehouses`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var name sql.NullString
		var lat, lng sql.NullFloat64
		var deletedAt sql.NullTime
		if err := rows.Scan(&id, &name, &lat, &lng, &deletedAt); err != nil {
			return err
		}
		if deletedAt.Valid {
			a.rep.Add("1", "warehouse", id, "skip", "warehouse terhapus/junk '"+name.String+"' -> produk terkait jatuh ke Cibinong")
			continue
		}
		tgtID := a.tgt.WarehouseBySlug[slugify(name.String)]
		if tgtID == "" {
			a.rep.Add("1", "warehouse", id, "fallback", "'"+name.String+"' tidak cocok slug manapun di v2 -> map ke Cibinong")
			tgtID = cibinong
		}
		a.whMap[id] = tgtID
		a.rep.Count("warehouse.map")

		// backfill koordinat Cibinong dari v1 bila v2 masih kosong
		if tgtID == cibinong && lat.Valid && lng.Valid {
			if err := a.exec(`UPDATE warehouse SET latitude = ?, longitude = ? WHERE id = ?::uuid AND latitude IS NULL`,
				lat.Float64, lng.Float64, tgtID); err != nil {
				return err
			}
			a.rep.Count("warehouse.backfill_koordinat")
		}
	}
	return rows.Err()
}

// §3.8 baris fallback untuk FK NOT NULL yang di v1 NULL
func (a *App) masterFallbackRows() error {
	kondisiSlugs := newSlugSpace(a.tgt.KondisiSlugs)
	kondisiSlugsEN := newSlugSpace(a.tgt.KondisiSlugsEN)
	paketSlugs := newSlugSpace(a.tgt.PaketSlugs)
	paketSlugsEN := newSlugSpace(a.tgt.PaketSlugsEN)
	if !a.tgt.KondisiIDs[fallbackKondisiID] {
		slug := kondisiSlugs.alloc("tidak-diketahui", false, fallbackKondisiID)
		slugEN := kondisiSlugsEN.alloc("unknown", false, fallbackKondisiID)
		if err := a.exec(`INSERT INTO kondisi_produk (id, nama_id, nama_en, slug, slug_id, slug_en, urutan, is_active)
			VALUES (?, 'Tidak Diketahui', 'Unknown', ?, ?, ?, 99, false) ON CONFLICT (id) DO NOTHING`,
			fallbackKondisiID, slug, slug, slugEN); err != nil {
			return err
		}
	}
	if !a.tgt.PaketIDs[fallbackKondisiPaketID] {
		slug := paketSlugs.alloc("tidak-diketahui", false, fallbackKondisiPaketID)
		slugEN := paketSlugsEN.alloc("unknown", false, fallbackKondisiPaketID)
		if err := a.exec(`INSERT INTO kondisi_paket (id, nama_id, nama_en, slug, slug_id, slug_en, urutan, is_active)
			VALUES (?, 'Tidak Diketahui', 'Unknown', ?, ?, ?, 99, false) ON CONFLICT (id) DO NOTHING`,
			fallbackKondisiPaketID, slug, slug, slugEN); err != nil {
			return err
		}
	}
	return nil
}

// §3.7 packaging_type -> tipe_produk
func (a *App) masterTipe() error {
	a.tipeMap = map[string]string{}
	for v1Enum, v2Slug := range map[string]string{
		"palet":      "palet-load",
		"truck_load": "truck-load",
		"container":  "container-load",
	} {
		id := a.tgt.TipeBySlug[v2Slug]
		if id == "" {
			return fmt.Errorf("tipe_produk slug '%s' tidak ada di target — pastikan seed 000008/000053 terpasang", v2Slug)
		}
		a.tipeMap[v1Enum] = id
	}
	return nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
