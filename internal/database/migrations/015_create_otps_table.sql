-- +goose Up
CREATE TABLE IF NOT EXISTS otps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone_number VARCHAR(20) NOT NULL,
    email VARCHAR(255) NOT NULL,
    otp_code VARCHAR(6) NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'consumer',
    purpose VARCHAR(20) DEFAULT 'signup',
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_otps_email ON otps(email);
CREATE INDEX idx_otps_otp_code ON otps(otp_code);
CREATE INDEX idx_otps_expires ON otps(expires_at);
CREATE INDEX idx_otps_is_used ON otps(is_used);

-- +goose Down
DROP TABLE IF EXISTS otps;