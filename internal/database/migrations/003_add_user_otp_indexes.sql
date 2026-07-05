-- +goose Up
CREATE INDEX IF NOT EXISTS idx_users_email_role ON users(email, role);
CREATE INDEX IF NOT EXISTS idx_users_phone_role ON users(phone_number, role);
CREATE INDEX IF NOT EXISTS idx_users_phone_email ON users(phone_number, email);

CREATE INDEX IF NOT EXISTS idx_otps_email_code ON otps(email, otp_code);
CREATE INDEX IF NOT EXISTS idx_otps_email_used_expiry ON otps(email, is_used, expires_at);
CREATE INDEX IF NOT EXISTS idx_otps_email_created ON otps(email, created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_users_email_role;
DROP INDEX IF EXISTS idx_users_phone_role;
DROP INDEX IF EXISTS idx_users_phone_email;
DROP INDEX IF EXISTS idx_otps_email_code;
DROP INDEX IF EXISTS idx_otps_email_used_expiry;
DROP INDEX IF EXISTS idx_otps_email_created;