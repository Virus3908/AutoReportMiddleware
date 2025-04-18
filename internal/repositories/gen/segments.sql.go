// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: segments.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const assignNewSpeakerToSegment = `-- name: AssignNewSpeakerToSegment :exec
UPDATE segments
SET
    speaker_id = $1
WHERE
    id = $2
`

type AssignNewSpeakerToSegmentParams struct {
	SpeakerID uuid.UUID `json:"speaker_id"`
	ID        uuid.UUID `json:"id"`
}

func (q *Queries) AssignNewSpeakerToSegment(ctx context.Context, arg AssignNewSpeakerToSegmentParams) error {
	_, err := q.db.Exec(ctx, assignNewSpeakerToSegment, arg.SpeakerID, arg.ID)
	return err
}

const createSegment = `-- name: CreateSegment :exec
INSERT INTO
    segments (
        diarize_id,
        start_time,
        end_time,
        speaker_id
    )
VALUES ($1, $2, $3, $4)
`

type CreateSegmentParams struct {
	DiarizeID uuid.UUID `json:"diarize_id"`
	StartTime float64   `json:"start_time"`
	EndTime   float64   `json:"end_time"`
	SpeakerID uuid.UUID `json:"speaker_id"`
}

func (q *Queries) CreateSegment(ctx context.Context, arg CreateSegmentParams) error {
	_, err := q.db.Exec(ctx, createSegment,
		arg.DiarizeID,
		arg.StartTime,
		arg.EndTime,
		arg.SpeakerID,
	)
	return err
}

const getCountSegmentsWithSpeakerID = `-- name: GetCountSegmentsWithSpeakerID :one
SELECT COUNT(*)
FROM segments
WHERE speaker_id = $1
`

func (q *Queries) GetCountSegmentsWithSpeakerID(ctx context.Context, speakerID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, getCountSegmentsWithSpeakerID, speakerID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getSegmentsByConversationsID = `-- name: GetSegmentsByConversationsID :many
SELECT 
    seg.id AS segment_id, 
    seg.start_time, 
    seg.end_time, 
    conv.file_url,  
    c.id AS conversation_id
FROM
    segments AS seg
    JOIN diarize AS d ON d.id = seg.diarize_id
    JOIN convert AS conv ON conv.id = d.convert_id
    JOIN conversations AS c ON c.id = conv.conversations_id
WHERE
    c.id = $1
ORDER BY seg.start_time
`

type GetSegmentsByConversationsIDRow struct {
	SegmentID      uuid.UUID `json:"segment_id"`
	StartTime      float64   `json:"start_time"`
	EndTime        float64   `json:"end_time"`
	FileUrl        *string   `json:"file_url"`
	ConversationID uuid.UUID `json:"conversation_id"`
}

func (q *Queries) GetSegmentsByConversationsID(ctx context.Context, id uuid.UUID) ([]GetSegmentsByConversationsIDRow, error) {
	rows, err := q.db.Query(ctx, getSegmentsByConversationsID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSegmentsByConversationsIDRow
	for rows.Next() {
		var i GetSegmentsByConversationsIDRow
		if err := rows.Scan(
			&i.SegmentID,
			&i.StartTime,
			&i.EndTime,
			&i.FileUrl,
			&i.ConversationID,
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
