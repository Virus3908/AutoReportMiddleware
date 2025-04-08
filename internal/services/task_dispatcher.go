package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/models"
	"main/internal/repositories"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TaskDispatcher struct {
	Repo      *repositories.RepositoryStruct
	Messenger MessageClient
	Storage StorageClient
	TxManager TxManager
	TaskFlow bool
}

func NewTaskDispatcher(
	repo *repositories.RepositoryStruct, 
	messenger MessageClient, 
	storage StorageClient,
	txManager TxManager,
	taskFlow bool,
) *TaskDispatcher {
	return &TaskDispatcher{
		Repo:      repo,
		Messenger: messenger,
		Storage: storage,
		TxManager: txManager,
		TaskFlow: taskFlow,
	}
}

func (s *TaskDispatcher) CreateConvertTask(ctx context.Context, conversationID uuid.UUID) error {
	fileURL, err := s.Repo.GetConversationFileURL(ctx, conversationID)
	if err != nil {
		return err
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		taskID, err := s.Repo.CreateTask(ctx, tx, models.ConvertTask)
		if err != nil {
			return err
		}
		err = s.Repo.CreateConvert(ctx, tx, taskID, conversationID)
		if err != nil {
			return err
		}
		convertMessage := models.Message{
			TaskID:          taskID,
			FileURL:         fileURL,
			CallbackPostfix: "/api/task/update/convert/",
		}
		convertMessageJSON, err := json.Marshal(convertMessage)
		if err != nil {
			return fmt.Errorf("failed to marshal convert message: %w", err)
		}
		return s.Messenger.SendMessage(ctx, models.ConvertTask, conversationID.String(), string(convertMessageJSON))
	})
}

func (s *TaskDispatcher) CreateDiarizeTask(ctx context.Context, conversationID uuid.UUID) error {
	response, err := s.Repo.GetConvertFileURLByConversationID(ctx, conversationID)
	if err != nil {
		return err
	}
	if response.FileUrl == nil {
		return fmt.Errorf("file is not converted")
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		taskID, err := s.Repo.CreateTask(ctx, tx, models.DiarizeTask)
		if err != nil {
			return err
		}
		err = s.Repo.CreateDiarize(ctx, tx, response.ID, taskID)
		if err != nil {
			return err
		}
		diarizeMessage := models.Message{
			TaskID:          taskID,
			FileURL:         *response.FileUrl,
			CallbackPostfix: "/api/task/update/diarize/",
		}
		diarizeMessageJSON, err := json.Marshal(diarizeMessage)
		if err != nil {
			return fmt.Errorf("failed to marshal convert message: %w", err)
		}

		return s.Messenger.SendMessage(ctx, models.DiarizeTask, conversationID.String(), string(diarizeMessageJSON))
	})
}

func (s *TaskDispatcher) CreateTranscribeTask(ctx context.Context, conversationID uuid.UUID) error {
	response, err := s.Repo.GetSegmentsByConversationsID(ctx, conversationID)
	if err != nil {
		return err
	}
	if len(response) == 0 {
		return fmt.Errorf("no segments")
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		for _, segment := range response {
			taskID, err := s.Repo.CreateTask(ctx, tx, models.TranscribeTask)
			if err != nil {
				return err
			}
			err = s.Repo.CreateTranscriptionWithTaskAndSegmentID(ctx, tx, taskID, segment.SegmentID)
			if err != nil {
				return err
			}
			transcribeMessage := models.Message{
				TaskID:          taskID,
				FileURL:         *segment.FileUrl,
				StartTime:       segment.StartTime,
				EndTime:         segment.EndTime,
				CallbackPostfix: "/api/task/update/transcription/",
			}
			transcribeMessageJSON, err := json.Marshal(transcribeMessage)
			if err != nil {
				return fmt.Errorf("failed to marshal convert message: %w", err)
			}
			err = s.Messenger.SendMessage(ctx, models.TranscribeTask, segment.ConversationID.String(), string(transcribeMessageJSON))
			if err != nil {
				return fmt.Errorf("failed to send message: %s", err)
			}
		}
		return err
	})
}

func (s *TaskDispatcher) HandleConvertCallback(
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
	conversationID, err := s.Repo.GetConversationIDByConvertTaskID(ctx, taskID)
	if err != nil {
		return err
	}
	err = s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		err := s.Repo.UpdateConvertByTaskID(ctx, tx, taskID, fileURL, audioLen)
		if err != nil {
			return err
		}
		err = s.Repo.UpdateTaskStatus(ctx, tx, taskID, models.StatusTaskOK)
		if err != nil {
			return err
		}
		return s.Repo.UpdateConversationStatusByID(ctx, tx, conversationID, models.StatusConverted)
	})
	if err != nil {
		return err
	}
	s.dispatchNext(models.ConvertTask, conversationID)
	return nil
}

func (s *TaskDispatcher) HandleDiarizeCallback(
	ctx context.Context,
	taskID uuid.UUID,
	payload models.SegmentsPayload,
) error {
	diarizeID, err := s.Repo.GetDiarizeIDByTaskID(ctx, taskID)
	if err != nil {
		return err
	}
	conversationID, err := s.Repo.GetConversationIDByDiarizeTaskID(ctx, taskID)
	if err != nil {
		return err
	}
	err = s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		err = s.Repo.UpdateTaskStatus(ctx, tx, taskID, models.StatusTaskOK)
		if err != nil {
			return err
		}
		for _, segment := range payload.Segments {
			err := s.Repo.CreateSegment(ctx, tx, diarizeID, segment.Start, segment.End, segment.Speaker)
			if err != nil {
				return err
			}
		}
		return s.Repo.UpdateConversationStatusByID(ctx, tx, conversationID, models.StatusDiarized)
	})
	if err != nil {
		return err
	}
	s.dispatchNext(models.DiarizeTask, conversationID)
	return nil
}

func (s *TaskDispatcher) HandleTransctiprionCallback(
	ctx context.Context,
	taskID uuid.UUID,
	payload models.TranscriptionPayload,
) error {
	conversationID, err := s.Repo.GetConversationIDByTranscriptionTaskID(ctx, taskID)
	if err != nil {
		return err
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		err := s.Repo.UpdateTaskStatus(ctx, tx, taskID, models.StatusTaskOK)
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
			return s.Repo.UpdateConversationStatusByID(ctx, tx, conversationID, models.StatusTranscribed)
		}
		return nil
	})
}

func (s *TaskDispatcher) dispatchNext(taskType int, conversationID uuid.UUID) {
	if !s.TaskFlow {
		return
	}

	go func() {
		var err error
		switch taskType {
		case models.ConvertTask:
			err = s.CreateDiarizeTask(context.Background(), conversationID)
		case models.DiarizeTask:
			err = s.CreateTranscribeTask(context.Background(), conversationID)
		}

		if err != nil {
			log.Printf("Failed to dispatch next task after type %d: %v\n", taskType, err)
		}
	}()
}