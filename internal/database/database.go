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
