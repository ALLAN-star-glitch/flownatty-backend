-- +goose Up
ALTER TABLE business_subcategories ADD COLUMN IF NOT EXISTS display_name VARCHAR(255);
UPDATE business_subcategories SET display_name = name WHERE display_name IS NULL;
ALTER TABLE business_subcategories ALTER COLUMN display_name SET NOT NULL;

-- +goose Down
ALTER TABLE business_subcategories DROP COLUMN IF EXISTS display_name;