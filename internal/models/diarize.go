package models

import (
	"time"
	"github.com/google/uuid"
)

type Diarize struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        int       `json:"end_time"`
	Speaker        uuid.UUID `json:"speaker"`
	TaskID         uuid.UUID `json:"task_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
