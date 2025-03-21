-- name: StoreRefreshToken :one
INSERT INTO refresh_tokens (
    token,
    created_at,
    updated_at,
    user_id,
    expires_at,
    revoked_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    NULL
)
RETURNING *;

-- name: GetUserByRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET
    updated_at = now(),
    revoked_at = now()
WHERE token = $1;