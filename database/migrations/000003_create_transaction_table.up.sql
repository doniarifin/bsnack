CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(100) PRIMARY KEY,
    is_new_customer BOOLEAN,
    customer_id VARCHAR(100) NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    product_id VARCHAR(100) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_size VARCHAR(100) NOT NULL,
    flavor VARCHAR(100) NOT NULL,
    quantity INTEGER NOT NULL,
    total_price NUMERIC(12,2),
    transaction_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);