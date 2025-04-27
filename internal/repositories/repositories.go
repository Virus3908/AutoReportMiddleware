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

func (r *RepositoryStruct) UpdateTranscriptionTextByTaskID(
	ctx context.Context,
	tx pgx.Tx,
	taskID uuid.UUID,
	text string,
) error {
	query := r.queries.WithTx(tx)
	return query.UpdateTranscriptionTextByTaskID(ctx, db.UpdateTranscriptionTextByTaskIDParams{
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

func (r *RepositoryStruct) UpdateTransctiptionTextByID(
	ctx context.Context,
	tx pgx.Tx,
	taskID uuid.UUID,
	text string,
) error {
	query := r.queries.WithTx(tx)
	return query.UpdateTranscriptionTextByID(ctx, db.UpdateTranscriptionTextByIDParams{
		ID:            taskID,
		Transcription: &text,
	})
}

func (r *RepositoryStruct) CreateParticipant(
	ctx context.Context,
	tx pgx.Tx,
	participantName *string,
	participantEmail string,
) error {
	query := r.queries.WithTx(tx)
	return query.CreateParticipant(ctx, db.CreateParticipantParams{
		Name:  participantName,
		Email: participantEmail,
	})
}

func (r *RepositoryStruct) GetParticipants(ctx context.Context) ([]db.Participant, error) {
	return r.queries.GetParticipants(ctx)
}

func (r *RepositoryStruct) DeleteParticipantByID(
	ctx context.Context,
	tx pgx.Tx,
	participantID uuid.UUID,
) error {
	query := r.queries.WithTx(tx)
	return query.DeleteParticipantByID(ctx, participantID)
}

func (r *RepositoryStruct) AssignParticipantToSpeaker(
	ctx context.Context,
	tx pgx.Tx,
	participantID *uuid.UUID,
	speakerID uuid.UUID,
) error {
	query := r.queries.WithTx(tx)
	err := query.AssignParticipantToSpeakerByID(ctx, db.AssignParticipantToSpeakerByIDParams{
		ID:            speakerID,
		ParticipantID: participantID,
	})
	return err
}

func (r *RepositoryStruct) CreateNewSpeakerForSegment(
	ctx context.Context,
	tx pgx.Tx,
	speaker int32,
	participantID *uuid.UUID,
	conversationID uuid.UUID,
) (uuid.UUID, error) {
	query := r.queries.WithTx(tx)
	return query.CreateNewSpeakerForSegment(ctx, db.CreateNewSpeakerForSegmentParams{
		ParticipantID:  participantID,
		ConversationID: conversationID,
		Speaker:        speaker,
	})
}

func (r *RepositoryStruct) GetSpeakerParticipantIDBySegmentID(
	ctx context.Context,
	tx pgx.Tx,
	speakerID uuid.UUID,
) (db.GetSpeakerIDAndParticipantIDBySegmentIDRow, error) {
	query := r.queries
	if tx != nil {
		query = r.queries.WithTx(tx)
	}
	return query.GetSpeakerIDAndParticipantIDBySegmentID(ctx, speakerID)
}

func (r *RepositoryStruct) GetSpeakerCountByConversationID(
	ctx context.Context,
	tx pgx.Tx,
	conversationID uuid.UUID,
) (int64, error) {
	query := r.queries
	if tx != nil {
		query = r.queries.WithTx(tx)
	}
	return query.GetSpeakerCountByConversationID(ctx, conversationID)
}

func (r *RepositoryStruct) AssignNewSpeakerToSegment(
	ctx context.Context,
	tx pgx.Tx,
	segmentID, speakerID uuid.UUID,
) error {
	query := r.queries.WithTx(tx)
	return query.AssignNewSpeakerToSegment(ctx, db.AssignNewSpeakerToSegmentParams{
		ID:        segmentID,
		SpeakerID: speakerID,
	})
}

func (r *RepositoryStruct) CountSegmentsWithSpeakerID(
	ctx context.Context,
	tx pgx.Tx,
	speakerID uuid.UUID,
) (int64, error) {
	query := r.queries
	if tx != nil {
		query = r.queries.WithTx(tx)
	}
	return query.GetCountSegmentsWithSpeakerID(ctx, speakerID)
}

func (r *RepositoryStruct) NullifySpeakerParticipantID(
	ctx context.Context,
	tx pgx.Tx,
	participantID *uuid.UUID,
) error {
	query := r.queries.WithTx(tx)
	return query.NullifySpeakerParticipantID(ctx, participantID)
}

func (r *RepositoryStruct) GetFullTranscriptionByConversationID(
	ctx context.Context,
	conversationID uuid.UUID,
) ([]db.GetFullTranscriptionByConversationIDRow, error) {
	return r.queries.GetFullTranscriptionByConversationID(ctx, conversationID)
}

func (r *RepositoryStruct) CreatePrompt(
	ctx context.Context,
	tx pgx.Tx,
	pormptName,
	prompt string,
) error {
	query := r.queries.WithTx(tx)
	return query.CreatePrompt(ctx,
		db.CreatePromptParams{
			PromptName: pormptName,
			Prompt:     prompt,
		})
}

func (r *RepositoryStruct) GetPrompts(
	ctx context.Context,
) ([]db.Prompt, error) {
	return r.queries.GetPrompts(ctx)
}

func (r *RepositoryStruct) DeletePromptByID(
	ctx context.Context,
	tx pgx.Tx,
	promptID uuid.UUID,
) error {
	query := r.queries.WithTx(tx)
	return query.DeletePromptByID(ctx, promptID)
}

func (r *RepositoryStruct) UpdatePromptByID(
	ctx context.Context,
	tx pgx.Tx,
	promptID uuid.UUID,
	promptName,
	prompt string,
) error {
	query := r.queries.WithTx(tx)
	return query.UpdatePromptByID(ctx, db.UpdatePromptByIDParams{
		ID:         promptID,
		Prompt:     prompt,
		PromptName: promptName,
	})
}

func (r *RepositoryStruct) GetPromptByID(
	ctx context.Context,
	promptID uuid.UUID,
) (db.Prompt, error) {
	return r.queries.GetPromptByID(ctx, promptID)
}

func (r *RepositoryStruct) UpdateParticipantByID(
	ctx context.Context,
	tx pgx.Tx,
	participantID uuid.UUID,
	participantName *string,
	participantEmail string,
) error {
	query := r.queries.WithTx(tx)
	return query.UpdateParticipantByID(ctx, db.UpdateParticipantByIDParams{
		ID:    participantID,
		Name:  participantName,
		Email: participantEmail,
	})
}

func (r *RepositoryStruct) GetParticipantByID(
	ctx context.Context,
	participantID uuid.UUID,
) (db.Participant, error) {
	return r.queries.GetParticipantByID(ctx, participantID)
}

func (r *RepositoryStruct) UpdateConversationNameByID(
	ctx context.Context,
	tx pgx.Tx,
	conversationID uuid.UUID,
	conversationName string,
) error {
	query := r.queries.WithTx(tx)
	return query.UpdateConversationNameByID(ctx, db.UpdateConversationNameByIDParams{
		ID:               conversationID,
		ConversationName: conversationName,
	})
}

func (r *RepositoryStruct) GetPromptByName(
	ctx context.Context,
	tx pgx.Tx,
	promptName string,
) (db.Prompt, error) {
	query := r.queries
	if tx != nil {
		query = r.queries.WithTx(tx)
	}
	return query.GetPromptByName(ctx, promptName)
}

func (r *RepositoryStruct) CreateSemiReport(
	ctx context.Context,
	tx pgx.Tx,
	conversationID,
	taskID,
	promptID uuid.UUID,
	partNum int,
) error {
	query := r.queries.WithTx(tx)
	return query.CreateSemiReport(
		ctx,
		db.CreateSemiReportParams{
			ConversationID: conversationID,
			TaskID:         taskID,
			PromptID:       promptID,
			PartNum: int32(partNum),
		},
	)
}

func (r *RepositoryStruct) GetConversationIDBySemiReportTaskID(
	ctx context.Context,
	taskID uuid.UUID,
) (uuid.UUID, error) {
	return r.queries.GetConversationIDBySemiReportTaskID(
		ctx,
		taskID,
	)
}

func (r *RepositoryStruct) GetCountOfUnSemiReportedParts(
	ctx context.Context,
	tx pgx.Tx,
	conversationID uuid.UUID,
) (int64, error) {
	query := r.queries
	if tx != nil {
		query = r.queries.WithTx(tx)
	}
	return query.GetCountOfUnSemiReportedParts(ctx, conversationID)
}

func (r *RepositoryStruct) UpdateSemiReportByTaskID(
	ctx context.Context,
	tx pgx.Tx,
	taskID uuid.UUID,
	semiReport string,
) error {
	query := r.queries.WithTx(tx)
	return query.UpdateSemiReportByTaskID(
		ctx,
		db.UpdateSemiReportByTaskIDParams{
			TaskID: taskID,
			SemiReport: &semiReport,
		},
	)
}