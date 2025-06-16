-- +goose Up
CREATE TABLE users(
id uuid primary key,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
email TEXT not null UNIQUE
);

-- +goose Down
DROP TABLE users;
