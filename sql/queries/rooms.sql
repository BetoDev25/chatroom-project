-- name: CreateRoom :one
INSERT INTO rooms (room_id, owner_id, room_name, created_at)
VALUES (
	gen_random_uuid(),
    $1,
	$2,
	NOW()
)
RETURNING *;

-- name: DeleteRoom :exec
DELETE FROM rooms WHERE room_id = $1;

-- name: GetRoomByName :one
SELECT *
FROM rooms
WHERE room_name = $1;

-- name: GetRooms :many
SELECT *
FROM rooms
WHERE owner_id = $1;