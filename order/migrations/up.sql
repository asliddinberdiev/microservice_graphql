CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    total_price MONEY NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS order_products (
    order_id UUID REFERENCES orders (id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    quantity INT NOT NULL,
    PRIMARY KEY (order_id, product_id)
);