-- +goose Up
ALTER TABLE friendship DROP CONSTRAINT IF EXISTS friendship_friend_status_key;

-- +goose Down
ALTER TABLE friendship ADD CONSTRAINT friendship_friend_status_key UNIQUE (friend_status);