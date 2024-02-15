-- name: GetWorkspacesByUser :many
SELECT * FROM workspaces WHERE user_id = $1 LIMIT $2 OFFSET $3;

-- name: CreateWorkspace :one
INSERT INTO workspaces (
  name,
  description,
  currency,
  language,
  user_id
) VALUES ($1, $2, $3, $4, $5) RETURNING *;
