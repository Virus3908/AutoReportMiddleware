package services

import (
	"context"
	"main/internal/repositories"
	"mime/multipart"

	"github.com/google/uuid"
)

type TaskCallbackReceiver struct {
	Storage   StorageClient
	Repo      *repositories.RepositoryStruct
	TxManager TxManager
}

type Segment struct {
	Speaker int32   `json:"speaker"`
	Start   float64 `json:"start"`
	End     float64 `json:"end"`
}

type SegmentsPayload struct {
	Segments []Segment `json:"segments"`
}

const (
	StatusProcessing = 1
	StatusOK         = 2
	StatusError      = 3
)

func NewTaskCallbackReceiver(storage StorageClient, repo *repositories.RepositoryStruct, txManager TxManager) *TaskCallbackReceiver {
	return &TaskCallbackReceiver{
		Storage:   storage,
		Repo:      repo,
		TxManager: txManager,
	}
}

func (s *TaskCallbackReceiver) HandleConvertCallback(
	ctx context.Context,
	taskID uuid.UUID,
	file multipart.File,
	fileName string,
	audioLen float64,
) error {
	fileURL, err := s.Storage.UploadFileAndGetURL(ctx, file, fileName)
	if err != nil {
		return err
	}
	tx, err := s.TxManager.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer s.TxManager.RollbackTransactionIfExist(ctx, tx)
	convertID, err := s.Repo.UpdateConvertByTaskID(ctx, tx, taskID, fileURL, audioLen)
	if err != nil {
		return err
	}
	err = s.Repo.UpdateTaskStatus(ctx, tx, taskID, StatusOK)
	if err != nil {
		return err
	}
	err = s.Repo.UpdateConversationStatusByConvertID(ctx, tx, convertID)
	if err != nil {
		return err
	}
	return s.TxManager.CommitTransaction(ctx, tx)
}

func (s *TaskCallbackReceiver) HandleDiarizeCallback(
	ctx context.Context,
	taskID uuid.UUID,
	payload SegmentsPayload,
) error {
	diarizeID, err := s.Repo.GetDiarizeIDByTaskID(ctx, taskID)
	if err != nil {
		return err
	}
	tx, err := s.TxManager.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer s.TxManager.RollbackTransactionIfExist(ctx, tx)
	err = s.Repo.UpdateTaskStatus(ctx, tx, taskID, StatusOK)
	if err != nil {
		return err
	}
	for _, segment := range payload.Segments {
		err := s.Repo.CreateSegment(ctx, tx, diarizeID, segment.Start, segment.End, segment.Speaker)
		if err != nil {
			return err
		}
	}
	err = s.Repo.UpdateConversationStatusByDiarizeID(ctx, tx, diarizeID)
	if err != nil {
		return err
	}
	return s.TxManager.CommitTransaction(ctx, tx)
}
