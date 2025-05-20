package conversations

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"main/internal/common/interfaces"
	"main/internal/logger"
	"main/internal/models"
	"main/internal/repositories"
	"main/internal/repositories/gen"
)

type ConversationsService struct {
	Repo      *repositories.RepositoryStruct
	Storage   interfaces.StorageClient
	TxManager interfaces.TxManager
}

func NewConversationsService(
	repo *repositories.RepositoryStruct,
	storage interfaces.StorageClient,
	txManager interfaces.TxManager,
) *ConversationsService {
	return &ConversationsService{
		Repo:      repo,
		Storage:   storage,
		TxManager: txManager,
	}
}

func (s *ConversationsService) CreateConversation(ctx context.Context, conversation_name, fileName string, file multipart.File) error {
	fileURL, err := s.Storage.UploadFileAndGetURL(ctx, file, fileName)
	if err != nil {
		return fmt.Errorf("exec: Create Conversation\nfailed to upload file: %w", err)
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.CreateConversation(ctx, tx, fileURL, conversation_name)
	})
}

func (s *ConversationsService) GetConversations(ctx context.Context) ([]db.Conversation, error) {
	return s.Repo.GetConversations(ctx)
}

func (s *ConversationsService) DeleteConversationByID(ctx context.Context, conversationID uuid.UUID) error {

	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		fileURL, err := s.Repo.DeleteConversation(ctx, tx, conversationID)
		if err != nil {
			return fmt.Errorf("exec: Delete Conversation\nfailed to delete conversation: %w", err)
		}
		err = s.Storage.DeleteFileByURL(ctx, fileURL)
		if err != nil {
			logger := logger.GetLoggerFromContext(ctx)
			logger.Warn("exec: Delete Conversation\nfailed to delete file from storage", interfaces.LogField{
				Key:   "file_url",
				Value: fileURL,
			})
		}
		return nil
	})
}

func (s *ConversationsService) GetConversationDetails(ctx context.Context, conversationID uuid.UUID) (*models.ConversationDetail, error) {

	conv, err := s.Repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("exec: Get Conversation Details\nfailed to get conversation by ID: %w", err)
	}

	result := &models.ConversationDetail{
		ConversationID:   conv.ID,
		ConversationName: conv.ConversationName,
		FileURL:          conv.FileUrl,
		Status:           conv.Status,
	}

	if conv.Status >= models.StatusConverted {
		file, err := s.Repo.GetConvertFileURLByConversationID(ctx, conversationID)
		if err != nil {
			logger := logger.GetLoggerFromContext(ctx)
			logger.Error("exec: Get Conversation Details\nfailed to get convert by conversation ID", 
			interfaces.LogField{
				Key:   "execution",
				Value: "GetReportByConversationID",
			},
			interfaces.LogField{
				Key: "Parameter",
				Value: "Convert",
			},
			interfaces.LogField{
				Key:   "error",
				Value: err.Error(),
			})
		}
		if file.FileUrl != nil {
			result.ConvertedFileURL = *file.FileUrl
		}
	}

	if conv.Status >= models.StatusDiarized {
		rows, err := s.Repo.GetSegmentsWithTranscriptionByConversationID(ctx, conversationID)
		if err != nil {
			logger := logger.GetLoggerFromContext(ctx)
			logger.Error("exec: Get Conversation Details\nfailed to get diarization by conversation ID", 
			interfaces.LogField{
				Key:   "execution",
				Value: "GetReportByConversationID",
			},
			interfaces.LogField{
				Key: "Parameter",
				Value: "Diarization",
			},
			interfaces.LogField{
				Key:   "error",
				Value: err.Error(),
			})
		}

		segments := make([]models.SegmentDetail, 0, len(rows))
		for _, row := range rows {
			seg := models.SegmentDetail{
				SegmentID: row.SegmentID,
				StartTime: row.StartTime,
				EndTime:   row.EndTime,
				Speaker:   row.Speaker,
			}
			if row.TranscriptionID != nil {
				seg.TranscriptionID = *row.TranscriptionID
			}
			if row.Transcription != nil {
				seg.Transcription = *row.Transcription
			}
			if row.ParticipantName != nil {
				seg.ParticipantName = *row.ParticipantName
			}
			if row.ParticipantID != nil {
				seg.ParticipantID = *row.ParticipantID
			}
			segments = append(segments, seg)
		}
		result.Segments = segments
	}

	if conv.Status >= models.StatusSemiReported {
		rows, err := s.Repo.GetSemiReportByConversationID(ctx, conversationID)
		if err != nil {
			logger := logger.GetLoggerFromContext(ctx)
			logger.Error("exec: Get Conversation Details\nfailed to get semi_report by conversation ID",
				interfaces.LogField{
					Key:   "execution",
					Value: "GetReportByConversationID",
				},
				interfaces.LogField{
					Key:   "Parameter",
					Value: "SemiReport",
				},
				interfaces.LogField{
					Key:   "error",
					Value: err.Error(),
				})
		}
		if len(rows) > 0 {
			for _, row := range rows {
				if row.SemiReport != nil {
					result.SemiReport += *row.SemiReport + "\n"
				}
			}
		}
	}
	if conv.Status >= models.StatusReported {
		report, err := s.Repo.GetReportByConversationID(ctx, conversationID)
		if err != nil {
			logger := logger.GetLoggerFromContext(ctx)
			logger.Error("exec: Get Conversation Details\nfailed to get report by conversation ID",
				interfaces.LogField{
					Key:   "execution",
					Value: "GetReportByConversationID",
				},
				interfaces.LogField{
					Key:   "Parameter",
					Value: "Report",
				},
				interfaces.LogField{
					Key:   "error",
					Value: err.Error(),
				})
		} else {
			if report.Report != nil {
				result.Report = *report.Report
			}
		}

	}

	return result, nil
}

func (s *ConversationsService) UpdateTranscriptionTextByID(
	ctx context.Context,
	segmentID uuid.UUID,
	transcription models.Transcription,
) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.UpdateTranscriptionTextBySegmentID(ctx, tx, segmentID, transcription.Transcription)
	})
}

func (s *ConversationsService) AssignParticipantToSegment(
	ctx context.Context,
	segmentID uuid.UUID,
	idPair models.ConnectParticipantToConversationType,
) error {
	conversationID, err := uuid.Parse(idPair.ConversationID)
	if err != nil {
		return fmt.Errorf("exec: Assign Participant to Segment\nfailed to parse conversation ID: %w", err)
	}
	var participantID *uuid.UUID
	if idPair.ParticipantID != "" {
		tempID, err := uuid.Parse(idPair.ParticipantID)
		if err != nil {
			return fmt.Errorf("exec: Assign Participant to Segment\nfailed to parse participant ID: %w", err)
		}
		participantID = &tempID
	}
	speakerIDs, err := s.Repo.GetSpeakerParticipantIDBySegmentID(ctx, nil, segmentID)
	if err != nil {
		return fmt.Errorf("exec: Assign Participant to Segment\nfailed to get speaker participant ID by segment ID: %w", err)
	}
	if speakerIDs.ParticipantID != nil {
		countSegmentsWithSpeaker, err := s.Repo.CountSegmentsWithSpeakerID(ctx, nil, speakerIDs.SpeakerID)
		if err != nil {
			return fmt.Errorf("exec: Assign Participant to Segment\nfailed to count segments with speaker ID: %w", err)
		}
		if countSegmentsWithSpeaker > 1 {
			return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
				speakerCount, err := s.Repo.GetSpeakerCountByConversationID(ctx, tx, conversationID)
				if err != nil {
					return fmt.Errorf("exec: Assign Participant to Segment\nfailed to get speaker count by conversation ID: %w", err)
				}
				newSpeakerID, err := s.Repo.CreateNewSpeakerForSegment(ctx, tx, int32(speakerCount), participantID, conversationID)
				if err != nil {
					return fmt.Errorf("exec: Assign Participant to Segment\nfailed to create new speaker for segment: %w", err)
				}
				return s.Repo.AssignNewSpeakerToSegment(ctx, tx, segmentID, newSpeakerID)
			})
		}
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.AssignParticipantToSpeaker(ctx, tx, participantID, speakerIDs.SpeakerID)
	})
}

func (s *ConversationsService) UpdateConversationNameByID(
	ctx context.Context,
	id uuid.UUID,
	conversationName models.ConversationName,
) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.UpdateConversationNameByID(ctx, tx, id, conversationName.ConversationName)
	})
}
