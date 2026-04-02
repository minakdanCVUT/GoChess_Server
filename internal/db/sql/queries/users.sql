-- name: CreateUser :one
INSERT INTO users (first_name, last_name, username, email, password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: UsernameExists :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE username = $1
) AS exists;

-- name: EmailExists :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE email = $1
) AS EXISTS;

-- name: VerifyEmail :one
UPDATE users
SET email_verified = TRUE,
    email_verified_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetUserByLogin :one
select id, username, password from users
where email = $1 or username = $1
limit 1;

-- name: GetUserStats :one
SELECT
    COUNT(*) AS total_games,
    COUNT(*) FILTER (WHERE winner_id = $1) AS wins,
    COUNT(*) FILTER (WHERE winner_id != $1 AND winner_id IS NOT NULL) AS losses,
    COUNT(*) FILTER (WHERE winner_id IS NULL AND status = 'finished') AS draws
FROM games
WHERE player_white_id = $1 OR player_black_id = $1;