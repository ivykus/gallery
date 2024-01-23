CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    token_hash TEXT NOT NULL UNIQUE
);