-- migrations/000055_add_unique_constraint_ppn_persentase.down.sql

-- Remove unique constraint from persentase field in ppn table
ALTER TABLE ppn DROP CONSTRAINT IF EXISTS unique_ppn_persentase;
