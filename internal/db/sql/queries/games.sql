-- name: CreateGame :one
INSERT INTO games (player_white_id, player_black_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetGameByID :one
SELECT * FROM games
WHERE id = $1;

-- name: GetGamesByUserID :many
SELECT * FROM games
WHERE player_white_id = $1 OR player_black_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: StartGame :one
UPDATE games
SET status = 'ongoing',
    started_at = NOW()
WHERE id = $1
RETURNING *;

-- name: FinishGame :one
UPDATE games
SET status     = 'finished',
    winner_id  = $2,
    end_reason = $3,
    finished_at = NOW()
WHERE id = $1
RETURNING *;

-- name: AbandonGame :one
UPDATE games
SET status = 'abandoned',
    finished_at = NOW()
WHERE id = $1
RETURNING *;