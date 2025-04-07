package services

import (
	"context"
	"main/internal/repositories"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

type TranscriptionPayload struct {
	Text string `json:"text"`
}

const (
	StatusProcessing = 1
	StatusOK         = 2
	StatusError      = 3
)

const (
	StatusConverted       = 1
	StatusDiarized        = 2
	StatusTranscribed     = 3
	StatusReported        = 4
	StatusSmthngWentWrong = 5
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
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		convertID, err := s.Repo.UpdateConvertByTaskID(ctx, tx, taskID, fileURL, audioLen)
		if err != nil {
			return err
		}
		err = s.Repo.UpdateTaskStatus(ctx, tx, taskID, StatusOK)
		if err != nil {
			return err
		}
		return s.Repo.UpdateConversationStatusByConvertID(ctx, tx, convertID, StatusConverted)
	})
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
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
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
		return s.Repo.UpdateConversationStatusByDiarizeID(ctx, tx, diarizeID, StatusDiarized)
	})
}

func (s *TaskCallbackReceiver) HandleTransctiprionCallback(
	ctx context.Context,
	taskID uuid.UUID,
	payload TranscriptionPayload,
) error {
	conversationID, err := s.Repo.GetConversationIDByTranscriptionTaskID(ctx, taskID)
	if err != nil {
		return err
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		err := s.Repo.UpdateTaskStatus(ctx, tx, taskID, StatusOK)
		if err != nil {
			return err
		}
		err = s.Repo.UpdateTransctiptionTextByID(ctx, tx, taskID, payload.Text)
		if err != nil {
			return err
		}
		count_of_untranscribed, err := s.Repo.GetCountOfUntranscribedSegments(ctx, tx, conversationID)
		if err != nil {
			return err
		}
		if count_of_untranscribed == 0 {
			return s.Repo.UpdateConversationStatusByID(ctx, tx, conversationID, StatusTranscribed)
		}
		return nil
	})
}
