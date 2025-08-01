// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: sessions.sql

package gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (user_id, refresh_token, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, refresh_token, expires_at, created_at
`

type CreateSessionParams struct {
	UserID       pgtype.UUID
	RefreshToken string
	ExpiresAt    int64
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, createSession, arg.UserID, arg.RefreshToken, arg.ExpiresAt)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteSession = `-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1
`

func (q *Queries) DeleteSession(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteSession, id)
	return err
}

const getSession = `-- name: GetSession :one
SELECT id, user_id, refresh_token, expires_at, created_at FROM sessions WHERE id = $1
`

func (q *Queries) GetSession(ctx context.Context, id int32) (Session, error) {
	row := q.db.QueryRow(ctx, getSession, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}
