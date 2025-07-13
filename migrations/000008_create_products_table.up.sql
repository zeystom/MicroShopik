CREATE TABLE products (
                          id            SERIAL PRIMARY KEY,
                          seller_id     INTEGER REFERENCES users(id) ON DELETE CASCADE,
                          title         VARCHAR(100) NOT NULL,
                          description   TEXT,
                          price         DECIMAL(10, 2) NOT NULL CHECK (price > 0),
                          category      VARCHAR(50),
                          is_active     BOOLEAN DEFAULT TRUE,
                          disposable    BOOLEAN DEFAULT TRUE,
                          max_sales     INTEGER DEFAULT NULL,
                          created_at    TIMESTAMP DEFAULT NOW(),
                          updated_at    TIMESTAMP DEFAULT NOW()
);