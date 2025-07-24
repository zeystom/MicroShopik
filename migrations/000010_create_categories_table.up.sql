CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(50) NOT NULL UNIQUE,
                            created_at TIMESTAMP DEFAULT NOW()
);
