package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

// phaseProduk memigrasi products -> produk (+gambar, dokumen, pivot merek). Dokumen §4-§5.
func (a *App) phaseProduk() error {
	for id := range a.tgt.ProdukIDs {
		a.produkKnown[id] = true
	}
	slugs := newSlugSpace(a.tgt.ProdukSlugs)
	slugsEN := newSlugSpace(a.tgt.ProdukSlugEN)
	cargoUsed := map[string]bool{}
	for c := range a.tgt.IDCargo {
		cargoUsed[c] = true
	}

	rows, err := a.my.Query(`SELECT id, wms_id, images, name, name_trans, slug, id_pallet,
			price, price_before_discount, total_quantity, packaging_type, pdf_file, is_active,
			warehouse_id, product_category_id, product_condition_id, product_status_id, status_package_id,
			note_discrepancy, length_cm, width_cm, height_cm, weight_kg,
			deleted_at, created_at, updated_at, sold_out
		FROM products ORDER BY created_at, id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, packagingType string
		var name, trans, slug, idPallet, images, pdfFile sql.NullString
		var whID, catID, condID, statusID, pakID sql.NullString
		var wmsID, price, priceBefore, totalQty, noteDisc, isActive, soldOut sql.NullInt64
		var length, width, height, weight sql.NullFloat64
		var deletedAt, createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&id, &wmsID, &images, &name, &trans, &slug, &idPallet,
			&price, &priceBefore, &totalQty, &packagingType, &pdfFile, &isActive,
			&whID, &catID, &condID, &statusID, &pakID,
			&noteDisc, &length, &width, &height, &weight,
			&deletedAt, &createdAt, &updatedAt, &soldOut); err != nil {
			return err
		}

		if a.tgt.ProdukIDs[id] {
			a.rep.Count("produk.sudah_ada_di_v2")
			continue
		}
		deleted := deletedAt.Valid

		// nama (dokumen §4)
		tid, ten := parseTrans(trans)
		namaID := truncate(firstNonEmpty(tid, name.String, "Produk "+id8(id)), 255)
		namaEN := truncate(firstNonEmpty(ten, namaID), 255)

		// slug + slug_id (+slug_en)
		base := firstNonEmpty(slug.String, slugify(namaID))
		finalSlug := slugs.alloc(base, deleted, id)
		if finalSlug != base {
			a.rep.Add("2", "produk", id, "slug_diubah", base+" -> "+finalSlug)
		}
		// slug_en: selalu dialokasikan (dedup terpisah dari slug_id), jangan pernah kosong
		baseEN := slugify(namaEN)
		finalSlugEN := slugsEN.alloc(baseEN, deleted, id)
		if finalSlugEN != baseEN {
			a.rep.Add("2", "produk", id, "slug_en_diubah", baseEN+" -> "+finalSlugEN)
		}
		slugEN := &finalSlugEN

		// id_cargo (unik di v2; v1 id_pallet tidak unik)
		var idCargo *string
		if v := firstNonEmpty(idPallet.String); v != "" {
			v = truncate(v, 50)
			if cargoUsed[v] {
				a.rep.Add("2", "produk", id, "id_cargo_duplikat", "'"+v+"' sudah dipakai baris lebih lama -> NULL")
			} else {
				cargoUsed[v] = true
				idCargo = &v
			}
		}

		var refID *string
		if wmsID.Valid {
			refID = strPtr(strconv.FormatInt(wmsID.Int64, 10))
		}

		// FK dengan fallback (dokumen §3.8)
		kategoriID := a.fallbackKategoriID
		if catID.Valid && a.kategoriKnown[catID.String] {
			kategoriID = catID.String
		} else {
			a.rep.Add("2", "produk", id, "fallback_kategori", "kategori v1 '"+catID.String+"' -> Lainnya")
		}
		kondisiID := fallbackKondisiID
		if condID.Valid {
			if m := a.condMap[condID.String]; m != "" {
				kondisiID = m
			} else {
				a.rep.Add("2", "produk", id, "fallback_kondisi", "kondisi v1 '"+condID.String+"' -> Tidak Diketahui")
			}
		} else {
			a.rep.Add("2", "produk", id, "fallback_kondisi", "kondisi v1 NULL -> Tidak Diketahui")
		}
		paketID := fallbackKondisiPaketID
		if pakID.Valid {
			if m := a.pakMap[pakID.String]; m != "" {
				paketID = m
			} else {
				a.rep.Add("2", "produk", id, "fallback_kondisi_paket", "status_package v1 '"+pakID.String+"' -> Tidak Diketahui")
			}
		} else {
			a.rep.Add("2", "produk", id, "fallback_kondisi_paket", "status_package v1 NULL -> Tidak Diketahui")
		}
		var sumberID *string
		if statusID.Valid {
			if m := a.sumMap[statusID.String]; m != "" {
				sumberID = &m
			}
		}
		warehouseID := a.fallbackWarehouseID
		if whID.Valid && a.whMap[whID.String] != "" {
			warehouseID = a.whMap[whID.String]
		} else if whID.Valid {
			a.rep.Add("2", "produk", id, "fallback_warehouse", "warehouse v1 '"+whID.String+"' -> Cibinong")
		} else {
			a.rep.Add("2", "produk", id, "fallback_warehouse", "warehouse v1 NULL -> Cibinong")
		}
		tipeID := a.tipeMap[packagingType]
		if tipeID == "" {
			return fmt.Errorf("produk %s: packaging_type '%s' tidak dikenal", id, packagingType)
		}

		// harga (dokumen §4: harga_sebelum_diskon wajib > 0)
		hargaSesudah := int64(0)
		if price.Valid {
			hargaSesudah = price.Int64
		} else if priceBefore.Valid {
			hargaSesudah = priceBefore.Int64
		}
		hargaSebelum := priceBefore.Int64
		if hargaSebelum <= 0 {
			if price.Valid && price.Int64 > 0 {
				hargaSebelum = price.Int64
			} else {
				hargaSebelum = 1
				a.rep.Add("2", "produk", id, "harga_anomali", "price & price_before_discount kosong -> harga_sebelum_diskon=1")
			}
		}

		qty := totalQty.Int64
		if qty < 0 {
			a.rep.Add("2", "produk", id, "quantity_negatif", fmt.Sprintf("%d -> 0", qty))
			qty = 0
		}
		discPct := noteDisc.Int64
		if discPct > 100 {
			a.rep.Add("2", "produk", id, "discrepancy_clamp", fmt.Sprintf("%d -> 100", discPct))
			discPct = 100
		}

		if err := a.exec(`INSERT INTO produk (id, nama_id, nama_en, slug, slug_id, slug_en, id_cargo, reference_id,
				kategori_id, kondisi_id, kondisi_paket_id, sumber_id, warehouse_id, tipe_produk_id,
				harga_sebelum_diskon, harga_sesudah_diskon, quantity, is_sold, discrepancy_percentage,
				panjang, lebar, tinggi, berat, is_active, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			id, namaID, namaEN, finalSlug, finalSlug, slugEN, idCargo, refID,
			kategoriID, kondisiID, paketID, sumberID, warehouseID, tipeID,
			hargaSebelum, hargaSesudah, qty, boolFromTiny(soldOut), discPct,
			floatOrZero(length), floatOrZero(width), floatOrZero(height), floatOrZero(weight),
			boolFromTiny(isActive), nullTimePtr(createdAt), nullTimePtr(updatedAt), nullTimePtr(deletedAt)); err != nil {
			return err
		}
		a.produkKnown[id] = true
		a.rep.Count("produk.insert")

		// gambar (dokumen §5.1)
		files, perr := parseImages(images)
		if perr != nil {
			a.rep.Add("2", "produk_gambar", id, "images_rusak", "JSON images tidak bisa diparse: "+truncate(images.String, 80))
		}
		for i, f := range files {
			if err := a.exec(`INSERT INTO produk_gambar (produk_id, gambar_url, urutan, is_primary, created_at)
				VALUES (?, ?, ?, ?, COALESCE(?, NOW()))`,
				id, "product-images/"+f, i, i == 0, nullTimePtr(createdAt)); err != nil {
				return err
			}
			a.rep.Count("produk_gambar.insert")
		}
		if len(files) == 0 && perr == nil {
			a.rep.Add("2", "produk_gambar", id, "tanpa_gambar", "")
		}

		// dokumen PDF (dokumen §5.2)
		if f := firstNonEmpty(pdfFile.String); f != "" {
			if err := a.exec(`INSERT INTO produk_dokumen (produk_id, nama_dokumen, file_url, tipe_file, created_at)
				VALUES (?, ?, ?, 'pdf', COALESCE(?, NOW()))`,
				id, truncate(namaID, 255), "product-documents/"+f, nullTimePtr(createdAt)); err != nil {
				return err
			}
			a.rep.Count("produk_dokumen.insert")
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return a.produkMerekPivot()
}

// §5.4 product_brand_pivot -> produk_merek
func (a *App) produkMerekPivot() error {
	rows, err := a.my.Query(`SELECT product_id, product_brand_id FROM product_brand_pivot`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var produkID, merekID string
		if err := rows.Scan(&produkID, &merekID); err != nil {
			return err
		}
		if !a.produkKnown[produkID] || !a.merekKnown[merekID] {
			a.rep.Add("2", "produk_merek", produkID, "skip", "produk/merek tidak termigrasi (merek: "+merekID+")")
			continue
		}
		if err := a.exec(`INSERT INTO produk_merek (produk_id, merek_id) VALUES (?, ?)
			ON CONFLICT (produk_id, merek_id) DO NOTHING`, produkID, merekID); err != nil {
			return err
		}
		a.rep.Count("produk_merek.insert")
	}
	return rows.Err()
}
