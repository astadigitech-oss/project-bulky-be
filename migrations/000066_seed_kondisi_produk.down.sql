-- Rollback seed data untuk tabel kondisi_produk
DELETE FROM kondisi_produk 
WHERE slug IN ('baru', 'bekas-seperti-baru', 'bekas-baik', 'bekas-cukup-baik', 'rusak');
