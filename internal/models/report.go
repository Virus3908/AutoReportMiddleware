package models

import (
	"time"
	"github.com/google/uuid"
)

type Report struct {
	ID            uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Report        string    `json:"report"`
	PromptID      uuid.UUID `json:"prompt_id"`
	TaskID        uuid.UUID `json:"task_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}