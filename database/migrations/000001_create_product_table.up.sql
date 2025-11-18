CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255),
    type VARCHAR(100),
    flavor VARCHAR(100),
    size VARCHAR(50),
    price NUMERIC(12,2),
    stock INTEGER,
    production_date TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
