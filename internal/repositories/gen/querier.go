// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	ASD(ctx context.Context) (ASDRow, error)
	CreateConversation(ctx context.Context, arg CreateConversationParams) error
	CreateConvert(ctx context.Context, arg CreateConvertParams) error
	CreateParticipant(ctx context.Context, arg CreateParticipantParams) error
	CreatePrompt(ctx context.Context, prompt string) error
	CreateTask(ctx context.Context, taskType int32) (uuid.UUID, error)
	DeleteConversationByID(ctx context.Context, id uuid.UUID) (string, error)
	DeleteConvertByForgeinID(ctx context.Context, conversationsID uuid.UUID) (uuid.UUID, error)
	DeleteConvertByID(ctx context.Context, id uuid.UUID) error
	DeleteParticipantByID(ctx context.Context, id uuid.UUID) error
	DeletePromptByID(ctx context.Context, id uuid.UUID) error
	DeleteTaskByID(ctx context.Context, id uuid.UUID) error
	GetConversationByID(ctx context.Context, id uuid.UUID) (Conversation, error)
	GetConversationFileURL(ctx context.Context, id uuid.UUID) (string, error)
	GetConversations(ctx context.Context) ([]Conversation, error)
	GetConvert(ctx context.Context) ([]Convert, error)
	GetConvertByID(ctx context.Context, id uuid.UUID) (Convert, error)
	GetParticipantByID(ctx context.Context, id uuid.UUID) (Participant, error)
	GetParticipants(ctx context.Context) ([]Participant, error)
	GetPromptByID(ctx context.Context, id uuid.UUID) (Prompt, error)
	GetPrompts(ctx context.Context) ([]Prompt, error)
	GetTaskByID(ctx context.Context, id uuid.UUID) (Task, error)
	GetTasks(ctx context.Context) ([]Task, error)
	UpdateConversationNameByID(ctx context.Context, arg UpdateConversationNameByIDParams) error
	UpdateConversationStatusByID(ctx context.Context, arg UpdateConversationStatusByIDParams) error
	UpdateConvertByTaskID(ctx context.Context, arg UpdateConvertByTaskIDParams) error
	UpdateParticipantByID(ctx context.Context, arg UpdateParticipantByIDParams) error
	UpdatePromptByID(ctx context.Context, arg UpdatePromptByIDParams) error
	UpdateTaskStatus(ctx context.Context, arg UpdateTaskStatusParams) error
}

var _ Querier = (*Queries)(nil)
