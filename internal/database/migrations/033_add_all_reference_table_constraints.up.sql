-- internal/database/migrations/033_add_all_reference_table_constraints.up.sql
-- +goose Up

-- ================================================
-- Add display_name columns (if not already added)
-- ================================================

ALTER TABLE business_types ADD COLUMN IF NOT EXISTS display_name VARCHAR(255);
UPDATE business_types SET display_name = name WHERE display_name IS NULL;
ALTER TABLE business_types ALTER COLUMN display_name SET NOT NULL;

ALTER TABLE business_sectors ADD COLUMN IF NOT EXISTS display_name VARCHAR(255);
UPDATE business_sectors SET display_name = name WHERE display_name IS NULL;
ALTER TABLE business_sectors ALTER COLUMN display_name SET NOT NULL;

COMMENT ON COLUMN business_types.display_name IS 'Display name for UI (e.g., "Private Limited Company (Ltd)")';
COMMENT ON COLUMN business_sectors.display_name IS 'Display name for UI (e.g., "Financial Services")';

-- ================================================
-- Add indexes for performance
-- ================================================

CREATE INDEX IF NOT EXISTS idx_business_types_display_name ON business_types(display_name);
CREATE INDEX IF NOT EXISTS idx_business_sectors_display_name ON business_sectors(display_name);
CREATE INDEX IF NOT EXISTS idx_business_subcategories_sector_id ON business_subcategories(sector_id);
CREATE INDEX IF NOT EXISTS idx_product_service_subcategories_category_id ON product_service_subcategories(category_id);

-- +goose Down

-- ================================================
-- Drop indexes
-- ================================================

DROP INDEX IF EXISTS idx_business_types_display_name;
DROP INDEX IF EXISTS idx_business_sectors_display_name;
DROP INDEX IF EXISTS idx_business_subcategories_sector_id;
DROP INDEX IF EXISTS idx_product_service_subcategories_category_id;

-- ================================================
-- Drop display_name columns
-- ================================================

ALTER TABLE business_types DROP COLUMN IF EXISTS display_name;
ALTER TABLE business_sectors DROP COLUMN IF EXISTS display_name;