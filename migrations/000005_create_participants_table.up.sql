CREATE TABLE participants (
                              conversation_id INTEGER REFERENCES conversations(id) ON DELETE CASCADE,
                              user_id         INTEGER REFERENCES users(id) ON DELETE CASCADE,
                              PRIMARY KEY (conversation_id, user_id)
);