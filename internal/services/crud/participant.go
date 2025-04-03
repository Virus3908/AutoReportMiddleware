package crud

import (
	"context"
	"main/internal/database"
	"main/internal/repositories"
	"github.com/google/uuid"
)

// вот эти все файлы у тебя по сути репозитория, причем однотипнейшная

type ParticipantCRUD struct {
	DB database.Database
}

func NewParticipantCRUD(db database.Database) *ParticipantCRUD {
	return &ParticipantCRUD{DB: db}
}

func (s *ParticipantCRUD) GetByID(ctx context.Context, id uuid.UUID) (repositories.Participant, error) {
	return s.DB.NewQuery().GetParticipantByID(ctx, id)
}

func (s *ParticipantCRUD) GetAll(ctx context.Context) ([]repositories.Participant, error) {
	return s.DB.NewQuery().GetParticipants(ctx)
}

func (s *ParticipantCRUD) Create(ctx context.Context, payload repositories.CreateParticipantParams) error {
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.CreateParticipant(ctx, payload)
	})
}

func (s *ParticipantCRUD) Update(ctx context.Context, id uuid.UUID, payload repositories.UpdateParticipantByIDParams) error {
	payload.ID = id
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.UpdateParticipantByID(ctx, payload)
	})
}

func (s *ParticipantCRUD) Delete(ctx context.Context, id uuid.UUID) error {
	return s.DB.WithTx(ctx, func(q *repositories.Queries) error {
		return q.DeleteParticipantByID(ctx, id)
	})
}