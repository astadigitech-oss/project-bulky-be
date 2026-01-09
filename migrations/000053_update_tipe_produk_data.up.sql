-- migrations/000053_update_tipe_produk_data.up.sql
-- Update tipe produk data dengan nama, slug, dan deskripsi yang baru

-- Update Paletbox -> Palet Load
UPDATE tipe_produk 
SET 
    nama = 'Palet Load',
    slug = 'palet-load',
    deskripsi = 'Layanan muatan kargo kecil, estimasi 1 palet berbagai kategori. Jangkauan: JABODETABEK dan BANDUNG.',
    urutan = 1
WHERE slug = 'paletbox';

-- Update Truckload -> Truck Load
UPDATE tipe_produk 
SET 
    nama = 'Truck Load',
    slug = 'truck-load',
    deskripsi = 'Layanan muatan kargo sedang, per kategori atau campuran. Jangkauan: Pulau Sumatera, Jawa, dan Bali.',
    urutan = 2
WHERE slug = 'truckload';

-- Update Container -> Container Load
UPDATE tipe_produk 
SET 
    nama = 'Container Load',
    slug = 'container-load',
    deskripsi = 'Layanan muatan kargo besar, kategori campuran untuk usaha menengah ke atas. Jangkauan: seluruh Indonesia termasuk luar pulau.',
    urutan = 3
WHERE slug = 'container';
