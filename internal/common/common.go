package common

import (
	"context"
	// "fmt"
	"main/internal/database"
	"main/internal/database/queries"

	"github.com/google/uuid"
	// "github.com/jackc/pgx/v5/pgtype"
)

func StrToPGUUID(strID string) (uuid.UUID, error) {
	id, err := uuid.Parse(strID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
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