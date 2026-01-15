-- Rollback seed data untuk tabel kondisi_paket
DELETE FROM kondisi_paket 
WHERE slug IN ('baik', 'rusak-ringan', 'rusak-sedang', 'rusak-berat');
