-- +goose Up
-- Add new columns to businesses table
ALTER TABLE businesses 
    ADD COLUMN IF NOT EXISTS location VARCHAR(255),
    ADD COLUMN IF NOT EXISTS latitude DECIMAL(10,8),
    ADD COLUMN IF NOT EXISTS longitude DECIMAL(11,8),
    ADD COLUMN IF NOT EXISTS opening_hours JSONB,
    ADD COLUMN IF NOT EXISTS social_links JSONB;

-- Add indexes for location-based queries
CREATE INDEX idx_businesses_location ON businesses(location);
CREATE INDEX idx_businesses_is_active ON businesses(is_active);

-- +goose Down
ALTER TABLE businesses 
    DROP COLUMN IF EXISTS location,
    DROP COLUMN IF EXISTS latitude,
    DROP COLUMN IF EXISTS longitude,
    DROP COLUMN IF EXISTS opening_hours,
    DROP COLUMN IF EXISTS social_links;