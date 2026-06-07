-- +goose Up
CREATE TABLE rooms (
    room_id UUID PRIMARY KEY,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    room_name TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL
);

CREATE TABLE message (
    message_id UUID PRIMARY KEY,
    room_id UUID NOT NULL REFERENCES rooms(room_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    message_type TEXT NOT NULL DEFAULT 'text', --for future implementation of image/video
    sent_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_room_history_room_sent ON message(room_id, sent_at);

-- +goose Down
DROP TABLE message;
DROP TABLE rooms;