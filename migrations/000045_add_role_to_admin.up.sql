-- migrations/000045_add_role_to_admin.up.sql

-- Add role_id column to admin table
ALTER TABLE admin 
ADD COLUMN role_id UUID REFERENCES role(id);

-- Set default role to ADMIN for existing admins
UPDATE admin 
SET role_id = (SELECT id FROM role WHERE kode = 'ADMIN')
WHERE role_id IS NULL;

-- Make role_id NOT NULL after setting defaults
ALTER TABLE admin 
ALTER COLUMN role_id SET NOT NULL;

CREATE INDEX idx_admin_role_id ON admin(role_id);
