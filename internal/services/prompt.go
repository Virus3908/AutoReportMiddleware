package services

import (
	"context"
	"main/internal/database"
	"main/internal/repositories"

	"github.com/google/uuid"
)

type PromptService struct {
	DB database.Database
}

func NewPromptService(db database.Database) *PromptService {
	return &PromptService{DB: db}
}

func (s *PromptService) GetByID(ctx context.Context, id uuid.UUID) (repositories.Prompt, error) {
	return s.DB.NewQuery().GetPromptByID(ctx, id)
}

func (s *PromptService) GetAll(ctx context.Context) ([]repositories.Prompt, error) {
	return s.DB.NewQuery().GetPrompts(ctx)
}

func (s *PromptService) Create(ctx context.Context, payload string) error {
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.CreatePrompt(ctx, payload)
	})
}

func (s *PromptService) Update(ctx context.Context, id uuid.UUID, payload repositories.UpdatePromptByIDParams) error {
	payload.ID = id
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.UpdatePromptByID(ctx, payload)
	})
}

func (s *PromptService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.DeletePromptByID(ctx, id)
	})
}
