-- +goose Up
-- Remove user_id column from businesses table since we now use business_members
ALTER TABLE businesses DROP COLUMN IF EXISTS user_id;

-- +goose Down
-- Add user_id column back (if needed for rollback)
ALTER TABLE businesses ADD COLUMN user_id UUID REFERENCES users(id);