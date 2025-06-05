-- +goose Up
CREATE TABLE articles (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id),
    title TEXT NOT NULL,
    body JSONB NOT NULL,
    image_url TEXT

);

-- +goose Down
DROP TABLE IF EXISTS articles;