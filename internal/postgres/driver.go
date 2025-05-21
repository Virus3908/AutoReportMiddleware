package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type PGDatabase struct {
	config   DBConfig
	connPool *pgxpool.Pool
}

func New(ctx context.Context, cfg DBConfig) (*PGDatabase, error) {
	database := PGDatabase{
		config: cfg,
	}
	err := database.createConnection(ctx)
	if err != nil {
		return nil, err
	}
	return &database, nil
}

func (d *PGDatabase) createConnection(ctx context.Context) error {
	connectionInfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		d.config.User,
		d.config.Password,
		d.config.Host,
		d.config.Port,
		d.config.Database)

	pool, err := pgxpool.New(ctx, connectionInfo)
	if err != nil {
		return fmt.Errorf("database connection error: %s", err)

	}
	if err = pool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping error: %s", err)
	}
	d.connPool = pool
	return nil
}

func (d *PGDatabase) GetPool() *pgxpool.Pool {
	return d.connPool
}

func (d *PGDatabase) Close() {
	if d.connPool != nil {
		d.connPool.Close()
	}
}

func (d *PGDatabase) StartTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := d.connPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (d *PGDatabase) StartNestedTransaction(ctx context.Context, tx pgx.Tx) (pgx.Tx, error) {
	return tx.Begin(ctx)
}

func (d *PGDatabase) CommitTransaction(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (d *PGDatabase) RollbackTransactionIfExist(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (d *PGDatabase) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := d.StartTransaction(ctx)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		d.RollbackTransactionIfExist(ctx, tx)
		return err
	}
	return d.CommitTransaction(ctx, tx)
}