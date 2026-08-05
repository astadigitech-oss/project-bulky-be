-- Snapshot alamat pengiriman pada pesanan.
-- Diperlukan oleh migrasi data v1: tabel `orders` v1 menyimpan alamat sebagai teks
-- (shipping_address, name, phone_number, latitude, longitude) tanpa FK ke `addresses`,
-- sehingga pesanan hasil migrasi tidak bisa mengisi `alamat_buyer_id`.
-- Snapshot juga berguna untuk pesanan baru: alamat buyer bisa diubah/dihapus setelah
-- pesanan dibuat, sementara dokumen pesanan harus mencatat alamat saat transaksi.

ALTER TABLE pesanan
    ADD COLUMN IF NOT EXISTS alamat_snapshot TEXT,
    ADD COLUMN IF NOT EXISTS nama_penerima_snapshot VARCHAR(100),
    ADD COLUMN IF NOT EXISTS telepon_penerima_snapshot VARCHAR(20),
    ADD COLUMN IF NOT EXISTS latitude_snapshot DECIMAL(10, 8),
    ADD COLUMN IF NOT EXISTS longitude_snapshot DECIMAL(11, 8);

COMMENT ON COLUMN pesanan.alamat_snapshot IS 'Alamat lengkap saat pesanan dibuat; fallback tampilan bila alamat_buyer_id NULL (mis. data hasil migrasi v1)';
COMMENT ON COLUMN pesanan.nama_penerima_snapshot IS 'Nama penerima saat pesanan dibuat';
COMMENT ON COLUMN pesanan.telepon_penerima_snapshot IS 'Telepon penerima saat pesanan dibuat';
COMMENT ON COLUMN pesanan.latitude_snapshot IS 'Latitude titik antar saat pesanan dibuat';
COMMENT ON COLUMN pesanan.longitude_snapshot IS 'Longitude titik antar saat pesanan dibuat';

-- Backfill snapshot untuk pesanan existing yang masih punya relasi alamat,
-- agar riwayat lama tetap tampil walau alamat buyer dihapus di kemudian hari.
UPDATE pesanan p
SET alamat_snapshot = COALESCE(p.alamat_snapshot, a.alamat_lengkap),
    nama_penerima_snapshot = COALESCE(p.nama_penerima_snapshot, a.nama_penerima),
    telepon_penerima_snapshot = COALESCE(p.telepon_penerima_snapshot, a.telepon_penerima),
    latitude_snapshot = COALESCE(p.latitude_snapshot, a.latitude),
    longitude_snapshot = COALESCE(p.longitude_snapshot, a.longitude)
FROM alamat_buyer a
WHERE a.id = p.alamat_buyer_id
  AND p.alamat_snapshot IS NULL;

-- Longgarkan constraint: pengiriman non-PICKUP boleh tanpa alamat_buyer_id
-- selama snapshot alamat terisi (kasus data migrasi v1).
ALTER TABLE pesanan DROP CONSTRAINT IF EXISTS chk_alamat_required;

ALTER TABLE pesanan ADD CONSTRAINT chk_alamat_required CHECK (
    (delivery_type = 'PICKUP') OR
    (delivery_type IN ('DELIVEREE', 'FORWARDER', 'FORWARDER_LCL')
        AND (alamat_buyer_id IS NOT NULL OR NULLIF(TRIM(alamat_snapshot), '') IS NOT NULL))
);
