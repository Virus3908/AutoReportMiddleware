package services

import (
	"context"
	"fmt"
	"strings"

	"main/internal/models"
	"main/internal/repositories"
	"main/pkg/messages/proto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TaskDispatcher struct {
	Repo      *repositories.RepositoryStruct
	Messenger MessageClient
	Storage   StorageClient
	TxManager TxManager
	TaskFlow  bool
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
		Storage:   storage,
		TxManager: txManager,
		TaskFlow:  taskFlow,
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
		convertMessage := &messages.WrapperTask{
			TaskId: taskID.String(),
			Task: &messages.WrapperTask_Convert{
				Convert: &messages.MessageConvertTask{
					FileUrl: fileURL,
				},
			},
		}
		return s.Messenger.SendMessage(ctx, conversationID.String(), convertMessage)
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
		diarizeMessage := &messages.WrapperTask{
			TaskId: taskID.String(),
			Task: &messages.WrapperTask_Diarize{
				Diarize: &messages.MessageDiarizeTask{
					ConvertedFileUrl: *response.FileUrl,
				},
			},
		}

		return s.Messenger.SendMessage(ctx, conversationID.String(), diarizeMessage)
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
			transcribeMessage := &messages.WrapperTask{
				TaskId: taskID.String(),
				Task: &messages.WrapperTask_Transcription{
					Transcription: &messages.MessageTranscriptionTask{
						FileUrl:   *segment.FileUrl,
						StartTime: segment.StartTime,
						EndTime:   segment.EndTime,
					},
				},
			}

			err = s.Messenger.SendMessage(
				ctx,
				segment.ConversationID.String(),
				transcribeMessage)
			if err != nil {
				return fmt.Errorf("failed to send message: %s", err)
			}
		}
		return err
	})
}

func (s *TaskDispatcher) CreateSemiReportTask(
	ctx context.Context,
	conversationID uuid.UUID,
	prompt models.Prompt,
) error {
	promptFromDB, err := s.Repo.GetPromptByName(ctx, nil, prompt.PromptName)
	if err != nil {
		return fmt.Errorf("prompt find error: %s", err)
	}
	transcriptionWithSpeakerText, audioLen, err := s.getTranscriptionWithSpeakerAndAudioLen(ctx, conversationID)
	if err != nil {
		return err
	}
	textParts, err := s.splitTextByAudioLen(*transcriptionWithSpeakerText, audioLen)
	if err != nil {
		return err
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		taskID, err := s.Repo.CreateTask(ctx, tx, models.SemiReportTask)
		if err != nil {
			return err
		}
		for partIndex, part := range textParts {
			err := s.Repo.CreateSemiReport(
				ctx, 
				tx, 
				conversationID, 
				taskID, 
				promptFromDB.ID,
				partIndex,
			)
			if err != nil {
				return err
			}
			semiReportMessage := &messages.WrapperTask{
				TaskId: taskID.String(),
				Task: &messages.WrapperTask_SemiReport{
					SemiReport: &messages.MessageReportTask{
						Message: part,
						Prompt:  promptFromDB.Prompt,
					},
				},
			}
			err = s.Messenger.SendMessage(
				ctx,
				conversationID.String(),
				semiReportMessage,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *TaskDispatcher) getTranscriptionWithSpeakerAndAudioLen(
	ctx context.Context,
	conversationID uuid.UUID,
) (*string, *float64, error) {
	transcriptionWithSpeaker, err := s.Repo.GetFullTranscriptionByConversationID(ctx, conversationID)
	if err != nil {
		return nil, nil, err
	}
	if len(transcriptionWithSpeaker) == 0 {
		return nil, nil, fmt.Errorf("not found transcriptions")
	}
	transcriptionWithSpeakerText := ""
	for _, transcription := range transcriptionWithSpeaker {
		if transcription.Transcription == nil {
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
	return &transcriptionWithSpeakerText, transcriptionWithSpeaker[0].AudioLen, nil
}

func (s *TaskDispatcher) splitTextByAudioLen(
	text string,
	audioLen *float64,
) (map[int]string, error) {
	if audioLen == nil || *audioLen <= 0 {
		return nil, fmt.Errorf("wrong audio len")
	}
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("empty transcription text")
	}

	numChunks := int(*audioLen / 600)
	if numChunks < 1 {
		numChunks = 1
	}
	partLength := float64(len(text)) / float64(numChunks)

	parts := make(map[int]string)
	var current strings.Builder
	partIndex := 1

	for _, line := range strings.Split(text, "\n") {
		current.WriteString(line)
		current.WriteString("\n")

		if float64(current.Len()) >= partLength {
			parts[partIndex] = strings.TrimSpace(current.String())
			current.Reset()
			partIndex++
		}
	}

	if strings.TrimSpace(current.String()) != "" {
		parts[partIndex] = strings.TrimSpace(current.String())
	}

	return parts, nil
}

func (s *TaskDispatcher) HandleTask(
	ctx context.Context,
	task *messages.WrapperResponse,
) error {
	taskID, err := uuid.Parse(task.TaskId)
	if err != nil {
		return err
	}
	switch t := task.Payload.(type) {
	case *messages.WrapperResponse_Convert:
		return s.handleConvertTask(ctx, taskID, t.Convert)
	case *messages.WrapperResponse_Diarize:
		return s.handleDiarizeTask(ctx, taskID, t.Diarize)
	case *messages.WrapperResponse_Transcription:
		return s.handleTransctiprionTask(ctx, taskID, t.Transcription)
	case *messages.WrapperResponse_SemiReport:
		return s.handleSemiReportTask(ctx, taskID, t.SemiReport)
	case *messages.WrapperResponse_Report:
		return nil
		// return s.handleReportTask(ctx, taskID, t.Report)
	case *messages.WrapperResponse_Error:
		return s.handleErrorTask(ctx, taskID, t.Error)
	default:
		return fmt.Errorf("unknown task type")
	}
}

func (s *TaskDispatcher) handleConvertTask(
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

func (s *TaskDispatcher) handleDiarizeTask(
	ctx context.Context,
	taskID uuid.UUID,
	segments *messages.DiarizeTaskResponse,
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

func (s *TaskDispatcher) handleTransctiprionTask(
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
		countOfUntranscribed, err := s.Repo.GetCountOfUntranscribedSegments(ctx, tx, conversationID)
		if err != nil {
			return err
		}
		if countOfUntranscribed == 0 {
			return s.Repo.UpdateConversationStatusByID(ctx, tx, conversationID, models.StatusTranscribed)
		}
		return nil
	})
}

func (s *TaskDispatcher) handleSemiReportTask(
	ctx context.Context,
	taskID uuid.UUID,
	semiReport *messages.ReportTaskResponse,
) error {
	conversationID, err := s.Repo.GetConversationIDBySemiReportTaskID(ctx, taskID)
	if err != nil {
		return err
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		err := s.Repo.UpdateTaskStatus(ctx, tx, taskID, models.StatusTaskOK)
		if err != nil {
			return err
		}
		err = s.Repo.UpdateSemiReportByTaskID(ctx, tx, taskID, semiReport.GetText())
		if err != nil {
			return err
		}
		countOfUnReported, err := s.Repo.GetCountOfUnSemiReportedParts(ctx, tx, conversationID)
		if err != nil {
			return err
		}
		if countOfUnReported == 0 {
			return s.Repo.UpdateConversationStatusByID(ctx, tx, conversationID, models.StatusSemiReported)
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
			// log.Printf("Failed to dispatch next task after type %d: %v\n", taskType, err)
		}
	}()
}

func (s *TaskDispatcher) handleErrorTask(
	ctx context.Context,
	taskID uuid.UUID,
	_ *messages.ErrorTaskResponse,
) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.UpdateTaskStatus(ctx, tx, taskID, models.StatusTaskError)
	})
}
