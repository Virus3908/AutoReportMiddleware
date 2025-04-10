package repositories

import (
	"context"
	"main/internal/models"
	"main/internal/repositories/gen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *RepositoryStruct) CreateTask(ctx context.Context, tx pgx.Tx, taskType models.TaskType) (uuid.UUID, error) {
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
) error {
	query := r.queries.WithTx(tx)
	return query.UpdateConvertByTaskID(ctx, db.UpdateConvertByTaskIDParams{
		TaskID:   taskID,
		FileUrl:  &fileURL,
		AudioLen: &audioLen,
	})
}

func (r *RepositoryStruct) UpdateTaskStatus(ctx context.Context, tx pgx.Tx, taskID uuid.UUID, status models.TaskStatus) error {
	query := r.queries.WithTx(tx)
	return query.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
		ID:     taskID,
		Status: status,
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
	speakerID uuid.UUID,
) error {
	query := r.queries.WithTx(tx)
	return query.CreateSegment(ctx, db.CreateSegmentParams{
		DiarizeID: diarizeID,
		StartTime: startTime,
		EndTime:   endTime,
		SpeakerID: speakerID,
	})
}

func (r *RepositoryStruct) GetSegmentsByConversationsID(
	ctx context.Context,
	conversationID uuid.UUID,
) ([]db.GetSegmentsByConversationsIDRow, error) {
	return r.queries.GetSegmentsByConversationsID(ctx, conversationID)
}

func (r *RepositoryStruct) CreateTranscriptionWithTaskAndSegmentID(
	ctx context.Context,
	tx pgx.Tx,
	taskID uuid.UUID,
	segmentID uuid.UUID,
) error {
	query := r.queries.WithTx(tx)
	return query.CreateTranscriptionWithTaskAndSegmentID(
		ctx, db.CreateTranscriptionWithTaskAndSegmentIDParams{
			TaskID:    taskID,
			SegmentID: segmentID,
		},
	)
}

func (r *RepositoryStruct) UpdateTransctiptionTextByID(
	ctx context.Context,
	tx pgx.Tx,
	taskID uuid.UUID,
	text string,
) error {
	query := r.queries.WithTx(tx)
	return query.UpdateTransctiptionTextByID(ctx, db.UpdateTransctiptionTextByIDParams{
		TaskID:        taskID,
		Transcription: &text,
	})
}

func (r *RepositoryStruct) GetCountOfUntranscribedSegments(
	ctx context.Context,
	tx pgx.Tx,
	conversationID uuid.UUID,
) (int64, error) {
	query := r.queries
	if tx != nil {
		query = r.queries.WithTx(tx)
	}
	return query.GetCountOfUntranscribedSegments(ctx, conversationID)
}

func (r *RepositoryStruct) GetConversationIDByTranscriptionTaskID(
	ctx context.Context,
	taskID uuid.UUID,
) (uuid.UUID, error) {
	return r.queries.GetConversationIDByTranscriptionTaskID(ctx, taskID)
}

func (r *RepositoryStruct) UpdateConversationStatusByID(
	ctx context.Context,
	tx pgx.Tx,
	conversationID uuid.UUID,
	status models.ConversationStatus,
) error {
	query := r.queries.WithTx(tx)
	return query.UpdateConversationStatusByID(ctx, db.UpdateConversationStatusByIDParams{
		ID:     conversationID,
		Status: status,
	})
}

func (r *RepositoryStruct) GetConversationIDByConvertTaskID(
	ctx context.Context,
	taskID uuid.UUID,
) (uuid.UUID, error) {
	return r.queries.GetConversationIDByConvertTaskID(ctx, taskID)
}

func (r *RepositoryStruct) GetConversationIDByDiarizeTaskID(
	ctx context.Context,
	taskID uuid.UUID,
) (uuid.UUID, error) {
	return r.queries.GetConversationIDByDiarizeTaskID(ctx, taskID)
}

func (r *RepositoryStruct) GetSegmentsWithTranscriptionByConversationID(
	ctx context.Context,
	conversationID uuid.UUID,
) ([]db.GetSegmentsWithTranscriptionByConversationIDRow, error) {
	return r.queries.GetSegmentsWithTranscriptionByConversationID(ctx, conversationID)
}

func (r *RepositoryStruct) CreateSpeakerWithConversationsID(
	ctx context.Context,
	tx pgx.Tx,
	conversationID uuid.UUID,
	speaker int32,
) (uuid.UUID, error) {
	query := r.queries.WithTx(tx)
	return query.CreateSpeakerWithConversationsID(ctx,
		db.CreateSpeakerWithConversationsIDParams{
			ConversationID: conversationID,
			Speaker:        speaker,
		})
}
