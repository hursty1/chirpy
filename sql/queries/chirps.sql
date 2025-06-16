-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
values($1, $2, $3, $4, $5)
returning *;

-- name: GetAllChirps :many
SELECT
*
from chirps
order by created_at asc;


-- name: GetChirpById :one
SELECT
*
from chirps
where id = $1;

-- name: DeleteChripById :exec
DELETE from chirps where id = $1;


-- name: GetAuthorChirps :many
SELECT
*
from chirps
where user_id = $1
ORDER BY created_at asc;