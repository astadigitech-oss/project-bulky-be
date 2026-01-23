-- migrations/000091_seed_warehouse.down.sql
-- Rollback: Hapus seed data warehouse

DELETE FROM warehouse WHERE slug = 'gudang-cilodong';
