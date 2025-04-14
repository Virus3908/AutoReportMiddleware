package services

import (
	"context"
	"main/internal/models"
	"main/internal/repositories"
	"main/internal/repositories/gen"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ConversationsService struct {
	Repo      *repositories.RepositoryStruct
	Storage   StorageClient
	TxManager TxManager
}

func NewConversationsService(repo *repositories.RepositoryStruct, storage StorageClient, txManager TxManager) *ConversationsService {
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

func (s *ConversationsService) DeleteConversation(ctx context.Context, conversationID uuid.UUID) error {
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
				SegmentID:     row.SegmentID,
				StartTime:     row.StartTime,
				EndTime:       row.EndTime,
				Speaker:       row.Speaker,
				Transcription: "",
			}
			if row.TranscriptionID != nil {
				seg.TranscriptionID = *row.TranscriptionID
			}
			if row.Transcription != nil {
				seg.Transcription = *row.Transcription
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

func (s *ConversationsService) CreateParticipant(
	ctx context.Context,
	participantPayload models.ParticipantData,
) (error) {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.CreateParticipant(ctx, tx, participantPayload)
	})
}

func (s *ConversationsService) GetParticipants(
	ctx context.Context,
) ([]models.ParticipantData, error) {
	participants, err := s.Repo.GetParticipants(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.ParticipantData, 0, len(participants))
	for _, participant := range participants {
		result = append(result, models.ParticipantData{
			Name:  *participant.Name,
			Email: participant.Email,
		})
	}
	return result, nil
}