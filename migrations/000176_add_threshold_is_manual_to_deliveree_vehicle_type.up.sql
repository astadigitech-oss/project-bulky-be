-- Tambah penanda threshold manual pada deliveree_vehicle_type.
--
-- Latar belakang: threshold_kubikasi & threshold_berat awalnya selalu dihitung
-- otomatis oleh Sync (recomputeThresholds). Tim ops membutuhkan threshold yang
-- DIISI MANUAL sebagai batas aman operasional (barang yang hampir mentok kapasitas
-- harus otomatis naik ke kendaraan level atas). Kolom threshold_is_manual menandai
-- bahwa threshold sudah di-set manual oleh ops sehingga Sync TIDAK menimpanya.
--
-- Nilai default false (otomatis). Saat ops menyimpan threshold via API Update,
-- kolom ini di-set true.

ALTER TABLE deliveree_vehicle_type
    ADD COLUMN IF NOT EXISTS threshold_is_manual boolean NOT NULL DEFAULT false;

-- Backfill: threshold existing (dari logika lama = kapasitas kendaraan di bawahnya)
-- di-reset ke kapasitas penuh & ditandai non-manual. Dengan default threshold =
-- kapasitas penuh, perilaku SelectVehicle tetap sama seperti sebelumnya (kendaraan
-- terkecil yang muat) sampai ops mengisi threshold manual.
UPDATE deliveree_vehicle_type
   SET threshold_kubikasi = kubikasi_max,
       threshold_berat    = berat_max,
       threshold_is_manual = false
 WHERE threshold_is_manual = false;
