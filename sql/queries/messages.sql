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
SELECT *
FROM message
WHERE room_id = $1
ORDER BY sent_at ASC
LIMIT 50;

