package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type DataBase struct {
	Config DBConfig
	Pool   *pgxpool.Pool
}

func New(cfg DBConfig) (*DataBase, error) {

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

	return &DataBase{
		Pool:   pool,
		Config: cfg,
	}, nil
}

func (d *DataBase) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}
