package services

import (
	"context"
	"main/internal/repositories"
	"mime/multipart"

	"github.com/jackc/pgx/v5"
	"main/internal/models"
	"google.golang.org/protobuf/proto"
)



type TxManager interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type StorageClient interface {
	UploadFileAndGetURL(ctx context.Context, file multipart.File, originalFilename string) (string, error)
	DeleteFileByURL(ctx context.Context, fileURL string) error
}

type MessageClient interface {
	SendMessage(ctx context.Context, taskType models.TaskType, key string, message proto.Message) error
}

type ServiceStruct struct {
	Conversations *ConversationsService
	Tasks *TaskDispatcher
}

func New(
	repo *repositories.RepositoryStruct, 
	storage StorageClient, 
	messenger MessageClient,
	txManager TxManager,
	taskFlow bool,
	host string,
	port int,
) *ServiceStruct {
	return &ServiceStruct{
		Conversations: NewConversationsService(repo, storage, txManager),
		Tasks: NewTaskDispatcher(repo, messenger, storage, txManager, taskFlow, host, port),
	}
}



