package services

import (
	"context"
	"main/internal/repositories"
	"mime/multipart"
)

// type Repositories interface {

// }

type StorageClient interface {
	UploadFileAndGetURL(ctx context.Context, file multipart.File, originalFilename string) (string, error)
	DeleteFileByURL(ctx context.Context, fileURL string) error
}

type MessageClient interface {
	SendMessage(ctx context.Context, key string, message string) error
}

type ServiceStruct struct {
	Messenger MessageClient
	Conversations *ConversationsService
	Tasks *TaskDispatcher
}

func New(repo *repositories.RepositoryStruct, storage StorageClient, messenger MessageClient) *ServiceStruct {
	return &ServiceStruct{
		Conversations: NewConversationsService(repo, storage),
		Tasks: NewTaskDispatcher(repo, messenger),
	}
}



