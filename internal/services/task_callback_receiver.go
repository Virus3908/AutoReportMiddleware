package services

import (
	"context"
	"main/internal/repositories"
	"mime/multipart"

	"github.com/google/uuid"
)

type TaskCallbackReceiver struct {
	Storage StorageClient
	Repo *repositories.RepositoryStruct
	TxManager TxManager
}

const (
	StatusProcessing = 1
	StatusOK = 2
	StatusError = 3
)

func NewTaskCallbackReceiver(storage StorageClient, repo *repositories.RepositoryStruct, txManager TxManager) *TaskCallbackReceiver {
	return &TaskCallbackReceiver{
		Storage: storage,
		Repo: repo,
		TxManager: txManager,
	}
}

func (s *TaskCallbackReceiver) HandleConvertCallback(ctx context.Context, taskID uuid.UUID, file multipart.File, fileName string) error {
	fileURL, err := s.Storage.UploadFileAndGetURL(ctx, file, fileName)
	if err != nil {
		return err
	}
	tx, err := s.TxManager.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer s.TxManager.RollbackTransactionIfExist(ctx, tx)
	err = s.Repo.UpdateConvertFileURL(ctx, tx, taskID, fileURL)
	if err != nil {	
		return err
	}
	err = s.Repo.UpdateTaskStatus(ctx, tx, taskID, StatusProcessing)
	if err != nil {
		return err
	}
	err = s.Repo.UpdateConversationStatus(ctx, tx, taskID, StatusProcessing)
	if err != nil {
		return err
	}
	return s.TxManager.CommitTransaction(ctx, tx)
}