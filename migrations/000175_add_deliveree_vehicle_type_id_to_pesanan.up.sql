-- Tambah kolom deliveree_vehicle_type_id pada pesanan.
--
-- Latar belakang: storefront menghitung ongkir Deliveree (quote) berdasarkan
-- kendaraan terkecil yang muat (logika kubikasi & berat dari master data
-- deliveree_vehicle_type). Agar ongkir yang ditagih ke buyer saat checkout
-- SELALU konsisten dengan kendaraan yang dipilih saat create booking, storefront
-- menyimpan id_deliveree kendaraan hasil quote ke kolom ini. Saat booking,
-- BE memakai kolom ini jika terisi & masih aktif, dengan fallback ke
-- seleksi kubikasi/berat (SelectVehicle) untuk order lama yang tidak punya nilai.

ALTER TABLE pesanan
    ADD COLUMN IF NOT EXISTS deliveree_vehicle_type_id INTEGER;
