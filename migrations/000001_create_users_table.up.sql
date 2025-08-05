CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    pfp_url VARCHAR(255),
    balance DECIMAL(12, 2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    last_login TIMESTAMP,
    CONSTRAINT uni_users_username UNIQUE (username),
    CONSTRAINT uni_users_email UNIQUE (email)
)