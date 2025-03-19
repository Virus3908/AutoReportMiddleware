package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	StartTransaction() (*repositories.Queries, func(), func() error, error)
	NewQuerry() *repositories.Queries 
}

type PGdatabase struct {
	Config DBConfig
	Pool   *pgxpool.Pool
}

func NewDatabase(cfg DBConfig) (*PGdatabase, error) {

	connectionInfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database)

	pool, err := pgxpool.New(context.Background(), connectionInfo)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Database connection error: %s", err))

	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Database ping error: %s", err))
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

func (d *PGdatabase) StartTransaction() (*repositories.Queries, func(), func() error, error) {
	tx, err := d.Pool.Begin(context.Background())
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

func (d *PGdatabase) NewQuerry() *repositories.Queries {
	return repositories.New(d.Pool)
}