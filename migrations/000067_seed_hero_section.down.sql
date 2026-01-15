-- Rollback seed data untuk tabel hero_section
DELETE FROM hero_section 
WHERE nama IN ('Hero Banner 1', 'Hero Banner 2', 'Hero Banner 3');
