-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSessionByToken :one
SELECT *
FROM sessions
WHERE token_hash = $1
LIMIT 1;

-- name: RevokeSession :exec
UPDATE sessions
SET revoked_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: RevokeAllUserSessions :exec
UPDATE sessions
SET revoked_at = NOW(), updated_at = NOW()
WHERE user_id = $1 AND revoked_at IS NULL;

-- name: CreateRecoveryCode :one
INSERT INTO recovery_codes (user_id, code_hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetUnusedRecoveryCodes :many
SELECT *
FROM recovery_codes
WHERE user_id = $1 AND used_at IS NULL
ORDER BY created_at DESC;

-- name: MarkRecoveryCodeUsed :exec
UPDATE recovery_codes
SET used_at = NOW(), updated_at = NOW()
WHERE id = $1;
