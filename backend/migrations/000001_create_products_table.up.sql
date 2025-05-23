-- Add UUID extension if not already installed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    stripe_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image_url TEXT,
    origin VARCHAR(255),
    roast_level VARCHAR(50),
    stock_level INTEGER DEFAULT 0,
    weight INTEGER DEFAULT 1,
    flavor_notes TEXT,
    active BOOLEAN DEFAULT TRUE,
    options JSONB DEFAULT '{}',
    allow_subscription BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT NOW (),
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT NOW ()
);
