package services

import (
	"main/internal/database"
	"main/internal/storage"
)

type ServicesStruct struct {
	Conversation *ConversationsService
	Participant  *ParticipantService
	Promt        *PromtService
	Convert      *ConvertService
}

func NewService(db database.Database, storage storage.Storage) *ServicesStruct {
	return &ServicesStruct{
		Conversation: NewConversationService(db, storage),
		Participant:  NewParticipantService(db),
		Promt:        NewPromtService(db),
	}
}
