package models

import (
	"time"
	"github.com/google/uuid"
)

type Conversation struct {
	ID               uuid.UUID `json:"id"`
	ConversationName string    `json:"conversation_name"`
	FileURL          string    `json:"file_url,omitempty"`
	Status           int       `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}