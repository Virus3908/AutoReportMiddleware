package services

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TaskDispatcher struct {
	Repo      *repositories.RepositoryStruct
	Messenger MessageClient
	TxManager TxManager
}

const (
	ConvertTask    = 1
	DiarizeTask    = 2
	TranscribeTask = 3
)

type MessageWithFileURL struct {
	TaskID          uuid.UUID `json:"task_id"`
	FileURL         string    `json:"file_url"`
	CallbackPostfix string    `json:"callback_postfix"`
}

func NewTaskDispatcher(repo *repositories.RepositoryStruct, messenger MessageClient, txManager TxManager) *TaskDispatcher {
	return &TaskDispatcher{
		Repo:      repo,
		Messenger: messenger,
		TxManager: txManager,
	}
}

func (s *TaskDispatcher) CreateConvertTask(ctx context.Context, conversationID uuid.UUID) error {
	fileURL, err := s.Repo.GetConversationFileURL(ctx, conversationID)
	if err != nil {
		return err
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		taskID, err := s.Repo.CreateTask(ctx, tx, ConvertTask)
		if err != nil {
			return err
		}
		err = s.Repo.CreateConvert(ctx, tx, taskID, conversationID)
		if err != nil {
			return err
		}
		convertMessage := MessageWithFileURL{
			TaskID:          taskID,
			FileURL:         fileURL,
			CallbackPostfix: "/api/task/update/convert/",
		}
		convertMessageJSON, err := json.Marshal(convertMessage)
		if err != nil {
			return fmt.Errorf("failed to marshal convert message: %w", err)
		}
		return s.Messenger.SendMessage(ctx, ConvertTask, conversationID.String(), string(convertMessageJSON))
	})
}

func (s *TaskDispatcher) CreateDiarizeTask(ctx context.Context, conversationID uuid.UUID) error {
	response, err := s.Repo.GetConvertFileURLByConversationID(ctx, conversationID)
	if err != nil {
		return err
	}
	if response.FileUrl == nil {
		return fmt.Errorf("file is not converted")
	}
	return s.TxManager.WithTx(ctx, func(tx pgx.Tx) error {
		taskID, err := s.Repo.CreateTask(ctx, tx, DiarizeTask)
		if err != nil {
			return err
		}
		err = s.Repo.CreateDiarize(ctx, tx, response.ID, taskID)
		if err != nil {
			return err
		}
		diarizeMessage := MessageWithFileURL{
			TaskID:          taskID,
			FileURL:         *response.FileUrl,
			CallbackPostfix: "/api/task/update/diarize/",
		}
		diarizeMessageJSON, err := json.Marshal(diarizeMessage)
		if err != nil {
			return fmt.Errorf("failed to marshal convert message: %w", err)
		}

		return s.Messenger.SendMessage(ctx, DiarizeTask, conversationID.String(), string(diarizeMessageJSON))
	})
}
