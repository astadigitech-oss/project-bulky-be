-- Rollback seed data untuk tabel banner_event_promo
DELETE FROM banner_event_promo 
WHERE nama IN ('Promo Tahun Baru', 'Flash Sale Weekend', 'Diskon Elektronik');
