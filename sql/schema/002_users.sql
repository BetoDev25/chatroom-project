-- +goose Up
ALTER TABLE users 
DROP COLUMN email,
ADD COLUMN username TEXT NOT NULL UNIQUE,
ADD COLUMN password TEXT NOT NULL;

-- +goose Down
ALTER TABLE users
ADD COLUMN email TEXT NOT NULL UNIQUE,
DROP COLUMN username,
DROP COLUMN password;

