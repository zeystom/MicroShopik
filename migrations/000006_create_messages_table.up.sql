CREATE TABLE messages (
                          id            SERIAL PRIMARY KEY,
                          conversation_id INTEGER REFERENCES conversations(id) ON DELETE CASCADE,
                          sender_id     INTEGER REFERENCES users(id) ON DELETE SET NULL,
                          order_id      INTEGER REFERENCES orders(id) ON DELETE SET NULL,
                          text          TEXT,
                          is_system     BOOLEAN DEFAULT FALSE,
                          created_at    TIMESTAMP DEFAULT NOW()
);