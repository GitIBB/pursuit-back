-- +goose Up
CREATE TABLE articles (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    TITLE TEXT NOT NULL,
    body TEXT NOT NULL
);

-- +goose Down
DROP TABLE articles;