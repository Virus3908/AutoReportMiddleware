package services

import (
	"context"
	"mime/multipart"
)

type Repositories interface {
	
}

type StorageClient interface {
	UploadFileAndGetURL(ctx context.Context, file multipart.File, originalFilename string) (string, error)
	DeleteFileByURL(ctx context.Context, fileURL string) error
}

type MessageClient interface {
	SendMessage(ctx context.Context, key string, message string) error
}

type ServiceStruct struct {
	Repo      Repositories
	Storage   StorageClient
	Messenger MessageClient
}

func New(repo Repositories, storage StorageClient, messenger MessageClient) *ServiceStruct {
	return &ServiceStruct{
		Repo:      repo,
		Storage:   storage,
		Messenger: messenger,
	}
}

// func (s *ServiceStruct) GetConversations(ctx context.Context) ([]gen.Conversation, error) {
// 	return s.DB.GetQuery().GetConversations(ctx)
// }
