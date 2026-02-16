-- =====================================================
-- SEED SAMPLE DATA UNTUK TESTING ADMIN PANEL
-- Menggunakan produk dan buyer yang sudah di-seed
-- =====================================================

DO $$
DECLARE
    v_buyer_id UUID;
    v_admin_id UUID;
    v_produk_id UUID;
    v_produk_nama VARCHAR;
    v_pesanan_id UUID;
    v_pesanan_item_id UUID;
    v_kupon_id UUID;
    v_metode_pembayaran_id UUID;
    v_alamat_buyer_id UUID;
BEGIN
    -- Get first buyer
    SELECT id INTO v_buyer_id FROM buyer WHERE deleted_at IS NULL ORDER BY created_at LIMIT 1;
    
    -- Get first admin
    SELECT id INTO v_admin_id FROM admin WHERE deleted_at IS NULL ORDER BY created_at LIMIT 1;
    
    -- Get first produk
    SELECT id, nama_id INTO v_produk_id, v_produk_nama FROM produk WHERE deleted_at IS NULL ORDER BY created_at LIMIT 1;
    
    -- Get first kupon (if exists)
    SELECT id INTO v_kupon_id FROM kupon WHERE deleted_at IS NULL ORDER BY created_at LIMIT 1;
    
    -- Get first metode pembayaran (if exists)
    SELECT id INTO v_metode_pembayaran_id FROM metode_pembayaran WHERE deleted_at IS NULL ORDER BY created_at LIMIT 1;
    
    -- Get first alamat buyer (if exists)
    SELECT id INTO v_alamat_buyer_id FROM alamat_buyer WHERE buyer_id = v_buyer_id AND deleted_at IS NULL ORDER BY created_at LIMIT 1;
    
    -- Skip if no buyer or produk found
    IF v_buyer_id IS NULL OR v_produk_id IS NULL THEN
        RAISE NOTICE 'No buyer or produk found, skipping seed';
        RETURN;
    END IF;
    
    -- Create alamat buyer if not exists (required for DELIVEREE/FORWARDER)
    IF v_alamat_buyer_id IS NULL THEN
        INSERT INTO alamat_buyer (
            buyer_id, label, nama_penerima, telepon_penerima,
            provinsi, kota, alamat_lengkap, is_default
        ) VALUES (
            v_buyer_id,
            'Alamat Sample',
            (SELECT nama FROM buyer WHERE id = v_buyer_id),
            '081234567890',
            'DKI Jakarta',
            'Jakarta Selatan',
            'Jl. Sample No. 123, Jakarta Selatan',
            true
        ) RETURNING id INTO v_alamat_buyer_id;
        
        RAISE NOTICE 'Created sample alamat_buyer: %', v_alamat_buyer_id;
    END IF;
    
    -- =====================================================
    -- 1. CREATE SAMPLE PESANAN (COMPLETED)
    -- =====================================================
    INSERT INTO pesanan (
        id, kode, buyer_id, delivery_type, payment_type, 
        payment_status, order_status,
        biaya_produk, biaya_pengiriman, biaya_ppn, total,
        created_at, updated_at
    ) VALUES (
        uuid_generate_v4(),
        'ORD-SAMPLE001',
        v_buyer_id,
        'PICKUP',
        'REGULAR',
        'PAID',
        'COMPLETED',
        5000000.00,
        0,
        550000.00,
        5550000.00,
        NOW() - INTERVAL '7 days',
        NOW() - INTERVAL '1 day'
    ) RETURNING id INTO v_pesanan_id;
    
    -- Create pesanan item
    INSERT INTO pesanan_item (
        id, pesanan_id, produk_id, nama_produk, qty, 
        harga_satuan, diskon_satuan, subtotal
    ) VALUES (
        uuid_generate_v4(),
        v_pesanan_id,
        v_produk_id,
        v_produk_nama,
        1,
        5000000.00,
        0,
        5000000.00
    ) RETURNING id INTO v_pesanan_item_id;
    
    -- Create pesanan pembayaran (if metode pembayaran exists)
    IF v_metode_pembayaran_id IS NOT NULL THEN
        INSERT INTO pesanan_pembayaran (
            pesanan_id, buyer_id, metode_pembayaran_id, jumlah, status, paid_at
        ) VALUES (
            v_pesanan_id,
            v_buyer_id,
            v_metode_pembayaran_id,
            5550000.00,
            'PAID',
            NOW() - INTERVAL '6 days'
        );
    ELSE
        -- Create without metode pembayaran if not exists
        INSERT INTO pesanan_pembayaran (
            pesanan_id, buyer_id, jumlah, status, paid_at
        ) VALUES (
            v_pesanan_id,
            v_buyer_id,
            5550000.00,
            'PAID',
            NOW() - INTERVAL '6 days'
        );
    END IF;
    
    -- Create status history
    INSERT INTO pesanan_status_history (
        pesanan_id, status_from, status_to, status_type, created_at
    ) VALUES 
    (v_pesanan_id, NULL, 'PENDING', 'ORDER', NOW() - INTERVAL '7 days'),
    (v_pesanan_id, 'PENDING', 'PAID', 'PAYMENT', NOW() - INTERVAL '6 days'),
    (v_pesanan_id, 'PENDING', 'PROCESSING', 'ORDER', NOW() - INTERVAL '6 days'),
    (v_pesanan_id, 'PROCESSING', 'READY', 'ORDER', NOW() - INTERVAL '4 days'),
    (v_pesanan_id, 'READY', 'SHIPPED', 'ORDER', NOW() - INTERVAL '2 days'),
    (v_pesanan_id, 'SHIPPED', 'COMPLETED', 'ORDER', NOW() - INTERVAL '1 day');
    
    -- =====================================================
    -- 2. CREATE SAMPLE ULASAN
    -- =====================================================
    INSERT INTO ulasan (
        pesanan_id, pesanan_item_id, buyer_id, produk_id,
        rating, komentar, is_approved, approved_at, approved_by,
        created_at
    ) VALUES (
        v_pesanan_id,
        v_pesanan_item_id,
        v_buyer_id,
        v_produk_id,
        5,
        'Produk sangat bagus! Sesuai dengan deskripsi. Packing rapi dan pengiriman cepat. Recommended seller!',
        true,
        NOW() - INTERVAL '1 day',
        v_admin_id,
        NOW() - INTERVAL '2 days'
    );
    
    -- =====================================================
    -- 3. CREATE KUPON USAGE (if kupon exists)
    -- =====================================================
    IF v_kupon_id IS NOT NULL THEN
        INSERT INTO kupon_usage (
            kupon_id, buyer_id, pesanan_id, kode_kupon, nilai_potongan
        ) VALUES (
            v_kupon_id,
            v_buyer_id,
            v_pesanan_id,
            (SELECT kode FROM kupon WHERE id = v_kupon_id),
            50000.00
        );
    END IF;
    
    -- =====================================================
    -- 4. CREATE ADDITIONAL PESANAN (VARIOUS STATUS)
    -- =====================================================
    
    -- Pesanan PROCESSING
    INSERT INTO pesanan (
        kode, buyer_id, delivery_type, alamat_buyer_id, payment_type, 
        payment_status, order_status,
        biaya_produk, biaya_pengiriman, biaya_ppn, total,
        created_at
    ) VALUES (
        'ORD-SAMPLE002',
        v_buyer_id,
        'DELIVEREE',
        v_alamat_buyer_id,
        'REGULAR',
        'PAID',
        'PROCESSING',
        10000000.00,
        500000.00,
        1155000.00,
        11655000.00,
        NOW() - INTERVAL '2 days'
    ) RETURNING id INTO v_pesanan_id;
    
    -- Create pesanan item for SAMPLE002
    INSERT INTO pesanan_item (
        pesanan_id, produk_id, nama_produk, qty, 
        harga_satuan, diskon_satuan, subtotal
    ) VALUES (
        v_pesanan_id,
        v_produk_id,
        v_produk_nama,
        2,
        5000000.00,
        0,
        10000000.00
    );
    
    -- Create pembayaran for SAMPLE002
    IF v_metode_pembayaran_id IS NOT NULL THEN
        INSERT INTO pesanan_pembayaran (
            pesanan_id, buyer_id, metode_pembayaran_id, jumlah, status, paid_at
        ) VALUES (
            v_pesanan_id,
            v_buyer_id,
            v_metode_pembayaran_id,
            11655000.00,
            'PAID',
            NOW() - INTERVAL '2 days'
        );
    ELSE
        INSERT INTO pesanan_pembayaran (
            pesanan_id, buyer_id, jumlah, status, paid_at
        ) VALUES (
            v_pesanan_id,
            v_buyer_id,
            11655000.00,
            'PAID',
            NOW() - INTERVAL '2 days'
        );
    END IF;
    
    -- Create status history for SAMPLE002
    INSERT INTO pesanan_status_history (
        pesanan_id, status_from, status_to, status_type, created_at
    ) VALUES 
    (v_pesanan_id, NULL, 'PENDING', 'ORDER', NOW() - INTERVAL '2 days'),
    (v_pesanan_id, 'PENDING', 'PAID', 'PAYMENT', NOW() - INTERVAL '2 days'),
    (v_pesanan_id, 'PENDING', 'PROCESSING', 'ORDER', NOW() - INTERVAL '2 days');
    
    -- Pesanan PENDING (waiting payment)
    INSERT INTO pesanan (
        kode, buyer_id, delivery_type, payment_type, 
        payment_status, order_status,
        biaya_produk, biaya_pengiriman, biaya_ppn, total,
        created_at
    ) VALUES (
        'ORD-SAMPLE003',
        v_buyer_id,
        'PICKUP',
        'REGULAR',
        'PENDING',
        'PENDING',
        7500000.00,
        0,
        825000.00,
        8325000.00,
        NOW() - INTERVAL '1 hour'
    ) RETURNING id INTO v_pesanan_id;
    
    -- Create pesanan item for SAMPLE003
    INSERT INTO pesanan_item (
        pesanan_id, produk_id, nama_produk, qty, 
        harga_satuan, diskon_satuan, subtotal
    ) VALUES (
        v_pesanan_id,
        v_produk_id,
        v_produk_nama,
        1,
        7500000.00,
        0,
        7500000.00
    );
    
    -- Create pembayaran for SAMPLE003 (PENDING)
    IF v_metode_pembayaran_id IS NOT NULL THEN
        INSERT INTO pesanan_pembayaran (
            pesanan_id, buyer_id, metode_pembayaran_id, jumlah, status, expired_at
        ) VALUES (
            v_pesanan_id,
            v_buyer_id,
            v_metode_pembayaran_id,
            8325000.00,
            'PENDING',
            NOW() + INTERVAL '23 hours'
        );
    ELSE
        INSERT INTO pesanan_pembayaran (
            pesanan_id, buyer_id, jumlah, status, expired_at
        ) VALUES (
            v_pesanan_id,
            v_buyer_id,
            8325000.00,
            'PENDING',
            NOW() + INTERVAL '23 hours'
        );
    END IF;
    
    -- Create status history for SAMPLE003
    INSERT INTO pesanan_status_history (
        pesanan_id, status_from, status_to, status_type, created_at
    ) VALUES (
        v_pesanan_id, NULL, 'PENDING', 'ORDER', NOW() - INTERVAL '1 hour'
    );
    
    RAISE NOTICE 'Sample data seeded successfully';
END $$;
