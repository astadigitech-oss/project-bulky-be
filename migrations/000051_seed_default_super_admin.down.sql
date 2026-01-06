-- migrations/000051_seed_default_super_admin.down.sql

-- Remove default super admin account
DELETE FROM admin WHERE email = 'admin@admin.com';
