-- Rollback seed data untuk tabel kondisi_paket
DELETE FROM kondisi_paket
WHERE slug IN ('sedang', 'bagus');
