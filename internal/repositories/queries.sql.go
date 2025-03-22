// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package repositories

import (
	"context"

	"github.com/google/uuid"
)

const CreateConversation = `-- name: CreateConversation :exec
INSERT INTO Conversations (conversation_name, file_url) VALUES ($1, $2)
`

type CreateConversationParams struct {
	ConversationName string `db:"conversation_name" json:"conversation_name"`
	FileUrl          string `db:"file_url" json:"file_url"`
}

func (q *Queries) CreateConversation(ctx context.Context, arg CreateConversationParams) error {
	_, err := q.db.Exec(ctx, CreateConversation, arg.ConversationName, arg.FileUrl)
	return err
}

const CreateConvertTask = `-- name: CreateConvertTask :exec
INSERT INTO Convert (conversations_id, task_id) 
VALUES ($1, $2)
ON CONFLICT (conversations_id) DO UPDATE
SET task_id = $2
`

type CreateConvertTaskParams struct {
	ConversationsID uuid.UUID `db:"conversations_id" json:"conversations_id"`
	TaskID          uuid.UUID `db:"task_id" json:"task_id"`
}

func (q *Queries) CreateConvertTask(ctx context.Context, arg CreateConvertTaskParams) error {
	_, err := q.db.Exec(ctx, CreateConvertTask, arg.ConversationsID, arg.TaskID)
	return err
}

const CreateDiarizeTask = `-- name: CreateDiarizeTask :exec
INSERT INTO diarize (conversation_id, task_id)
VALUES ($1, $2)
ON CONFLICT (conversation_id) DO UPDATE
SET task_id = $2
`

type CreateDiarizeTaskParams struct {
	ConversationID uuid.UUID `db:"conversation_id" json:"conversation_id"`
	TaskID         uuid.UUID `db:"task_id" json:"task_id"`
}

func (q *Queries) CreateDiarizeTask(ctx context.Context, arg CreateDiarizeTaskParams) error {
	_, err := q.db.Exec(ctx, CreateDiarizeTask, arg.ConversationID, arg.TaskID)
	return err
}

const CreateParticipant = `-- name: CreateParticipant :exec
INSERT INTO Participants (name, email) VALUES ($1, $2)
`

type CreateParticipantParams struct {
	Name  *string `db:"name" json:"name"`
	Email string  `db:"email" json:"email"`
}

func (q *Queries) CreateParticipant(ctx context.Context, arg CreateParticipantParams) error {
	_, err := q.db.Exec(ctx, CreateParticipant, arg.Name, arg.Email)
	return err
}

const CreatePromt = `-- name: CreatePromt :exec
INSERT INTO Promts (promt) VALUES ($1)
`

func (q *Queries) CreatePromt(ctx context.Context, promt string) error {
	_, err := q.db.Exec(ctx, CreatePromt, promt)
	return err
}

const CreateSegments = `-- name: CreateSegments :exec
INSERT INTO segments (conversation_id, start_time, end_time, speaker)
VALUES ($1, $2, $3, $4)
`

type CreateSegmentsParams struct {
	ConversationID uuid.UUID `db:"conversation_id" json:"conversation_id"`
	StartTime      float64   `db:"start_time" json:"start_time"`
	EndTime        float64   `db:"end_time" json:"end_time"`
	Speaker        int32     `db:"speaker" json:"speaker"`
}

func (q *Queries) CreateSegments(ctx context.Context, arg CreateSegmentsParams) error {
	_, err := q.db.Exec(ctx, CreateSegments,
		arg.ConversationID,
		arg.StartTime,
		arg.EndTime,
		arg.Speaker,
	)
	return err
}

const DeleteConversationByID = `-- name: DeleteConversationByID :exec
DELETE FROM Conversations WHERE id = $1
`

func (q *Queries) DeleteConversationByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, DeleteConversationByID, id)
	return err
}

const DeleteParticipantByID = `-- name: DeleteParticipantByID :exec
DELETE FROM Participants WHERE id = $1
`

func (q *Queries) DeleteParticipantByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, DeleteParticipantByID, id)
	return err
}

const DeletePromtByID = `-- name: DeletePromtByID :exec
DELETE FROM Promts WHERE id = $1
`

func (q *Queries) DeletePromtByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, DeletePromtByID, id)
	return err
}

const GetConversationByID = `-- name: GetConversationByID :one
SELECT id, conversation_name, file_url, status, created_at, updated_at FROM Conversations WHERE id = $1
`

func (q *Queries) GetConversationByID(ctx context.Context, id uuid.UUID) (Conversation, error) {
	row := q.db.QueryRow(ctx, GetConversationByID, id)
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

const GetConversationFileURL = `-- name: GetConversationFileURL :one
SELECT file_url FROM conversations WHERE id = $1
`

func (q *Queries) GetConversationFileURL(ctx context.Context, id uuid.UUID) (string, error) {
	row := q.db.QueryRow(ctx, GetConversationFileURL, id)
	var file_url string
	err := row.Scan(&file_url)
	return file_url, err
}

const GetConversationIDByDiarizeTaskID = `-- name: GetConversationIDByDiarizeTaskID :one
SELECT conversation_id FROM diarize WHERE task_id = $1
`

func (q *Queries) GetConversationIDByDiarizeTaskID(ctx context.Context, taskID uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, GetConversationIDByDiarizeTaskID, taskID)
	var conversation_id uuid.UUID
	err := row.Scan(&conversation_id)
	return conversation_id, err
}

const GetConversations = `-- name: GetConversations :many
SELECT id, conversation_name, file_url, status, created_at, updated_at FROM Conversations
`

func (q *Queries) GetConversations(ctx context.Context) ([]Conversation, error) {
	rows, err := q.db.Query(ctx, GetConversations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Conversation{}
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

const GetConvertFileURL = `-- name: GetConvertFileURL :one
SELECT file_url FROM convert WHERE conversations_id = $1
`

func (q *Queries) GetConvertFileURL(ctx context.Context, conversationsID uuid.UUID) (*string, error) {
	row := q.db.QueryRow(ctx, GetConvertFileURL, conversationsID)
	var file_url *string
	err := row.Scan(&file_url)
	return file_url, err
}

const GetParticipantByID = `-- name: GetParticipantByID :one
SELECT id, name, email, created_at, updated_at FROM Participants WHERE id = $1
`

func (q *Queries) GetParticipantByID(ctx context.Context, id uuid.UUID) (Participant, error) {
	row := q.db.QueryRow(ctx, GetParticipantByID, id)
	var i Participant
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetParticipants = `-- name: GetParticipants :many
SELECT id, name, email, created_at, updated_at FROM Participants
`

func (q *Queries) GetParticipants(ctx context.Context) ([]Participant, error) {
	rows, err := q.db.Query(ctx, GetParticipants)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Participant{}
	for rows.Next() {
		var i Participant
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
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

const GetPromtByID = `-- name: GetPromtByID :one
SELECT id, promt, created_at, updated_at FROM Promts WHERE id = $1
`

func (q *Queries) GetPromtByID(ctx context.Context, id uuid.UUID) (Promt, error) {
	row := q.db.QueryRow(ctx, GetPromtByID, id)
	var i Promt
	err := row.Scan(
		&i.ID,
		&i.Promt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetPromts = `-- name: GetPromts :many
SELECT id, promt, created_at, updated_at FROM Promts
`

func (q *Queries) GetPromts(ctx context.Context) ([]Promt, error) {
	rows, err := q.db.Query(ctx, GetPromts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Promt{}
	for rows.Next() {
		var i Promt
		if err := rows.Scan(
			&i.ID,
			&i.Promt,
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

const UpdateConversationNameByID = `-- name: UpdateConversationNameByID :exec
UPDATE Conversations SET conversation_name = $1 WHERE id = $2
`

type UpdateConversationNameByIDParams struct {
	ConversationName string    `db:"conversation_name" json:"conversation_name"`
	ID               uuid.UUID `db:"id" json:"id"`
}

func (q *Queries) UpdateConversationNameByID(ctx context.Context, arg UpdateConversationNameByIDParams) error {
	_, err := q.db.Exec(ctx, UpdateConversationNameByID, arg.ConversationName, arg.ID)
	return err
}

const UpdateConvertTask = `-- name: UpdateConvertTask :exec
UPDATE Convert SET file_url = $1, audio_len = $2 WHERE task_id = $3
`

type UpdateConvertTaskParams struct {
	FileUrl  *string   `db:"file_url" json:"file_url"`
	AudioLen *float64  `db:"audio_len" json:"audio_len"`
	TaskID   uuid.UUID `db:"task_id" json:"task_id"`
}

func (q *Queries) UpdateConvertTask(ctx context.Context, arg UpdateConvertTaskParams) error {
	_, err := q.db.Exec(ctx, UpdateConvertTask, arg.FileUrl, arg.AudioLen, arg.TaskID)
	return err
}

const UpdateParticipantByID = `-- name: UpdateParticipantByID :exec
UPDATE Participants SET name = $1, email = $2 WHERE id = $3
`

type UpdateParticipantByIDParams struct {
	Name  *string   `db:"name" json:"name"`
	Email string    `db:"email" json:"email"`
	ID    uuid.UUID `db:"id" json:"id"`
}

func (q *Queries) UpdateParticipantByID(ctx context.Context, arg UpdateParticipantByIDParams) error {
	_, err := q.db.Exec(ctx, UpdateParticipantByID, arg.Name, arg.Email, arg.ID)
	return err
}

const UpdatePromtByID = `-- name: UpdatePromtByID :exec
UPDATE Promts SET promt = $1 WHERE id = $2
`

type UpdatePromtByIDParams struct {
	Promt string    `db:"promt" json:"promt"`
	ID    uuid.UUID `db:"id" json:"id"`
}

func (q *Queries) UpdatePromtByID(ctx context.Context, arg UpdatePromtByIDParams) error {
	_, err := q.db.Exec(ctx, UpdatePromtByID, arg.Promt, arg.ID)
	return err
}
