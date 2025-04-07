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

func (r *RepositoryStruct) UpdateConvertByTaskID(
	ctx context.Context,
	tx pgx.Tx,
	taskID uuid.UUID,
	fileURL string,
	audioLen float64,
) (uuid.UUID, error) {
	query := r.queries.WithTx(tx)
	return query.UpdateConvertByTaskID(ctx, db.UpdateConvertByTaskIDParams{
		TaskID:   taskID,
		FileUrl:  &fileURL,
		AudioLen: &audioLen,
	})
}

func (r *RepositoryStruct) UpdateTaskStatus(ctx context.Context, tx pgx.Tx, taskID uuid.UUID, status int32) error {
	query := r.queries.WithTx(tx)
	return query.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
		ID:     taskID,
		Status: status,
	})
}

func (r *RepositoryStruct) UpdateConversationStatusByConvertID(ctx context.Context, tx pgx.Tx, convertID uuid.UUID) error {
	query := r.queries.WithTx(tx)
	return query.UpdateConversationStatusByConvertID(ctx, db.UpdateConversationStatusByConvertIDParams{
		ID:     convertID,
		Status: StatusConverted,
	})
}

func (r *RepositoryStruct) UpdateConversationStatusByDiarizeID(ctx context.Context, tx pgx.Tx, diarizeID uuid.UUID) error {
	query := r.queries.WithTx(tx)
	return query.UpdateConversationStatusByDiarizeID(ctx, db.UpdateConversationStatusByDiarizeIDParams{
		ID:     diarizeID,
		Status: StatusDiarized,
	})
}

func (r *RepositoryStruct) GetConvertFileURLByConversationID(ctx context.Context, conversationID uuid.UUID) (db.GetConvertFileURLByConversationIDRow, error) {
	return r.queries.GetConvertFileURLByConversationID(ctx, conversationID)
}

func (r *RepositoryStruct) CreateDiarize(ctx context.Context, tx pgx.Tx, convertID, taskID uuid.UUID) error {
	query := r.queries.WithTx(tx)
	return query.CreateDiarize(ctx, db.CreateDiarizeParams{
		TaskID:    taskID,
		ConvertID: convertID,
	})
}

func (r *RepositoryStruct) GetDiarizeIDByTaskID(ctx context.Context, taskID uuid.UUID) (uuid.UUID, error) {
	return r.queries.GetDiarizeIDByTaskID(ctx, taskID)
}

func (r *RepositoryStruct) CreateSegment(
	ctx context.Context,
	tx pgx.Tx,
	diarizeID uuid.UUID,
	startTime, endTime float64,
	speaker int32,
) error {
	query := r.queries.WithTx(tx)
	return query.CreateSegment(ctx, db.CreateSegmentParams{
		DiarizeID: diarizeID,
		StartTime: startTime,
		EndTime:   endTime,
		Speaker:   speaker,
	})
}
