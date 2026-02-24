-- name: CreateCategory :one
INSERT INTO categories (
  id, parent_id, name, slug, description
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name;

-- name: ListSubcategories :many
SELECT * FROM categories
WHERE parent_id = $1
ORDER BY name;

-- name: UpdateCategory :one
UPDATE categories
SET parent_id = $1, name = $2, slug = $3, description = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $5
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;