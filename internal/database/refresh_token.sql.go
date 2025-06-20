// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: refresh_token.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const addRefreshToken = `-- name: AddRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
values ($1, $2, $3, $4, $5)
returning token, created_at, updated_at, user_id, expires_at, revoked_at
`

type AddRefreshTokenParams struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func (q *Queries) AddRefreshToken(ctx context.Context, arg AddRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, addRefreshToken,
		arg.Token,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.ExpiresAt,
	)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const fetchFreshToken = `-- name: FetchFreshToken :one
SELECT
token, created_at, updated_at, user_id, expires_at, revoked_at
FROM refresh_tokens
where token = $1
`

func (q *Queries) FetchFreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, fetchFreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = now(),
updated_at = now()
returning token, created_at, updated_at, user_id, expires_at, revoked_at
`

func (q *Queries) RevokeRefreshToken(ctx context.Context) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, revokeRefreshToken)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}
