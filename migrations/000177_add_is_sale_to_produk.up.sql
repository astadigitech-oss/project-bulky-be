-- Tambah penanda ribbon "SALE" pada tabel produk.
--
-- Latar belakang: FE (mobile & web) membutuhkan flag boolean untuk menampilkan
-- ribbon "SALE" di item produk palet pada list product. Admin panel dapat
-- enable/disable flag ini per produk melalui menu triple-dot di list product.
--
-- Nilai default false (ribbon tidak tampil).

ALTER TABLE produk
    ADD COLUMN IF NOT EXISTS is_sale boolean NOT NULL DEFAULT false;
