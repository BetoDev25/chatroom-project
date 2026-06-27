-- name: CreateConversation :one
INSERT INTO conversations (conversation_id, friendship_id, created_at, updated_at)
VALUES (
	gen_random_uuid(),
    $1,
	NOW(),
	NOW()
)
RETURNING *;

-- name: GetConvoBetweenUsers :one
SELECT c.*
FROM conversations c
JOIN friendship f ON c.friendship_id = f.friendship_id
WHERE (f.sender_id = $1 AND f.receiver_id = $2)
   OR (f.sender_id = $2 AND f.receiver_id = $1);

