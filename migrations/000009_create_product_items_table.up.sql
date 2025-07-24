CREATE TABLE product_items (
                               id            SERIAL PRIMARY KEY,
                               product_id    INTEGER REFERENCES products(id) ON DELETE CASCADE,
                               data          TEXT NOT NULL,
                               is_used       BOOLEAN DEFAULT FALSE,

                               created_at    TIMESTAMP DEFAULT NOW()
);