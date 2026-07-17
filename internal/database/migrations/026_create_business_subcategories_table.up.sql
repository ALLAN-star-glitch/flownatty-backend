-- internal/database/migrations/026_create_business_subcategories_table.up.sql
-- +goose Up
-- Business Subcategories (Detailed)
CREATE TABLE IF NOT EXISTS business_subcategories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sector_id UUID NOT NULL REFERENCES business_sectors(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    UNIQUE(sector_id, name)
);

CREATE INDEX IF NOT EXISTS idx_business_subcategories_sector_id ON business_subcategories(sector_id);
CREATE INDEX IF NOT EXISTS idx_business_subcategories_name ON business_subcategories(name);
CREATE INDEX IF NOT EXISTS idx_business_subcategories_is_active ON business_subcategories(is_active);
CREATE INDEX IF NOT EXISTS idx_business_subcategories_deleted_at ON business_subcategories(deleted_at);

COMMENT ON TABLE business_subcategories IS 'Detailed business subcategories within sectors';
COMMENT ON COLUMN business_subcategories.sector_id IS 'Foreign key to business_sectors';
COMMENT ON COLUMN business_subcategories.name IS 'Name of the subcategory (e.g., Supermarket, Boutique)';
COMMENT ON COLUMN business_subcategories.description IS 'Description of the subcategory';
COMMENT ON COLUMN business_subcategories.icon IS 'Icon name for UI display';

-- +goose Down
DROP TABLE IF EXISTS business_subcategories;