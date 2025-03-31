package services

import (
	"context"
	"main/internal/database"
	"main/internal/repositories"
	"github.com/google/uuid"
)

type ConvertService struct {
	DB database.Database
}

func NewConvertService(db database.Database) *ConvertService {
	return &ConvertService{
		DB: db,
	}
}

func (s *ConvertService) Create(ctx context.Context, payload uuid.UUID) error {
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.CreateConvert(ctx, payload)
	})
}

func (s *ConvertService) Update(ctx context.Context, task_id uuid.UUID, payload repositories.UpdateConvertByTaskIDParams) error {
	payload.TaskID = task_id
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.UpdateConvertByTaskID(ctx, payload)
	})
}
