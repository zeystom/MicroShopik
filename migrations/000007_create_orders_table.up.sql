CREATE TABLE orders (
                        id          SERIAL PRIMARY KEY,
                        customer_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                        product_id  INTEGER REFERENCES products(id) ON DELETE SET NULL,
                        product_item_id INTEGER REFERENCES product_items(id) ON DELETE SET NULL,
                        status      VARCHAR(20) DEFAULT 'pending',  -- 'pending', 'completed', 'refunded'
                        created_at  TIMESTAMP DEFAULT NOW()
);