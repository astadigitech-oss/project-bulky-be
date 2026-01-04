-- Add foto_url column to buyer table
ALTER TABLE buyer ADD COLUMN IF NOT EXISTS foto_url TEXT;
