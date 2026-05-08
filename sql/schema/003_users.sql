-- +goose Up
ALTER TABLE users
RENAME COLUMN password TO hashed_password;

-- +goose Down
ALTER TABLE users
RENAME COLUMN hashed_password TO password;