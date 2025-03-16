package models

import (
	"time"
	"github.com/google/uuid"
)

type Transcribe struct {
	ID            uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SegmentID     uuid.UUID `json:"segment_id"`
	Transcription string    `json:"transcription"`
	TaskID        uuid.UUID `json:"task_id"`
	Complete      bool      `json:"complete"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
