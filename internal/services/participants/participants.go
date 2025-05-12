package participants

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"main/internal/common/interfaces"
	"main/internal/models"
	"main/internal/repositories"
	db "main/internal/repositories/gen"
)

type ParticipantService struct {
	Repo      *repositories.RepositoryStruct
	TxManager interfaces.TxManager
}

func NewParticipantService(repo *repositories.RepositoryStruct, txManager interfaces.TxManager) *ParticipantService {
	return &ParticipantService{
		Repo:      repo,
		TxManager: txManager,
	}
}

func (s *ParticipantService) CreateParticipant(
	ctx context.Context,
	participantPayload models.Participant,
) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.CreateParticipant(ctx, tx, &participantPayload.Name, participantPayload.Email)
	})
}

func (s *ParticipantService) GetParticipants(
	ctx context.Context,
) ([]db.Participant, error) {
	return s.Repo.GetParticipants(ctx)
}

func (s *ParticipantService) DeleteParticipantByID(
	ctx context.Context,
	participantID uuid.UUID,
) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		err := s.Repo.NullifySpeakerParticipantID(ctx, tx, &participantID)
		if err != nil {
			return fmt.Errorf("exec: Delete Participant\nfailed to nullify speaker participant ID: %w", err)
		}
		return s.Repo.DeleteParticipantByID(ctx, tx, participantID)
	})
}

func (s *ParticipantService) UpdateParticipantByID(
	ctx context.Context,
	participantID uuid.UUID,
	participantPayload models.Participant,
) error {
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.Repo.UpdateParticipantByID(ctx, tx, participantID, &participantPayload.Name, participantPayload.Email)
	})
}

func (s *ParticipantService) GetParticipantByID(
	ctx context.Context,
	participantID uuid.UUID,
) (db.Participant, error) {
	return s.Repo.GetParticipantByID(ctx, participantID)
}
