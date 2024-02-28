-- name: GetWorkspacesByUser :many
SELECT * FROM workspaces WHERE user_id = $1 and deleted_at IS NULL LIMIT $2 OFFSET $3;

-- name: CreateWorkspace :one
INSERT INTO workspaces (
  name,
  description,
  currency,
  language,
  user_id
) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateWorkspace :one
UPDATE workspaces SET
  name = coalesce(sqlc.narg('name'), name),
  description = coalesce(sqlc.narg('description'), description),
  currency = coalesce(sqlc.narg('currency'), currency),
  language = coalesce(sqlc.narg('language'), language)
WHERE id = sqlc.arg('id') and user_id = sqlc.arg('user_id') RETURNING *;

-- name: DeleteWorkspace :exec
UPDATE workspaces SET deleted_at = now() WHERE id = $1 and user_id = $2;
