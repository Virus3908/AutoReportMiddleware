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

var conversationStatusPriority = map[models.ConversationStatus]int{
	models.StatusError:       -1,
	models.StatusCreated:     0,
	models.StatusConverted:   1,
	models.StatusDiarized:    2,
	models.StatusTranscribed: 3,
	models.StatusReported:    4,
}

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

	if conversationStatusPriority[conv.Status] >= conversationStatusPriority[models.StatusConverted] {
		file, err := s.Repo.GetConvertFileURLByConversationID(ctx, conversationID)
		if err != nil {
			return nil, err
		}
		if file.FileUrl != nil {
			result.ConvertedFileURL = *file.FileUrl
		}
	}

	if conversationStatusPriority[conv.Status] >= conversationStatusPriority[models.StatusDiarized] {
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
