-- Deploy plant-nursery:2_plant to pg

BEGIN;

CREATE TABLE IF NOT EXISTS plants (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    image_url VARCHAR(100) NOT NULL,
    stock INTEGER NOT NULL,
    price NUMERIC NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by VARCHAR(100) DEFAULT 'system' NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by VARCHAR(100)
);

COMMIT;
