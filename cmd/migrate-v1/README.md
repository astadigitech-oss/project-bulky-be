# migrate-v1 — Migrasi Data Bulky v1 → v2 (Master, Transaksi & Buyer)

Tool sekali-jalan (idempotent, aman diulang) untuk memindahkan data dari database Bulky v1 (Laravel/MySQL) ke database v2 (Postgres): master produk & referensi, buyer & admin, alamat, lalu **transaksi** (pesanan, pembayaran, kupon, ulasan, consent, keranjang). Migrasi ini **hanya additive** — re-run menambah data yang belum ada, tidak menimpa yang sudah ada (lihat Catatan penting). Setelah migrasi, ada satu langkah pembersihan terpisah (cleanup produk LQD, lihat "Langkah 4").

Acuan lengkap mapping & keputusan:
- Tahap 1 (master produk & buyer): [docs/old DB/migrasi-data-v1-produk-buyer.md](../../docs/old%20DB/migrasi-data-v1-produk-buyer.md) — keputusan #1–#12
- Tahap 2 (transaksi): [docs/old DB/migrasi-data-v1-transaksi.md](../../docs/old%20DB/migrasi-data-v1-transaksi.md) — keputusan #13–#24

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

Urutan fase (lihat `main.go`): (1) master referensi → (2) produk + gambar + dokumen + pivot merek → (3) buyer + admin → (4) alamat buyer → (5) master transaksi (disclaimer, metode bayar, settings) → (6) pesanan + item → (7) pembayaran (invoices) → (8) kupon, ulasan, consent → (9) keranjang → validasi. Tiap fase satu transaksi; gagal di tengah = rollback fase itu.

**File fisik (gambar/PDF/foto profil)** tidak ditangani tool ini: file menyusul dari file export server prod v1 (jendela maintenance). Zip folder `storage` Laravel v1, lalu upload lewat endpoint admin:

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

### ZIP backup v1 yang besar (>500MB) — chunk upload

Endpoint `/import-v1` (multipart sekali kirim) dibatasi `BodyLimit` Fiber 500MB,
jadi ZIP backup `storage` v1 yang bisa mencapai beberapa GB **tidak bisa diupload
sekali jalan**. Gunakan chunk upload (pola sama seperti upload video):

```
POST /api/panel/assets/import-v1/chunk      (multipart, per bagian)
  upload_id:     UUID yang sama untuk semua chunk
  chunk_index:   0-based
  total_chunks:  jumlah seluruh chunk
  chunk_data:    file bagian ZIP

POST /api/panel/assets/import-v1/finalize   (multipart)
  upload_id:     UUID yang sama
  total_chunks:  jumlah seluruh chunk
```

Alur di sisi client (FE):

1. Buat `upload_id` (mis. `crypto.randomUUID()`), potong ZIP menjadi N bagian
   (mis. 100–200MB per bagian agar aman di bawah BodyLimit).
2. Kirim tiap bagian ke `/import-v1/chunk` dengan `chunk_index` berurutan 0..N-1.
   Chunk disimpan ke `<UPLOAD_PATH>/chunks/v1-<upload_id>/chunk_XXXXX`.
3. Setelah semua terkirim, panggil `/import-v1/finalize`. Server menggabungkan
   chunk menjadi ZIP sementara, lalu memprosesnya **persis seperti** `/import-v1`
   (mapping tabel di atas). Folder chunk otomatis dihapus sesudahnya.

Respons `finalize` sama persis dengan `/import-v1` (`imported`/`skipped`/
`unmatched`/`files`). ZIP dibaca streaming dari disk — tidak dimuat penuh ke
memori, jadi aman untuk ukuran beberapa GB.

### File yatim (tidak direferensikan DB) — pruning

Saat migrasi DB, produk LQD yang tidak pernah terjual di-skip (filter `phase_produk.go`)
atau di-soft-delete (`cleanup_lqd.sql`). File gambar/PDF milik produk tersebut tetap
ikut dalam ZIP backup v1 dan akan ikut diekstrak — menjadi **file yatim** yang tidak
dipanggil aplikasi mana pun. Endpoint pruning menghapus file di storage v2 yang tidak
direferensikan kolom file DB v2:

```
POST /api/panel/assets/prune-orphans   (JSON)
  { "dry_run": true }    // default: hanya laporan, TIDAK menghapus
  { "dry_run": false }   // benar-benar hapus
```

Daftar referensi diambil dari `collectFilePaths` (tabel & kolom yang sama dengan
endpoint `/export`), mencakup `buyer.foto_url` — jadi foto profil buyer aman.
`logo_value` metode pembayaran berisi kode teks, bukan path file, jadi tidak ikut.
Folder `chunks/` (temp upload) dilewati.

Respons `dry_run=true`:

```json
{
  "success": true,
  "message": "Pruning selesai",
  "data": {
    "dry_run": true,
    "total_files": 520,
    "total_size": 4294967296,
    "deleted": 0,
    "orphans": [
      { "path": "product-images/lqd-sample.jpg", "size": 102400 }
    ]
  }
}
```

**Urutan aman:** migrasi DB → import ZIP v1 → jalankan `prune-orphans` dengan
`dry_run=true` dulu, tinjau daftar `orphans`, baru eksekusi `dry_run=false`.
Jangan dijalankan sebelum import aset selesai — file yang belum direferensikan DB
(hanya ada di ZIP) akan dianggap yatim dan terhapus permanen.

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
