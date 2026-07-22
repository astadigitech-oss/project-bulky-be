# migrate-v1 — Migrasi Data Bulky v1 → v2 (Produk & Buyer)

Tool sekali-jalan (idempotent, aman diulang) untuk memindahkan data master **produk** dan **buyer** dari database Bulky v1 (Laravel/MySQL) ke database v2 (Postgres).

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

**Fase 5 (file fisik)** tidak ditangani tool ini: file gambar/PDF/foto profil menyusul dari file export server prod v1 (jendela maintenance) dan tinggal ditaruh di `<UPLOAD_PATH>/product-images/`, `product-documents/`, `profile/` — path di DB sudah final.

## Catatan penting

- **Dump per 2026-07-19 hanya untuk latihan/dry-run.** Eksekusi final harus memakai dump segar yang diambil pada jendela maintenance yang sama dengan file storage (lihat dokumen §5.3).
- UUID v1 dipertahankan; baris yang sudah ada di target di-skip, sehingga tool aman dijalankan berulang.
- Admin hasil migrasi (`users.is_admin=1` + tabel `admins` v1) diberi role `ADMIN`; yang tanpa password ditandai di report dan perlu reset password manual.
