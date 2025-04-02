package crud

import (
	"context"
	"main/internal/database"
	"main/internal/repositories"

	"github.com/google/uuid"
)

type PromptCRUD struct {
	DB database.Database
}

func NewPromptCRUD(db database.Database) *PromptCRUD {
	return &PromptCRUD{DB: db}
}

func (s *PromptCRUD) GetByID(ctx context.Context, id uuid.UUID) (repositories.Prompt, error) {
	return s.DB.NewQuery().GetPromptByID(ctx, id)
}

func (s *PromptCRUD) GetAll(ctx context.Context) ([]repositories.Prompt, error) {
	return s.DB.NewQuery().GetPrompts(ctx)
}

func (s *PromptCRUD) Create(ctx context.Context, payload string) error {
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.CreatePrompt(ctx, payload)
	})
}

func (s *PromptCRUD) Update(ctx context.Context, id uuid.UUID, payload repositories.UpdatePromptByIDParams) error {
	payload.ID = id
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.UpdatePromptByID(ctx, payload)
	})
}

func (s *PromptCRUD) Delete(ctx context.Context, id uuid.UUID) error {
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.DeletePromptByID(ctx, id)
	})
}
