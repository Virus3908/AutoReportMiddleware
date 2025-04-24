package services

import (
	"context"
	"fmt"
	"log"
	"main/internal/models"
	"main/internal/repositories"
	"main/pkg/messages/proto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	ConvertCallbackPostfix    = "/api/task/update/convert/"
	DiarizeCallbackPostfix    = "/api/task/update/diarize/"
	TranscribeCallbackPostfix = "/api/task/update/transcription/"
	ErrorCallbackPostfix      = "/api/task/update/error/"
)

type TaskDispatcher struct {
	Repo        *repositories.RepositoryStruct
	Messenger   MessageClient
	Storage     StorageClient
	TxManager   TxManager
	CallbackURL string
	TaskFlow    bool
}

func NewTaskDispatcher(
	repo *repositories.RepositoryStruct,
	messenger MessageClient,
	storage StorageClient,
	txManager TxManager,
	taskFlow bool,
	host string,
	port int,
) *TaskDispatcher {

	return &TaskDispatcher{
		Repo:        repo,
		Messenger:   messenger,
		Storage:     storage,
		TxManager:   txManager,
		TaskFlow:    taskFlow,
		CallbackURL: fmt.Sprintf("http://%s:%d", host, port),
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
		convertMessage := &messages.MessageConvertTask{
			TaskId:               taskID.String(),
			FileUrl:              fileURL,
			CallbackUrl:          s.CallbackURL,
			CallbackPostfix:      ConvertCallbackPostfix,
			ErrorCallbackPostfix: ErrorCallbackPostfix,
		}
		return s.Messenger.SendMessage(ctx, models.ConvertTask, conversationID.String(), convertMessage)
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
		diarizeMessage := &messages.MessageDiarizeTask{
			TaskId:               taskID.String(),
			ConvertedFileUrl:     *response.FileUrl,
			CallbackUrl:          s.CallbackURL,
			CallbackPostfix:      DiarizeCallbackPostfix,
			ErrorCallbackPostfix: ErrorCallbackPostfix,
		}

		return s.Messenger.SendMessage(ctx, models.DiarizeTask, conversationID.String(), diarizeMessage)
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
			transcribeMessage := &messages.MessageTranscriptionTask{
				TaskId:               taskID.String(),
				FileUrl:              *segment.FileUrl,
				StartTime:            segment.StartTime,
				EndTime:              segment.EndTime,
				CallbackUrl:          s.CallbackURL,
				CallbackPostfix:      TranscribeCallbackPostfix,
				ErrorCallbackPostfix: ErrorCallbackPostfix,
			}
			err = s.Messenger.SendMessage(ctx, models.TranscribeTask, segment.ConversationID.String(), transcribeMessage)
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
	converted *messages.ConvertTaskResponse,
) error {
	conversationID, err := s.Repo.GetConversationIDByConvertTaskID(ctx, taskID)
	if err != nil {
		return err
	}
	err = s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		err := s.Repo.UpdateConvertByTaskID(ctx, tx, taskID, converted.GetConvertedFileUrl(), converted.GetAudioLen())
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
	segments *messages.SegmentsTaskResponse,
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
		speakerMap := make(map[int32]uuid.UUID)
		for speaker := 0; speaker < int(segments.GetNumOfSpeakers()); speaker++ {
			speakerID, err := s.Repo.CreateSpeakerWithConversationsID(ctx, tx, conversationID, int32(speaker))
			if err != nil {
				return err
			}
			speakerMap[int32(speaker)] = speakerID
		}
		for _, segment := range segments.GetSegments() {
			err := s.Repo.CreateSegment(ctx, tx, diarizeID, segment.GetStartTime(), segment.GetEndTime(), speakerMap[segment.GetSpeaker()])
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
	transcription *messages.TranscriptionTaskResponse,
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
		err = s.Repo.UpdateTranscriptionTextByTaskID(ctx, tx, taskID, transcription.GetTranscription())
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

func (s *TaskDispatcher) HandleErrorCallback(
	ctx context.Context,
	taskID uuid.UUID,
	_ *messages.ErrorTaskResponse,
) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.UpdateTaskStatus(ctx, tx, taskID, models.StatusTaskError)
	})
}

func (s *TaskDispatcher) CreateSemiReportTask(
	ctx context.Context,
	conversationID uuid.UUID,
) error {
	transcriptionWithSpeaker, err := s.Repo.GetFullTranscriptionByConversationID(ctx, conversationID)
	if err != nil {
		return err
	}
	transcriptionWithSpeakerText := ""
	for _, transcription := range transcriptionWithSpeaker {
		if (transcription.Transcription == nil) {
			continue
		}
		if *transcription.Transcription != "" {
			if transcription.ParticipantName != nil {
				transcriptionWithSpeakerText += *transcription.ParticipantName + ": "
			} else {
				transcriptionWithSpeakerText += fmt.Sprintf("Speaker %d: ", transcription.Speaker)
			}
			transcriptionWithSpeakerText += *transcription.Transcription + "\n"
		}
	}
	log.Printf("Transcription with speaker: %s\n", transcriptionWithSpeakerText)
	return nil
}