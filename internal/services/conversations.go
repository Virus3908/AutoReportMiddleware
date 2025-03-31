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

func (s *ConversationsService) GetConversationByID(ctx context.Context, id uuid.UUID) (repositories.Conversation, error) {
	return s.DB.NewQuerry().GetConversationByID(ctx, id)
}

func (s *ConversationsService) GetConversations(ctx context.Context)([]repositories.Conversation, error) {
	return s.DB.NewQuerry().GetConversations(ctx)
}

func (s *ConversationsService) CreateConversation(ctx context.Context, file multipart.File, conversationName, originalFilename string) error {
	defer file.Close()

	fileURL, err := s.Storage.UploadFile(file, originalFilename)
	if err != nil {
		return err
	}

	payload := repositories.CreateConversationParams{
		ConversationName: conversationName,
		FileUrl: fileURL,
	}

	tx, rollback, commit, err := s.DB.StartTransaction()
	if err != nil {
		return err
	}
	defer rollback()
	err = tx.CreateConversation(ctx, payload)
	if err != nil {
		return err
	}
	commit()
	return nil
}

func (s *ConversationsService) UpdateConversation(ctx context.Context, id uuid.UUID, payload repositories.UpdateConversationNameByIDParams) error {
	payload.ID = id
	tx, rollback, commit, err := s.DB.StartTransaction()
	if err != nil {
		return err
	}
	defer rollback()
	err = tx.UpdateConversationNameByID(ctx, payload)
	if err != nil {
		return err
	}
	commit()
	return nil
}

func (s *ConversationsService) DeleteConversation(ctx context.Context, id uuid.UUID) error {
	tx, rollback, commit, err := s.DB.StartTransaction()
	if err != nil {
		return err
	}
	defer rollback()
	fileURL, err := tx.DeleteConversationByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.Storage.DeleteFileByURL(fileURL); err != nil {
		return err
	}
	commit()
	return nil
}