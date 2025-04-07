// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: conversations.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createConversation = `-- name: CreateConversation :exec
INSERT INTO
    Conversations (conversation_name, file_url)
VALUES ($1, $2)
`

type CreateConversationParams struct {
	ConversationName string `json:"conversation_name"`
	FileUrl          string `json:"file_url"`
}

func (q *Queries) CreateConversation(ctx context.Context, arg CreateConversationParams) error {
	_, err := q.db.Exec(ctx, createConversation, arg.ConversationName, arg.FileUrl)
	return err
}

const deleteConversationByID = `-- name: DeleteConversationByID :one
DELETE FROM Conversations WHERE id = $1 RETURNING file_url
`

func (q *Queries) DeleteConversationByID(ctx context.Context, id uuid.UUID) (string, error) {
	row := q.db.QueryRow(ctx, deleteConversationByID, id)
	var file_url string
	err := row.Scan(&file_url)
	return file_url, err
}

const getConversationByID = `-- name: GetConversationByID :one
SELECT id, conversation_name, file_url, status, created_at, updated_at FROM Conversations WHERE id = $1
`

func (q *Queries) GetConversationByID(ctx context.Context, id uuid.UUID) (Conversation, error) {
	row := q.db.QueryRow(ctx, getConversationByID, id)
	var i Conversation
	err := row.Scan(
		&i.ID,
		&i.ConversationName,
		&i.FileUrl,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getConversationFileURL = `-- name: GetConversationFileURL :one
SELECT file_url FROM conversations WHERE id = $1
`

func (q *Queries) GetConversationFileURL(ctx context.Context, id uuid.UUID) (string, error) {
	row := q.db.QueryRow(ctx, getConversationFileURL, id)
	var file_url string
	err := row.Scan(&file_url)
	return file_url, err
}

const getConversations = `-- name: GetConversations :many
SELECT id, conversation_name, file_url, status, created_at, updated_at FROM Conversations
`

func (q *Queries) GetConversations(ctx context.Context) ([]Conversation, error) {
	rows, err := q.db.Query(ctx, getConversations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Conversation
	for rows.Next() {
		var i Conversation
		if err := rows.Scan(
			&i.ID,
			&i.ConversationName,
			&i.FileUrl,
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

const updateConversationStatusByConvertID = `-- name: UpdateConversationStatusByConvertID :exec
UPDATE conversations
SET
    status = $1
FROM convert
WHERE
    convert.id = $2
    AND convert.conversations_id = conversations.id
`

type UpdateConversationStatusByConvertIDParams struct {
	Status int32     `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) UpdateConversationStatusByConvertID(ctx context.Context, arg UpdateConversationStatusByConvertIDParams) error {
	_, err := q.db.Exec(ctx, updateConversationStatusByConvertID, arg.Status, arg.ID)
	return err
}

const updateConversationStatusByDiarizeID = `-- name: UpdateConversationStatusByDiarizeID :exec
UPDATE conversations
SET status = $1
FROM convert, diarize
WHERE
    diarize.id = $2
    AND diarize.convert_id = convert.id
    AND convert.conversations_id = conversations.id
`

type UpdateConversationStatusByDiarizeIDParams struct {
	Status int32     `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) UpdateConversationStatusByDiarizeID(ctx context.Context, arg UpdateConversationStatusByDiarizeIDParams) error {
	_, err := q.db.Exec(ctx, updateConversationStatusByDiarizeID, arg.Status, arg.ID)
	return err
}

const updateConversationStatusByID = `-- name: UpdateConversationStatusByID :exec
UPDATE conversations SET status = $1 WHERE id = $2
`

type UpdateConversationStatusByIDParams struct {
	Status int32     `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) UpdateConversationStatusByID(ctx context.Context, arg UpdateConversationStatusByIDParams) error {
	_, err := q.db.Exec(ctx, updateConversationStatusByID, arg.Status, arg.ID)
	return err
}
