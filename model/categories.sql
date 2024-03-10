-- name: CreateCategory :exec
INSERT INTO categories (
  name, user_id, workspace_id, cat_type, description, icon, color
) VALUES (
  $1,
  $2,
  $3,
  $4,
  sqlc.narg('description'),
  sqlc.narg('icon'),
  sqlc.narg('color')
);

-- name: UpdateCategory :exec
UPDATE categories SET
  name = sqlc.arg('name'),
  description = sqlc.narg('description'),
  icon = sqlc.narg('icon'),
  color = sqlc.narg('color')
WHERE id = sqlc.arg('id');

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1;

-- name: GetCategoriesByWorkspace :many
SELECT * FROM categories WHERE workspace_id = $1;

-- name: GetCategoriesByUser :many
SELECT * FROM categories WHERE user_id = $1;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: DeleteCategoriesByWorkspace :exec
DELETE FROM categories WHERE workspace_id = $1;
