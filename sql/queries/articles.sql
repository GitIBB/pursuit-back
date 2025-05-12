-- name: CreateArticle :one
INSERT INTO articles (id, created_at, updated_at, user_id, title, body)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetArticles :many
SELECT * FROM articles
ORDER BY created_at ASC;

-- name: GetArticle :one
Select * FROM articles
WHERE id = $1;