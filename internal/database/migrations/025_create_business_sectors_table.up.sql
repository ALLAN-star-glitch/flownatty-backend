-- internal/database/migrations/025_create_business_sectors_table.up.sql
-- +goose Up
-- Business Sectors (Industry)
CREATE TABLE IF NOT EXISTS business_sectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_business_sectors_name ON business_sectors(name);
CREATE INDEX IF NOT EXISTS idx_business_sectors_is_active ON business_sectors(is_active);
CREATE INDEX IF NOT EXISTS idx_business_sectors_deleted_at ON business_sectors(deleted_at);

COMMENT ON TABLE business_sectors IS 'Business industry sectors (Retail, Health, Technology, etc.)';
COMMENT ON COLUMN business_sectors.name IS 'Name of the sector';
COMMENT ON COLUMN business_sectors.description IS 'Description of the sector';
COMMENT ON COLUMN business_sectors.icon IS 'Icon name for UI display';
COMMENT ON COLUMN business_sectors.sort_order IS 'Order for display in UI';

-- +goose Down
DROP TABLE IF EXISTS business_sectors;