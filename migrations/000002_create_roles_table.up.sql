CREATE TABLE roles(
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL,
    description TEXT
);

INSERT INTO roles (name, description) VALUES
('admin', 'System administrator'),
('seller', 'Goods seller'),
('customer', 'Good buyer');