package crud

import (
	"context"
	"main/internal/database"
	"main/internal/repositories"

	"github.com/google/uuid"
)

type TaskCRUD struct {
	DB database.Database
}

const (
	TaskTypeConvert    int32 = 1
	TaskTypeDiarize    int32 = 2
	TaskTypeTranscribe int32 = 3
)

func NewTaskCRUD(db database.Database) *TaskCRUD {
	return &TaskCRUD{
		DB: db,
	}
}

func (s *TaskCRUD) GetAll(ctx context.Context) ([]repositories.Task, error) {
	return s.DB.NewQuery().GetTasks(ctx)
}

func (s *TaskCRUD) GetByID(ctx context.Context, id uuid.UUID) (repositories.Task, error) {
	return s.DB.NewQuery().GetTaskByID(ctx, id)
}

func (s *TaskCRUD) CreateConvertTask(
	ctx context.Context,
	conversationID uuid.UUID,
	fn func(id uuid.UUID, taskType int32) error,
) error {
	return s.DB.WithTx(ctx, func(tx *repositories.Queries) error {
		taskID, err := tx.CreateTask(ctx, TaskTypeConvert)
		if err != nil {
			return err
		}
		err = tx.CreateConvert(ctx, repositories.CreateConvertParams{
			ConversationsID: conversationID,
			TaskID:          taskID,
		})
		if err != nil {
			return err
		}
		return fn(taskID, TaskTypeConvert)
	})
}
