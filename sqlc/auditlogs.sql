-- name: CreateLogs :one
INSERT INTO audit_logs (action) VALUES ($1) RETURNING *;

-- name: ListLogs :many
SELECT * FROM audit_logs ORDER BY account_id LIMIT $1 OFFSET $2;


