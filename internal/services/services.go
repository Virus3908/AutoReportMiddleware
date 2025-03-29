package services

import (
	"main/internal/database"
	"main/internal/storage"
)

type ServicesStruct struct {
	Conversations *ConversationsService
}

func NewService(db database.Database, storage storage.Storage) *ServicesStruct {
	return &ServicesStruct{
		Conversations: NewConversationService(db, storage),
	}
}