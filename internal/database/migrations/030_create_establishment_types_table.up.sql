-- internal/database/migrations/030_create_establishment_types_table.up.sql
-- +goose Up
-- Create establishment types table
CREATE TABLE IF NOT EXISTS establishment_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    category VARCHAR(20) DEFAULT 'physical', -- physical, digital, hybrid, mobile, home
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_establishment_types_name ON establishment_types(name);
CREATE INDEX IF NOT EXISTS idx_establishment_types_category ON establishment_types(category);
CREATE INDEX IF NOT EXISTS idx_establishment_types_is_active ON establishment_types(is_active);

COMMENT ON TABLE establishment_types IS 'Types of business establishments (physical, digital, hybrid, etc.)';
COMMENT ON COLUMN establishment_types.category IS 'Category of establishment: physical, digital, hybrid, mobile, home';

-- Add establishment_type_id to businesses table
ALTER TABLE businesses 
ADD COLUMN IF NOT EXISTS establishment_type_id UUID REFERENCES establishment_types(id),
ADD COLUMN IF NOT EXISTS market_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS stall_number VARCHAR(50),
ADD COLUMN IF NOT EXISTS website VARCHAR(255),
ADD COLUMN IF NOT EXISTS social_media VARCHAR(255),
ADD COLUMN IF NOT EXISTS is_remote BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS is_delivery BOOLEAN DEFAULT false;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_businesses_establishment_type_id ON businesses(establishment_type_id);
CREATE INDEX IF NOT EXISTS idx_businesses_market_name ON businesses(market_name);
CREATE INDEX IF NOT EXISTS idx_businesses_is_remote ON businesses(is_remote);
CREATE INDEX IF NOT EXISTS idx_businesses_is_delivery ON businesses(is_delivery);

COMMENT ON COLUMN businesses.establishment_type_id IS 'Type of establishment';
COMMENT ON COLUMN businesses.website IS 'Website URL for digital/remote businesses';
COMMENT ON COLUMN businesses.social_media IS 'Social media handles (JSON or comma-separated)';
COMMENT ON COLUMN businesses.is_remote IS 'Whether the business operates remotely';
COMMENT ON COLUMN businesses.is_delivery IS 'Whether the business offers delivery';

-- +goose Down
ALTER TABLE businesses 
DROP COLUMN IF EXISTS establishment_type_id,
DROP COLUMN IF EXISTS market_name,
DROP COLUMN IF EXISTS stall_number,
DROP COLUMN IF EXISTS website,
DROP COLUMN IF EXISTS social_media,
DROP COLUMN IF EXISTS is_remote,
DROP COLUMN IF EXISTS is_delivery;

DROP INDEX IF EXISTS idx_businesses_establishment_type_id;
DROP INDEX IF EXISTS idx_businesses_market_name;
DROP INDEX IF EXISTS idx_businesses_is_remote;
DROP INDEX IF EXISTS idx_businesses_is_delivery;

DROP TABLE IF EXISTS establishment_types;