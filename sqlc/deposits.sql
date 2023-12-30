-- name: CreateDeposit :one
INSERT INTO deposits (
  account_id,
  amount,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4
)

RETURNING *;

