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

const (
	TaskTypeConvert    int32 = 1
	TaskTypeDiarize    int32 = 2
	TaskTypeTranscribe int32 = 3
)

type ConvertTask struct {
	ConvertID       uuid.UUID `json:"convert_id"`
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

	return s.CRUD.Convert.CreateWithFN(ctx, conversationID, func(convertID uuid.UUID) error {
		convertTask := ConvertTask{
			FileURL:         fileURL,
			ConvertID:       convertID,
			TaskType:        TaskTypeConvert,
			CallbackPostfix: "/api/task/update/convert/",
		}
		payload, err := json.Marshal(convertTask)
		if err != nil {
			return fmt.Errorf("error marshaling payload: %w", err)
		}
		return s.Kafka.SendMessage(ctx, conversationID.String(), string(payload))
	})
}
