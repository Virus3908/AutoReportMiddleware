// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"time"

	"github.com/google/uuid"
	"main/internal/models"
)

type Conversation struct {
	ID               uuid.UUID                 `json:"id"`
	ConversationName string                    `json:"conversation_name"`
	FileUrl          string                    `json:"file_url"`
	Status           models.ConversationStatus `json:"status"`
	Processed        bool                      `json:"processed"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
}

type ConversationSpeaker struct {
	ID             uuid.UUID  `json:"id"`
	ParticipantID  *uuid.UUID `json:"participant_id"`
	Speaker        int32      `json:"speaker"`
	ConversationID uuid.UUID  `json:"conversation_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Convert struct {
	ID              uuid.UUID `json:"id"`
	ConversationsID uuid.UUID `json:"conversations_id"`
	FileUrl         *string   `json:"file_url"`
	AudioLen        *float64  `json:"audio_len"`
	TaskID          uuid.UUID `json:"task_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Diarize struct {
	ID        uuid.UUID `json:"id"`
	ConvertID uuid.UUID `json:"convert_id"`
	TaskID    uuid.UUID `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Participant struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Prompt struct {
	ID         uuid.UUID `json:"id"`
	PromptName string    `json:"prompt_name"`
	Prompt     string    `json:"prompt"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Report struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Report         *string   `json:"report"`
	PromptID       uuid.UUID `json:"prompt_id"`
	TaskID         uuid.UUID `json:"task_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Segment struct {
	ID        uuid.UUID `json:"id"`
	DiarizeID uuid.UUID `json:"diarize_id"`
	StartTime float64   `json:"start_time"`
	EndTime   float64   `json:"end_time"`
	SpeakerID uuid.UUID `json:"speaker_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SemiReport struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SemiReport     *string   `json:"semi_report"`
	PartNum        int32     `json:"part_num"`
	PromptID       uuid.UUID `json:"prompt_id"`
	TaskID         uuid.UUID `json:"task_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Task struct {
	ID        uuid.UUID         `json:"id"`
	Status    models.TaskStatus `json:"status"`
	TaskType  models.TaskType   `json:"task_type"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Transcription struct {
	ID            uuid.UUID `json:"id"`
	SegmentID     uuid.UUID `json:"segment_id"`
	Transcription *string   `json:"transcription"`
	TaskID        uuid.UUID `json:"task_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
