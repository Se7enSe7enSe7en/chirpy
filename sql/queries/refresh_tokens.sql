-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(
        token,
        created_at,
        updated_at,
        user_id,
        expires_at,
        revoked_at
    )
VALUES (
        $1,
        now(),
        now(),
        $2,
        $3,
        $4
    )
RETURNING *;
-- name: GetRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1
LIMIT 1;
-- name: GetUserFromRefreshToken :one
SELECT u.*
FROM refresh_tokens r
    JOIN users u ON r.user_id = u.id
WHERE r.token = $1
LIMIT 1;
-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = $2,
    updated_at = now()
WHERE token = $1
RETURNING *;
-- name: ResetRefreshTokensTable :exec
DELETE FROM refresh_tokens;