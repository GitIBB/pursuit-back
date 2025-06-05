-- +goose Up
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS categories;