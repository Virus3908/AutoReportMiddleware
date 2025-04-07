package services

import (
	"context"
	"main/internal/repositories"
	"main/internal/repositories/gen"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	// "mime/multipart"
)

type ConversationsService struct {
	Repo      *repositories.RepositoryStruct
	Storage   StorageClient
	TxManager TxManager
}

func NewConversationsService(repo *repositories.RepositoryStruct, storage StorageClient, txManager TxManager) *ConversationsService {
	return &ConversationsService{
		Repo:      repo,
		Storage:   storage,
		TxManager: txManager,
	}
}

func (s *ConversationsService) CreateConversation(ctx context.Context, conversation_name, fileName string, file multipart.File) error {
	fileURL, err := s.Storage.UploadFileAndGetURL(ctx, file, fileName)
	if err != nil {
		return err
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.CreateConversation(ctx, tx, fileURL, conversation_name)
	})
}

func (s *ConversationsService) GetConversations(ctx context.Context) ([]db.Conversation, error) {
	return s.Repo.GetConversations(ctx)
}

func (s *ConversationsService) GetConversationDetails(ctx context.Context, conversationID uuid.UUID) (*db.Conversation, error) {
	return s.Repo.GetConversationDetails(ctx, conversationID)
}

func (s *ConversationsService) DeleteConversation(ctx context.Context, conversationID uuid.UUID) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		fileURL, err := s.Repo.DeleteConversation(ctx, tx, conversationID)
		if err != nil {
			return err
		}
		return s.Storage.DeleteFileByURL(ctx, fileURL)
	})
}
