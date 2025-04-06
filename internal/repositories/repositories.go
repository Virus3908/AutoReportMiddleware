package repositories

import (
	"context"
	"main/internal/repositories/gen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	StatusConverted   = 1
	StatusDiarized    = 2
	StatusTranscribed = 3
	StatusReported    = 4
	StatusError       = 5
)

type RepositoryStruct struct {
	queries *db.Queries
}

func New(pool *pgxpool.Pool) *RepositoryStruct {
	return &RepositoryStruct{
		queries: db.New(pool),
	}
}

func (r *RepositoryStruct) GetConversations(ctx context.Context) ([]db.Conversation, error) {
	return r.queries.GetConversations(ctx)
}

func (r *RepositoryStruct) CreateConversation(ctx context.Context, tx pgx.Tx, fileURL, conversation_name string) error {
	query := r.queries.WithTx(tx)
	return query.CreateConversation(ctx, db.CreateConversationParams{
		FileUrl:          fileURL,
		ConversationName: conversation_name,
	})
}

func (r *RepositoryStruct) GetConversationDetails(ctx context.Context, conversationID uuid.UUID) (*db.Conversation, error) {
	conversation, err := r.queries.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *RepositoryStruct) GetConversationFileURL(ctx context.Context, conversationID uuid.UUID) (string, error) {
	return r.queries.GetConversationFileURL(ctx, conversationID)
}

func (r *RepositoryStruct) CreateTask(ctx context.Context, tx pgx.Tx, taskType int32) (uuid.UUID, error) {
	query := r.queries.WithTx(tx)
	return query.CreateTask(ctx, taskType)
}

func (r *RepositoryStruct) CreateConvert(ctx context.Context, tx pgx.Tx, taskID, conversationID uuid.UUID) error {
	query := r.queries.WithTx(tx)
	return query.CreateConvert(ctx, db.CreateConvertParams{
		ConversationsID: conversationID,
		TaskID:          taskID,
	})
}

func (r *RepositoryStruct) DeleteConversation(ctx context.Context, tx pgx.Tx, conversationID uuid.UUID) (string, error) {
	query := r.queries.WithTx(tx)
	return query.DeleteConversationByID(ctx, conversationID)
}

func (r *RepositoryStruct) UpdateConvertFileURL(ctx context.Context, tx pgx.Tx, taskID uuid.UUID, fileURL string) error {
	query := r.queries.WithTx(tx)
	return query.UpdateConvertByTaskID(ctx, db.UpdateConvertByTaskIDParams{
		TaskID:  taskID,
		FileUrl: &fileURL,
	})
}

func (r *RepositoryStruct) UpdateTaskStatus(ctx context.Context, tx pgx.Tx, taskID uuid.UUID, status int32) error {
	query := r.queries.WithTx(tx)
	return query.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
		ID:     taskID,
		Status: status,
	})
}

func (r *RepositoryStruct) UpdateConversationStatus(ctx context.Context, tx pgx.Tx, taskID uuid.UUID, status int32) error {
	query := r.queries.WithTx(tx)
	switch status {
	case StatusConverted:
		return query.UpdateConversationStatusByConvertID(ctx, db.UpdateConversationStatusByConvertIDParams{
			ID:     taskID,
			Status: status,
		})
	}
	return query.UpdateConversationStatusByID(ctx, db.UpdateConversationStatusByIDParams{
		ID:     taskID,
		Status: StatusError})
}
