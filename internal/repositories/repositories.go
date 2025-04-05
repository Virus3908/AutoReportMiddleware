package repositories

import (
	"context"
	"main/internal/repositories/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

type Database interface {
	GetPool() *pgxpool.Pool
	StartTransaction(ctx context.Context) (pgx.Tx, error)
	StartNestedTransaction(ctx context.Context, tx pgx.Tx) (pgx.Tx, error)
	CommitTransaction(ctx context.Context, tx pgx.Tx) error
	RollbackTransactionIfExist(ctx context.Context, tx pgx.Tx) error
}

type RepositoryStruct struct {
	database Database
	queries  *db.Queries
}

func New(dbDriver Database) *RepositoryStruct {
	return &RepositoryStruct{
		database: dbDriver,
		queries:  db.New(dbDriver.GetPool()),
	}
}

func (r *RepositoryStruct) GetConversations(ctx context.Context) ([]db.Conversation, error) {
	return r.queries.GetConversations(ctx)
}

func (r *RepositoryStruct) CreateConversation(ctx context.Context, fileURL, conversation_name string) error {
	tx, err := r.database.StartTransaction(ctx)
	if err != nil {
		return err
	}
	query := r.queries.WithTx(tx)
	err = query.CreateConversation(ctx, db.CreateConversationParams{
		FileUrl:          fileURL,
		ConversationName: conversation_name,
	})
	if err != nil {
		r.database.RollbackTransactionIfExist(ctx, tx)
		return err
	}
	return r.database.CommitTransaction(ctx, tx)
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

func (r *RepositoryStruct) CreateTask(
	ctx context.Context, 
	conversationID uuid.UUID, 
	fileURL string, 
	taskType int32,
	fn func(taskID uuid.UUID) error,
) error {
	tx, err := r.database.StartTransaction(ctx)
	if err != nil {
		return err
	}
	query := r.queries.WithTx(tx)
	taskID, err := query.CreateTask(ctx, taskType)
	if err != nil {
		r.database.RollbackTransactionIfExist(ctx, tx)
		return err
	}
	err = query.CreateConvert(ctx, db.CreateConvertParams{
		ConversationsID: conversationID,
		TaskID: taskID,
	})
	if err != nil {
		r.database.RollbackTransactionIfExist(ctx, tx)
		return err
	}
	err = fn(taskID)
	if err != nil {
		r.database.RollbackTransactionIfExist(ctx, tx)
		return err
	}
	return r.database.CommitTransaction(ctx, tx)

}