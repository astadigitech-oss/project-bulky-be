-- Rollback seed data untuk tabel kondisi_produk
DELETE FROM kondisi_produk
WHERE slug IN ('mediocre-60-80', 'brandnew', '90-100');
