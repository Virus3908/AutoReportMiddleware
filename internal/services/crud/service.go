package crud

import (
	"context"
	"github.com/google/uuid"
	"main/internal/database"
	"main/internal/repositories"
	"main/internal/storage"
)

type CrudService[T any, CreateDTO any, UpdateDTO any] interface {
	GetAll(ctx context.Context) ([]T, error)
	GetByID(ctx context.Context, id uuid.UUID) (T, error)
	Create(ctx context.Context, payload CreateDTO) error
	Update(ctx context.Context, id uuid.UUID, payload UpdateDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CrudConversations interface {
	GetByID(ctx context.Context, id uuid.UUID) (repositories.Conversation, error)
	GetAll(ctx context.Context) ([]repositories.Conversation, error)
	GetFileURLByID(ctx context.Context, id uuid.UUID) (string, error)
	Create(ctx context.Context, payload repositories.CreateConversationParams) error
	Update(ctx context.Context, id uuid.UUID, payload repositories.UpdateConversationNameByIDParams) error
	UpdateStatus(ctx context.Context, id uuid.UUID, payload repositories.UpdateConversationStatusByIDParams) error
	Delete(ctx context.Context, id uuid.UUID) (*string, error)
}

type CrudServicesStruct struct {
	Conversation CrudConversations
	Participant  CrudService[repositories.Participant, repositories.CreateParticipantParams, repositories.UpdateParticipantByIDParams]
	Prompt       CrudService[repositories.Prompt, string, repositories.UpdatePromptByIDParams]
	Convert      *ConvertCRUD
	Task         *TaskCRUD
}

func NewService(db database.Database, storage storage.Storage) *CrudServicesStruct {
	return &CrudServicesStruct{
		Conversation: NewConversationCRUD(db, storage),
		Participant:  NewParticipantCRUD(db),
		Prompt:       NewPromptCRUD(db),
		Convert:      NewConvertCRUD(db),
		Task:         NewTaskCRUD(db),
	}
}
