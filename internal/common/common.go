package common

import (
	"github.com/jackc/pgx/v5/pgtype"
	"context"
	"main/internal/database"
	"main/internal/database/queries"
	"github.com/google/uuid"
)

func StrToPGUUID(strID string) (pgtype.UUID, error) {
	id, err := uuid.Parse(strID)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: id, Valid: true}, nil
}

func StartTransaction(db *database.DataBase) (*queries.Queries, func(), func() (error), error) {
	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}
	rollback := func ()  {
		tx.Rollback(context.Background())
	}
	commit := func() (error) {
		return tx.Commit(context.Background())
	}
	return (&queries.Queries{}).WithTx(tx), rollback, commit, nil
}