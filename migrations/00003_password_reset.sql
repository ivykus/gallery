-- +goose Up
-- +goose StatementBegin
CREATE TABLE passwords_reset (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT UNIQUE NOT NULL, 
    expires_at TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE passwords_reset;
-- +goose StatementEnd


