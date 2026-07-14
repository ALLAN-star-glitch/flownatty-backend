-- +goose Up
-- Add is_active column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;

-- Update existing users to be active
UPDATE users SET is_active = true WHERE is_active IS NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS is_active;