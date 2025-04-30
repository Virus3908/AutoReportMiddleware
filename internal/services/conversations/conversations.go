package conversations

import (
	"context"
	"main/internal/common/interfaces"
	"main/internal/models"
	"main/internal/repositories"
	"main/internal/repositories/gen"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		return err
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
			return err
		}
		return s.Storage.DeleteFileByURL(ctx, fileURL)
	})
}

func (s *ConversationsService) GetConversationDetails(ctx context.Context, conversationID uuid.UUID) (*models.ConversationDetail, error) {
	conv, err := s.Repo.GetConversationDetails(ctx, conversationID)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		if file.FileUrl != nil {
			result.ConvertedFileURL = *file.FileUrl
		}
	}

	if conv.Status >= models.StatusDiarized {
		rows, err := s.Repo.GetSegmentsWithTranscriptionByConversationID(ctx, conversationID)
		if err != nil {
			return nil, err
		}

		segments := make([]models.SegmentDetail, 0, len(rows))
		for _, row := range rows {
			seg := models.SegmentDetail{
				SegmentID: row.SegmentID,
				StartTime: row.StartTime,
				EndTime:   row.EndTime,
				Speaker:   row.Speaker,
				// Transcription:   "",
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

	return result, nil
}

func (s *ConversationsService) UpdateTranscriptionTextByID(
	ctx context.Context,
	id uuid.UUID,
	transcription models.Transcription,
) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.UpdateTransctiptionTextByID(ctx, tx, id, transcription.Transcription)
	})
}

func (s *ConversationsService) AssignParticipantToSegment(
	ctx context.Context,
	segmentID uuid.UUID,
	idPair models.ConnectParticipantToConversationType,
) error {
	conversationID, err := uuid.Parse(idPair.ConversationID)
	if err != nil {
		return err
	}
	var participantID *uuid.UUID
	if idPair.ParticipantID != "" {
		tempID, err := uuid.Parse(idPair.ParticipantID)
		if err != nil {
			return err
		}
		participantID = &tempID
	} else {
		participantID = nil
	}
	speakerIDs, err := s.Repo.GetSpeakerParticipantIDBySegmentID(ctx, nil, segmentID)
	if err != nil {
		return err
	}
	if speakerIDs.ParticipantID != nil {
		countSegmentsWithSpeaker, err := s.Repo.CountSegmentsWithSpeakerID(ctx, nil, speakerIDs.SpeakerID)
		if err != nil {
			return err
		}
		if countSegmentsWithSpeaker > 1 {
			return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
				speakerCount, err := s.Repo.GetSpeakerCountByConversationID(ctx, tx, conversationID)
				if err != nil {
					return err
				}
				newSpeakerID, err := s.Repo.CreateNewSpeakerForSegment(ctx, tx, int32(speakerCount), participantID, conversationID)
				if err != nil {
					return err
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
