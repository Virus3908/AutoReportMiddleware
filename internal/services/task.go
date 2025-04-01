package services

import "main/internal/database"

type TaskService struct {
	DB database.Database
}

func NewTaskService(db database.Database) *TaskService{
	return &TaskService{
		DB: db,
	}
}
