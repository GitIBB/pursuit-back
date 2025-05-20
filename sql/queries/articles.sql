-- name: CreateArticle :one
INSERT INTO articles (id, created_at, updated_at, user_id, title, body, image_url)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetArticles :many
SELECT * FROM articles
ORDER BY created_at ASC
LIMIT $1 OFFSET $2;

-- name: GetArticle :one
Select * FROM articles
WHERE id = $1;

-- name: DeleteArticle :exec
DELETE FROM articles
where id = $1;

-- name: GetTotalArticlesCount :one
SELECT COUNT(*) FROM articles;

-- name: GetArticlesByUserId :many
SELECT * FROM articles
WHERE user_id = $1
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;