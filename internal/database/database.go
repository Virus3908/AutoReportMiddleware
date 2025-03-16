package database

import (
	"context"
	"fmt"
	"main/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/jackc/pgx/v5/pgxpool"
)

// var db *pgxpool.Pool

type DataBase struct {
	Pool *pgxpool.Pool
}

func New(cfg config.DBConnection) (*DataBase, error) {

	connectionInfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database)

	pool, err := pgxpool.New(context.Background(), connectionInfo)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Ошибка подключения к базе данных: %s", err))
		
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Ошибка пинга базы данных: %s", err))
	}

	return &DataBase{Pool: pool}, nil
}

func (d *DataBase) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}