package taskDispatcher

import (
	"context"
	"fmt"
	"strings"

	"main/internal/common/interfaces"
	"main/internal/logger"
	"main/internal/models"
	"main/internal/repositories"
	"main/pkg/messages/proto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TaskDispatcher struct {
	Repo      *repositories.RepositoryStruct
	Messenger interfaces.MessageClient
	Storage   interfaces.StorageClient
	TxManager interfaces.TxManager
	TaskFlow  bool
}

func NewTaskDispatcher(
	repo *repositories.RepositoryStruct,
	messenger interfaces.MessageClient,
	storage interfaces.StorageClient,
	txManager interfaces.TxManager,
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
		return fmt.Errorf("exec: Create Convert Task\nfailed to get conversation file URL: %w", err)
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		taskID, err := s.Repo.CreateTask(ctx, tx, models.ConvertTask)
		if err != nil {
			return fmt.Errorf("exec: Create Convert Task\nfailed to create task: %w", err)
		}
		err = s.Repo.CreateConvert(ctx, tx, taskID, conversationID)
		if err != nil {
			return fmt.Errorf("exec: Create Convert Task\nfailed to create convert: %w", err)
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
		return fmt.Errorf("exec: Create Diarize Task\nfailed to get convert file URL: %w", err)
	}
	if response.FileUrl == nil {
		return fmt.Errorf("file is not converted")
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		taskID, err := s.Repo.CreateTask(ctx, tx, models.DiarizeTask)
		if err != nil {
			return fmt.Errorf("exec: Create Diarize Task\nfailed to create task: %w", err)
		}
		err = s.Repo.CreateDiarize(ctx, tx, response.ID, taskID)
		if err != nil {
			return fmt.Errorf("exec: Create Diarize Task\nfailed to create diarize: %w", err)
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
		return fmt.Errorf("exec: Create Transcribe Task\nfailed to get segments: %w", err)
	}
	if len(response) == 0 {
		return fmt.Errorf("exec: Create Transcribe Task\nno segments found")
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		for _, segment := range response {
			taskID, err := s.Repo.CreateTask(ctx, tx, models.TranscribeTask)
			if err != nil {
				return fmt.Errorf("exec: Create Transcribe Task\nfailed to create task: %w", err)
			}
			err = s.Repo.CreateTranscriptionWithTaskAndSegmentID(ctx, tx, taskID, segment.SegmentID)
			if err != nil {
				return fmt.Errorf("exec: Create Transcribe Task\nfailed to create transcription: %w", err)
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
				return fmt.Errorf("exec: Create Transcribe Task\nfailed to send message: %s", err)
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
		return fmt.Errorf("exec: Create Semi Report Task\nfailed to get prompt by name: %w", err)
	}
	transcriptionWithSpeakerText, audioLen, err := s.getTranscriptionWithSpeakerAndAudioLen(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("exec: Create Semi Report Task\nfailed to get transcription with speaker and audio length: %w", err)
	}
	textParts, err := s.splitTextByAudioLen(*transcriptionWithSpeakerText, audioLen)
	if err != nil {
		return fmt.Errorf("exec: Create Semi Report Task\nfailed to split text by audio length: %w", err)
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		taskID, err := s.Repo.CreateTask(ctx, tx, models.SemiReportTask)
		if err != nil {
			return fmt.Errorf("exec: Create Semi Report Task\nfailed to create task: %w", err)
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
				return fmt.Errorf("exec: Create Semi Report Task\nfailed to create semi report: %w", err)
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
				return fmt.Errorf("exec: Create Semi Report Task\nfailed to send message: %w", err)
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
		return nil, nil, fmt.Errorf("exec: Get Transcription With Speaker\nfailed to get transcription with speaker: %w", err)
	}
	if len(transcriptionWithSpeaker) == 0 {
		return nil, nil, fmt.Errorf("exec: Get Transcription With Speaker\nno transcription found")
	}
	transcriptionWithSpeakerText := ""
	for _, transcription := range transcriptionWithSpeaker {
		if transcription.Transcription == nil || *transcription.Transcription == "" {
			continue
		}
		if transcription.ParticipantName != nil {
			transcriptionWithSpeakerText += *transcription.ParticipantName + ": "
		} else {
			transcriptionWithSpeakerText += fmt.Sprintf("Speaker %d: ", transcription.Speaker)
		}
		transcriptionWithSpeakerText += *transcription.Transcription + "\n"
	}
	return &transcriptionWithSpeakerText, transcriptionWithSpeaker[0].AudioLen, nil
}

func (s *TaskDispatcher) splitTextByAudioLen(
	text string,
	audioLen *float64,
) (map[int]string, error) {
	if audioLen == nil || *audioLen <= 0 {
		return nil, fmt.Errorf("exec: Split Text By Audio Length\ninvalid audio length")
	}
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("exec: Split Text By Audio Length\nempty text")
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
		return fmt.Errorf("exec: Handle Task\nfailed to parse task ID: %w", err)
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
		return fmt.Errorf("exec: Handle Task\nunknown task type: %T", t)
	}
}

func (s *TaskDispatcher) handleConvertTask(
	ctx context.Context,
	taskID uuid.UUID,
	converted *messages.ConvertTaskResponse,
) error {
	conversationID, err := s.Repo.GetConversationIDByConvertTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("exec: Handle Convert Task\nfailed to get conversation ID by convert task ID: %w", err)
	}
	err = s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		err := s.Repo.UpdateConvertByTaskID(ctx, tx, taskID, converted.GetConvertedFileUrl(), converted.GetAudioLen())
		if err != nil {
			return fmt.Errorf("exec: Handle Convert Task\nfailed to update convert by task ID: %w", err)
		}
		err = s.Repo.UpdateTaskStatus(ctx, tx, taskID, models.StatusTaskOK)
		if err != nil {
			return fmt.Errorf("exec: Handle Convert Task\nfailed to update task status: %w", err)
		}
		return s.Repo.UpdateConversationStatusByID(ctx, tx, conversationID, models.StatusConverted)
	})
	if err != nil {
		return fmt.Errorf("exec: Handle Convert Task\nfailed to update convert: %w", err)
	}
	s.dispatchNext(ctx, models.ConvertTask, conversationID)
	return nil
}

func (s *TaskDispatcher) handleDiarizeTask(
	ctx context.Context,
	taskID uuid.UUID,
	segments *messages.DiarizeTaskResponse,
) error {
	diarizeID, err := s.Repo.GetDiarizeIDByTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("exec: Handle Diarize Task\nfailed to get diarize ID by task ID: %w", err)
	}
	conversationID, err := s.Repo.GetConversationIDByDiarizeTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("exec: Handle Diarize Task\nfailed to get conversation ID by diarize task ID: %w", err)
	}
	err = s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		err = s.Repo.UpdateTaskStatus(ctx, tx, taskID, models.StatusTaskOK)
		if err != nil {
			return fmt.Errorf("exec: Handle Diarize Task\nfailed to update task status: %w", err)
		}
		speakerMap := make(map[int32]uuid.UUID)
		for speaker := 0; speaker < int(segments.GetNumOfSpeakers()); speaker++ {
			speakerID, err := s.Repo.CreateSpeakerWithConversationsID(ctx, tx, conversationID, int32(speaker))
			if err != nil {
				return fmt.Errorf("exec: Handle Diarize Task\nfailed to create speaker: %w", err)
			}
			speakerMap[int32(speaker)] = speakerID
		}
		for _, segment := range segments.GetSegments() {
			err := s.Repo.CreateSegment(ctx, tx, diarizeID, segment.GetStartTime(), segment.GetEndTime(), speakerMap[segment.GetSpeaker()])
			if err != nil {
				return fmt.Errorf("exec: Handle Diarize Task\nfailed to create segment: %w", err)
			}
		}
		return s.Repo.UpdateConversationStatusByID(ctx, tx, conversationID, models.StatusDiarized)
	})
	if err != nil {
		return fmt.Errorf("exec: Handle Diarize Task\nfailed to update diarize: %w", err)
	}
	s.dispatchNext(ctx, models.DiarizeTask, conversationID)
	return nil
}

func (s *TaskDispatcher) handleTransctiprionTask(
	ctx context.Context,
	taskID uuid.UUID,
	transcription *messages.TranscriptionTaskResponse,
) error {
	conversationID, err := s.Repo.GetConversationIDByTranscriptionTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("exec: Handle Transcription Task\nfailed to get conversation ID by transcription task ID: %w", err)
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		err := s.Repo.UpdateTaskStatus(ctx, tx, taskID, models.StatusTaskOK)
		if err != nil {
			return fmt.Errorf("exec: Handle Transcription Task\nfailed to update task status: %w", err)
		}
		err = s.Repo.UpdateTranscriptionTextByTaskID(ctx, tx, taskID, transcription.GetTranscription())
		if err != nil {
			return fmt.Errorf("exec: Handle Transcription Task\nfailed to update transcription text by task ID: %w", err)
		}
		countOfUntranscribed, err := s.Repo.GetCountOfUntranscribedSegments(ctx, tx, conversationID)
		if err != nil {
			return fmt.Errorf("exec: Handle Transcription Task\nfailed to get count of untranscribed segments: %w", err)
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
		return fmt.Errorf("exec: Handle Semi Report Task\nfailed to get conversation ID by semi report task ID: %w", err)
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		err := s.Repo.UpdateTaskStatus(ctx, tx, taskID, models.StatusTaskOK)
		if err != nil {
			return fmt.Errorf("exec: Handle Semi Report Task\nfailed to update task status: %w", err)
		}
		err = s.Repo.UpdateSemiReportByTaskID(ctx, tx, taskID, semiReport.GetText())
		if err != nil {
			return fmt.Errorf("exec: Handle Semi Report Task\nfailed to update semi report by task ID: %w", err)
		}
		countOfUnReported, err := s.Repo.GetCountOfUnSemiReportedParts(ctx, tx, conversationID)
		if err != nil {
			return fmt.Errorf("exec: Handle Semi Report Task\nfailed to get count of unreported parts: %w", err)
		}
		if countOfUnReported == 0 {
			return s.Repo.UpdateConversationStatusByID(ctx, tx, conversationID, models.StatusSemiReported)
		}
		return nil
	})
}

func (s *TaskDispatcher) dispatchNext(ctx context.Context, taskType int, conversationID uuid.UUID) {
	if !s.TaskFlow {
		return
	}
	logger := logger.GetLoggerFromContext(ctx)

	var err error
	switch taskType {
	case models.ConvertTask:
		logger.Info("dispatching next task: Diarize", interfaces.LogField{Key: "conversation_id", Value: conversationID})
		err = s.CreateDiarizeTask(ctx, conversationID)
	case models.DiarizeTask:
		logger.Info("dispatching next task: Transcribe", interfaces.LogField{Key: "conversation_id", Value: conversationID})
		err = s.CreateTranscribeTask(ctx, conversationID)
	}

	if err != nil {
		logger.Error("failed to dispatch next task", interfaces.LogField{Key: "error", Value: err})
	}

}

func (s *TaskDispatcher) handleErrorTask(
	ctx context.Context,
	taskID uuid.UUID,
	msg *messages.ErrorTaskResponse,
) error {
	logger := logger.GetLoggerFromContext(ctx)
	logger.Error("task failed", interfaces.LogField{Key: "task_id", Value: taskID}, interfaces.LogField{Key: "error", Value: msg.GetError()})
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.UpdateTaskStatus(ctx, tx, taskID, models.StatusTaskError)
	})
}
