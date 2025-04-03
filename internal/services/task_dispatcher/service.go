package taskdispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"main/internal/kafka"
	"main/internal/services/crud"
)

type TaskDispatcher struct {
	CRUD  crud.CrudServicesStruct
	Kafka *kafka.Producer
}

type ConvertTask struct {
	TaskID          uuid.UUID `json:"task_id"`
	FileURL         string    `json:"file_url"`
	TaskType        int32     `json:"task_type"`
	CallbackPostfix string    `json:"callback_postfix"`
}

func NewTaskDispatcher(crud crud.CrudServicesStruct, kafka *kafka.Producer) *TaskDispatcher {
	return &TaskDispatcher{
		CRUD:  crud,
		Kafka: kafka,
	}
}

func (s *TaskDispatcher) CreateConvertTask(ctx context.Context, conversationID uuid.UUID) error {
	fileURL, err := s.CRUD.Conversation.GetFileURLByID(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("error retrieving file URL: %w", err)
	}

	return s.CRUD.Task.CreateConvertTask(ctx, conversationID, func(taskID uuid.UUID, taskType int32) error {
		convertTask := ConvertTask{
			FileURL:         fileURL,
			TaskID:          taskID,
			TaskType:        taskType,
			CallbackPostfix: "/api/task/update/convert/",
		}
		payload, err := json.Marshal(convertTask)
		if err != nil {
			return fmt.Errorf("error marshaling payload: %w", err)
		}
		return s.Kafka.SendMessage(ctx, conversationID.String(), string(payload))
	})
}
