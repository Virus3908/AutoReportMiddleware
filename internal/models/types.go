package models

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID               uuid.UUID `json:"id"`
	ConversationName string    `json:"conversation_name"`
	FileUrl          string    `json:"file_url"`
	Status           int32     `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
