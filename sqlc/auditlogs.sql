-- name: CreateLogs :one
INSERT INTO audit_logs (
  account_id,
  action,
  action_type,
  action_description
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: ListLogs :many
SELECT * FROM audit_logs ORDER BY account_id LIMIT $1 OFFSET $2;


