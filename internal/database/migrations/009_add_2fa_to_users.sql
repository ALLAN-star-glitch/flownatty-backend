-- +goose Up
ALTER TABLE users ADD COLUMN two_factor_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN two_factor_secret VARCHAR(255);

-- +goose Down
ALTER TABLE users DROP COLUMN two_factor_enabled;
ALTER TABLE users DROP COLUMN two_factor_secret;