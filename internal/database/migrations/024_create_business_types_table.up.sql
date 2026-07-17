-- internal/database/migrations/024_create_business_types_table.up.sql
-- +goose Up
-- Business Types (Legal Structure)
CREATE TABLE IF NOT EXISTS business_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_business_types_name ON business_types(name);
CREATE INDEX IF NOT EXISTS idx_business_types_is_active ON business_types(is_active);
CREATE INDEX IF NOT EXISTS idx_business_types_deleted_at ON business_types(deleted_at);

COMMENT ON TABLE business_types IS 'Business legal structures (Sole Proprietorship, Partnership, etc.)';
COMMENT ON COLUMN business_types.name IS 'Name of the business type';
COMMENT ON COLUMN business_types.description IS 'Description of the business type';
COMMENT ON COLUMN business_types.icon IS 'Icon name for UI display';

-- +goose Down
DROP TABLE IF EXISTS business_types;