package repositories

import (
	"context"
	"main/internal/models"
	"main/internal/repositories/gen"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *RepositoryStruct) GetConversations(ctx context.Context) ([]models.Conversation, error) {
	conversations, err := r.queries.GetConversations(ctx)
	if err != nil {
		return nil, err
	}
	var result []models.Conversation
	for _, conversation := range conversations {
		result = append(result, models.Conversation{
			ID:          conversation.ID,
			ConversationName:       conversation.ConversationName,
			Status: conversation.	Status,
			CreatedAt:   conversation.CreatedAt,
			UpdatedAt:   conversation.UpdatedAt,
		})
	}
	return result, nil
}

func (r *RepositoryStruct) CreateConversation(ctx context.Context, fileURL, conversation_name string ) error {
	tx, err := r.database.StartTransaction(ctx)
	if err != nil {	
		return err
	}
	query := r.queries.WithTx(tx)
	err = query.CreateConversation(ctx, db.CreateConversationParams{
		FileUrl: 		fileURL,
		ConversationName: conversation_name,
	})
	if err != nil {
		r.database.RollbackTransactionIfExist(ctx, tx)
		return err
	}
	return r.database.CommitTransaction(ctx, tx)
}