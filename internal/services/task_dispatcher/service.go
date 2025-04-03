package taskdispatcher

import (
	"context"
	"fmt"
	"encoding/json"
	"main/internal/kafka"
	"main/internal/services/crud"

	"github.com/google/uuid"
)

type TaskDispatcher struct {
	CRUD crud.CrudServicesStruct
	Kafka *kafka.Producer
}

const (
	TaskTypeConvert int32 = 1
	TaskTypeDiarize int32 = 2
	TaskTypeTranscribe int32 = 3
)

type ConvertTask struct {
	ConvertID uuid.UUID `json:"convert_id"`
	FileURL string `json:"file_url"`
	TaskType int32 `json:"task_type"`
	CallbackPostfix string `json:"callback_postfix"`
}

func NewTaskDispatcher(crud crud.CrudServicesStruct, kafka *kafka.Producer) *TaskDispatcher {
	return &TaskDispatcher{
		CRUD: crud,
		Kafka: kafka,
	}
}

func (s *TaskDispatcher) CreateConvertTask(ctx context.Context, conversation_id uuid.UUID) error {
	fileURL, err := s.CRUD.Conversation.GetFileURLByID(ctx, conversation_id)
	if err != nil {
		return fmt.Errorf("error creating convert task: %s", err)
	}

	convert_id, err := s.CRUD.Convert.Create(ctx, conversation_id)
	if err != nil {
		return fmt.Errorf("error creating row in db: %s", err)
	}

	convertTask := ConvertTask {
		FileURL: fileURL,
		ConvertID: *convert_id,
		TaskType: TaskTypeConvert,
		CallbackPostfix: "/api/task/update/convert/",
	}
	payload, err := json.Marshal(convertTask)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}
	
	err = s.Kafka.SendMessage(ctx, conversation_id.String(), string(payload))
	if err != nil {
		return fmt.Errorf("error sending Kafka message: %w", err)
	}

	return nil
}