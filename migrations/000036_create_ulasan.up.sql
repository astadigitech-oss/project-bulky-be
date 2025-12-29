-- migrations/000035_create_ulasan.up.sql

-- =====================================================
-- TABEL: ulasan
-- =====================================================
-- Review/ulasan dari buyer untuk produk yang dibeli
-- 
-- Business Rules:
-- 1. Hanya bisa review setelah pesanan COMPLETED
-- 2. 1 buyer hanya bisa 1 review per item pesanan
-- 3. Rating 1-5 (required), komentar (opsional)
-- 4. Admin harus approve agar tampil di publik
-- =====================================================

CREATE TABLE ulasan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Relations
    pesanan_id UUID NOT NULL REFERENCES pesanan(id),
    pesanan_item_id UUID NOT NULL REFERENCES pesanan_item(id),
    buyer_id UUID NOT NULL REFERENCES buyer(id),
    produk_id UUID NOT NULL REFERENCES produk(id),
    
    -- Review content
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    komentar TEXT,
    gambar VARCHAR(255),
    
    -- Approval
    is_approved BOOLEAN DEFAULT false,
    approved_at TIMESTAMP,
    approved_by UUID REFERENCES admin(id),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    -- Constraints
    CONSTRAINT unique_ulasan_per_item UNIQUE (pesanan_item_id)
);

-- Indexes
CREATE INDEX idx_ulasan_pesanan_id ON ulasan(pesanan_id);
CREATE INDEX idx_ulasan_buyer_id ON ulasan(buyer_id);
CREATE INDEX idx_ulasan_produk_id ON ulasan(produk_id);
CREATE INDEX idx_ulasan_is_approved ON ulasan(is_approved) WHERE deleted_at IS NULL;
CREATE INDEX idx_ulasan_rating ON ulasan(rating);
CREATE INDEX idx_ulasan_created_at ON ulasan(created_at DESC);

-- Trigger updated_at
CREATE TRIGGER trg_ulasan_updated_at
    BEFORE UPDATE ON ulasan
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- TRIGGER: Validate order is completed before review
-- =====================================================

CREATE OR REPLACE FUNCTION fn_validate_ulasan_order_completed()
RETURNS TRIGGER AS $$
DECLARE
    v_order_status VARCHAR(20);
BEGIN
    -- Check if order is completed
    SELECT order_status INTO v_order_status
    FROM pesanan
    WHERE id = NEW.pesanan_id;
    
    IF v_order_status != 'COMPLETED' THEN
        RAISE EXCEPTION 'Ulasan hanya dapat diberikan untuk pesanan yang sudah selesai (COMPLETED)';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_validate_ulasan_order
    BEFORE INSERT ON ulasan
    FOR EACH ROW
    EXECUTE FUNCTION fn_validate_ulasan_order_completed();

-- =====================================================
-- TRIGGER: Auto-set approved_at when approved
-- =====================================================

CREATE OR REPLACE FUNCTION fn_set_ulasan_approved_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_approved = true AND OLD.is_approved = false THEN
        NEW.approved_at := NOW();
    END IF;
    
    IF NEW.is_approved = false THEN
        NEW.approved_at := NULL;
        NEW.approved_by := NULL;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_ulasan_approved_at
    BEFORE UPDATE OF is_approved ON ulasan
    FOR EACH ROW
    EXECUTE FUNCTION fn_set_ulasan_approved_at();

-- Comments
COMMENT ON TABLE ulasan IS 'Review/ulasan dari buyer. Hanya untuk pesanan COMPLETED. Admin harus approve agar tampil publik.';
COMMENT ON COLUMN ulasan.pesanan_item_id IS 'Item spesifik yang di-review. 1 ulasan per item (unique constraint).';
COMMENT ON COLUMN ulasan.rating IS 'Rating bintang 1-5';
COMMENT ON COLUMN ulasan.komentar IS 'Isi ulasan (opsional)';
COMMENT ON COLUMN ulasan.gambar IS 'Foto ulasan (opsional, max 1)';
COMMENT ON COLUMN ulasan.is_approved IS 'Status approval. Harus true agar tampil di publik.';
