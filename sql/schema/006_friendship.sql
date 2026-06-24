-- +goose Up
CREATE TABLE friendship (
    friendship_id UUID PRIMARY KEY,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_status TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CHECK (sender_id != receiver_id),
    CONSTRAINT unique_friendship UNIQUE (sender_id, receiver_id)
);

CREATE INDEX idx_friendship_sender ON friendship(sender_id);
CREATE INDEX idx_friendship_receiver ON friendship(receiver_id);
CREATE INDEX idx_friendship_status ON friendship(friend_status);
CREATE INDEX idx_friendship_updated ON friendship(updated_at);

-- +goose Down
DROP TABLE friendship