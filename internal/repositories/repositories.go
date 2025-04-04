package repositories

import (
	"context"
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

func (r *RepositoryStruct) GetConversations(ctx context.Context) ([]db.Conversation, error) {
	return r.queries.GetConversations(ctx)
}