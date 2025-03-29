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