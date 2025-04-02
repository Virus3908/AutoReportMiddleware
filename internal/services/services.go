package services

import (
	"main/internal/database"
	"main/internal/kafka"
	"main/internal/services/conversations"
	"main/internal/services/crud"
	"main/internal/storage"
)

type ServicesStruct struct {
	CrudService *crud.CrudServicesStruct
	ConversationsService *conversations.ConversationsService
}

func NewServices(db database.Database, storage storage.Storage, kafka *kafka.Producer) *ServicesStruct{
	crudService := crud.NewService(db, storage)
	conversationsService := conversations.NewService(crudService.Conversation, storage)
	return &ServicesStruct{
		CrudService: crudService,
		ConversationsService: conversationsService,
	}
}