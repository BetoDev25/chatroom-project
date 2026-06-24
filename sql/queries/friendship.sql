-- name: CreateFriendRequest :one
INSERT INTO friendship (friendship_id, sender_id, receiver_id, friend_status, created_at, updated_at)
VALUES (
	gen_random_uuid(),
    $1,
    $2,
    'pending',
	NOW(),
	NOW()
)
RETURNING *;

-- name: UpdateFriendStatus :exec
UPDATE friendship
SET friend_status = $2, updated_at = NOW()
WHERE friendship_id = $1;

-- name: DeleteFriendship :exec
DELETE FROM friendship WHERE friendship_id = $1;