package services

import (
	"main/internal/database"
	"main/internal/storage"
)

type ConversationsService struct {
	DB database.Database
	Storage storage.Storage
}

func NewConversationService(db database.Database, storage storage.Storage) *ConversationsService{
	return &ConversationsService{
		DB: db,
		Storage: storage,
	}
}