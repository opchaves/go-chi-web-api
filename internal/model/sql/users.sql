-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUsers :many
SELECT * FROM users ORDER BY id;

-- name: IsEmailTaken :one
SELECT 1 FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (
  id,
  first_name,
  last_name,
  email,
  password,
  verified,
  verification_token,
  avatar,
  created_at,
  updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
  first_name = coalesce(sqlc.narg('first_name'), first_name),
  last_name = coalesce(sqlc.narg('last_name'), last_name),
  email = coalesce(sqlc.narg('email'), email),
  password = coalesce(sqlc.narg('password'), password),
  verified = coalesce(sqlc.narg('verified'), verified),
  verification_token = coalesce(sqlc.narg('verification_token'), verification_token),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;
