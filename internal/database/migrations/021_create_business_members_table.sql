-- +goose Up
-- Create business_members table for many-to-many relationship
CREATE TABLE business_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'staff',
    is_active BOOLEAN DEFAULT true,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    invited_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(business_id, user_id)
);

-- Create indexes
CREATE INDEX idx_business_members_business_id ON business_members(business_id);
CREATE INDEX idx_business_members_user_id ON business_members(user_id);
CREATE INDEX idx_business_members_role ON business_members(role);
CREATE INDEX idx_business_members_is_active ON business_members(is_active);

-- Composite indexes for common queries
CREATE INDEX idx_business_members_business_user ON business_members(business_id, user_id);
CREATE INDEX idx_business_members_business_role ON business_members(business_id, role);

-- Migrate existing data: Add owners from businesses table
INSERT INTO business_members (business_id, user_id, role, is_active, joined_at)
SELECT id, user_id, 'business_owner', true, created_at
FROM businesses
WHERE user_id IS NOT NULL
ON CONFLICT (business_id, user_id) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS business_members;