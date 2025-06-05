-- name: GetCategories :many
SELECT * FROM categories ORDER BY name ASC;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1;
