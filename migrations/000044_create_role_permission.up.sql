-- migrations/000044_create_role_permission.up.sql

CREATE TABLE role_permission (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES role(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permission(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(role_id, permission_id)
);

CREATE INDEX idx_role_permission_role_id ON role_permission(role_id);
CREATE INDEX idx_role_permission_permission_id ON role_permission(permission_id);

-- Assign all permissions to SUPER_ADMIN
INSERT INTO role_permission (role_id, permission_id)
SELECT 
    (SELECT id FROM role WHERE kode = 'SUPER_ADMIN'),
    id
FROM permission;

-- Assign limited permissions to ADMIN (exclude admin:manage, role:manage, system:manage)
INSERT INTO role_permission (role_id, permission_id)
SELECT 
    (SELECT id FROM role WHERE kode = 'ADMIN'),
    id
FROM permission
WHERE kode NOT IN ('admin:manage', 'role:manage', 'system:manage');

-- Assign read-only + basic permissions to STAFF
INSERT INTO role_permission (role_id, permission_id)
SELECT 
    (SELECT id FROM role WHERE kode = 'STAFF'),
    id
FROM permission
WHERE kode LIKE '%:read' 
   OR kode IN ('pesanan:update_status', 'ulasan:approve', 'produk:stock');

COMMENT ON TABLE role_permission IS 'Pivot table untuk relasi role ↔ permission';
