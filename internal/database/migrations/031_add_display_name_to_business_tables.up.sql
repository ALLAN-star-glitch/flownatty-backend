-- internal/database/migrations/031_add_display_name_to_business_tables.up.sql
-- +goose Up

-- Add display_name to business_types
ALTER TABLE business_types ADD COLUMN IF NOT EXISTS display_name VARCHAR(255);
UPDATE business_types SET display_name = name WHERE display_name IS NULL;
ALTER TABLE business_types ALTER COLUMN display_name SET NOT NULL;

-- Add display_name to business_sectors
ALTER TABLE business_sectors ADD COLUMN IF NOT EXISTS display_name VARCHAR(255);
UPDATE business_sectors SET display_name = name WHERE display_name IS NULL;
ALTER TABLE business_sectors ALTER COLUMN display_name SET NOT NULL;

-- Add comments for documentation
COMMENT ON COLUMN business_types.display_name IS 'Display name for UI (e.g., "Private Limited Company (Ltd)")';
COMMENT ON COLUMN business_sectors.display_name IS 'Display name for UI (e.g., "Financial Services")';

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_business_types_display_name ON business_types(display_name);
CREATE INDEX IF NOT EXISTS idx_business_sectors_display_name ON business_sectors(display_name);

-- +goose Down

-- Remove indexes
DROP INDEX IF EXISTS idx_business_types_display_name;
DROP INDEX IF EXISTS idx_business_sectors_display_name;

-- Remove display_name columns
ALTER TABLE business_types DROP COLUMN IF EXISTS display_name;
ALTER TABLE business_sectors DROP COLUMN IF EXISTS display_name;