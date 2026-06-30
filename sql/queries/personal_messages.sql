-- name: CreatePersonalMessage :one
INSERT INTO personal_messages (message_id, conversation_id, user_id, encrypted_content, message_type, sent_at)
VALUES (
	gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
	NOW()
)
RETURNING *;

-- name: GetRecentConvo :many
SELECT 
    p.*,
    users.username
FROM personal_messages p
JOIN users ON p.user_id = users.id
WHERE p.conversation_id = $1
ORDER BY p.sent_at DESC
LIMIT $2 OFFSET $3;

