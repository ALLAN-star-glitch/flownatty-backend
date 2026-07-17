-- internal/database/migrations/029_update_businesses_add_classification.up.sql
-- +goose Up
-- Add classification columns to businesses table
ALTER TABLE businesses 
ADD COLUMN IF NOT EXISTS business_type_id UUID REFERENCES business_types(id),
ADD COLUMN IF NOT EXISTS sector_id UUID REFERENCES business_sectors(id),
ADD COLUMN IF NOT EXISTS subcategory_id UUID REFERENCES business_subcategories(id),
ADD COLUMN IF NOT EXISTS employee_count INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS year_established INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS business_size VARCHAR(20) DEFAULT 'micro';

-- Drop the old category column if it exists
ALTER TABLE businesses DROP COLUMN IF EXISTS category;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_businesses_business_type_id ON businesses(business_type_id);
CREATE INDEX IF NOT EXISTS idx_businesses_sector_id ON businesses(sector_id);
CREATE INDEX IF NOT EXISTS idx_businesses_subcategory_id ON businesses(subcategory_id);
CREATE INDEX IF NOT EXISTS idx_businesses_business_size ON businesses(business_size);

-- Add check constraint for business_size (using simpler syntax)
ALTER TABLE businesses ADD CONSTRAINT chk_business_size 
    CHECK (business_size IN ('micro', 'small', 'medium', 'large'));

-- Add comments
COMMENT ON COLUMN businesses.business_type_id IS 'Legal structure of the business';
COMMENT ON COLUMN businesses.sector_id IS 'Industry sector of the business';
COMMENT ON COLUMN businesses.subcategory_id IS 'Detailed subcategory of the business';
COMMENT ON COLUMN businesses.employee_count IS 'Number of employees';
COMMENT ON COLUMN businesses.year_established IS 'Year the business was established';
COMMENT ON COLUMN businesses.business_size IS 'Business size: micro, small, medium, large';

-- +goose Down
-- Remove classification columns from businesses table
ALTER TABLE businesses 
DROP COLUMN IF EXISTS business_type_id,
DROP COLUMN IF EXISTS sector_id,
DROP COLUMN IF EXISTS subcategory_id,
DROP COLUMN IF EXISTS employee_count,
DROP COLUMN IF EXISTS year_established,
DROP COLUMN IF EXISTS business_size;

-- Re-add the category column (since we dropped it in up)
ALTER TABLE businesses ADD COLUMN IF NOT EXISTS category VARCHAR(100) DEFAULT 'general';

-- Drop indexes
DROP INDEX IF EXISTS idx_businesses_business_type_id;
DROP INDEX IF EXISTS idx_businesses_sector_id;
DROP INDEX IF EXISTS idx_businesses_subcategory_id;
DROP INDEX IF EXISTS idx_businesses_business_size;

-- Drop constraint
ALTER TABLE businesses DROP CONSTRAINT IF EXISTS chk_business_size;