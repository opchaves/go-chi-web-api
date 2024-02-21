-- name: CreateAccount :one
INSERT INTO accounts (
  name, balance, user_id, workspace_id
) VALUES (
  $1, $2, $3, $4
) returning *;

-- name: UpdateAccount :exec
UPDATE accounts SET
  name = sqlc.arg('name'),
  balance = coalesce(sqlc.narg('balance'), balance)
WHERE id = sqlc.arg('id');

-- name: GetAccountByID :one
SELECT * FROM accounts WHERE id = $1;

-- name: GetAccountsByWorkspace :many
SELECT * FROM accounts WHERE workspace_id = $1;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;
