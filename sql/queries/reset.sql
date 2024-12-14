-- name: Reset :exec
DELETE FROM users
USING chirps;
