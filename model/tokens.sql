-- name: GetTokenById :one
SELECT * FROM tokens WHERE id = $1;

-- name: GetTokensByUser :many
SELECT * FROM tokens WHERE user_id = $1;

-- name: GetToken :one
SELECT * FROM tokens WHERE token = $1;

-- name: CreateToken :one
INSERT INTO tokens (
  token,
  identifier,
  mobile,
  expires_at,
  user_id
) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateToken :exec
UPDATE tokens SET
  token = sqlc.arg('token'),
  expires_at = sqlc.arg('expires_at'),
  updated_at = now()
WHERE id = sqlc.arg('id');

-- name: DeleteTokenByID :exec
DELETE FROM tokens WHERE id = $1;
