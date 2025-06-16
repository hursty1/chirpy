-- name: AddRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
values ($1, $2, $3, $4, $5)
returning *;

-- name: FetchFreshToken :one
SELECT
*
FROM refresh_tokens
where token = $1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = now(),
updated_at = now()
returning *;