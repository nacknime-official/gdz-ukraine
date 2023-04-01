-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;
-- name: GetUserByTelegramID :one
SELECT *
FROM users
WHERE telegram_id = $1;