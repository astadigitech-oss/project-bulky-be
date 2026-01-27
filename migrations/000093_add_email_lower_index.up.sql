-- Add case-insensitive email index for admin table
CREATE INDEX idx_admin_email_lower ON admin(LOWER(email)) WHERE deleted_at IS NULL;

-- Add case-insensitive email index for buyer table
CREATE INDEX idx_buyer_email_lower ON buyer(LOWER(email)) WHERE deleted_at IS NULL;
