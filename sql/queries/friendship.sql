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

-- name: GetPendingRequests :many
SELECT 
    f.*,
    u.username
FROM friendship f
JOIN users u ON f.sender_id = u.id
WHERE f.receiver_id = $1 AND f.friend_status = 'pending';

-- name: GetFriends :many
SELECT 
    u.id,
    u.username
FROM friendship f
JOIN users u ON (
    (f.sender_id = $1 AND f.receiver_id = u.id) OR 
    (f.receiver_id = $1 AND f.sender_id = u.id)
)
WHERE f.friend_status = 'accepted'
AND u.id != $1;