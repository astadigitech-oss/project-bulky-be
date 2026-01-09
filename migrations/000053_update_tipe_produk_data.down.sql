-- migrations/000053_update_tipe_produk_data.down.sql
-- Rollback tipe produk data ke versi sebelumnya

-- Rollback Palet Load -> Paletbox
UPDATE tipe_produk 
SET 
    nama = 'Paletbox',
    slug = 'paletbox',
    deskripsi = 'Produk dalam kemasan paletbox',
    urutan = 1
WHERE slug = 'palet-load';

-- Rollback Truck Load -> Truckload
UPDATE tipe_produk 
SET 
    nama = 'Truckload',
    slug = 'truckload',
    deskripsi = 'Produk dalam kemasan truckload',
    urutan = 3
WHERE slug = 'truck-load';

-- Rollback Container Load -> Container
UPDATE tipe_produk 
SET 
    nama = 'Container',
    slug = 'container',
    deskripsi = 'Produk dalam kemasan container',
    urutan = 2
WHERE slug = 'container-load';
