-- name: CreateRoom :one
INSERT INTO rooms (room_id, room_name, created_at)
VALUES (
	gen_random_uuid(),
    $1,
	NOW()
)
RETURNING *;

-- name: GetRoomByName :one
SELECT *
FROM rooms
WHERE room_name = $1;