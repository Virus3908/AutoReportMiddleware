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
	Repo      *repositories.RepositoryStruct
	Storage   StorageClient
	Messenger MessageClient
	Conversations *ConversationsService
}

func New(repo *repositories.RepositoryStruct, storage StorageClient, messenger MessageClient) *ServiceStruct {
	return &ServiceStruct{
		Messenger: messenger,
		Conversations: NewConversationsService(repo, storage),
	}
}



