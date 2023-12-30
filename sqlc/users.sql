-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email,
    phone_number,
    password_changed_at,
    created_at,
    last_login_at,
    login_attempts,
    locked_until,
    reset_token,
    reset_token_expires_at
    ) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
    )
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE username = $1 LIMIT 1
FOR UPDATE;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: UpdateUser :exec
UPDATE users
set hashed_password = $2,
    full_name = $3,
    email = $4,
    phone_number = $5,
    password_changed_at = $6,
    created_at = $7,
    last_login_at = $8,
    login_attempts = $9,
    locked_until = $10,
    reset_token = $11,
    reset_token_expires_at = $12
WHERE username = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;

