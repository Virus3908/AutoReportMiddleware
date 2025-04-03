package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"main/internal/repositories"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Database interface {
	CloseConnection()
	// startTransaction(context.Context) (*Tx, error)
	WithTx(ctx context.Context, fn func(q *repositories.Queries) error) error
	NewQuery() *repositories.Queries
}

type Tx struct {
	Queries  *repositories.Queries
	Rollback func() error
	Commit   func() error
} // тут могу ошибаться, но у пгх уже есть все эти методы, вообще транзакция сделана тут плохо прямо, 

type Transactioned interface { // эту хрень по умолчанию в хендлеры можно инклудить например
	StartTransaction(context.Context) (pgx.Tx, error)
	StartNestedTransaction(ctx context.Context, tx pgx.Tx) (pgx.Tx, error)
	CommitTransaction(context.Context, pgx.Tx) error
	RollbackTransactionIfExist(context.Context, pgx.Tx) error
}

import (
	db "generated from sqlc"
)

func New(meter *metric.Meter, connPool *pgxpool.Pool, jwtGenerator JWTGenerator, storageServiceBaseURL string) *Repository {
	return &Repository{
		// meter:                 meter, ну вдруг тебе надо как-то
		queries:               db.New(connPool), // 
		connPool:              connPool,
		storageServiceBaseURL: storageServiceBaseURL,
	}
}

func (r *Repository) StartTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.connPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *Repository) StartNestedTransaction(ctx context.Context, tx pgx.Tx) (pgx.Tx, error) {
	return tx.Begin(ctx)
}

func (r *Repository) CommitTransaction(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (r *Repository) RollbackTransactionIfExist(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (r *Repository) SomeRequestWithTx( // это в репе уже
	ctx context.Context, tx pgx.Tx, ...,
) error {
	queries := r.queries.WithTx(tx)

	_, err := queries.SomeRequest(ctx, params)
	return err
}

func (r *Repository) SomeRequestWhereYouNeedTransaction( // это в репе уже
	ctx context.Context, ...,
) error {
	tx, err := r.StartTransaction(ctx)
	if err != nil {
		return err
	}

	queries := r.queries.WithTx(tx)

	_, err := queries.SomeRequest(ctx, params)
	return err
}

func (r *Repository) SomeRequestWhichCanBeWithOrWithoutTransaction( // это в репе уже
	ctx context.Context, tx pgx.Tx, ...,
) error {
	queries := r.queries
	if tx != nil {
		queries = r.queries.WithTx(tx)
	}

	_, err := queries.SomeRequest(ctx, params)
	return err
}

type PGdatabase struct {
	Config DBConfig
	Pool   *pgxpool.Pool
}

func NewDatabase(ctx context.Context, cfg DBConfig) (*PGdatabase, error) {

	connectionInfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database)

	pool, err := pgxpool.New(ctx, connectionInfo)

	if err != nil {
		return nil, fmt.Errorf("Database connection error: %s", err)

	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Database ping error: %s", err)
	}

	return &PGdatabase{
		Pool:   pool,
		Config: cfg,
	}, nil
}

func (d *PGdatabase) CloseConnection() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}

func (d *PGdatabase) startTransaction(ctx context.Context) (*Tx, error) {
	tx, err := d.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Tx{
		Queries: (&repositories.Queries{}).WithTx(tx),
		Rollback: func() error {
			return tx.Rollback(ctx)
		},
		Commit: func() error {
			return tx.Commit(ctx)
		},
	}, nil
}

func (d *PGdatabase) NewQuery() *repositories.Queries {
	return repositories.New(d.Pool)
}

func (d *PGdatabase) WithTx(ctx context.Context, fn func(q *repositories.Queries) error) error {
	tx, err := d.startTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx.Queries); err != nil {
		return err
	}

	return tx.Commit()
}
