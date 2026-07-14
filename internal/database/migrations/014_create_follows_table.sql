-- +goose Up
-- Create follows table for consumer-business relationships
CREATE TABLE follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    business_id UUID NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, business_id)
);

-- Create indexes
CREATE INDEX idx_follows_user_id ON follows(user_id);
CREATE INDEX idx_follows_business_id ON follows(business_id);

-- Create composite index for checking follow status
CREATE INDEX idx_follows_user_business ON follows(user_id, business_id);

-- +goose Down
DROP TABLE IF EXISTS follows;