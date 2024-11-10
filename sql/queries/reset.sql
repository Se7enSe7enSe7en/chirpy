-- name: ResetUserTable :exec
DELETE FROM users;
-- name: ResetChirpsTable :exec
DELETE FROM chirps;
-- name: ResetRefreshTokensTable :exec
DELETE FROM refresh_tokens;