CREATE TABLE products (
                          id            SERIAL PRIMARY KEY,
                          seller_id     INTEGER REFERENCES users(id) ON DELETE CASCADE,
                          title         VARCHAR(100) NOT NULL,
                          description   TEXT,
                          price         DECIMAL(10, 2) NOT NULL CHECK (price > 0),
                          category_id   INTEGER REFERENCES categories(id) ON DELETE SET NULL,
                          is_active     BOOLEAN DEFAULT TRUE,
                          disposable    BOOLEAN DEFAULT TRUE,
                          max_sales     INTEGER DEFAULT NULL,
                          sold_count    INTEGER DEFAULT 0,
                          created_at    TIMESTAMP DEFAULT NOW(),
                          updated_at    TIMESTAMP DEFAULT NOW()
);