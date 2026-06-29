-- name: CreateConversation :one
INSERT INTO conversations (conversation_id, friendship_id, created_at, updated_at)
VALUES (
	gen_random_uuid(),
    $1,
	NOW(),
	NOW()
)
RETURNING *;

-- name: GetConvoByFriendshipID :one
SELECT * FROM conversations
WHERE friendship_id = $1;
