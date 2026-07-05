-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_active_at TIMESTAMP;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS last_active_at;