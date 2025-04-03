package services

import (
	"main/internal/database"
	"main/internal/kafka"
	"main/internal/services/conversations"
	"main/internal/services/crud"
	taskdispatcher "main/internal/services/task_dispatcher"
	"main/internal/storage"
)

type ServicesStruct struct {
	CrudService *crud.CrudServicesStruct
	ConversationsService *conversations.ConversationsService
	TaskService *taskdispatcher.TaskDispatcher
}

func NewServices(db database.Database, storage storage.Storage, kafka *kafka.Producer) *ServicesStruct{
	crudService := crud.NewService(db, storage)
	conversationsService := conversations.NewService(crudService.Conversation, storage)
	taskService := taskdispatcher.NewTaskDispatcher(*crudService, kafka)
	return &ServicesStruct{
		CrudService: crudService,
		ConversationsService: conversationsService,
		TaskService: taskService,
	}
}