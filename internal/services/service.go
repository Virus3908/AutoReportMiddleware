package services

import (
	"context"
	"main/internal/repositories"
	"mime/multipart"

	"github.com/jackc/pgx/v5"
)



type TxManager interface {
	StartTransaction(ctx context.Context) (pgx.Tx, error)
	StartNestedTransaction(ctx context.Context, tx pgx.Tx) (pgx.Tx, error)
	CommitTransaction(ctx context.Context, tx pgx.Tx) error
	RollbackTransactionIfExist(ctx context.Context, tx pgx.Tx) error
}

type StorageClient interface {
	UploadFileAndGetURL(ctx context.Context, file multipart.File, originalFilename string) (string, error)
	DeleteFileByURL(ctx context.Context, fileURL string) error
}

type MessageClient interface {
	SendMessage(ctx context.Context, taskType int32, key string, msg string) error
}

type ServiceStruct struct {
	Conversations *ConversationsService
	Tasks *TaskDispatcher
	TaskCallbackReceiver *TaskCallbackReceiver
}

func New(repo *repositories.RepositoryStruct, storage StorageClient, messenger MessageClient, txManager TxManager) *ServiceStruct {
	return &ServiceStruct{
		Conversations: NewConversationsService(repo, storage, txManager),
		Tasks: NewTaskDispatcher(repo, messenger, txManager),
		TaskCallbackReceiver: NewTaskCallbackReceiver(storage, repo, txManager),
	}
}



