-- Delete sample data
DELETE FROM kupon_usage WHERE pesanan_id IN (
    SELECT id FROM pesanan WHERE kode LIKE 'ORD-SAMPLE%'
);

DELETE FROM ulasan WHERE pesanan_id IN (
    SELECT id FROM pesanan WHERE kode LIKE 'ORD-SAMPLE%'
);

DELETE FROM pesanan_status_history WHERE pesanan_id IN (
    SELECT id FROM pesanan WHERE kode LIKE 'ORD-SAMPLE%'
);

DELETE FROM pesanan_pembayaran WHERE pesanan_id IN (
    SELECT id FROM pesanan WHERE kode LIKE 'ORD-SAMPLE%'
);

DELETE FROM pesanan_item WHERE pesanan_id IN (
    SELECT id FROM pesanan WHERE kode LIKE 'ORD-SAMPLE%'
);

DELETE FROM pesanan WHERE kode LIKE 'ORD-SAMPLE%';

-- Delete sample alamat_buyer if created
DELETE FROM alamat_buyer WHERE label = 'Alamat Sample' AND alamat_lengkap LIKE 'Jl. Sample No. 123%';
