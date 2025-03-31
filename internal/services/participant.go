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
	return s.DB.NewQuery().GetParticipantByID(ctx, id)
}

func (s *ParticipantService) GetAll(ctx context.Context) ([]repositories.Participant, error) {
	return s.DB.NewQuery().GetParticipants(ctx)
}

func (s *ParticipantService) Create(ctx context.Context, payload repositories.CreateParticipantParams) error {
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.CreateParticipant(ctx, payload)
	})
}

func (s *ParticipantService) Update(ctx context.Context, id uuid.UUID, payload repositories.UpdateParticipantByIDParams) error {
	payload.ID = id
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.UpdateParticipantByID(ctx, payload)
	})
}

func (s *ParticipantService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.DeleteParticipantByID(ctx, id)
	})
}