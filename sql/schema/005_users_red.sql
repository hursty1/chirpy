-- +goose Up
ALTER TABLE users add COLUMN is_chirpy_red boolean NOT NULL default false;
UPDATE users set is_chirpy_red = false where is_chirpy_red is NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN is_chirpy_red;