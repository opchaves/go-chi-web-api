-- name: GetTransactionById :one
SELECT * FROM transactions WHERE id = $1 and user_id = $2;

-- name: GetTransactionsByUser :many
SELECT * FROM transactions WHERE user_id = $1;

-- name: CreateTransaction :one
INSERT INTO transactions (
  title,
  "note",
  amount,
  paid,
  t_type,
  handled_at,
  workspace_id,
  user_id,
  category_id,
  account_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateTransaction :one
UPDATE transactions SET
  title = coalesce(sqlc.narg('title'), title),
  "note" = coalesce(sqlc.narg('note'), "note"),
  amount = coalesce(sqlc.narg('amount'), amount),
  paid = coalesce(sqlc.narg('paid'), paid),
  t_type = coalesce(sqlc.narg('t_type'), t_type),
  handled_at = coalesce(sqlc.narg('handled_at'), handled_at),
  category_id = coalesce(sqlc.narg('category_id'), category_id),
  account_id = coalesce(sqlc.narg('account_id'), account_id)
WHERE id = sqlc.arg('id') and user_id = sqlc.arg('user_id') RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions WHERE id = $1 and user_id = $2;

