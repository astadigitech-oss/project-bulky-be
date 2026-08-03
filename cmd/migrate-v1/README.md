# migrate-v1 — Migrasi Data Bulky v1 → v2 (Produk & Buyer)

Tool sekali-jalan (idempotent, aman diulang) untuk memindahkan data master **produk** dan **buyer** dari database Bulky v1 (Laravel/MySQL) ke database v2 (Postgres). Migrasi ini **hanya additive** — re-run menambah data yang belum ada, tidak menimpa yang sudah ada (lihat Catatan penting). Setelah migrasi, ada satu langkah pembersihan terpisah (cleanup produk LQD, lihat "Langkah 4").

Acuan lengkap mapping & keputusan: [docs/old DB/migrasi-data-v1-produk-buyer.md](../../docs/old%20DB/migrasi-data-v1-produk-buyer.md).

## Prasyarat

1. **Postgres v2** sudah menjalankan seluruh migrasi golang-migrate (s.d. 000157+), dan **sudah di-backup** (`pg_dump`).
2. **Dump v1 di-restore ke MySQL lokal**:
   ```
   mysql -u root -p -e "CREATE DATABASE bulky_v1"
   mysql -u root -p bulky_v1 < "docs/old DB/db_bulky_old.sql"
   ```
   Verifikasi restore: baris terakhir file dump tidak terpotong, dan
   `SELECT COUNT(*) FROM products; SELECT COUNT(*) FROM users;` masuk akal.
3. Env di `.env` (atau environment): `DB_*` (Postgres v2, sama dengan aplikasi) + `V1_DB_*` (MySQL v1, lihat `.env.example`).

## Menjalankan

```bash
# 1. Dry-run (default) — tidak menulis apa pun, hasilkan report
go run ./cmd/migrate-v1

# 2. Tinjau report
#    migrate-v1-report.json — semua dedup/fallback/skip/anomali tercatat di sini

# 3. Eksekusi sungguhan
go run ./cmd/migrate-v1 -execute

# opsi: -report <path>  (default migrate-v1-report.json)
```

Urutan fase: (1) master referensi → (2) produk + gambar + dokumen + pivot merek → (3) buyer + admin → (4) alamat buyer → (6) validasi. Tiap fase satu transaksi; gagal di tengah = rollback fase itu.

**Fase 5 (file fisik)** tidak ditangani tool ini: file gambar/PDF/foto profil menyusul dari file export server prod v1 (jendela maintenance). Zip folder `storage` Laravel v1, lalu upload lewat endpoint admin:

```
POST /api/panel/assets/import-v1   (Super Admin, multipart field "file")
```

ZIP v1 diharapkan berstruktur Laravel (`storage/app/public/...`) dan dipetakan ulang otomatis ke storage v2:

| Sumber di ZIP v1 | Tujuan di storage v2 |
|---|---|
| `storage/app/public/products/<file>` | `product-images/<file>` |
| `storage/app/public/<file>.pdf` | `product-documents/<file>` |
| `storage/app/public/reviews/<file>` | `reviews/<file>` |
| `storage/app/public/public/profile/<file>` | `profile/<file>` |

File yang tidak cocok mapping di atas di-skip dan dilaporkan di respons (`unmatched`). Endpoint `/import` (tanpa `-v1`) tetap untuk ZIP hasil `/export` v2 yang path-nya sudah relatif.

Respons `POST /api/panel/assets/import-v1`:

```json
{
  "success": true,
  "message": "Import v1 selesai",
  "data": {
    "imported": 4,          // file berhasil diekstrak & dipetakan
    "skipped": 1,           // file di-skip
    "unmatched": [],        // contoh path (maks 20) yang tidak dikenali
    "files": [              // laporan detail per-file
      { "source": "storage/app/public/products/a.jpg", "dest": "product-images/a.jpg", "status": "imported" },
      { "source": "storage/app/public/banners/b.jpg", "status": "skipped", "reason": "path tidak dikenali" }
    ]
  }
}
```

Status per file: `imported` | `skipped` (dengan `reason`: `path tidak dikenali`,
`path traversal ditolak`, `gagal buat folder`, `gagal baca zip`, `gagal buat file`,
`gagal tulis file`).

## Langkah 4: Cleanup produk LQD (wajib bila v2 sudah terisi)

Produk **LQD** ("Mix Product LQDSLE###" / "Palet LQDSLE###") adalah artefak percobaan
sync WMS yang gagal di v1 — ikut termigrasi tapi mayoritas tidak pernah terjual.
Filter di `phase_produk.go` hanya mencegah LQD **baru** masuk saat migrasi fresh;
produk yang **sudah terlanjur ada di v2** (termasuk bila re-run setelah filter
ditambahkan) **tidak dihapus oleh re-run** karena idempotensi me-*skip* baris yang ada.

Oleh karena itu jalankan skrip cleanup **setelah** migrasi selesai:

```bash
psql "postgresql://bulky_db:<password>@<host>:5432/bulky_db?sslmode=disable" \
  -f cmd/migrate-v1/cleanup_lqd.sql
```

Skrip ini **soft-delete** (set `deleted_at`) produk LQD yang tidak pernah terjual —
483 yang tidak pernah direferensikan `pesanan_item` + 37 yang hanya muncul di
pesanan EXPIRED/FAILED/CANCELLED = 520 produk — sambil **mempertahankan 106 produk
LQD yang pernah terjual** (order PAID+COMPLETED, `is_sold=true` mereka valid).
Soft-delete (bukan hard-delete) menjaga FK `pesanan_item` & riwayat pesanan tetap
utuh; rollback cukup `SET deleted_at = NULL`. Skrip idempotent (`deleted_at IS NULL`
guard) dan menampilkan laporan sebelum/sesudah di dalam transaksi.

## Catatan penting

- **Re-run bersifat additive, bukan sinkronisasi.** Semua INSERT memakai
  `ON CONFLICT (id) DO NOTHING` dan `LoadTargetState` me-*skip* ID yang sudah ada —
  re-run hanya menambah data yang belum ada dan **tidak pernah menimpa/memperbarui**
  data yang sudah dimigrasi. Perbaikan data di v1 tidak akan tersinkron ke v2 oleh
  re-run; kalau perlu perbaikan, lakukan lewat UPDATE/SQL manual di v2.
- `LoadTargetState` membaca `produk` tanpa filter `deleted_at`, jadi produk yang
  sudah di-soft-delete tetap "dikenal" dan tidak dimigrasi ulang (aman terhadap
  cleanup LQD).
- **Dump per 2026-07-19 hanya untuk latihan/dry-run.** Eksekusi final harus memakai dump segar yang diambil pada jendela maintenance yang sama dengan file storage (lihat dokumen §5.3).
- UUID v1 dipertahankan; baris yang sudah ada di target di-skip, sehingga tool aman dijalankan berulang.
- Admin hasil migrasi (`users.is_admin=1` + tabel `admins` v1) diberi role `ADMIN`; yang tanpa password ditandai di report dan perlu reset password manual.

## Catatan penting

- **Dump per 2026-07-19 hanya untuk latihan/dry-run.** Eksekusi final harus memakai dump segar yang diambil pada jendela maintenance yang sama dengan file storage (lihat dokumen §5.3).
- UUID v1 dipertahankan; baris yang sudah ada di target di-skip, sehingga tool aman dijalankan berulang.
- Admin hasil migrasi (`users.is_admin=1` + tabel `admins` v1) diberi role `ADMIN`; yang tanpa password ditandai di report dan perlu reset password manual.
