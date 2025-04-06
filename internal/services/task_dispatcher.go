package services

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/repositories"
	"github.com/google/uuid"
)

type TaskDispatcher struct {
	Repo *repositories.RepositoryStruct
	Messenger MessageClient
	TxManager TxManager
}

const (
	ConvertTask = 1
	DiarizeTask = 2
	TranscribeTask = 3
)

type ConvertMessage struct {
	TaskID          uuid.UUID `json:"task_id"`
	FileURL         string    `json:"file_url"`
	TaskType        int32     `json:"task_type"`
	CallbackPostfix string    `json:"callback_postfix"`
}

func NewTaskDispatcher(repo *repositories.RepositoryStruct, messenger MessageClient, txManager TxManager) *TaskDispatcher {
	return &TaskDispatcher{
		Repo: repo,
		Messenger: messenger,
		TxManager: txManager,
	}
}

func (s *TaskDispatcher) CreateConvertTask(ctx context.Context, conversationID uuid.UUID) error {
	fileURL, err := s.Repo.GetConversationFileURL(ctx, conversationID)
	if err != nil {
		return err
	}
	tx, err := s.TxManager.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer s.TxManager.RollbackTransactionIfExist(ctx, tx)

	taskID, err := s.Repo.CreateTask(ctx, tx, ConvertTask)
	if err != nil {
		return err
	}
	err = s.Repo.CreateConvert(ctx, tx, taskID, conversationID)
	if err != nil {
		return err
	}
	convertMessage := ConvertMessage{
		TaskID:          taskID,
		FileURL:         fileURL,
		TaskType:        ConvertTask,
		CallbackPostfix: "/api/task/update/convert/",
	}
	convertMessageJSON, err := json.Marshal(convertMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal convert message: %w", err)
	}
	err = s.Messenger.SendMessage(ctx, conversationID.String(), string(convertMessageJSON))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return s.TxManager.CommitTransaction(ctx, tx)
}