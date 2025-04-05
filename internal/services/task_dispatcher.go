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

func NewTaskDispatcher(repo *repositories.RepositoryStruct, messenger MessageClient) *TaskDispatcher {
	return &TaskDispatcher{
		Repo: repo,
		Messenger: messenger,
	}
}

func (s *TaskDispatcher) CreateConvertTask(ctx context.Context, conversationID uuid.UUID) error {
	fileURL, err := s.Repo.GetConversationFileURL(ctx, conversationID)
	if err != nil {
		return err
	}
	return s.Repo.CreateTask(
		ctx,
		conversationID,
		fileURL,
		ConvertTask,
		func(taskID uuid.UUID) error {
			taskMSG := ConvertMessage{
				TaskID: taskID,
				FileURL: fileURL,
				TaskType: ConvertTask,
				CallbackPostfix: "/api/task/update/convert/",
			}
			payload, err := json.Marshal(taskMSG)
			if err != nil {
				return fmt.Errorf("error marshaling payload: %w", err)
			}
			return s.Messenger.SendMessage(ctx, conversationID.String(), string(payload))
		},
	)
}