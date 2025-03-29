package services

import (
	"context"
	"main/internal/database"
	"main/internal/repositories"
	"github.com/google/uuid"
)

type ParticipantService struct {
	DB database.Database
}

func NewParticipantService(db database.Database) *ParticipantService {
	return &ParticipantService{DB: db}
}

func (s *ParticipantService) GetByID(ctx context.Context, id uuid.UUID) (repositories.Participant, error) {
	return s.DB.NewQuerry().GetParticipantByID(ctx, id)
}

func (s *ParticipantService) Get(ctx context.Context) ([]repositories.Participant, error) {
	return s.DB.NewQuerry().GetParticipants(ctx)
}

func (s *ParticipantService) Create(ctx context.Context, payload repositories.CreateParticipantParams) error {
	tx, rollback, commit, err := s.DB.StartTransaction()
	if err != nil {
		return err
	}
	defer rollback()
	err = tx.CreateParticipant(ctx, payload)
	if err != nil {
		return err
	}
	commit()
	return nil
}

func (s *ParticipantService) Update(ctx context.Context, id uuid.UUID, payload repositories.UpdateParticipantByIDParams) error {
	payload.ID = id
	tx, rollback, commit, err := s.DB.StartTransaction()
	if err != nil {
		return err
	}
	defer rollback()
	err = tx.UpdateParticipantByID(ctx, payload)
	if err != nil {
		return err
	}
	commit()
	return nil
}

func (s *ParticipantService) Delete(ctx context.Context, id uuid.UUID) error {
	tx, rollback, commit, err := s.DB.StartTransaction()
	if err != nil {
		return err
	}
	defer rollback()
	err = tx.DeleteParticipantByID(ctx, id)
	if err != nil {
		return err
	}
	commit()
	return nil
}