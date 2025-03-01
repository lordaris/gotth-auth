-- Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password_hash VARCHAR(100) NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL,
    plaintext_token VARCHAR(255) NOT NULL,
    expiry TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS tokens_token_hash_idx ON tokens(token_hash);
CREATE INDEX IF NOT EXISTS tokens_user_id_idx ON tokens(user_id);

-- Down
--DROP TABLE IF EXISTS tokens;
--DROP TABLE IF EXISTS users;

-- INSERT INTO users (name, email) VALUES 
   -- ('User 1', 'user1@example.com'),
   -- ('User 2', 'user2@example.com'),
    -- ('User 3', 'user3@example.com');

