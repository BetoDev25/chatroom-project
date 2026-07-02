-- name: CreateMessage :one
INSERT INTO message (message_id, room_id, user_id, content, message_type, sent_at)
VALUES (
	gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
	NOW()
)
RETURNING *;

-- name: GetRecentMessages :many
SELECT 
    message.*,
    users.username
FROM message
JOIN users ON message.user_id = users.id
WHERE message.room_id = $1
ORDER BY message.sent_at DESC
LIMIT $2 OFFSET $3;

