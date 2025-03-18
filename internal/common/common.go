package common

import (
	"context"
	"main/internal/database"
	"main/internal/repositories"
)

func StartTransaction(db *database.DataBase) (*repositories.Queries, func(), func() error, error) {
	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}
	rollback := func() {
		tx.Rollback(context.Background())
	}
	commit := func() error {
		return tx.Commit(context.Background())
	}
	return (&repositories.Queries{}).WithTx(tx), rollback, commit, nil
}
