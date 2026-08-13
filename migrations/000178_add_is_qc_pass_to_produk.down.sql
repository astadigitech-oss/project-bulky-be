-- Hapus penanda "QC PASS" pada tabel produk.
ALTER TABLE produk
    DROP COLUMN IF EXISTS is_qc_pass;
