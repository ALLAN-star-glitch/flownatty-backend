-- +goose Up
-- Update default values for email verification and 2FA
ALTER TABLE users ALTER COLUMN is_email_verified SET DEFAULT true;
ALTER TABLE users ALTER COLUMN two_factor_enabled SET DEFAULT true;

-- Update existing users to match the new defaults
UPDATE users SET is_email_verified = true WHERE is_email_verified = false;
UPDATE users SET two_factor_enabled = true WHERE two_factor_enabled = false;

-- If email_verified_at is null, set it to verified_at for existing users
UPDATE users 
SET email_verified_at = verified_at 
WHERE email_verified_at IS NULL AND is_email_verified = true;

-- +goose Down
-- Revert to previous defaults
ALTER TABLE users ALTER COLUMN is_email_verified SET DEFAULT false;
ALTER TABLE users ALTER COLUMN two_factor_enabled SET DEFAULT false;

-- Note: We don't revert the data as it would lose the verification status