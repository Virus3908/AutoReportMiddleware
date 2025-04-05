package services

import (
	"context"
	"main/internal/repositories"
	"mime/multipart"
	// "mime/multipart"
)

type ConversationsService struct {
	Repo      *repositories.RepositoryStruct
	Storage   StorageClient
}

func NewConversationsService(repo *repositories.RepositoryStruct, storage StorageClient) *ConversationsService {
	return &ConversationsService{
		Repo:    repo,
		Storage: storage,
	}
}

func (s *ConversationsService) CreateConversation(ctx context.Context, conversation_name, fileName string, file multipart.File) (error) {
	fileURL, err := s.Storage.UploadFileAndGetURL(ctx, file, fileName)
	if err != nil {
		return err
	}
	err = s.Repo.CreateConversation(ctx, fileURL, conversation_name)
	if err != nil {
		return err
	}
	return nil
}

