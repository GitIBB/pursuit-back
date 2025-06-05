-- name: CreateArticle :one
INSERT INTO articles (id, created_at, updated_at, user_id, category_id, title, body, image_url)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetArticles :many
SELECT a.*, users.username
FROM articles a
JOIN users ON a.user_id = users.id
ORDER BY a.created_at ASC
LIMIT $1 OFFSET $2;

-- name: GetArticle :one
Select a.*, users.username
FROM articles a
JOIN users on a.user_id = users.id
WHERE a.id = $1;

-- name: DeleteArticle :exec
DELETE FROM articles
where id = $1;

-- name: GetTotalArticlesCount :one
SELECT COUNT(*) FROM articles;

-- name: GetArticlesByUserId :many
SELECT a.*, users.username
FROM articles a
JOIN users ON a.user_id = users.id
WHERE a.user_id = $1
ORDER BY a.created_at ASC
LIMIT $2 OFFSET $3;