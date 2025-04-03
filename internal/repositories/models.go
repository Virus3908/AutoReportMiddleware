// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repositories

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID               uuid.UUID `db:"id" json:"id"`
	ConversationName string    `db:"conversation_name" json:"conversation_name"`
	FileUrl          string    `db:"file_url" json:"file_url"`
	Status           int32     `db:"status" json:"status"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type ConversationsParticipant struct {
	ID             uuid.UUID `db:"id" json:"id"`
	UserID         uuid.UUID `db:"user_id" json:"user_id"`
	Speaker        *int32    `db:"speaker" json:"speaker"`
	ConversationID uuid.UUID `db:"conversation_id" json:"conversation_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type Convert struct {
	ID              uuid.UUID `db:"id" json:"id"`
	ConversationsID uuid.UUID `db:"conversations_id" json:"conversations_id"`
	FileUrl         *string   `db:"file_url" json:"file_url"`
	AudioLen        *float64  `db:"audio_len" json:"audio_len"`
	TaskID          uuid.UUID `db:"task_id" json:"task_id"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type Diarize struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ConverID  uuid.UUID `db:"conver_id" json:"conver_id"`
	TaskID    uuid.UUID `db:"task_id" json:"task_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Participant struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      *string   `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Prompt struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Prompt    string    `db:"prompt" json:"prompt"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Report struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	ConversationID uuid.UUID  `db:"conversation_id" json:"conversation_id"`
	Report         *string    `db:"report" json:"report"`
	PromptID       *uuid.UUID `db:"prompt_id" json:"prompt_id"`
	TaskID         uuid.UUID  `db:"task_id" json:"task_id"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}

type Segment struct {
	ID        uuid.UUID `db:"id" json:"id"`
	DiarizeID uuid.UUID `db:"diarize_id" json:"diarize_id"`
	StartTime float64   `db:"start_time" json:"start_time"`
	EndTime   float64   `db:"end_time" json:"end_time"`
	Speaker   int32     `db:"speaker" json:"speaker"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Task struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Status    int32     `db:"status" json:"status"`
	TaskType  int32     `db:"task_type" json:"task_type"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Transcription struct {
	ID            uuid.UUID `db:"id" json:"id"`
	SegmentID     uuid.UUID `db:"segment_id" json:"segment_id"`
	Transcription *string   `db:"transcription" json:"transcription"`
	TaskID        uuid.UUID `db:"task_id" json:"task_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
