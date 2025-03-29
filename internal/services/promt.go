package services

import (
	"context"
	"main/internal/database"
	"main/internal/repositories"
	"github.com/google/uuid"
)


type PromtService struct {
	DB database.Database
}

func NewPromtService(db database.Database) *PromtService {
	return &PromtService{DB: db}
}

func (s *PromtService) GetByID(ctx context.Context, id uuid.UUID) (repositories.Promt, error) {
	return s.DB.NewQuerry().GetPromtByID(ctx, id)
}

func (s *PromtService) Get(ctx context.Context) ([]repositories.Promt, error) {
	return s.DB.NewQuerry().GetPromts(ctx)
}

func (s *PromtService) Create(ctx context.Context, payload string) error {
	tx, rollback, commit, err := s.DB.StartTransaction()
	if err != nil {
		return err
	}
	defer rollback()
	err = tx.CreatePromt(ctx, payload)
	if err != nil {
		return err
	}
	commit()
	return nil
}

func (s *PromtService) Update(ctx context.Context, id uuid.UUID, payload repositories.UpdatePromtByIDParams) error {
	payload.ID = id
	tx, rollback, commit, err := s.DB.StartTransaction()
	if err != nil {
		return err
	}
	defer rollback()
	err = tx.UpdatePromtByID(ctx, payload)
	if err != nil {
		return err
	}
	commit()
	return nil
}

func (s *PromtService) Delete(ctx context.Context, id uuid.UUID) error {
	tx, rollback, commit, err := s.DB.StartTransaction()
	if err != nil {
		return err
	}
	defer rollback()
	err = tx.DeletePromtByID(ctx, id)
	if err != nil {
		return err
	}
	commit()
	return nil
}