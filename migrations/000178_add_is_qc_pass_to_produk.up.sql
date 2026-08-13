-- Tambah penanda "QC PASS" pada tabel produk.
--
-- Latar belakang: Admin panel membutuhkan flag boolean untuk menandai produk
-- yang sudah lolos quality control (QC PASS). Admin dapat enable/disable flag
-- ini per produk melalui menu triple-dot di list product, sama seperti pola
-- toggle is_sale.
--
-- Nilai default false (belum QC PASS).

ALTER TABLE produk
    ADD COLUMN IF NOT EXISTS is_qc_pass boolean NOT NULL DEFAULT false;
