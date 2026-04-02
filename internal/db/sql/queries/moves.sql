-- name: CreateMove :one
INSERT INTO moves (game_id, player_id, move_number, from_x, from_y, to_x, to_y, piece)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetMovesByGameID :many
SELECT * FROM moves
WHERE game_id = $1
ORDER BY move_number ASC;

-- name: GetLastMoveByGameID :one
SELECT * FROM moves
WHERE game_id = $1
ORDER BY move_number DESC
LIMIT 1;