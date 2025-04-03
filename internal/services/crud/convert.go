package crud

import (
	"context"
	"github.com/google/uuid"
	"main/internal/database"
	"main/internal/repositories"
)

type ConvertCRUD struct {
	DB database.Database
}

func NewConvertCRUD(db database.Database) *ConvertCRUD {
	return &ConvertCRUD{
		DB: db,
	}
}

func (s *ConvertCRUD) GetAll(ctx context.Context) ([]repositories.Convert, error) {
	return s.DB.NewQuery().GetConvert(ctx)
}

func (s *ConvertCRUD) GetByID(ctx context.Context, id uuid.UUID) (repositories.Convert, error) {
	return s.DB.NewQuery().GetConvertByID(ctx, id)
}

func (s *ConvertCRUD) Create(ctx context.Context, payload uuid.UUID) (*uuid.UUID, error) {
	var id uuid.UUID

	err := s.DB.WithTx(ctx, func(tx *repositories.Queries) error {
		var err error
		id, err = tx.CreateConvert(ctx, payload)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *ConvertCRUD) Update(ctx context.Context, id uuid.UUID, payload repositories.UpdateConvertByTaskIDParams) error {
	payload.ID = id
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.UpdateConvertByTaskID(ctx, payload)
	})
}

func (s *ConvertCRUD) Delete(ctx context.Context, id uuid.UUID) error {
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.DeleteConvertByID(ctx, id)
	})
}

func (s *ConvertCRUD) DeleteByForgeinKey(ctx context.Context, forgeinKeyID uuid.UUID) (*uuid.UUID, error) {
	var id uuid.UUID

	err := s.DB.WithTx(ctx, func(tx *repositories.Queries) error {
		var err error
		id, err = tx.DeleteConvertByForgeinID(ctx, forgeinKeyID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &id, nil
}