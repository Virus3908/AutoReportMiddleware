package repository

import (
	"context"
	"main/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConversationsRepo struct {
	DB *pgxpool.Pool
}

func NewConversationsRepo(db *pgxpool.Pool) *ConversationsRepo {
	return &ConversationsRepo{DB: db}
}

func (r *ConversationsRepo) GetConversationsByID(ctx context.Context, userID uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	query := `SELECT id, conversation_name, file_url, status, created_at, updated_at FROM conversations WHERE conversation_name = $1`
	err := r.DB.QueryRow(ctx, query, userID).Scan(&conversation.ID, &conversation.ConversationName, &conversation.FileURL, &conversation.Status, &conversation.CreatedAt, &conversation.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationsRepo) GetConversations(ctx context.Context) ([]models.Conversation, error) {
	var conversations []models.Conversation
	query := `SELECT id, conversation_name, file_url, status, created_at, updated_at FROM conversations`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var conversation models.Conversation
		err := rows.Scan(&conversation.ID, &conversation.ConversationName, &conversation.FileURL, &conversation.Status, &conversation.CreatedAt, &conversation.UpdatedAt)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}
	return conversations, nil
}

func (r *ConversationsRepo) CreateConversation(ctx context.Context, fileURL, convName string) (*models.Conversation, error) {
	conversation := &models.Conversation{
		ConversationName: convName,
		FileURL: fileURL,
		Status: 0,
	}
	query := `INSERT INTO conversations (conversation_name, file_url, status) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.DB.QueryRow(ctx, query, conversation.ConversationName, conversation.FileURL, conversation.Status).Scan(&conversation.ID, &conversation.CreatedAt, &conversation.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

