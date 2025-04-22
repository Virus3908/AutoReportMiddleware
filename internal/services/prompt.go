package services

import (
	"context"
	"main/internal/models"
	"main/internal/repositories"
	db "main/internal/repositories/gen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PromptService struct {
	Repo      *repositories.RepositoryStruct
	TxManager TxManager
}

func NewPromptService(repo *repositories.RepositoryStruct, txManager TxManager) *PromptService {
	return &PromptService{
		Repo:      repo,
		TxManager: txManager,
	}
}

func (s *PromptService) GetPrompts(ctx context.Context) ([]db.Prompt, error) {
	return s.Repo.GetPrompts(ctx)
}

func (s *PromptService) CreatePrompt(ctx context.Context, prompt models.Prompt) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.CreatePrompt(ctx, tx, prompt.PromptName, prompt.Prompt)
	})
}

func (s *PromptService) DeletePromptByID(ctx context.Context, promptID uuid.UUID) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.DeletePromptByID(ctx, tx, promptID)
	})
}

func (s *PromptService) UpdatePromptByID(ctx context.Context, promptID uuid.UUID, prompt models.Prompt) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.UpdatePromptByID(ctx, tx, promptID, prompt.PromptName, prompt.Prompt)
	})
}