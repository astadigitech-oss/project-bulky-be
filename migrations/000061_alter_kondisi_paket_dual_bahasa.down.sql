ALTER TABLE kondisi_paket DROP COLUMN IF EXISTS nama_en;
ALTER TABLE kondisi_paket RENAME COLUMN nama_id TO nama;
