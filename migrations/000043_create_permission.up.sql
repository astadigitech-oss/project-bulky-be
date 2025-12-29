-- migrations/000043_create_permission.up.sql

CREATE TABLE permission (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    kode VARCHAR(50) NOT NULL UNIQUE,
    modul VARCHAR(50) NOT NULL,
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_permission_kode ON permission(kode);
CREATE INDEX idx_permission_modul ON permission(modul);

CREATE TRIGGER trg_permission_updated_at
    BEFORE UPDATE ON permission
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Seed default permissions
INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    -- Dashboard
    ('View Dashboard', 'dashboard:read', 'dashboard', 'Melihat dashboard'),
    
    -- Master Data
    ('View Kategori', 'kategori:read', 'master', 'Melihat kategori produk'),
    ('Manage Kategori', 'kategori:manage', 'master', 'CRUD kategori produk'),
    ('View Tipe Produk', 'tipe_produk:read', 'master', 'Melihat tipe produk'),
    ('Manage Tipe Produk', 'tipe_produk:manage', 'master', 'CRUD tipe produk'),
    ('View Brand', 'brand:read', 'master', 'Melihat brand/merek'),
    ('Manage Brand', 'brand:manage', 'master', 'CRUD brand/merek'),
    ('View Kondisi', 'kondisi:read', 'master', 'Melihat kondisi produk'),
    ('Manage Kondisi', 'kondisi:manage', 'master', 'CRUD kondisi produk'),
    
    -- Produk
    ('View Produk', 'produk:read', 'produk', 'Melihat produk'),
    ('Create Produk', 'produk:create', 'produk', 'Membuat produk baru'),
    ('Update Produk', 'produk:update', 'produk', 'Mengupdate produk'),
    ('Delete Produk', 'produk:delete', 'produk', 'Menghapus produk'),
    ('Manage Stock', 'produk:stock', 'produk', 'Mengatur stok produk'),
    
    -- Pesanan
    ('View Pesanan', 'pesanan:read', 'pesanan', 'Melihat pesanan'),
    ('Update Pesanan', 'pesanan:update', 'pesanan', 'Mengupdate pesanan'),
    ('Update Status Pesanan', 'pesanan:update_status', 'pesanan', 'Mengubah status pesanan'),
    
    -- Buyer
    ('View Buyer', 'buyer:read', 'buyer', 'Melihat data buyer'),
    ('Manage Buyer', 'buyer:manage', 'buyer', 'Mengelola data buyer'),
    
    -- Ulasan
    ('View Ulasan', 'ulasan:read', 'ulasan', 'Melihat ulasan'),
    ('Approve Ulasan', 'ulasan:approve', 'ulasan', 'Approve/reject ulasan'),
    ('Delete Ulasan', 'ulasan:delete', 'ulasan', 'Menghapus ulasan'),
    
    -- Marketing
    ('View Marketing', 'marketing:read', 'marketing', 'Melihat hero section, banner'),
    ('Manage Marketing', 'marketing:manage', 'marketing', 'Mengelola hero section, banner'),
    
    -- Operasional
    ('View Operasional', 'operasional:read', 'operasional', 'Melihat info pickup, dokumen'),
    ('Manage Operasional', 'operasional:manage', 'operasional', 'Mengelola operasional'),
    
    -- Pembayaran
    ('View Pembayaran', 'pembayaran:read', 'pembayaran', 'Melihat metode pembayaran'),
    ('Manage Pembayaran', 'pembayaran:manage', 'pembayaran', 'Mengelola metode pembayaran'),
    
    -- Admin Management
    ('View Admin', 'admin:read', 'admin', 'Melihat data admin'),
    ('Manage Admin', 'admin:manage', 'admin', 'Mengelola admin'),
    ('Manage Role', 'role:manage', 'admin', 'Mengelola role & permission'),
    
    -- System
    ('View System', 'system:read', 'system', 'Melihat force update, maintenance'),
    ('Manage System', 'system:manage', 'system', 'Mengelola system settings'),
    ('View Activity Log', 'activity_log:read', 'system', 'Melihat activity log');

COMMENT ON TABLE permission IS 'Permission granular untuk setiap aksi di sistem';
