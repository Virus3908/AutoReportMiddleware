package services

import (
	"context"
	"main/internal/repositories"
	db "main/internal/repositories/gen"
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