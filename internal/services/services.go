package services

import (
	"context"
	"main/internal/database"
	"main/internal/repositories"
	"main/internal/storage"

	"github.com/google/uuid"
)

type CrudService[T any, CreateDTO any, UpdateDTO any] interface {
	GetAll(ctx context.Context) ([]T, error)
	GetByID(ctx context.Context, id uuid.UUID) (T, error)
	Create(ctx context.Context, payload CreateDTO) error
	Update(ctx context.Context, id uuid.UUID, payload UpdateDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ServicesStruct struct {
	Conversation *ConversationsService
	Participant  CrudService[repositories.Participant, repositories.CreateParticipantParams, repositories.UpdateParticipantByIDParams]
	Prompt       CrudService[repositories.Prompt, string, repositories.UpdatePromptByIDParams]
	Convert      *ConvertService
}

func NewService(db database.Database, storage storage.Storage) *ServicesStruct {
	return &ServicesStruct{
		Conversation: NewConversationService(db, storage),
		Participant:  NewParticipantService(db),
		Prompt:       NewPromptService(db),
		Convert:      NewConvertService(db),
	}
}
