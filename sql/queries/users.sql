-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
values ($1, $2, $3, $4, $5)
returning *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserFromEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateEmailAndPassword :one
UPDATE users
set email = $1,
hashed_password = $2,
updated_at = now()
where id = $3
returning *;


-- name: UpdateIsRed :one
UPDATE users
set is_chirpy_red = true,
updated_at = now()
where id = $1
returning *;