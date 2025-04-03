package crud

import (
	"context"
	"github.com/google/uuid"
	"main/internal/database"
	"main/internal/repositories"
	"main/internal/storage"
)

type ConversationsCRUDStruct struct {
	DB      database.Database
	Storage storage.Storage
}

func NewConversationCRUD(db database.Database, storage storage.Storage) *ConversationsCRUDStruct {
	return &ConversationsCRUDStruct{
		DB:      db,
		Storage: storage,
	}
}

func (s *ConversationsCRUDStruct) GetByID(ctx context.Context, id uuid.UUID) (repositories.Conversation, error) {
	return s.DB.NewQuery().GetConversationByID(ctx, id)
}

func (s *ConversationsCRUDStruct) GetFileURLByID(ctx context.Context, id uuid.UUID) (string, error) {
	return s.DB.NewQuery().GetConversationFileURL(ctx, id)
}

func (s *ConversationsCRUDStruct) GetAll(ctx context.Context) ([]repositories.Conversation, error) {
	return s.DB.NewQuery().GetConversations(ctx)
}

func (s *ConversationsCRUDStruct) Create(ctx context.Context, payload repositories.CreateConversationParams) error {
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.CreateConversation(ctx, payload)
	})
}

func (s *ConversationsCRUDStruct) Update(ctx context.Context, id uuid.UUID, payload repositories.UpdateConversationNameByIDParams) error {
	payload.ID = id
	return s.DB.WithTx(ctx, func(tx *repositories.Queries) error {
		return tx.UpdateConversationNameByID(ctx, payload)
	})
}

func (s *ConversationsCRUDStruct) UpdateStatus(ctx context.Context, id uuid.UUID, payload repositories.UpdateConversationStatusByIDParams) error {
	payload.ID = id
	return s.DB.WithTx(ctx, func(tx *repositories.Queries) error {
		return tx.UpdateConversationStatusByID(ctx, payload)
	})
}

func (s *ConversationsCRUDStruct) Delete(ctx context.Context, id uuid.UUID) (*string, error) {
	var fileURL string

	err := s.DB.WithTx(ctx, func(tx *repositories.Queries) error {
		var err error
		fileURL, err = tx.DeleteConversationByID(ctx, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &fileURL, nil
}
