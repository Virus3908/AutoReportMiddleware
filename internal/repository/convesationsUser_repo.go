package repository

import (
	"main/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConversationsUserRepo struct {
	Pool *pgxpool.Pool
}

func NewConversationsUserRepo(pool *pgxpool.Pool) *ConversationsUserRepo {
	return &ConversationsUserRepo{Pool: pool}
}
