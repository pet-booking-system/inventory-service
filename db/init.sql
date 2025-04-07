ALTER DATABASE invservice SET timezone TO 'Asia/Almaty';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE resource_status AS ENUM ('available', 'booked', 'unavailable');

CREATE TABLE IF NOT EXISTS resources (
    resource_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    type TEXT,
    status resource_status NOT NULL DEFAULT 'available',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_resources_type ON resources (type);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_resources_updated_at
BEFORE UPDATE ON resources
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
