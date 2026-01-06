-- migrations/000051_seed_default_super_admin.up.sql

-- Insert default super admin account
-- Email: admin@admin.com
-- Password: password (hashed with bcrypt)
-- Note: Password hash is pre-generated bcrypt hash for "password"

INSERT INTO admin (nama, email, password, role_id, is_active, created_at, updated_at)
SELECT 
    'Super Admin',
    'superadmin@admin.com',
    '$2a$10$DeB03oVv/58Fr8RcdLllEOvJQpY7jsO/r5.iCJk6ADAsHHDQgvGa.',
    (SELECT id FROM role WHERE kode = 'SUPER_ADMIN'),
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM admin WHERE email = 'superadmin@admin.com'
);

COMMENT ON COLUMN admin.password IS 'Password (hashed bcrypt) - Default super admin password: "password"';
