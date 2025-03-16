package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"main/internal/models"
)

type ConvertRepo struct {
	DB *pgxpool.Pool
}

func NewConvertRepo(db *pgxpool.Pool) *ConvertRepo {
	return &ConvertRepo{DB: db}
}

func (r *ConvertRepo) GetConvertByID(ctx context.Context, ID uuid.UUID) (*models.Convert, error) {
	var convert models.Convert
	query := `SELECT id, conversation_id, file_url, task_id, created_at, updated_at FROM convert WHERE id = $1`
	err := r.DB.QueryRow(ctx, query, ID).Scan(&convert.ID, &convert.ConversationsID, &convert.FileURL, &convert.TaskID, &convert.CreatedAt, &convert.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &convert, nil
}

func (r *ConvertRepo) GetConvertByTaskID(ctx context.Context, TaskID uuid.UUID) (*models.Convert, error) {
	var convert models.Convert
	query := `SELECT id, conversation_id, file_url, task_id, created_at, updated_at FROM convert WHERE task_id = $1`
	err := r.DB.QueryRow(ctx, query, TaskID).Scan(&convert.ID, &convert.ConversationsID, &convert.FileURL, &convert.TaskID, &convert.CreatedAt, &convert.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &convert, nil
}

func (r *ConvertRepo) GetConvertByConversationID(ctx context.Context, ConversationID uuid.UUID) (*models.Convert, error) {
	var convert models.Convert
	query := `SELECT id, conversation_id, file_url, task_id, created_at, updated_at FROM convert WHERE conversation_id = $1`
	err := r.DB.QueryRow(ctx, query, ConversationID).Scan(&convert.ID, &convert.ConversationsID, &convert.FileURL, &convert.TaskID, &convert.CreatedAt, &convert.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &convert, nil
}

func (r *ConvertRepo) GetConvert(ctx context.Context) ([]models.Convert, error) {
	var converts []models.Convert
	query := `SELECT id, conversation_id, file_url, task_id, created_at, updated_at FROM convert`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var convert models.Convert
		err := rows.Scan(&convert.ID, &convert.ConversationsID, &convert.FileURL, &convert.TaskID, &convert.CreatedAt, &convert.UpdatedAt)
		if err != nil {
			return nil, err
		}
		converts = append(converts, convert)
	}
	return converts, nil
}

func (r *ConvertRepo) CreateConvert(ctx context.Context, fileURL string, conversationID, taskID uuid.UUID) (*models.Convert, error) {
	convert := &models.Convert{
		ConversationsID: conversationID,
		FileURL:         fileURL,
		TaskID:          taskID,
	}
	query := `INSERT INTO convert (conversation_id, file_url, task_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.DB.QueryRow(ctx, query, convert.ConversationsID, convert.FileURL, convert.TaskID).Scan(&convert.ID, &convert.CreatedAt, &convert.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return convert, nil
}

func (r *ConvertRepo) UpdateConvert(ctx context.Context, ID uuid.UUID, fileURL string) (*models.Convert, error) {
	convert := &models.Convert{
		FileURL: fileURL,
	}
	query := `UPDATE convert SET file_url = $1 WHERE id = $2 RETURNING id, conversation_id, file_url, task_id, created_at, updated_at`
	err := r.DB.QueryRow(ctx, query, convert.FileURL, ID).Scan(&convert.ID, &convert.ConversationsID, &convert.FileURL, &convert.TaskID, &convert.CreatedAt, &convert.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return convert, nil
}
