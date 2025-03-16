package models

import (
	"time"
	"github.com/google/uuid"
)

type Prompt struct {
	ID        uuid.UUID `json:"id"`
	Prompt    string    `json:"prompt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}