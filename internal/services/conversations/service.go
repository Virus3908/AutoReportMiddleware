package conversations

import (
	"context"
	"main/internal/repositories"
	"main/internal/services/crud"
	"main/internal/storage"
	"mime/multipart"
)

type ConversationsService struct {
	CRUD    crud.CrudConversations
	Storage storage.Storage
}

func NewService(crud crud.CrudConversations, storage storage.Storage) *ConversationsService {
	return &ConversationsService{CRUD: crud, Storage: storage}
}

func (service *ConversationsService) CreateFromMultipart(
	ctx context.Context,
	file multipart.File,
	filename string,
	conversationName string,
) error {
	defer file.Close()

	fileURL, err := service.Storage.UploadFile(file, filename)
	if err != nil {
		return err
	}

	payload := repositories.CreateConversationParams{
		ConversationName: conversationName,
		FileUrl:          fileURL,
	}

	return service.CRUD.Create(ctx, payload)
}