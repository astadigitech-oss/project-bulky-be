-- Kembalikan constraint ke bentuk migrasi 000161 (wajib alamat_buyer_id untuk non-PICKUP).
-- Baris pesanan yang hanya punya snapshot akan menahan perintah ini; hapus/perbaiki
-- baris tersebut lebih dulu bila rollback benar-benar dibutuhkan.
ALTER TABLE pesanan DROP CONSTRAINT IF EXISTS chk_alamat_required;

ALTER TABLE pesanan ADD CONSTRAINT chk_alamat_required CHECK (
    (delivery_type = 'PICKUP') OR
    (delivery_type IN ('DELIVEREE', 'FORWARDER', 'FORWARDER_LCL') AND alamat_buyer_id IS NOT NULL)
);

ALTER TABLE pesanan
    DROP COLUMN IF EXISTS longitude_snapshot,
    DROP COLUMN IF EXISTS latitude_snapshot,
    DROP COLUMN IF EXISTS telepon_penerima_snapshot,
    DROP COLUMN IF EXISTS nama_penerima_snapshot,
    DROP COLUMN IF EXISTS alamat_snapshot;
