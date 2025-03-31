package services

import (
	"context"
	"main/internal/database"
	"main/internal/repositories"
	"main/internal/storage"
	"mime/multipart"
	"github.com/google/uuid"
)

type ConversationsService struct {
	DB database.Database
	Storage storage.Storage
}

func NewConversationService(db database.Database, storage storage.Storage) *ConversationsService{
	return &ConversationsService{
		DB: db,
		Storage: storage,
	}
}

func (s *ConversationsService) GetByID(ctx context.Context, id uuid.UUID) (repositories.Conversation, error) {
	return s.DB.NewQuery().GetConversationByID(ctx, id)
}

func (s *ConversationsService) GetAll(ctx context.Context)([]repositories.Conversation, error) {
	return s.DB.NewQuery().GetConversations(ctx)
}

func (s *ConversationsService) Create(ctx context.Context, file multipart.File, conversationName, originalFilename string) error {
	defer file.Close()

	fileURL, err := s.Storage.UploadFile(file, originalFilename)
	if err != nil {
		return err
	}

	payload := repositories.CreateConversationParams{
		ConversationName: conversationName,
		FileUrl: fileURL,
	}

	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.CreateConversation(ctx, payload)
	})
}

func (s *ConversationsService) Update(ctx context.Context, id uuid.UUID, payload repositories.UpdateConversationNameByIDParams) error {
	payload.ID = id
	return s.DB.WithTx(ctx, func(tx *repositories.Queries) error {
		return tx.UpdateConversationNameByID(ctx, payload)
	})
}


func (s *ConversationsService) Delete(ctx context.Context, id uuid.UUID) error {
	var fileURL string

	err :=  s.DB.WithTx(ctx, func(tx *repositories.Queries) error {
		var err error
		fileURL, err = tx.DeleteConversationByID(ctx, id)
		return err
	})
	if err != nil {
		return err
	}

	return s.Storage.DeleteFileByURL(fileURL)
}