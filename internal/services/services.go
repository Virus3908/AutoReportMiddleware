package services

import (
	"main/internal/database"
	"main/internal/storage"

)

// type Services interface {

// }



type ServicesStruct struct {
	Conversations *ConversationsService
	Participant   *ParticipantService
	Promt         *PromtService
}

func NewService(db database.Database, storage storage.Storage) *ServicesStruct {
	return &ServicesStruct{
		Conversations: NewConversationService(db, storage),
		Participant:   NewParticipantService(db),
		Promt:         NewPromtService(db),
	}
}
