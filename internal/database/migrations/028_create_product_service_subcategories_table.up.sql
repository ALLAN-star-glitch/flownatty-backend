-- internal/database/migrations/028_create_product_service_subcategories_table.up.sql
-- +goose Up
-- Product/Service Subcategories
CREATE TABLE IF NOT EXISTS product_service_subcategories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES product_service_categories(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    UNIQUE(category_id, name)
);

CREATE INDEX IF NOT EXISTS idx_product_service_subcategories_category_id ON product_service_subcategories(category_id);
CREATE INDEX IF NOT EXISTS idx_product_service_subcategories_name ON product_service_subcategories(name);
CREATE INDEX IF NOT EXISTS idx_product_service_subcategories_is_active ON product_service_subcategories(is_active);
CREATE INDEX IF NOT EXISTS idx_product_service_subcategories_deleted_at ON product_service_subcategories(deleted_at);

COMMENT ON TABLE product_service_subcategories IS 'Subcategories for products and services';
COMMENT ON COLUMN product_service_subcategories.category_id IS 'Foreign key to product_service_categories';
COMMENT ON COLUMN product_service_subcategories.name IS 'Subcategory name (e.g., Women''s Clothing, Phones)';
COMMENT ON COLUMN product_service_subcategories.description IS 'Description of the subcategory';
COMMENT ON COLUMN product_service_subcategories.icon IS 'Icon name for UI display';

-- +goose Down
DROP TABLE IF EXISTS product_service_subcategories;