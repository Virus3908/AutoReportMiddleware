// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: convert.sql

package repositories

import (
	"context"

	"github.com/google/uuid"
)

const CreateConvert = `-- name: CreateConvert :one
INSERT INTO convert (conversations_id) VALUES ($1)
RETURNING id
`

func (q *Queries) CreateConvert(ctx context.Context, conversationsID uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, CreateConvert, conversationsID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const DeleteConvertByForgeinID = `-- name: DeleteConvertByForgeinID :one
DELETE FROM convert 
WHERE conversations_id = $1
RETURNING id
`

func (q *Queries) DeleteConvertByForgeinID(ctx context.Context, conversationsID uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, DeleteConvertByForgeinID, conversationsID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const DeleteConvertByID = `-- name: DeleteConvertByID :exec
DELETE FROM convert WHERE id = $1
`

func (q *Queries) DeleteConvertByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, DeleteConvertByID, id)
	return err
}

const GetConvert = `-- name: GetConvert :many
SELECT id, conversations_id, file_url, audio_len, status, created_at, updated_at FROM convert
`

func (q *Queries) GetConvert(ctx context.Context) ([]Convert, error) {
	rows, err := q.db.Query(ctx, GetConvert)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Convert{}
	for rows.Next() {
		var i Convert
		if err := rows.Scan(
			&i.ID,
			&i.ConversationsID,
			&i.FileUrl,
			&i.AudioLen,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetConvertByID = `-- name: GetConvertByID :one
SELECT id, conversations_id, file_url, audio_len, status, created_at, updated_at FROM convert WHERE id = $1
`

func (q *Queries) GetConvertByID(ctx context.Context, id uuid.UUID) (Convert, error) {
	row := q.db.QueryRow(ctx, GetConvertByID, id)
	var i Convert
	err := row.Scan(
		&i.ID,
		&i.ConversationsID,
		&i.FileUrl,
		&i.AudioLen,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const UpdateConvertByTaskID = `-- name: UpdateConvertByTaskID :exec
UPDATE convert SET file_url = $1, status = $2 WHERE id = $3
`

type UpdateConvertByTaskIDParams struct {
	FileUrl *string   `db:"file_url" json:"file_url"`
	Status  int32     `db:"status" json:"status"`
	ID      uuid.UUID `db:"id" json:"id"`
}

func (q *Queries) UpdateConvertByTaskID(ctx context.Context, arg UpdateConvertByTaskIDParams) error {
	_, err := q.db.Exec(ctx, UpdateConvertByTaskID, arg.FileUrl, arg.Status, arg.ID)
	return err
}
