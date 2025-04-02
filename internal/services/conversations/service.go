package conversations

import (
	"context"
	"fmt"
	"main/internal/repositories"
	"main/internal/services/crud"
	"main/internal/storage"
	"mime/multipart"

	"github.com/google/uuid"
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

func (service *ConversationsService) DeleteConversations(ctx context.Context, id uuid.UUID) error {
	fileURL, err := service.CRUD.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delet error: %s", err)
	}
	return service.Storage.DeleteFileByURL(*fileURL)
}