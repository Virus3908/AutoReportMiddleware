package services

import (
	"main/internal/common/interfaces"
	"main/internal/repositories"
	"main/internal/services/conversations"
	"main/internal/services/participants"
	"main/internal/services/prompt"
	"main/internal/services/task_dispatcher"
)

type ServiceStruct struct {
	Conversations *conversations.ConversationsService
	Tasks *taskDispatcher.TaskDispatcher
	Prompts *prompt.PromptService
	Participants *participants.ParticipantService
}

func New(
	repo *repositories.RepositoryStruct, 
	storage interfaces.StorageClient, 
	messenger interfaces.MessageClient,
	txManager interfaces.TxManager,
	taskFlow bool,
) *ServiceStruct {
	return &ServiceStruct{
		Conversations: conversations.NewConversationsService(repo, storage, txManager),
		Tasks: taskDispatcher.NewTaskDispatcher(repo, messenger, storage, txManager, taskFlow),
		Prompts: prompt.NewPromptService(repo, txManager),
		Participants: participants.NewParticipantService(repo, txManager),
	}
}



