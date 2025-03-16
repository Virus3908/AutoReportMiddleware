package models

import (
	"time"
	"github.com/google/uuid"
)

type Convert struct {
	ID             uuid.UUID `json:"id"`
	ConversationsID uuid.UUID `json:"conversations_id"`
	FileURL        string    `json:"file_url,omitempty"`
	TaskID         uuid.UUID `json:"task_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
