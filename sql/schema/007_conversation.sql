-- +goose Up
ALTER TABLE rooms 
ADD COLUMN type TEXT NOT NULL DEFAULT 'public',
ADD COLUMN hashed_password TEXT NOT NULL DEFAULT '';

CREATE TABLE conversations (
    conversation_id UUID PRIMARY KEY,
    friendship_id UUID NOT NULL REFERENCES friendship(friendship_id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_friendship_conversation UNIQUE (friendship_id)
);

CREATE TABLE personal_messages (
    message_id UUID PRIMARY KEY,
    conversation_id UUID NOT NULL REFERENCES conversations(conversation_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    encrypted_content TEXT NOT NULL,
    message_type TEXT NOT NULL DEFAULT 'text',
    sent_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE conversations;
DROP TABLE personal_messages;
ALTER TABLE rooms DROP COLUMN type;
ALTER TABLE rooms DROP COLUMN hashed_password;
