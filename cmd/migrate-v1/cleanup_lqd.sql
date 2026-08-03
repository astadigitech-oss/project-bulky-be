-- ============================================================================
-- Cleanup produk LQD di v2 (Postgres bulky_db)
-- ============================================================================
-- Latar belakang:
--   Produk "Mix Product LQDSLE###" / "Palet LQDSLE###" adalah hasil percobaan
--   sync data ke WMS yang GAGAL di v1. Produk ini terlanjur ikut termigrasi ke
--   v2 (626 produk, semua is_sold=true padahal sebagian tidak pernah terjual).
--
-- Strategi (disepakati, paling aman):
--   1. PERTAHANKAN 106 produk LQD yang pernah dijual (order PAID+COMPLETED)
--      beserta riwayat pesanan_item-nya. is_sold=true mereka VALID.
--   2. Soft-delete (bukan hard-delete!) produk LQD yang tidak pernah terjual:
--      - 483 produk yang tidak pernah direferensikan pesanan_item sama sekali
--      - 37 produk yang hanya muncul di pesanan EXPIRED/FAILED/CANCELLED
--
-- Kenapa soft-delete, bukan hard-delete?
--   - pesanan_item (164 baris) tetap bisa dibaca via join yang menghormati
--     deleted_at -> riwayat pesanan utuh.
--   - Menghindari konflik FK bila ada tabel lain yang masih menunjuk produk.
--   - Aman untuk rollback (cukup SET deleted_at = NULL).
--
-- Cara pakai (harus dijalankan via psql, JANGAN lewat aplikasi):
--   psql "postgresql://bulky_db:<password>@195.85.19.129:5432/bulky_db?sslmode=disable" -f cleanup_lqd.sql
--
-- Catatan: skrip ini IDEMPOTENT — dijalankan ulang tidak merusak apa pun.
-- ============================================================================

BEGIN;

-- ----------------------------------------------------------------------------
-- 0) Laporan sebelum eksekusi
-- ----------------------------------------------------------------------------
SELECT 'sebelum: total produk LQD' AS langkah, COUNT(*) FROM produk
WHERE nama_id ILIKE '%LQDSLE%';
SELECT 'sebelum: LQD pernah terjual (PAID+COMPLETED)' AS langkah,
       COUNT(DISTINCT pi.produk_id)
FROM pesanan_item pi
JOIN pesanan p ON p.id = pi.pesanan_id
JOIN produk pr ON pr.id = pi.produk_id
WHERE pr.nama_id ILIKE '%LQDSLE%'
  AND p.payment_status = 'PAID' AND p.order_status = 'COMPLETED';
SELECT 'sebelum: LQD tidak pernah diorder' AS langkah, COUNT(*)
FROM produk pr
WHERE pr.nama_id ILIKE '%LQDSLE%'
  AND NOT EXISTS (SELECT 1 FROM pesanan_item pi WHERE pi.produk_id = pr.id);

-- ----------------------------------------------------------------------------
-- 1) Soft-delete produk LQD yang TIDAK PERNAH diorder sama sekali (483)
--    (tidak ada pesanan_item yang menunjuk -> 100% aman)
-- ----------------------------------------------------------------------------
UPDATE produk SET deleted_at = now()
WHERE deleted_at IS NULL
  AND nama_id ILIKE '%LQDSLE%'
  AND NOT EXISTS (SELECT 1 FROM pesanan_item pi WHERE pi.produk_id = produk.id);

-- ----------------------------------------------------------------------------
-- 2) Soft-delete produk LQD yang HANYA muncul di pesanan GAGAL (37)
--    (EXPIRED/FAILED payment atau CANCELLED order — tidak pernah terjual)
--    pesanan_item-nya tetap ada (soft-delete produk tidak memutus FK)
-- ----------------------------------------------------------------------------
UPDATE produk SET deleted_at = now()
WHERE deleted_at IS NULL
  AND nama_id ILIKE '%LQDSLE%'
  AND EXISTS (
    SELECT 1
    FROM pesanan_item pi
    JOIN pesanan p ON p.id = pi.pesanan_id
    WHERE pi.produk_id = produk.id
      AND (p.payment_status IN ('EXPIRED','FAILED')
           OR p.order_status = 'CANCELLED')
  )
  AND NOT EXISTS (
    SELECT 1
    FROM pesanan_item pi2
    JOIN pesanan p2 ON p2.id = pi2.pesanan_id
    WHERE pi2.produk_id = produk.id
      AND p2.payment_status = 'PAID' AND p2.order_status = 'COMPLETED'
  );

-- ----------------------------------------------------------------------------
-- 3) Laporan sesudah eksekusi
-- ----------------------------------------------------------------------------
SELECT 'sesudah: total produk LQD masih aktif' AS langkah, COUNT(*)
FROM produk WHERE nama_id ILIKE '%LQDSLE%' AND deleted_at IS NULL;
SELECT 'sesudah: produk LQD di-soft-delete' AS langkah, COUNT(*)
FROM produk WHERE nama_id ILIKE '%LQDSLE%' AND deleted_at IS NOT NULL;
SELECT 'sesudah: LQD terjual masih utuh' AS langkah,
       COUNT(DISTINCT pi.produk_id)
FROM pesanan_item pi
JOIN pesanan p ON p.id = pi.pesanan_id
JOIN produk pr ON pr.id = pi.produk_id
WHERE pr.nama_id ILIKE '%LQDSLE%' AND pr.deleted_at IS NULL
  AND p.payment_status = 'PAID' AND p.order_status = 'COMPLETED';
-- pesanan_item yang menunjuk produk LQD soft-deleted (harus 0 yang PAID+COMPLETED)
SELECT 'sesudah: pesanan_item LQD PAID+COMPLETED ke produk soft-deleted (harus 0)' AS langkah,
       COUNT(*)
FROM pesanan_item pi
JOIN pesanan p ON p.id = pi.pesanan_id
JOIN produk pr ON pr.id = pi.produk_id
WHERE pr.nama_id ILIKE '%LQDSLE%' AND pr.deleted_at IS NOT NULL
  AND p.payment_status = 'PAID' AND p.order_status = 'COMPLETED';

COMMIT;
