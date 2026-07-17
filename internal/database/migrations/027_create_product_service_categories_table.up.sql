-- internal/database/migrations/027_create_product_service_categories_table.up.sql
-- +goose Up
-- Product/Service Categories
CREATE TABLE IF NOT EXISTS product_service_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    type VARCHAR(20) NOT NULL, -- 'product' or 'service'
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    UNIQUE(name, type)
);

CREATE INDEX IF NOT EXISTS idx_product_service_categories_type ON product_service_categories(type);
CREATE INDEX IF NOT EXISTS idx_product_service_categories_is_active ON product_service_categories(is_active);
CREATE INDEX IF NOT EXISTS idx_product_service_categories_deleted_at ON product_service_categories(deleted_at);

COMMENT ON TABLE product_service_categories IS 'Categories for products and services';
COMMENT ON COLUMN product_service_categories.name IS 'Category name (e.g., Fashion, Electronics)';
COMMENT ON COLUMN product_service_categories.type IS 'Type: product or service';
COMMENT ON COLUMN product_service_categories.icon IS 'Icon name for UI display';
COMMENT ON COLUMN product_service_categories.sort_order IS 'Order for display in UI';

-- +goose Down
DROP TABLE IF EXISTS product_service_categories;