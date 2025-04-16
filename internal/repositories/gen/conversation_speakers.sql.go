// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: conversation_speakers.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const assignParticipantToSpeaker = `-- name: AssignParticipantToSpeaker :exec
UPDATE conversation_speakers SET participant_id = $1 WHERE id = $2
`

type AssignParticipantToSpeakerParams struct {
	ParticipantID *uuid.UUID `json:"participant_id"`
	ID            uuid.UUID  `json:"id"`
}

func (q *Queries) AssignParticipantToSpeaker(ctx context.Context, arg AssignParticipantToSpeakerParams) error {
	_, err := q.db.Exec(ctx, assignParticipantToSpeaker, arg.ParticipantID, arg.ID)
	return err
}

const createNewSpeakerForSegment = `-- name: CreateNewSpeakerForSegment :one
INSERT INTO conversation_speakers (speaker, participant_id, conversation_id) VALUES ($1, $2, $3)
RETURNING id
`

type CreateNewSpeakerForSegmentParams struct {
	Speaker        int32      `json:"speaker"`
	ParticipantID  *uuid.UUID `json:"participant_id"`
	ConversationID uuid.UUID  `json:"conversation_id"`
}

func (q *Queries) CreateNewSpeakerForSegment(ctx context.Context, arg CreateNewSpeakerForSegmentParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createNewSpeakerForSegment, arg.Speaker, arg.ParticipantID, arg.ConversationID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createSpeakerWithConversationsID = `-- name: CreateSpeakerWithConversationsID :one
INSERT INTO conversation_speakers (conversation_id, speaker) VALUES ($1, $2)
RETURNING id
`

type CreateSpeakerWithConversationsIDParams struct {
	ConversationID uuid.UUID `json:"conversation_id"`
	Speaker        int32     `json:"speaker"`
}

func (q *Queries) CreateSpeakerWithConversationsID(ctx context.Context, arg CreateSpeakerWithConversationsIDParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createSpeakerWithConversationsID, arg.ConversationID, arg.Speaker)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getSpeakerCountByConversationID = `-- name: GetSpeakerCountByConversationID :one
SELECT COUNT(*)
FROM conversation_speakers
WHERE conversation_id = $1
`

func (q *Queries) GetSpeakerCountByConversationID(ctx context.Context, conversationID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, getSpeakerCountByConversationID, conversationID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getSpeakerParticipantIDBySegmentID = `-- name: GetSpeakerParticipantIDBySegmentID :one
SELECT cs.participant_id 
FROM conversation_speakers AS cs 
JOIN segments AS s ON cs.id = s.speaker_id
WHERE s.id = $1
`

func (q *Queries) GetSpeakerParticipantIDBySegmentID(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	row := q.db.QueryRow(ctx, getSpeakerParticipantIDBySegmentID, id)
	var participant_id *uuid.UUID
	err := row.Scan(&participant_id)
	return participant_id, err
}
