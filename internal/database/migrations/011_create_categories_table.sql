-- +goose Up
-- Create categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    icon VARCHAR(100),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index for faster lookups
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_is_active ON categories(is_active);

-- Seed initial categories
INSERT INTO categories (id, name, slug, icon, description, is_active) VALUES
    (gen_random_uuid(), 'All', 'all', '', 'All categories', false),
    (gen_random_uuid(), 'Retail', 'retail', '', 'Retail stores and shops', true),
    (gen_random_uuid(), 'Fashion', 'fashion', '', 'Clothing and accessories', true),
    (gen_random_uuid(), 'Beauty', 'beauty', '', 'Beauty products and cosmetics', true),
    (gen_random_uuid(), 'Electronics', 'electronics', '', 'Electronics and gadgets', true),
    (gen_random_uuid(), 'Food', 'food', '', 'Food and beverages', true),
    (gen_random_uuid(), 'Health', 'health', '', 'Health and wellness products', true);

-- +goose Down
DROP TABLE IF EXISTS categories;