-- migrations/000055_add_unique_constraint_ppn_persentase.up.sql

-- Add unique constraint to persentase field in ppn table
ALTER TABLE ppn ADD CONSTRAINT unique_ppn_persentase UNIQUE (persentase);

COMMENT ON CONSTRAINT unique_ppn_persentase ON ppn IS 'Memastikan tidak ada duplikasi nilai persentase PPN';
